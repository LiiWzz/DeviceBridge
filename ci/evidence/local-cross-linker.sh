#!/bin/sh
# Evidence helper for macOS only; CI installs GNU cross-linkers directly.
set -eu
: "${TMPDIR:?Use a shared /private/tmp directory for Rust linker temporaries}"
compiler=$(basename "$0")
agent_root=$(git rev-parse --show-toplevel)
rust_root=$(rustc --print sysroot)
exec docker run --rm --network none --user "$(id -u):$(id -g)" \
  --mount "type=bind,source=$agent_root,target=$agent_root" \
  --mount "type=bind,source=$rust_root,target=$rust_root,readonly" \
  --mount "type=bind,source=$TMPDIR,target=$TMPDIR" \
  --workdir "$PWD" --entrypoint "$compiler" \
  devicebridge-m01-linkers:bookworm "$@"
