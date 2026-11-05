# M01 CI (DB-006 / #6)

The [workflow](../.github/workflows/ci.yml) checks pull requests into `develop`,
pushes to `develop` and `feat/**`, and manual workflow dispatch where GitHub makes
it available. Draft PRs receive the same checks. There are no path filters that
silently skip a module. No deployment, release, PR approval or merge is performed.

| Job | Checks and outputs |
| --- | --- |
| Workflow and shell lint | actionlint 1.7.7, Bash syntax, ShellCheck, four CI failure-gate tests |
| Go, one job per module | formatting, vet, uncached verbose tests, race tests, native executable |
| Rust agent | rustfmt, Clippy with warnings denied, locked tests and release build |
| Linux cross-build, ARMv7 and AArch64 | three Go applications with CGO disabled, Rust agent with GNU cross-linkers, ELF details and checksums |

All jobs use GitHub-hosted `ubuntu-24.04` with explicit timeouts and cancel older
runs for the same ref. Go modules are independent: `server/control`, `server/relay`,
`cli`, and `packaging`. The packager is checked/built as build-host tooling; it is
not added to device dependencies or ARM application bundles. Compose integration remains DB-008; the security job is supplied by DB-007.

## Reproduce checks

Run from the repository root (the script also accepts invocation from another
directory). Build-host requirements: Bash, Git, Go, Rust/rustup, Python 3 for CI
gate tests, ShellCheck, `file` for ARM inspection, and SHA-256 utilities.

```sh
rustup toolchain install 1.89.0 --profile minimal --component rustfmt,clippy
bash ci/check.sh workflow
bash ci/check.sh go server/control
bash ci/check.sh go server/relay
bash ci/check.sh go cli
bash ci/check.sh go packaging
bash ci/check.sh rust
```

Go 1.23.4 and Rust 1.89.0 reproduce the existing foundation's tested baselines.
Go uses `GOTOOLCHAIN=local` and `GOWORK=off`; Rust uses `RUSTUP_TOOLCHAIN=1.89.0`
by default. Rustfmt/Clippy are installed explicitly. Go caching is disabled because
the four modules have no external dependencies or `go.sum`; actionlint's versioned
Go module downloads are separate build-host tooling. These pinned compiler versions
are not a certification of current support or vulnerability status. Toolchain
upgrades and dependency security evaluation require their own validation.

On an Ubuntu build host, install GNU cross-compilers/sysroots explicitly, then run:

```sh
sudo apt-get update
sudo apt-get install --no-install-recommends -y gcc-arm-linux-gnueabihf gcc-aarch64-linux-gnu libc6-dev-armhf-cross libc6-dev-arm64-cross
rustup toolchain install 1.89.0 --profile minimal --target armv7-unknown-linux-gnueabihf,aarch64-unknown-linux-gnu
bash ci/check.sh arm armv7
bash ci/check.sh arm arm64
```

Only provisioning on the disposable build host uses elevated privileges. Neither
applications nor embedded devices need root, compilers, Python, Docker or these
CI tools. GNU Rust binaries require compatible target libc, loader, CPU and kernel;
they are not universal Buildroot/Yocto binaries. These jobs cross-build binaries;
they do not execute ARM code or measure device memory/CPU.

## Logs and provenance

Each script invocation writes under ignored `dist/ci/<job>/`:

- `commands.txt`: working directories, exact command arguments, output and exit codes.
- `SOURCE.txt`: actual checked-out commit, event SHA, run ID/attempt, platform,
  job, selected Rust toolchain and dirty-workspace flag.
- `bin/` and `SHA256SUMS`: successful build outputs and hashes of binaries plus
  source metadata. Re-running a job replaces its log/metadata and removes any old
  checksum manifest before executing checks.

Artifact names include the event SHA and run attempt; retention is 14 days.
For PR events, the checkout can be GitHub's synthetic merge commit. Compare
`checked_out_commit` and `event_commit`; neither is assumed to be the PR head.
Tool versions and cross-compiler identity are captured in command output.
Verify a downloaded successful build from the directory containing its
`SOURCE.txt` and `SHA256SUMS` (Go/ARM ZIPs have a job subdirectory) with
`sha256sum -c SHA256SUMS` (or `shasum -a 256 -c SHA256SUMS`).

Failed checks still upload available logs using `if: always()`. A failed job can
contain partial or previous local binaries, but has no new success manifest.
Use only outputs from a successful job with matching run/source metadata.
Setup failures before script startup have workflow logs but may have no artifact.
Artifacts are development builds, with no release signature or independent
source attestation. No public binary release is published by this workflow.

## Security, resource effects and rollback

Token permission is `contents: read`; all unspecified permissions are disabled.
Checkout uses `persist-credentials: false`. Only `pull_request` executes PR code;
there is no privileged `pull_request_target`, repository secret injection,
self-hosted runner, shared build cache, arbitrary user command or deployment.
Module/target arguments are allowlisted before filesystem output. Failed commands
retain nonzero exit codes through `tee`; nonempty gofmt output fails the job.

Actions are pinned to explicit majors, as required by #6. Verified upstream
input/runtime definitions: [checkout v4](https://github.com/actions/checkout/blob/v4/action.yml),
[setup-go v5](https://github.com/actions/setup-go/blob/v5/action.yml), and
[upload-artifact v4](https://github.com/actions/upload-artifact/blob/v4/action.yml).
Major tags, hosted runner images and apt package repositories can change; builds
are traceable but not claimed byte-for-byte reproducible across runner revisions.
There is no dedicated secret scanner or dependency vulnerability audit in this
issue. DB-007 owns the security baseline; existing application negative tests run
here without changing security/transport/configuration contracts.

Application source/dependencies are unchanged. No device allocation, runtime
dependency or flash write is added. CI consumes hosted CPU/storage and retains
artifacts temporarily; no numeric performance or 128 MB acceptance is claimed.
No physical ARM hardware, fleet, RSS/CPU, Buildroot/Yocto or SSH tunnel validation
was performed. ADR-001 still says Proposed although merged foundation PR #11
records human approval; this issue preserves its decisions and status text.

Rollback is a reviewed revert of the DB-006 commits on a feature branch. It
removes the workflow/scripts and leaves applications/data unchanged. Local
generated outputs may be removed from `dist/ci/` deliberately; retained GitHub
artifacts expire after 14 days. No branch rules, approvals, auto-merge settings,
databases or deployed services are changed. Human Tech Lead review and manual
merge into `develop` remain mandatory.

See [EVIDENCE.md](EVIDENCE.md) for executed results and host/runner limitations.
# Security baseline (DB-007)

The `security` job scans full fetched Git history with Gitleaks v8.24.3 and
checks detection/redaction using a disposable synthetic fixture. Run locally:

```sh
go install github.com/zricethezav/gitleaks/v8@v8.24.3
bash ci/security.sh
```

Add the Go bin directory to `PATH`. Commit changes before scanning; this command
does not scan uncommitted files. Reports in `dist/ci/security/` are fully redacted.
See [security baseline](../security/README.md) for limitations and evidence.
