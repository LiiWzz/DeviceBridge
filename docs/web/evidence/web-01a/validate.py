"""Validate the discovery packet's local references and documentation scope."""
from pathlib import Path
import re
import subprocess

BASE = "8a6307a96ea0c46b2d91e089925c7dddb21d6377"
documents = [Path("docs/ROADMAP.md"), Path("docs/web/WEB-01a.md")]
targets = 0
for source in documents:
    for target in re.findall(r"\]\(([^)]+)\)", source.read_text()):
        if "://" in target or target.startswith("#"):
            continue
        assert (source.parent / target.split("#", 1)[0]).is_file(), (source, target)
        targets += 1
print(f"PASS: {targets} local Markdown targets exist; external URLs/fragments not checked.")

changed = subprocess.check_output(
    ["git", "diff", BASE, "--name-only"], text=True).splitlines()
for path in changed:
    assert path == "docs/ROADMAP.md" or path.startswith("docs/web/"), path
print(f"PASS: tracked diff against {BASE} is limited to roadmap and docs/web.")
untracked = subprocess.check_output(
    ["git", "ls-files", "--others", "--exclude-standard"], text=True).splitlines()
assert not untracked, untracked  # Stage the explicit documentation paths first.
subprocess.run(["git", "diff", "--check", BASE], check=True)
print("PASS: git diff --check " + BASE)
print("LIMIT: source/document inspection only; no application, browser, ARM or security acceptance.")
