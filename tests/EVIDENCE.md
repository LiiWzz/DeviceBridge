# DB-008 executed evidence

Issue: [#8](https://github.com/LiiWzz/DeviceBridge/issues/8). Base:
`af49bad5270988d533ebb9e50f925274c13d67b7` (full base in `evidence/local-SOURCE.txt`), latest `origin/develop`
when this feature branch was created. DB-003/004/005/006/007 implementations are
merged into this base, although their issue trackers remain open. No contracts,
security policy or ADR status change is made here.

Before: CI had native/unit/cross-build/security checks; no Linux integration
job. DB-005 had a manual database smoke. After: one CI job builds native binaries,
runs process probes and signal shutdown checks, then reuses the isolated database
smoke. Independent review and human acceptance remain separate.

## Local execution

Host: macOS ARM64, Go 1.23.4, Cargo/Rust 1.89.0, Docker server 29.8.2,
Compose v5.5.1; PostgreSQL 17.11 Debian container in the Linux Docker VM.

| Executed command | Actual result |
| --- | --- |
| `bash tests/smoke.sh` | Exit 0; four harness failure gates, three Go builds, locked Rust release build, eight command exit checks, both signal cycles, Compose migration/negative/restart/shutdown and isolated volume cleanup |
| `bash -n ci/check.sh ci/security.sh tests/smoke.sh` | Exit 0 |
| `docker run --rm --network none -v "$PWD:/src:ro" -w /src koalaman/shellcheck:v0.10.0 ci/check.sh ci/security.sh tests/smoke.sh` | Exit 0 after correcting initial EXIT-trap variable warning |
| `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 -shellcheck= .github/workflows/ci.yml` | Exit 0 |
| `python3 ci/test_check.py` | Exit 0; four existing CI failure-gate tests passed |
| `git diff --check` | Exit 0 |

Captured smoke output (trailing whitespace trimmed) and tool versions: [local-smoke.txt](evidence/local-smoke.txt).
Exact lint commands/output/exit codes: [tooling.txt](evidence/tooling.txt).
Local run was against uncommitted source on the stated base; tested script and
workflow hashes: [tested-source-sha256.txt](evidence/tested-source-sha256.txt).
Application code is unchanged. Existing component unit/race/lint and ARM builds
are left to their unchanged CI jobs; no new local execution of those suites or
ARM cross-linkers is claimed.

## Security, performance, limitations and rollback

Negative cases cover public bind, insecure CLI URL, unimplemented connect and
transport, invalid agent config, HTTP route/method guards, missing/empty password,
unauthenticated DB TCP and invalid/drifted/failing migrations. Disposable database
password is generated and never printed. No credentials, device keys or production
configuration are added. Gitleaks v8.24.3 `dir tests --no-banner --redact=100` and
`PATH="$(go env GOPATH)/bin:$PATH" bash ci/security.sh` exited 0. The latter
proved synthetic detection/redaction, then scanned 20 reachable local commits
without findings; output: [security.txt](evidence/security.txt).
pattern scanning cannot certify absence of every secret.

Build-host CPU/disk/network and one bounded CI job are added; exact resource
consumption is not benchmarked. No device runtime/dependency/allocation/flash/root
change. Process probes are local; the control service does not depend on the DB,
and there is no functional SSH/tunnel/authentication test or production guarantee.
No physical ARM board, 128 MiB RAM/RSS/CPU acceptance, Buildroot/Yocto, Windows or
fleet test is available. Docker VM ARM64 execution is not physical device evidence.

Rollback: reviewed revert PR into `develop`, removing the integration job and
harness. No schema or device changes; smoke cleanup removes only its disposable
volume. If killed before cleanup, identify the exact recorded test project before
manual cleanup; never delete the ordinary development database volume. Human
Tech Lead alone approves/manual-merges the Draft PR.
