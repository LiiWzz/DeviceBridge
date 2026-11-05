# DB-007 executed evidence

Date: 2026-11-05 (Asia/Bangkok). Base: `origin/develop` at
`74cda23`. Host: macOS ARM64; Go 1.23.4, Rust 1.89.0. Tests ran on the
working tree; [tested-source-sha256.txt](evidence/tested-source-sha256.txt)
identifies tested implementation/configuration files. The Draft PR records its
exact final commit SHA. Command output is captured in the linked files; trailing whitespace is normalized.

| Command | Actual result / output |
| --- | --- |
| `go test ./cmd/... -run TestMalformedArgumentsDoNotLeakValues -count=1 -v` in `cli` and `server/relay`, before fix | Both exit 1; sensitive-input disclosure reproduced: [CLI](evidence/before-cli.txt), [relay](evidence/before-server-relay.txt) |
| `bash ci/check.sh go server/control` | Exit 0: formatting, vet, verbose tests, race tests, native build and hashes: [output](evidence/server-control.txt) |
| `bash ci/check.sh go server/relay` | Exit 0, including new diagnostic regression tests: [output](evidence/server-relay.txt) |
| `bash ci/check.sh go cli` | Exit 0, including new diagnostic regression tests: [output](evidence/cli.txt) |
| `bash ci/check.sh rust` | Exit 0: fmt, Clippy warnings denied, locked tests and release build: [output](evidence/rust.txt) |
| `bash -n ci/check.sh ci/security.sh`; `python3 ci/test_check.py` | Exit 0; four CI failure-gate tests pass: [output](evidence/tooling.txt) |
| `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 -shellcheck= .github/workflows/ci.yml` | Exit 0: [output](evidence/tooling.txt) |
| `docker run --rm -v "$PWD:/mnt" koalaman/shellcheck:v0.10.0 ci/check.sh ci/security.sh` | Exit 0; host ShellCheck is unavailable, so Linux container used: [output](evidence/tooling.txt) |
| `POSTGRES_PASSWORD=<random disposable value> POSTGRES_PORT=55437 sh deploy/smoke.sh` | Exit 0 on Docker Linux/aarch64; negative missing/empty/wrong password checks, loopback binding, nonroot/capability/resource inspection, migration drift and transactional failure checks, restart/shutdown and isolated cleanup: [output](evidence/compose.txt) |
| `PATH="$PWD/dist/security-tools:$PATH" bash ci/security.sh` | See [scan output](evidence/security.txt) and [redacted report](evidence/gitleaks.json); installed with `GOBIN="$PWD/dist/security-tools" go install github.com/zricethezav/gitleaks/v8@v8.24.3` |

No suppression/allowlist was added. The first synthetic fixture was too short
for the built-in private-key rule and the gate correctly failed; the final
fixture uses a longer explicitly nonfunctional marker. No real key was generated
or stored. The scan covers committed reachable local history, not uncommitted or
ignored files. Hosted checks and cross-build outcomes, if available, are recorded
on the PR; local execution alone does not establish their status.

No new ARM cross-linker validation was performed locally for this issue. CI
retains ARMv7/AArch64 jobs; no physical ARMv7/AArch64 embedded board with 128 MiB
RAM is attached. macOS ARM64 and a Linux Docker VM do not satisfy device RAM,
idle RSS, CPU, Buildroot/Yocto/libc compatibility or reconnect acceptance. There
is no implemented remote connectivity/authentication to validate. Packaging code
was unchanged and its local suite was not rerun; hosted CI still includes it.

Security impact: fixed Go flag diagnostics that could echo sensitive input;
added threat/config/secret inventory, an unimplemented lifecycle design checklist,
and a redacted history scan with detection self-check. Existing parsing acceptance,
loopback bindings and disabled transport are preserved. No credentials or
production configuration introduced. See [baseline](README.md) for residual
risks and approval gates.

Performance: scanner and its Go dependencies run on the build host/CI only.
No agent runtime dependencies, idle allocations, flash writes, root requirement,
transport buffers or database schema changes. Error diagnostics use fixed text.

Rollback: reviewed revert through a feature branch/Draft PR into `develop`;
no data rollback or secret rotation is needed for this change. Reverting loses
the scan and diagnostic protection. No merge, self-approval or auto-merge is
performed. Independent review and human Tech Lead approval are pending.
