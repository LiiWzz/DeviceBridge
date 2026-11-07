# ADR-010: Propose one authoritative commit for security outcomes and audit

**Status: Proposed — explicit human Tech Lead approval required.**
Issue [#65](https://github.com/LiiWzz/DeviceBridge/issues/65).
Repository source baseline: `b25d382076facb1f7d4e68da79f392ed93a6caae`.

Propose one authoritative server-side commit boundary for a protected state
change, its operation outcome and required audit record, with read/export views
derived separately. This adds storage coordination and reduces availability when
required durable evidence cannot be established, but avoids a successful change
whose attribution exists only in a process log or undelivered queue. This D02/P02
architecture proposal remains unapproved and selects no SQL schema, API/wire
contract, provider, library, cryptography or numerical policy.

## Repository context and alternatives

[ADR-005](ADR-005-credential-lifecycle-audit.md) accepts durable attributable
outcomes and invalidation precedence but leaves consistency/storage mechanisms
open. [ADR-002](ADR-002-device-record-identity-binding.md) proposes one-use
binding/reconciliation; [ADR-004](ADR-004-session-authority-binding.md) accepts
exact session binding. [ADR-007](ADR-007-operator-trust-and-authority.md) proposes
current-authority outcome disclosure and protected provisioning evidence.
[ADR-008](ADR-008-device-bootstrap-trust.md) accepts limited device-bootstrap
principles. [ADR-009](ADR-009-active-session-validity.md) now records limited
owner agreement with stop/closure policy, without approving mechanisms/budgets.

The current [control handler](../../server/control/internal/control/server.go)
serves probes only. The [baseline migration](../../deploy/migrations/0001_baseline.up.sql)
contains metadata and no application tables; [migration documentation](../../deploy/README.md)
describes a development migration runner, not a production security/audit store.
The [security baseline](../../security/README.md) describes future lifecycle
obligations. These are source observations, not executed storage or security
tests. No audit producer, durable operation journal or projection is implemented
by this proposal; no external database/provider guarantee is asserted.

| Alternative | Trade-off |
| --- | --- |
| Commit mutation, then write audit independently on a best-effort basis | Simple and available during audit outages, but a crash can lose the only attributable outcome. Not proposed. |
| Place authoritative mutation and audit in separately coordinated stores | Can separate custody/failure domains, but requires an approved cross-store commit/recovery design; an exporter retry alone cannot make the original writes atomic. Not the proposed initial boundary. |
| Proposed shared authoritative commit, derived read/export views | Keeps control-state outcome and required evidence together, at the cost of one availability/custody boundary. Separate projections need duplicate/recovery handling; stronger independent integrity must be designed explicitly. |

## Proposed commit and reconciliation boundary

1. **Required audit evidence is an authoritative record, not an intent.** For a
   protected server-side mutation, require the current authorized state decision,
   original operation attribution and required audit outcome to commit as one
   logical unit. A queued intention to produce audit later is insufficient.
   The record itself must satisfy approved durability, integrity, access and
   retention requirements before completion or dependent usable access is
   exposed. Read/export projection delivery is separate: neither an exporter ACK
   nor a populated dashboard proves the authoritative commit. Exact transactional
   storage, acknowledgement semantics and failure model require approval and
   executed evidence; do not assume the development database meets them.
2. **Bind retries to the original operation without granting rights.** Establish
   correlation scoped to verified actor/domain, exact action/resource/current
   binding and the original intended change. Reject substitution of different
   contents or targets under the same operation reference; caller-selected IDs
   cannot confer access or collide across principals. Concurrent duplicate
   requests must resolve to one authoritative state transition and attributable
   result. Missing response/cancellation is no proof of rollback. Reconciliation
   requires current scoped authority before disclosing protected results and
   never reexecutes a committed grant merely to recreate a response. Original
   attribution survives changes of operator; an independently authorized scoped
   responder may reconcile without impersonating the former caller.
3. **Separate committed control state from remote effects.** An issued rotation
   decision is not proof of device installation; a session-start decision is not
   proof of SSH attachment, and recorded invalidation is not proof every remote
   session has closed. Record subsequent verified effects or uncertainty as
   attributable follow-up outcomes linked to the original operation, rather
   than rewriting history. Approve remote proof/correlation/retry boundaries
   separately. This proposes neither a database transaction across devices nor
   exactly-once network delivery or execution. Gate usable credentials/access on
   all required proofs and durable audit, not only a local commit.
4. **Preserve invalidation before reconciling uncertain grants.** Current scoped
   revoke/decommission must not wait for a pending grant/rotation to finish.
   Its authoritative denial covers existing and pending authority; a late grant,
   restored snapshot, duplicate event or projection replay cannot activate it.
   Confirm only the state/effects whose required evidence is durable. If the
   authoritative store is unavailable, report unknown and withhold new access
   while validity is uncertain; do not claim completed revocation. Local traffic
   stop/cleanup under ADR-009 must not wait for successful audit-store delivery.
   Its bounded evidence/reconciliation failure path remains a separate approved
   obligation, not permission to continue forwarding or claim complete audit.
5. **Record attempts and denials without inventing a principal.** At trusted
   ingress/decision boundaries, distinguish observed attempt, authenticated
   denial, committed state outcome and later reconciliation/remote observation.
   Record the trusted producer and which actor/resource/time facts were verified;
   a forged subject or claimed device is unverified data, not authoritative
   attribution. Sanitize and bound untrusted input. Never include credentials,
   private material, full request bodies, terminal input/output or pasted secrets
   by default. A rejected request changes no protected state even when its audit
   cannot be persisted; report/alert the audit gap through the approved bounded
   failure path. Denial must not turn into access to make logging succeed, and
   unauthenticated floods must not exhaust mutation/audit capacity. Approve
   attempted/denied-event durability, admission and loss-accounting semantics;
   this design makes no universal claim to record every arriving packet.
6. **Derive views without making them authority.** Proposed bounded read/export
   projections reference stable authoritative outcomes and tolerate duplicate,
   delayed and reordered delivery with explicit gaps/checkpoints. Rebuild from
   authoritative evidence under approved retention/recovery rules; projection
   lag is not mutation failure and a missing row is not proof of rejection.
   Reads/search/export/reconciliation require current resource/action scope and
   protected cursor/result handling, even for original callers. No projection
   consumer may mint permissions or replay mutations. Worker bounds, pagination,
   serialization and projection recovery mechanisms remain unselected.
7. **Preserve provenance and define integrity limits.** Normal writers append
   attributable corrections/follow-up evidence rather than silently rewriting
   prior outcomes. Separate ordinary mutation, audit-read/export and retention/
   administration privileges. Append-only application behavior alone cannot
   prevent a privileged storage administrator or compromised writer from
   falsifying/removing history. Security must choose custody separation,
   integrity verification, tamper detection and recovery evidence before making
   stronger claims; no signing/hash-chain/WORM mechanism is selected here.
8. **Restore and retention cannot revive authority.** Approve coordinated
   recovery of state, correlation and authoritative audit against a demonstrably
   trusted recovery boundary. Missing, corrupt or mutually inconsistent records
   block activation of affected access until independently authorized repair;
   an internally consistent old snapshot alone cannot prove no later revocation.
   Recover subsequent invalidations/consumption before enabling access or keep
   affected authority denied. Audit retention, privacy redaction and deletion
   must not remove the denial/replay-prevention evidence needed for surviving
   credentials or pending operations. Define approved residual evidence and
   terminal operation handling before expiry/deletion; never treat a forgotten
   operation reference as permission to execute its grant again. Preserve
   attributable retention/repair decisions without storing prohibited secrets.

## Decisions and future failure evidence

Tech Lead assigns named people to these proposed roles. These are future
executions, not tests passed or mechanisms approved by this document.

| Area / proposed owners | Required approval and evidence |
| --- | --- |
| Commit/durability — Backend/Security/DevOps | Choose authoritative store, commit/isolation/acknowledgement and failure model. Crash before/during/after each write/commit/response; partial write, disk full and lost acknowledgement must yield no false completion or usable grant without required evidence. |
| Correlation and concurrency — Backend/QA/Security | Select request equivalence, operation scope, conflict serialization and scoped reconciliation. Exercise same reference/different payload, cross-actor/domain/resource collisions, duplicate concurrent delivery, original caller disabled, lost response and authorized responder; no duplicate mutation or outcome leakage. |
| Remote effects/invalidation — Agent/Backend/Relay/QA | Define producer trust, follow-up proof, ordering and denial precedence. Test lost device install ACK, late start/rotation after revoke, remote closure unknown and stop during audit outage; preserve authority denial without claiming distributed atomicity or remote completion. |
| Audit ingress/privacy — Security/Backend/DevOps | Approve attempted/denied/unauthenticated durability and bounded failure/loss reporting. Test forged attribution, hostile/oversized payloads, clock uncertainty, log injection, secret redaction, input/output exclusion, floods and producer/store outage. No allow on audit failure. |
| Views and integrity — Security/Backend/QA | Approve projection duplicate/order/gap semantics, scoped list/search/export/cursors, administration separation and tamper detection. Test stale/corrupt views, cross-resource queries, missing/altered authoritative evidence and compromised-writer/storage-admin limits. |
| Recovery/retention/bounds — DevOps/Security/Hardware QA | Approve trusted restore boundary, invalidation recovery, retention/redaction/deletion and residual replay-denial evidence. Test old/mismatched backup, missing later revoke, retired operation reuse, full queues/store and 100-device correlated events; measure raw server CPU/RAM/storage/latency and any device write/traffic impact. |

P01 protected provisioning prerequisites and P02 accepted audit foundation remain
ordered before enrollment/grant/policy mutations. Provisioning evidence created
before ordinary audit must use ADR-007/008's separately approved trust/custody
path and reconcile before activation; no self-authorized first administrator or
unaudited bootstrap exception. This proposal completes neither P02 nor D02.

## Impact, gates, rollback and approval

Documentation only: no schema/migration, runtime/prototype, dependency, listener,
credential, production configuration, device allocation/flash write or root use
changes. Future authoritative commits, protected audit retention and projections
add server CPU/RAM/storage/latency and outage costs. Approve bounded event sizes,
queues, workers, queries/exports, retries, retention and capacity reservation for
security invalidation; no numeric policy is selected. Avoid unnecessary ARM
allocations/dependencies and flash churn; do not install a database/web runtime
on devices. Actual ARMv7/AArch64 Buildroot/Yocto/Debian evidence remains required
against 128 MB total RAM / <=16 MB idle agent RSS, including active/outage and
100-device workloads. macOS ARM and cross-builds do not establish acceptance.

ADR-009's limited approval does not approve this new consistency architecture.
Other ADR statuses remain unchanged. P00/M01 acceptance, ADR-001 disposition,
remaining D01 including ADR-002/007, D02 proof/time/storage mechanisms, D03
contracts/budgets and ordered P01–P09 gates stay open. Full WEB-01,
[ADR-003](ADR-003-browser-terminal-scope.md) and browser architecture/UX/M02
acceptance remain separate; no document merge authorizes production access.

Executed local document/link/scope/whitespace and redacted scan evidence is in
[validation.txt](../evidence/m02-d02b/validation.txt) and
[secrets.txt](../evidence/m02-d02b/secrets.txt). External URLs/fragments and
database/provider/protocol facts are unverified. No application build/lint,
database durability/failover, live audit/security/enrollment/SSH/PTY/browser,
Windows/ARM toolchain, physical ARM or hosted CI pass is claimed. Pattern scans
cannot certify absence of every secret.

Revert/supersede documents through a reviewed feature branch and Draft PR into
`develop`, preserving owner approval provenance. No deployed state needs
recovery now. Future rollback must preserve invalidation, consumed authority,
operation outcomes and audit; replaying exports or restoring a snapshot cannot
grant access or restore withdrawn trust. Do not delete audit to make an older
deployment compatible. Code rollback cannot undo commands or restore secrets.

**Human decision: PENDING; intentionally unsigned.** Tech Lead records exact
reviewed revision, scope, accept/revise/reject decision, owners and authoritative
identity/time reference. “ok next” concerns ADR-009's limited policy, not this
new architecture. Agent review and merge supply no contract/security approval
under [AGENTS.md](../../AGENTS.md).
