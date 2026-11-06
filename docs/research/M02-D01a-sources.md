# M02-D01a — Identity-binding source findings

Issue: [#43](https://github.com/LiiWzz/DeviceBridge/issues/43).
Inspected source base: `e32cd6630b72452dd7c50bc0c448eef70e30d515`,
branch `docs/m02-d01a-identity-binding`, 2026-11-06; verified with
`git rev-parse HEAD` and `git status --short --branch` before writing.

**Status: local primary-source research for a proposed ADR; no approval.**
This report follows the source/symbol convention of
[WEB-01a](../web/WEB-01a.md#observed-foundation-and-its-limits).
Observed behavior below is source inspection, not executed runtime evidence.

## Implemented foundation versus desired behavior

| Area / owning source | Observed at the source base | Design implication, not implemented capability |
| --- | --- | --- |
| Record storage: [baseline migration](../../deploy/migrations/0001_baseline.up.sql); [migration runner](../../deploy/migrate.sh), `devicebridge_migrations.history` | Baseline explicitly has no application tables. Runner transaction/advisory lock protects migration bookkeeping. | No persistent device record, identity-binding or audit store exists. Migration atomicity does not establish enrollment atomicity. |
| Control API: [`service.handle`, `service.serve`, `newService`](../../server/control/internal/control/server.go) | Exact health/readiness probes only; other paths return 404. Lifecycle logging describes process serving/shutdown. | No create/enroll/reconcile/authentication/authorization endpoint or authoritative enrollment audit producer exists. A ready process proves none of these. |
| Agent storage/runtime: [`Config`, `Config::load`, `Config::read`](../../agent/src/config.rs); [`execute`, Unix `run`](../../agent/src/main.rs) | Read-only bounded local configuration contains logging/polling settings. Runtime idles without network or disk writes; no identity field or persistent identity store. | Device identity provisioning, protected storage and crash-safe binding remain future work; do not treat installation/configuration as enrollment. |
| CLI/relay: [CLI `run`, `usage`](../../cli/cmd/devicebridge/main.go); [relay `Connector`, `Disabled.Open`](../../server/relay/internal/transport/transport.go) | CLI validates local settings without contacting a server; remote commands unavailable. Disabled transport rejects opening streams; its interface explicitly defines no auth/enrollment protocol. | No enrollment command, authenticated connection or native SSH handoff can be reused. No runnable bootstrap syntax should be inferred. |
| UX record creation: [`devices`, `state.draft`, `finishCreate`, submit event handler](../ux/prototype/prototype.js); [UX-04 form/results](../ux/UX-04.md#form-results-and-inventory-behavior) | Synthetic metadata becomes an in-memory `demo-new-…` record, Not enrolled, observation Unknown, no last-seen evidence. Name/site/model are display input. | This demonstrates record-before-enrollment UX only. A display label or demo identifier is not an authenticated device identity, online observation or SSH permission. |
| UX uncertainty/audit: [`state.draft.pending`, reconciliation click handler, `events`, `audit`](../ux/prototype/prototype.js); [UX-04 handoff](../ux/UX-04.md#enrollment-and-backend-handoff) | Page-local pending lock and simulated lookup; reload loses state. Audit rows are fixed synthetic examples, not creation events. | No durable deduplication, concurrent-client exclusion, authoritative reconciliation or audit guarantees follow from the prototype. |
| SSH host keys: [ADR-001 decisions/trade-offs](../adr/ADR-001-foundation.md); [security lifecycle table](../../security/README.md#proposed-credential-and-token-lifecycle--unimplemented) | Verification is required/proposed for future native SSH; compatibility and pinning remain unverified. | Device enrollment identity and SSH server host-key trust are separate questions. Enrollment must not silently approve a host key or authorize SSH. |

The [security baseline opening](../../security/README.md) explicitly confirms no
enrollment, credential store, authentication service, session authorization, TLS
transport or remote SSH. Its lifecycle table proposes operator-authorized,
device-bound, one-use bootstrap and unique device identities; it selects no
credential mechanism. [M02 D01–D03 and P01–P03](../M02_EXECUTION.md#dependency-graph-and-accountable-evidence)
likewise describe future design/enforcement/audit/enrollment prerequisites.

## Conclusions for the proposed ADR-002

The following are recommendations derived from those gaps, not existing policy
or proof that a binding mechanism exists. ADR-002 should separate a management
record from the enrolled device identity and require explicit operator authority
for binding the intended device to the intended record. Metadata creation alone
must confer no enrollment, connectivity or SSH capability. This carries forward
[UX-04's handoff](../ux/UX-04.md#enrollment-and-backend-handoff) and
[M02's D01 obligations](../M02_EXECUTION.md#decisions-the-adrs-must-resolve).

Propose atomic binding and authoritative reconciliation as requirements, without
selecting transactions, uniqueness keys, persistence fields or wire semantics.
Unknown outcomes must remain unresolved until authoritative evidence determines
the original result; a timeout or absent audit event cannot establish failure.
This extends the prototype's review intent to a future durable design, not its
page-local implementation; see [WEB-01a negative questions](../web/WEB-01a.md#negative-security-and-resource-questions-to-retain).

| Threat / failure | Required design question and future evidence |
| --- | --- |
| Wrong device/record or unauthorized operator | How does independently verified device identity bind to operator-authorized record scope, including requests outside the UI? Display metadata is not proof. |
| Replay/expiry | How is an expired, consumed or wrong-device enrollment rejected without a shared fleet credential? Secure bootstrap delivery and one-use enforcement remain undecided. |
| Concurrent enrollment / record edits | How does one authoritative result prevent conflicting device/record bindings and inconsistent audit outcomes across clients/restarts? |
| Lost response / restart | How is the original result reconciled without duplicate enrollment, blind retry or granting authority while outcome is unknown? |
| Cloned device material | What proves the intended physical device and detects/contains clones? Unique labels alone cannot supply clone resistance; protected storage and recovery require later decisions. |
| Replacement / re-enrollment | How is replacement explicitly authorized, old identity invalidated and history preserved without silently rebinding or inheriting SSH host-key trust? Lifecycle/host-key recovery remains D02/D03 work. |

These are proposed acceptance questions grounded in the
[security lifecycle](../../security/README.md#proposed-credential-and-token-lifecycle--unimplemented)
and [M02 evidence table](../M02_EXECUTION.md#validation-and-evidence-gates).
ADR-002 must leave credential format, bootstrap wire/delivery mechanism, provider,
cryptography, API and schema choices open. Detailed user/session authorization,
audit/lifecycle/outage semantics and host-key trust need their scoped decisions.

## Limits, approvals and impact

Only repository sources were inspected; no external research, live GitHub
governance audit, browser checks, builds or security/runtime tests were run.
Existing test names/evidence are historical, not fresh passes. No Windows or ARM
toolchain execution or physical ARM measurements establish compatibility here.

[AGENTS.md rule 9](../../AGENTS.md), [ADR-001 status](../adr/ADR-001-foundation.md)
and [M02 blocking gates](../M02_EXECUTION.md#existing-foundation-and-blocking-gates)
require explicit human Tech Lead approval of contract/security ADRs before
dependent implementation, plus M01 acceptance. Submission/merge of a proposal
does not record that approval; audit foundations precede enrollment mutations.

This report adds no runtime dependency, allocation, flash write, privilege,
credential or production configuration. Future binding/storage changes need
security review and measured ARM resource/crash-recovery evidence. Rollback is
superseding this report or a reviewed documentation revert through a Draft PR
to `develop`; no device, database or credential state changes.
