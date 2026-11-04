# DB-002 implementation evidence

Issue: https://github.com/LiiWzz/DeviceBridge/issues/2
Recorded 2026-11-04 Asia/Bangkok (2026-11-03 UTC).

Worktree: `/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-002-rust-agent`
Branch: `feat/m01-db-002-rust-agent`
Base: `origin/develop` at `c8fae1481fe657f385d7747878780d1f4485e986`.
All changes are owned by `agent/`. The exact delivery commit SHA is reported to
the coordinator after committing; resolve it locally with `git rev-parse HEAD`.
No push or PR creation was performed; coordinator review and Draft PR into
`develop` remain separate delivery steps.

Read README.md, AGENTS.md, docs/M01_EXECUTION.md and ADR-001 before coding.
DB-001 dependency and human approval in merged PR #11 were supplied by the
coordinator. ADR-001's checked-in status still says **Proposed**; this discrepancy
is disclosed, and the shared ADR is unchanged. The coordinator also reported
main/develop currently unprotected and auto-merge disabled; no repository
settings were changed.

## Result and acceptance limits

Implemented the Rust binary, optional bounded JSON configuration, CLI overrides,
JSON event logging, robust exit codes, Unix SIGINT/SIGTERM shutdown, release size
optimizations and local ARM GNU cross-linker configuration. No networking,
enrollment, SSH/tunneling, NAT traversal, remote commands or authentication exists.

Host checks pass. Both requested Linux ARM builds were actually attempted and
**failed** due to missing GNU cross-linkers. Linux ARM executable builds and real
target runtime/resource validation remain outstanding; this is not a fully
validated ARM release. The macOS host uses an ARM64 CPU but its measurements are
**host measurements**, not Linux embedded ARM measurements.

All Cargo commands below ran from `agent/`. Output is captured from the actual
commands (build timings and process identifiers are specific to this run).
The lockfile records dependency versions and checksums.

## Environment

```sh
rustc -Vv
cargo -V
cargo fmt --version
cargo clippy --version
rustup target list --installed
uname -srm
sw_vers
id -u
command -v arm-linux-gnueabihf-gcc aarch64-linux-gnu-gcc
```

The final `command -v` exited 1 with no paths; the other environment commands succeeded.

```text
rustc 1.89.0 (29483883e 2025-08-04)
binary: rustc
commit-hash: 29483883eed69d5fb4db01964cdf2af4d86e9cb2
commit-date: 2025-08-04
host: aarch64-apple-darwin
release: 1.89.0
LLVM version: 20.1.7
cargo 1.89.0 (c24e10642 2025-06-23)
rustfmt 1.8.0-stable (29483883ee 2025-08-04)
clippy 0.1.89 (29483883ee 2025-08-04)
aarch64-apple-darwin
aarch64-unknown-linux-gnu
armv7-unknown-linux-gnueabihf
Darwin 25.6.0 arm64
ProductName:		macOS
ProductVersion:		26.6.2
BuildVersion:		25G83
501
```

## Formatting

```sh
cargo fmt --all -- --check
```

Exit 0; no output. `cargo fmt --all` had been used to format sources before this check.

```text
(no output)
```

## Clippy

```sh
cargo clippy --all-targets -- -D warnings
```

Exit 0.

```text
    Checking devicebridge-agent v0.1.0 (/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-002-rust-agent/agent)
    Finished `dev` profile [unoptimized + debuginfo] target(s) in 0.24s
```

## Tests

```sh
cargo test
```

Exit 0; 10 unit tests and 5 subprocess integration tests.

```text
   Compiling devicebridge-agent v0.1.0 (/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-002-rust-agent/agent)
    Finished `test` profile [unoptimized + debuginfo] target(s) in 0.35s
     Running unittests src/main.rs (target/debug/deps/devicebridge_agent-74eca58069441435)

running 10 tests
test cli::tests::invalid_utf8_options_and_values_are_errors ... ok
test cli::tests::defaults_actions_and_overrides ... ok
test config::tests::read_errors_are_reported ... ok
test config::tests::defaults_and_boundaries_are_valid ... ok
test config::tests::input_bound_is_enforced_and_exact_limit_is_accepted ... ok
test config::tests::malformed_unknown_duplicate_and_invalid_values_fail_closed ... ok
test cli::tests::rejects_bad_duplicate_or_conflicting_arguments ... ok
test logging::tests::logs_are_one_json_line_with_escaped_input ... ok
test logging::tests::write_failure_is_returned ... ok
test cli::tests::file_values_are_overridden_by_cli ... ok

test result: ok. 10 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out; finished in 0.00s

     Running tests/process.rs (target/debug/deps/process-936e39427d23adba)

running 5 tests
test unix::fifo_config_is_rejected_without_waiting_for_a_writer ... ok
test file_cli_precedence_and_invalid_file_cannot_be_masked ... ok
test informational_commands_and_invalid_options ... ok
test missing_directory_oversized_and_malformed_config_are_rejected ... ok
test unix::sigint_and_sigterm_exit_successfully_with_structured_shutdown ... ok

test result: ok. 5 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out; finished in 0.45s

```

