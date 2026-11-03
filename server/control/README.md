# Control API foundation (DB-003)

Independent Go module for [issue #3](https://github.com/LiiWzz/DeviceBridge/issues/3).
Requires Go 1.23 or later and uses only the standard library. Run from this directory:

```sh
go run ./cmd/control
curl --fail http://127.0.0.1:8080/healthz
curl --fail http://127.0.0.1:8080/readyz
```

`GET` and `HEAD /healthz` return 200 while the HTTP server serves requests.
`GET` and `HEAD /readyz` return 200 when serving and 503 when draining.
Responses use JSON; HEAD responses have no body. Unknown paths return 404;
other methods on the probe paths return 405 with `Allow: GET, HEAD`.
Readiness describes only this process: there are no database or external dependencies.
Shutdown clears readiness and stops accepting connections immediately; external probes
may see connection refusal instead of 503. Existing requests drain within the timeout.

Environment values are optional; explicitly empty or malformed values fail startup
with exit code 1. Configuration values are not echoed in validation errors.

| Variable | Default | Validation |
| --- | --- | --- |
| `DEVICEBRIDGE_CONTROL_ADDR` | `127.0.0.1:8080` | Loopback IP literal, canonical decimal port 1024–65535; IPv6 such as `[::1]:8080` supported; no hostname, zone or wildcard |
| `DEVICEBRIDGE_CONTROL_READ_HEADER_TIMEOUT` | `5s` | Go duration, 1ms–5m, no longer than read timeout |
| `DEVICEBRIDGE_CONTROL_READ_TIMEOUT` | `10s` | Go duration, 1ms–5m |
| `DEVICEBRIDGE_CONTROL_WRITE_TIMEOUT` | `10s` | Go duration, 1ms–5m |
| `DEVICEBRIDGE_CONTROL_IDLE_TIMEOUT` | `60s` | Go duration, 1ms–5m |
| `DEVICEBRIDGE_CONTROL_SHUTDOWN_TIMEOUT` | `10s` | Go duration, 1ms–5m |

Headers are limited to 8 KiB plus the Go HTTP parser's read-ahead allowance.
Logs are structured JSON (`slog`) on stderr. Startup, shutdown and failures are
logged; requests, headers, bodies and environment values are not application-logged.
SIGINT and SIGTERM trigger graceful shutdown. Startup/serve failures or exhausted
shutdown deadlines exit nonzero; deadline exhaustion forcibly closes connections.
The binary needs no root privileges, persistent storage or configuration files.

This module provides local HTTP probes only. It has no authentication, enrollment,
remote commands, SSH forwarding, tunnels, NAT traversal, debug endpoints or production
transport. Do not expose it through a proxy as a production control plane. It makes
no claim of remote connectivity or production security readiness.

```sh
gofmt -l cmd internal
go vet ./...
go test ./...
go test -race ./...
go build -trimpath -o bin/control ./cmd/control
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -o bin/control-linux-armv7 ./cmd/control
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o bin/control-linux-arm64 ./cmd/control
```

See [EVIDENCE.md](EVIDENCE.md) for executed checks, limitations, resource/security
assessment and rollback. All generated binaries belong in ignored `bin/`.
