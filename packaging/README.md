# Installable M01 preview packages (DB-011 / #16)

DeviceBridge components can be shipped as separate portable archives. End users
extract an archive and run **`sh install.sh`**; no compiler is needed on the
destination. These are M01 foundation previews: the agent idles, Control serves
local probes, relay transport is disabled, and the CLI has no remote operations.
Installing a preview does not provide working remote SSH.

คู่มือติดตั้งภาษาไทย: [README.th.md](README.th.md).

## Install and run

Obtain the complete archive and checksum from a trusted publisher. This PR does
not publish a GitHub Release or set up a download endpoint. Verify its SHA-256
using a separately trusted checksum, then extract it. For example, for a macOS
ARM64 CLI preview (replace the version/target with your actual artifact):

```sh
shasum -a 256 -c devicebridge-cli-0.1.0-preview.4-darwin-arm64.tar.gz.sha256
tar -xzf devicebridge-cli-0.1.0-preview.4-darwin-arm64.tar.gz
cd devicebridge-cli-0.1.0-preview.4-darwin-arm64
sh install.sh
"$HOME/.local/devicebridge/cli/bin/devicebridge" help
"$HOME/.local/devicebridge/cli/bin/devicebridge" check
```

On Linux, use `sha256sum -c ARCHIVE.sha256` for the outer checksum. The installer
independently verifies the payload against its embedded SHA-256 before writing.
Checksums detect accidental corruption, **not publisher authenticity**: an attacker
who replaces the script/archive and checksum can replace both. Signed release
distribution is separate future work; do not download and pipe scripts into a shell.

The default private prefix is `$HOME/.local/devicebridge/COMPONENT`. To choose
another dedicated prefix, use `sh install.sh --prefix /absolute/private/path`.
Do not use `sudo`. Existing unrelated directories/files and symlinked destination
paths are rejected. Use physical paths without `..`, `.`, repeated slashes, trailing
slashes or newlines. Symlinked HOME is resolved to its physical directory. Each
component has its own prefix; prefixes cannot be shared between components/targets.
Choose a filesystem that permits executable files (`noexec` mounts cannot run the
binary). The installer deliberately does not execute a payload to probe this.

| Component | Installed command | Run after installation |
|---|---|---|
| agent | `.../agent/bin/devicebridge-agent` | Start in foreground; Ctrl-C stops it |
| control | `.../control/bin/devicebridge-control` | Start in foreground; probe `http://127.0.0.1:8080/healthz` |
| relay | `.../relay/bin/devicebridge-relay` | Start the idle skeleton in foreground; Ctrl-C stops it |
| cli | `.../cli/bin/devicebridge` | `help`, `version`, `check` |

