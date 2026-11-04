# DB-003 execution evidence

Scope: issue #3 only, `server/control/` only. Base: `origin/develop` at
`c8fae1481fe657f385d7747878780d1f4485e986`; branch:
`feat/m01-db-003-control-api`. Required README, AGENTS, M01 execution plan and
ADR-001 were read before implementation. DB-001 PR #11 merge and human ADR
approval are supplied by the coordinator; ADR-001 still says **Proposed**.
That stale status is disclosed here and the shared ADR was not edited.

## Executed checks

Commands below ran from `server/control/` on the local macOS ARM64 host.
Empty output is recorded explicitly. Cross-builds do not imply target execution.

Recorded UTC: 2026-11-03T20:46:55.924541+00:00

```text
$ go version
go version go1.23.4 darwin/arm64
exit: 0
```

```text
$ go env GOOS GOARCH CGO_ENABLED
darwin
arm64
1
exit: 0
```

```text
$ uname -a
Darwin Mac.lan 25.6.0 Darwin Kernel Version 25.6.0: Fri Jul 31 19:18:43 PDT 2026; root:xnu-12377.161.14~5/RELEASE_ARM64_T8112 arm64
exit: 0
```

```text
$ python3 --version
Python 3.14.5
exit: 0
```

```text
$ gofmt -w cmd internal
(no output)
exit: 0
```

```text
$ gofmt -l cmd internal
(no output)
exit: 0
```

```text
$ go vet ./...
(no output)
exit: 0
```

