#!/bin/sh
# Embedded into a trusted offline package by the build-host Go tool.
set -eu
set -f
umask 077
component='@COMPONENT@'
version='@VERSION@'
target='@TARGET@'
binary='@BINARY@'
payload_sha='@SHA256@'
action=install
prefix=
prefix_seen=no
action_seen=no
stage=
link_tmp=
locked=no
remove_root=no

die() { printf 'DeviceBridge: %s\n' "$*" >&2; exit 1; }
usage() {
    cat <<EOF
DeviceBridge $component $version ($target) — M01 foundation only
Usage: sh install.sh [--prefix ABSOLUTE_PATH]
       sh PREFIX/manage.sh [--status|--rollback|--uninstall] [--prefix PATH]
Default install prefix: \$HOME/.local/devicebridge/$component
Offline; no root, downloads, service enablement, configuration or PATH changes.
SHA-256 detects corruption; obtain the entire package from a trusted publisher.
EOF
}
while [ "$#" -gt 0 ]; do
    case "$1" in
        --help|-h) [ "$#" -eq 1 ] || die 'help must be used alone'; usage; exit 0 ;;
        --prefix)
            if [ "$prefix_seen" != no ] || [ "$#" -lt 2 ]; then die 'prefix requires one value and cannot be repeated'; fi
            [ -n "$2" ] || die 'prefix cannot be empty'
            prefix=$2; prefix_seen=yes; shift 2 ;;
        --status|--rollback|--uninstall)
            [ "$action_seen" = no ] || die 'choose only one action'
            action=${1#--}; action_seen=yes; shift ;;
        *) die 'unknown argument; use --help' ;;
    esac
done
case "$component:$binary" in
    agent:devicebridge-agent|control:devicebridge-control|relay:devicebridge-relay|cli:devicebridge) ;;
    *) die 'invalid package metadata' ;;
esac
case "$version" in ''|*[!A-Za-z0-9.-]*|*..*|[.-]*) die 'invalid version' ;; esac
[ "${#version}" -le 64 ] || die 'invalid version length'
case "$target" in linux-armv7|linux-arm64|linux-amd64|darwin-arm64|darwin-amd64) ;; *) die 'invalid package target' ;; esac
[ "${#payload_sha}" -eq 64 ] || die 'invalid package checksum'
case "$payload_sha" in *[!0-9a-f]*) die 'invalid package checksum' ;; esac
for utility in uname dirname mkdir cp mv ln readlink chmod rm rmdir cat wc mktemp; do
    command -v "$utility" >/dev/null 2>&1 || die "missing utility: $utility"
done
if command -v sha256sum >/dev/null 2>&1; then hash_tool=sha256sum
elif command -v shasum >/dev/null 2>&1; then hash_tool=shasum
elif command -v openssl >/dev/null 2>&1; then hash_tool=openssl
else die 'SHA-256 utility required: sha256sum, shasum, or openssl'
fi
hash_file() {
    case "$hash_tool" in
        sha256sum) hash_output=$(sha256sum < "$1") || die 'cannot hash payload'; hash_output=${hash_output%% *} ;;
        shasum) hash_output=$(shasum -a 256 < "$1") || die 'cannot hash payload'; hash_output=${hash_output%% *} ;;
        openssl) hash_output=$(openssl dgst -sha256 < "$1") || die 'cannot hash payload'; hash_output=${hash_output##* } ;;
    esac
    printf '%s\n' "$hash_output"
}
regular() { [ -f "$1" ] && [ ! -L "$1" ]; }
source_dir=$(CDPATH='' cd -P "$(dirname "$0")" && pwd -P) || die 'cannot locate installer'
if [ -z "$prefix" ]; then
    if [ "${0##*/}" = manage.sh ]; then prefix=$source_dir
    else
        [ -n "${HOME:-}" ] || die 'HOME is unset; specify --prefix'
        home_dir=$(CDPATH='' cd -P "$HOME" && pwd -P) || die 'HOME must be an existing directory'
        prefix=$home_dir/.local/devicebridge/$component
    fi
fi
case "$prefix" in /*) ;; *) die 'prefix must be absolute' ;; esac
case "$prefix" in /|*'/../'*|*/..|*'/./'*|*/.|*'//'*|*/|*'
'*) die 'prefix must be a normalized absolute path without newlines' ;; esac
# Never follow a symlink in a destination path. Dedicated prefixes must have
# trusted parents; this unprivileged tool is not a boundary against the same UID.
saved_ifs=$IFS; IFS=/; current=
for part in $prefix; do
    [ -n "$part" ] || continue
    current=$current/$part
    [ ! -L "$current" ] || die 'symlink in installation path'
    if [ -e "$current" ] && [ ! -d "$current" ]; then die 'non-directory in installation path'; fi
