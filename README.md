# DeviceBridge

Secure, lightweight remote maintenance infrastructure for embedded Linux / ARM devices behind NAT/CGNAT.

**Status:** Milestone 01 foundation only. No functional tunnel, enrollment, or remote SSH is implemented.

## Design goals
- Device RAM: 128 MB total; agent idle RSS target <= 16 MB (to be validated on real ARM hardware)
- Mixed Linux distributions: Buildroot, Yocto, Debian/Ubuntu variants
- Fleet: 100 devices; outbound-only TLS relay as future M02 capability
- SSH-native access, least privilege, relay cannot decrypt SSH traffic

## Repository roadmap
- `agent/`: Rust embedded agent (DB-002)
- `server/control/`: Go control API (DB-003)
- `server/relay/` and `cli/`: Go relay and CLI (DB-004)
- `deploy/`: local infrastructure (DB-005)
- `tests/`: integration/security/performance tests (DB-008)
- `docs/`: architectural decisions, threat models, acceptance criteria

Directory ownership is reserved; empty folders are created by their owning issues.

## Installation previews

Portable packages install with `sh install.sh` after extraction, without compiling
on the destination. They provide explicit upgrade, rollback and uninstall in a
private user prefix. See the [installation guide](packaging/README.md) or
[คู่มือติดตั้งภาษาไทย](packaging/README.th.md). These are M01 foundation previews;
remote SSH/tunneling remains unimplemented. No public binary release is published
by this change. Package tooling consumes explicit prebuilt artifacts; application
PR approval/manual merge and release validation remain separate gates.

## Workflow
Feature branch -> Draft PR targeting `develop` -> CI -> independent review -> explicit human Tech Lead approval -> manual merge. Never automatically merge into `main` or `develop`.

See [AGENTS.md](AGENTS.md), [M01 execution plan](docs/M01_EXECUTION.md) and [ADR-001](docs/adr/ADR-001-foundation.md).
