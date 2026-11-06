# ADR-005: Propose credential lifecycle and authoritative audit outcomes

**Status: Proposed — Tech Lead approval required.** Issue [#51](https://github.com/LiiWzz/DeviceBridge/issues/51).
Source baseline: `4e74760eaf7080f07bdb01f1740a8d94c1ac6b0f`.

Propose treating credential changes and access decisions as recoverable, attributable operations whose authoritative outcome is durably auditable before success is exposed. This costs storage, availability and coordination during outages. It avoids reporting access as successful when the credential/grant state or required audit evidence is uncertain. The owner approved ADR-004's limited session-authority domain constraints in the project conversation; the provenance and scope are recorded in [PR #53](https://github.com/LiiWzz/DeviceBridge/pull/53). No lifecycle, storage or audit mechanism is approved here.

## Proposed lifecycle

| Operation | Proposed domain rule | Decisions still open |
| --- | --- | --- |
| Issue / enroll | Separate operator authorization, intended-device proof, credential creation and identity binding. Never accept a first claimant or shared fleet credential. Do not expose success until required audit outcome is durable. | Trusted bootstrap/recovery, proof, credential/verifier format and custody; consistency boundary with ADR-002/P02. |
| Renew | Authenticate the current identity and renew only for the same device/domain; renewal does not change SSH account or grant session permission. | Trusted time, lifetime, skew, offline behavior and retry/overlap bounds. |
| Rotate | Establish and verify a replacement before retiring the old credential; bound any overlap. A crash or lost response reconciles the original rotation and cannot revive retired material. | Atomicity, overlap, flash-write budget, backup/restore and proof of replacement. |
| Revoke | Invalidate the scoped user/device/grant credential for new authorization; preserve attribution. Never imply an active SSH process is immediately stopped. | Propagation SLA, active-session input/closure behavior, cache/outage semantics and maximum delay. |
| Decommission / recover | Revoke pending authority and credentials, record the outcome and retain required history. Recovery cannot reactivate revoked identity or erase attribution. | Retention, deletion/backup policy, legal/operational holds and hardware reset/replacement flow. |

## Time, uncertainty and outages

All operation decisions must distinguish **accepted**, **rejected**, and
**unknown pending reconciliation**; a timeout, missing local audit display or
disconnected client is not proof of failure. Retries bind to the original actor,
resource and operation, cannot consume authority twice, and disclose no other
principal's result. Cancellation cannot undo an accepted operation.

Expired, revoked, unverifiable or wrong-device credentials fail closed for new
admission and attachment. Do not invent a grace period or fall back to local
cached authority when control, trusted time or required audit storage is
unavailable. Existing session handling must be explicitly chosen and measured;
this ADR does not promise immediate termination, command rollback or descendant
cleanup. A missing/uncertain outcome blocks a new conflicting lifecycle action
until reconciled under current authority.

The system needs a trusted time basis and explicit skew/rollback behavior before
expiry can be meaningful. Decisions must account for clock rollback, reboot,
offline devices, restore from backup and control-plane partitions. Relative
timeouts alone do not establish trusted credential age across processes or
devices. Source, synchronization, monotonic/wall-clock roles and safe behavior
remain unselected pending design evidence.

## Audit obligations

Audit must preserve attributable attempted, denied, accepted, rejected and
uncertain/reconciled outcomes for credential issue/renew/rotate/revoke,
enrollment binding and session authorization. Each event links the authenticated
actor, affected identity/resource, operation correlation, outcome and trusted
time context without recording secret material. Define who may read/export it,
retention/redaction, integrity/access controls and recovery independently from
ordinary application logs. Never record SSH private keys, bearer credentials,
terminal keystrokes, pasted commands or screen output by default.

Required audit storage failure must not be converted into an unlogged success.
Decide per operation whether to reject before mutation or durably reconcile a
committed mutation; no response may claim completion until the authoritative
state and required audit outcome agree. Producer restart, duplicate delivery,
partial write, backup/restore, disk full and audit-store outage need explicit
recovery evidence. Audit completeness must not be inferred from a UI history
view or process log.

## Alternatives and consequences

| Alternative | Consequence |
| --- | --- |
| Best-effort event logging after mutation | High availability and simple flow, but can leave an irreversible credential change with no authoritative attribution after crash/outage. Not proposed. |
| Block every security mutation unless required audit is durable | Clear fail-closed evidence, but reduces availability and requires a reviewed atomicity/reconciliation design. Proposed domain requirement, exact implementation remains open. |
| Local device logs as authoritative lifecycle history | Works during disconnection, but cannot alone enforce fleet-wide revocation or ensure server/session decisions share one outcome. Not proposed as the sole authority. |

## Future evidence and owners

Security/Backend/Agent define approved mechanisms and named owners before
implementation. QA must later cover concurrent issue/renew/revoke, expired and
replayed credentials, rotation crash at every boundary, clock rollback/skew,
restore/decommission, lost responses, control/audit outage and recovery, duplicate
events, disk full, redaction and unauthorized audit reads. Hardware QA measures
flash writes and endurance on supported ARM devices. Values for lifetimes,
overlap, revocation delay, retries, storage retention and resource budgets require
Tech Lead-approved measured profiles; none are selected here.

P00/M01 acceptance and ADR-001 approval, agreed remaining D01 including ADR-002
and ADR-004 scope, and approved D02 decisions precede dependent implementation.
D03 still owns authenticated transport, framing, host-key lifecycle and bounded
relay/session costs. P02 audit foundation precedes sensitive mutations. Full
WEB-01 and browser-specific decisions remain separate prerequisites.

## Impact, rollback and decision record

Executed documentation checks and limitations are captured in
[validation.txt](../evidence/m02-d02/validation.txt); scan commands/output are in
[secrets.txt](../evidence/m02-d02/secrets.txt). Pattern scans cannot certify absence
of every secret.

This proposal changes documentation only. Future trusted time, durable audit,
credential state and reconciliation add flash/storage, memory, CPU, latency and
availability costs. Device database/web runtimes, root and unmeasured flash
churn remain prohibited assumptions. No live browser, crypto, credential,
application build, Windows/ARM toolchain, physical ARM or hosted CI evidence is
claimed. Revert/supersede through a reviewed PR; there is no deployed state to
recover now. Future rollback must preserve revocation, operation history and
audit integrity; reverting code cannot undo remote commands or restore secrets.

**Human decision: PENDING; intentionally unsigned.** Tech Lead records exact
reviewed commit, accept/revise/reject, decision scope, remaining owners, identity
and authoritative UTC reference. Merge and agent review are not approval.
