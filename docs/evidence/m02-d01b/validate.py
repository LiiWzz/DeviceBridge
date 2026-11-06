"""Documentation consistency and scope only; no authorization/runtime validation."""
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[3]
BASE = "c3f4c99e508a6466f48bd3efd31c6eb1a1ec64dc"
DOCS = ["docs/adr/ADR-004-session-authority-binding.md", "docs/M02_EXECUTION.md", "docs/ROADMAP.md"]
ALLOWED = set(DOCS + [f"docs/evidence/m02-d01b/{name}" for name in ("validate.py", "validation.txt", "secrets.txt")])

def git(*args):
    return subprocess.check_output(["git", *args], cwd=ROOT, text=True)

count = 0
for file in DOCS:
    path = ROOT / file
    for target in re.findall(r"\]\(([^)]+)\)", path.read_text()):
        if "://" in target or target.startswith("#"):
            continue
        target = target.split("#")[0]
        if target:
            assert (path.parent / target).exists(), (file, target)
            count += 1
print(f"PASS: {count} local Markdown file targets; external URLs/fragments unverified.")
changed = set(git("diff", "--name-only", BASE).splitlines())
# The repository ignores docs; include only this issue's known new file locations.
changed |= {file for file in ALLOWED if (ROOT / file).exists() and git("ls-files", "--", file) == ""}
assert changed == ALLOWED, sorted(changed ^ ALLOWED)
subprocess.run(["git", "diff", "--check", BASE], cwd=ROOT, check=True)
print("PASS: exactly six scoped documentation/evidence files; whitespace clean.")
protected = ["agent", "cli", "server", "deploy", "packaging", "ci", "security", "AGENTS.md",
             "README.md", "GLOSSARY.md", "docs/ux", "docs/M01_EXECUTION.md",
             "docs/adr/ADR-001-foundation.md", "docs/adr/ADR-002-device-record-identity-binding.md",
             "docs/adr/ADR-003-browser-terminal-scope.md"]
assert not git("diff", BASE, "--", *protected)
print("PASS: runtime, prototype, security baseline, ADR-001/002/003 and M01 unchanged.")
proposal = (ROOT / DOCS[0]).read_text()
assert "Status: Proposed" in proposal and "PENDING; intentionally unsigned" in proposal
for file in DOCS[1:]:
    assert "adr/ADR-004-session-authority-binding.md" in (ROOT / file).read_text()
print("PASS: proposed unsigned status and navigable handoffs retained; no approval inferred.")