For example, run the installed agent with
`"$HOME/.local/devicebridge/agent/bin/devicebridge-agent"`. Defaults require no
config file. Installation never starts a process, edits your shell profile/PATH,
enables boot services, opens ports, creates credentials, or overwrites configuration.
You may add a component's `bin` directory to PATH yourself. M01 Control remains
loopback-only; relay/agent/CLI perform no remote transport. PostgreSQL development
setup remains [DB-005 / PR #15](https://github.com/LiiWzz/DeviceBridge/pull/15),
separate from these binary packages and never installed on embedded devices.

## Upgrade, rollback and uninstall

Extract the new archive and run its `install.sh` with the **same prefix**. It
checks the payload/host before mutation, stages the file, hashes the staged copy,
then activates by replacing one symlink on the same filesystem. No partially
copied executable is activated. Reinstalling identical version/content is a no-op
for the binary; reusing a version with different content fails. Installation never
executes the supplied binary.
Upgrades refresh the installed management script so fixes reach existing users.
If executable activation fails, the previous executable and rollback target remain
available; if history replacement fails, the installer restores the old executable.
The management script may already have been refreshed after such a failed upgrade.
Atomicity covers the executable link, not a journaled multi-file transaction;
unrecoverable filesystem errors or power loss require inspection.

```sh
# The installed management script knows its own prefix (including custom paths).
sh "$HOME/.local/devicebridge/agent/manage.sh" --status
sh "$HOME/.local/devicebridge/agent/manage.sh" --rollback
# Stop the running agent with Ctrl-C/SIGTERM before uninstalling.
sh "$HOME/.local/devicebridge/agent/manage.sh" --uninstall
```

Rollback verifies and swaps the current/previous executable links; another rollback
switches back. Earlier release files are retained until uninstall (check flash
capacity before repeated upgrades). Processes already running keep their old
executable mapping; stop and restart them deliberately to use a new/rolled-back
version. No process is killed by the installer.

Uninstall verifies all managed releases before removing anything. It removes only
managed binaries/checksums/links/manager/marker. User configuration and other files
outside release directories are preserved; extra files *inside* release directories
cause a refusal. No database/volume/config is deleted. A modified or corrupt
installed state fails closed and needs manual inspection rather than forced removal.

Only install/upgrade/rollback/uninstall write to the prefix. Status only reads.
An atomic `mkdir` lock excludes other installers. If interrupted, normal signal
cleanup removes temporary files/lock. After power loss/SIGKILL inspect the prefix
and confirm no installer is running before removing **only the empty**
`.install-lock` with `rmdir`. Staging files or corrupt state need inspection; no
automatic stale-lock takeover or crash/power-loss durability guarantee is made.
An initial disk-full/permission failure can leave a partial private prefix; inspect
it before retrying. No global/system directories are managed.

## Platforms and destination prerequisites

| Target | Host detection | Typical use |
|---|---|---|
| `linux-armv7` | Linux + `armv7l` | Embedded ARMv7 agent |
| `linux-arm64` | Linux + `aarch64` / `arm64` | Agent, server or CLI |
| `linux-amd64` | Linux + `x86_64` / `amd64` | Server or CLI |
| `darwin-arm64` | Darwin + `arm64` | Developer CLI |
| `darwin-amd64` | Darwin + `x86_64` | Developer CLI |

The packager checks ELF/Mach-O architecture/class/type before creating archives;
the installer checks OS/CPU without an override. ARMv7 ISA attributes, CPU features,
libc/loader versions, distribution and kernel compatibility need separate validation.
In particular, a GNU-linked Rust build is **not** a universal musl/Buildroot/Yocto
binary. No ABI auto-detection or automatic runtime installation is attempted.
The local `0.1.0-preview.4` agent artifacts use statically linked musl instead;
those ran in ARM64 Linux containers and ARMv7 emulation. Neither result establishes
compatibility with every embedded kernel, CPU or distribution.

Destination: POSIX shell, standard file utilities including `mktemp` and `readlink`,
and at least one of `sha256sum`, `shasum -a 256`, or `openssl dgst -sha256`.
BusyBox builds must enable those required applets; full BusyBox/container results
are not a promise about every customized embedded image. Only a binary, small
management script and checksum are installed. No Rust/Go compiler, Python, Node,
Docker or package manager is required on an embedded device; no daemon or periodic
installation work is added. Flash writes occur during explicit installation and
version changes, never in an installer idle loop (there is none).

Use a private user-owned prefix with trusted parents. Unprivileged shell filesystem
checks are not a security boundary against the same UID or someone who can replace
parent directories. The installer does not need root or create a privileged service.

## Build packages (maintainers only)

Run from `packaging/` on a workstation with Go 1.23+:

```sh
go run . --component agent --version 0.1.0-preview.4 --target linux-arm64 \
  --binary /path/to/prebuilt/devicebridge-agent \
  --revision FULL_40_CHARACTER_SOURCE_COMMIT --output ../dist
```

Use the same command with `control`, `relay` or `cli` and the matching binary/target.
The input is never executed. The tool streams a portable archive with fixed paths,
an embedded installer, README, payload and SHA256SUMS, plus an outer `.sha256`.
It refuses existing output names and symlink/non-executable/mismatched inputs.
Version/revision fields are validated. Revision is an operator declaration; this
tool does not attest that a binary was built from that source. CI provenance and
release signing require separate reviewed work.

The tool consumes explicit artifacts and does not import, rebuild, merge or copy
the implementations from Draft PRs #12–#15. Preview integration against their
exact commits is labeled in [EVIDENCE.md](EVIDENCE.md). Production release wiring
and DB-006/DB-007 remain dependent on human-approved/manual-merged implementations.
No architecture or security contract changes are made here.

```sh
gofmt -l .
go vet ./...
go test -count=1 -v ./...
go test -race ./...
go build -trimpath -o bin/devicebridge-package .
sh -n install.sh
```

See [EVIDENCE.md](EVIDENCE.md) for executed results, security review and limitations.
