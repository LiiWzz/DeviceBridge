#!/usr/bin/env bash
# Build-host only; use the pinned Gitleaks build documented in security/README.md.
set -euo pipefail
root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"
unset GITLEAKS_CONFIG GITLEAKS_CONFIG_TOML GITLEAKS_CONFIG_JSON
out=$root/dist/ci/security
mkdir -p "$out"

# Prove detection and full redaction against a disposable synthetic fixture.
fixture=$(mktemp -d)
trap 'rm -rf "$fixture"' EXIT
python3 - "$fixture" <<'PY'
from pathlib import Path
import sys
marker = "DB007_SYNTHETIC_NOT_A_REAL_KEY"
Path(sys.argv[1], "fixture.txt").write_text(
    "-----BEGIN " + "RSA PRIVATE KEY-----\n" + marker * 4 + "\n-----END " + "RSA PRIVATE KEY-----\n"
)
PY
status=0
gitleaks dir "$fixture" --no-banner --redact=100 \
  --report-format=json --report-path="$fixture/report.json" > "$fixture/output.txt" 2>&1 || status=$?
python3 - "$fixture" "$status" <<'PY'
import json
from pathlib import Path
import sys
root = Path(sys.argv[1])
assert sys.argv[2] == "1", "scanner did not return its finding exit code"
report = json.loads((root / "report.json").read_text())
assert any(item["RuleID"] == "private-key" for item in report), "fixture not detected"
for name in ("output.txt", "report.json"):
    assert "DB007_SYNTHETIC_NOT_A_REAL_KEY" not in (root / name).read_text(), "redaction failed"
print("PASS: synthetic private key detected; console and JSON report redacted")
PY

# Scan committed history reachable from all locally available refs. CI fetches
# full history; ignored files and uncommitted workspace changes are not covered.
# Use built-in rules; no project allowlist or secret suppression is introduced.
gitleaks git --no-banner --redact=100 --log-opts=--all \
  --report-format=json --report-path="$out/gitleaks.json"