```text
$ go test -count=1 -v ./...
?   	github.com/LiiWzz/DeviceBridge/server/control/cmd/control	[no test files]
=== RUN   TestConfigDefaultsAndOverrides
--- PASS: TestConfigDefaultsAndOverrides (0.00s)
=== RUN   TestRejectMalformedConfig
=== RUN   TestRejectMalformedConfig/addr=
=== RUN   TestRejectMalformedConfig/addr=:8080
=== RUN   TestRejectMalformedConfig/addr=0.0.0.0:8080
=== RUN   TestRejectMalformedConfig/addr=[::]:8080
=== RUN   TestRejectMalformedConfig/addr=192.0.2.1:8080
=== RUN   TestRejectMalformedConfig/addr=localhost:8080
=== RUN   TestRejectMalformedConfig/addr=[::1%lo0]:8080
=== RUN   TestRejectMalformedConfig/addr=127.0.0.1
=== RUN   TestRejectMalformedConfig/addr=127.0.0.1:0
=== RUN   TestRejectMalformedConfig/addr=127.0.0.1:80
=== RUN   TestRejectMalformedConfig/addr=127.0.0.1:65536
=== RUN   TestRejectMalformedConfig/addr=127.0.0.1:http
=== RUN   TestRejectMalformedConfig/addr=127.0.0.1:+8080
=== RUN   TestRejectMalformedConfig/addr=127.0.0.1:08080
=== RUN   TestRejectMalformedConfig/addr=127.0.0.1:8080_
=== RUN   TestRejectMalformedConfig/READ_HEADER_TIMEOUT=
=== RUN   TestRejectMalformedConfig/READ_HEADER_TIMEOUT=secret-do-not-log
=== RUN   TestRejectMalformedConfig/READ_HEADER_TIMEOUT=0
=== RUN   TestRejectMalformedConfig/READ_HEADER_TIMEOUT=-1s
=== RUN   TestRejectMalformedConfig/READ_HEADER_TIMEOUT=1ns
=== RUN   TestRejectMalformedConfig/READ_HEADER_TIMEOUT=301s
=== RUN   TestRejectMalformedConfig/READ_HEADER_TIMEOUT=999999999999999999h
=== RUN   TestRejectMalformedConfig/READ_TIMEOUT=
=== RUN   TestRejectMalformedConfig/READ_TIMEOUT=secret-do-not-log
=== RUN   TestRejectMalformedConfig/READ_TIMEOUT=0
=== RUN   TestRejectMalformedConfig/READ_TIMEOUT=-1s
=== RUN   TestRejectMalformedConfig/READ_TIMEOUT=1ns
=== RUN   TestRejectMalformedConfig/READ_TIMEOUT=301s
=== RUN   TestRejectMalformedConfig/READ_TIMEOUT=999999999999999999h
=== RUN   TestRejectMalformedConfig/WRITE_TIMEOUT=
=== RUN   TestRejectMalformedConfig/WRITE_TIMEOUT=secret-do-not-log
=== RUN   TestRejectMalformedConfig/WRITE_TIMEOUT=0
=== RUN   TestRejectMalformedConfig/WRITE_TIMEOUT=-1s
=== RUN   TestRejectMalformedConfig/WRITE_TIMEOUT=1ns
=== RUN   TestRejectMalformedConfig/WRITE_TIMEOUT=301s
=== RUN   TestRejectMalformedConfig/WRITE_TIMEOUT=999999999999999999h
=== RUN   TestRejectMalformedConfig/IDLE_TIMEOUT=
=== RUN   TestRejectMalformedConfig/IDLE_TIMEOUT=secret-do-not-log
=== RUN   TestRejectMalformedConfig/IDLE_TIMEOUT=0
=== RUN   TestRejectMalformedConfig/IDLE_TIMEOUT=-1s
=== RUN   TestRejectMalformedConfig/IDLE_TIMEOUT=1ns
=== RUN   TestRejectMalformedConfig/IDLE_TIMEOUT=301s
=== RUN   TestRejectMalformedConfig/IDLE_TIMEOUT=999999999999999999h
=== RUN   TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=
=== RUN   TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=secret-do-not-log
=== RUN   TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=0
=== RUN   TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=-1s
=== RUN   TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=1ns
=== RUN   TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=301s
=== RUN   TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=999999999999999999h
--- PASS: TestRejectMalformedConfig (0.00s)
    --- PASS: TestRejectMalformedConfig/addr= (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=:8080 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=0.0.0.0:8080 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=[::]:8080 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=192.0.2.1:8080 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=localhost:8080 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=[::1%lo0]:8080 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=127.0.0.1 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=127.0.0.1:0 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=127.0.0.1:80 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=127.0.0.1:65536 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=127.0.0.1:http (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=127.0.0.1:+8080 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=127.0.0.1:08080 (0.00s)
    --- PASS: TestRejectMalformedConfig/addr=127.0.0.1:8080_ (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_HEADER_TIMEOUT= (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_HEADER_TIMEOUT=secret-do-not-log (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_HEADER_TIMEOUT=0 (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_HEADER_TIMEOUT=-1s (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_HEADER_TIMEOUT=1ns (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_HEADER_TIMEOUT=301s (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_HEADER_TIMEOUT=999999999999999999h (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_TIMEOUT= (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_TIMEOUT=secret-do-not-log (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_TIMEOUT=0 (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_TIMEOUT=-1s (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_TIMEOUT=1ns (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_TIMEOUT=301s (0.00s)
    --- PASS: TestRejectMalformedConfig/READ_TIMEOUT=999999999999999999h (0.00s)
    --- PASS: TestRejectMalformedConfig/WRITE_TIMEOUT= (0.00s)
    --- PASS: TestRejectMalformedConfig/WRITE_TIMEOUT=secret-do-not-log (0.00s)
    --- PASS: TestRejectMalformedConfig/WRITE_TIMEOUT=0 (0.00s)
    --- PASS: TestRejectMalformedConfig/WRITE_TIMEOUT=-1s (0.00s)
    --- PASS: TestRejectMalformedConfig/WRITE_TIMEOUT=1ns (0.00s)
    --- PASS: TestRejectMalformedConfig/WRITE_TIMEOUT=301s (0.00s)
    --- PASS: TestRejectMalformedConfig/WRITE_TIMEOUT=999999999999999999h (0.00s)
    --- PASS: TestRejectMalformedConfig/IDLE_TIMEOUT= (0.00s)
    --- PASS: TestRejectMalformedConfig/IDLE_TIMEOUT=secret-do-not-log (0.00s)
    --- PASS: TestRejectMalformedConfig/IDLE_TIMEOUT=0 (0.00s)
    --- PASS: TestRejectMalformedConfig/IDLE_TIMEOUT=-1s (0.00s)
    --- PASS: TestRejectMalformedConfig/IDLE_TIMEOUT=1ns (0.00s)
    --- PASS: TestRejectMalformedConfig/IDLE_TIMEOUT=301s (0.00s)
    --- PASS: TestRejectMalformedConfig/IDLE_TIMEOUT=999999999999999999h (0.00s)
    --- PASS: TestRejectMalformedConfig/SHUTDOWN_TIMEOUT= (0.00s)
    --- PASS: TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=secret-do-not-log (0.00s)
    --- PASS: TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=0 (0.00s)
    --- PASS: TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=-1s (0.00s)
    --- PASS: TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=1ns (0.00s)
    --- PASS: TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=301s (0.00s)
    --- PASS: TestRejectMalformedConfig/SHUTDOWN_TIMEOUT=999999999999999999h (0.00s)
=== RUN   TestProbeRoutes
=== RUN   TestProbeRoutes/GET/healthz{"status":"ok"}_
=== RUN   TestProbeRoutes/HEAD/healthz
=== RUN   TestProbeRoutes/GET/readyz{"status":"ready"}_
=== RUN   TestProbeRoutes/GET/readyz{"status":"shutting_down"}_
=== RUN   TestProbeRoutes/HEAD/readyz
=== RUN   TestProbeRoutes/POST/healthz{"error":"method_not_allowed"}_
=== RUN   TestProbeRoutes/DELETE/readyz{"error":"method_not_allowed"}_
=== RUN   TestProbeRoutes/OPTIONS/healthz{"error":"method_not_allowed"}_
=== RUN   TestProbeRoutes/GET/{"error":"not_found"}_
=== RUN   TestProbeRoutes/GET/debug/pprof/{"error":"not_found"}_
=== RUN   TestProbeRoutes/GET/healthz/{"error":"not_found"}_
=== RUN   TestProbeRoutes/GET/x/../healthz{"error":"not_found"}_
--- PASS: TestProbeRoutes (0.00s)
    --- PASS: TestProbeRoutes/GET/healthz{"status":"ok"}_ (0.00s)
    --- PASS: TestProbeRoutes/HEAD/healthz (0.00s)
    --- PASS: TestProbeRoutes/GET/readyz{"status":"ready"}_ (0.00s)
    --- PASS: TestProbeRoutes/GET/readyz{"status":"shutting_down"}_ (0.00s)
    --- PASS: TestProbeRoutes/HEAD/readyz (0.00s)
    --- PASS: TestProbeRoutes/POST/healthz{"error":"method_not_allowed"}_ (0.00s)
    --- PASS: TestProbeRoutes/DELETE/readyz{"error":"method_not_allowed"}_ (0.00s)
    --- PASS: TestProbeRoutes/OPTIONS/healthz{"error":"method_not_allowed"}_ (0.00s)
    --- PASS: TestProbeRoutes/GET/{"error":"not_found"}_ (0.00s)
    --- PASS: TestProbeRoutes/GET/debug/pprof/{"error":"not_found"}_ (0.00s)
    --- PASS: TestProbeRoutes/GET/healthz/{"error":"not_found"}_ (0.00s)
    --- PASS: TestProbeRoutes/GET/x/../healthz{"error":"not_found"}_ (0.00s)
=== RUN   TestServeReadinessAndShutdown
--- PASS: TestServeReadinessAndShutdown (0.00s)
=== RUN   TestShutdownDrainsOrTimesOut
=== RUN   TestShutdownDrainsOrTimesOut/drain
=== RUN   TestShutdownDrainsOrTimesOut/deadline
--- PASS: TestShutdownDrainsOrTimesOut (0.19s)
    --- PASS: TestShutdownDrainsOrTimesOut/drain (0.01s)
    --- PASS: TestShutdownDrainsOrTimesOut/deadline (0.18s)
=== RUN   TestMalformedHTTPAndHeaderLimit
=== RUN   TestMalformedHTTPAndHeaderLimit/bad-version
=== RUN   TestMalformedHTTPAndHeaderLimit/bad-header
=== RUN   TestMalformedHTTPAndHeaderLimit/ambiguous-length
=== RUN   TestMalformedHTTPAndHeaderLimit/oversized-header
--- PASS: TestMalformedHTTPAndHeaderLimit (0.57s)
    --- PASS: TestMalformedHTTPAndHeaderLimit/bad-version (0.01s)
    --- PASS: TestMalformedHTTPAndHeaderLimit/bad-header (0.00s)
    --- PASS: TestMalformedHTTPAndHeaderLimit/ambiguous-length (0.00s)
    --- PASS: TestMalformedHTTPAndHeaderLimit/oversized-header (0.00s)
=== RUN   TestSlowHeadersAreClosed
--- PASS: TestSlowHeadersAreClosed (0.03s)
=== RUN   TestStartupAndServeFailures
--- PASS: TestStartupAndServeFailures (0.00s)
PASS
ok  	github.com/LiiWzz/DeviceBridge/server/control/internal/control	1.252s
exit: 0
```

