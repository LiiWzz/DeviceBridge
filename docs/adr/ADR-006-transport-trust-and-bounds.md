# ADR-006: Propose independently verified SSH over a bounded authenticated relay

**Status: Proposed — explicit human Tech Lead approval required.**
Issue [#57](https://github.com/LiiWzz/DeviceBridge/issues/57).
Repository source baseline: `661e191` (merge of PR #56).

Propose carrying authorized SSH byte streams through an outbound authenticated
TLS relay while keeping relay-peer trust, device identity, session authority and
SSH host-key trust distinct. Independent checks and aggregate resource accounting
cost availability, coordination and memory, but prevent a trusted relay endpoint
or an inventory record from becoming proof of the intended SSH destination.
This is a D03 design input, not a completed transport contract or implementation.

## Source context and alternatives

[ADR-001](ADR-001-foundation.md) supplies the proposed TLS/TCP 443, SSH-native,
separate control/data-plane and no-relay-SSH-private-key boundaries.
[ADR-002](ADR-002-device-record-identity-binding.md) separates device records
from verified enrollment binding. Accepted limited domain constraints in
[ADR-004](ADR-004-session-authority-binding.md) and
[ADR-005](ADR-005-credential-lifecycle-audit.md) require exact session binding,
current authority, invalidation precedence, independent recovery and durable
audit outcomes. Their unresolved mechanisms and approval gates remain open.

Source inspection at the baseline establishes only foundation behavior:
[relay `Disabled.Open`](../../server/relay/internal/transport/transport.go)
rejects without opening a socket; [relay config](../../server/relay/internal/config/config.go)
reserves ceilings without allocating streams; [CLI `run`](../../cli/cmd/devicebridge/main.go)
validates local settings without connecting; [agent `run`](../../agent/src/main.rs)
idles without networking or disk writes. These are source observations, not
executed connectivity/security tests or external library/API claims.

| Alternative | Trade-off |
| --- | --- |
| Trust a relay-selected device/address or first observed SSH key | Simple onboarding, but a wrong destination or substituted first key can become trusted. A relay connection alone does not prove SSH endpoint identity. Not proposed. |
| Terminate SSH on the relay | Simplifies some client handoff, but introduces plaintext/credential custody and changes the opaque relay boundary. Not proposed for M02; browser termination needs its separate approved architecture. |
| Proposed opaque relay plus independent intended-device host-key verification | Preserves SSH encryption and separate authorities; adds authenticated key provisioning/change, client integration and failure handling. Mechanisms remain approval blockers. |

## Proposed trust and lifecycle constraints

1. **Authenticate each boundary before forwarding.** The agent initiates outbound
   connectivity to the approved relay; no public device listener or unauthenticated
   forwarding is introduced. Verify the relay service peer under an approved TLS
   trust policy and authenticate device/client identities under approved D01/D02
   mechanisms. TLS server trust does not establish device identity, user authority
   or SSH host-key trust. Certificate validation, names, trust anchors, provisioning,
   expiry/time bootstrap and rotation/recovery require separate approved choices;
   no insecure verification switch, expired-credential fallback or fleet secret.
2. **Bind admission and attachment to the authorized stream.** Carry the exact
   operator, current intended device binding, approved SSH account/action and
   original session/client context from ADR-004 into trusted enforcement. A
   claimed stream/device ID or relay assertion is insufficient. Do not allocate
   usable forwarding or reach an arbitrary host/port from a caller-selected
   address. D01/D03 must approve the agent's permitted local SSH destination and
   prove account/channel restrictions at trusted SSH endpoints; an opaque relay
   cannot inspect encrypted login or confine arbitrary shell commands. Required
   durable authorization/audit outcomes precede usable access under ADR-005.
3. **Verify the intended SSH host independently.** Provision an authenticated
   association between the current intended device identity and its permitted
   SSH host keys through a reviewed trust path independent of the candidate
   stream's unverified key claim. A scan, displayed fingerprint, record name,
   reachable relay or successful TLS handshake alone is insufficient. Unknown,
   mismatched, revoked or unverifiable keys block SSH handoff; no automatic
   first-use acceptance or bypass prompt that converts uncertainty to trust.
   Exact proof, custody, multi-key/algorithm support and distribution remain open.
4. **Treat host-key change as an authorized trust change.** Planned rotation
   needs authenticated intended-device binding, bounded approved overlap,
   independently verified replacement and durable outcomes. Unexpected change,
   device replacement, lost keys or compromise require separately authorized
   recovery/rebinding under ADR-002/005, not deletion of a warning or reuse of
   stale inventory/key caches. Invalidation covers pending replacement trust;
   a delayed update, restart or restored backup cannot revive revoked keys or
   an old device binding. Revalidate current trust before new admission/attachment.
5. **Version and bound framing before interpreting it.** Propose explicit
   compatibility/version handling and bounded parsing of lengths, message types,
   stream references and state transitions. Reject unsupported versions,
   malformed/truncated/oversized messages, illegal sequencing and cross-session
   references without forwarding data or allocating from untrusted lengths.
   No silent downgrade. Format, negotiation/authentication coverage, multiplexing,
   error disclosure and per-stream versus connection failure scope need approved
   protocol decisions and fuzz/negative evidence; this ADR chooses no wire bytes.
6. **Handoff preserves SSH verification and least privilege.** The native client
   must verify the intended SSH host and independently authenticate to the approved
   account before an interactive session is exposed. Define how the handoff binds
   the verified identity to the exact stream and prevents concurrent substitution.
   No relay SSH private keys, custom shell RPC, root requirement, secrets in
   arguments/URLs/logs or unauthenticated local forwarding endpoint. If a local
   handoff endpoint is later chosen, approve its access control, lifetime and
   race protection explicitly; loopback alone is not authorization. Unsupported
   OpenSSH/Dropbear configurations fail closed rather than weakening checks.
7. **Separate stream closure from remote process outcome.** Specify ordered
   delivery, half-close, cancellation and drain/error semantics with bounded
   buffers and deadlines. A closed transport, timeout or cancellation is not
   confirmation of remote command rollback, process termination or descendant
   cleanup. Expiry/revocation and dependency outages must follow approved active
   session/input/closure policy; neither unlimited continuation nor an invented
   grace period is selected here. Pending outcomes remain explicit and auditable.
8. **Reconnect never replays access or input.** Backoff/retry must be bounded
   across the fleet. Lost starts reconcile the original operation under current
   authority before a conflicting new start; invalidation takes precedence.
   Reattachment, if approved, repeats current authority, target and host-key
   checks and binds the original session. Drop unsent keystrokes/paste/commands;
   never replay them automatically. TLS reconnection is no session authorization.
   Normal connectivity and independently authenticated credential recovery remain
   separate; recovery cannot grant shell or arbitrary forwarding.

## Aggregate admission and backpressure

Propose accounting for live **and in-flight** work before accepting it: service
connections, TLS handshakes, unauthenticated attempts, authenticated idle agents,
pending authorization/start/reconciliation, active streams, both directions of
buffers, queues, timers, retry state and bounded audit/log production. Enforce
reviewed server-wide and per-principal/device/client limits with fairness and
bounded rejection; per-stream limits alone are insufficient. Define reservation,
release and crash/restart reconciliation so cancellation and duplicate requests
cannot leak capacity or double-allocate it. Rate/admission enforcement must bound
work before expensive identity checks too, without granting access on saturation.

Slow consumers need bounded queues, backpressure and approved stall/deadline
policy. Do not move backlog to another unbounded queue or device flash. Account
for shared-connection head-of-line blocking and isolation, including a peer that
never completes a handshake/frame or never drains output. Select multiplexing,
buffer ownership/reuse and concurrency only after measured profiles; no per-read
allocation strategy or new device runtime dependency is selected here.

The foundation relay maxima are not production budgets: `1024` streams times
`1048576` bytes already permits 1 GiB before other costs if naively allocated.
The device baseline is **128 MB total RAM, agent idle RSS target <=16 MB**;
retain MB units and report actual units in evidence. Idle RSS alone excludes
active SSH/PTY and reconnect overhead. Hardware QA must measure supported
ARMv7/AArch64, Buildroot/Yocto/Debian cases and 100-device fleet profiles,
including TLS state, simultaneous retries, saturation, flash writes, CPU and
P50/P95/P99 latency. Cross-builds and a macOS ARM host do not establish acceptance.

## Future negative evidence and accountable decisions

These are required future executions, not passed tests. Tech Lead assigns named
people to these proposed roles and approves mechanisms/budgets before coding.

| Area / proposed owners | Decision and required evidence |
| --- | --- |
| TLS/device trust — Security/Agent/Relay | Select peer identity verification, trust provisioning/update, authenticated device/client proof, trusted-time and recovery bootstrap. Reject wrong names/peers, invalid/untrusted/expired identities, time rollback, stolen/revoked identity and unavailable trust/time services. No circular dependency on expired normal credentials. |
| SSH host trust — Security/CLI/Agent | Select authoritative key binding/storage/distribution, freshness, rotation/overlap/invalidation and recovery. Test substituted first key, wrong device with valid relay TLS, key change after authorization, multiple keys, lost rotation response, concurrent revoke, stale backup and host-key-service outage. |
| Framing and stream binding — Relay/Agent/CLI/QA | Approve version/negotiation/downgrade policy, state machine, malformed-input failure scope and stream binding. Test wrong operator/device/account/session, duplicate IDs, late frames for closed streams, truncated/oversized/unknown frames and restart/reconnect races without allocation spikes or cross-stream data. |
| Native handoff — CLI/Agent/Security | Select supported client/server versions/configurations, credential custody, exact account/action/channel enforcement and safe stream handoff. Execute OpenSSH/Dropbear interoperability, different-account login, extra channels/forwarding, local handoff theft/races if applicable, unknown/mismatched keys and unsupported-server rejection. Shell permission still carries the account's actual powers. |
| Session effects — Security/Backend/Relay/Agent | Approve expiry/revocation propagation and active-input/closure deadlines, control/audit/time outage behavior, half-close/cancel/drain and reattachment policy. Exercise lost start, disconnect, uncertain cleanup and audit-store failure; never claim completed revoke/closure without authoritative evidence. |
| Measured bounds — Hardware QA/DevOps/Relay/Agent | Approve connection/handshake/stream/buffer/queue/rate limits, timeouts, fairness, retry/backoff/jitter, flash budget and operational alerts. Measure unauthenticated floods, slow peers, output floods, concurrent 100-device reconnect, exhaustion and capacity release after cancellation/crash; retain source/binary hashes and raw per-platform logs. |

No credential scheme, crypto suite, provider/library, API/schema, wire format,
numeric lifetime/limit or default is chosen. D03 remains incomplete until these
decisions and their conflicts with D01/D02 are resolved and explicitly approved.

## Gates, impact, rollback and approval

[M02](../M02_EXECUTION.md) allows design review inputs before P00; implementation
still requires signed M01 acceptance/ADR-001 disposition, remaining agreed D01
including ADR-002, approved D02 mechanisms and D03 contracts/budgets. Limited
ADR-004/005 acceptance does not close these gates. P01–P09 remain ordered; audit
foundation precedes security mutations. Full WEB-01 and
[browser-terminal architecture](ADR-003-browser-terminal-scope.md) remain separate;
this native SSH proposal selects no browser plaintext or key custody design.

Documentation only: no listener, runtime/prototype, dependency, credential,
production config, device allocation, flash write or root use changes. Future TLS,
host-key distribution, admission accounting and audit add CPU/RAM/latency and
availability costs; validate on real ARM before rollout. No successful transport,
authentication, security, SSH/PTY, browser, application build, Windows/ARM toolchain,
physical ARM or hosted CI execution is claimed. Local document/scope checks and
secret-scan commands/output are in [validation.txt](../evidence/m02-d03/validation.txt)
and [secrets.txt](../evidence/m02-d03/secrets.txt). Pattern scans cannot certify
absence of every secret; external URLs/fragments and external API facts unverified.

Revert/supersede documents through a reviewed feature branch and Draft PR into
`develop`; no deployed session or data needs recovery now. Future rollout needs
reviewed disable-new-access and reconcile/close procedures. Deployment rollback
must preserve revoked credentials/keys, current device bindings, original
operation outcomes and durable audit history. It cannot bypass host-key checks,
revive retired trust or undo remote commands by reverting code.

**Human decision: PENDING; intentionally unsigned.** Tech Lead must record the
exact reviewed revision, scope, accept/revise/reject decision, accountable owners
and authoritative identity/time reference. This proposal, agent review and merge
are not contract/security approval under [AGENTS.md](../../AGENTS.md).
