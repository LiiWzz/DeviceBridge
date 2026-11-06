# ADR-003: Propose an interactive SSH terminal in the dashboard

**Status: Proposed — human Tech Lead scope/security approval required.**
Issue: [#45](https://github.com/LiiWzz/DeviceBridge/issues/45).
Source baseline: `1090dbdfd401bfb07d85d8fdb12dec08b4a1045c`.
Owner request in the project conversation: “ในหน้าเว็บเลย” — select an in-page
terminal after being offered desktop/native CLI versus browser terminal.
This establishes the requested experience, not an approved security architecture.

## Proposed decision and baseline

Propose a separately gated M03 browser-terminal workstream: an authorized user
selects a device, chooses Connect, types commands in the dashboard and explicitly
disconnects. The target is an interactive SSH session with a pseudo-terminal,
not a form that runs one command at a time. Device shell behavior and permissions
come from its OpenSSH/Dropbear account, not from a privileged web-shell agent.

The [current roadmap](../ROADMAP.md) explicitly defers browser terminal, and
[UX-02](../ux/UX-02.md) / [UX-03](../ux/UX-03.md) provide native CLI guidance only.
Those artifacts remain the historical/current baseline pending this proposal's
approval and a reviewed UX extension. [ADR-001](ADR-001-foundation.md) places
the dashboard outside M01/M02 and requires SSH-native access, verified host keys,
outbound device connectivity and no relay SSH private keys.
[ADR-002](ADR-002-device-record-identity-binding.md) remains Proposed; an inventory
record or enrollment never grants session authority. No ADR status is changed here.

Existing source supplies no terminal backend: [control `service.handle`](../../server/control/internal/control/server.go)
serves probes, [relay `Disabled.Open`](../../server/relay/internal/transport/transport.go)
rejects streams, [CLI `run`](../../cli/cmd/devicebridge/main.go) performs local
checks, and [agent `run`](../../agent/src/main.rs) idles without networking.
This is source inspection, not executed terminal or connectivity evidence.

If approved, this ADR changes the proposed M03 product scope only. It does not
choose the SSH client location, credentials, provider, terminal library, protocol,
routes, roles, timeouts or resource limits. Those architecture/security choices
must be approved before implementation; M01/M02 acceptance and full WEB-01 remain
prerequisites. No browser terminal is implemented by this documentation.

## Proposed operator experience

```text
Devices -> device detail -> Connect -> authorized session preparation
  -> verify intended device and trusted SSH host key
  -> terminal attached -> type commands -> Disconnect -> closed
  -> denied / expired / revoked / unavailable -> no usable terminal
  -> uncertain start or lost connection -> reconcile original session
```

The future UI should show the selected device, remote account and session state
before accepting input. Acceptance must demonstrate typing, Enter, Ctrl+C,
interactive `top`/`vi`, terminal resize and a clear Disconnect control on supported
device/browser combinations. These are desired behaviors, not runnable device
commands or a claim those programs exist on every image. Missing shell/programs,
unsupported PTY behavior and denied SSH accounts must yield explicit failures.
No clipboard, file transfer, recording, multi-user sharing, unattended command
execution or root escalation is implied by the interactive-terminal request.

## Architecture alternatives requiring a separate decision

| Option | Proposed evaluation / consequence |
| --- | --- |
| Native CLI only | Retains the existing planned trust boundary and client host-key verification, but does not meet the owner's in-page interaction request. Keep it as the separately authorized fallback; never bypass a browser denial by automatically opening it. |
| Dedicated server-side SSH terminal gateway | Candidate: browser terminal input/output connects to a separately isolated gateway acting as the SSH client over the accepted M02 relay. Because SSH terminates at that gateway, it sees terminal plaintext and needs a reviewed device-account authentication method. This is a new trusted component, not the existing relay/control probe service. No relay private-key custody, user private-key upload, shared fleet credential or plaintext logging is authorized. Exact gateway credentials/custody, isolation, origin/session binding and terminal transport remain blockers; this ADR does not approve that architecture. |
| Browser-side SSH client with opaque transport bridge | Candidate: SSH terminates in the browser, with a bridge carrying encrypted SSH bytes. Requires feasibility/interoperability evidence plus browser-side identity/key handling, trusted host-key delivery, application integrity/XSS review and bounded bridge authorization. A browser terminal renderer alone is not an SSH implementation. No library or ability to protect keys from a compromised browser is assumed. |

The browser endpoint necessarily sees displayed output and typed input. The
server-side option would additionally expose them to the SSH gateway; an opaque
relay would still see encrypted SSH payloads. Do not describe an SSH-terminating
gateway as blind, move SSH termination into the relay by renaming it, or infer
end-to-end browser/device encryption from browser HTTPS alone.

The architecture decision must select a concrete, reviewed implementation of
SSH/terminal transport, credential lifecycle and host-key trust before code. It
must preserve the [no custom device shell protocol / unauthenticated forwarding
rules](../../AGENTS.md). No generic command-execution endpoint or shell RPC may
replace the device's authenticated SSH server. Any required exception needs its
own explicit policy ADR approval; none is granted here.

## Security, failure and acceptance obligations

All rows below are proposed obligations for follow-up design and executed tests,
not existing controls or a selected wire/error contract.

| Case / boundary | Required design and future evidence |
| --- | --- |
| User -> selected device/account | Separate session authorization from visibility/enrollment. Direct requests, attachment and input must stay bound to the authorized user, intended device/account and original session; reject swapped IDs and replayed/cross-user handles. Bound SSH account choices and destinations; never accept arbitrary host/port/command from the browser. |
| Browser origin/session -> terminal transport | Authenticate start and attach/upgrade, enforce approved origin and session protections, reject cross-site/cross-origin abuse and token exposure in URLs/logs. Choose CSRF/session/transport mechanisms in full WEB-01, including concurrent tabs and logout behavior. UI controls alone are not authorization. |
| SSH peer / credentials | Independently verify the intended SSH host key before shell input. Unknown/mismatched keys, invalid credentials and expired/revoked grants fail closed without bypass controls. Define secure provisioning/custody/cleanup; no SSH private keys on relay, no private-key upload to control/relay, no fleet credential. |
| Expired, revoked, denied or unavailable | Keep distinct UI states and reject new unauthorized starts/attachments. Define approved detection/active-session closure deadlines and input rejection behavior, including control-plane outages; fresh authorization is required for a new session. Never reuse revoked authority on reconnect. |
| Start response lost / reconnect / duplicate click | Disable blind duplicate starts; reconcile the authorized original session. Never replay keystrokes, pasted text or buffered commands on reconnect. Choose whether explicit reattach to the same session is supported; otherwise close/reconcile it before a fresh request. |
| Disconnect, tab close, network loss or gateway restart | Specify server-owned bounded cleanup even when the browser cannot signal closure, including PTY/channel/account credentials and orphan sessions. Distinguish closing the SSH session from terminating all remote descendants: background processes may outlive it, so do not promise command rollback or process-tree cleanup without proved semantics. |
| Untrusted device output / browser compromise | Treat terminal escape sequences/output as hostile data, never executable HTML. Review link handling, clipboard-control sequences, renderer/add-ons, CSP, supply-chain and XSS. No automatic clipboard read/write or output-triggered browser actions; define safe bounded scrollback and paste/multiline confirmation. Device text must not masquerade as trusted authorization UI. |
| Audit and privacy | Record authorized session lifecycle/outcomes and actor/device/account attribution under approved audit access/retention/failure policies. Do not record keystrokes, screen contents or pasted secrets by default. Any recording feature needs separate explicit privacy/security approval; audit failure and lost responses require reconciliable outcomes before access is enabled. |
| Slow peers, output floods and fleet load | Approve aggregate/per-user/per-device limits for sessions, pending starts, handshake state, buffers, output/scrollback, queues, retries and idle/absolute lifetimes. Resize storms, blocked consumers, binary/invalid output and concurrent 100-device workloads require bounded behavior and measurements. |

## Dependency-ordered follow-up and owners

The names below are proposed work packages, not created implementation issues.
The Tech Lead assigns accountable owners and scopes issues one at a time.

| Proposed package | Owner roles / start gate / deliverable |
| --- | --- |
| BTERM-DESIGN | Architect / Security / Backend / Agent / Frontend: review this scope, choose SSH termination and authentication/custody, threat model, host-key verification, origin/session/input controls and lifecycle/audit semantics in approved ADRs. Depends on agreed M02 D01/D02/D03 inputs; cannot infer their approval from this proposal. |
| BTERM-UX | UX / Frontend: reviewed interactive-terminal journey and keyboard/accessibility states, account/device context, safe paste, errors, reconnect and disconnect; extend UX-03 explicitly. Synthetic design preview may precede production gates, but no active shell is authorized by a prototype. |
| BTERM-SERVICE | Backend / Security / Relay: after M01/M02 acceptance, full WEB-01 and architecture/security approvals, implement the selected bounded authorized session/transport integration with audit, host-key verification and cleanup evidence. Browser renderer is not a substitute for these prerequisites. |
| BTERM-UI | Frontend / QA: after matching accepted backend and UX/contracts, implement the in-page terminal and state handling. Execute real-browser/keyboard/assistive-tech and supported-device PTY tests; do not ship a mock as connected functionality. |
| BTERM-ACCEPT | Independent reviewer / QA / Hardware QA / Human Tech Lead: validate deny/revoke/replay/origin/output/failure/resource cases, actual SSH/Dropbear interoperability, accepted audit behavior, deployment recovery and signed supported scope before rollout. |

This supplements the [M03 roadmap](../ROADMAP.md) if approved, rather than
waiving [M02](../M02_EXECUTION.md) or [WEB-01](../web/WEB-01a.md) prerequisites.
Native CLI access remains separately gated and is not removed by the proposal.

## Impact, rollout, evidence and rollback

The present change adds documents only: no listener, browser script/library,
agent runtime, configuration, credential or executable command. Future gateway
and browser resource costs, attack surface and plaintext/key custody depend on
the selected architecture; no security or performance guarantee follows here.
Device agent still forwards authorized encrypted SSH streams; browser assets and
terminal gateway runtime stay off the 128 MB device. The <=16 MB idle RSS target
and session/PTY overhead require real ARM measurements in the supported matrix.

Future rollout must be disabled until approved acceptance, allow controlled
enablement for explicitly authorized users/devices, and support stopping new
browser starts while closing/reconciling existing sessions under approved rules.
Rollback disables browser access and restores the approved deployment without
reviving grants, bypassing host keys or forgetting sessions. Closing the page or
reverting code cannot undo commands already run. Preserve audit and secrets
revocation/cleanup obligations; CLI fallback retains its independent checks.

Document/link/scope and secret-scan executions are captured in
[validation.txt](../evidence/browser-terminal/validation.txt) and
[secrets.txt](../evidence/browser-terminal/secrets.txt). This proposal uses local
repository sources, not verified library/provider/API facts. No live browser,
SSH/PTY, build/security/cryptographic test, Windows/ARM toolchain or physical ARM
result is claimed. BrowserOS neo MCP remains unavailable. Pattern scanning cannot
certify absence of every secret. Revert/supersede these documentation additions
through a reviewed PR into `develop`; no runtime/session state is created now.

## Approval record

**PENDING; intentionally unsigned.** Tech Lead must record the exact ADR commit,
accept/revise/reject decision, scope/architecture approval boundaries, accountable
follow-up owners and authoritative identity/time reference. The owner's product
preference, merge or agent review is not that approval. No dependent contract or
security implementation proceeds without the required decisions under
[AGENTS.md rule 9](../../AGENTS.md).