done
IFS=$saved_ifs
set +f
marker="devicebridge-install-v1:$component:$target"
check_marker() {
    regular "$prefix/.devicebridge-owner" || die 'prefix exists but is not a managed installation'
    [ "$(wc -c < "$prefix/.devicebridge-owner")" -le 128 ] || die 'invalid owner marker'
    [ "$(cat "$prefix/.devicebridge-owner")" = "$marker" ] || die 'prefix belongs to a different component/target'
}
valid_id() {
    case "$1" in ''|*[!A-Za-z0-9.-]*|*..*|[.-]*) return 1 ;; esac
    [ "${#1}" -le 100 ]
}
check_release() {
    valid_id "$1" || die 'invalid installed release identifier'
    release_path=$prefix/releases/$1
    if [ ! -d "$release_path" ] || [ -L "$release_path" ]; then die 'invalid installed release directory'; fi
    if ! regular "$release_path/$binary" || ! regular "$release_path/payload.sha256"; then die 'invalid installed payload/checksum'; fi
    [ "$(wc -c < "$release_path/payload.sha256")" -eq 65 ] || die 'invalid installed checksum size'
    stored_sha=$(cat "$release_path/payload.sha256")
    case "$stored_sha" in *[!0-9a-f]*) die 'invalid installed checksum' ;; esac
    [ "${#stored_sha}" -eq 64 ] || die 'invalid installed checksum'
    [ "$(hash_file "$release_path/$binary")" = "$stored_sha" ] || die 'installed payload checksum mismatch'
}
active_id=; previous_id=
check_state() {
    for directory in bin releases; do
        if [ ! -d "$prefix/$directory" ] || [ -L "$prefix/$directory" ]; then die 'invalid managed directory'; fi
    done
    if [ -e "$prefix/bin/$binary" ] || [ -L "$prefix/bin/$binary" ]; then
        [ -L "$prefix/bin/$binary" ] || die 'refusing to replace unrelated binary'
        active_link=$(readlink "$prefix/bin/$binary") || die 'cannot read active link'
        case "$active_link" in ../releases/*/"$binary") ;; *) die 'unmanaged active link' ;; esac
        active_id=${active_link#../releases/}; active_id=${active_id%/"$binary"}
        check_release "$active_id"
    fi
    if [ -e "$prefix/previous" ] || [ -L "$prefix/previous" ]; then
        [ -L "$prefix/previous" ] || die 'unmanaged previous path'
        previous_link=$(readlink "$prefix/previous") || die 'cannot read previous link'
        case "$previous_link" in releases/*/"$binary") ;; *) die 'unmanaged previous link' ;; esac
        previous_id=${previous_link#releases/}; previous_id=${previous_id%/"$binary"}
        check_release "$previous_id"
    fi
    regular "$prefix/manage.sh" || die 'invalid installed management script'
}
cleanup() {
    [ -z "$link_tmp" ] || rm -f "$link_tmp"
    [ -z "$stage" ] || rm -rf "$stage"
    if [ "$locked" = yes ]; then rmdir "$prefix/.install-lock" 2>/dev/null || :; fi
    if [ "$remove_root" = yes ]; then rmdir "$prefix" 2>/dev/null || :; fi
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP
atomic_link() {
    candidate=$(dirname "$2")/.devicebridge-next-$$
    if [ -e "$candidate" ] || [ -L "$candidate" ]; then return 1; fi
    ln -s "$1" "$candidate" || return 1
    link_tmp=$candidate
    if ! mv -f "$link_tmp" "$2"; then
        rm -f "$link_tmp"
        link_tmp=
        return 1
    fi
    link_tmp=
}
activate_release() {
    # Updating history first would lose the last rollback target if the binary
    # directory is unwritable. Commit the executable before touching history.
    atomic_link "../releases/$1/$binary" "$prefix/bin/$binary" || die 'activation failed; existing version and rollback target unchanged'
    if [ -n "$active_id" ]; then
        if ! atomic_link "releases/$active_id/$binary" "$prefix/previous"; then
            atomic_link "../releases/$active_id/$binary" "$prefix/bin/$binary" || die 'history update and restoration failed; inspect installation before retrying'
            die 'history update failed; existing version and rollback target restored'
        fi
    fi
}
if [ "$action" = install ]; then
    [ "${0##*/}" = install.sh ] || die 'install/upgrade using install.sh from an extracted package'
    os=$(uname -s); cpu=$(uname -m)
    case "$os:$cpu" in
        Linux:armv7l) host=linux-armv7 ;;
        Linux:aarch64|Linux:arm64) host=linux-arm64 ;;
        Linux:x86_64|Linux:amd64) host=linux-amd64 ;;
        Darwin:arm64) host=darwin-arm64 ;;
        Darwin:x86_64) host=darwin-amd64 ;;
        *) die 'unsupported host OS/CPU' ;;
    esac
    [ "$host" = "$target" ] || die "wrong package for this host; need $host"
    if [ ! -d "$source_dir/bin" ] || [ -L "$source_dir/bin" ]; then die 'invalid package binary directory'; fi
    if ! regular "$source_dir/install.sh" || ! regular "$source_dir/bin/$binary"; then die 'missing or non-regular package payload'; fi
    [ "$(hash_file "$source_dir/bin/$binary")" = "$payload_sha" ] || die 'package checksum mismatch; installation unchanged'
    if [ -e "$prefix" ]; then check_marker
    else
        mkdir -p "$(dirname "$prefix")"
        mkdir "$prefix" || die 'cannot create private installation prefix'
        printf '%s\n' "$marker" > "$prefix/.devicebridge-owner"
        mkdir "$prefix/bin" "$prefix/releases"
        cp "$source_dir/install.sh" "$prefix/manage.sh"
        chmod 700 "$prefix/manage.sh"
    fi
