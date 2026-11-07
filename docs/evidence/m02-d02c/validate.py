"""Check documentation scope only; no clock/expiry/security acceptance."""
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[3]
BASE = "f3ec71ee387aab3e6bb3e9173446d83e96e5dbfa"
ADR10 = "docs/adr/ADR-010-authoritative-audit-consistency.md"
ADR11 = "docs/adr/ADR-011-trusted-time-validity.md"
DOCS = [ADR10, ADR11, "docs/M02_EXECUTION.md", "docs/ROADMAP.md"]
ALLOWED = set(DOCS + [f"docs/evidence/m02-d02c/{name}" for name in
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
            assert (path.parent / target).is_file(), (file, target)
            count += 1
print(f"PASS: {count} local Markdown file targets; external URLs/fragments unverified.")
changed = set(git("diff", "--name-only", BASE).splitlines())
# New docs can be globally ignored on this host; include scoped additions.
changed |= {file for file in ALLOWED if (ROOT / file).exists()
            and not git("ls-files", "--", file)}
assert changed == ALLOWED, sorted(changed ^ ALLOWED)
subprocess.run(["git", "diff", "--check", BASE], cwd=ROOT, check=True)
for file in ALLOWED:
    for line in (ROOT / file).read_text().splitlines():
        assert line == line.rstrip(), (file, "trailing whitespace")
print("PASS: exactly seven scoped document/evidence files; whitespace clean.")
protected = ["agent", "cli", "server", "deploy", "packaging", "ci", "security",
             "AGENTS.md", "README.md", "GLOSSARY.md", "docs/ux", "docs/M01_EXECUTION.md"]
protected += [str(p.relative_to(ROOT)) for p in (ROOT / "docs/adr").glob("*.md")
              if str(p.relative_to(ROOT)) not in (ADR10, ADR11)]
assert not git("diff", BASE, "--", *protected)
print("PASS: runtime, migrations, prototype, glossary, security baseline, ADR-001–009 and M01 unchanged.")
old10 = git("show", f"{BASE}:{ADR10}")
new10 = (ROOT / ADR10).read_text()
old_body = old10.split("**Human decision:")[0]
new_body = new10.split("**Human decision:")[0]
assert re.sub(r"\*\*Status:.*?\*\*", "STATUS", old_body) == re.sub(
    r"\*\*Status:.*?\*\*", "STATUS", new_body)
assert "Status: Accepted (architecture direction only)" in new10
assert "656ad1209d50d3c7e8fe7aec0c141276ca396927" in new10
assert "approve and next" in new10
print("PASS: ADR-010 reviewed body preserved; architecture approval provenance recorded.")
new11 = (ROOT / ADR11).read_text()
assert "Status: Proposed" in new11 and "PENDING; intentionally unsigned" in new11
for file in DOCS[2:]:
    assert "adr/ADR-011-trusted-time-validity.md" in (ROOT / file).read_text()
print("PASS: ADR-011 remains Proposed/unsigned; M02/roadmap handoffs present.")
