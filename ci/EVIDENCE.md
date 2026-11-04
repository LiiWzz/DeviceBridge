# DB-006 executed evidence

Issue: https://github.com/LiiWzz/DeviceBridge/issues/6

Branch: `feat/m01-db-006`, based on fetched `origin/develop` at
`afc6cdf2bf3f7b8d028df64b766a66ffbf910393`. Dependencies are implemented in this
base: merged PRs #13 (agent), #14 (control), #12 (relay/CLI). No pending application
PR code is imported. Read README.md, AGENTS.md, docs/M01_EXECUTION.md and ADR-001
before implementation. No architecture/security/configuration contract changed.

## Scope

Added GitHub Actions format/lint/test/native-build jobs, GNU Linux ARM cross-build
jobs, traceable artifact metadata/checksums, CI failure-gate tests, and reproduction
instructions. Checks include all four independent Go modules, including the merged
packager. ARM jobs build only the three Go applications plus the Rust agent.

## Local verification

Recorded on macOS Darwin 25.6.0 ARM64, using Go 1.23.4, Rust/Cargo 1.89.0,
actionlint 1.7.7, and ShellCheck 0.10.0 in Docker. Python runs only build-host CI
failure-gate tests. Baseline compiler versions are reproduced, not certified as
currently supported or vulnerability-free.

The final invocation runs the following independent modes, each with output/exit
codes saved by `ci/check.sh` and the invoking process:

```sh
bash ci/check.sh workflow
bash ci/check.sh go server/control
bash ci/check.sh go server/relay
bash ci/check.sh go cli
bash ci/check.sh go packaging
bash ci/check.sh rust
bash ci/check.sh arm armv7
bash ci/check.sh arm arm64
```

All eight modes exited 0. Full results are recorded in
`evidence/local-results.json`. Native Go formatting/vet/tests/race/builds passed
for all four modules. Rust passed 10 unit and 5 process tests, formatting, Clippy
and release build. Both architectures built all four application binaries.
Workflow lint and all four failure-gate tests passed. All seven build output
checksum manifests were verified with `shasum -a 256 -c SHA256SUMS`, exit 0.
Each corresponding `evidence/<job>.txt` contains actual command directories,
arguments, outputs and exit codes. Tabs in committed text logs are rendered as
literal `\t` to avoid whitespace-check errors. The working copy was uncommitted
when local checks ran: metadata records the base SHA and `workspace_dirty=true`;
this is not represented as a clean committed-source build. The delivered script
and workflow hashes are recorded alongside these logs. Hosted CI separately checks
committed source.

macOS has no native GNU Linux cross-compilers. For these local cross-builds, only
compiler invocations use the existing build-host Docker image
`devicebridge-m01-linkers:bookworm`, image ID
`sha256:0c869e5d03a5a18dc6983598386a28133f47ca754ade1eb2c0deb7c432032812`
(ARM64 Linux VM). It contains Debian GCC 12.2 GNU ARMv7/AArch64 linkers/sysroots;
its setup is recorded in `packaging/evidence/linkers.Dockerfile` from DB-011.
The exact wrapper is [local-cross-linker.sh](evidence/local-cross-linker.sh).
It runs with the developer UID/GID, networking disabled, and a read-only Rust
sysroot. The repository and explicitly created `/private/tmp/db006-rust-*`
linker temporaries are mounted at matching absolute paths.

```sh
rustup target add --toolchain 1.89.0 armv7-unknown-linux-gnueabihf aarch64-unknown-linux-gnu
# Symlink both compiler names in a temporary PATH directory to local-cross-linker.sh.
# TMPDIR points to the explicitly shared /private/tmp/db006-rust-* directory.
PATH="$task_tools:$PATH" TMPDIR="$task_linktmp" RUSTUP_TOOLCHAIN=1.89.0 bash ci/check.sh arm armv7
PATH="$task_tools:$PATH" TMPDIR="$task_linktmp" RUSTUP_TOOLCHAIN=1.89.0 bash ci/check.sh arm arm64
```

