# M01 security baseline — DB-007 / #7

M01 is a local foundation preview. There is no enrollment, user/device
credential store, authentication service, session authorization, TLS transport,
or remote SSH. A successful config check or health response establishes none of
these. No demo login, shared fleet token, shell RPC or port-forward is provided.

This baseline documents existing behavior under ADR-001. It does not approve
production deployment or change configuration/wire contracts. ADR-001 remains
marked Proposed; foundation components already landed in `develop`. Any M02
security contract requires a separate ADR and explicit human Tech Lead approval
before implementation. Proposed lifecycle steps below are review inputs only.

## Assets and trust boundaries

| Boundary / asset | Attacker or failure | Existing control and evidence | Residual risk / gate |
| --- | --- | --- | --- |
| Operator files/arguments → agent process | Malformed JSON, huge input, FIFO blocking, diagnostic injection | 4096-byte regular-file read; duplicate/unknown field and numeric checks; sanitized errors; `agent/src/config.rs`, `agent/tests/process.rs` | Config owner must be trusted; symlinks to regular files allowed; OS argument/environment limits still apply |
| Local HTTP client → control | Public exposure, request exhaustion, debug routes mistaken for auth | Loopback IP literal only; ports 1024–65535; bounded timeouts and headers; exact health routes; `server/control/internal/control/*_test.go` | Local users can call probes; plain HTTP; no aggregate connection/rate cap; never expose through a proxy as a production API |
| CLI settings → future relay | Credentials in URL, insecure endpoint, unauthorized remote operations | HTTPS URL validation rejects userinfo/query/fragment; `check` makes no connection; remote commands unavailable; `cli/*` tests | URL validation is not certificate verification or host-key pinning; both need M02 design |
| Relay settings → transport | Resource exhaustion or accidental listener | Only `disabled`; stream/buffer ranges validated; unimplemented transport errors; relay tests | Ceilings are reserved settings, not allocated/enforced transport resources; maxima multiply to 1 GiB if naively allocated; review aggregate bounds before M02 |
| Local host → development PostgreSQL | Missing password, external exposure, privilege escalation | Required nonempty password; IPv4 loopback binding; TCP SCRAM; postgres user, capabilities dropped, no-new-privileges and resource limits; `deploy/smoke.sh` | Docker users can read secrets and administer DB; local Unix socket trusted; development role is an administrator; never use as production role |
| Repository/CI → published artifacts | Committed credentials, leaked build logs, malicious dependencies | Ignored env/key files; redacted full-history Gitleaks; lockfile/pinned tool versions; CI read-only token; explicit package inputs | Pattern scans miss unknown/encoded secrets; ignored files can be force-added; mutable action/image tags and dependency provenance remain review concerns |
| Compromised device or workstation → fleet | Stolen future identity, impersonation, replay | No operational device credentials or remote capability in M01 | Host compromise is outside M01 protection; M02 needs revocation, per-session authorization and SSH host-key verification |

No successful remote connectivity or production security guarantees are claimed.
DB-002/003/004/005 merge commits are present in base `74cda23`; their issue
trackers still show open. DB-008 smoke and DB-009 independent review remain
separate work, not completed by this baseline.

## Configuration inventory and validation

The code is authoritative. [agent-config.schema.json](agent-config.schema.json)
is an editor/documentation schema; no schema runtime is installed on ARM. JSON
Schema cannot detect duplicate keys or cap serialized file bytes. The Rust parser
additionally rejects duplicate fields, invalid UTF-8, trailing documents and
non-object roots. No fields are required; `{}` uses defaults. JSON Schema's
integer type can accept `10.0` while the Rust `u64` parser rejects it.

| Component/input | Default | Accepted values / precedence |
| --- | --- | --- |
| Agent `--config PATH` JSON | No file; defaults | Regular file ≤4096 bytes, read once; only `log_level`, `shutdown_poll_ms`; invalid files fail before CLI overrides |
| Agent `log_level` / `--log-level` | `info` | `error`, `info`, `debug`; CLI overrides valid file; no environment settings |
| Agent `shutdown_poll_ms` / `--shutdown-poll-ms` | `100` | Integer 10–1000; CLI unsigned decimal digits; repeated flags rejected; values ≤1024 bytes; help/version skip file loading |
| Control `DEVICEBRIDGE_CONTROL_ADDR` | `127.0.0.1:8080` | Loopback IP literal, no zone/DNS/wildcard; canonical decimal port 1024–65535; IPv4-mapped loopback accepted |
| Control prefix + `READ_HEADER_TIMEOUT`, `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT`, `SHUTDOWN_TIMEOUT` | `5s`, `10s`, `10s`, `60s`, `10s` | Go durations 1ms–5m; header ≤ read timeout; present-but-empty fails; unrelated environment variables ignored |
| Relay `--transport` | `disabled` | Only `disabled`; no network operation |
| Relay `--max-streams` / `--max-buffer-bytes` | `16` / `32768` | Integers 1–1024 / 1–1048576; no buffers allocated in M01 |
| CLI `check --relay-url` | `https://localhost:443` | Absolute HTTPS; valid ASCII DNS or IP, no userinfo/query/fragment/non-root path, no IPv6 zone; optional port 1–65535 |
| CLI `check --timeout` | `5s` | Go duration >0 and ≤1m; no request performed |
| Go relay/CLI flags | Above | Unknown/missing/malformed flags fail with exit 2; repeated flags use Go flag semantics (last wins); parser diagnostics do not echo arguments |
| Compose `POSTGRES_PASSWORD` / `POSTGRES_PORT` | None / `5432` | Password must be nonempty (not a password-strength validator); Compose validates port; bind remains `127.0.0.1` |

