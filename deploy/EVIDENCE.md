# DB-005 (#5) implementation evidence

Recorded 2026-11-03 UTC (2026-11-04 Asia/Bangkok). Scope is `deploy/` only.
Worktree: `/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-005-database`.
Branch: `feat/m01-db-005-database`, based on `origin/develop`
`c8fae1481fe657f385d7747878780d1f4485e986`. No shared files or contracts changed.
No push, PR, merge, approval, protection change, or auto-merge action performed.
The coordinator owns independent review and delivery of the Draft PR to `develop`.

## Dependency and architecture status

Read root README, AGENTS, M01 execution contract and ADR-001 before implementation.
The dispatch confirms DB-001 merged in PR #11 and human ADR approval recorded
there. ADR-001's file still says **Proposed**; that stale text is disclosed here
and intentionally not edited. PostgreSQL remains workstation/server-side only.
No enrollment, application domain tables, remote commands, tunneling, transport,
NAT traversal, production config, or security contract changes were introduced.

## Executed toolchain inspection

Commands and captured output:

```text
$ docker --version
Docker version 29.8.2, build 7fc2dff
$ docker compose version
Docker Compose version v5.5.1
$ docker info --format 'server={{.ServerVersion}} os={{.OSType}} arch={{.Architecture}}'
server=29.8.2 os=linux arch=aarch64
$ uname -srm
Darwin 25.6.0 arm64
$ git --version
git version 2.50.1 (Apple Git-155)
$ openssl version
OpenSSL 3.6.3 9 Jun 2026 (Library: OpenSSL 3.6.3 9 Jun 2026)
$ docker image inspect postgres:17-bookworm --format 'id={{.Id}} digests={{json .RepoDigests}} os={{.Os}} arch={{.Architecture}}'
id=sha256:f92e98998515ee76149a895785ea1fbece6dec5657bf9ab9204037d6f9849fef digests=["postgres@sha256:839fbb5a4942a7758dd93e18b42ea5649807491a0545b5eeeb8f0b8bb11857a5"] os=linux arch=arm64
```

`docker buildx imagetools inspect postgres:17-bookworm` exited 0 and reported
index digest `sha256:839fbb5a4942a7758dd93e18b42ea5649807491a0545b5eeeb8f0b8bb11857a5`,
including `linux/amd64`, `linux/arm/v7`, and `linux/arm64/v8` manifests. Only the
Linux ARM64 image was executed. Container output was
`postgres (PostgreSQL) 17.11 (Debian 17.11-1.pgdg12+2)`.

Initially `docker info` and `docker image ls` failed with:

```text
failed to connect to the docker API at unix:///Users/liiwz/.docker/run/docker.sock; check if the path is correct and if the daemon is running: dial unix /Users/liiwz/.docker/run/docker.sock: connect: no such file or directory
```

`open -a Docker` exited 0; the subsequent daemon inspection succeeded. No Docker
installation, host privilege change, or root shell was needed. Docker Desktop was
left running; test containers, networks and volumes were removed.

## Actual acceptance run

From the worktree root:

```sh
sh -n deploy/migrate.sh && sh -n deploy/smoke.sh && git diff --check
POSTGRES_PASSWORD=$(openssl rand -hex 32) sh deploy/smoke.sh > /tmp/db005-smoke-final.log 2>&1
```

Both commands exited **0**. The password was generated only for the disposable
test environment, was never printed, and was not written to the repository.
Compose configuration was validated with **`config --quiet`**, never rendered with
credentials. [Full captured smoke output](evidence/smoke.txt) includes each command
and actual output; trailing whitespace alone was removed when preserving the log.
`compose`, `migrate`, `sql`, and assertion functions are defined in
[smoke.sh](smoke.sh), with the exact project/directory arguments and SQL commands.

Verified by that executed run:

| Check | Actual result |
| --- | --- |
| Missing and empty password configuration | Both rejected: required variable `POSTGRES_PASSWORD` missing a value |
| Valid `docker compose config --quiet` | Exit 0, no output |
| `up -d --wait --wait-timeout 120` | Healthy PostgreSQL container |
| Published port | `127.0.0.1:55435` (asserted exactly) |
| Runtime user | `uid=999(postgres) gid=999(postgres)`; UID asserted nonzero |
| Limits and isolation inspection | `health=healthy user=postgres memory=268435456 nanoCPUs=1000000000 pids=128 capDrop=[ALL] security=[no-new-privileges:true]` |
| TCP connection with no password | Rejected: `fe_sendauth: no password supplied` |
| Invalid migration action | Rejected with usage output |
| Apply, repeat apply, status | Exactly one baseline tracking row; repeat preserves application timestamp |
| Application schema | Zero tables in `public` |
| Rollback, repeated rollback | Zero tracking rows; both commands succeed |
| Changed rollback file | Rejected: `Applied migration missing or changed; restore original files` |
| Comment moved between paired files without changing concatenated bytes | `cmp` exits 0; migration rejects checksum drift |
| Missing rollback pair | Rejected: `Missing rollback: 0001_baseline.down.sql` |
| Failing forward SQL after DDL | `division by zero`; probe table absent and original tracking count still one |
| Two simultaneous forward runners | Both exit 0, exactly one baseline row |
| Full Compose down / up | Returns healthy and retains baseline history in named volume |
| Final migration rollback / shutdown | Zero tracking rows, no service containers; isolated volume removed |

