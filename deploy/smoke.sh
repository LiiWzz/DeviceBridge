#!/bin/sh
# Isolated, destructive ONLY to its own temporary Compose project and volume.
set -eu
: "${POSTGRES_PASSWORD:?Provide a unique disposable test password}"
export POSTGRES_PASSWORD
export POSTGRES_PORT="${POSTGRES_PORT:-55435}"
source_dir=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
test_dir=$(mktemp -d)
project=db005-smoke-$$
cp "$source_dir/compose.yaml" "$source_dir/migrate.sh" "$test_dir/"
cp -R "$source_dir/migrations" "$test_dir/"

compose() { docker compose --project-directory "$test_dir" -p "$project" "$@"; }
cleanup() {
    result=$?
    trap - EXIT HUP INT TERM
    compose down --volumes || result=1
    rm -rf "$test_dir"
    exit "$result"
}
trap cleanup EXIT HUP INT TERM
run() { printf '\n+ '; printf '%s ' "$@"; printf '\n'; "$@"; }
reject() {
    printf '\n+ expect failure: '; printf '%s ' "$@"; printf '\n'
    if "$@"; then echo 'FAIL: unexpected success' >&2; exit 1; else echo 'Expected failure observed'; fi
}
migrate() { compose exec -T db sh /opt/devicebridge/migrate.sh "$@"; }
sql() { compose exec -T db psql -X -U devicebridge -d devicebridge -v ON_ERROR_STOP=1 -Atqc "$1"; }
assert_sql() {
    actual=$(sql "$1")
    [ "$actual" = "$2" ] || { echo "FAIL SQL: $1 returned $actual, expected $2" >&2; exit 1; }
    printf 'PASS SQL: %s => %s\n' "$1" "$actual"
}

reject env -u POSTGRES_PASSWORD docker compose --project-directory "$test_dir" -p "$project" config --quiet
# Test empty separately; no .env is copied into the temporary directory.
reject env POSTGRES_PASSWORD= docker compose --project-directory "$test_dir" -p "$project" config --quiet
run compose config --quiet
run compose up -d --wait --wait-timeout 120
run compose ps
run compose exec -T db id
run compose exec -T db postgres --version
run compose port db 5432
[ "$(compose port db 5432)" = "127.0.0.1:$POSTGRES_PORT" ]
container=$(compose ps -q db)
run docker inspect --format 'health={{.State.Health.Status}} user={{.Config.User}} memory={{.HostConfig.Memory}} nanoCPUs={{.HostConfig.NanoCpus}} pids={{.HostConfig.PidsLimit}} capDrop={{.HostConfig.CapDrop}} security={{.HostConfig.SecurityOpt}}' "$container"
[ "$(compose exec -T db id -u)" != 0 ]
reject compose exec -T -e PGPASSWORD= db psql -X -h 127.0.0.1 -U devicebridge -d devicebridge -w -c 'SELECT 1'
reject migrate invalid
run migrate up
run migrate up
run migrate status
assert_sql 'SELECT count(*) FROM devicebridge_migrations.history' 1
assert_sql "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public'" 0
run migrate down
run migrate down
assert_sql 'SELECT count(*) FROM devicebridge_migrations.history' 0
run migrate up

# Reject drift of either direction, and missing rollback pairs.
printf '\n-- deliberate test drift\n' >> "$test_dir/migrations/0001_baseline.down.sql"
reject migrate up
cp "$source_dir/migrations/0001_baseline.down.sql" "$test_dir/migrations/"
# Moving a trailing comment across the pair boundary preserves concatenated bytes
# but must still change the tracked checksum for the individual files.
forward=$test_dir/migrations/0001_baseline.up.sql
rollback=$test_dir/migrations/0001_baseline.down.sql
cat "$forward" "$rollback" > "$test_dir/pair.before"
sed '$d' "$forward" > "$test_dir/forward.moved"
tail -n 1 "$forward" > "$test_dir/rollback.moved"
cat "$rollback" >> "$test_dir/rollback.moved"
mv "$test_dir/forward.moved" "$forward"
mv "$test_dir/rollback.moved" "$rollback"
cat "$forward" "$rollback" > "$test_dir/pair.after"
run cmp "$test_dir/pair.before" "$test_dir/pair.after"
reject migrate up
cp "$source_dir/migrations/0001_baseline.up.sql" "$forward"
cp "$source_dir/migrations/0001_baseline.down.sql" "$rollback"
mv "$test_dir/migrations/0001_baseline.down.sql" "$test_dir/rollback.saved"
reject migrate up
mv "$test_dir/rollback.saved" "$test_dir/migrations/0001_baseline.down.sql"

# A failing pending migration must leave neither its DDL nor a tracking row.
printf 'CREATE TABLE public.db005_failure_probe (id integer);\nSELECT 1 / 0;\n' > "$test_dir/migrations/0002_failure.up.sql"
printf 'DROP TABLE public.db005_failure_probe;\n' > "$test_dir/migrations/0002_failure.down.sql"
reject migrate up
assert_sql "SELECT to_regclass('public.db005_failure_probe') IS NULL" t
assert_sql 'SELECT count(*) FROM devicebridge_migrations.history' 1
rm "$test_dir/migrations/0002_failure.up.sql" "$test_dir/migrations/0002_failure.down.sql"

run migrate down
printf '\n+ migrate up (two concurrent processes)\n'
migrate up > "$test_dir/concurrent-a.log" 2>&1 &
first=$!
migrate up > "$test_dir/concurrent-b.log" 2>&1 &
second=$!
wait "$first"
wait "$second"
cat "$test_dir/concurrent-a.log" "$test_dir/concurrent-b.log"
assert_sql 'SELECT count(*) FROM devicebridge_migrations.history' 1
run compose down
run compose up -d --wait --wait-timeout 120
assert_sql 'SELECT count(*) FROM devicebridge_migrations.history' 1
run migrate down
assert_sql 'SELECT count(*) FROM devicebridge_migrations.history' 0
run compose down
[ -z "$(compose ps -q)" ]
echo 'PASS: DB-005 smoke; cleanup will remove only the isolated test volume'
