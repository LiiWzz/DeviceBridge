"""Check this documentation scope; not runtime/security acceptance."""
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[3]
BASE = "697045e113b5a8b5ff834181e023c37eef6bc2a7"
ADR8 = "docs/adr/ADR-008-device-bootstrap-trust.md"
ADR9 = "docs/adr/ADR-009-active-session-validity.md"
DOCS = [ADR8, ADR9, "GLOSSARY.md", "docs/M02_EXECUTION.md", "docs/ROADMAP.md"]
ALLOWED = set(DOCS + [f"docs/evidence/m02-d02a/{name}" for name in
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
              if str(p.relative_to(ROOT)) not in (ADR8, ADR9)]
assert not git("diff", BASE, "--", *protected)
print("PASS: runtime, prototype, security baseline, ADR-001–007 and M01 unchanged.")
old8 = git("show", f"{BASE}:{ADR8}")
new8 = (ROOT / ADR8).read_text()
# Approval changes status and decision record only, retaining reviewed constraints.
old_body = old8.split("**Human decision:")[0]
new_body = new8.split("**Human decision:")[0]
assert re.sub(r"\*\*Status:.*?\*\*", "STATUS", old_body) == re.sub(
    r"\*\*Status:.*?\*\*", "STATUS", new_body)
assert "Status: Accepted (domain principles only)" in new8
assert "44fd7fa18527ea3dc3583e1e4402e2cafbfc86a8" in new8
assert "confirm and next" in new8
print("PASS: ADR-008 reviewed body preserved; limited confirmation provenance recorded.")
new9 = (ROOT / ADR9).read_text()
assert "Status: Proposed" in new9 and "PENDING; intentionally unsigned" in new9
for file in DOCS[3:]:
    assert "adr/ADR-009-active-session-validity.md" in (ROOT / file).read_text()
print("PASS: ADR-009 remains Proposed/unsigned; M02/roadmap handoffs present.")
