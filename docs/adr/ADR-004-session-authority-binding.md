# ADR-004: Bind session authority to the intended user, device and SSH account

**Status: Accepted (limited domain constraints only); implementation remains blocked by unresolved design and M01/M02 gates.**
Issue: [#49](https://github.com/LiiWzz/DeviceBridge/issues/49).
Source baseline: `c3f4c99e508a6466f48bd3efd31c6eb1a1ec64dc`.

Propose an explicit per-session authorization decision for the authenticated
operator, intended device identity/current record binding, SSH account and access
context. This adds authoritative admission/reconciliation work, but prevents
inventory visibility, enrollment or a reusable device-wide permission from being
mistaken for authority to attach to any session/account. Retrofitting that boundary
after issuing access credentials would require credential and audit migration.
Only the domain constraints below are accepted; no implementation mechanism is approved.

## Source context and alternatives

[ADR-001](ADR-001-foundation.md) requires per-session checks, SSH-native access,
verified host keys and no relay private keys. [ADR-002](ADR-002-device-record-identity-binding.md)
separates device records/identity/enrollment from SSH authority, but leaves the
operator and session boundary unresolved. [M02 D01](../M02_EXECUTION.md) permits
review inputs before M01 acceptance; D02/D03 and implementation remain gated.
[ADR-003](ADR-003-browser-terminal-scope.md) adds the proposed browser context,
not a choice of SSH termination or credential custody.

The [security lifecycle inputs](../../security/README.md#proposed-credential-and-token-lifecycle--unimplemented)
already require fresh user/device/resource authorization. Actual source has no
such enforcement: [control `service.handle`](../../server/control/internal/control/server.go)
serves probes only and [relay `Disabled.Open`](../../server/relay/internal/transport/transport.go)
always rejects. This is source inspection, not live authentication evidence.

| Alternative | Trade-off |
| --- | --- |
| Infer access from visibility or enrollment | Fewer decisions, but a viewer/enroller becomes an unintended SSH operator; conflicts with the existing separation. |
| Reusable broad device access with caller-selected account/session | Easier reconnect, but a copied authority can be retargeted or survive changes without an explicit session decision; account/resource scope and attribution become ambiguous. |
| Proposed explicitly bound per-session decision | Adds admission, freshness, replay/concurrency and recovery obligations; makes the authorized target and attachment context explicit. Credential representation and enforcement mechanisms still need approval. |

## Accepted domain constraints

All constraints below are accepted at domain level only, not API fields, roles or
selected proof/credential algorithms. The [glossary](../../GLOSSARY.md) meanings stand.

1. **Authenticate the operator independently.** An approved trusted identity
   source must establish the operator; a display name, submitted user ID, device
   identity, relay assertion or browser state is insufficient. Authenticate the
   device independently too. Trust establishment/recovery, identity domains and
   mapping external identities to authoritative principals remain D01 blockers.
2. **Authorize the exact target/action/context.** The authoritative decision must
   cover that operator, intended record and its current verified device identity,
   approved remote SSH account, interactive-access action and original session
   context. Bind the admitted client/attachment context without allowing copied
   handles to substitute for operator proof. Native client and proposed browser
   access need reviewed context-specific checks; one is no automatic fallback or
   escalation path for denial of the other. No arbitrary caller-selected host,
   port, account, user, device or session may expand the decision's scope.
3. **Keep permissions separate.** Viewing, record creation, enrollment,
   administration and audit access do not imply session access. Session authority
   does not create an SSH login credential, verify a host key or grant OS rights.
   Account mapping and credential provisioning must support the approved scope;
   control metadata naming an account cannot enforce it inside encrypted SSH.
   The opaque relay must not acquire SSH plaintext/private keys to inspect the
   login. D01/D03 must locate and prove account/action enforcement at reviewed
   trusted SSH client/server boundaries before any access is enabled.
4. **Enforce again at admission and attachment.** Verification must establish
   current user/device authority, intended binding, account/context and validity
   before accepting a new session or attachment; the exact atomicity/freshness
   mechanism belongs to D02/D03. A changed/replaced identity is a changed target:
   reject the stale authorization, even if the record name/ID is unchanged. A
   grant or UI control alone is no proof of completed checks. No usable path is
   enabled if required proof, host-key trust or required audit outcome is absent.
5. **Reconcile the original operation; never retarget or replay.** A lost response
   does not prove a start failed. Correlated duplicate starts must converge on
   one authoritative outcome, without a second session or consumption of the
   same authority for a different one. Reconciliation requires current authority
   for the original operator/resource/session and must not disclose another
   user's result. Reattachment, if later supported, requires fresh checks and
   original-session binding; otherwise close/reconcile before a fresh request.
   Never replay keystrokes, paste or buffered commands on reconnect. Concurrent
   distinct sessions need distinct decisions and approved aggregate admission
   limits, not an undocumented one-session-per-device rule.
6. **Fail closed on missing or invalid authority.** Denied, expired, revoked,
   unverifiable or wrong-target authority cannot admit a new session/attachment.
   Control/audit outages cannot become an allow decision; no unaudited bypass or
   automatic reuse of cached authority is authorized here. Active-session
   detection, input rejection/closure deadlines, outage/cache/time semantics and
   cleanup need explicit D02 approval and bounded D03 enforcement. This proposal
   promises no immediate revocation, process termination or command rollback.

Interactive shell permission carries the remote account's real powers, not a
per-command allowlist enforced by this proposal. Root escalation, file transfer,
additional SSH channels and forwarding are not implicitly authorized by an
interactive-access decision. D01/D03 must define achievable restrictions and
account configuration, including what a permitted shell can itself do, and
demonstrate them with OpenSSH/Dropbear. Do not claim SSH shell access confines
arbitrary user commands or descendant processes. No custom shell RPC is proposed.

## Cases and required future evidence

These are future executed acceptance obligations, not tests passed here.

| Case | Required result / owner evidence |
| --- | --- |
| Viewer/enroller/admin requests SSH; forged user identity | Backend/Security: direct requests cannot derive session authority from unrelated capabilities or claimant IDs. |
| Valid authority copied across user/device/account/session/client context | Backend/Agent/CLI: reject each substitution, including direct attach/reconcile requests; no cross-resource result disclosure. |
| Record rebound or credential revoked between decision and attach | Backend/Agent: stale target/authority rejected at the approved transition boundary; D02 must specify races, time and active-session effects. |
| Duplicate request, concurrent start or crash around admission/audit | Backend/QA: original operation reconciles without duplicate access; distinct sessions separately authorized/bounded; no completed access with missing required audit evidence. P02 precedes grant mutations. |
| Lost start response, lost transport, expired/revoked session then reconnect | CLI/Frontend/QA: preserve original reference, no blind new start/input replay, fresh checks and explicit uncertain/closed outcomes; no reuse of invalid authority. |
| Control, identity, policy or audit dependency unavailable | Security/QA: no new start/attach based on uncertainty; prove approved D02 handling of already active sessions without inventing a local grace period. |
| SSH login chooses a different account or extra channel | CLI/Agent/Security: prove approved SSH-side account/action restrictions, independent host-key checks and unsupported-server rejection. An opaque relay checking only a device identifier is insufficient evidence. |
| Browser versus native client; unknown/mismatched SSH host key | Frontend/CLI/Security: independent authenticated context and intended-device host-key verification; no downgrade, bypass or automatic fallback after denial. |

## Remaining gates and accountable decisions

Tech Lead assigns named people to the proposed owner roles below. No provider,
tenant/role hierarchy, credential format, crypto, API/schema, wire protocol,
correlation key or numerical budget is selected by this domain proposal.

| Gate / proposed owners | Decision needed before dependent implementation |
| --- | --- |
| Remaining D01 — Architect/Security/Backend/Agent | Operator/device trust and recovery, domain/scoped capabilities/delegation, intended-device bootstrap, account/action mapping and enforcement. ADR-002/004 agreement does not complete these. |
| D02 — Security/Backend/Agent | Grant proof/custody/lifetimes/time, consumption/duplicate/outcome persistence, races/revocation/active-session effects, outages, producer/store audit coupling/retention/access/redaction and cleanup obligations. Audit content must not default to terminal keystrokes/screens. |
| D03 — Relay/Agent/CLI/Security | Authenticated bounded stream/client binding, admission/attachment/reconnect proof, host-key bootstrap/change, SSH interoperability/account/channel restrictions and measurable aggregate session/handshake/buffer/retry costs. |
| Browser/full WEB-01 — Architect/Frontend/Backend/Security | After agreed M02 inputs: select SSH termination/custody, authenticated web sessions/origin/attachment/input controls and approved UI/API/error contracts; browser renderer cannot enforce backend authority. |
| P00, P01–P09 — Human Tech Lead/QA/Hardware QA | M01 signed acceptance and ADR disposition, approved D01–D03, implementation/evidence in dependency order, real ARM matrix and independent review before production. No merge closes a gate. |

## Impact, evidence, rollback and decision record

This change is documents only: no credential, listener, config, permission grant,
runtime/prototype code, root use, device allocation or flash write. Future checks,
reconciliation/audit storage and authenticated admission add latency, memory and
availability costs; approve bounds and measure the initial 100-device fleet and
128 MB ARM / <=16 MB idle RSS target rather than assuming free checks or polling.

Reproducible local document/scope/whitespace checks and secret-scan outputs are
in [validation.txt](../evidence/m02-d01b/validation.txt),
[validate.py](../evidence/m02-d01b/validate.py) and
[secrets.txt](../evidence/m02-d01b/secrets.txt). External sources/URLs/fragments,
application builds, security/crypto enforcement, real browser/SSH/PTY, Windows/ARM
toolchains, physical ARM and hosted CI were not validated for this proposal.
Pattern scans cannot certify absence of all secrets. No library/API facts chosen.

Revert/supersede via a reviewed Draft PR into `develop`; no deployed session,
credential or data needs recovery now. Future rollback must preserve invalidated
authority, original outcomes and audit history; reverting code cannot undo remote
commands or revive revoked credentials. ADR-001/002/003 and M01/M02 stay unchanged.

## Human decision record

**Accepted with limited scope by the project owner in the conversation on
2026-11-07 (Asia/Bangkok).** The approval covers only these proposed domain
constraints: bind a session authorization decision to the authenticated operator,
intended/current device identity, approved SSH account and session/client
context; keep management/enrollment/session/SSH/OS authority separate; recheck
the binding at admission/attachment; reconcile the original operation without
retargeting or replay; fail closed for new access when required authority is
invalid or uncertain. The exact approved text is this ADR at the commit linked
from the review PR for issue #52. The user explicitly confirmed the limited
scope and that unresolved credential/provider/mechanism details still require
separate decisions. No verbatim signing identity, UTC timestamp or external
signature was captured, so none is asserted here.

This approval does **not** approve credentials, protocols, providers, account
provisioning/restrictions, API/schema, cryptography, numeric budgets, active
session termination, audit implementation, or production access. Those remain
separate D01/D02/D03 and M01/M02 gates under [AGENTS.md rule 9](../../AGENTS.md).
