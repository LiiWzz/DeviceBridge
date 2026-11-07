# Milestone 02 — Proposed execution plan

Planning issue: [#35](https://github.com/LiiWzz/DeviceBridge/issues/35).
Local primary-source baseline: `8d258fdacad19a896026dccaf8e67a38218dbdb7`,
obtained with `git rev-parse HEAD`; prepared 2026-11-06.

**Status: proposed for human review; no implementation or acceptance approval.**
This document decomposes the [roadmap's M02 scope](ROADMAP.md#milestones-and-release-gates)
into reviewable work. Package IDs below are proposed planning references, **not
created GitHub issues**. They select no API/schema, protocol framing, credential
format, cryptography, provider, role model or numerical security/resource policy.
Only a separately scoped issue and its approved prerequisites authorize coding.

## Existing foundation and blocking gates

The [README](../README.md) describes M01 foundation only; enrollment, functional
tunnels and remote SSH are absent. Source inspection supports that boundary:
[agent `run`](../agent/src/main.rs) idles without networking or disk writes;
[control `service.handle`](../server/control/internal/control/server.go) serves
health/readiness probes; [relay `Disabled.Open`](../server/relay/internal/transport/transport.go)
rejects opening a stream; [CLI `run`](../cli/cmd/devicebridge/main.go) validates
local settings without connecting. These observations are not runtime tests.

[M01's execution contract](M01_EXECUTION.md) still requires independent review,
hardware evidence and human GO/NO-GO. [DB-010](acceptance/DB-010.md#open-gates-and-owner-specific-closure)
records pending DB-009 review, ADR approval, physical ARM acceptance and repository
protections; its historical snapshot is not a fresh governance audit. Reconcile
historical findings with later maintenance PRs before recording M01 closure.
Merged documents or closed issues cannot replace those decisions. M02 planning
may proceed, but M02 implementation follows explicit M01 acceptance.

[ADR-001](adr/ADR-001-foundation.md) remains **Proposed — Tech Lead approval
required**. Its recorded boundaries guide this proposal without changing its
status: Rust agent; Go control/relay/CLI; outbound-first TLS/TCP 443; separate
control/data planes; versioned transport framing; OpenSSH/Dropbear encrypted SSH
streams; no relay SSH private keys, default fleet credentials or unauthenticated
forwarding. No custom shell protocol or authentication bypass is permitted.
Under [AGENTS.md](../AGENTS.md), contract/security changes need separate ADRs and
explicit human Tech Lead approval before dependent implementation.

## Dependency graph and accountable evidence

Owners below are proposed roles; the Tech Lead assigns accountable people and
scopes each package into small issues before execution. If a package needs
multiple issues, preserve its internal ordering and gate each dependent slice.
Execute one GitHub issue at a time on an isolated feature branch based on
`develop`; submit Draft PRs into `develop`, citing the issue. Independent review,
explicit human approval and manual merge remain required; never self-approve,
merge, push directly to `main`/`develop` or enable auto-merge.

| Proposed ID | Deliverable / dependency | Owner | Evidence required before dependent work |
| --- | --- | --- | --- |
| M02-P00 | M01 closure; no implementation | Human Tech Lead / Hardware QA / repository administrator | DB-009 disposition, ADR-001 approval/status, actual ARM baseline, current protection evidence, DB-010 signed GO for the intended scope. NO-GO leaves implementation blocked. |
| M02-D01 | Identity/enrollment and authorization design; review inputs may precede P00 | Architect / Security / Backend / Agent | Proposed ADRs covering device/user trust boundaries, bootstrap binding/replay, capabilities/resources and per-session checks; explicit Tech Lead approval before implementation. |
| M02-D02 | Lifecycle, time/outage and audit design; depends on D01 agreement | Security / Backend / Agent | Separate ADR decisions for issue/store/renew/rotate/revoke/decommission, active-session effects, audit producer/storage/access/failure semantics and recovery; explicit Tech Lead approvals. |
| M02-D03 | Transport, host-key and bounds design; depends on D01/D02 agreement | Architect / Relay / Agent / CLI / Security | Proposed ADRs for versioned framing, TLS trust, admission/aggregate bounds, reconnect, SSH host-key trust/recovery and native handoff; approved measurable budgets and Tech Lead approval. |
| M02-P01 | Identity verification and resource/action enforcement foundations; P00 + approved D01/D02 | Backend / Security | Direct allow/deny tests for user/device/action/resource boundaries; protected identity provisioning prerequisites; no permissive stub or enrollment mutation exposed. |
| M02-P02 | Audit foundation; P01 + approved D02 | Backend / Security | Producer/store integrity, redaction, attempted/denied/completed outcomes and outage/reconciliation evidence; accepted before exposing security-sensitive mutations. |
| M02-P03 | Device-bound enrollment; P01/P02 + approved D01/D02 | Backend / Agent | Approved bootstrap delivery and protected identity storage; replay/expiry/wrong-device, duplicate/lost-response and restart tests; no fleet-wide shared credential. |
| M02-P04 | Credential lifecycle and revocation; P03 + approved D02 | Backend / Agent / Security | Renewal/rotation/revocation/decommission, atomic update/recovery, stale-state/outage and flash-write evidence; revoked identities never regain validity on rollback. |
| M02-P05 | Per-session authorization; P04/P02 + approved D01/D02 | Backend / Agent / CLI / Security | Fresh user/device/resource authorization for each session; allow/deny/expiry/revoked/replayed/wrong-session evidence; approved active-session and outage behavior. |
| M02-P06 | Bounded outbound TLS relay; P05 + approved D03 | Relay / Agent / Security | TLS peer verification, version rejection, authenticated stream binding and aggregate admission/backpressure evidence; no public device listener or unauthenticated forwarding. |
| M02-P07 | CLI host-key verification and SSH-native handoff; P06 + approved D03 | CLI / Agent / Security | OpenSSH/Dropbear interoperability, independently verified host-key binding, mismatch rejection and credential-safe handoff; no SSH private-key custody on relay. |
| M02-P08 | Integrated outage/security/resource and recovery validation; P04–P07 | QA / Security / Hardware QA / DevOps | Revision-bound test matrix, real ARM measurements, fleet/load profiles, raw failure/recovery evidence and runbook; critical/high findings resolved before acceptance. |
| M02-P09 | Independent review and acceptance packet; P08 | Independent reviewer, then Human Tech Lead | Review exact source/evidence, disposition findings, confirm every gate and record signed GO/NO-GO with supported platforms/exclusions. |

D01–D03 are design/review packages, not implied approvals or assigned ADR numbers.
Their boundaries may require several ADRs. Resolve cross-document inconsistencies
before approval; no package may implement an unresolved contract by inference.

The first scoped D01 design slice is [#43](https://github.com/LiiWzz/DeviceBridge/issues/43):
[Proposed ADR-002](adr/ADR-002-device-record-identity-binding.md) separates device
records, identities and authorized enrollment binding, with
[repository source findings](research/M02-D01a-sources.md) and a
[domain glossary](../GLOSSARY.md). It remains pending Tech Lead approval and
does not complete D01 or authorize P01/P03 or full WEB-01. Bootstrap mechanism,
credential/lifecycle/audit and transport decisions remain separately gated.

The separate [D01b / ADR-004](adr/ADR-004-session-authority-binding.md),
[#49](https://github.com/LiiWzz/DeviceBridge/issues/49), accepts binding session
authority to the user, current intended device, SSH account and session/client
context as limited domain constraints, with owner approval recorded in
[PR #53](https://github.com/LiiWzz/DeviceBridge/pull/53). Remaining D01 trust and
enforcement mechanisms, D02/D03 approvals and P00 remain prerequisites to implementation.

[D02 / ADR-005](adr/ADR-005-credential-lifecycle-audit.md),
[#51](https://github.com/LiiWzz/DeviceBridge/issues/51), records owner-accepted
domain principles for lifecycle, trusted-time, revocation, uncertain outcomes and
audit, including the recovery extension in #55. Credential/storage mechanisms,
remaining design decisions and measured budgets still require approval.
Executed document checks and scan output are in
[D02 validation](evidence/m02-d02/validation.txt) and
[D02 secret scan](evidence/m02-d02/secrets.txt).

Follow-up [#55](https://github.com/LiiWzz/DeviceBridge/issues/55) extends ADR-005
with deployed-device renewal, credential separation, bounded verified rotation
and independently authenticated recovery. Approved offline/recovery limits,
pre-provisioned trust and tested outage/power-loss paths are rollout blockers.
No recovery credential or protocol is selected. Evidence:
[validation](evidence/m02-recovery/validation.txt),
[secret scan](evidence/m02-recovery/secrets.txt).

[D03 / Proposed ADR-006](adr/ADR-006-transport-trust-and-bounds.md),
[#57](https://github.com/LiiWzz/DeviceBridge/issues/57), proposes separate TLS,
device/session and SSH host-key trust, opaque authorized stream handling,
native handoff and aggregate admission/backpressure obligations. Mechanisms,
framing, active-session policies and measured budgets remain approval blockers;
this design input completes no D03 or implementation gate. Evidence:
[validation](evidence/m02-d03/validation.txt), [scan](evidence/m02-d03/secrets.txt).

[D01c / Proposed ADR-007](adr/ADR-007-operator-trust-and-authority.md),
[#59](https://github.com/LiiWzz/DeviceBridge/issues/59), proposes trusted operator
identity mapping, explicit resource/action scope and non-amplifying delegation,
including initial trust/audit bootstrap obligations. Identity/provider, policy,
delegation and provisioning mechanisms remain approval blockers; it completes
no D01 or implementation gate. Evidence:
[validation](evidence/m02-d01c/validation.txt), [scan](evidence/m02-d01c/secrets.txt).

[D01d / ADR-008](adr/ADR-008-device-bootstrap-trust.md),
[#61](https://github.com/LiiWzz/DeviceBridge/issues/61), records owner-accepted
domain principles for independently
established intended-device trust, mutual bootstrap and fresh operation-bound
enrollment proof, with protected provisioning/custody and durable one-use
binding/reconciliation. Proof, service/time trust, delivery/storage mechanisms,
clone assumptions and measured budgets remain approval blockers; it completes
no D01 or implementation gate. Evidence:
[validation](evidence/m02-d01d/validation.txt), [scan](evidence/m02-d01d/secrets.txt).
The limited confirmation after PR #62 is recorded in
[#63](https://github.com/LiiWzz/DeviceBridge/issues/63); original proposal evidence
above remains historical. Approval-record checks are included in
[D02a validation](evidence/m02-d02a/validation.txt).

[D02a / ADR-009](adr/ADR-009-active-session-validity.md),
[#63](https://github.com/LiiWzz/DeviceBridge/issues/63), records owner agreement
with the limited policy of stopping both
directions of usable session traffic and bounded closure when validity is lost
or cannot be established. No automatic resume or replay; local/remote/audit
outcomes stay distinct. Detection, freshness, propagation, cutoff/closure budgets
and mechanisms remain approval blockers; this completes no D02/D03 or
implementation gate. Evidence:
[validation](evidence/m02-d02a/validation.txt), [scan](evidence/m02-d02a/secrets.txt).
The limited confirmation after PR #64 is recorded in
[#65](https://github.com/LiiWzz/DeviceBridge/issues/65); original proposal evidence
above remains historical. Approval-record checks are included in
[D02b validation](evidence/m02-d02b/validation.txt).

[D02b / ADR-010](adr/ADR-010-authoritative-audit-consistency.md),
[#65](https://github.com/LiiWzz/DeviceBridge/issues/65), records owner-approved
architecture direction for a shared
authoritative commit for protected control-state mutation, operation outcome
and required audit record, with separately derived read/export views. Remote
effects need their own evidence; no exporter acknowledgement replaces durable
audit or resurrects invalidated authority. Store/consistency/integrity/recovery
mechanisms and measured bounds remain approval blockers; this completes no
D02/P02 or implementation gate. Evidence:
[validation](evidence/m02-d02b/validation.txt), [scan](evidence/m02-d02b/secrets.txt).
Approval after PR #66 is recorded in
[#67](https://github.com/LiiWzz/DeviceBridge/issues/67); original proposal evidence
above remains historical. Approval-record checks are included in
[D02c validation](evidence/m02-d02c/validation.txt).

[D02c / Proposed ADR-011](adr/ADR-011-trusted-time-validity.md),
[#67](https://github.com/LiiWzz/DeviceBridge/issues/67), proposes independently
trusted current-time bounds wholly inside authorized validity windows, with
conservative local enforcement, non-circular bootstrap and uncertainty-aware
audit. Time evidence grants no identity/policy/revocation authority. Source,
verification/clock/continuity mechanisms and measured budgets remain approval
blockers; this completes no D02/D03 or implementation gate. Evidence:
[validation](evidence/m02-d02c/validation.txt), [scan](evidence/m02-d02c/secrets.txt).

P01 must scope bootstrap/provisioning explicitly so audit does not depend on an
unaudited enrollment shortcut. P02 precedes enrollment, credential changes,
revocation and session-grant mutations; lifecycle logs alone do not satisfy it.
This ordering follows [WEB-01a's audit prerequisite](web/WEB-01a.md#dependency-ordered-proposed-follow-up-slices)
and the [security lifecycle review inputs](../security/README.md#proposed-credential-and-token-lifecycle--unimplemented).

## Decisions the ADRs must resolve

The [security baseline](../security/README.md) and
[DB-010 recommendations](acceptance/DB-010.md#scope-exclusions-and-m02-recommendations)
identify these unresolved questions. The following are decision obligations,
not chosen contracts or newly accepted security policy:

- D01: operator/device bootstrap trust and recovery channel, user identity,
  device uniqueness/clone resistance, credential/verifier storage, authorization
  boundaries and binding of an authorized request to its intended device/session.
- D02: issuance, expiry, trusted time/skew, authenticated renewal, rotation overlap,
  revocation propagation and handling of active sessions, decommission/backup
  recovery, offline/cache semantics, crash consistency and bounded flash writes.
- D02: authoritative audit producers, outcome distinction, visibility, integrity,
  retention/redaction, duplicate/concurrent mutation reconciliation and behavior
  when the producer/store fails or a response is lost after a mutation.
- D03: versioned framing and compatibility, TLS trust/certificate validation,
  stream binding, cancellation/half-close, backpressure, retry/backoff, reconnect
  authorization, SSH host-key verification/bootstrap/change/recovery and native
  client handoff. No executable connect syntax is established by this plan.
- D02/D03: measurable lifetimes, overlap windows, revocation/outage guarantees,
  connection/handshake/session/query/queue/log bounds and CPU/memory/network
  budgets. Tech Lead-approved ADRs must choose values and test workloads before
  implementation; this proposal supplies none.

Bounds must cover aggregate live and in-flight work, not only per-stream payloads:
TLS/handshake state, connections, queues, buffers in both directions, pending
authorization, reconnect attempts and slow consumers. The reserved
[relay settings](../server/relay/internal/config/config.go) allocate nothing in
M01; multiplying their maxima permits 1 GiB if naively allocated, as recorded in
[security trust boundaries](../security/README.md#assets-and-trust-boundaries).
Do not adopt those maxima as M02 budgets. Review bounded allocation and avoid
unnecessary ARM runtime dependencies, flash writes and root privileges.

## Validation and evidence gates

Each implementation PR must bind its evidence to exact source/binary revisions,
approved ADRs and linked issue. Record exact commands, outputs, tool versions,
test fixtures, limitations, security/performance implications and rollback steps.
Select new transport/lifecycle tests after approval; source inspection, test
names and healthy probes cannot establish passing security or connectivity.
Retain [M01 component gates](M01_EXECUTION.md#test-gates) and
[CI reproduction commands](../ci/README.md) where applicable; execute and capture
the relevant Rust/Go formatting, lint, tests, builds, secret scans and integration
checks for changed components. Do not claim unsupported toolchains passed.

| Evidence area | Required future coverage / owner gate |
| --- | --- |
| Identity and enrollment | QA/Security: unauthenticated, unauthorized, wrong-user/device/resource, expired/revoked bootstrap, replay/duplicate/concurrent submission, cloned identity, malformed input and lost-response reconciliation. |
| Session authorization | QA/Security: direct requests bypassing CLI/UI, wrong-device/session binding, privilege separation, replay, expiry during use, revoked/stale grants, reconnect after revocation and control-plane/time-source outages under approved semantics. |
| Credential lifecycle / audit | QA/Security: failed renewal, rotation overlap boundaries, interrupted/power-loss update, stale caches, decommission/recovery, audit producer/store failure, denied/attempted/completed outcomes and secret redaction. |
| TLS / transport | QA/Relay: invalid or untrusted/expired TLS identity, wrong peer, unsupported/malformed/truncated framing, unauthorized forwarding, slow peers, backpressure, aggregate saturation, cancellation, disconnect/reconnect and server restart. |
| Native SSH / host keys | QA/CLI: correct authorized device with OpenSSH/Dropbear; missing/untrusted/changed/mismatched host key and wrong-resource substitution; no bypass on recovery, no custom shell protocol and no credentials in URLs/arguments/logs/audit. |
| Outages and recovery | QA/DevOps: control/relay/database failure and restart, network loss/flapping, uncertain mutation outcome and retained revoked state; prove approved fail-closed behavior rather than success inferred from missing errors. |

The [README baseline](../README.md#design-goals) is **128 MB total device RAM;
agent idle RSS target <= 16 MB**, requiring real ARM evidence. Report units
explicitly; do not silently substitute MiB. Hardware QA must confirm supported
ARMv7/AArch64 and Buildroot/Yocto/Debian coverage with the Tech Lead, retaining
[DB-010's unmeasured matrix](acceptance/DB-010.md#hardware-baseline--measurements-unavailable).
Cross-builds and an ARM macOS host do not establish target execution/acceptance.

For each supported board record actual model/ISA/RAM, distribution/kernel/libc/
loader, ordinary-user setup, revision/binary hashes, tools, sampling method,
warm-up/window and raw RSS/CPU/restart logs. Measure idle, authorized sessions,
aggregate saturation, reconnect/outage and recovery; record flash writes and
allocation impact. Under [ADR-001 validation](adr/ADR-001-foundation.md#validation),
report P50/P95/P99 for approved network/load profiles and the initial 100-device
fleet. Unmeasured/unsupported combinations remain explicit blockers or human
dispositions under the approved process, never cross-build-derived passes.

## Recovery, review and human decision

[Proposed ADR-003](adr/ADR-003-browser-terminal-scope.md), issue
[#45](https://github.com/LiiWzz/DeviceBridge/issues/45), requests browser terminal
as a separately gated M03 extension. It changes no M02 deliverable or approval:
M02 retains SSH-native device transport and independent host-key/per-session
checks. Browser production integration requires accepted M02 and the additional
architecture, web-session, credential-custody, audit and UX decisions in ADR-003.

P08's runbook must cover disabling new access safely, restoring an approved
deployment, compatible schema/state recovery, interrupted identity updates and
audit preservation. Rollback must never revive revoked credentials, bypass
authorization or host-key verification, or silently weaken approved bounds.
Irreversible migration/credential consequences require reviewed recovery steps;
a code revert alone cannot revoke secrets or restore trust. These obligations
derive from the [security rollback guidance](../security/README.md#reproduction-impact-and-rollback).

P09 must record exact reviewed commits/PRs and ADR approval references, M01 gate
closure, hardware/toolchain limitations, independent findings/dispositions and
supported operational scope before requesting human acceptance. **Human decision:
PENDING; intentionally unsigned.** Only the human Tech Lead records GO/NO-GO,
rationale, exclusions, identity and authoritative signature/reference with UTC
timestamp. This plan signs nothing and authorizes no merge, release or production
SSH access. No gate is satisfied merely by merging this document.

M03 remains separately gated: [WEB-01a](web/WEB-01a.md#capability-separation-and-owner-gates)
is preliminary discovery, full WEB-01 needs agreed M02 identity/authorization,
and UX-03 review, approved web/API/security contracts and accepted M01/M02 remain
production integration prerequisites. Browser terminal, dashboard runtime on ARM,
P2P/QUIC, video, bulk commands and OTA remain outside this M02 plan under the
[roadmap](ROADMAP.md#explicitly-deferred) and ADR-001 boundaries.

## Documentation evidence and impact

Documentation validation and changed-file secret-scan commands/output belong in
[validation.txt](evidence/m02-plan/validation.txt) and
[secrets.txt](evidence/m02-plan/secrets.txt); only captured executions establish
their outcomes. This planning work claims no new application build/test, live
connectivity, physical ARM, ARM GNU toolchain or Windows runtime validation and
no fresh governance audit. Pattern scans cannot guarantee absence of all secrets.

This documentation adds no code, dependency, production configuration, schema,
credential, allocation, flash write or privilege requirement. Future security
and performance effects remain approval and measurement obligations. Supersede
or revert this plan through a feature branch and Draft PR into `develop`, keeping
source/decision provenance; human review and manual merge remain required. No
service/device/database/credential recovery is needed to revert this document.
