// Reproducible exploration trace, not a prototype test suite or browser runner.
import { readFileSync } from 'node:fs';
import { Script, createContext } from 'node:vm';
const html = readFileSync(new URL('../../ux/prototype/terminal-logic.html', import.meta.url), 'utf8');
const scripts = [...html.matchAll(/<script id="([^"]+)">([\s\S]*?)<\/script>/g)];
for (const [, id, source] of scripts) {
  new Script(source, { filename: id });
  console.log(`Syntax compiled: ${id}`);
}
const context = createContext({});
new Script(scripts[0][2] + '\nglobalThis.demo = TerminalDemo;').runInContext(context);
// Reuse the demo's exact guided scenarios without running its DOM shell.
const scenarios = scripts[1][2].split('const scenarios = ')[1].split('\nlet selected = ')[0];
new Script('globalThis.scenarios = ' + scenarios).runInContext(context);
for (const scenario of context.scenarios) {
  let state = context.demo.initial();
  console.log(`\nWalkthrough: ${scenario.name}`);
  for (const [label, action] of scenario.steps) {
    state = context.demo.reduce(state, action);
    console.log(JSON.stringify({ action: label, phase: state.phase, authority: state.authority,
      session: state.session, starts: state.starts, draft: state.draft, staged: state.staged,
      outputEntries: state.output.length, notice: state.notice }));
  }
}
for (const outcome of ['denied', 'expired', 'unavailable']) {
  let state = context.demo.initial();
  const types = outcome === 'expired' ? ['connect', 'attach', outcome, 'send', 'connect'] : ['connect', outcome, 'send'];
  for (const type of types) state = context.demo.reduce(state, { type });
  console.log(`\nFree-play outcome ${outcome}: ${JSON.stringify(state)}`);
}
let state = context.demo.initial();
for (const type of ['connect', 'attach']) state = context.demo.reduce(state, {type});
for (const action of [ {type:'paste-edit',text:'<img src=x onerror=example>\necho inert'}, {type:'stage'},
  {type:'paste-edit',text:'edited sample'}, {type:'confirm'}, {type:'stage'}, {type:'confirm'}, {type:'confirm'} ]) {
  state = context.demo.reduce(state, action);
  console.log(`Paste exploration: ${JSON.stringify(state)}`);
}
for (let i = 0; i < 50; i++) {
  state = context.demo.reduce(state, {type:'draft', text:'x'.repeat(2048)});
  state = context.demo.reduce(state, {type:'send'});
}
state = context.demo.reduce(state, {type:'resize',columns:999,rows:0});
console.log(`\nBounded exploration: output entries=${state.output.length}; maximum entry characters=${Math.max(...state.output.map(x => x.length))}; dimensions=${state.columns}x${state.rows}; ${state.notice}`);
console.log('LIMIT: pure model and script syntax only; no DOM, browser, SSH, PTY or security acceptance.');
