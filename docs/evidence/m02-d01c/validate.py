"""Check this proposal's documents and scope; no transport/security validation."""
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[3]
BASE = "168578d8e9e0ff5ecada58b282d69fdf32852fc7"
DOCS = ["docs/adr/ADR-007-operator-trust-and-authority.md",
        "docs/M02_EXECUTION.md", "docs/ROADMAP.md", "GLOSSARY.md"]
ALLOWED = set(DOCS + [f"docs/evidence/m02-d01c/{name}" for name in
                     ("validate.py", "validation.txt", "secrets.txt")])


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
# docs is globally ignored on this host; explicitly include new scoped files.
changed |= {file for file in ALLOWED if (ROOT / file).exists()
            and not git("ls-files", "--", file)}
assert changed == ALLOWED, sorted(changed ^ ALLOWED)
subprocess.run(["git", "diff", "--check", BASE], cwd=ROOT, check=True)
for file in ALLOWED:
    for line in (ROOT / file).read_text().splitlines():
        assert line == line.rstrip(), (file, "trailing whitespace")
print("PASS: exactly seven scoped document/evidence files; whitespace clean.")
protected = ["agent", "cli", "server", "deploy", "packaging", "ci", "security",
             "AGENTS.md", "README.md", "docs/ux",
             "docs/M01_EXECUTION.md"]
protected += [str(p.relative_to(ROOT)) for p in (ROOT / "docs/adr").glob("*.md")
              if p.name != "ADR-007-operator-trust-and-authority.md"]
assert not git("diff", BASE, "--", *protected)
print("PASS: runtime, prototype, security baseline, prior ADRs and M01 unchanged.")
proposal = (ROOT / DOCS[0]).read_text()
assert "Status: Proposed" in proposal and "PENDING; intentionally unsigned" in proposal
for file in DOCS[1:3]:
    assert "adr/ADR-007-operator-trust-and-authority.md" in (ROOT / file).read_text()
print("PASS: proposed unsigned status and linked M02/roadmap handoffs retained.")
