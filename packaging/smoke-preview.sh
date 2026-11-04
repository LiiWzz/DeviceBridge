#!/bin/sh
# Executes trusted preview artifacts only, in an isolated test/container prefix.
# Requires BusyBox wget for local probes; not a device runtime dependency.
set -eu
package_dir=${1:?Usage: smoke-preview.sh PACKAGE_DIR TARGET VERSION}
package_target=${2:?target required}
package_version=${3:?version required}
[ "$#" -eq 3 ] || exit 2
case "$package_target" in linux-arm64|linux-armv7|linux-amd64) ;; *) exit 2 ;; esac
case "$package_version" in ''|*[!A-Za-z0-9.-]*|*..*|[.-]*) exit 2 ;; esac
test_root=$(mktemp -d)
child=
cleanup() {
    if [ -n "$child" ]; then kill -KILL "$child" 2>/dev/null || :; wait "$child" 2>/dev/null || :; fi
    rm -rf "$test_root"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' HUP TERM
stop_child() {
    kill -TERM "$child"
    attempt=0
    while kill -0 "$child" 2>/dev/null; do
        attempt=$((attempt + 1))
        [ "$attempt" -lt 100 ] || { echo 'shutdown deadline exceeded' >&2; exit 1; }
        sleep 0.1
    done
    wait "$child"
    child=
}
wait_log() {
    attempt=0
    until grep -q "$1" "$2"; do
        kill -0 "$child" 2>/dev/null || { cat "$2"; exit 1; }
        attempt=$((attempt + 1))
        [ "$attempt" -lt 100 ] || exit 1
        sleep 0.1
    done
}
printf 'Host: '; uname -sm
printf 'UID: '; id -u
for component in agent control relay cli; do
    archive=devicebridge-$component-$package_version-$package_target
    (cd "$package_dir" && sha256sum -c "$archive.tar.gz.sha256")
    tar -xzf "$package_dir/$archive.tar.gz" -C "$test_root"
    prefix=$test_root/installed-$component
    sh "$test_root/$archive/install.sh" --prefix "$prefix"
    sh "$test_root/$archive/install.sh" --prefix "$prefix"
    sh "$prefix/manage.sh" --status
    case "$component" in
        cli)
            "$prefix/bin/devicebridge" version
            "$prefix/bin/devicebridge" check
            if "$prefix/bin/devicebridge" check --relay-url http://localhost; then exit 1; fi
            ;;
        agent)
            "$prefix/bin/devicebridge-agent" --version
            "$prefix/bin/devicebridge-agent" --check-config
            "$prefix/bin/devicebridge-agent" > "$test_root/agent.log" 2>&1 &
            child=$!
            wait_log agent_started "$test_root/agent.log"
            stop_child
            grep agent_stopped "$test_root/agent.log"
            ;;
        control)
            DEVICEBRIDGE_CONTROL_ADDR=127.0.0.1:18081 "$prefix/bin/devicebridge-control" > "$test_root/control.log" 2>&1 &
            child=$!
            wait_log 'control API listening' "$test_root/control.log"
            response=$(busybox wget -qO- http://127.0.0.1:18081/healthz)
            [ "$response" = '{"status":"ok"}' ]
            response=$(busybox wget -qO- http://127.0.0.1:18081/readyz)
            [ "$response" = '{"status":"ready"}' ]
            printf 'Control: health=ok ready=ready\n'
            stop_child
            grep 'control API stopped' "$test_root/control.log"
            ;;
        relay)
            "$prefix/bin/devicebridge-relay" --version
            "$prefix/bin/devicebridge-relay" --check
            "$prefix/bin/devicebridge-relay" > "$test_root/relay.log" 2>&1 &
            child=$!
            wait_log 'relay idle' "$test_root/relay.log"
            stop_child
            grep 'relay stopped' "$test_root/relay.log"
            ;;
    esac
    sh "$prefix/manage.sh" --uninstall
    [ ! -e "$prefix" ]
done
echo 'PASS: trusted preview install/reinstall/status/run/negative/shutdown/uninstall'
