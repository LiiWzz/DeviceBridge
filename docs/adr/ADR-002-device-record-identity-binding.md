# ADR-002: Separate device records from enrollment identity binding

**Status: Proposed — explicit human Tech Lead approval required.**
Issue: [#43](https://github.com/LiiWzz/DeviceBridge/issues/43).
Source baseline: `e32cd6630b72452dd7c50bc0c448eef70e30d515`.
This is a limited M02-D01 design proposal, not an implemented or accepted
contract. M01 acceptance and the remaining D01/D02/D03 approvals remain required
before dependent implementation. A merge alone does not change this status.

## Context and trade-off

[UX-04](../ux/UX-04.md) starts with creating a record from zero inventory, then
enrolling its agent later. A display name/site/model therefore cannot serve as
proof of identity. The [source inventory](../research/M02-D01a-sources.md) traces
the present foundation: control probes, a metadata-only baseline migration,
an idle agent and disabled remote access provide no enrollment authority.
The [security lifecycle review inputs](../../security/README.md#proposed-credential-and-token-lifecycle--unimplemented)
already separate device-bound bootstrap, identity and per-session authorization.

We propose keeping the administrative record and the security identity distinct,
and treating their association as a separately authorized, reconciliable
transition. This supports staging inventory before installing an agent, at the
cost of an additional state boundary and explicit recovery/replacement work.
Conflating these concepts would make later recovery and access history difficult
to interpret. The [glossary](../../GLOSSARY.md) gives their domain meanings.

## Proposed decisions

All requirements in this section are proposals for approval. They specify domain
constraints, not routes, database tables, field names, wire states or algorithms.

1. **Creating a record establishes metadata only.** It supplies no device
   credential, enrollment, connectivity observation or SSH permission. Existing
   record names may coincide; metadata, network address and local config checks
   cannot authenticate an agent. Installation and enrollment remain separate.
2. **Enrollment requires independent authority and device proof.** The control
   boundary must authenticate the operator, authorize enrollment for the exact
   intended record/resource scope, and verify the intended agent's identity.
   View/create rights do not imply enrollment authority; enrollment authority
   does not imply permission to connect or administer other resources. A relay,
   browser assertion or claimant-supplied record ID is not sufficient authority.
3. **Bootstrap must bind to the intended device before acceptance.** Propose a
   device-bound, one-use enrollment capability, restricted to its intended
   record and verified device identity. A design accepting whichever claimant
   arrives first is insufficient. The trusted way to establish that binding,
   proof mechanism, secure delivery and capability expiry/revocation are
   unresolved approval blockers below. No enrollment command or secret is
   defined here. Issuing and consuming that authority both require protected
   authorization and the approved audit foundation.
4. **A binding has an unambiguous current target.** Propose at most one current
   device identity per record and at most one current record per device identity
   within the deployment's authoritative enrollment domain. A renamed record,
   reinstallation, competing claim or discovered device cannot silently replace
   that association. Cross-domain identity uniqueness and moves remain unresolved;
   this does not select a tenant model or claim physical clone resistance.
5. **Consumption and binding must be recoverable as one logical outcome.**
   Authoritative processing must prevent two successful conflicting bindings or
   reuse of consumed enrollment authority. A committed outcome survives retries
   and process restart; a lost response does not justify a new enrollment.
   Reconciliation requires current authorization to the original record and
   operation and returns its authoritative outcome, without disclosing a
   different operator's device or credential. Exact correlation, persistence,
   transactions and audit coupling require the follow-up designs below.
6. **Replacement is an explicit security operation.** A copied credential,
   factory reset or replacement agent cannot acquire an existing record from
   matching metadata. Recovery must establish fresh authority and the intended
   identity; historical associations must remain attributable. Handling the old
   credential, active sessions, pending enrollment authority and audit continuity
   requires D02 approval before replacement is exposed. A software identity proof
   alone cannot distinguish a device from a clone holding identical secret
   material; do not claim hardware authenticity without separate evidence.
7. **Enrollment success does not authorize SSH.** Fresh per-session checks and
   independent SSH host-key trust remain required under
   [ADR-001](ADR-001-foundation.md). A device identity, its credential and its SSH
   host key are distinct concepts; binding one does not verify the others.

## Alternatives considered

| Alternative | Why it is not the proposed choice / consequence |
| --- | --- |
| Create a record only after enrollment | Reduces placeholder records, but cannot support the owner's create-from-zero then enroll journey without conflating its stages. |
| Treat a record ID or name as identity | Simpler onboarding, but a known identifier or duplicate metadata cannot prove the intended device; risks wrong-device binding. |
| Shared fleet bootstrap or first claimant wins | Reduces provisioning effort but does not establish intended-device binding; shared fleet credentials also violate ADR-001. |
| Allow multiple current identities for one record | Could support multi-agent devices, but makes the authorized SSH target and replacement ambiguous. Such a future need requires its own reviewed model. |
| Proposed distinct record, identity and authorized binding | Preserves staging and attribution; adds reconciliation, conflict detection and replacement design costs. Credential scheme and secure bootstrap remain undecided. |

## Boundary cases and required future evidence

These are requirements to test after approved designs and scoped implementation;
they are not executed tests or selected API errors/status enums.

| Case | Proposed required outcome / evidence |
| --- | --- |
| Unauthenticated operator, viewer or record creator without enrollment rights | Direct control requests cannot issue authority, bind or query protected operation outcomes; no UI-only enforcement. |
| Claimant submits another record ID or copies a display name | Reject intended-record/device mismatch even with otherwise valid proof; no cross-resource disclosure. |
| Replayed, expired, revoked or stolen enrollment capability | Reject reuse/invalid authority; possession without the approved intended-device proof cannot establish a binding. Time and revocation decisions require D02. |
| Concurrent claims, including attempts against different records | Enforce both uniqueness constraints and one-use consumption; competing claims do not overwrite the winner. Exact duplicate retries reconcile the original operation without a second binding. |
| Lost response, cancellation or timeout | Treat result as uncertain until authorized reconciliation. Cancellation does not undo a committed outcome; absence of a local response/audit display is not proof of failure. |
| Crash before/after commit or audit producer/store failure | Approved D02 reconciliation distinguishes rejected, committed and uncertain outcomes; never acknowledge a completed enrollment with incomplete required audit evidence. Atomicity/recovery must prevent authority reuse. |
| Record removed/revoked or authority changed while enrollment is in flight | Recheck the authoritative resource/authority at the accepted transition boundary; prevent an orphaned binding or resurrected access. Exact race/time rules require follow-up approval. |
| Reset, restored backup, replacement or copied credential | No silent rebind; retain attribution and invalidated-authority decisions. Document software-only clone limitations and any hardware-backed evidence. |
| Enrolled but never observed or SSH not authorized | Report enrollment separately from observation and connection authority; no connected status or session grant inferred from the binding. |

## Remaining decisions and implementation gates

[M02's execution plan](../M02_EXECUTION.md#dependency-graph-and-accountable-evidence)
permits D01 review inputs before P00. This slice completes neither all D01 design
nor WEB-01's agreed identity/authorization prerequisite. Proposed owner roles
must be assigned accountable people by the Tech Lead.

| Gate / proposed owner | Required decision before dependent implementation |
| --- | --- |
| P00 — Human Tech Lead / Hardware QA / repository administrator | Explicit M01 acceptance, ADR-001 disposition, hardware and governance evidence; current merge history is not acceptance. |
| Remaining D01 — Architect / Security / Agent / Backend | Trusted intended-device bootstrap and recovery channels, operator identity trust, authorization scope/delegation, deployment identity domain, uniqueness/clone assumptions, enrollment proofs and secure delivery. Compare credential/crypto/storage mechanisms and approve concrete choices separately. |
| D02 — Security / Backend / Agent | Credential issue/store/rotate/revoke/decommission; trusted time/lifetimes; in-flight changes, replacement/active-session effects, uncertain-result correlation/retention, transaction/recovery and audit producer/store coupling. Decide safe behavior on control/store/audit outage. |
| D02/D03 — Security / Agent / Relay / CLI | Bounded attempts/handshakes/pending work, retries, buffers and flash writes; approved measured budgets, authenticated transport and independent SSH host-key bootstrap/change/recovery. |
| Full WEB-01 — Backend / Security / Architect | After agreed M02 identity/authorization: API/schema/errors, record creation/query scope, authenticated web sessions, mutation concurrency/reconciliation and audit contracts. UX fields and local demo IDs are not these contracts. |

Approving this ADR would approve only its proposed domain separation, authority,
cardinality and logical outcome constraints. It would not approve an enrollment
mechanism or enable implementation while any necessary gate remains unresolved.
No provider, credential encoding, cryptographic algorithm, transport frame,
endpoint, database schema or numerical resource/security policy is selected here.

## Impact, validation and rollback

This change is documentation only. Future implementation adds authoritative
binding/reconciliation/audit storage and verification costs that must be bounded
and measured; it must not introduce device web/database runtimes or assume root.
The 128 MB total RAM and <=16 MB agent idle RSS target from the
[README](../../README.md#design-goals) require real ARM measurements. Protected
identity storage, flash writes, proof overhead and clone resistance remain design
and measurement obligations, not achieved guarantees.

Executed document/source/scope and secret-scan evidence belongs in
[validation.txt](../evidence/m02-d01a/validation.txt) and
[secrets.txt](../evidence/m02-d01a/secrets.txt). No application build/test,
live enrollment, cryptographic/security acceptance, browser focus/accessibility,
Windows/ARM toolchain or physical ARM validation is claimed. BrowserOS neo MCP
tools were unavailable; repository source research is not external protocol
verification or a fresh hardware/governance audit. Pattern scans cannot certify
absence of every secret.

Rollback this proposal via a reviewed documentation revert/superseding ADR and
Draft PR into `develop`, preserving decision provenance. There is no deployed
binding, credential, service or schema to recover. Future implementation rollback
must preserve authoritative history and invalidated credentials; reverting code
must not re-enable a consumed capability or revoked identity.

## Human decision record

**PENDING; intentionally unsigned.** Tech Lead records accept/revise/reject,
reviewed ADR commit/PR, precise decision scope, remaining gate owners and an
authoritative decision reference with identity and UTC time. No approval is
inferred from the owner's instruction to continue, agent review or a merge.
