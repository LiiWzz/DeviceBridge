# DB-011 / #16 — installation evidence

Recorded 2026-11-04. Feature branch `feat/m01-db-011-easy-install`, isolated
worktree `/private/tmp/devicebridge-db011`, based on fetched `origin/develop`
`c8fae1481fe657f385d7747878780d1f4485e986`. Required README, AGENTS, execution plan
and architecture ADRs read before coding. DB-001 is merged in PR #11; its recorded
ADR approval and stale checked-in Proposed status remain disclosed, not rewritten.

This issue adds packaging only. No application implementation was copied or changed,
and DB-006/DB-007 were not started against unapproved source PRs. Four existing
application PRs remain drafts awaiting human approval/manual merge. No issue closure,
PR merge/approval, auto-merge, public release publication, production deployment,
system service, host user installation or shell-profile change was performed.

## Executed checks

Build host: macOS ARM64, Go 1.23.4, Rust 1.89.0, Docker 29.8.2.
The new Go module uses only the standard library (`go list -m all` prints only
`github.com/LiiWzz/DeviceBridge/packaging`). Go is required on the build host only.

| Executed command / context | Actual result |
|---|---|
| `gofmt -l .` in packaging | Exit 0, no output |
| `go vet ./...` | Exit 0, no output |
| `go test -count=1 -v ./...` | Exit 0; 10 top-level groups, all host cases passed |
| `go test -race ./...` | Exit 0 |
| `go build -trimpath -o bin/devicebridge-package .` | Exit 0 |
| `sh -n install.sh` and `sh -n smoke-preview.sh` | Both exit 0 |
| `git diff --check` | Exit 0 |
| ShellCheck 0.10.0, both shell files, `-s sh` | Exit 0, no output |
| Linux ARM64 compiled Go test binary, Debian dash, UID 1000 | Exit 0; shasum/OpenSSL absent, those fallback cases explicitly skipped |
| Same Go tests with BusyBox shell/utilities, UID 1000 | Exit 0; absent shasum/OpenSSL skipped; missing-hash and history-mv injection simulations skipped because ash exposes built-in applets |
| 13 final package builder invocations | All exit 0; every outer archive checksum independently matched |
| Actual macOS ARM64 CLI archive install/run/check/status/uninstall | Exit 0; isolated temporary prefix removed |
| Actual four-component ARM64 archives, Debian dash | Exit 0, complete install/run/shutdown/uninstall smoke |
| Actual four-component ARM64 archives, minimal BusyBox 1.37.0 | Exit 0, complete smoke, Docker Linux ARM64 VM |
| Actual four-component ARMv7 archives, BusyBox 1.37.0 | Exit 0, complete smoke, Docker ARMv7 **emulation** |

Captured output: [host tests](evidence/mac-tests.txt), [race](evidence/mac-race.txt),
[Linux dash tests](evidence/linux-dash-tests.txt), [BusyBox tests](evidence/linux-busybox-tests.txt),
[native CLI](evidence/preview-mac-cli.txt), [ARM64 dash](evidence/preview-arm64-dash.txt),
[ARM64 BusyBox](evidence/preview-arm64-busybox.txt), [ARMv7 emulation](evidence/preview-armv7-emulated.txt),
[final checks](evidence/final-checks.txt), [all package commands](evidence/package-build.txt).
Trailing whitespace alone was removed when preserving command output.
Go output tabs are rendered as literal `\t` to avoid Git whitespace errors.

The installer lifecycle tests exercise first install, identical repeat without
rewriting the executable, two versions, repeated rollback, custom paths with
spaces, uninstall and preservation of user configuration. Negative tests cover
corrupted/missing/symlink payloads, symlink package directories, wrong targets,
invalid metadata/arguments, empty/relative/traversal prefixes, unrelated existing
files/directories, symlink ancestors/active paths, conflicting version contents,
foreign locks, extra release files, malformed installed state and missing hash
utilities. The build tool rejects invalid metadata, scripts in place of ELF/Mach-O,
symlink inputs, wrong CPU/type, existing archive/checksum paths, and preserves
preexisting output on failure. All three hashing implementations passed on macOS.

## Exact platform verification commands

Linux installer test suite is cross-compiled from the **new packaging module**:

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go test -c -o ../dist/packaging-tests-linux-arm64 .
```

It was executed as a static test binary in the container below, with `/tests`
bound read-only to this worktree's `dist` directory. The BusyBox variation created
`/tmp/busybox-bin` with `busybox --install -s`, prepended that directory to PATH,
and ran the same test binary. No Go compiler is needed in the runtime container.

Final archive smoke used these arguments (all commands exit 0):

```sh
docker run --rm --network none --read-only \
  --tmpfs /tmp:rw,exec,nosuid,nodev,size=128m \
  --cap-drop ALL --security-opt no-new-privileges --user 1000:1000 \
  --mount type=bind,source=/private/tmp/devicebridge-db011/dist/preview.4/packages,target=/packages,readonly \
  --mount type=bind,source=/private/tmp/devicebridge-db011/packaging,target=/verification,readonly \
  --platform linux/arm64 --entrypoint /bin/busybox busybox@sha256:d82c2ab94640ded77cf76514ce6a84870761105058a4a9e51b05a8a79be97a6c \
  timeout 30 sh /verification/smoke-preview.sh /packages linux-arm64 0.1.0-preview.4