Both Rust standard-library target installations exited 0. Local ShellCheck is a
PATH wrapper around this actual command (CI uses the hosted runner's ShellCheck):

```sh
docker run --rm --mount "type=bind,source=$PWD,target=/mnt,readonly" --workdir /mnt koalaman/shellcheck:v0.10.0 ci/check.sh
```

An initial packaging run was invalidated by editing the running shell script: its
race test/build exited 0 but the shell later failed on an unrelated line after
reading the changed file. That attempt is not used as final evidence. All modes
were restarted together after the executable scripts stopped changing. Earlier
local checks also found that a failed rerun could retain an old success manifest;
the script now removes it up front and the failure regression seeds/checks this
case explicitly.

## Negative cases and review

The four CI gate tests use disposable repositories. They verify:

- Unsupported/path-traversing module/target selectors and surplus arguments exit 2
  before creating output directories.
- A subprocess exiting 23 remains exit 23 through `tee`, stderr and the exit code
  are captured, source metadata remains available, and an old success manifest is
  removed.
- Actual unformatted Go source causes exit 1 before vet/tests instead of merely
  printing a filename and passing.
- Actual malformed Go source preserves gofmt's parse-error exit 2 and stops before
  vet. This covers shell conditional contexts where `set -e` alone is insufficient.

Existing Go/Rust application tests cover malformed config, fail-closed transport,
HTTP limits/shutdown and installer safety. No synthetic test claims connectivity
or authentication. Manual source review covers permissions, credentials, events,
command injection, output provenance and module boundaries. Independent/human
review remains outstanding. A targeted added-file secret pattern search is
recorded in `evidence/source-checks.txt`; no dedicated secret scanner or dependency
vulnerability audit is claimed.

## Hosted execution

Both initial hosted runs completed successfully for delivered source commit
`a4ed46c0ee0761a1e732ee610d04dfa4b764db0b`:

- [Push run 41911341485](https://github.com/LiiWzz/DeviceBridge/actions/runs/41911341485): success.
- [Draft PR #18 run 41911387490](https://github.com/LiiWzz/DeviceBridge/actions/runs/41911387490): all eight jobs and artifact uploads succeeded.

The PR checked out synthetic merge commit
`f91559e3d0290b12a8cdb8ebf1ce53f6d47156e6` on hosted Linux x86_64. All eight downloaded
ZIPs matched their GitHub-reported SHA-256 digest. All seven build manifests
matched every listed binary/source file. Every artifact's source metadata matched
the PR run, merge commit and clean workspace. The workflow-only artifact correctly
contains logs/source metadata without binaries or a checksum manifest. Go/ARM ZIPs
preserve a job subdirectory, so verification locates SOURCE.txt before reading
relative manifest paths; the initial assumption that every ZIP was flat was
corrected before the final verification. No downloaded ARM binaries were executed.

Actual command sections from all eight jobs are captured in `evidence/hosted-*.txt`
(timestamp prefixes removed and tabs rendered as literal \t). API run/job/step
conclusions and artifact IDs/digests/expiry are in `evidence/hosted-results.json`.
ZIP, manifest and provenance assertions are recorded in
`evidence/hosted-artifact-verification.json`. The exact executed verifier is
`evidence/verify-hosted-artifacts.py`; it uses Python standard-library urllib,
zipfile and hashlib, reads archives without extraction/execution, and never saves
signed download URLs in its results. Invocation:

```sh
python3 ci/evidence/verify-hosted-artifacts.py '<connector artifact metadata JSON with current download URLs>'
```

All eight artifact verifications exited 0. Signed download URLs were supplied by
GitHub connector tools at execution time and are not committed. These are
checksums/provenance consistency checks, not independent publisher attestation.
This evidence update changes only docs/reports; its resulting PR head checks are
reported in the PR after they complete.

GitHub CLI has no authentication in this environment. The first Git push was
rejected because the local OAuth App lacks workflow scope. The authorized GitHub
connector successfully created the same complete source tree and feature branch;
its tree SHA `0eea1aea6f679baff1a989a420e6a69c435b8eb5` matched the tested local commit
exactly. Local branch history was aligned only after confirming identical trees
and a clean worktree. No token, repository protection or protected branch was
changed. Draft PR creation and hosted inspection used the same connector.

## Security, performance, platform limits and rollback

`contents: read` only, credential persistence disabled, hosted runners, no injected
repository secrets or privileged PR events, no shared caches or release/deployment
rights. Major-pinned official actions and versioned actionlint are build-host
supply-chain dependencies; hosted images/apt packages remain mutable. No tokens,
private keys, production configuration, application dependency or runtime change
is introduced. Branch protections/merge permissions are not changed.

Artifacts are temporary development outputs, not signed production releases.
Cross-built Go binaries are static ELF ARMv7/AArch64. Rust binaries are GNU dynamic
ELF with matching ARM loaders, not universal Buildroot/Yocto packages. No ARM
binary is executed by this issue's cross-build jobs. No physical 128 MB device,
Linux ARM runtime/RSS/CPU/restarts, musl or Windows toolchain validation was
performed. The macOS development host and Linux linker VM do not satisfy target
hardware acceptance. CI adds no embedded allocation, flash writes, root
requirement or runtime dependencies; hosted CPU/storage costs are unmeasured.

Rollback is a reviewed revert of DB-006 commits on a feature branch; no deployed
services, persistent device state or databases require changes. Generated local
`dist/ci/` outputs can be removed deliberately. GitHub artifacts expire after 14
days. See [README.md](README.md) for full reproduction and security/resource
assessment. Draft PR targets `develop`; human Tech Lead alone approves/manual
merges. No self-approval, merge, auto-merge or direct protected-branch push.
