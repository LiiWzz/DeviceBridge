#!/usr/bin/env bash
# Build-host tooling only. Invoke with bash from any working directory.
set -euo pipefail

root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"
mode=${1:-}
case "$mode:$#" in
  workflow:1|rust:1) label=$mode ;;
  go:2)
    case "$2" in
      server/control|server/relay|cli|packaging) module=$2 ;;
      *) echo 'unsupported Go module' >&2; exit 2 ;;
    esac
    label=go-${module//\//-}
    ;;
  arm:2)
    case "$2" in
      armv7) arch=arm; target=armv7-unknown-linux-gnueabihf; compiler=arm-linux-gnueabihf-gcc ;;
      arm64) arch=arm64; target=aarch64-unknown-linux-gnu; compiler=aarch64-linux-gnu-gcc ;;
      *) echo 'unsupported ARM target' >&2; exit 2 ;;
    esac
    label=arm-$2
    ;;
  *) echo 'usage: bash ci/check.sh workflow|rust|go MODULE|arm armv7|arm64' >&2; exit 2 ;;
esac

export GOTOOLCHAIN=local GOWORK=off
export RUSTUP_TOOLCHAIN=${RUSTUP_TOOLCHAIN:-1.89.0}
out=$root/dist/ci/$label
mkdir -p "$out/bin"
# A failed rerun must not retain a previous success manifest.
rm -f "$out/SHA256SUMS"
log=$out/commands.txt
: > "$log"
{
  printf 'checked_out_commit=%s\n' "$(git rev-parse HEAD)"
  printf 'event_commit=%s\n' "${GITHUB_SHA:-local}"
  printf 'run_id=%s\nrun_attempt=%s\nevent=%s\n' "${GITHUB_RUN_ID:-local}" "${GITHUB_RUN_ATTEMPT:-local}" "${GITHUB_EVENT_NAME:-local}"
  printf 'job=%s\nrunner=%s\nrust_toolchain=%s\n' "$label" "$(uname -srm)" "$RUSTUP_TOOLCHAIN"
  if [[ -n $(git status --porcelain) ]]; then
    printf 'workspace_dirty=true\n'
  else
    printf 'workspace_dirty=false\n'
  fi
} > "$out/SOURCE.txt"

run() {
  local status
  printf '[cwd %s]\n$' "$PWD" | tee -a "$log"
  printf ' %q' "$@" | tee -a "$log"
  printf '\n' | tee -a "$log"
  if "$@" 2>&1 | tee -a "$log"; then
    status=0
  else
    status=${PIPESTATUS[0]}
    # A failed log write is also a failed check.
    if [[ $status == 0 ]]; then status=1; fi
  fi
  printf '[exit %s]\n' "$status" | tee -a "$log"
  return "$status"
}

check_go_format() {
  local unformatted
  unformatted=$(gofmt -l .) || return "$?"
  if [[ -n $unformatted ]]; then
    printf '%s\n' "$unformatted"
    return 1
  fi
}

checksums() {
  cd "$out"
  if command -v sha256sum >/dev/null 2>&1; then
    run sha256sum bin/* SOURCE.txt
    sha256sum bin/* SOURCE.txt > SHA256SUMS
  else
    run shasum -a 256 bin/* SOURCE.txt
    shasum -a 256 bin/* SOURCE.txt > SHA256SUMS
  fi
}

case "$mode" in
  workflow)
    run go version
    run go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 -shellcheck= .github/workflows/ci.yml
    run bash -n ci/check.sh ci/security.sh tests/smoke.sh
    run shellcheck ci/check.sh ci/security.sh tests/smoke.sh
    run python3 ci/test_check.py
    ;;
  go)
    cd "$root/$module"
    run go version
    run go env GOOS GOARCH CGO_ENABLED GOTOOLCHAIN GOWORK
    run check_go_format
    run go vet ./...
    run go test -count=1 -v ./...
    run go test -race -count=1 ./...
    run go build -trimpath -o "$out/bin/" ./...
    checksums
    ;;
  rust)
    cd "$root/agent"
    run rustc --version
    run cargo --version
    run cargo fmt --all -- --check
    run cargo clippy --locked --all-targets -- -D warnings
    run cargo test --locked
    run cargo build --release --locked
    run cp target/release/devicebridge-agent "$out/bin/"
    checksums
    ;;
  arm)
    run go version
    run rustc --version
    run cargo --version
    run "$compiler" --version
    export CGO_ENABLED=0 GOOS=linux GOARCH=$arch
    if [[ $arch == arm ]]; then export GOARM=7; fi
    run go env GOOS GOARCH GOARM CGO_ENABLED GOTOOLCHAIN GOWORK
    for module in server/control server/relay cli; do
      cd "$root/$module"
      run go build -trimpath -o "$out/bin/" ./...
    done
    cd "$root/agent"
    run cargo build --release --locked --target "$target"
    run cp "target/$target/release/devicebridge-agent" "$out/bin/"
    run file "$out"/bin/*
    checksums
    ;;
esac