```text
$ go test -race ./...
?   	github.com/LiiWzz/DeviceBridge/server/control/cmd/control	[no test files]
ok  	github.com/LiiWzz/DeviceBridge/server/control/internal/control	2.773s
exit: 0
```

```text
$ go build -trimpath -o bin/control ./cmd/control
(no output)
exit: 0
```

```text
$ CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -o bin/control-linux-armv7 ./cmd/control
(no output)
exit: 0
```

```text
$ CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o bin/control-linux-arm64 ./cmd/control
(no output)
exit: 0
```

```text
$ file bin/control bin/control-linux-armv7 bin/control-linux-arm64
bin/control:             Mach-O 64-bit executable arm64
bin/control-linux-armv7: ELF 32-bit LSB executable, ARM, EABI5 version 1 (SYSV), statically linked, Go BuildID=rgrj7li0r9UZKKEb7bLq/ZDYlXzmcgTRpUmqGi4F7/ThIAV9nEpzvKTBDdEZPC/6tddyjYDFIJuzntWoGor, with debug_info, not stripped
bin/control-linux-arm64: ELF 64-bit LSB executable, ARM aarch64, version 1 (SYSV), statically linked, Go BuildID=FvXGmj3wUq2LmEhUOFfX/UwXf7dFatO7k5HFyzWvD/PokKiSHviYYTn_bg9yWt/O0qRwoHitTHO2hirT2Io, with debug_info, not stripped
exit: 0
```

