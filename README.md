# DeviceBridge

Secure, lightweight remote maintenance infrastructure for embedded Linux / ARM devices behind NAT/CGNAT.

**Status:** Milestone 01 foundation only. No functional tunnel, enrollment, or remote SSH is implemented.

## Design goals
- Device RAM: 128 MB total; agent idle RSS target <= 16 MB (to be validated on real ARM hardware)
- Mixed Linux distributions: Buildroot, Yocto, Debian/Ubuntu variants
- Fleet: 100 devices; outbound-only TLS relay as future M02 capability
- SSH-native access, least privilege, relay cannot decrypt SSH traffic

## Repository roadmap

Delivery plan: **M01 foundation -> M02 secure connectivity -> M03 management
dashboard**. UX/UI design can start before M02 completes; production dashboard
integration depends on accepted connectivity and approved API/security contracts.
See the [product roadmap and UX/UI plan](docs/ROADMAP.md) for screens, work
packages, dependencies and acceptance gates. M02/M03 are planned work, not
implemented capabilities.

An [in-page interactive SSH terminal proposal](docs/adr/ADR-003-browser-terminal-scope.md)
is tracked separately for M03 following the owner's request. It is unimplemented
and awaits scope/security approval; the existing native CLI plan remains the
baseline until the relevant decisions and acceptance gates are satisfied.

Try the [local browser-terminal state demo](docs/ux/UX-05.md) by opening its single
HTML file. It simulates typing, paste and session outcomes; it connects to no device
and executes no commands. Human flow review and production approval remain pending.

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

Build/test automation and local reproduction commands are documented in
[CI checks](ci/README.md). ARM cross-build outputs do not establish physical
device compatibility or memory acceptance.

The [security baseline](security/README.md) records M01 trust boundaries, config
validation, secret handling, and proposed credential lifecycle review gates.
