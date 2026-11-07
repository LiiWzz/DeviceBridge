"""Check documentation scope only; no audit/storage/security acceptance."""
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[3]
BASE = "b25d382076facb1f7d4e68da79f392ed93a6caae"
ADR9 = "docs/adr/ADR-009-active-session-validity.md"
ADR10 = "docs/adr/ADR-010-authoritative-audit-consistency.md"
DOCS = [ADR9, ADR10, "GLOSSARY.md", "docs/M02_EXECUTION.md", "docs/ROADMAP.md"]
ALLOWED = set(DOCS + [f"docs/evidence/m02-d02b/{name}" for name in
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
print("PASS: exactly eight scoped document/evidence files; whitespace clean.")
protected = ["agent", "cli", "server", "deploy", "packaging", "ci", "security",
             "AGENTS.md", "README.md", "docs/ux", "docs/M01_EXECUTION.md"]
protected += [str(p.relative_to(ROOT)) for p in (ROOT / "docs/adr").glob("*.md")
              if str(p.relative_to(ROOT)) not in (ADR9, ADR10)]
assert not git("diff", BASE, "--", *protected)
print("PASS: runtime, migrations, prototype, security baseline, ADR-001–008 and M01 unchanged.")
old9 = git("show", f"{BASE}:{ADR9}")
new9 = (ROOT / ADR9).read_text()
old_body = old9.split("**Human decision:")[0]
new_body = new9.split("**Human decision:")[0]
assert re.sub(r"\*\*Status:.*?\*\*", "STATUS", old_body) == re.sub(
    r"\*\*Status:.*?\*\*", "STATUS", new_body)
assert "Status: Accepted (domain policy only)" in new9
assert "25fae2fb0fc8d5ff85c3caf79331f73e0dbab3ec" in new9
assert "ok next" in new9
print("PASS: ADR-009 reviewed body preserved; limited confirmation provenance recorded.")
new10 = (ROOT / ADR10).read_text()
assert "Status: Proposed" in new10 and "PENDING; intentionally unsigned" in new10
for file in DOCS[3:]:
    assert "adr/ADR-010-authoritative-audit-consistency.md" in (ROOT / file).read_text()
print("PASS: ADR-010 remains Proposed/unsigned; M02/roadmap handoffs present.")
