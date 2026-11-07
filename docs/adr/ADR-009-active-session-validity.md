# ADR-009: Propose bounded closure when active-session validity is lost

**Status: Proposed — explicit human Tech Lead approval required.**
Issue [#63](https://github.com/LiiWzz/DeviceBridge/issues/63).
Repository source baseline: `697045e113b5a8b5ff834181e023c37eef6bc2a7`.

Propose stopping usable session traffic when current authority becomes invalid
or cannot be established, then closing the affected session within approved
bounds. This sacrifices session continuity during dependency outages, but
prevents an open connection from preserving withdrawn permission. Exact time,
propagation, enforcement and cleanup mechanisms and budgets need approval;
this D02/D03 policy proposal is not implemented or accepted.

## Context and alternatives

[ADR-004](ADR-004-session-authority-binding.md) accepts limited exact
operator/device/account/session binding. [ADR-005](ADR-005-credential-lifecycle-audit.md)
accepts lifecycle, independent agent renewal, invalidation precedence and audit
principles, leaving active-session behavior open. Proposed
[ADR-006](ADR-006-transport-trust-and-bounds.md) leaves active-input/closure
deadlines unresolved. [ADR-007](ADR-007-operator-trust-and-authority.md) proposes
current scoped policy checks; [ADR-008](ADR-008-device-bootstrap-trust.md) now
records the owner's limited bootstrap domain approval, not session policy.

Source inspection establishes no live session enforcement: [relay `Disabled.Open`](../../server/relay/internal/transport/transport.go)
rejects all streams, [control `service.handle`](../../server/control/internal/control/server.go)
serves probes and [agent `run`](../../agent/src/main.rs) idles without networking.
These repository observations are not executed security/connectivity tests.
The [glossary](../../GLOSSARY.md) distinguishes session authorization and current
session validity from an open connection or agent connectivity.

| Alternative | Trade-off |
| --- | --- |
| Continue until disconnect after the original authorization | Fewer interruptions, but stale permission may sustain access indefinitely through expiry, revocation or outage. Not proposed. |
| Pause indefinitely and automatically resume when dependencies return | Preserves session state, but retains resources and risks replay or late decisions reviving withdrawn access. Not proposed. |
| Proposed stop traffic and close within approved bounds | Limits continued access and retained resources; transient uncertainty can end useful sessions. Requires explicit detection/cutoff/closure evidence and current-authority reconciliation. |

## Proposed validity and closure policy

1. **Current validity applies beyond initial admission.** New admission and
   attachment require current exact authority under ADR-004. Established sessions
   remain subject to operator/device/grant validity, current resource binding,
   approved account/action/channel scope and required control/policy/time/audit
   evidence. Expiry, revocation/decommission, operator disablement, relevant
   policy/delegation withdrawal, device rebind or invalid host-key trust must not
   be defeated by keeping an existing connection open. Approve which trusted
   boundaries observe and enforce each condition; UI state is no enforcement.
2. **Stop application traffic on observed invalidity or uncertainty.** At a
   trusted enforcement boundary, once approved evidence no longer establishes
   validity, reject new input and stop forwarding application bytes in both
   directions for the affected session. Do not continue output as an implicitly
   authorized read-only session. Discard unsent input/paste and queued application
   data rather than replay it, route it elsewhere or persist it to device flash.
   Bound transport teardown separately; only required teardown/control evidence
   may continue under its own authority, without creating a usable data path.
   A UI input lock alone is insufficient. Already delivered or in-flight bytes
   cannot be recalled; D03 must define the exact observable cutoff boundary,
   buffers under its control and the maximum residual exposure.
3. **Close with bounded work and distinguish outcomes.** After traffic stops,
   close/release the affected session resources within separately approved
   deadlines; do not pause indefinitely waiting for dependencies or remote ACKs.
   Local forwarding stopped, local resources released, remote SSH channel closed
   and required audit durably reconciled are distinct evidence claims. A timeout
   may leave remote outcome unknown. Do not report complete remote cleanup from
   a local socket close. Closing SSH cannot guarantee terminating background
   processes/descendants or undo commands already executed; stronger process
   control would require a separately approved design.
4. **Outages do not extend access.** If a required authority/time/audit dependency
   can no longer establish validity under the approved freshness rules, apply
   stop/closure rather than stale-cache permission, an invented grace period or
   acceptance of expired credentials. A transient packet loss is not itself a
   chosen validity test: detection, evidence freshness and maximum partition
   exposure must be approved and measured. A partitioned relay or agent cannot
   instantly observe a remote revoke; independent local enforcement must bound
   continued traffic if remote evidence stops arriving. No immediate fleet-wide
   revoke guarantee or polling/lease/push mechanism is selected here.
5. **Trusted time loss fails closed without circular bootstrap.** Approve the
   authority for time, uncertainty/skew, evidence age and restart behavior.
   Clock rollback, reboot, missing trusted time or a late response cannot reset
   expiry or lengthen validity. A local monotonic timer can bound elapsed work
   but alone proves neither credential age nor persistent validity after reboot.
   Do not skip TLS/credential validation to reach a time service. Bootstrap and
   independent recovery must use their separately approved trust paths; no
   time source, lifetime, timer algorithm or persistent state format is chosen.
6. **Late allows cannot resume a stopped session.** Once invalidity/uncertainty
   triggers stop/closure, a delayed grant, renewal, dependency recovery, duplicate
   request or transport reconnect cannot resume that session or flush queued
   input. Reconcile the original session under current scoped authority before
   conflicting starts. A later authorized session requires a fresh decision and
   binding plus independent SSH authentication/host-key verification. Preserve
   revoke/decommission precedence over pending grants, retries and restored
   state. Never claim that grant expiry alone terminated a remote process.
7. **Invalidation and cleanup cannot wait for successful audit delivery.** Require
   approved durable authorization/audit before initial usable access. If required
   audit later fails, stop/close as above; an audit-store outage must not hold
   forwarding open to obtain a closure event. Capture/reconcile closure provenance
   through an approved bounded failure path and report audit outcome unknown
   until durable confirmation. Decide producer/store consistency, retry/retention,
   overflow and crash recovery before implementation; no unbounded local queue,
   secret/plaintext logging or claim of completed revoke/audit without evidence.
8. **Keep agent continuity and session permission separate.** User-session expiry
   or closure must not revoke an independently valid device identity or disable
   authenticated agent renewal/recovery. Conversely, agent connectivity or
   successful renewal cannot keep an invalid user session open. If a transport
   is shared, isolate affected session traffic and bound cleanup without
   treating a surviving control path as a shell bypass. Shared-connection failure
   may still disrupt other valid work; select isolation/multiplexing separately.

## Required decisions and future evidence

Tech Lead assigns named owners and approves mechanisms and numerical budgets
before implementation. These are future executions, not tests passed here.

| Area / proposed owners | Approval and required negative/failure evidence |
| --- | --- |
| Authority and races — Security/Backend/Agent/Relay | Define authoritative changes, admission/attachment ordering and each enforcement boundary. Exercise expiry/revoke/operator disablement/rebind/delegation withdrawal just before/after attach and during traffic, wrong session/account/device and concurrent pending grant. No broader session survives through a stale allow. |
| Partition and time — Security/Backend/Agent/QA | Approve freshness, trusted time, detection, propagation, traffic cutoff and maximum partition-exposure budgets. Test control/policy/audit/time partitions independently, delayed/lost evidence, clock rollback/reboot, reconnect storms and revoked device isolated from control. Retain origin-to-each-enforcer timing; no zero-delay claim. |
| Both directions and cleanup — Relay/Agent/CLI/QA | Specify byte cutoff, queued/in-flight data, teardown and local versus remote acknowledgement. Test blocked readers, output floods, half-close, stalled SSH endpoint, cancellation and process/relay restart. Prove stopped application traffic and bounded capacity release; distinguish remote process outcome from channel closure. |
| No resurrection — Security/Backend/QA | Exercise late allow/renewal after stop, stale cache/backup, duplicate stop/start, lost closure response, authority lost before reconciliation and policy recovery. No automatic resume, replayed input or conflicting new session; a new authorized session has fresh binding/SSH checks. |
| Audit failure — Security/Backend/DevOps | Approve bounded producer/store failure recovery, attribution, scoped outcome reads and redaction. Test disk full, duplicate/lost events, producer crash and unavailable store during revoke/closure; no usable access retained for logging and no false completion. |
| Agent independence and bounds — Agent/Hardware QA/DevOps | Test user grant expiry while valid agent renewal continues; device invalidation still blocks normal identity use. Approve aggregate sessions, timers, checks, queues, teardown/reconciliation retries and flash budgets. Measure 100-device correlated expiry/reconnect/outage, raw CPU/RAM/latency and flash effects on supported real ARM. |

This proposal selects policy direction only, not an API/error enum, wire state
machine, credential/provider/crypto, polling/push/lease mechanism, schema or
numerical cutoff/closure/retention value. Without approved measurable budgets,
no maximum access-after-revoke or cleanup guarantee can be claimed.

## Gates, impact, rollback and decision record

P00/M01 acceptance, ADR-001 disposition, remaining D01 including ADR-002/007,
D02 proof/time/audit mechanisms and D03 transport/enforcement/bounds approvals
remain prerequisites. Limited ADR-004/005/008 acceptance authorizes no runtime
implementation. P01–P09 remain ordered; audit foundation precedes sensitive
mutations. This completes neither D02 nor D03. Native SSH is the M02 baseline;
[ADR-003](ADR-003-browser-terminal-scope.md), full WEB-01 and browser architecture/
UX/M02 acceptance remain separate. A later browser service must enforce approved
backend policy; this document changes no prototype or browser contract.

Documentation only: no runtime, dependency, listener, credential, production
configuration, device allocation/flash write or root use changes. Future active
validity checks and cleanup add CPU/RAM/network/latency and outage costs. Avoid
unnecessary ARM allocations/dependencies, flash churn and device web/database
runtimes. Measure supported ARMv7/AArch64 Buildroot/Yocto/Debian workloads against
128 MB total RAM / <=16 MB idle agent RSS, including active and outage overhead;
macOS ARM and cross-builds do not prove physical-device acceptance.

Exact local document/link/scope/whitespace and redacted secret-scan output is in
[validation.txt](../evidence/m02-d02a/validation.txt) and
[secrets.txt](../evidence/m02-d02a/secrets.txt). External URLs/fragments and
provider/protocol facts are unverified. No application build/lint, live session/
revocation/security/SSH/PTY/browser, Windows/ARM toolchain, physical ARM or hosted
CI pass is claimed. Pattern scans cannot certify absence of every secret.

Revert/supersede through a reviewed feature branch and Draft PR into `develop`,
preserving the owner's approval provenance. No deployed session needs recovery
now. Future rollback must retain stopped/invalidated authority, original outcomes
and audit; an older deployment or restored database must not resume a stopped
session or bypass host keys. Reverting code cannot undo completed commands.

**Human decision: PENDING; intentionally unsigned.** Tech Lead records exact
reviewed revision, scope, accept/revise/reject decision, accountable owners and
authoritative identity/time reference. The preceding “confirm and next” approves
ADR-008's limited domain principles only; it does not approve this new proposal.
Agent review and merge supply no approval under [AGENTS.md](../../AGENTS.md).
