# ADR-001: Foundation and deployment boundaries

**Status:** Proposed — Tech Lead approval required
**Context:** Embedded Linux/IoT ARM, mixed Linux variants, 128 MB system RAM, initial fleet of 100 devices behind CGNAT.

## Decisions
- Rust agent; Go control service, relay and developer CLI.
- Remote shell stays OpenSSH/Dropbear; agent will later forward authorized encrypted SSH byte streams only.
- Outbound-first, relay-first over TLS/TCP 443 for future M02; no port-forward or public IP required on devices.
- PostgreSQL server side; no database, Python, Node or Docker required on device.
- Control and data planes separate; transport protocol framing is explicitly versioned.
- P2P, QUIC, browser dashboard, video streaming deferred.
- Device identity, per-session authorization and host-key verification required before production SSH access.

## Security invariants
- No SSH private keys stored on relay.
- No default fleet credentials or shared long-lived device tokens.
- No root permissions for agent unless justified and approved.
- Fail closed on authorization/token expiry.
- Relay enforces bounded streams and buffers.
- Never ship unverified security-sensitive stubs as production features.

## Trade-offs
- Relay-first increases network path length but reduces NAT complexity.
- Rust improves memory safety but does not itself guarantee <=16 MB RSS.
- SSH reuses audited functionality but device SSH server compatibility and host-key pinning need verification.

## Validation
Test on real ARMv7/AArch64 samples, with 128 MB RAM and Buildroot/Yocto/Debian cases; report P50/P95/P99, memory, CPU and reconnect under controlled network profiles. Results unknown until hardware validation.

## Follow-up
Implementation contracts and threats reviewed in DB-007; delivery acceptance by DB-010.
