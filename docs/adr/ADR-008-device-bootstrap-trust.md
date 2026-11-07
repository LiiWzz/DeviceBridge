# ADR-008: Propose independently established device bootstrap trust

**Status: Accepted (domain principles only); mechanisms, budgets and implementation gates remain pending.**
Issue [#61](https://github.com/LiiWzz/DeviceBridge/issues/61).
Repository source baseline: `77976bdf60270fe187393155c4856833432c3fd4`.

Propose establishing trusted evidence of the intended device independently of
its enrollment claim, then requiring fresh proof bound to the authorized
enrollment operation. This adds provisioning, custody and reconciliation work,
but prevents a copied record identifier or stolen enrollment capability from
choosing the device that becomes trusted. This remaining D01 design input
selects no credential, cryptography, hardware, provider or protocol mechanism.

## Context and alternatives

[ADR-002](ADR-002-device-record-identity-binding.md) proposes distinct records,
identities and one-use device-bound enrollment authority, leaving how intended
device trust is established unresolved. [ADR-007](ADR-007-operator-trust-and-authority.md)
proposes independent operator trust and scoped authority. The limited accepted
domain principles in [ADR-004](ADR-004-session-authority-binding.md) and
[ADR-005](ADR-005-credential-lifecycle-audit.md) keep session permission separate
and require durable outcomes, invalidation precedence and independent recovery.
[ADR-006](ADR-006-transport-trust-and-bounds.md) proposes distinct TLS peer,
device and SSH host-key trust; none can be inferred from another.

At the baseline, [agent `run`](../../agent/src/main.rs) idles without networking
or disk writes, while [control `service.handle`](../../server/control/internal/control/server.go)
serves health/readiness probes only. The [security baseline](../../security/README.md)
lists future bootstrap/lifecycle review obligations. These are repository source
observations, not executed enrollment tests or external protocol facts. No
provisioning service, protected key store or hardware authenticity exists by
virtue of this document.

| Alternative | Trade-off |
| --- | --- |
| Trust names, serial labels, network location or whichever claimant presents an enrollment capability first | Cheap installation, but establishes no independent intended-device proof and permits substitution. Not proposed. |
| Shared fleet bootstrap credential | Simplifies distribution, but one compromise crosses devices and violates ADR-001. Not proposed. |
| Proposed independently verified per-device provisioning and operation-bound proof | Supports explicit intended-device binding and attributable recovery; adds trusted provisioning and per-device custody costs. Concrete evidence/delivery mechanisms remain open. |
| Require hardware-backed device evidence universally | May improve extraction/clone resistance, but depends on board capabilities, supply-chain trust and verifier lifecycle. No such support is established for mixed Buildroot/Yocto/Debian devices; assess separately rather than claim universal hardware authenticity. |

## Proposed domain constraints

1. **Establish intended-device trust before accepting a claim.** An independently
   authenticated provisioning authority must establish which device proof belongs
   to the intended record and authoritative enrollment domain through a reviewed
   verification/delivery path. A claimant's new key, asserted serial number,
   installation log, browser confirmation or relay assertion alone cannot create
   that trust. An operator may approve enrollment only within current scoped
   authority; approval cannot turn an unverified claimant into the intended
   device. Record evidence provenance and the accountable provisioning actor.
   Who verifies the physical or supply-chain association and how that evidence
   reaches the verifier are explicit unresolved approval blockers.
2. **Bootstrap trust is mutual and precedes normal credentials.** The device
   must verify the intended service and deployment domain before disclosing
   protected proof or accepting installed trust; the service independently
   verifies the intended device. A reachable TLS endpoint alone cannot establish
   the device's identity. Initial service anchors, device verification material,
   trusted-time bootstrap and authenticated updates need a protected path that
   does not depend on a credential the enrollment is about to issue. Never
   disable certificate/time checks or accept a first observed peer to escape
   this dependency. No TLS mode, anchor format or clock mechanism is selected.
3. **Separate device proof from enrollment permission.** Proof establishes the
   claimant's relationship to independently trusted device material; enrollment
   authority permits one exact binding. Require both at the trusted control
   boundary, scoped to the current operator authority, deployment domain,
   intended record, verified identity and original enrollment operation. Bind
   fresh proof to that same context and intended service so captured evidence,
   relayed claims, cross-domain requests or another operation cannot substitute
   the target or reuse acceptance. Challenge/transcript design, freshness,
   cryptography and any channel-binding mechanism require separate approval;
   this is not a claim to defeat every live relay or physical cloning attack.
4. **Protect provisioning and custody throughout their lifecycle.** Approve
   per-device proof generation, delivery, verification, installation, storage,
   activation, update, revocation and destruction with least privilege. Do not
   place secrets in images shared across the fleet, arguments, URLs, logs, audit
   or repository configuration. Distribution intermediaries and installers
   cannot mint broader enrollment authority. Substituted material, incomplete
   installation or missing evidence blocks activation; reconcile power loss and
   lost acknowledgement before retrying. Compromised provisioning authority,
   stolen device material or a restored image requires attributable invalidation
   and separately authorized fresh trust, not another first-claim attempt.
   Software possession alone cannot distinguish two devices holding identical
   copied secrets; define supported extraction/clone assumptions and any
   hardware evidence before making stronger guarantees.
5. **One-use consumption and binding form one durable logical outcome.** Follow
   ADR-002's one-current-binding constraints and ADR-005's audit obligations.
   Recheck current enrollment authority, intended-device trust, resource binding
   and validity at the approved transition boundary. Consume permission and
   establish the binding without exposing two successful conflicting outcomes.
   Repeated delivery reconciles the original operation under current authority;
   a lost response, timeout or local cancellation is no proof of rejection and
   cannot justify reissuing authority or enrolling a different device. Preserve
   attributable provisioning evidence and required durable enrollment audit
   before reporting success. Atomicity, correlation, retention and crash/audit
   recovery mechanisms remain open.
6. **Invalidation and uncertainty cannot become bootstrap exceptions.** Expired,
   revoked, consumed, wrong-device or unverifiable authority cannot admit a new
   binding. Required control, trusted-time, provisioning verification or audit
   outages yield no stale-cache allow or unlogged success. Report unknown
   outcomes explicitly and withhold conflicting grants pending reconciliation.
   Independently authorized revoke/decommission still takes precedence over
   pending enrollment or replacement; late proof, retries, restart and backup
   restore cannot activate invalidated authority. Do not claim completed
   invalidation without its required durable evidence. Approve propagation,
   serialization and any active-session consequences separately.
7. **Recovery and replacement establish fresh authority, not ambient access.**
   Initial enrollment, renewal and independently authenticated recovery are
   distinct operations. Under ADR-005, design/provision surviving recovery trust
   before rollout; neither expired normal credentials nor an enrollment receipt
   alone authenticate recovery. Recovery repairs scoped credentials only and
   cannot grant SSH, arbitrary forwarding or impersonation. A replaced or
   compromised device requires ADR-002's explicit rebind, fresh independently
   verified intended-device trust and disposition of old/pending authority.
   Retain historical attribution and invalidation. If no independent proof,
   trusted service/time path or authorized owner survives, require physical
   service or a separately reviewed independent management path; no permanent
   emergency credential or promise of universal remote recovery.

Enrollment remains separate from connectivity, per-session authority, SSH login
and independently verified host keys. This proposal neither binds host keys nor
selects a root-of-trust product, credential encoding, API/schema, wire format,
proof algorithm or numerical lifetime/resource/security policy.

## Future evidence and accountable decisions

These are future execution requirements, not passing tests. Tech Lead assigns
named accountable people to the proposed roles and approves decisions first.

| Area / proposed owners | Decision and negative evidence required |
| --- | --- |
| Intended-device provenance — Security/DevOps/Agent | Select authenticated provisioning actor, physical/supply-chain verification and evidence custody. Reject forged labels/record IDs, substituted new keys, compromised installer/provisioner and missing provenance. Establish limits for copied software proof and hardware-backed claims. |
| Mutual bootstrap — Security/Backend/Agent | Select initial service/domain trust, protected delivery/update and trusted time without circular normal-credential dependencies. Test wrong service/domain, untrusted peer, substituted anchors, clock rollback/no trusted time and partial provisioning. |
| Proof and permission — Security/Backend/QA | Approve fresh context-bound proof and one-use authority. Test stolen enrollment permission without device proof, valid device proof without operator authority, wrong device/record/domain/operation/service, replay and relayed/substituted evidence under the stated threat model. |
| Durable binding — Backend/QA/Security | Select serialization, correlation, evidence/audit coupling and current-authority reconciliation. Exercise duplicate/concurrent claims across records, resource rebind/operator disablement during enrollment, power loss before/after install/commit, lost responses and audit/control/store outages; no second binding or unauthorized outcome disclosure. |
| Invalidation/recovery — Security/Agent/Backend/DevOps | Select revocation of provisioned trust and pending authority, replacement and independent recovery provisioning. Test revoke/decommission during unknown enrollment, late proof, expired/revoked/missing recovery evidence, restored images/backups, provisioner compromise and loss of all trust. Preserve denial and history; record physical-service limits. |
| Bounds/custody — Hardware QA/Agent/Backend | Approve proof sizes, parsing/verification work, unauthenticated admission, pending operations, retries, durable retention and flash budgets. Measure slow/malformed/oversized proofs and 100-device provisioning/recovery bursts with raw CPU/RAM/latency/write evidence; bound work before expensive verification and never spill backlog to device flash. |

P01 must specify protected provisioning prerequisites; P02 audit foundation
precedes enrollment/security mutations. Initial provisioning before ordinary
audit must use separately approved durable evidence and reconcile it before
activation, as proposed in ADR-007. There is no unaudited bootstrap exception.
P03 follows P01/P02 and approved D01/D02; this design input bypasses none of them.

## Impact, gates, rollback and decision record

Documentation only: no runtime/prototype, dependency, listener, production config,
credential, allocation, device flash write or root use changes. Future proof,
protected storage, audit/reconciliation and provisioning add CPU/RAM/storage,
latency and operational costs. Avoid unnecessary allocations/dependencies and
flash churn on ARM; do not assume root or device web/database runtimes. Measure
real ARMv7/AArch64 Buildroot/Yocto/Debian profiles against **128 MB total RAM,
agent idle RSS target <=16 MB**, reporting actual units and active/bootstrap
overhead. Cross-builds and macOS ARM hosts do not establish target acceptance.

ADR-001/002/006/007 remain Proposed; ADR-004/005 retain only their recorded limited
domain acceptance. Remaining D01/D02 mechanisms, D03 contracts/budgets,
P00/M01 acceptance and ordered P01–P09 gates remain open. This completes no
M01/M02, full WEB-01 or browser-terminal approval;
[ADR-003](ADR-003-browser-terminal-scope.md) and its separate
architecture/UX/production gates are unchanged.

Executed document/link/scope/whitespace checks and redacted scan output are in
[validation.txt](../evidence/m02-d01d/validation.txt) and
[secrets.txt](../evidence/m02-d01d/secrets.txt). External URLs/fragments and
external protocol/provider facts are unverified. No application build/lint,
live enrollment/authentication/crypto/security, SSH/PTY/browser, Windows/ARM
toolchain, physical ARM or hosted CI pass is claimed. Pattern scans cannot
certify absence of every secret.

Revert/supersede documents through a reviewed feature branch and Draft PR into
`develop`; no deployed credentials or state need recovery now. Future rollback
must preserve consumed authority, invalidated trust, original binding/outcomes
and audit history; restoring an image/database cannot revive withdrawn material.
Code rollback cannot restore stolen secrets or undo completed remote commands.

**Human decision: accepted limited domain principles.** In the project
conversation on 2026-11-07 (Asia/Bangkok), the owner replied “confirm and next”
after the summary of [PR #62](https://github.com/LiiWzz/DeviceBridge/pull/62),
reviewed text at `44fd7fa18527ea3dc3583e1e4402e2cafbfc86a8`. This records approval
of the proposed domain constraints: independent intended-device trust, mutual
bootstrap, separate device proof and enrollment authority, protected per-device
provisioning/custody, durable one-use binding/reconciliation, invalidation
precedence and independently authenticated recovery with physical-service limits.
The original proposal above is retained as the reviewed text. No signing
identity, exact approval time or external signature was captured; none is
invented here. This is the owner's conversation confirmation, not an agent
signature or approval inferred from merge.

Proof/credential/crypto/hardware/provider choices, service/time bootstrap,
delivery/storage, protocol/API/schema, consistency/audit mechanisms, numerical
lifetimes/resource budgets and implementation remain unapproved. ADR-001/002/
006/007, remaining D01/D02/D03, P00/M01/M02 and browser gates stay open under
[AGENTS.md](../../AGENTS.md). This approval does not authorize implementation.
Approval-record validation/scan evidence is linked from
[ADR-009](ADR-009-active-session-validity.md); original proposal evidence above
remains historical and is not evidence of this later status change.