Negative cases are executable in the existing component tests: invalid types,
unknown/duplicate JSON keys, overflow, file size boundaries, FIFOs, insecure URLs,
URL credentials, non-loopback addresses, invalid timeouts and disabled transport.
DB-007 adds regression tests for accidental disclosure through Go parser errors.
Help still prints known flag names/defaults; it never prints supplied values.

## Secret handling and least privilege

- M01 agent/relay/CLI accept no device token or private key. Do not add credentials
  to CLI arguments, URLs, examples, test evidence or logs. OS process inspection
  can expose arguments even when application diagnostics are sanitized.
- Copy `deploy/.env.example` locally, use `chmod 600`, and supply a unique local
  password. The checked-in password is intentionally empty. Do not collect
  rendered `docker compose config`, container environments or shell tracing;
  `docker compose config --quiet` validates without printing the password.
- Environment variables and `.env` are development conveniences, not a production
  secret store. Docker access implies privileged host/database access. Rotation
  must update the stored database password deliberately; changing `.env` alone
  leaves the old volume credential and causes healthcheck failure.
- Run agent/control/relay as ordinary users without capabilities; use private
  operator-owned directories. No device Docker, database, root, state directory
  or flash writes are required. Collect logs using bounded operator-managed
  retention; a blocked sink can delay shutdown.
- Go parser errors use fixed text; agent JSON errors use location only; control
  errors name the setting without its supplied value. Do not log request headers,
  credentials or arbitrary payloads when adding future functionality.
- Scan before publishing; if a real secret is found, revoke/rotate it first,
  remove it from the proposed change and coordinate history cleanup with the
  Tech Lead. Removing a file or reverting a commit does not revoke a secret.

## Proposed credential and token lifecycle — unimplemented

The following is a design checklist for a future ADR, not a new accepted policy,
API, schema, token format or TTL. No automatic enrollment or rotation ships here.

| Stage | Proposed behavior to review | Failure / acceptance question |
| --- | --- | --- |
| Bootstrap | Operator-authorized, device-bound, one-use enrollment capability delivered outside logs/arguments; no default fleet credential | Reject replay, expired enrollment and identity mismatch; decide secure bootstrap and recovery channels |
| Issue/store identity | Unique per-device identity with least privilege; keep private material on device in protected storage; server stores only necessary verifiers/public material | Decide certificate vs opaque token, protected storage, permissions and clone resistance across Buildroot/Yocto/Debian |
| Authorize session | Authenticate user and device; verify device scope and fresh authorization on each session; client verifies SSH host key; relay never stores SSH private keys | Reject wrong-device, wrong-user, expired/revoked grants and host-key mismatch; no fallback auth or shell protocol |
| Expire/renew | Bounded lifetimes; renew authenticated identities before expiry; fail closed on expiry, renewal failure or uncertain authorization | Decide TTLs, trusted time/skew, offline behavior, retry bounds and aggregate resource limits |
| Rotate | Reissue per-device identity with explicitly bounded overlap; update atomically and verify before retiring old identity | Decide crash/power-loss recovery and flash write budget; old identity cannot become permanently valid |
| Revoke/respond | Revoke a single device/user/grant without rotating the fleet; reject further sessions and define handling of active sessions | Decide revocation propagation SLA and control-plane outage semantics; test disconnect/reconnect and stale cache cases |
| Decommission | Disable server identity, revoke pending grants, remove local secrets through an explicit operator flow | Decide auditable cleanup, backup retention and recovery; never reinstate revoked identity during rollback |

Production blockers include the approved ADR, actual authentication and expiry
checks, revocation behavior, authenticated transport, SSH host-key verification,
aggregate resource enforcement and physical ARM validation. No implementation
should infer approval from this checklist.

## Reproduction, impact and rollback

Install Gitleaks on the build host only:

```sh
go install github.com/zricethezav/gitleaks/v8@v8.24.3
bash ci/security.sh
bash ci/check.sh rust
bash ci/check.sh go server/control
bash ci/check.sh go server/relay
bash ci/check.sh go cli
```

Put the Go bin directory on `PATH`. `ci/security.sh` first proves detection and
redaction against a synthetic private-key-shaped fixture in a temporary directory,
then scans committed history reachable from all local refs. CI fetches full
history. Local shallow clones cover less history. Uncommitted/ignored files are
not covered; commit and rerun before pushing. Built-in rules are used without a
project allowlist. Pattern scans cannot certify absence of all secrets.

Exact executed commands, captured outcomes, limitations and source identity live
in [EVIDENCE.md](EVIDENCE.md). No independent review is claimed. The new scanner
adds build-host CPU/network/disk cost only; no runtime dependency, agent allocation,
flash write or privilege is added. Parser diagnostics use constant messages;
configuration acceptance and transport remain unchanged. The agent's ≤16 MiB
idle RSS target and 128 MiB ARM device acceptance need real hardware evidence.

Rollback is a normal reviewed revert of this issue's commit(s), via a feature
branch and PR to `develop`. No data migration, credential change, persistent
state or volume deletion is involved. Reverting reintroduces the diagnostic
exposure and removes the added CI check; never restore a revoked credential.
