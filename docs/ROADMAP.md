# DeviceBridge product roadmap and UX/UI plan

Planning issue: [DB-012 / #19](https://github.com/LiiWzz/DeviceBridge/issues/19).

UX-01 proposal: [operator journeys and information architecture](ux/UX-01.md),
tracked in [#27](https://github.com/LiiWzz/DeviceBridge/issues/27). The owner
approved its journeys, inventory and capability separation as inputs to UX-02
on 2026-11-05 in the project conversation; see its limited review record.
The [UX-02 prototype packet](ux/UX-02.md), tracked in
[#29](https://github.com/LiiWzz/DeviceBridge/issues/29), received owner flow
approval and selection of A (fleet sidebar) on 2026-11-05 in the project
conversation. The [UX-03 design system and handoff](ux/UX-03.md), tracked in
[#31](https://github.com/LiiWzz/DeviceBridge/issues/31), is proposed for human
review. Production integration remains gated separately.

The preliminary [WEB-01a backend gap inventory](web/WEB-01a.md), tracked in
[#33](https://github.com/LiiWzz/DeviceBridge/issues/33), maps source capabilities
to dashboard needs and records contract prerequisites. This discovery slice
does not satisfy WEB-01's agreed M02 identity/authorization dependency or select
API/security contracts. UX-03 review and all implementation gates remain open.

The proposed [M02 execution plan](M02_EXECUTION.md), tracked in
[#35](https://github.com/LiiWzz/DeviceBridge/issues/35), breaks secure connectivity
into dependency-ordered planning packages and approval/evidence gates. It is
documentation for review; it neither closes M01 acceptance nor authorizes M02
implementation or selects web/security contracts.

## Status and boundaries

The repository currently provides M01 foundations; enrollment, working tunnels,
remote SSH and a management dashboard are not implemented. This document adds
explicit UX/UI and dashboard delivery work. It does not certify M01 completion
or approve new API/security contracts.

[M01 execution](M01_EXECUTION.md) remains the foundation acceptance contract.
[ADR-001](adr/ADR-001-foundation.md) continues to govern architecture: the browser
dashboard remains outside M01/M02 delivery and is scheduled here as M03. Its
existing status and decisions are unchanged. Any future API, identity,
authorization or other security contract changes need a Proposed ADR and human
Tech Lead approval before implementation, as required by [AGENTS.md](../AGENTS.md).

Dates and effort estimates are not assigned. The Tech Lead sets them after work
packages are scoped into issues and prerequisite decisions are approved. The
work-package IDs below are planning references, not created GitHub issues or
claims that work has started. Execute one linked GitHub issue at a time in an
isolated feature branch; submit Draft PRs into `develop` for human review and
manual merge.

## Milestones and release gates

| Milestone | Scope | Exit gate |
| --- | --- | --- |
| M01 — Foundation | Rust agent, Go control/relay/CLI skeletons, PostgreSQL/Compose, CI, security baseline and integration smoke; portable installation previews are tracked separately in DB-011 / #16 | DB-001..010 gates from the execution plan, independent review and human Tech Lead GO/NO-GO; packaging has its own issue acceptance |
| M02 — Secure connectivity | Planned device identity/enrollment, session authorization, outbound TLS relay and SSH-native CLI access behind NAT/CGNAT | Separate execution plan and approved contracts; demonstrated allow/deny/expiry behavior, host-key verification, reconnect handling and real ARM resource evidence; independent review and Tech Lead acceptance |
| M03 — Management dashboard | UX/UI design, authenticated web management, fleet/device visibility, session management, users/permissions and audit history | UX/UI approval, accepted M02 prerequisites, approved web/API/security contracts, dashboard functional/security/accessibility evidence and Tech Lead acceptance |

M02 scope above and the linked execution plan are proposals for human review.
M02 implementation follows M01 acceptance. M03 discovery, wireframes and a
clearly labeled prototype using synthetic data may start earlier. Production
integration and rollout wait for M01/M02 acceptance and the web contract gate.

## Dashboard MVP screens

| Screen / flow | Operator outcome | Required behavior |
| --- | --- | --- |
| Sign-in and access denied | Enter the management system and understand access restrictions | Define login, logout, session-expired and forbidden states; authentication method is decided at the contract gate |
| Fleet overview | See fleet availability and recent operational activity | Counts agree with authorized device results; show freshness and unavailable/unknown data explicitly |
| Devices | Find an authorized device quickly | Search, status filter and pagination; online/offline/unknown status with last-seen time; no device credentials in results |
| Device detail and onboarding | Identify a device, inspect connectivity and follow enrollment/revocation steps | Show agreed metadata and last-seen time; design safe enrollment instructions and confirmation/result states for access revocation; credentials and token delivery await the approved contract |
| Sessions | Inspect authorized active sessions and session history; request access through the native CLI | Provide safe CLI handoff without embedded credentials; show denied/expired requests; termination is included only with backend authorization, confirmation and audit support |
| Users and permissions | Manage who can inspect devices and request access | Define a permission matrix for viewing, connecting and administering; show denied actions; backend enforces every permission; role names/model await contract approval |
| Audit history | Trace access requests and administrative changes | Filter by actor, device, action, time and result; redact credentials; audit fields, retention and visibility are agreed at the contract gate |

Browser terminal/SSH key custody is outside the MVP. Remote access continues
through OpenSSH/Dropbear and the native CLI; the web interface must not introduce
a custom shell protocol or move private SSH keys onto the control service/relay.

## UX/UI and implementation work packages

| ID | Deliverable | Dependencies / start gate | Owner role | Acceptance |
| --- | --- | --- | --- | --- |
| UX-01 | Operator journeys, information architecture, screen inventory and draft permission matrix | Roadmap review; can begin before M02 delivery | Product / UX | Tech Lead reviews find-device, onboarding, access-request, revocation and audit journeys, including unavailable/denied paths |
| UX-02 | Wireframes and clickable prototype for every MVP screen | UX-01 | UX/UI | Clearly labeled synthetic data; loading, empty, error, expired, forbidden and stale-data states; desktop/mobile layouts; human flow review recorded |
| UX-03 | UI design system and implementation handoff | UX-02 review | UX/UI / Frontend | Navigation, typography, spacing, color/status semantics, tables, forms and dialogs documented; keyboard/focus behavior, contrast and non-color status cues specified |
| WEB-01 | Backend gap analysis, API/schema/error definitions, authentication/session design, permission enforcement and audit contract | UX-01 plus agreed M02 identity/authorization design | Backend / Security / Architect | Identify missing device, enrollment, session, user and audit APIs; submit ADRs for contract/security changes; record Tech Lead approval before implementing them |
| WEB-02 | Approved backend capabilities and web application foundation | WEB-01 approval, UX-03 and M02 acceptance | Backend / Frontend | Backend gaps implemented in separate dependency-ordered issues; authenticated shell, routing, authorized API client and reusable UI states; framework/hosting chosen with recorded rationale |
| WEB-03 | Fleet overview, device list/detail and onboarding/revocation flows | WEB-02 and matching accepted backend APIs | Frontend / Backend | Authorized results, search/filter/pagination and freshness work; enrollment/revocation success and failure states follow approved contracts; responsive and keyboard flows verified |
| WEB-04 | Session management, users/permissions and audit screens | WEB-02 and matching accepted backend APIs | Frontend / Backend | Approved access matrix enforced server-side; CLI handoff, expiry/deny paths, supported session termination and redacted audit entries verified |
| WEB-05 | Integration/security/accessibility/performance validation and release runbook | WEB-03 and WEB-04 | QA / Security / DevOps | Reproducible evidence, negative cases, measured server/browser resource behavior for a 100-device fleet, rollout/rollback instructions and independent review |
| WEB-06 | M03 acceptance report | WEB-05; critical/high findings closed | Human Tech Lead | Reviewed evidence, limitations and explicit GO/NO-GO; manual merge/release decisions only |

Each package must be decomposed into small GitHub issues before implementation;
packages that cover backend and frontend work are not a single coding assignment.
Approve UX-01/02/03 deliverables before implementation of their dependent UI.
Backend contract design and visual design may inform each other; unresolved
contract decisions remain visible blockers to production integration.

## Validation and operational acceptance

- Functional: trace each MVP screen and journey to its issue, approved design,
  backend contract and integration evidence. Counts, filters and pagination must
  remain correct for authorized results in a 100-device test fleet.
- Negative/security: verify unauthenticated and unauthorized access, direct API
  calls that bypass hidden UI controls, access to unauthorized device/session/user
  IDs, expired/revoked access, invalid inputs, duplicate submissions and upstream
  outages. Review session protection, CSRF where applicable, XSS, secret redaction
  and audit completeness under the approved design.
- UX/accessibility: cover loading, empty, error, offline, unknown/stale, forbidden
  and expired states. Verify keyboard navigation, focus restoration, form labels,
  contrast, non-color status cues and responsive layouts. Adopt WCAG 2.2 AA as the
  planned accessibility target; capture automated and human review evidence.
- Resources: serve the dashboard on server/browser infrastructure. Add no web
  runtime or UI assets to embedded devices. Use bounded/paginated APIs and define
  refresh/backoff behavior so visibility does not create unbounded load. Set
  measurable browser/API latency, memory and refresh budgets at WEB-01/02, then
  record actual results at WEB-05. Any agent-side telemetry changes require their
  own issue and ARM measurements against the <=16 MB idle RSS target on 128 MB
  devices; cross-builds do not establish hardware acceptance.
- Delivery evidence: record exact build/lint/test commands and outputs for the
  selected stack, browser coverage, toolchain/hardware gaps, security review,
  performance implications and rollback steps in every implementation PR.
  Synthetic prototypes never establish working connectivity or security.
- Rollout/rollback: WEB-05 must document how to disable dashboard access and
  restore the prior approved deployment without bypassing authentication. Any
  schema migration needs a reviewed compatibility and rollback/recovery plan;
  preserve audit integrity and existing CLI access.

## Explicitly deferred

P2P, QUIC, video streaming, browser terminal, bulk remote command execution,
firmware/OTA management and advanced analytics are outside this dashboard MVP.
They need separately scoped issues and, where contracts or security change,
approved ADRs. No frontend framework, identity provider, production deployment
configuration or delivery date is selected by this planning change.
