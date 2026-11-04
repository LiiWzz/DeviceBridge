# Developer CLI foundation (DB-004 / #4)

This separate Go module requires Go 1.23 or later and uses only the standard
library. It offers help, version, and local configuration validation. No network
connection, shell command, enrollment, authentication, or forwarding is performed.

From this directory:

```sh
go build -trimpath -o bin/devicebridge ./cmd/devicebridge
./bin/devicebridge --help
./bin/devicebridge version
./bin/devicebridge check
./bin/devicebridge check --relay-url https://relay.example:443 --timeout 5s
```

No arguments prints help and exits 0. `help` / `-h` / `--help` and `version` /
`--version` exit 0. `check --help` prints check flags. Unknown commands (including
`connect`, `ssh`, and `shell`), unknown flags, extra arguments, and invalid
configuration exit 2. `check` exits 0 only for locally valid configuration;
**this does not indicate server readiness, authentication, or connectivity**.

Settings are flags only; no files or environment credentials are read:

| `check` flag | Default | Validation |
| --- | --- | --- |
| `--relay-url` | `https://localhost:443` | Absolute HTTPS URL with ASCII DNS host or IP address; optional port 1..65535; root path only; no userinfo, query, fragment, or IPv6 zone |
| `--timeout` | `5s` | Go duration greater than zero and at most one minute |

The endpoint and timeout reserve future settings, never initiate a request, and
are not a transport contract. Invalid URL errors do not echo endpoint credentials.
International DNS names must be supplied as ASCII/punycode. Transport remains
unimplemented; no insecure mode or authentication bypass exists.

Run checks within this module (there is no shared root Go module):

```sh
gofmt -l cmd internal
go vet ./...
go test ./...
go test -race ./...
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -o bin/devicebridge-linux-armv7 ./cmd/devicebridge
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o bin/devicebridge-linux-arm64 ./cmd/devicebridge
```

See [EVIDENCE.md](EVIDENCE.md) for actual output and unsupported hardware checks.