```

ARMv7 changed platform to `linux/arm/v7`, timeout to 60 and target to `linux-armv7`.
It used manifest `sha256:91aa2a66a996afa9821eed6758ed0489ee891c50f4b8940b7b3db500b79f029f`.
The dash variation used `--entrypoint /bin/sh devicebridge-m01-linkers:bookworm`
followed by `/verification/smoke-preview.sh /packages linux-arm64 0.1.0-preview.4`.
The test binary used the same isolation flags with a 256 MiB writable tmpfs.

The smoke script checks outer checksums, installs/reinstalls/status for every
component, runs CLI version/check and a rejected insecure endpoint, runs agent
version/config validation and real SIGTERM lifecycle, asserts Control health and
readiness JSON over loopback, starts/stops the idle relay, and removes every
installation. The script exercises SIGTERM; the source component's prior tests
cover SIGINT. It never calls remote operations. All test writes occur in temporary
container storage; no embedded device or actual user prefix was modified.

BusyBox tested manifest digest:
`sha256:a7dbc319b1438115a207e11930dc4bb3e643dcb3ec0bdd694bb36cd8fad289ec`.
Debian tools/test image:
`sha256:0c869e5d03a5a18dc6983598386a28133f47ca754ade1eb2c0deb7c432032812` (ARM64).
ShellCheck image `koalaman/shellcheck:v0.10.0`, digest
`sha256:2097951f02e735b613f4a34de20c40f937a6c8f18ecb170612c88c34517221fb`.

An initial smoke exited 126: Docker's test tmpfs did not allow execution, so the
installed binary returned Permission denied. The dedicated test tmpfs was then
explicitly mounted `exec`; all final smokes passed. No production mount was changed.
Documentation tells users to choose an executable filesystem. An initial BusyBox
test failed because its built-in SHA applet cannot be hidden by removing it from
PATH; the test harness now detects that capability and explicitly skips only the
inapplicable missing-utility simulation. Actual missing-utility rejection passed
on macOS and Debian dash. Initial ShellCheck predicate/style warnings were fixed
with explicit conditionals and an explicitly empty CDPATH before final checks.
An initial concurrent mutable-tag run selected an ARMv7 BusyBox image for the
ARM64 smoke; the installer correctly rejected the architecture mismatch. Final
smokes pin the architecture-specific manifests above and both exited 0.

Independent Spec review reproduced a P2 rollback-history loss after a failed
activation. It was fixed by committing the executable before history and restoring
the old executable if history replacement fails. New tests simulate unwritable
binary storage and a failed history rename, then verify the original rollback
still works. Both run on unprivileged macOS/Debian dash; BusyBox cannot intercept
its built-in mv through PATH, so only the injected-rename case is skipped there.
An additional upgrade test verifies replacement of an obsolete installed manager.
The reviewer checked both fixes in follow-ups with no remaining normal-failure
finding. Final full suites and regenerated preview.4 smokes include these fixes.

## Source provenance and ARM build gap resolution

Preview binary source checkouts remained clean at these commits:

| Component | Source commit / draft PR |
|---|---|
| agent | `dcefd3b2f8b1de53bfd180d9613d49d6a2539b54` / #13 |
| control | `5fca32be4d6cb67d6a3e013be6be7b855c0c1475` / #14 |
| relay and CLI | `4518fc8d3edb4a172cb763cca7222a8afdacb2cb` / #12 |

The earlier DB-002 GNU builds failed because macOS had no GNU cross-linkers.
This task subsequently **completed both GNU Linux links** using build-host Docker
GCC 12.2 and matching Debian sysroots. The exact linker wrappers are preserved in
[ARM64 wrapper](evidence/gnu-link-arm64.sh), [ARMv7 wrapper](evidence/gnu-link-armv7.sh),
and [build-host Dockerfile](evidence/linkers.Dockerfile). Docker receives source
and Rust standard-library mounts; the latter is read-only, linking is non-root
and network-disabled. GCC is never installed on a destination device.

Actual commands from the unchanged agent checkout's `agent/` directory:

```sh
TMPDIR=/private/tmp/db011-rust-tmp CARGO_TARGET_AARCH64_UNKNOWN_LINUX_GNU_LINKER=/private/tmp/db011-link-arm64 cargo build --release --locked --target aarch64-unknown-linux-gnu
TMPDIR=/private/tmp/db011-rust-tmp CARGO_TARGET_ARMV7_UNKNOWN_LINUX_GNUEABIHF_LINKER=/private/tmp/db011-link-armv7 cargo build --release --locked --target armv7-unknown-linux-gnueabihf
```

Both exit 0, release build completion times 6.66s / 8.81s in this run. `file`
confirmed GNU dynamic ELF with `/lib/ld-linux-aarch64.so.1` /
`/lib/ld-linux-armhf.so.3`. Those preliminary GNU artifacts are not the final
portable preview package payloads.

For fewer destination dependencies, the **final preview.4 agents are static musl**:

```sh
rustup target add aarch64-unknown-linux-musl armv7-unknown-linux-musleabihf
task_rust_sysroot=$(rustc --print sysroot)
CARGO_TARGET_AARCH64_UNKNOWN_LINUX_MUSL_LINKER="$task_rust_sysroot/lib/rustlib/aarch64-apple-darwin/bin/rust-lld" RUSTFLAGS='-C linker-flavor=ld.lld' cargo build --release --locked --target aarch64-unknown-linux-musl
CARGO_TARGET_ARMV7_UNKNOWN_LINUX_MUSLEABIHF_LINKER="$task_rust_sysroot/lib/rustlib/aarch64-apple-darwin/bin/rust-lld" RUSTFLAGS='-C linker-flavor=ld.lld' cargo build --release --locked --target armv7-unknown-linux-musleabihf
```

All exit 0. Actual release completion times 8.51s / 12.37s. `file` confirmed
statically linked, stripped ARM64 and ARM EABI5 ELF executables. The standard
libraries and rust-lld were on the build host only. No Rust source, Cargo contract
or dependency version was changed. Both resulting archives passed the real binary
smoke above; ARMv7 execution used emulation, ARM64 execution used Docker's Linux VM.

New Go AMD64 builds used `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath`
for Control, relay and CLI, with respective `./cmd/control`,
`./cmd/devicebridge-relay`, `./cmd/devicebridge` and ignored `bin/*-linux-amd64`
outputs. CLI macOS Intel used `CGO_ENABLED=0 GOOS=darwin GOARCH=amd64` with the same
CLI command and `bin/devicebridge-darwin-amd64`. All exit 0. Linux ARMv7/ARM64 and
native Mac CLI inputs were the previously tested unchanged component artifacts.
Linux AMD64 and macOS Intel binaries were packaged/format-checked, **not executed**.

Thirteen final archives and outer checksums exist locally under
`/private/tmp/devicebridge-db011/dist/preview.4/packages`; they are ignored build
artifacts, not Git commits or public release assets. Archive sizes and hashes:
[ARTIFACTS.md](evidence/ARTIFACTS.md). The assembled local ZIP is 25,047,200 bytes,
SHA-256 `b5c358eead9d165d0dbf85e87829ac4f8ca6df19a9defc0a38cfa850b20d237d`.
The package version labels distribution previews; underlying skeleton version
commands still report their existing 0.1.0 / 0.1.0-dev values. Declared source
revision plus payload hash identify inputs; no signed build attestation is claimed.

## Security, resources and rollback

Offline trusted-bundle operation only. No URL downloads, shell eval, application
payload execution during installation, credential generation, permissive auth,
root requirement, automatic service, global PATH/config mutation or remote port
exposure. Checksums detect corruption, not publisher identity; signed production
distribution is explicitly outside this preview. The builder validates metadata
before shell embedding, accepts only known component/target names, streams binary
content and verifies the copied bytes. It refuses output collisions atomically.

Installer validates source hash before destination writes, rejects destination
symlinks/unrelated state, uses a private per-component prefix and mkdir lock,
stages/hash-checks bytes before activation, then renames a symlink to a regular
executable on the same filesystem. No truncated payload becomes active. No daemon,
network or periodic work is added. Explicit installations/upgrades write executable
files/checksums to storage; this is deliberate installation flash I/O, not runtime
flash logging. Old versions consume disk/flash until uninstall. Hash/copy/install
time, peak installer RSS, CPU, physical flash behavior and device performance were
not benchmarked; no numeric resource guarantee is claimed.

Private prefixes require trusted parents. Shell checks do not defend against the
same UID or an attacker who can replace parent directories. Interrupted updates
may leave an empty lock/staging directory; no forced takeover or fsync/power-loss
durability claim. Root/global installation, package-manager/boot integration,
Windows and customized BusyBox applet sets are not validated. Static ELF still
requires compatible CPU/kernel; no physical ARMv7/AArch64, 128 MB device,
Buildroot/Yocto/Debian fleet, idle RSS target or latency acceptance is established.

Rollback: stop the running component explicitly, run installed `manage.sh --rollback`,
restart it to use the previous image. Uninstall verifies/removes managed releases
and preserves user files outside release directories; it never deletes database
volumes. Reverting this issue's commits removes tooling/docs, not user installations.
No host user installation was performed by this task. Temporary test prefixes and
containers were removed; build-host Docker images/Rust targets and local download
artifacts are retained for reproducibility. Human review/manual merge remains the
delivery gate. Dedicated secret-scanner certification is not claimed.

Independent Standards and Spec reviews of the pinned implementation are recorded
separately in [REVIEW.md](REVIEW.md); final evidence/documentation additions are
the coordinator's responsibility.
