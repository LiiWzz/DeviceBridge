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
cleanup. A missing/uncertain outcome blocks a conflicting action that creates,
extends or restores authority until reconciled under current authority.

**Invalidation takes precedence over unresolved issuance or rotation.** An
independently authorized revoke or decommission request must not wait for an
uncertain rotation to finish. It must establish an authoritative denial for the
affected identity/grant scope covering both existing credentials and any pending
replacement authority. A late rotation result, retry, restart or reconciliation
must preserve that denial and cannot activate a replacement or revive retired
material. Decommission's invalidation step takes precedence; destructive cleanup
and removal of attribution still require safe reconciliation and retention rules.

This is no outage or audit bypass: invalidation still requires current authority
and the required durable outcome/audit evidence. If those dependencies prevent
confirmation, report the outcome as unknown, withhold new access while validity
is uncertain, and never claim revocation completed. D02/D03 must approve the
serialization, durable denial, propagation and recovery mechanisms before
implementation; active-session behavior and numerical deadlines remain open.

The system needs a trusted time basis and explicit skew/rollback behavior before
expiry can be meaningful. Decisions must account for clock rollback, reboot,
offline devices, restore from backup and control-plane partitions. Relative
timeouts alone do not establish trusted credential age across processes or
devices. Source, synchronization, monotonic/wall-clock roles and safe behavior
remain unselected pending design evidence.

## Deployed-device continuity and recovery proposal

Follow-up [#55](https://github.com/LiiWzz/DeviceBridge/issues/55) captures the
owner's selected direction: devices already installed must have an engineered
remote recovery path. This is selection of a design direction, not acceptance of
this ADR, a recovery credential scheme or a guarantee of perpetual reachability.

1. **Agent renewal runs independently of an SSH session.** Schedule authenticated
   renewal sufficiently before expiry, with bounded retry/backoff and early
   operator alerts. A user need not connect first. Define server-side detection
   of missed renewals, escalation owners and actionable expiry/recovery status;
   an offline agent cannot be relied on to report its own impending expiry.
2. **Separate agent identity from operator SSH authority.** Expiry or revocation
   of a user's session grant must not expire the agent's independently valid
   maintenance identity. Agent connectivity never grants the user SSH access.
   Device-identity revocation or decommission still blocks normal agent access;
   this separation must not evade intentional invalidation.
3. **Rotate with verified durable replacement.** Confirm the intended device can
   authenticate with its securely installed replacement before retiring the old
   credential, within an explicitly bounded overlap. Reconcile interrupted
   install/acknowledgement across power loss and reboot without unbounded overlap
   or choosing old material just because a backup contains it. If replacement
   cannot be verified before expiry, raise an explicit recovery-required state;
   neither a silent lifetime extension nor a permanent fallback is permitted.
4. **Expired credentials do not grant normal access.** A device returning after
   the approved offline window must use a separately authenticated, narrowly
   scoped recovery flow. The expired credential, record ID or network address
   alone is no recovery proof. Recovery authorizes credential repair only, not
   shell, arbitrary forwarding or user impersonation. It remains outbound from
   the device through an authenticated, bounded path; no unauthenticated listener
   or custom device shell protocol is introduced.
5. **Independent recovery trust must be designed and provisioned before rollout.**
   Establish intended-device proof and independently authorized operator/control
   approval, bind the recovery to its identity and operation, restrict attempts,
   and make authority short-lived/one-use as appropriate to the approved scheme.
   Specify secure provisioning, storage, renewal, rotation, revocation and loss
   of the recovery trust itself, including the bootstrap of server trust and
   trusted time when the normal credential/time source cannot be used. No
   evergreen emergency key, default fleet secret, acceptance of expired normal
   credentials, or circular dependency on an expired credential is chosen here.
6. **Recovery cannot revive invalid authority.** Reject automatic recovery of a
   revoked/decommissioned identity and suppress pending replacement authority.
   Compromise recovery, if needed, is a separately authorized replacement/rebind
   decision under ADR-002 with fresh trust, audit and explicit disposition of old
   authority. It must not silently reactivate the revoked identity. Replay,
   concurrent recovery, lost responses and late rotation results reconcile one
   operation under current authority and preserve invalidation precedence.

Before deployment, Security/Agent/Backend/DevOps must select the recovery proof,
custody, protected delivery, trusted server/time bootstrap, authoritative outcome
and audit coupling, and demonstrate the supported offline/recovery envelope.
Tech Lead must approve maximum offline duration, renewal lead time/lifetimes,
bounded overlap/retries/recovery attempts, server-side alert/escalation thresholds,
revocation propagation, flash budget and recovery service ownership. These are
explicit rollout blockers; a plan that simply leaves deployed devices to expire
does not satisfy this proposal. No numeric values or mechanisms are selected.

QA must later demonstrate user expiry while a valid agent renews, loss of
connectivity around renewal/rotation, power loss before/after installation and
acknowledgement, return within/beyond the offline window, expired identity with
valid versus missing/revoked recovery proof, wrong-device/replayed recovery,
clock rollback/no trusted time, audit/control/recovery-service outage, compromise
during uncertain rotation, decommission and backup restore. Record real ARM
flash/RAM/CPU costs and server aggregate retry/recovery load, not just success.

Remote recovery cannot work without surviving device proof, power, an eventual
network route, trusted service/time bootstrap and available operator authority.
Loss of these may require physical service or a separately reviewed independent
management path. Product/operations owners must explicitly accept those limits
and the response plan before shipping; no design can promise access under every
failure. Revocation, audit, credential secrecy and host-key checks stay required.

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
events, disk full, redaction and unauthorized audit reads. Include lost rotation
response followed by compromise/revoke or decommission: existing and pending
replacement credentials must remain denied through late responses, retries,
restart and restore, with authorized reconciliation and durable audit attribution.
Hardware QA measures
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
Deployed-device recovery follow-up checks are in
[validation.txt](../evidence/m02-recovery/validation.txt) and
[secrets.txt](../evidence/m02-recovery/secrets.txt).

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
