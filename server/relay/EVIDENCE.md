# DB-004 (#4) executed evidence

Captured UTC: 2026-11-03T20:47:26.299234+00:00

Module: `/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-004-relay-cli/server/relay`. Commands below ran in this directory. Exit codes and stdout/stderr are captured from actual subprocesses; empty output is shown explicitly and tab characters are rendered as \t.

## Scope and approval disclosure

Only `server/relay/` and `cli/` are owned by DB-004. Two independent Go modules, no shared root module or changed contracts. DB-001 PR #11 merge and human ADR approval were supplied by the coordinator; ADR-001 still says **Proposed** in tracked text. That stale status is disclosed here and was not edited. Feature is foundation only.

Coordinator reports auto-merge disabled and main/develop unprotected; no protection settings were changed. No push or PR was performed by this worker. Independent review and a Draft PR into develop citing #4 remain coordinator work; human Tech Lead approval and manual merge are required.

## Toolchain and checks

~~~text
$ go version
go version go1.23.4 darwin/arm64
[exit 0; expected 0]
~~~

~~~text
$ go env GOOS GOARCH CGO_ENABLED GOTOOLCHAIN
darwin
arm64
1
auto
[exit 0; expected 0]
~~~

~~~text
$ uname -srm
Darwin 25.6.0 arm64
[exit 0; expected 0]
~~~

~~~text
$ python3 --version
Python 3.14.5
[exit 0; expected 0]
~~~

~~~text
$ git branch --show-current
feat/m01-db-004-relay-cli
[exit 0; expected 0]
~~~

~~~text
$ git rev-parse HEAD
c8fae1481fe657f385d7747878780d1f4485e986
[exit 0; expected 0]
~~~

~~~text
$ gofmt -w cmd internal
(no output)
[exit 0; expected 0]
~~~

~~~text
$ gofmt -l cmd internal
(no output)
[exit 0; expected 0]
~~~

~~~text
$ go vet ./...
(no output)
[exit 0; expected 0]
~~~

~~~text
$ go test ./...
ok  \tgithub.com/LiiWzz/DeviceBridge/server/relay/cmd/devicebridge-relay\t0.956s
ok  \tgithub.com/LiiWzz/DeviceBridge/server/relay/internal/config\t1.126s
ok  \tgithub.com/LiiWzz/DeviceBridge/server/relay/internal/transport\t1.312s
[exit 0; expected 0]
~~~

~~~text
$ go test -race ./...
ok  \tgithub.com/LiiWzz/DeviceBridge/server/relay/cmd/devicebridge-relay\t2.508s
ok  \tgithub.com/LiiWzz/DeviceBridge/server/relay/internal/config\t2.098s
ok  \tgithub.com/LiiWzz/DeviceBridge/server/relay/internal/transport\t2.917s
[exit 0; expected 0]
~~~

~~~text
$ go build -trimpath -o bin/devicebridge-relay ./cmd/devicebridge-relay
(no output)
[exit 0; expected 0]
~~~

~~~text
$ env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -o bin/devicebridge-relay-linux-armv7 ./cmd/devicebridge-relay
(no output)
[exit 0; expected 0]
~~~

~~~text
$ env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o bin/devicebridge-relay-linux-arm64 ./cmd/devicebridge-relay
(no output)
[exit 0; expected 0]
~~~

~~~text
$ file bin/devicebridge-relay bin/devicebridge-relay-linux-armv7 bin/devicebridge-relay-linux-arm64
bin/devicebridge-relay:             Mach-O 64-bit executable arm64
bin/devicebridge-relay-linux-armv7: ELF 32-bit LSB executable, ARM, EABI5 version 1 (SYSV), statically linked, Go BuildID=kKKHJaNEK-9CowJmZ3Un/eC_KHhnL8__w-eHysotK/E-zmnnJEk34d6JqpHuBr/2yWDoc0gSXBb13W5ILCG, with debug_info, not stripped
bin/devicebridge-relay-linux-arm64: ELF 64-bit LSB executable, ARM aarch64, version 1 (SYSV), statically linked, Go BuildID=0Zt0wJk_FrbcOaMlXrLh/jgtmRTj1JtOdMRxzq3KI/VUFrUKkq07afCEiObxR7/sMVb4aS9CikF4SgZrqqB, with debug_info, not stripped
[exit 0; expected 0]
~~~

