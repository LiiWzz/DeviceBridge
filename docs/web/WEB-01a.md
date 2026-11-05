# WEB-01a — Preliminary backend gap inventory

Issue: [#33](https://github.com/LiiWzz/DeviceBridge/issues/33).
Local source inspected: `develop` base
`8a6307a96ea0c46b2d91e089925c7dddb21d6377`, 2026-11-05.

**Status: preliminary discovery for review; no contract or security approval.**
This packet maps the seven proposed dashboard screens to existing source and
open backend questions. It does not complete WEB-01: the
[roadmap's WEB-01 row](../ROADMAP.md#uxui-and-implementation-work-packages)
requires agreed M02 identity/authorization design. No route, application schema,
token format, role, grant, security policy or production implementation is chosen.

The [UX-01 review record](../ux/UX-01.md#human-review-and-handoff) and
[UX-02 flow record](../ux/UX-02.md#human-flow-review-record) clear discovery inputs
only. [UX-03](../ux/UX-03.md#decision-and-primary-source) is proposed with human
review pending; merging its documentation does not supply that review.
[DB-010](../acceptance/DB-010.md#evidence-identity-and-limits) records pending
M01 human, governance and physical ARM gates. Neither packet authorizes M02.
[ADR-001](../adr/ADR-001-foundation.md) remains Proposed, and
[AGENTS.md rule 9](../../AGENTS.md) requires an ADR and Tech Lead approval for
contract/security changes.

## Observed foundation and its limits

These observations come from local first-party source, not fresh executions or
an external service audit. Symbols identify the implementation reviewed; existing
tests below establish intended coverage only, without claiming they ran here.

| Source / symbol | Observed behavior | Backend implication |
| --- | --- | --- |
| [Control `service.handle`, `service.serve`, `newService`](../../server/control/internal/control/server.go) | Exact `/healthz` and `/readyz` paths; GET/HEAD only; other paths 404 and unsupported probe methods 405. Readiness reflects process serving/shutdown, not database/device reachability. Headers and request timeouts are bounded. | Probes cannot supply authorized inventory, fleet counts, freshness or connected-session evidence. Their small JSON bodies are existing probe behavior, not a dashboard error model. |
| [Control `Config.Validate`, `LoadConfig`](../../server/control/internal/control/config.go) | Loopback IP literal, unprivileged port and bounded timeout settings; no user/session configuration. | This local plain-HTTP listener is not a production authenticated API boundary. Do not proxy it into a dashboard as an authentication substitute; see [security trust boundaries](../../security/README.md#assets-and-trust-boundaries). |
| [Relay `run`](../../server/relay/cmd/devicebridge-relay/main.go), [`Disabled.Open`, `Connector`](../../server/relay/internal/transport/transport.go) | Process idles without listeners; `Disabled.Open` returns nil plus `ErrUnimplemented`. `Connector` is an internal provisional seam with no framing or authentication. | No transport, authorized stream, tunnel or session telemetry exists to integrate. Retain disabled behavior until approved M02 implementation. |
| [Relay `Config.Default`, `Config.Validate`](../../server/relay/internal/config/config.go) | Only disabled transport; reserved stream/buffer ranges validated; no actual buffers or streams allocated. | Settings are not evidence of runtime admission/resource enforcement. Multiplying maxima would permit 1 GiB if naively allocated; aggregate bounds need review, also recorded in [security baseline](../../security/README.md#assets-and-trust-boundaries). |
| [CLI `run`, `usage`](../../cli/cmd/devicebridge/main.go) | Help/version and `check` only. Check validates local settings without contacting a server; unknown remote commands return unavailable. | No existing enrollment/connect command, handoff protocol or native-session result to expose in the UI. Do not invent executable syntax. |
| [Agent `execute`, Unix/non-Unix `run`](../../agent/src/main.rs), [`Config`, `Config.load`, `Config.read`](../../agent/src/config.rs) | Local help/version/config check and Unix signal-driven idle loop; config has only logging and shutdown polling. Idle loop has no network/disk writes; non-Unix run reports Unsupported. | No device registration, heartbeat, credential lifecycle or status source. Agent config validity is not enrollment or connectivity. Future telemetry requires a separately scoped agent issue and ARM evidence. |
| [Compose `services.db`](../../deploy/compose.yaml), [migration runner](../../deploy/migrate.sh) (`devicebridge_migrations.history`, advisory transaction lock), [baseline migration](../../deploy/migrations/0001_baseline.up.sql) | Development PostgreSQL only in Compose; loopback binding and required password. Runner tracks migration versions/checksums; baseline adds metadata comment only, explicitly no application tables. | Migration bookkeeping is not device/user/session/audit storage. Development DB configuration and administrator role are not production deployment/auth design; see [security trust boundaries](../../security/README.md#assets-and-trust-boundaries). |
| [Security baseline opening and proposed lifecycle](../../security/README.md) | Explicitly no enrollment, credential store, authentication service, session authorization, TLS transport or remote SSH. Lifecycle table is unimplemented review input. | Identity, web sessions, resource authorization and auditable lifecycle behavior remain missing. Logs/probes/migration timestamps are not an audit event contract. |

Existing source checks include [control `TestProbeRoutes` and `TestMalformedHTTPAndHeaderLimit`](../../server/control/internal/control/server_test.go),
[relay `TestDisabledAlwaysFailsClosed`](../../server/relay/internal/transport/transport_test.go)
and [CLI `TestRun` and `TestMalformedArgumentsDoNotLeakValues`](../../cli/cmd/devicebridge/main_test.go).
They concern foundation behavior, not dashboard authentication or M02 acceptance.

## Seven-screen traceability

The needs below derive from [UX-01's inventory and journeys](../ux/UX-01.md#mvp-screen-inventory),
[UX-02's synthetic flow examples](../ux/UX-02.md) and
[UX-03's screen/state handoff](../ux/UX-03.md#screen-and-state-handoff).
The current-source column points back to source-owned observations; the final
column contains future questions, not selected fields, APIs or semantics.

| Screen / intended need | Current source and gap | Questions for full WEB-01 / acceptance |
| --- | --- | --- |
| Sign-in / access denied | [Security baseline opening](../../security/README.md): no auth service or credential store; [control `service.handle`](../../server/control/internal/control/server.go) serves probes only. | Which approved user identity and web-session mechanisms, login/logout/expiry protections, return navigation and forbidden/missing disclosure rules? What protection against session fixation, CSRF where applicable and stale protected caches? |
| Fleet overview | [Control `service.handle` / `service.serve`](../../server/control/internal/control/server.go): process probes only; [agent `run`](../../agent/src/main.rs) produces no connectivity observations. | What authoritative authorized subset supplies consistent counts/activity? How distinguish unknown, unavailable, offline and stale, and record observation freshness without claiming SSH readiness? |
| Devices | [Baseline migration](../../deploy/migrations/0001_baseline.up.sql): no application tables; [control `service.handle`](../../server/control/internal/control/server.go): no inventory/query handler. | What reviewed identifiers/metadata, scoped search/filter/page behavior, bounds and count consistency? How prevent direct-ID and query-based disclosure? |
| Device detail / onboarding / revoke | [Agent `Config` / `run`](../../agent/src/config.rs), [agent runtime](../../agent/src/main.rs) and [security lifecycle table](../../security/README.md#proposed-credential-and-token-lifecycle--unimplemented): no identity/enrollment/revocation implementation. | How bind intended device identity, bootstrap/recover safely, handle replay/expiry and reconcile ambiguous mutation results? What are future-request versus active-session revocation effects and outage/propagation guarantees? No credential instructions until approved. |
| Sessions / native CLI handoff | [CLI `run`](../../cli/cmd/devicebridge/main.go) and [relay `Disabled.Open`](../../server/relay/internal/transport/transport.go): remote operations unavailable; [baseline migration](../../deploy/migrations/0001_baseline.up.sql): no session store. | What separately authorized request/grant/established/history evidence and own/other visibility? What accepted native-client interface verifies host keys and avoids credentials in URLs/commands? Is termination supported, with independent authorization and audit? |
| Users / permissions | [Security baseline opening](../../security/README.md) and [baseline migration](../../deploy/migrations/0001_baseline.up.sql): no user identity, permissions enforcement or application store. | What capability/resource scope model, administrative boundaries, self-escalation safeguards, concurrent-change reconciliation and audit evidence? Role names and grants remain unselected. |
| Audit history | [Control `service.serve`](../../server/control/internal/control/server.go) logs process lifecycle; [migration runner](../../deploy/migrate.sh) records migration history only; [baseline migration](../../deploy/migrations/0001_baseline.up.sql) has no audit store. | Which authoritative event producers, attempted/denied/completed outcomes, actor/device/action/time/result semantics, access/redaction, retention, integrity and failure behavior? Empty history must remain distinct from unavailable evidence. |

## Capability separation and owner gates

Carry forward [UX-01's draft matrix](../ux/UX-01.md#draft-permission-discussion-matrix):
viewing devices, requesting SSH, inspecting sessions, enrolling, revoking,
changing user permissions, reading audit and terminating a session (if supported)
are separate decisions. Viewing cannot imply connection authority; connection
cannot imply administration; administration cannot imply unrestricted SSH/audit.
No row defines a role or grants authority. UI visibility must not become backend
enforcement, including requests made outside the browser.

The owner assignments below are proposed responsibilities for planning, derived
from the [roadmap owners and start gates](../ROADMAP.md#uxui-and-implementation-work-packages).
No approval is recorded by this report.

| Gate / decision | Proposed owner | Required outcome before dependent work |
| --- | --- | --- |
| M01 acceptance and architecture disposition | Human Tech Lead / Architect / independent reviewer | Resolve [DB-010 pending gates](../acceptance/DB-010.md) and ADR status; explicit acceptance before M02 implementation. |
| M02 identity and per-session authorization agreement | M02 backend/agent/CLI owners + Security / Architect | Agreed design available to full WEB-01; approved ADRs before contract implementation. Select no token/identity scheme through UX discovery. |
| Full WEB-01 web/API/security contracts | Backend / Security / Architect; Human Tech Lead approval | Define and review identity integration, session protections, resource/action authorization, data/query/error/mutation contracts and audit. Approved ADRs and explicit Tech Lead approval before implementation. |
| UX-03 review | Human owner / UX | Approve or revise [proposed design handoff](../ux/UX-03.md#implementation-responsibilities-and-acceptance); merge is not review evidence. |
| Production integration | Backend / Frontend / QA / Human Tech Lead | WEB-01 approval, UX-03 approval and M01/M02 acceptance, then matching accepted APIs and WEB-02+ evidence. No production readiness inference from this inventory. |

## Dependency-ordered proposed follow-up slices

These are discussion slices, **not created GitHub issues, committed API designs
or permission to start implementation**. Tech Lead should scope each into an
issue, follow dependencies and execute one issue at a time under
[AGENTS.md](../../AGENTS.md). The order refines the existing
[WEB-01..06 plan](../ROADMAP.md#uxui-and-implementation-work-packages).

1. Resolve M01 gates and agree the M02 identity/authorization input. Record open
   enrollment, host-key, revocation, outage and resource decisions with accountable
   owners; approved ADRs precede implementing those contracts.
2. Complete full WEB-01 after that agreement: review web identity/session and
   capability/resource enforcement, then data/freshness/query/error and mutation
   reconciliation contracts together with the audit producer/storage/query
   contract. Submit ADRs and wait for explicit Tech Lead approval. Audit is a
   prerequisite of security-sensitive mutations, not a later UI-only addition.
3. After approvals, scope shared backend identity/enforcement and audit foundations
   into separate issues with direct authorization-negative tests. Build approved
   M02 identity/enrollment/authorized transport/native CLI behavior in its own
   dependency graph; demonstrate acceptance before production web integration.
4. Scope authorized read capabilities (inventory/observations/consistent overview,
   then device detail/query bounds) and administrative lifecycle capabilities
   separately. Enrollment/revocation and permission mutations depend on accepted
   identity, enforcement, reconciliation and audit support. Session/history/native
   handoff depend on accepted M02 evidence; termination remains conditional.
5. With matching accepted backends, UX-03 review and M02 acceptance, scope WEB-02's
   authenticated shell/client and WEB-03/04's seven screens. Rewrite production UI
   separately from [UX-02's throwaway prototype](../ux/UX-02.md#human-flow-review-record).
6. WEB-05 verifies integration/security/accessibility/resource behavior and
   deployment recovery; WEB-06 supplies independent review and human GO/NO-GO.
   No slice waives the preceding approval or hardware gates.

## Negative, security and resource questions to retain

These are proposed evidence obligations grounded in the
[roadmap validation plan](../ROADMAP.md#validation-and-operational-acceptance),
[UX-01 journeys](../ux/UX-01.md#journeys-for-human-review),
[security lifecycle questions](../../security/README.md#proposed-credential-and-token-lifecycle--unimplemented)
and [ADR-001 invariants](../adr/ADR-001-foundation.md#security-invariants):

- Can unauthenticated, wrong-user/device/session and direct-ID requests disclose
  records or mutate state? How are expiry, revocation, stale authorization and
  control-plane outages rejected without fallback or self-escalation?
- How are duplicate/replayed submissions, concurrent changes, lost responses,
  uncertain success and restart/power loss reconciled? How does audit distinguish
  attempt, denial and completion when an event producer/store fails?
- Which limits bound queries, page sizes, requests, connections, buffers, history
  growth and refresh/backoff? How are counts consistent with authorized results
  in a 100-device fleet? Choose measurable budgets in approved WEB-01/02 work;
  this packet supplies no polling interval, TTL, page size or performance claim.
- How do tests demonstrate host-key mismatch rejection, no relay SSH private-key
  custody, no custom shell protocol, no unauthenticated forwarding, and no secrets
  in UI/URLs/arguments/logs/audit? How do session protection, CSRF where applicable
  and XSS checks follow the selected approved design?
- What happens when history, identity provider, inventory or telemetry is
  unavailable? How does UI withhold protected/stale content and avoid declaring
  connectivity or success from a probe, local check or missing audit event?
- Do later agent changes retain least privilege, bounded allocation and flash
  writes on 128 MB devices with the <=16 MB idle RSS target? Measure actual ARM
  model/kernel/distribution, reconnect and restart; cross-build is insufficient.
  Keep dashboard assets/runtime on server/browser infrastructure as required by
  the [roadmap resource boundary](../ROADMAP.md#validation-and-operational-acceptance).

## Evidence, limitations, impact and rollback

Executed document checks and changed-file secret-scan scope/output are retained
in [validation.txt](evidence/web-01a/validation.txt) and
[secrets.txt](evidence/web-01a/secrets.txt). These establish only the checks named
there. Source inspection and test names do not establish new build/test passes.
No application/runtime, integration, browser/accessibility, penetration or
performance execution is claimed for this documentation slice. No physical ARM
hardware evidence, ARM GNU toolchain validation or Windows runtime validation is
available for this task. External links identify issues/inputs; no external
research or fresh governance/hardware audit is part of this local inventory.

The change adds documentation only: no dependency, production configuration,
credential, route, schema migration, agent allocation, flash write or root
requirement. It changes no security policy and establishes no production security
guarantee. Pattern scans cannot certify absence of every secret. Resource and
security effects of future implementation remain review/measurement questions.

Rollback: supersede this report or revert its documentation/roadmap reference
through an isolated feature branch and Draft PR into `develop`, with human review
and manual merge. Preserve source/approval provenance when superseding it. No
service, device, database or credential state is changed; no data migration or
credential recovery action is needed. Never restore revoked credentials through
a future implementation rollback, per the [security baseline](../../security/README.md#reproduction-impact-and-rollback).
