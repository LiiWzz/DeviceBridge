"""Documentation and source-scope validation for issue #43; no runtime tests."""

from pathlib import Path
import re
import subprocess

BASE = 'e32cd6630b72452dd7c50bc0c448eef70e30d515'
ROOT = Path(__file__).resolve().parents[3]
DOCS = [
    'GLOSSARY.md',
    'docs/ROADMAP.md',
    'docs/M02_EXECUTION.md',
    'docs/adr/ADR-002-device-record-identity-binding.md',
    'docs/research/M02-D01a-sources.md',
]
ALLOWED = set(DOCS + [
    'docs/evidence/m02-d01a/validate.py',
    'docs/evidence/m02-d01a/validation.txt',
    'docs/evidence/m02-d01a/secrets.txt',
])


def git(*args):
    return subprocess.check_output(['git', *args], cwd=ROOT, text=True)


count = 0
for name in DOCS:
    doc = ROOT / name
    for target in re.findall(r'\]\(([^)]+)\)', doc.read_text()):
        if '://' in target or target.startswith('#'):
            continue
        path = target.split('#', 1)[0].split('?', 1)[0]
        assert (doc.parent / path).is_file(), (name, target)
        count += 1
print(f'PASS: {count} local Markdown file targets; external URLs/fragments unverified.')

# Staging includes explicitly approved documentation otherwise ignored by git.
changed = set(git('diff', BASE, '--name-only').splitlines())
changed.update(git('ls-files', '--others', '--exclude-standard').splitlines())
assert changed == ALLOWED, (sorted(changed - ALLOWED), sorted(ALLOWED - changed))
subprocess.run(['git', 'diff', '--check', BASE], cwd=ROOT, check=True)
assert git('diff', BASE, '--', 'agent', 'server', 'cli', 'deploy', 'ci',
           'security', 'tests', 'packaging', 'docs/ux', 'docs/adr/ADR-001-foundation.md',
           'docs/M01_EXECUTION.md', 'AGENTS.md', 'README.md') == ''
print('PASS: exactly eight scoped documentation/evidence files; whitespace clean.')
print('PASS: application/prototype code, security baseline, ADR-001 and M01 contract unchanged.')

# Preserve source identity for the research claims without executing applications.
for name in [
    'server/control/internal/control/server.go',
    'deploy/migrations/0001_baseline.up.sql',
    'agent/src/config.rs',
    'agent/src/main.rs',
    'cli/cmd/devicebridge/main.go',
    'server/relay/internal/transport/transport.go',
    'docs/ux/prototype/prototype.js',
]:
    baseline = subprocess.check_output(['git', 'show', f'{BASE}:{name}'], cwd=ROOT)
    assert (ROOT / name).read_bytes() == baseline, name
print('PASS: seven inspected owning source files match the research baseline byte-for-byte.')
