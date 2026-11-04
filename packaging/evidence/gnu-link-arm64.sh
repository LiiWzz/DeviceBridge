#!/bin/sh
set -eu
agent_root=$(git rev-parse --show-toplevel)
rust_root=$(rustc --print sysroot)
exec docker run --rm --network none --user "$(id -u):$(id -g)" \
  --mount "type=bind,source=$agent_root,target=$agent_root" \
  --mount "type=bind,source=$rust_root,target=$rust_root,readonly" \
  --mount 'type=bind,source=/private/tmp/db011-rust-tmp,target=/private/tmp/db011-rust-tmp' \
  --workdir "$PWD" --entrypoint aarch64-linux-gnu-gcc \
  devicebridge-m01-linkers:bookworm "$@"