```text
$ ls -l bin/control bin/control-linux-armv7 bin/control-linux-arm64
-rwxr-xr-x@ 1 liiwz  staff  7547458 Nov  4 03:47 bin/control
-rwxr-xr-x@ 1 liiwz  staff  7526953 Nov  4 03:47 bin/control-linux-arm64
-rwxr-xr-x@ 1 liiwz  staff  7527337 Nov  4 03:47 bin/control-linux-armv7
exit: 0
```

```text
$ go list -m all
github.com/LiiWzz/DeviceBridge/server/control
exit: 0
```

```text
$ git diff --check
(no output)
exit: 0
```

## Native executable smoke

The exact self-contained check below uses Python only as a local test driver;
the service runtime has no Python dependency.

```sh
python3 - <<'PY'
import http.client, json, os, signal, socket, subprocess, time

base = {k: v for k, v in os.environ.items() if not k.startswith("DEVICEBRIDGE_CONTROL_")}

def environment(**values):
    return dict(base, **{"DEVICEBRIDGE_CONTROL_" + k: v for k, v in values.items()})

def request(port, method, path):
    conn = http.client.HTTPConnection("127.0.0.1", port, timeout=1)
    try:
        conn.request(method, path)
        response = conn.getresponse()
        return response.status, response.read().decode()
    finally:
        conn.close()

for sig in (signal.SIGINT, signal.SIGTERM):
    with socket.socket() as reservation:
        reservation.bind(("127.0.0.1", 0))
        port = reservation.getsockname()[1]
    proc = subprocess.Popen(["./bin/control"], env=environment(ADDR=f"127.0.0.1:{port}"), stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    try:
        deadline = time.monotonic() + 5
        while True:
            try:
                assert request(port, "GET", "/readyz")[0] == 200
                break
            except (ConnectionError, OSError):
                if time.monotonic() >= deadline or proc.poll() is not None:
                    raise
                time.sleep(0.01)
        for method, path, expected in (("GET", "/healthz", 200), ("HEAD", "/readyz", 200), ("POST", "/healthz", 405), ("GET", "/debug/pprof/", 404)):
            status, body = request(port, method, path)
            assert status == expected, (method, path, status, body)
            if method == "HEAD":
                assert body == ""
            print(f"{sig.name} {method} {path}: {status} {body.strip()}".rstrip())
        proc.send_signal(sig)
        stdout, stderr = proc.communicate(timeout=5)
        assert proc.returncode == 0, (proc.returncode, stderr)
        logs = [json.loads(line) for line in stderr.splitlines()]
        assert [entry["msg"] for entry in logs] == ["control API listening", "control API shutting down", "control API stopped"], logs
        print(f"{sig.name} exit={proc.returncode}; JSON logs:")
        print(stderr, end="")
        try:
            with socket.create_connection(("127.0.0.1", port), timeout=1):
                raise AssertionError("listener still open")
        except ConnectionRefusedError:
            print(f"{sig.name}: listener closed")
    finally:
        if proc.poll() is None:
            proc.kill()
            proc.communicate()

for values in ({"ADDR": "0.0.0.0:8080"}, {"READ_TIMEOUT": "0"}, {"ADDR": ""}):
    result = subprocess.run(["./bin/control"], env=environment(**values), capture_output=True, text=True, timeout=5)
    assert result.returncode == 1, result
    log = json.loads(result.stderr)
    assert log["msg"] == "invalid configuration" and log["level"] == "ERROR", log
    print(f"Invalid config {values!r}: exit={result.returncode}; {result.stderr.strip()}")

with socket.socket() as occupied:
    occupied.bind(("127.0.0.1", 0))
    occupied.listen()
    port = occupied.getsockname()[1]
    result = subprocess.run(["./bin/control"], env=environment(ADDR=f"127.0.0.1:{port}"), capture_output=True, text=True, timeout=5)
    assert result.returncode == 1, result
    assert json.loads(result.stderr)["msg"] == "control API failed"
    print(f"Occupied port: exit={result.returncode}; {result.stderr.strip()}")
print("PASS: native binary probes, SIGINT/SIGTERM, JSON logs, closed listeners, config and bind failures")
PY
```

