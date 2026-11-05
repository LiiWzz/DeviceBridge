#!/usr/bin/env bash
# Build-host smoke only; never installs on devices or uses deploy/.env.
set -euo pipefail
root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"
export GOTOOLCHAIN=local GOWORK=off
export RUSTUP_TOOLCHAIN=${RUSTUP_TOOLCHAIN:-1.89.0}
out=$root/dist/ci/integration
mkdir -p "$out/bin"
rm -f "$out/PASS"
exec > >(tee "$out/commands.txt") 2>&1
finish() {
    local result=$?
    echo "[exit $result]"
    exit "$result"
}
trap finish EXIT
{
    printf 'checked_out_commit=%s\n' "$(git rev-parse HEAD)"
    printf 'event_commit=%s\nrun_id=%s\n' "${GITHUB_SHA:-local}" "${GITHUB_RUN_ID:-local}"
    printf 'runner=%s\n' "$(uname -srm)"
    printf 'workspace_dirty=%s\n' "$(test -z "$(git status --porcelain)" && echo false || echo true)"
} > "$out/SOURCE.txt"
run() { printf '\n$'; printf ' %q' "$@"; printf '\n'; "$@"; }
run go version
run cargo --version
run docker version --format '{{.Server.Version}}'
run docker compose version
run python3 tests/test_process_smoke.py
for module in server/control server/relay cli; do
    (
        cd "$root/$module"
        run go build -trimpath -o "$out/bin/" ./...
    )
done
(
    cd "$root/agent"
    run cargo build --release --locked
)
run cp agent/target/release/devicebridge-agent "$out/bin/"
run python3 tests/process_smoke.py "$out/bin" "$out"
# Reuse the accepted DB-005 migration/health/restart/cleanup smoke unchanged.
# Only a newly generated disposable password enters the child environment.
POSTGRES_PASSWORD=$(python3 -c 'import secrets; print(secrets.token_hex(32))')
export POSTGRES_PASSWORD
run sh deploy/smoke.sh
unset POSTGRES_PASSWORD
printf 'PASS: DB-008 process and isolated Compose smoke\n' > "$out/PASS"
cat "$out/PASS"