## Host release build

```sh
cargo build --release --locked
```

Exit 0.

```text
   Compiling devicebridge-agent v0.1.0 (/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-002-rust-agent/agent)
    Finished `release` profile [optimized] target(s) in 2.54s
```

## ARM standard libraries

```sh
rustup target add armv7-unknown-linux-gnueabihf aarch64-unknown-linux-gnu
```

Exit 0; this installs standard libraries on the build host, not on a device.

```text
info: downloading component 'rust-std' for 'armv7-unknown-linux-gnueabihf'
info: installing component 'rust-std' for 'armv7-unknown-linux-gnueabihf'
info: downloading component 'rust-std' for 'aarch64-unknown-linux-gnu'
info: installing component 'rust-std' for 'aarch64-unknown-linux-gnu'
```

## ARMv7 Linux build attempt

```sh
cargo build --release --locked --target armv7-unknown-linux-gnueabihf
```

Exit 101; no linked ARMv7 executable produced.

```text
   Compiling libc v0.2.191
   Compiling serde_core v1.0.229
   Compiling zmij v1.0.23
   Compiling serde v1.0.229
   Compiling serde_json v1.0.154
   Compiling signal-hook v0.3.18
   Compiling memchr v2.8.3
   Compiling itoa v1.0.18
   Compiling errno v0.3.14
   Compiling signal-hook-registry v1.4.8
   Compiling devicebridge-agent v0.1.0 (/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-002-rust-agent/agent)
error: linker `arm-linux-gnueabihf-gcc` not found
  |
  = note: No such file or directory (os error 2)

error: could not compile `devicebridge-agent` (bin "devicebridge-agent") due to 1 previous error
```

## AArch64 Linux build attempt

```sh
cargo build --release --locked --target aarch64-unknown-linux-gnu
```

Exit 101; no linked AArch64 Linux executable produced.

```text
   Compiling libc v0.2.191
   Compiling serde_core v1.0.229
   Compiling zmij v1.0.23
   Compiling signal-hook v0.3.18
   Compiling serde_json v1.0.154
   Compiling serde v1.0.229
   Compiling memchr v2.8.3
   Compiling itoa v1.0.18
   Compiling errno v0.3.14
   Compiling signal-hook-registry v1.4.8
   Compiling devicebridge-agent v0.1.0 (/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-002-rust-agent/agent)
error: linker `aarch64-linux-gnu-gcc` not found
  |
  = note: No such file or directory (os error 2)

error: could not compile `devicebridge-agent` (bin "devicebridge-agent") due to 1 previous error
```

## Host resource and shutdown baseline

Exact commands, run from `agent/` (the surrounding capture used
`set -o pipefail` and `tee /tmp/db002-measure.log`):

```sh
file target/release/devicebridge-agent
stat -f '%z bytes' target/release/devicebridge-agent
otool -L target/release/devicebridge-agent
./target/release/devicebridge-agent --config config.example.json --check-config
./target/release/devicebridge-agent --config config.example.json > /tmp/db002-idle-events.log 2>&1 &
db002_agent_pid=$!
trap 'kill -TERM "$db002_agent_pid" 2>/dev/null' EXIT
sleep 1
ps -o pid=,rss=,%cpu=,etime= -p "$db002_agent_pid"
sleep 1
ps -o pid=,rss=,%cpu=,etime= -p "$db002_agent_pid"
sleep 1
ps -o pid=,rss=,%cpu=,etime= -p "$db002_agent_pid"
kill -TERM "$db002_agent_pid"
wait "$db002_agent_pid"
db002_exit=$?
trap - EXIT
printf 'agent exit: %s\n' "$db002_exit"
cat /tmp/db002-idle-events.log
```

Exit 0. Actual output (ps columns: PID, RSS in KiB, %CPU, elapsed time):

```text
target/release/devicebridge-agent: Mach-O 64-bit executable arm64
352768 bytes
target/release/devicebridge-agent:
	/usr/lib/libiconv.2.dylib (compatibility version 7.0.0, current version 7.0.0)
	/usr/lib/libSystem.B.dylib (compatibility version 1.0.0, current version 1356.0.0)
configuration valid
29952   1456   0.0 00:01
29952   1456   0.0 00:02
29952   1456   0.0 00:03
agent exit: 0
{"timestamp_ms":1793738874209,"level":"info","event":"agent_started","message":"M01 foundation; no transport"}
{"timestamp_ms":1793738877301,"level":"info","event":"agent_stopped","message":"shutdown requested"}
```

Host Mach-O release file: **352,768 bytes** (344.5 KiB).
Three short idle samples: **1,456 KiB RSS**, **0.0% CPU** at ps precision.
These are neither peak/steady-state benchmarks nor Linux ELF sizes, and do not
prove the <=16 MB RSS target on a 128 MB ARM device. No real ARMv7/AArch64 Linux,
Buildroot, Yocto, Debian/Ubuntu device, model/kernel, latency distribution, CPU
baseline under load, restart behavior on hardware, or flash activity was measured.
The log timestamps demonstrate a clean host SIGTERM exit; subprocess tests also
exercise SIGINT. The host ran as uid 501, without root.

