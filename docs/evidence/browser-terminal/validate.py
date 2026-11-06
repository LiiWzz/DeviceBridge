"""Issue #45 document/scope checks; no browser or SSH tests."""

from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[3]
BASE = '1090dbdfd401bfb07d85d8fdb12dec08b4a1045c'
DOCS = ['README.md', 'GLOSSARY.md', 'docs/ROADMAP.md', 'docs/M02_EXECUTION.md',
        'docs/adr/ADR-003-browser-terminal-scope.md']
ALLOWED = set(DOCS + ['docs/evidence/browser-terminal/validate.py',
                      'docs/evidence/browser-terminal/validation.txt',
                      'docs/evidence/browser-terminal/secrets.txt'])


def git(*args):
    return subprocess.check_output(['git', *args], cwd=ROOT, text=True)


count = 0
for name in DOCS:
    doc = ROOT / name
    for target in re.findall(r'\]\(([^)]+)\)', doc.read_text()):
        if '://' in target or target.startswith('#'):
            continue
        assert (doc.parent / target.split('#')[0].split('?')[0]).is_file(), (name, target)
        count += 1
print(f'PASS: {count} local Markdown file targets; external URLs/fragments unverified.')

changed = set(git('diff', BASE, '--name-only').splitlines())
changed.update(git('ls-files', '--others', '--exclude-standard').splitlines())
assert changed == ALLOWED, (sorted(changed - ALLOWED), sorted(ALLOWED - changed))
subprocess.run(['git', 'diff', '--check', BASE], cwd=ROOT, check=True)
assert git('diff', BASE, '--', 'agent', 'server', 'cli', 'deploy', 'ci',
           'security', 'tests', 'packaging', 'docs/ux', 'AGENTS.md',
           'docs/adr/ADR-001-foundation.md',
           'docs/adr/ADR-002-device-record-identity-binding.md',
           'docs/M01_EXECUTION.md') == ''
print('PASS: exactly eight documentation/evidence files; whitespace clean.')
print('PASS: runtime/prototype/security baseline, ADR-001/002 and M01 contract unchanged.')

# Assert scope remains proposed; these are document consistency checks only.
adr = (ROOT / DOCS[-1]).read_text()
assert '**Status: Proposed' in adr and '**PENDING; intentionally unsigned.**' in adr
assert 'No browser terminal is implemented' in adr
assert 'not from a privileged web-shell agent' in adr
for name in ['README.md', 'docs/ROADMAP.md', 'docs/M02_EXECUTION.md']:
    assert 'ADR-003-browser-terminal-scope.md' in (ROOT / name).read_text()
print('PASS: proposal status and navigable handoffs retained; no approval/runtime acceptance inferred.')
