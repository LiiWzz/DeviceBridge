# ADR-011: Propose conservative time validity with independently trusted bounds

**Status: Proposed — explicit human Tech Lead approval required.**
Issue [#67](https://github.com/LiiWzz/DeviceBridge/issues/67).
Repository source baseline: `f3ec71ee387aab3e6bb3e9173446d83e96e5dbfa`.

Propose accepting a time-limited permission only when independently trusted
bounds establish that the whole current-time interval lies inside its authorized
window. Use separately bounded local enforcement without treating elapsed time
as initial trust. This can reject access near expiry or during time outages,
but avoids silently extending permission through uncertainty or clock rollback.
This is a D02 design input, not an approved clock or credential implementation.

## Context and alternatives

[ADR-005](ADR-005-credential-lifecycle-audit.md) accepts expiry/renewal and
independent recovery principles but leaves trusted time/skew open.
[ADR-008](ADR-008-device-bootstrap-trust.md) requires non-circular service/time
bootstrap; [ADR-009](ADR-009-active-session-validity.md) accepts stopping and
closing sessions when validity cannot be established. Approved architecture
direction in [ADR-010](ADR-010-authoritative-audit-consistency.md) distinguishes
authoritative outcomes from observed/uncertain remote effects; it selects no
timestamp representation or time provider.

Source inspection shows [agent `run`](../../agent/src/main.rs) using a shutdown
poll interval while idling without networking, and [control `service.handle`](../../server/control/internal/control/server.go)
serving probes only. The [security baseline](../../security/README.md) reserves
trusted time/lifetimes as future decisions. These are source observations, not
tests of synchronization, clock APIs or expiry enforcement. Existing shutdown/
request timeouts are no credential lifetime or trusted-time evidence.

| Alternative | Trade-off |
| --- | --- |
| Trust the current wall-clock value or a caller's timestamp | Simple, but lacks independent authority/freshness and can accept replay or rollback. Not proposed. |
| Start a local elapsed timer on receipt/reboot | Bounds some local work but cannot alone prove a credential's prior age or time spent offline. Not proposed as initial or persistent time trust. |
| Proposed independently trusted uncertainty bounds plus local enforcement | Makes ambiguous validity explicit; adds bootstrap, evidence-age, drift and availability costs. Concrete source, verification and clock mechanisms remain approval blockers. |

## Proposed policy

1. **Time evidence needs authority, scope and freshness.** Approve the time
   trust source, protected provisioning, verification and permitted deployment/
   purpose before accepting evidence. An authenticated peer is not automatically
   an authorized time authority. Bind evidence to its intended context and
   establish age independently of a caller-controlled timestamp or delivery
   arrival. Delayed/replayed samples cannot reset freshness. Define bounded
   source error, delivery uncertainty, local elapsed-time error and evidence-age
   limits before deriving current bounds. A plausible date or successful
   handshake alone is insufficient; no source/protocol/crypto is chosen here.
2. **Require the entire trusted interval to fit.** Where approved policy defines
   a not-before boundary `N` and expiry `E`, and trusted current time is in closed
   interval `[L, U]`, time validity requires `L >= N` and `U < E`. The instant
   `E` is excluded; equality at `N` is permitted. A bound crossing either boundary
   cannot yield a temporal allow. This is only one necessary condition: identity,
   intended binding, current policy/revocation, proof freshness and required audit
   must independently pass. No trusted bounds, reversed/malformed bounds or an
   invalid window yields no temporal allow. Representation, rounding, overflow,
   optional-boundary policy and comparison implementation need separate approval.
3. **Local deadlines enforce but do not create validity.** Derive conservative
   local stop/revalidation obligations from approved trusted bounds and error
   assumptions; account for evidence aging and elapsed time before use, not only
   when initially sampled. Define suspend/resume and clock-change detection,
   maximum check delay and scheduling/error budgets. Do not reset a credential's
   original window on retry, transport reconnect or clock adjustment. A local
   monotonic interval alone is no trusted absolute age across reboot or restore.
   This selects no OS clock API, timer algorithm, polling/lease protocol or
   numerical drift/freshness/check/closure policy.
4. **Clock/source uncertainty cannot extend access.** Backward or forward jumps,
   source disagreement, invalid verification, excessive uncertainty, source
   compromise or loss of required evidence trigger conservative revalidation;
   if validity cannot be established, withhold new access and apply ADR-009's
   session stop/closure. Do not pick whichever source grants the longest window,
   silently average conflicting evidence or invent a grace period. Approve source
   quorum/disagreement and trust-change policy separately. A time correction may
   support a fresh independently authorized decision; it cannot resume a stopped
   session or restore revoked/consumed authority. Never claim instant fleet-wide
   clock/revoke detection without approved measured propagation bounds.
5. **Reboot/restore requires fresh justified trust.** Persisted timestamps, a
   restored database/image or a restarted timer cannot prove time continuity or
   absence of later invalidation. Until an approved continuity or independent
   rebootstrap path establishes trustworthy current bounds and current authority,
   deny time-sensitive access. Do not write a new flash checkpoint for every
   sample/request. Select protected state, continuity assumptions and rollback/
   tamper detection only after hardware and failure evidence; no storage format
   or hardware clock capability is assumed across supported devices.
6. **Bootstrap and recovery must break the trust cycle explicitly.** Provision
   independent service/time verification prerequisites before rollout. An
   expired normal credential cannot authenticate its own replacement/time proof;
   an unverified time response cannot authorize the peer supplying it. Design a
   protected initial trust path under ADR-008 and independent narrowly scoped
   recovery under ADR-005 without disabling TLS/certificate/time checks or
   accepting the first observed source. If no valid path survives, withhold access
   and use physical service or separately reviewed independent management.
   This obligation selects no bypass, permanent fallback or concrete protocol.
7. **Renewal remains independent and audit preserves uncertainty.** Start valid
   agent renewal early enough for approved uncertainty/outage/retry margins;
   neither user SSH grant expiry nor a user's session closure expires a valid
   device identity. Missing trusted time still cannot authorize renewal/access
   with expired or unverifiable normal evidence; follow approved recovery instead.
   Under ADR-010, record trusted producer, available time bounds/provenance and
   uncertainty without asserting an exact verified instant when none exists.
   Preserve original audit observations and append attributable corrections.
   A timestamp alone is neither proof of causal ordering nor of remote completion;
   approve operation correlation and ordering separately. No secret/raw credential
   or terminal content belongs in time/audit diagnostics.

## Decisions and future evidence

Tech Lead assigns named owners and approves mechanisms/budgets first. These are
future executed obligations, not passing clock/security tests.

| Area / proposed owners | Required decision and negative evidence |
| --- | --- |
| Source/provenance — Security/Backend/Agent | Select trust source, protected provisioning/change, scope, freshness, error bounds and disagreement policy. Test forged/replayed/delayed samples, wrong domain/purpose, compromised/withdrawn source and sources that disagree; authenticated identity alone cannot confer time authority. |
| Window boundaries — Security/Backend/QA | Specify representation/rounding/overflow and any absent boundaries. Exercise `L = N`, `L < N`, `U = E`, `U > E`, wholly valid versus crossing intervals, reversed bounds, empty/expired windows and no evidence; only wholly justified windows can pass the temporal condition. |
| Local enforcement — Agent/Relay/CLI/QA | Approve elapsed-time/error accounting, drift, suspension and maximum check/cutoff/closure delay. Test queued/delayed decisions, clock jumps, scheduler stalls, suspend/resume and uncertainty growth while traffic runs; apply stop/closure, never automatic resume/replay. |
| Restart/bootstrap/recovery — Security/Agent/DevOps | Choose independent service/time trust bootstrap and any protected continuity state. Test no RTC/network/trusted time, reboot during renewal, stale image/backup, normal credentials after expiry, absent/revoked recovery proof and loss of all trust. No circular bootstrap or skipped certificate checks; record physical-service limits. |
| Authority/audit — Backend/Security/QA | Exercise revoked grant with otherwise valid time, valid agent renewal after user expiry, invalidated source during mutation, audit with uncertain time and later correction. No time-derived permission, history rewriting, invented total event order or false remote completion. |
| Bounds — Hardware QA/Agent/DevOps | Approve sample/check/state sizes, timers, retries/backoff, diagnostic queues and flash budgets. Measure 100-device synchronized time outage/recovery, malformed evidence, stalled sources and repeated reboot; retain raw CPU/RAM/network/latency/write evidence for supported real ARM. |

## Gates, impact, rollback and decision record

Documentation only: no time service, schema/migration, runtime/prototype,
dependency, listener, credential, production config, device allocation/flash
write or root use changes. Future verification, sampling, enforcement and audit
add CPU/RAM/network/latency and outage costs. Avoid unnecessary ARM allocations/
dependencies and flash churn; no device database/web engine. Measure supported
ARMv7/AArch64 Buildroot/Yocto/Debian workloads against 128 MB total RAM /
<=16 MB idle agent RSS, including active/outage and 100-device profiles.
Cross-builds/macOS ARM do not establish physical-device acceptance.

Other ADR statuses remain unchanged except ADR-010's recorded architecture
approval. This proposal completes no D02/D03 or implementation gate. P00/M01
acceptance, ADR-001 disposition, remaining D01 including ADR-002/007,
D02 proof/storage/time mechanisms and numerical budgets, D03 contracts and
ordered P01–P09 remain open. Full WEB-01,
[ADR-003](ADR-003-browser-terminal-scope.md) and browser architecture/UX/M02
acceptance remain separate; no merge supplies production authorization.

Exact local document/link/scope/whitespace and redacted scan evidence is in
[validation.txt](../evidence/m02-d02c/validation.txt) and
[secrets.txt](../evidence/m02-d02c/secrets.txt). External URLs/fragments and
time-provider/protocol/OS API facts are unverified. No application build/lint,
live time/expiry/security/SSH/PTY/browser, Windows/ARM toolchain, physical ARM
or hosted CI pass is claimed. Pattern scans cannot certify absence of all secrets.

Revert/supersede documents through a reviewed feature branch and Draft PR into
`develop`, preserving owner approval provenance. No deployed time/state needs
recovery now. Future rollback must preserve expiry, stopped/invalidated and
consumed authority, original outcomes and audit; restoring a clock/image/database
cannot extend access or resume a stopped session. Code rollback cannot undo
remote commands or restore lost trust.

**Human decision: PENDING; intentionally unsigned.** Tech Lead records exact
reviewed revision, scope, accept/revise/reject decision, owners and authoritative
identity/time reference. “approve and next” approves ADR-010's architecture
direction, not this new time policy. Agent review and merge supply no approval
under [AGENTS.md](../../AGENTS.md).