## Negative cases and self-review

The first `cargo test` exited 101: Serde accepted an empty positional JSON array
as a default struct. Actual failure excerpt:

```text
---- config::tests::malformed_unknown_duplicate_and_invalid_values_fail_closed stdout ----

thread 'config::tests::malformed_unknown_duplicate_and_invalid_values_fail_closed' panicked at src/config.rs:136:13:
accepted []
```

Fixed by requiring the root JSON value to begin with an object before typed
deserialization. The complete final passing test output appears above.
Self-review also caught FIFO blocking during open, before metadata validation;
Unix opens now use O_NONBLOCK and then require a regular file. The subprocess
FIFO test has a three-second deadline and proves rejection without a writer.
No unresolved implementation findings were identified in this self-review;
independent coordinator review is still required.

Coverage includes config exact-size/oversize boundary, integer limits/overflow,
bad/unknown/duplicate fields, missing and wrong types, empty/invalid UTF-8 JSON,
non-object roots, trailing data, missing files/directories/FIFO, injected read
errors, malformed/duplicate/conflicting CLI flags, UTF-8 errors, file/CLI
precedence, invalid config despite an overriding CLI value, JSON escaping, log
write failure, help/version/check-config, and real Unix INT/TERM subprocess exits.

Source review inspected all new files and searched production sources for unsafe
blocks, network types, process execution and panicking unwrap/expect calls:

```sh
git diff --cached --check
rg -n 'unsafe|net::|TcpStream|UdpSocket|Command::|unwrap\(|expect\(' src
```

`git diff --cached --check` exited 0 without output. The rg matches were only unwraps
inside unit-test modules; there are no unsafe blocks, network types or child
commands in production source. Test helpers use host `kill`/`mkfifo`; these
are not device runtime dependencies. Manual review found no secrets, tokens,
private keys or production configuration in the changed files. This is a
scope/source review, not a dependency vulnerability audit or formal security proof.

## Security and resource assessment

Configuration reads at most 4097 bytes to distinguish the 4096-byte limit.
CLI option/value length is limited to 1024 bytes with a bounded option set and
duplicate rejection. Typed deserialization denies unknown and duplicate fields;
polling is constrained to 10–1000 ms. Config errors exclude user-provided field
names and values, and logs escape messages as JSON.

Four direct dependency declarations: serde/serde_json reuse one parser/serializer
for config and logs; signal-hook supplies safe Unix signal registration; libc
supplies O_NONBLOCK to reject special files without FIFO blocking. No async
runtime, CLI framework, transport, TLS, file logger or extra runtime threads.
Serde's proc-macro dependencies execute only on the build host. Review the
lockfile and reassess dependencies as features grow.

Startup allocates for CLI/config/JSON parsing and signal registration. Normal
idle work is one atomic load and sleep, with no allocations, network or disk
operations in that loop. Default sleep is 100 ms (roughly 10 wakeups/second).
Release enables size optimization, LTO, one codegen unit and stripping.
No filesystem writes, pidfiles or persistent device state are created by the
runtime. Lifecycle events use stderr; the operator must choose a sink that
does not write to device flash if flash writes are prohibited.

No socket binds/connects, commands, credential defaults or authorization stubs
are exposed. Unprivileged operation needs only read access to an optional config
and writable stdout/stderr. Signals only set an atomic flag. Interrupted sleep
is checked on the next poll; scheduling and a blocked stderr sink can delay
shutdown. Signal-handler registration/log I/O failures produce runtime exit 1;
invalid input exits 2, successful check/shutdown exits 0. Failure to report an
error to a closed stderr sink is ignored without panicking.

GNU ARM targets require matching cross-linkers plus libc/sysroot, device loader
and ABI. Rust standard libraries alone do not supply these. Non-Unix runtime
returns an explicit unsupported error; non-Unix checks were not executed.
Musl/other libc variants, Linux dynamic-link dependencies and Linux signal
runtime behavior remain untested. There is no claim of remote connectivity or
production security guarantees and no architecture/security contract change.

## Rollback and remaining delivery

Before deployment, remove the foundation binary/config or restore the previous
artifact; there is no device state/database migration to undo. After review and
manual merge, a human can revert this issue's commit via a normal feature branch
and Draft PR into develop. Do not reset shared branches or change protections.
Stopping an installed instance uses SIGTERM/SIGINT; no root requirement is added.

Coordinator: independently review this commit, preserve these failed cross-build
results in the Draft PR, and provide the exact commit SHA and issue #2 link.
Complete both Linux ARM builds on a host with GNU C cross-toolchains and matching
sysroots, then run on representative 128 MB ARMv7/AArch64 devices before claiming
hardware resource acceptance. Human Tech Lead alone approves/manual-merges.
