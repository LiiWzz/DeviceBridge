# M01 integration smoke — DB-008 / #8

From the repository root:

```sh
bash tests/smoke.sh
```

Requires a Unix build host with Go 1.23.4, Rust/Cargo 1.89.0, Python 3,
Docker and Compose supporting `up --wait`. CI runs on Ubuntu 24.04. No Python,
Docker or compiler is installed on embedded devices. Build binaries from this
checkout, run the harness failure gates, then check the actual processes:

- Control: GET/HEAD health/readiness, invalid methods (405), absent debug/device
  endpoints (404), reject public bind configuration; after SIGTERM and SIGINT,
  require exit 0, a stopped log and closed listener.
- Relay/agent: wait for startup logs before signalling; require clean exit and
  stopped logs for both signals. Invalid transport/config fails closed.
- CLI: local config validation succeeds; insecure URL and unimplemented connect
  fail. These checks do not attempt remote connectivity.
- Compose: reuse `deploy/smoke.sh` for SQL health, missing/empty password,
  unauthenticated TCP rejection, migration idempotence/drift/transaction failure,
  concurrent migration serialization, persistence after restart and shutdown.
  Only its disposable project and volume are removed by its cleanup trap.

Stages are sequential because M01 control readiness has no database dependency,
and relay/agent have no transport. This is component lifecycle integration, not
proof of a connected system, authentication, SSH or M01 acceptance. Readiness may
become connection refusal during shutdown; no observable 503 window is required.
Startup/shutdown waits are bounded. Unexpected exits, HTTP results or migration
results fail the job. Harness failure gates check wrong exit codes, premature
exit, unclean shutdown and cleanup on probe failure. Ephemeral control port
selection has a small bind race; a conflict fails visibly.

Output is in ignored `dist/ci/integration/`: command log, process logs, source/run
identity and native binaries. `PASS` is removed before each run and recreated
only after all stages pass. CI uploads available evidence even on failure, with
14-day retention. Exact executed results and limits: [EVIDENCE.md](EVIDENCE.md).

Operator control configuration is cleared from child environments. HTTP probes
use a direct loopback connection without proxies. A fresh disposable database
password is generated without printing it; `deploy/.env` is not copied. Never run
with shell tracing or collect rendered Compose configuration/container
environments. `POSTGRES_PORT` optionally selects the DB test port (default 55435).
Concurrent smokes need distinct ports. Kill/host loss may bypass cleanup; inspect
the recorded `db005-smoke-*` project name before manually removing only that
disposable project's resources. Ordinary developer volumes are not targets.

Changes affect build-host CPU/storage/network and CI duration only. No runtime
dependencies, device allocation/flash writes, root requirement or application/API/
security contract changes. Docker access retains its existing privileged-host
trust boundary. Real ARM, 128 MiB RSS/CPU and Buildroot/Yocto validation remain
unavailable; native macOS and Linux CI results are not physical ARM acceptance.
Rollback through a reviewed revert PR into `develop`; remove ignored outputs if
desired. No database schema change, device change or production rollout occurs.
