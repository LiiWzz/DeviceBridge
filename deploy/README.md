# Local PostgreSQL foundation (DB-005 / #5)

This is a developer workstation database, with PostgreSQL 17 and Docker Compose.
It does not deploy anything to embedded devices or production. Docker and a
Compose release supporting `up --wait` are required; no host PostgreSQL client,
language runtime, or root shell is required. Docker access itself is privileged
host access; use an existing developer Docker installation.

The repeatable acceptance smoke uses a temporary copy and its own project/volume.
It deletes only that disposable volume on exit. From the repository root, run
`POSTGRES_PASSWORD=$(openssl rand -hex 32) sh deploy/smoke.sh` (OpenSSL is only
needed for this example's disposable password generation). The smoke uses host
port 55435 by default; override `POSTGRES_PORT` if necessary. It does not read
`deploy/.env` or touch the normal development volume.

Run these commands from `deploy/`:

```sh
cp .env.example .env
chmod 600 .env
# Edit .env and set POSTGRES_PASSWORD to your own unique password.
# Use a single-quoted dotenv value if it contains $ or #; do not commit .env.
docker compose config --quiet
docker compose up -d --wait --wait-timeout 120
docker compose ps
docker compose exec -T db sh /opt/devicebridge/migrate.sh up
docker compose exec -T db sh /opt/devicebridge/migrate.sh status
```

An empty or missing password causes configuration validation to fail. Never use
`docker compose config` without `--quiet` when collecting evidence: rendered
configuration includes the password. Container inspection and Docker access can
also expose development environment secrets. No real or default password is
provided. Set `POSTGRES_PORT` to another free port if 5432 is already in use;
the binding is always `127.0.0.1`, never a public interface. The database/user
are both `devicebridge`; this is a development administrator, not a production
application role. Host connections use SCRAM authentication. The SQL healthcheck
checks the configured password over TCP, rather than only checking a listening
socket. Local container Unix-socket access is trusted for migration administration.

## Migrations

Migration commands use the running database container and the image's POSIX shell,
`sha256sum`, and `psql`. SQL and the runner are mounted read-only. There is no
automatic startup migration, so schema changes remain explicit.

```sh
# Apply all pending versions in filename order; repeat is a no-op.
docker compose exec -T db sh /opt/devicebridge/migrate.sh up
# Roll back the latest applied version; repeat with no versions is a no-op.
docker compose exec -T db sh /opt/devicebridge/migrate.sh down
# Show versions, checksums, and application times.
docker compose exec -T db sh /opt/devicebridge/migrate.sh status
```

`0001_baseline` only annotates migration metadata. No enrollment, device, session,
or other application tables are defined. Tracking lives in
`devicebridge_migrations.history`; rolling the baseline back leaves this empty
metadata table/schema in place. Forward and rollback files together have a tracked
SHA-256 checksum of their individual hashes and filenames, preserving file
boundaries. Changed/missing applied files and gaps in history fail closed.
Concurrent commands serialize on a PostgreSQL advisory transaction lock. Every
command, including all pending versions, is one transaction with `ON_ERROR_STOP`;
a SQL failure rolls back both the changes and their tracking rows.

Future migrations must be reviewed pairs named `NNNN_lowercase_name.up.sql` and
`NNNN_lowercase_name.down.sql`, using unique increasing four-digit prefixes. Append
versions; never modify an applied pair. SQL must be transactional, with no
`BEGIN`, `COMMIT`, `ROLLBACK`, psql commands, or operations such as
`CREATE INDEX CONCURRENTLY` that cannot run in a transaction. The runner executes
trusted repository SQL, not user-supplied SQL, and is a foundation rather than a
general-purpose production migration tool.

## Persistence, resources, and rollback

```sh
# Stop and remove containers/network; retain the named pgdata volume.
docker compose down
# Restart later, preserving data and migration history.
docker compose up -d --wait --wait-timeout 120
```

The default Compose project name is `devicebridge-dev`; use a distinct
`docker compose -p NAME ...` consistently for independent environments. Compose
stores the named `pgdata` volume outside the repository. Initialization settings
only apply to a new volume. Changing `.env` does not change a stored user's
password; a mismatch intentionally makes the healthcheck fail. Restore the old
value, or rotate the database password deliberately before updating `.env`.

To undo a migration, run `down` on the migration runner before stopping the stack.
To discard **all data for this local project**, explicitly run
`docker compose down --volumes`; this is irreversible without a separate backup.
No command here is a backup or production upgrade procedure.

The database runs as the image's `postgres` user with all Linux capabilities
dropped and no privilege escalation. Limits are 256 MiB container RAM, one CPU,
128 processes, 64 MiB shared memory, 32 MiB PostgreSQL shared buffers, and 20
connections. These are local development bounds, not performance guarantees.
The Docker VM adds host overhead. The PostgreSQL 17 Debian image supports Linux
ARM64 and AMD64; actual host/image validation is recorded in [EVIDENCE.md](EVIDENCE.md).
The major-version image tag receives upstream fixes and is mutable; a production
image/digest policy is out of scope. Do not change the major version against an
existing volume without a planned PostgreSQL upgrade.

ADR-001 still says “Proposed.” The dispatch records human approval in merged
foundation PR #11; this implementation relies on that approval and leaves the
shared ADR untouched. Security contracts remain unchanged. No ARMv7 embedded
hardware, 128 MiB device, fleet behavior, or remote connectivity is validated here.
