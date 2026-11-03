# Relay foundation (DB-004 / #4)

This independent Go module requires Go 1.23 or later and only uses the standard
library. M01 has **no network listener, forwarding, SSH, enrollment, or usable
transport**. The default process idles until SIGINT or SIGTERM and exits cleanly.

From this directory:

```sh
go build -trimpath -o bin/devicebridge-relay ./cmd/devicebridge-relay
./bin/devicebridge-relay --help
./bin/devicebridge-relay --version
./bin/devicebridge-relay --check
./bin/devicebridge-relay
```

Configuration is command-line flags only; no config files, environment secrets,
credentials, or production settings are read. Unknown flags and positional
arguments fail with exit 2. `--check` validates and exits 0, without starting.

| Flag | Default | Validation |
| --- | --- | --- |
| `--transport` | `disabled` | All other values fail closed with exit 2 |
| `--max-streams` | `16` | Integer 1..1024 |
| `--max-buffer-bytes` | `32768` | Integer 1..1048576 |

Stream and buffer settings reserve future ceilings only: no streams or buffers
are allocated, and no production enforcement is claimed. The internal,
provisional `transport.Connector` seam defines no wire protocol or authorization
contract; `transport.Disabled.Open` always returns a nil stream and
`ErrUnimplemented`. Reviewed identity, per-session authorization, host-key
verification, framing, bounds, and expiry handling are future work before any
transport can be enabled.

Run checks in this module (there is no shared root Go module):

```sh
gofmt -l cmd internal
go vet ./...
go test ./...
go test -race ./...
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -o bin/devicebridge-relay-linux-armv7 ./cmd/devicebridge-relay
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o bin/devicebridge-relay-linux-arm64 ./cmd/devicebridge-relay
```

See [EVIDENCE.md](EVIDENCE.md) for executed checks, limitations, security/resource
assessment, ADR status disclosure, and rollback.