The first run exited **3** after reaching healthy PostgreSQL: shell `printf`
interpreted `\e` inside a generated `\endif` as an escape character. Captured
[failure excerpt](evidence/initial-failure.txt) preserves that error (`\x1b` is
the escaped representation of the control byte). The runner now emits psql
commands through literal `%s` arguments. The full final run above passed after
that fix; the initial failed run also cleaned up its project and volume.

Additional commands executed:

```text
$ sh -n deploy/migrate.sh && sh -n deploy/smoke.sh && git diff --check
(no output; exit 0)
$ git check-ignore deploy/.env
deploy/.env
$ git check-ignore deploy/.env.example
(no output; exit 1, as expected: example is tracked)
$ docker ps -a --filter name=db005-smoke --format '{{.Names}}'
(no output; exit 0)
$ docker volume ls --filter name=db005-smoke --format '{{.Name}}'
(no output; exit 0)
$ command -v shellcheck
(no output; exit 1: unavailable)
$ command -v gitleaks
(no output; exit 1: unavailable)
$ rg -n --glob '!EVIDENCE.md' '(-----BEGIN ([A-Z]+ )?PRIVATE KEY-----|postgres(ql)?://[^[:space:]]+:[^[:space:]@]+@|^POSTGRES_PASSWORD=[a-zA-Z0-9]{12,})' deploy
(no output; exit 1: no matching credential patterns)
```

## Self-review, security and resource implications

Self-review covered shell quoting, SQL transaction boundaries, concurrent
bootstrap, checksum verification for both directions, filenames, ordering,
credential exposure, volume lifetime, and scope. The live failure above was fixed;
duplicate numeric prefixes and invalid migration names are rejected before SQL.
Coordinator independent review found that hashing concatenated file contents
failed to preserve the boundary between forward and rollback SQL. The fix hashes
each file and filename separately, then hashes those canonical lines. The final
smoke moves a comment across the boundary, proves identical concatenated content
with `cmp`, and verifies that the runner rejects drift. The full smoke was rerun
after this fix, and the preserved final log is from that rerun.
No unresolved implementation blocker was found. Manual inspection and a targeted
credential-pattern search supplement the unavailable dedicated secret scanner;
this does not claim a comprehensive security audit or ShellCheck pass.

The development administrator is a PostgreSQL superuser; no production role model
is invented. Host connections require SCRAM; health verifies TCP password auth.
Container-local Unix sockets are trusted for administration, so Docker access and
the local development environment must be trusted. Environment secrets are
inspectable by Docker operators. No secret is embedded in image configuration,
SQL, committed examples, or logs. A password mismatch against an existing volume
is expected to fail health; initialization does not rotate stored credentials.

All capabilities are dropped, escalation disabled, runtime user is non-root,
SQL/runner mounts are read-only, and the published port is fixed to loopback.
RAM is capped at 256 MiB, CPU at one, processes at 128, shared memory at 64 MiB,
shared buffers at 32 MiB, connections at 20. These constrain workstation resource
use; no throughput, idle RSS, latency, or device memory result is claimed. Named
volume writes stay on the developer's Docker host, never embedded flash. The image
tag is mutable within PostgreSQL major 17; the exact tested digest is above.

No physical ARMv7/AArch64 target, 128 MiB device, Buildroot, Yocto, Windows host,
or Linux AMD64 run was available/performed. The ARM64 Docker VM does not establish
target hardware acceptance. Rust/Go builds are outside this issue. No new runtime
dependencies are required on embedded devices; the migration runner uses tools
already in the PostgreSQL image. The smoke example alone uses host OpenSSL for
disposable credential generation. There is no production deployment claim.

## Rollback and remaining work

Use `docker compose exec -T db sh /opt/devicebridge/migrate.sh down` from `deploy/`
to undo the latest version, then `docker compose down` to stop while preserving
data. Empty migration metadata remains intentionally. To discard local data,
explicitly use `docker compose down --volumes` for the intended local project;
this is destructive and is not a backup procedure. The smoke only removed its
uniquely named disposable project/volume. Reverting this feature commit removes
the infrastructure files but does not delete an existing Docker volume.

Remaining: coordinator independent review, feature push and Draft PR citing #5
with base `develop`, followed by human Tech Lead review. No PR was opened by this
worker. Shared ADR text reconciliation and branch protections remain outside this
worker's ownership; the coordinator reported unprotected branches and disabled
auto-merge, and this worker made no GitHub changes.
