# Milestone 01 — Execution contract

**Goal:** reproducible, secure project foundation. M01 does not include remote SSH or production transport.

## Issue work graph
| Issue | Scope | Dependencies | Agent role |
|---|---|---|---|
| [DB-001](https://github.com/LiiWzz/DeviceBridge/issues/1) | foundation and ADR | none | Architect |
| [DB-002](https://github.com/LiiWzz/DeviceBridge/issues/2) | Rust ARM agent skeleton | DB-001 | Rust |
| [DB-003](https://github.com/LiiWzz/DeviceBridge/issues/3) | Go control API skeleton | DB-001 | Go |
| [DB-004](https://github.com/LiiWzz/DeviceBridge/issues/4) | Go relay and CLI skeleton | DB-001 | Go |
| [DB-005](https://github.com/LiiWzz/DeviceBridge/issues/5) | PostgreSQL and Compose | DB-001 | DevOps |
| [DB-006](https://github.com/LiiWzz/DeviceBridge/issues/6) | CI + cross-build | DB-002/003/004 | DevOps |
| [DB-007](https://github.com/LiiWzz/DeviceBridge/issues/7) | Security baseline | DB-002/003/004/005 | Security |
| [DB-008](https://github.com/LiiWzz/DeviceBridge/issues/8) | Integration smoke | DB-003/004/005/006/007 | QA |
| [DB-009](https://github.com/LiiWzz/DeviceBridge/issues/9) | Independent review | DB-001..008 | Reviewer |
| [DB-010](https://github.com/LiiWzz/DeviceBridge/issues/10) | Acceptance report | DB-009 | Human Tech Lead |

## Test gates
- Rust: cargo fmt, cargo clippy, cargo test; cross-build ARMv7/AArch64
- Go: gofmt, go vet, go test ./..., builds for Linux ARM where applicable
- Compose: docker compose config + service readiness and shutdown smoke
- Security: secret scan, malformed config negative tests, least-privilege review
- Target hardware: ARM device RAM 128 MB; measure idle RSS, CPU and restart handling, record actual model/kernel
- No successful remote connectivity or security guarantees claimed by M01

## Definition of done
All issues reviewed, checks passing, reproducible evidence attached, critical/high findings closed, Tech Lead GO/NO-GO recorded. All PR merges manual; no direct `main` changes.

## Tech Lead action
Set GitHub branch rulesets for `main` and `develop`; require review/status checks, disallow force push/deletions and direct commits, keep auto-merge disabled. Linked GitHub integration cannot configure protection.
