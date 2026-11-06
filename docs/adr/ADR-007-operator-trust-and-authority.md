# ADR-007: Propose trusted operator identities and scoped maintenance authority

**Status: Proposed — explicit human Tech Lead approval required.**
Issue [#59](https://github.com/LiiWzz/DeviceBridge/issues/59).
Repository source baseline: `168578d8e9e0ff5ecada58b282d69fdf32852fc7`.

Propose mapping verified operator identities to explicit action/resource authority
at trusted service boundaries, with delegation permitted only within independently
authorized limits. This adds policy, identity-mapping and reconciliation costs,
but avoids treating a login, device enrollment or an administrator label as
unrestricted access. This remaining D01 input proposes no provider, role hierarchy,
tenant model, credential scheme, crypto, API/schema or numerical policy.

## Repository context and alternatives

[ADR-002](ADR-002-device-record-identity-binding.md) requires independent operator
authority and intended-device proof for enrollment, leaving user trust and
delegation open. [ADR-004](ADR-004-session-authority-binding.md) accepts limited
exact user/device/account/session binding; it does not select operator trust or
roles. [ADR-005](ADR-005-credential-lifecycle-audit.md) accepts lifecycle,
invalidation precedence and durable outcomes at domain level only.
[ADR-006](ADR-006-transport-trust-and-bounds.md) remains a proposed transport
input, not operator proof. [M02](../M02_EXECUTION.md) permits design review inputs
before P00 but gates implementation on explicit prerequisite approvals.

Existing [control `service.handle`](../../server/control/internal/control/server.go)
serves health/readiness probes only, and [relay `Disabled.Open`](../../server/relay/internal/transport/transport.go)
rejects streams. The [security baseline](../../security/README.md) describes no
operator authentication or policy service. These are repository source
observations, not live enforcement tests or external identity-provider claims.
The [glossary](../../GLOSSARY.md) separates operator identity, maintenance authority,
device identity and session authorization.

| Alternative | Trade-off |
| --- | --- |
| Trust submitted IDs, UI controls or relay assertions | Simple, but direct requests and identity substitution evade the intended boundary. Not proposed. |
| Treat every authenticated administrator as an unrestricted maintenance operator | Fewer decisions, but silently combines inventory, enrollment, audit, policy and SSH powers. Not proposed; an approved administrator may later hold explicit scoped capabilities. |
| Proposed authoritative identity mapping and scoped decisions | Adds mapping, policy availability and audit work; makes identity collisions, delegation and per-action denial explicit. Exact mechanisms remain unselected. |

## Proposed constraints

1. **Establish operator trust before mapping authority.** Accept identity evidence
   only from an explicitly approved trust source and verification path. Bind its
   stable subject and authoritative identity domain to the intended local human
   principal; display name, email, request IDs, browser state and device proof are
   insufficient. Identical names/subjects in different domains are not the same
   operator. Do not merge/link identities or transfer rights on renaming,
   reassignment, source migration or account recreation without independently
   authorized verification and attributable history. Trusted identity-source
   onboarding/change, proof validation and disabled/deleted-account freshness
   remain concrete D01/D02 decisions; no implicit local fallback on outage.
2. **Authorize the exact action and current resource scope, denying by default.**
   Trusted services must establish the authenticated actor, requested action,
   authoritative domain/resource, current device binding where relevant, and
   applicable policy before returning protected data or committing a mutation.
   A caller-selected resource/domain is a claim to verify. List/search/export,
   individual reads, mutations and outcome reconciliation all need scoped checks;
   a client filter or hidden button is insufficient. Unknown scope, conflicting
   or missing policy and unverifiable authority cannot yield an allow decision.
   Policy composition/precedence and representation need separate approval.
3. **Keep capabilities distinct.** Inventory visibility/creation, enrollment,
   identity recovery, policy changes, delegation, session access and audit read/export
   are distinct powers, as is administration of credentials. None automatically
   grants the others. This is a capability separation proposal, not a named role
   set or tenant hierarchy. Session access still requires ADR-004's exact operator,
   intended/current device, SSH account/action and client/session context plus
   independent SSH login/host-key trust. Agent renewal remains independent of
   operator SSH grants; no service may substitute its own broader credentials
   for an operator's missing authority to become a confused deputy.
4. **Delegation cannot amplify authority.** If later supported, both permission
   to delegate and the delegated action/resource scope must be independently
   authorized. Merely possessing a permission does not imply a right to delegate
   it. Each effective delegated decision is bounded by the grant, every current
   upstream delegation limit and relevant current policy/identity validity;
   removal/expiry/revocation of an upstream right cannot leave an independent
   surviving copy. No implicit transitive redelegation or self-grant of missing
   authority. Self-directed policy edits cannot authorize themselves; an existing
   independently valid authority must approve them. Cycles, unknown ancestry,
   ambiguous scope and unbounded chains fail closed. Chain model, depth, lifetime,
   enforcement cost and any separately authorized permanent assignment remain
   unresolved choices, not a selected delegation feature or numeric policy.
5. **Policy changes are security operations.** Changing identity mappings,
   capability assignments, delegation, trust sources or scope requires current
   independent authority and durable attributable outcomes under ADR-005.
   Recheck validity/resource binding at the approved transition boundary for new
   admission/attachment and sensitive mutations. A stale allow decision cannot
   survive identity disablement, resource rebind, policy withdrawal or upstream
   invalidation. Exact serialization/freshness, propagation and active-session
   input/closure policy remain D02/D03 approval blockers; do not promise immediate
   remote process termination or rollback of completed actions.
6. **Uncertain outcomes reconcile the original operation under current authority.**
   A timeout or policy-service/audit failure is no proof of rejection or success.
   Preserve original actor/resource/action attribution and correlate duplicate
   mutations without duplicating authority; disclose protected outcomes only to
   a currently authorized principal. Losing authority can block the original
   caller from reconciliation; a separately authorized scoped responder must
   resolve it without impersonating that caller. Independently authorized revoke
   or decommission takes precedence over uncertain grants; late decisions,
   retries, restart or backup restore cannot resurrect withdrawn authority.
   No new allow from stale cached policy, missing trusted time or unavailable
   required identity/control/audit dependencies. Report unknown durable outcomes
   explicitly; never claim a completed revocation without confirming evidence.

## Initial trust and audit bootstrap

The first administrative identity cannot authorize its own trust or be established
by an unauthenticated first-claim endpoint. Propose an explicitly reviewed
provisioning ceremony with independently authenticated deployment-owner authority,
protected identity mapping/policy installation and durable attributable evidence
before any maintenance access is enabled. The same principle covers replacement
or loss of all administrative identities; an operator named in inventory, a device
credential, database access or a permanent emergency token is no implied recovery
authority. No provisioning command, production configuration or secret is defined.

P01 must specify the protected verification/provisioning prerequisites; P02 must
provide accepted audit producer/storage behavior before security mutations or
public enrollment/grant/policy administration are exposed. If initial identity or
policy installation itself precedes the ordinary audit service, its separately
approved provisioning evidence must be durable and authoritatively reconciled
with audit before enabling access. Missing evidence blocks activation; bootstrap
is no unaudited exception. Select who authenticates the initial owner, how that
trust is delivered/verified, evidence custody/integrity and audit reconciliation
in approved designs. This breaks no P01/P02 dependency by inventing a permissive
administrator, and does not assert the ceremony exists or is approved.

## Future negative evidence and decisions

These are future executed obligations. Tech Lead assigns named people to the
proposed roles; unresolved choices block dependent implementation.

| Area / proposed owners | Required decision and evidence |
| --- | --- |
| Operator trust — Security/Backend/Architect | Select trusted identity sources, stable-domain/subject mapping, verification, linking/reassignment/migration, disabled-user freshness and recovery. Reject forged IDs, same email/name across domains, recreated accounts inheriting old rights, untrusted source changes and unavailable/revoked identity evidence. |
| Scope and enforcement — Backend/Security/QA | Select resource/domain model, policy composition, deny/conflict rules and trusted enforcement boundaries. Test direct UI-bypass requests, wrong action/domain/device/account, filtered list versus detail/export/reconcile leakage and service credentials used as a confused deputy. No role/tenant model is approved here. |
| Delegation/policy lifecycle — Security/Backend | Decide whether delegation is supported; if so approve bounded scope/ancestry, redelegation, non-amplification, withdrawal and invalidation propagation. Test self-escalation, cycles, missing ancestors, rights removed after delegation, concurrent policy edit/rebind/revoke and uncertain grant followed by invalidation. Distinguish separately authorized assignments from delegated authority. |
| Provisioning and total lockout — DevOps/Security/Backend | Approve initial-owner proof, protected delivery/install, independently authorized recovery, evidence custody and P01/P02 audit reconciliation. Test first claimant, substituted owner, missing audit, partial provisioning/crash, duplicate/lost response, restored backup and loss of every admin identity. No default administrator or hidden bypass. |
| Time/outage/audit — Security/Backend/QA | Approve freshness, trusted time, atomic decisions and producer/store consistency, current-authority reconciliation, active-session behavior and scoped responder access. Exercise disabled operator or upstream revocation during mutation, dependency outages, late accepted results, redaction and unauthorized audit/outcome reads. |
| Bounds — Backend/Hardware QA/DevOps | Approve identity/policy/chain/query/cache/queue/retry bounds and workloads. Measure malformed requests, deep/cyclic delegation, large inventory/export, concurrent policy changes and 100-device activity; retain exact revisions, raw CPU/RAM/latency and audit storage costs. |

## Impact, gates, rollback and decision record

Documentation only; no runtime/prototype, dependency, permission, identity source,
credential, production config, device allocation/flash write or root use changes.
Future authoritative checks, mapping/policy storage and audit add server RAM/CPU,
latency and outage costs. Device runtime stays Rust without a web/database engine;
cross-platform proof and transport overhead still need real ARMv7/AArch64 and
Buildroot/Yocto/Debian evidence against 128 MB total RAM / <=16 MB idle agent RSS.
No delegation feature or resource guarantee is implemented by this proposal.

Remaining device bootstrap/proof, D01 mechanisms and ADR-002 approval, D02
mechanisms, Proposed ADR-006/D03 decisions and P00/M01 acceptance remain gates.
Limited ADR-004/005 domain acceptance does not authorize implementation. P01–P09
remain dependency ordered; full WEB-01 and browser architecture/UX/M02 acceptance
remain separate. ADR-001–006 statuses and decisions are unchanged.

Exact local document/link/scope/whitespace checks and secret-scan output are in
[validation.txt](../evidence/m02-d01c/validation.txt) and
[secrets.txt](../evidence/m02-d01c/secrets.txt). External URLs/fragments and provider
facts unverified. No application build/lint, live authentication/policy/security,
SSH/PTY/browser, Windows/ARM toolchain, physical ARM or hosted CI pass claimed.
Pattern scans cannot certify absence of every secret.

Revert/supersede documents through a reviewed feature branch and Draft PR into
`develop`; no deployed policy or identity needs recovery now. Future deployment
rollback must preserve disabled identities, withdrawn/delegated authority limits,
revocation, operation provenance and audit. Restoring a database backup or old
policy cannot restore invalidated authority. Reverting code cannot undo completed
remote commands or restore lost trust.

**Human decision: PENDING; intentionally unsigned.** Tech Lead records exact
reviewed revision, decision/scope, accountable owners and authoritative identity/
time reference. Neither instruction to continue, agent review nor merge changes
this status or supplies approval under [AGENTS.md](../../AGENTS.md).
