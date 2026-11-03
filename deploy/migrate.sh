#!/bin/sh
# Run inside the db container; requires only tools in the PostgreSQL image.
set -eu
export LC_ALL=C

action=${1:-}
case "$action" in
    up|down|status) ;;
    *) echo 'Usage: sh /opt/devicebridge/migrate.sh {up|down|status}' >&2; exit 2 ;;
esac
[ "$#" -eq 1 ] || { echo 'Expected exactly one action' >&2; exit 2; }

cd /opt/devicebridge/migrations
sql=$(mktemp)
trap 'rm -f "$sql"' EXIT HUP INT TERM

cat > "$sql" <<'SQL'
BEGIN;
-- All invocations share one transaction lock, including bootstrap and rollback.
SELECT pg_advisory_xact_lock(749251, 5);
CREATE SCHEMA IF NOT EXISTS devicebridge_migrations;
CREATE TABLE IF NOT EXISTS devicebridge_migrations.history (
    version text PRIMARY KEY,
    checksum text NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT current_timestamp
);
CREATE TEMP TABLE expected_migrations (version text PRIMARY KEY, checksum text NOT NULL);
SQL

# Only repository-owned SQL is executed. Validate identifiers before embedding them.
previous_prefix=
for forward in *.up.sql; do
    [ -f "$forward" ] || { echo 'No forward migrations found' >&2; exit 1; }
    version=${forward%.up.sql}
    case "$version" in
        [0-9][0-9][0-9][0-9]_[a-z]*) ;;
        *) echo "Invalid migration version: $version" >&2; exit 1 ;;
    esac
    case "$version" in
        *[!a-z0-9_]*) echo "Invalid migration name: $version" >&2; exit 1 ;;
    esac
    prefix=${version%%_*}
    [ "$prefix" != "$previous_prefix" ] || { echo "Duplicate migration prefix: $prefix" >&2; exit 1; }
    previous_prefix=$prefix
    rollback=$version.down.sql
    [ -f "$rollback" ] || { echo "Missing rollback: $rollback" >&2; exit 1; }
    # Hash each direction and its filename before hashing the pair. This preserves
    # file boundaries: moving bytes between directions must count as drift.
    checksum=$(sha256sum "$forward" "$rollback" | sha256sum)
    checksum=${checksum%% *}
    printf "INSERT INTO expected_migrations VALUES ('%s', '%s');\n" "$version" "$checksum" >> "$sql"
done
for rollback in *.down.sql; do
    [ -f "${rollback%.down.sql}.up.sql" ] || { echo "Missing forward pair: $rollback" >&2; exit 1; }
done

cat >> "$sql" <<'SQL'
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM devicebridge_migrations.history h
        LEFT JOIN expected_migrations e USING (version)
        WHERE e.version IS NULL OR e.checksum <> h.checksum
    ) THEN
        RAISE EXCEPTION 'Applied migration missing or changed; restore original files';
    END IF;
    IF EXISTS (
        SELECT 1 FROM expected_migrations e
        WHERE e.version < (SELECT max(version) FROM devicebridge_migrations.history)
          AND NOT EXISTS (SELECT 1 FROM devicebridge_migrations.history h WHERE h.version = e.version)
    ) THEN
        RAISE EXCEPTION 'Migration history has a gap; append new migrations in order';
    END IF;
END $$;
SELECT coalesce(max(version), '') AS rollback_version FROM devicebridge_migrations.history
\gset
SQL

for forward in *.up.sql; do
    version=${forward%.up.sql}
    case "$action" in
        up)
            printf "SELECT NOT EXISTS (SELECT 1 FROM devicebridge_migrations.history WHERE version = '%s') AS run_migration\n" "$version" >> "$sql"
            printf '%s\n' '\gset' '\if :run_migration' >> "$sql"
            cat "$forward" >> "$sql"
            printf "\nINSERT INTO devicebridge_migrations.history SELECT version, checksum, current_timestamp FROM expected_migrations WHERE version = '%s';\n" "$version" >> "$sql"
            printf '%s\n' '\endif' >> "$sql"
            ;;
        down)
            printf "SELECT :'rollback_version' = '%s' AS run_migration\n" "$version" >> "$sql"
            printf '%s\n' '\gset' '\if :run_migration' >> "$sql"
            cat "$version.down.sql" >> "$sql"
            printf "\nDELETE FROM devicebridge_migrations.history WHERE version = '%s';\n" "$version" >> "$sql"
            printf '%s\n' '\endif' >> "$sql"
            ;;
    esac
done

cat >> "$sql" <<'SQL'
TABLE devicebridge_migrations.history ORDER BY version;
COMMIT;
SQL
psql -X --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" --set ON_ERROR_STOP=1 --file "$sql"
