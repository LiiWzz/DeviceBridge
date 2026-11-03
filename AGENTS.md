# DeviceBridge AI coding agent policy

## Mandatory rules
1. Read `README.md`, `docs/M01_EXECUTION.md`, and architecture ADRs before coding.
2. Execute one GitHub issue at a time; follow dependencies; work in an isolated feature branch based on `develop`.
3. Never push to `main` or `develop`. Never merge any PR, approve your own PR, disable protections, or enable auto-merge.
4. Create **Draft PRs** with base `develop` and cite the GitHub issue. Human Tech Lead alone approves merging.
5. No claim of passing checks without executed commands and captured output. Explicitly state unsupported toolchains or missing ARM hardware.
6. Never commit credentials, device tokens, private keys, or production configuration.
7. On ARM: avoid allocations, unnecessary runtime dependencies, writes to flash, and root privileges.
8. No custom shell protocol, unsafe authentication bypass, or unauthenticated port-forwarding.
9. If contracts or security policies must change, submit an ADR and wait for Tech Lead approval.
10. Record test evidence, performance implications, security implications and rollback steps in every PR.

## PR checklist
- [ ] Scope limited to linked issue
- [ ] Build/lint/tests executed (include exact commands and outcome)
- [ ] Negative/security cases covered
- [ ] No secrets introduced
- [ ] Cross-platform and RAM considerations noted
- [ ] New docs/ADR where necessary
- [ ] Draft PR into `develop`; awaiting human approval
