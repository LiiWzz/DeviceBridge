# DB-004 (#4) executed evidence

Captured UTC: 2026-11-03T20:47:26.299227+00:00

Module: `/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-004-relay-cli/cli`. Commands below ran in this directory. Exit codes and stdout/stderr are captured from actual subprocesses; empty output is shown explicitly and tab characters are rendered as \t.

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
ok  \tgithub.com/LiiWzz/DeviceBridge/cli/cmd/devicebridge\t0.412s
ok  \tgithub.com/LiiWzz/DeviceBridge/cli/internal/config\t0.801s
[exit 0; expected 0]
~~~

~~~text
$ go test -race ./...
ok  \tgithub.com/LiiWzz/DeviceBridge/cli/cmd/devicebridge\t1.539s
ok  \tgithub.com/LiiWzz/DeviceBridge/cli/internal/config\t1.969s
[exit 0; expected 0]
~~~

~~~text
$ go build -trimpath -o bin/devicebridge ./cmd/devicebridge
(no output)
[exit 0; expected 0]
~~~

~~~text
$ env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -o bin/devicebridge-linux-armv7 ./cmd/devicebridge
(no output)
[exit 0; expected 0]
~~~

~~~text
$ env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o bin/devicebridge-linux-arm64 ./cmd/devicebridge
(no output)
[exit 0; expected 0]
~~~

~~~text
$ file bin/devicebridge bin/devicebridge-linux-armv7 bin/devicebridge-linux-arm64
bin/devicebridge:             Mach-O 64-bit executable arm64
bin/devicebridge-linux-armv7: ELF 32-bit LSB executable, ARM, EABI5 version 1 (SYSV), statically linked, Go BuildID=kpw19nMn-ci9faTKS11y/Bvtsfi4Qt3JROcqD2odS/NCrVxbIsGLCwwXOCOtMr/bFKvrnq37iXBtWPBIS1S, with debug_info, not stripped
bin/devicebridge-linux-arm64: ELF 64-bit LSB executable, ARM aarch64, version 1 (SYSV), statically linked, Go BuildID=Jnt5TPdum2_o9pSi_Tbz/dEsIsSjQtGQX4mDalBjB/6UjzjEV1LcoZnAR_gaVa/aJqbOOgwi8sZGiLc_tPN, with debug_info, not stripped
[exit 0; expected 0]
~~~

~~~text
$ ./bin/devicebridge --help
Usage: devicebridge <help|version|check>
M01 developer CLI: transport unimplemented; no remote operations.
  check [--relay-url https://localhost:443] [--timeout 5s]
    Validate local settings only; does not contact a server.
[exit 0; expected 0]
~~~

~~~text
$ ./bin/devicebridge version
devicebridge 0.1.0-dev
[exit 0; expected 0]
~~~

~~~text
$ ./bin/devicebridge check
configuration valid; transport unimplemented; no connection attempted
[exit 0; expected 0]
~~~

~~~text
$ ./bin/devicebridge check --relay-url https://relay.example:8443 --timeout 1s
configuration valid; transport unimplemented; no connection attempted
[exit 0; expected 0]
~~~

~~~text
$ ./bin/devicebridge check --relay-url http://relay.example
relay-url must be an absolute HTTPS URL
[exit 2; expected 2]
~~~

~~~text
$ ./bin/devicebridge check --relay-url https://user:example-not-a-credential@relay.example
relay-url must not contain credentials, query, fragment, or a non-root path
[exit 2; expected 2]
~~~

~~~text
$ ./bin/devicebridge check --timeout 0s
timeout must be greater than zero and at most 1m
[exit 2; expected 2]
~~~

~~~text
$ ./bin/devicebridge check --insecure
flag provided but not defined: -insecure
Usage of check:
  -relay-url string
    \tfuture relay HTTPS endpoint; never contacted (default "https://localhost:443")
  -timeout duration
    \treserved future request timeout (0 < timeout <= 1m) (default 5s)
[exit 2; expected 2]
~~~

~~~text
$ ./bin/devicebridge connect device
command unavailable; transport and remote operations are unimplemented (use help)
[exit 2; expected 2]
~~~

~~~text
$ ./bin/devicebridge shell device id
command unavailable; transport and remote operations are unimplemented (use help)
[exit 2; expected 2]
~~~

## Negative cases and self-review

Tests cover help/version/check and unsupported connect/shell; HTTPS-only endpoints; credential/query/fragment/path rejection; malformed ports, DNS/IP hosts and escapes; zero/negative/excessive timeouts; unknown flags and extra arguments. Credential URL validation does not echo the rejected value. Self-review identified net/url alone accepting invalid bracketed hosts; added netip and ASCII DNS validation plus negative tests before final verification. Self-review checked path ownership, no external dependencies, no shell execution or socket APIs, no credential/config files, and failure paths. No outstanding implementation findings were identified; independent review remains required.

## Security, resources, and limitations

CLI makes no DNS lookup or network connection, executes no remote/local shell commands, and rejects all unavailable operations. Endpoint/timeout settings are local validation only, not a connectivity or authentication result. Both modules use only Go's standard library, require no root privileges, have no runtime package dependency or application filesystem writes, and do not write flash. Startup parsing/Go runtime allocations exist; no numeric RAM, CPU, or latency claim is made. Native build/test tools write only development artifacts/caches. No credentials, private keys, device tokens, production configuration, insecure auth option, or public forwarding were introduced.

Linux ARMv7 and ARM64 binaries were cross-built, not executed. No real ARM hardware, 128 MB RAM measurement, Buildroot/Yocto/Debian runtime validation, reconnect metrics, P50/P95/P99, or production SSH/TLS security test was available/performed. Windows and other unlisted platforms/toolchains were not validated. Go 1.23.4 darwin/arm64 is the tested toolchain; current toolchain support/security lifecycle was not researched or certified. No production readiness is claimed.

## Rollback

Stop the idle relay with SIGINT or SIGTERM if running. Remove the optional local `bin/` build outputs and revert the DB-004 feature commit on a feature branch (manual reviewed process for any later merge). No service registration, database migration, persisted state, production settings, or credentials require rollback.

## Result

All recorded checks met their expected exit codes and formatting output was empty.

## Final packaging review

The first staged whitespace check flagged literal space-before-tab indentation in captured Go flag help output. Evidence now renders tabs as \t to preserve the output without introducing Git whitespace errors; source code was unaffected. No dedicated secret scanner (gitleaks/trufflehog) was installed; the targeted pattern scan below supplements manual review and does not certify absence of every possible secret.

~~~text
$ go list -m all
github.com/LiiWzz/DeviceBridge/cli
[exit 0; expected 0]
~~~

~~~text
$ rg -n --no-ignore 'BEGIN .*PRIVATE KEY|AKIA[0-9A-Z]{16}|gh[pousr]_[A-Za-z0-9]{20,}' -g '*.go' -g '*.mod' .
(no output)
[exit 1; expected 1]
~~~
