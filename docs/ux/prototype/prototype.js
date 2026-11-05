/* Throwaway UX-02. Question: fleet navigation, attention-first board, or device workbench?
 * Three structural variants on this local page via ?variant=A/B/C. No production route.
 * All mutations affect memory only. No real commands, identities, credentials or network.
 */
'use strict';
const localPreview = location.protocol === 'file:' || ['localhost', '127.0.0.1', '[::1]'].includes(location.hostname);
const app = document.querySelector('#app');
const variants = { A: 'Fleet sidebar', B: 'Operations board', C: 'Device workbench' };
const screens = { overview: 'Overview', devices: 'Devices', detail: 'Device detail', onboarding: 'Onboarding', sessions: 'Sessions', users: 'Users & permissions', audit: 'Audit history', signin: 'Sign-in' };
const scenarios = ['ready', 'loading', 'empty', 'error', 'expired', 'forbidden', 'stale'];
const devices = [
  { id: 'demo-001', name: 'North pump station', site: 'North · synthetic', status: 'online', seen: '05 Nov 2026, 15:58 UTC', model: 'ARMv7 sample' },
  { id: 'demo-002', name: 'Workshop gateway', site: 'Workshop · synthetic', status: 'offline', seen: '05 Nov 2026, 13:20 UTC', model: 'AArch64 sample' },
  { id: 'demo-003', name: 'East meter', site: 'East · synthetic', status: 'unknown', seen: 'No observation available', model: 'ARMv7 sample' },
  { id: 'demo-004', name: 'South pump station', site: 'South · synthetic', status: 'online', seen: '05 Nov 2026, 15:57 UTC', model: 'AArch64 sample' },
  { id: 'demo-005', name: 'Lab controller', site: 'Lab · synthetic', status: 'offline', seen: '04 Nov 2026, 09:00 UTC', model: 'ARMv7 sample' },
  { id: 'demo-006', name: 'West gateway', site: 'West · synthetic', status: 'unknown', seen: 'No observation available', model: 'AArch64 sample' }
];
const events = [
  { actor: 'Demo operator', device: 'demo-001', action: 'Access request', time: '2026-11-05', result: 'Denied' },
  { actor: 'Demo maintainer', device: 'demo-002', action: 'Enrollment review', time: '2026-11-04', result: 'Pending' },
  { actor: 'Demo operator', device: 'demo-004', action: 'Access request', time: '2026-11-05', result: 'Expired' }
];
const params = new URLSearchParams(location.search);
const state = { variant: Object.hasOwn(variants, params.get('variant')) ? params.get('variant') : 'A', screen: Object.hasOwn(screens, params.get('screen')) ? params.get('screen') : 'overview', scenario: 'ready', device: 'demo-001', query: '', filter: 'all', page: 0, session: 'Denied', outcome: '', audit: { actor: '', device: '', action: '', date: '', result: '' } };
const esc = value => String(value).replace(/[&<>"']/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
const button = (text, action, extra = '') => `<button type="button" data-action="${action}" ${extra}>${text}</button>`;
const badge = value => `<span class="status ${esc(value)}">${esc(value[0].toUpperCase() + value.slice(1))}</span>`;
const selectedDevice = () => devices.find(d => d.id === state.device);
const announce = text => { document.querySelector('#announcement').textContent = text; };
function nav() {
  return ['overview', 'devices', 'sessions', 'users', 'audit'].map(key => button(screens[key], 'navigate', `data-screen="${key}" ${state.screen === key || key === 'devices' && ['detail', 'onboarding'].includes(state.screen) ? 'aria-current="page"' : ''}`)).join('');
}
function VariantA(content) {
  return `<div class="layout-a"><aside class="sidebar"><div class="brand"><span class="mark">D</span>DeviceBridge</div><p class="eyebrow">Fleet workspace</p><nav aria-label="Main">${nav()}</nav><p class="note">Demo workspace<br>6 synthetic devices<br><br>No live connection</p>${button('Sign-in preview', 'navigate', 'data-screen="signin"')}</aside><main id="content" tabindex="-1">${content}</main></div>`;
}
function VariantB(content) {
  return `<div class="layout-b"><header class="topbar"><div class="brand"><span class="mark">D</span>DeviceBridge</div><nav aria-label="Main">${nav()}${button('Sign-in', 'navigate', 'data-screen="signin"')}</nav></header><section class="ops-intro" aria-label="Demo attention queue"><div><p class="eyebrow">Operations · synthetic snapshot</p><p class="attention">4 devices need a closer look</p><p>2 offline · 2 unknown · observations fixed at 16:00 UTC</p></div>${button('Inspect attention queue', 'attention')}</section><main id="content" tabindex="-1">${content}</main></div>`;
}
function VariantC(content) {
  return `<div class="layout-c"><header class="topbar"><div class="brand"><span class="mark">D</span>DeviceBridge / Workbench</div><div class="inline">${button('Fleet overview', 'navigate', 'data-screen="overview"')}${button('Sign-in', 'navigate', 'data-screen="signin"')}</div></header><div class="workbench"><aside class="device-rail"><p class="eyebrow">Choose a device</p><div class="rail-devices">${devices.map(d => button(`${esc(d.name)}<br><small>${esc(d.id)} · ${esc(d.status)}</small>`, 'device', `data-id="${d.id}" class="${state.device === d.id ? 'selected' : ''}"`)).join('')}</div><p class="review-note">Context persists across sessions and audit views. Every record here is invented.</p></aside><main id="content" tabindex="-1"><nav class="tabs" aria-label="Workspace">${nav()}</nav><hr>${content}</main></div></div>`;
}
function scenarioControls() {
  return `<div class="scenario"><label for="scenario">Review a synthetic state <select id="scenario">${scenarios.map(s => `<option ${s === state.scenario ? 'selected' : ''}>${s}</option>`).join('')}</select></label><p class="review-note">These controls change the demonstration, not permissions or a service.</p></div>`;
}
function stateContent() {
  const views = {
    loading: ['Loading demo records', 'Waiting for a simulated response. No current status is established.', 'Show loaded example'],
    empty: ['No demo records', 'This illustrates an empty authorized view. It is different from a failed request.', 'Return to populated example'],
    error: ['Information unavailable', 'A simulated request failed. Missing results do not mean there are no devices or events.', 'Retry demo'],
    expired: ['Session expired — example', 'In the product, sign in again through the approved authentication flow. This prototype collects no credentials.', 'View sign-in example'],
    forbidden: ['Access not permitted — example', 'No resource details are shown. Follow your organization’s approved access request process.', 'Back to overview']
  };
  const entry = views[state.scenario];
  if (!entry) return '';
  return `<section class="panel state-panel"><h2>${entry[0]}</h2><p>${entry[1]}</p>${button(entry[2], state.scenario === 'expired' ? 'signin' : state.scenario === 'forbidden' ? 'home' : 'ready', 'class="primary"')}</section>`;
}
function overview() {
  const metric = status => devices.filter(d => d.status === status).length;
  const stats = `<div class="grid">${['online', 'offline', 'unknown'].map(status => `<div class="panel metric"><span>${esc(status[0].toUpperCase() + status.slice(1))}</span><strong>${metric(status)}</strong>${button('View devices', 'status', `data-status="${status}"`)}</div>`).join('')}<div class="panel metric"><span>Demo inventory</span><strong>${devices.length}</strong><small>Fixed snapshot · 16:00 UTC</small></div></div>`;
  const queue = `<section class="panel"><h2>Inspect before requesting access</h2><p>Offline and unknown observations need context. Online status grants no access.</p>${devices.filter(d => d.status !== 'online').map(d => `<p class="inline">${badge(d.status)} ${button(esc(d.name), 'device', `data-id="${d.id}"`)}</p>`).join('')}</section>`;
  const recent = `<section class="panel"><h2>Recent request examples</h2><ul class="timeline">${events.map(e => `<li><strong>${esc(e.result)} · ${esc(e.action)}</strong><br><small>${esc(e.actor)} · ${esc(e.device)} · ${e.time}</small></li>`).join('')}</ul>${button('Inspect audit history', 'navigate', 'data-screen="audit"')}</section>`;
  if (state.variant === 'B') return `${queue}<div class="split">${recent}<section>${stats}</section></div>`;
  if (state.variant === 'C') return `${detail(true)}${stats}`;
  return `${stats}<div class="split">${queue}${recent}</div>`;
}
function deviceList() {
  const found = devices.filter(d => `${d.name} ${d.id} ${d.site}`.toLowerCase().includes(state.query.toLowerCase()) && (state.filter === 'all' || (state.filter === 'attention' ? d.status !== 'online' : d.status === state.filter)));
  const pages = Math.max(1, Math.ceil(found.length / 3));
  state.page = Math.min(state.page, pages - 1);
  return `<section class="panel"><div class="toolbar"><label for="search">Find a device<input id="search" type="search" placeholder="Name, ID or site" value="${esc(state.query)}"></label><label for="filter">Observation<select id="filter">${['all', 'online', 'offline', 'unknown', 'attention'].map(s => `<option ${s === state.filter ? 'selected' : ''}>${s}</option>`).join('')}</select></label>${button('Clear filters', 'clear')}</div><p class="muted">${found.length} of ${devices.length} synthetic devices · last-seen times in UTC</p>${found.length ? `<div class="table-wrap"><table><caption class="sr-only">Synthetic devices matching filters</caption><thead><tr><th>Device</th><th>Observation</th><th>Last seen</th></tr></thead><tbody>${found.slice(state.page * 3, state.page * 3 + 3).map(d => `<tr><td>${button(esc(d.name), 'device', `data-id="${d.id}"`)}<small>${d.id} · ${d.site}</small></td><td>${badge(d.status)}</td><td>${d.seen}</td></tr>`).join('')}</tbody></table></div>` : `<div class="state-panel"><h2>No matching devices</h2><p>Try another search or clear the filters. This is a successful query with no matches.</p></div>`}<div class="pagination">${button('Previous', 'page', `data-step="-1" ${state.page === 0 ? 'disabled' : ''}`)}<span>Page ${state.page + 1} of ${pages}</span>${button('Next', 'page', `data-step="1" ${state.page + 1 === pages ? 'disabled' : ''}`)}</div></section>`;
}
function detail(compact = false) {
  const d = selectedDevice();
  return `<section class="panel"><p class="eyebrow">Selected device · synthetic</p><h2>${esc(d.name)}</h2><dl><dt>Device ID</dt><dd>${d.id}</dd><dt>Site</dt><dd>${d.site}</dd><dt>Model</dt><dd>${d.model}</dd><dt>Observation</dt><dd>${badge(d.status)}</dd><dt>Last seen</dt><dd>${d.seen}</dd></dl><p class="notice">${d.status === 'offline' ? 'Offline example. The last observation does not establish current reachability.' : d.status === 'unknown' ? 'Unknown example. No observation is available; do not treat this as offline.' : 'Online observation is synthetic and does not authorize a connection.'}</p><div class="inline">${button('Access request preview', 'access', 'class="primary"')}${button('Onboarding guidance', 'navigate', 'data-screen="onboarding"')}${!compact ? button('Revoke access preview', 'revoke', 'class="danger"') : ''}${button('Back to device results', 'navigate', 'data-screen="devices"')}</div></section>${!compact ? `<section class="panel"><h2>Connectivity and identity</h2><p>Enrollment, identity verification and SSH access are not implemented. Native client host-key verification is required in the future approved flow. No ignore-verification shortcut is offered.</p></section>` : ''}`;
}
function onboarding() {
  return `<section class="panel"><h2>Onboard ${esc(selectedDevice().name)}</h2><ol class="steps"><li>Confirm the intended device identity and operator authority.</li><li>Follow enrollment instructions once the M02 contract is approved. No command or token is provided here.</li><li>Inspect an authoritative enrollment result separately from connectivity.</li></ol><p class="notice">Example only. A successful local health or configuration check does not establish enrollment.</p><div class="inline">${button('Preview pending enrollment', 'enroll')}${button('View device', 'navigate', 'data-screen="detail"')}</div></section>${state.outcome ? `<section class="panel" role="status"><h2>Simulated enrollment result</h2><p>${esc(state.outcome)}</p></section>` : ''}`;
}
function sessions() {
  return `<section class="panel"><h2>Access request example</h2><p>Device context: <strong>${esc(selectedDevice().name)}</strong> · ${state.device}</p><label class="field" for="request-result">Inspect simulated outcome<select id="request-result">${['Denied', 'Expired', 'Revoked', 'Unavailable', 'Handoff guidance'].map(s => `<option ${s === state.session ? 'selected' : ''}>${s}</option>`).join('')}</select></label><div class="notice" role="status">${esc({ Denied: 'Request denied. Viewing a device does not grant connection permission.', Expired: 'The example request expired. A future client must obtain fresh authorization.', Revoked: 'Access was revoked in this example. Do not reuse old credentials or silently restore identity.', Unavailable: 'Request outcome is unknown. Reconcile authoritative state before retrying.', 'Handoff guidance': 'Use the future approved native client flow with host-key verification. CLI syntax is not yet defined; no runnable command is shown.' }[state.session])}</div><p>No session is connected. Browser terminal, SSH private-key upload and session termination are outside this prototype.</p></section><section class="panel"><h2>Request history · synthetic</h2>${events.map(e => `<p>${badge(e.result.toLowerCase())} ${esc(e.device)} · ${e.time}</p>`).join('')}${button('Inspect audit examples', 'navigate', 'data-screen="audit"')}</section>`;
}
function users() {
  return `<section class="panel"><h2>Capabilities are separate decisions</h2><p>Illustrative people only. Role names, grants and resource scopes await approved contracts. These rows do not grant access.</p><div class="table-wrap"><table><caption class="sr-only">Proposed capability separation</caption><thead><tr><th>Example person</th><th>View devices</th><th>Request SSH</th><th>Administer</th><th>Inspect audit</th></tr></thead><tbody>${['Demo operator', 'Demo maintainer'].map(name => `<tr><th scope="row">${name}</th>${Array(4).fill('<td>To be decided</td>').join('')}</tr>`).join('')}</tbody></table></div><p>Viewing does not imply connecting. Administration does not imply SSH or unrestricted audit access. Enforcement belongs on the backend.</p>${button('Permission change preview', 'permission')}</section>${state.outcome ? `<p class="notice" role="status">${esc(state.outcome)}</p>` : ''}`;
}
function audit() {
  const rows = events.filter(e => Object.entries(state.audit).every(([key, value]) => (key === 'date' ? e.time : e[key]).toLowerCase().includes(value.toLowerCase())));
  return `<section class="panel"><h2>Trace an example request</h2><div class="audit-filters">${Object.entries({ actor: 'Actor', device: 'Device', action: 'Action', date: 'Date (UTC)', result: 'Result' }).map(([key, label]) => `<label class="field" for="audit-${key}">${label}<input id="audit-${key}" data-audit="${key}" type="${key === 'date' ? 'date' : 'search'}" value="${esc(state.audit[key])}"></label>`).join('')}</div><p class="review-note">Synthetic fields for design review. Visibility, retention and exact field contracts remain unresolved.</p><div class="table-wrap"><table><caption class="sr-only">Filtered synthetic audit history</caption><thead><tr><th>When (UTC)</th><th>Actor</th><th>Device</th><th>Action</th><th>Result</th></tr></thead><tbody>${rows.map(e => `<tr><td>${e.time}</td><td>${esc(e.actor)}</td><td>${e.device}</td><td>${esc(e.action)}</td><td>${badge(e.result.toLowerCase())}</td></tr>`).join('')}</tbody></table></div>${rows.length ? '' : '<p class="notice">No matching events. Missing evidence does not establish that nothing happened.</p>'}${button('Clear audit filters', 'clear-audit')}</section>`;
}
function signin() {
  return `<section class="panel state-panel"><p class="eyebrow">Sign-in concept</p><h2>Enter your fleet workspace</h2><p>Authentication provider and session design await approval. This demo accepts no username, password, key or token.</p>${button('Enter synthetic workspace', 'home', 'class="primary"')}<p>${button('Preview access denied', 'denied')}${button('Preview session expired', 'expired')}</p></section>`;
}
function render(focus = false) {
  if (!localPreview) {
    app.textContent = 'Local throwaway preview only. Open the HTML file or serve it on localhost. This is not a production dashboard.';
    return;
  }
  const views = { overview, devices: deviceList, detail, onboarding, sessions, users, audit, signin };
  const blocked = stateContent();
  const content = `<header class="page-head"><div><p class="eyebrow">Demo workspace / ${esc(screens[state.screen])}</p><h1>${esc(screens[state.screen])}</h1><p class="muted">Synthetic snapshot · 05 Nov 2026, 16:00 UTC</p></div><span class="status">Design review only</span></header>${scenarioControls()}${state.scenario === 'stale' ? '<p class="notice">Stale data example: last refresh failed. These historical observations do not establish current reachability.</p>' : ''}${blocked || views[state.screen]()}<details class="state-readout"><summary>Prototype state (for review)</summary><pre>${esc(JSON.stringify(state, null, 2))}</pre></details>`;
  app.innerHTML = ({ A: VariantA, B: VariantB, C: VariantC }[state.variant])(content);
  const switcher = document.querySelector('#switcher');
  switcher.hidden = false;
  switcher.innerHTML = `${button('←', 'variant', 'data-step="-1" aria-label="Previous layout"')}<span aria-live="polite">${state.variant} · ${variants[state.variant]}</span>${button('→', 'variant', 'data-step="1" aria-label="Next layout"')}`;
  if (focus) document.querySelector('#content').focus();
}
function syncURL() {
  const url = new URL(location.href);
  url.searchParams.set('variant', state.variant);
  url.searchParams.set('screen', state.screen);
  history.replaceState(null, '', url);
}
function navigate(screen) {
  state.screen = screen;
  state.scenario = 'ready';
  state.outcome = '';
  render(true);
  syncURL();
}
let confirmTrigger = '';
function confirmAction(kind) {
  const dialog = document.querySelector('#confirm');
  confirmTrigger = kind;
  const revoke = kind === 'revoke';
  dialog.innerHTML = `<h2 id="confirm-title">${revoke ? 'Revoke access preview' : 'Permission change preview'}</h2><p>${revoke ? `Target: ${esc(selectedDevice().name)} (${state.device}). In a real flow, verify scope and authority before submitting.` : 'Target: Demo operator. Exact grants and self-escalation controls require an approved policy.'}</p><p>All results below are simulated. No identity or permission changes.</p><label class="field" for="confirmation-outcome">Example outcome<select id="confirmation-outcome"><option>Unknown</option><option>Rejected</option><option>Confirmed</option></select></label><div class="toolbar">${button('Cancel', 'cancel')}${button('Simulate submission', 'submit-confirm', 'class="danger"')}</div>`;
  dialog.showModal();
  dialog.querySelector('[data-action="cancel"]').focus();
}
document.addEventListener('click', event => {
  const control = event.target.closest('button[data-action]');
  if (!control || !localPreview) return;
  const action = control.dataset.action;
  if (action === 'navigate') navigate(control.dataset.screen);
  else if (action === 'device') { state.device = control.dataset.id; navigate('detail'); }
  else if (action === 'variant') { const keys = Object.keys(variants); state.variant = keys[(keys.indexOf(state.variant) + Number(control.dataset.step) + keys.length) % keys.length]; syncURL(); render(true); }
  else if (action === 'status' || action === 'attention') { state.filter = control.dataset.status || 'attention'; state.query = ''; state.page = 0; navigate('devices'); }
  else if (action === 'clear') { state.query = ''; state.filter = 'all'; state.page = 0; render(); document.querySelector('#search').focus(); }
  else if (action === 'page') { state.page += Number(control.dataset.step); render(); document.querySelector('#search').focus(); announce(`Device results page ${state.page + 1}`); }
  else if (action === 'access') navigate('sessions');
  else if (action === 'signin') navigate('signin');
  else if (action === 'home') navigate('overview');
  else if (['ready', 'denied', 'expired'].includes(action)) { state.scenario = action === 'denied' ? 'forbidden' : action; render(true); }
  else if (action === 'enroll') { state.outcome = 'Pending example. Enrollment and connectivity are unconfirmed. Reconcile authoritative state before retrying.'; render(true); }
  else if (['revoke', 'permission'].includes(action)) confirmAction(action);
  else if (action === 'cancel') document.querySelector('#confirm').close();
  else if (action === 'submit-confirm') {
    const result = document.querySelector('#confirmation-outcome').value;
    state.outcome = `${result} — synthetic result only. ${result === 'Unknown' ? 'Response lost; reconcile authoritative state before retrying.' : result === 'Rejected' ? 'No change was accepted in this example.' : 'This illustrates confirmation, not an actual change. Active-session effects remain unresolved.'}`;
    document.querySelector('#confirm').close();
    if (confirmTrigger === 'revoke') { state.screen = 'detail'; render(); const panel = document.createElement('p'); panel.className = 'notice'; panel.setAttribute('role', 'status'); panel.textContent = state.outcome; document.querySelector('#content').append(panel); }
    else render();
    document.querySelector(`[data-action="${confirmTrigger}"]`)?.focus();
    announce(state.outcome);
  } else if (action === 'clear-audit') { Object.keys(state.audit).forEach(k => state.audit[k] = ''); render(); document.querySelector('#audit-actor').focus(); }
});
document.addEventListener('change', event => {
  if (event.target.id === 'scenario') { state.scenario = event.target.value; render(); document.querySelector('#scenario').focus(); announce(`Synthetic ${state.scenario} state`); }
  else if (event.target.id === 'filter') { state.filter = event.target.value; state.page = 0; render(); document.querySelector('#filter').focus(); }
  else if (event.target.id === 'request-result') { state.session = event.target.value; render(); document.querySelector('#request-result').focus(); }
});
document.addEventListener('input', event => {
  const target = event.target;
  if (target.id !== 'search' && !target.dataset.audit) return;
  const position = target.selectionStart;
  if (target.id === 'search') { state.query = target.value; state.page = 0; }
  else state.audit[target.dataset.audit] = target.value;
  render();
  const replacement = document.getElementById(target.id);
  replacement.focus();
  if (position !== null) replacement.setSelectionRange(position, position);
});
document.addEventListener('keydown', event => {
  if (!['ArrowLeft', 'ArrowRight'].includes(event.key) || event.altKey || event.ctrlKey || event.metaKey || event.shiftKey || event.target.closest('input,textarea,select,[contenteditable]') || document.querySelector('#confirm').open) return;
  event.preventDefault();
  document.querySelector(`#switcher [data-step="${event.key === 'ArrowLeft' ? '-1' : '1'}"]`)?.click();
});
document.querySelector('#confirm').addEventListener('close', () => { document.querySelector(`[data-action="${confirmTrigger}"]`)?.focus(); });
render();
