# DeviceBridge agent (DB-002)

M01 Rust binary foundation only. No enrollment, networking, tunnels, remote commands,
authentication, or device credentials are implemented. Run as an ordinary user.

From this directory:

```sh
cargo build --release --locked
./target/release/devicebridge-agent --help
./target/release/devicebridge-agent --config config.example.json --check-config
./target/release/devicebridge-agent --config config.example.json
```

SIGINT or SIGTERM requests a clean exit on Unix. Handlers only set an atomic flag;
the main loop checks every 100 ms by default (configurable 10–1000 ms). Scheduling
and a blocked log sink can delay exit; this is not a hard real-time guarantee.
Other platforms support validation/help but the runtime returns an unsupported error.

Optional configuration is a regular JSON file of at most 4096 bytes, read once at
startup. Defaults are `log_level: "info"` and `shutdown_poll_ms: 100`.
Unknown/duplicate fields, invalid UTF-8/JSON, invalid types and out-of-range values
are errors. CLI flags `--log-level` and `--shutdown-poll-ms` override valid file
values; an invalid file is always rejected. `--config PATH` selects the file;
without it defaults apply. No environment configuration is read. Option values
are separate arguments, at most 1024 bytes each; repeated flags and conflicting
actions are rejected. `--help`/`--version` do not open the configuration file.

Logs are JSON lines on stderr with `timestamp_ms` (Unix epoch, zero if clock is
before epoch), `level`, `event`, and `message`. Normal events are startup/shutdown;
the idle loop does no logging, allocation, disk I/O, or network I/O. Log level
`error` suppresses lifecycle events; `debug` currently adds no events. Exit codes
are 0 for success, 2 for invalid CLI/config, and 1 for runtime or output failures.
No daemonization, pidfile, state directory, flash writes, or root privileges are
needed. Direct stdout/stderr to a suitable sink; external logging retention is
the operator's responsibility.

Host checks (Rust toolchain with rustfmt and clippy installed):

```sh
cargo fmt --all -- --check
cargo clippy --all-targets -- -D warnings
cargo test
cargo build --release --locked
```

Linux ARM GNU cross-builds, from `agent/` so `.cargo/config.toml` is discovered:

```sh
rustup target add armv7-unknown-linux-gnueabihf aarch64-unknown-linux-gnu
cargo build --release --locked --target armv7-unknown-linux-gnueabihf
cargo build --release --locked --target aarch64-unknown-linux-gnu
```

Install build-host C cross-toolchains `arm-linux-gnueabihf-gcc` and
`aarch64-linux-gnu-gcc`, matching the target glibc/sysroot. Override the linker
with `CARGO_TARGET_ARMV7_UNKNOWN_LINUX_GNUEABIHF_LINKER` or
`CARGO_TARGET_AARCH64_UNKNOWN_LINUX_GNU_LINKER` if required. These are hard-float
ARMv7 and AArch64 glibc builds, not universal Buildroot/Yocto images. Match the
device CPU, loader, libc and ABI; other libc targets require separate build
configuration and validation. Build tools run on the host only; the device
requires neither Cargo nor Docker, Python or Node.

See [EVIDENCE.md](EVIDENCE.md) for actual results, host-only memory/size baseline,
toolchain gaps, remaining hardware validation, security assessment and rollback.
The repository's ADR-001 still says Proposed; the coordinator reports human
approval in merged foundation PR #11. That stale status is disclosed here without
changing the shared ADR. This issue does not change architecture or security contracts.