~~~text
$ ./bin/devicebridge-relay --help
Usage: devicebridge-relay [flags]
M01 idle skeleton: no listeners; transport unimplemented.
  -check
    \tvalidate settings and exit without starting
  -max-buffer-bytes int
    \treserved future buffer ceiling (1..1048576 bytes) (default 32768)
  -max-streams int
    \treserved future stream ceiling (1..1024) (default 16)
  -transport string
    \ttransport (disabled only; unimplemented) (default "disabled")
  -version
    \tprint version and exit
[exit 0; expected 0]
~~~

~~~text
$ ./bin/devicebridge-relay --version
devicebridge-relay 0.1.0-dev
[exit 0; expected 0]
~~~

~~~text
$ ./bin/devicebridge-relay --check
configuration valid; transport unimplemented and disabled
[exit 0; expected 0]
~~~

~~~text
$ ./bin/devicebridge-relay --check --transport tls
transport is unimplemented; only disabled is supported
[exit 2; expected 2]
~~~

~~~text
$ python3 - <<'PY'
import selectors
import signal
import subprocess

for sig in (signal.SIGTERM, signal.SIGINT):
    proc = subprocess.Popen(["./bin/devicebridge-relay"], stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE, text=True)
    try:
        with selectors.DefaultSelector() as ready:
            ready.register(proc.stdout, selectors.EVENT_READ)
            assert ready.select(timeout=3), "startup timeout"
        first = proc.stdout.readline()
        assert first == "relay idle; no listeners; transport unimplemented and disabled\n", first
        proc.send_signal(sig)
        stdout, stderr = proc.communicate(timeout=3)
        assert proc.returncode == 0, proc.returncode
        assert stdout == "relay stopped\n", stdout
        assert stderr == "", stderr
        print(sig.name + ": startup observed; stopped; exit=0; stderr empty")
    finally:
        if proc.poll() is None:
            proc.kill()
            proc.communicate(timeout=3)
PY
SIGTERM: startup observed; stopped; exit=0; stderr empty
SIGINT: startup observed; stopped; exit=0; stderr empty
[exit 0; expected 0]
~~~

## Negative cases and self-review

Tests cover disabled transport refusal even with cancelled context; zero, negative, and excessive stream/buffer bounds; malformed integers; unknown listener flags; positional commands; version/config interactions; idle cancellation and already-cancelled startup. Built process smoke verifies SIGTERM and SIGINT after startup, clean exit, and empty stderr. Self-review checked path ownership, no external dependencies, no shell execution or socket APIs, no credential/config files, and failure paths. No outstanding implementation findings were identified; independent review remains required.

## Security, resources, and limitations

Relay opens no listener and allocates no transport streams/buffers. Disabled.Open always returns a nil stream plus ErrUnimplemented. Resource flags are provisional ceilings, not implemented production bounds; no framing/authentication contract is introduced. Idle blocks on context cancellation without polling or periodic work. Both modules use only Go's standard library, require no root privileges, have no runtime package dependency or application filesystem writes, and do not write flash. Startup parsing/Go runtime allocations exist; no numeric RAM, CPU, or latency claim is made. Native build/test tools write only development artifacts/caches. No credentials, private keys, device tokens, production configuration, insecure auth option, or public forwarding were introduced.

Linux ARMv7 and ARM64 binaries were cross-built, not executed. No real ARM hardware, 128 MB RAM measurement, Buildroot/Yocto/Debian runtime validation, reconnect metrics, P50/P95/P99, or production SSH/TLS security test was available/performed. Windows and other unlisted platforms/toolchains were not validated. Go 1.23.4 darwin/arm64 is the tested toolchain; current toolchain support/security lifecycle was not researched or certified. No production readiness is claimed.

## Rollback

Stop the idle relay with SIGINT or SIGTERM if running. Remove the optional local `bin/` build outputs and revert the DB-004 feature commit on a feature branch (manual reviewed process for any later merge). No service registration, database migration, persisted state, production settings, or credentials require rollback.

## Result

All recorded checks met their expected exit codes and formatting output was empty.

## Final packaging review

The first staged whitespace check flagged literal space-before-tab indentation in captured Go flag help output. Evidence now renders tabs as \t to preserve the output without introducing Git whitespace errors; source code was unaffected. No dedicated secret scanner (gitleaks/trufflehog) was installed; the targeted pattern scan below supplements manual review and does not certify absence of every possible secret.

~~~text
$ go list -m all
github.com/LiiWzz/DeviceBridge/server/relay
[exit 0; expected 0]
~~~

~~~text
$ rg -n --no-ignore 'BEGIN .*PRIVATE KEY|AKIA[0-9A-Z]{16}|gh[pousr]_[A-Za-z0-9]{20,}' -g '*.go' -g '*.mod' .
(no output)
[exit 1; expected 1]
~~~
