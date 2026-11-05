"""Check local planning references and ensure the change stays documentation-only."""

from pathlib import Path
import re
import subprocess

BASE = "8d258fdacad19a896026dccaf8e67a38218dbdb7"
DOCUMENTS = (Path("docs/M02_EXECUTION.md"), Path("docs/ROADMAP.md"))

targets = 0
for document in DOCUMENTS:
    for target in re.findall(r"\]\(([^)]+)\)", document.read_text()):
        if "://" in target or target.startswith("#"):
            continue
        path = document.parent / target.split("#", 1)[0]
        if not path.is_file():
            raise SystemExit(f"Missing local target: {document}: {target}")
        targets += 1
print(f"PASS: {targets} local Markdown targets exist; external URLs/fragments not checked.")

changed = set(subprocess.check_output(
    ["git", "diff", BASE, "--name-only"], text=True).splitlines())
changed.update(subprocess.check_output(
    ["git", "ls-files", "--others", "--exclude-standard"], text=True).splitlines())
for name in sorted(changed):
    if name not in {str(path) for path in DOCUMENTS} and not name.startswith(
        "docs/evidence/m02-plan/"
    ):
        raise SystemExit(f"Outside issue documentation scope: {name}")
print(f"PASS: {len(changed)} changed/new files restricted to issue #35 documentation/evidence.")
subprocess.run(["git", "diff", "--check", BASE], check=True)
print(f"PASS: git diff --check {BASE}")
print("LIMIT: no application, browser, cross-build, physical ARM or security acceptance.")