```text
SIGINT GET /healthz: 200 {"status":"ok"}
SIGINT HEAD /readyz: 200
SIGINT POST /healthz: 405 {"error":"method_not_allowed"}
SIGINT GET /debug/pprof/: 404 {"error":"not_found"}
SIGINT exit=0; JSON logs:
{"time":"2026-11-04T03:48:21.029208+07:00","level":"INFO","msg":"control API listening","addr":"127.0.0.1:54985"}
{"time":"2026-11-04T03:48:21.041352+07:00","level":"INFO","msg":"control API shutting down"}
{"time":"2026-11-04T03:48:21.041383+07:00","level":"INFO","msg":"control API stopped"}
SIGINT: listener closed
SIGTERM GET /healthz: 200 {"status":"ok"}
SIGTERM HEAD /readyz: 200
SIGTERM POST /healthz: 405 {"error":"method_not_allowed"}
SIGTERM GET /debug/pprof/: 404 {"error":"not_found"}
SIGTERM exit=0; JSON logs:
{"time":"2026-11-04T03:48:21.047873+07:00","level":"INFO","msg":"control API listening","addr":"127.0.0.1:54993"}
{"time":"2026-11-04T03:48:21.05754+07:00","level":"INFO","msg":"control API shutting down"}
{"time":"2026-11-04T03:48:21.057565+07:00","level":"INFO","msg":"control API stopped"}
SIGTERM: listener closed
Invalid config {'ADDR': '0.0.0.0:8080'}: exit=1; {"time":"2026-11-04T03:48:21.062791+07:00","level":"ERROR","msg":"invalid configuration","error":"DEVICEBRIDGE_CONTROL_ADDR must use a loopback IP literal"}
Invalid config {'READ_TIMEOUT': '0'}: exit=1; {"time":"2026-11-04T03:48:21.067715+07:00","level":"ERROR","msg":"invalid configuration","error":"DEVICEBRIDGE_CONTROL_READ_TIMEOUT must be between 1ms and 5m"}
Invalid config {'ADDR': ''}: exit=1; {"time":"2026-11-04T03:48:21.072355+07:00","level":"ERROR","msg":"invalid configuration","error":"DEVICEBRIDGE_CONTROL_ADDR must be a loopback IP and port"}
Occupied port: exit=1; {"time":"2026-11-04T03:48:21.077156+07:00","level":"ERROR","msg":"control API failed","error":"listen: listen tcp 127.0.0.1:55002: bind: address already in use"}
PASS: native binary probes, SIGINT/SIGTERM, JSON logs, closed listeners, config and bind failures
exit: 0
```

