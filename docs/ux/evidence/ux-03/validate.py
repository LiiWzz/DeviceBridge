"""Documentation and prescribed palette checks; no browser/conformance claim."""
from pathlib import Path
import re
import subprocess

BASE = "3eb9a80a58899ff5bd2131b1d23b264740e3a764"
files = [Path(p) for p in (
    "docs/ROADMAP.md", "docs/ux/UX-02.md", "docs/ux/UX-03.md"
)]
targets = 0
for source in files:
    for target in re.findall(r"\]\(([^)]+)\)", source.read_text()):
        if "://" in target or target.startswith("#"):
            continue
        assert (source.parent / target.split("#", 1)[0]).is_file(), (source, target)
        targets += 1
print(f"PASS: {targets} local Markdown file targets exist; URLs/fragments not checked.")


def luminance(color):
    channels = [int(color[i:i + 2], 16) / 255 for i in (1, 3, 5)]
    linear = [v / 12.92 if v <= 0.04045 else ((v + 0.055) / 1.055) ** 2.4
              for v in channels]
    return sum(a * b for a, b in zip(linear, (0.2126, 0.7152, 0.0722)))


rows = re.findall(
    r"^\| ([^|]+) \| `(?P<fg>#[0-9A-F]{6})` \| `(?P<bg>#[0-9A-F]{6})` \|",
    files[-1].read_text(), re.M)
assert len(rows) == 14, f"Expected 14 prescribed pairs, got {len(rows)}"
for label, foreground, background in rows:
    light, dark = sorted((luminance(foreground), luminance(background)), reverse=True)
    ratio = (light + 0.05) / (dark + 0.05)
    threshold = 3.0 if label.startswith(("Control boundary", "Focus ring")) else 4.5
    assert ratio >= threshold, (label, ratio, threshold)
    print(f"PASS: {label}: {foreground}/{background} = {ratio:.2f}:1 >= {threshold}:1")

changed = subprocess.check_output(
    ["git", "diff", BASE, "--name-only"], text=True).splitlines()
for path in changed:
    assert path in {str(p) for p in files} or path.startswith("docs/ux/evidence/ux-03/"), path
protected = ["agent", "server", "cli", "deploy", "packaging", "ci", "security",
             "tests", ".github", "AGENTS.md", "README.md", "docs/adr",
             "docs/M01_EXECUTION.md", "docs/ux/prototype"]
assert not subprocess.check_output(["git", "diff", BASE, "--", *protected])
print(f"PASS: application/prototype/contracts unchanged against {BASE}.")
for args in (["git", "diff", "--check"], ["git", "diff", "--cached", "--check"]):
    subprocess.run(args, check=True)
    print("PASS: " + " ".join(args))
print("LIMIT: static documentation and opaque sRGB pairs only; no browser, ARM or security acceptance.")