else
    [ -d "$prefix" ] || die 'installation not found'
    check_marker
fi
if [ "$action" != status ]; then
    mkdir "$prefix/.install-lock" 2>/dev/null || die 'installation locked; check for another installer before removing a stale lock'
    locked=yes
fi
check_state
case "$action" in
    status)
        [ -n "$active_id" ] || die 'no active release'
        printf 'Component: %s\nTarget: %s\nActive: %s\nPrevious: %s\nExecutable: %s/bin/%s\n' "$component" "$target" "$active_id" "${previous_id:-none}" "$prefix" "$binary"
        ;;
    install)
        new_id=$version-$target
        if [ -e "$prefix/releases/$new_id" ] || [ -L "$prefix/releases/$new_id" ]; then
            check_release "$new_id"
            [ "$stored_sha" = "$payload_sha" ] || die 'version already exists with a different payload; use a new version'
        else
            stage=$(mktemp -d "$prefix/releases/.stage.XXXXXX") || die 'cannot stage installation'
            cp "$source_dir/bin/$binary" "$stage/$binary"
            chmod 755 "$stage/$binary"
            [ "$(hash_file "$stage/$binary")" = "$payload_sha" ] || die 'payload changed while copying'
            printf '%s\n' "$payload_sha" > "$stage/payload.sha256"
            mv "$stage" "$prefix/releases/$new_id"
            stage=
        fi
        if [ "$active_id" != "$new_id" ]; then
            activate_release "$new_id"
        fi
        printf 'Installed %s %s\nRun: "%s/bin/%s" --help\nManage: sh "%s/manage.sh" --status\n' "$component" "$version" "$prefix" "$binary" "$prefix"
        ;;
    rollback)
        if [ -z "$active_id" ] || [ -z "$previous_id" ] || [ "$active_id" = "$previous_id" ]; then die 'no different previous release available'; fi
        activate_release "$previous_id"
        printf 'Rolled back %s to %s\n' "$component" "$previous_id"
        ;;
    uninstall)
        # Preflight every release before any deletion. Extra files in a release
        # fail closed; unrelated user configuration elsewhere is preserved.
        for release in "$prefix/releases"/* "$prefix/releases"/.[!.]* "$prefix/releases"/..?*; do
            [ -e "$release" ] || [ -L "$release" ] || continue
            check_release "${release##*/}"
            for entry in "$release"/* "$release"/.[!.]* "$release"/..?*; do
                [ -e "$entry" ] || [ -L "$entry" ] || continue
                case "${entry##*/}" in "$binary"|payload.sha256) ;; *) die 'extra release files found; nothing removed' ;; esac
            done
        done
        for release in "$prefix/releases"/*; do
            [ -d "$release" ] || continue
            rm "$release/$binary" "$release/payload.sha256"
            rmdir "$release"
        done
        rm -f "$prefix/bin/$binary" "$prefix/previous" "$prefix/manage.sh" "$prefix/.devicebridge-owner"
        rmdir "$prefix/bin" "$prefix/releases" 2>/dev/null || :
        remove_root=yes
        printf 'Uninstalled %s; unrelated files/configuration preserved. Stop running processes separately.\n' "$component"
        ;;
esac