## Self-review and assessment

Self-review covered every source file and the module README, specifically route
exposure, configuration failures, error logging, shutdown races, HTTP resource
settings, dependencies and Git scope. A review finding was fixed before the
recorded final checks: the graceful-drain test now has a 1s shutdown allowance,
while the forced-timeout case separately uses 100ms, reducing scheduler-related
false failures. No unresolved correctness findings were identified in this review;
independent coordinator review remains required.

Security: loopback IP literals only, unprivileged ports only, no DNS resolution or
wildcard bind, no credentials, persistent configuration, authorization stubs or
debug endpoints. The only endpoints are health/readiness with GET/HEAD; other
methods and paths fail explicitly. Malformed framing, conflicting content lengths,
oversized headers and slow headers are exercised over real local TCP connections.
All five durations are positive and bounded, header timeout cannot exceed the
whole read timeout, header size is bounded, readiness is atomic, and shutdown
has a deadline plus forced connection close. No architecture/security contract
changes or shared-file edits were made. No production control or connectivity
security guarantee is claimed. A targeted local secret-pattern scan is recorded
below; it is not a comprehensive secret-scanner certification.

Resources: standard library only, fixed small response strings, no per-request
application JSON encoder, no request-body parsing or request logging, no database,
no background polling, and no application filesystem writes. JSON logs go to
stderr; deployment log storage is an operator decision. Go net/http still allocates
buffers and goroutines per connection; there is no global connection limit or
load/RAM guarantee. Binaries are roughly 7.5 MB with debug information; artifact
size is not RSS. This is a server-side service, not the embedded Rust agent.
No target ARM Linux hardware or Linux emulator is available in this worktree:
ARMv7/AArch64 binaries were cross-built but not executed. Real 128 MB device RSS,
CPU, latency, restart behavior, Buildroot/Yocto/Debian compatibility and all fleet
performance targets remain unmeasured. IPv6 configuration is unit-tested; runtime
smoke used IPv4 loopback. Go 1.23.4 is the installed/used toolchain; newer versions
were not validated. No containers or external services were required.

Rollback: stop the control process with SIGINT/SIGTERM, remove its deployment
entry if one was created by an operator, and revert this issue's commit. This
module creates no schema, persistent state or migrations, and the implementation
needs no root privileges. This worker did not deploy it, push, open a PR, merge,
approve, change branch protections or enable auto-merge. The coordinator reports
main/develop lack protection and auto-merge is disabled; this worker did not
independently inspect remote settings or change them. Coordinator owns independent
review, feature push and Draft PR into develop citing #3; human Tech Lead owns
approval and manual merge.

## Final test and targeted secret scan

```text
$ go test ./...
?   	github.com/LiiWzz/DeviceBridge/server/control/cmd/control	[no test files]
ok  	github.com/LiiWzz/DeviceBridge/server/control/internal/control	0.948s
exit: 0
```

```text
$ rg -n 'BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY|AKIA[0-9A-Z]{16}|gh[pousr]_[A-Za-z0-9]{20,}' cmd internal go.mod README.md EVIDENCE.md
(no output)
exit: 1
```

For rg, exit 1 means no matching secret pattern; no scanner installation is claimed.

Worktree: `/Users/liiwz/orca/workspaces/DeviceBridge/feat-m01-db-003-control-api`.
The first staged whitespace check found trailing spaces in HEAD smoke output;
the test-driver print was corrected and smoke evidence was captured again above.

## Staged review checks

```text
$ git diff --cached --check
(no output)
exit: 0
```

```text
$ git diff --cached --name-only
server/control/EVIDENCE.md
server/control/README.md
server/control/cmd/control/main.go
server/control/go.mod
server/control/internal/control/config.go
server/control/internal/control/config_test.go
server/control/internal/control/server.go
server/control/internal/control/server_test.go
exit: 0
```
