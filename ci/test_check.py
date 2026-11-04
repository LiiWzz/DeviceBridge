"""Exercise CI failure gates in disposable repositories; no application mutation."""

import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import unittest


SCRIPT = Path(__file__).with_name("check.sh")


class FailureGates(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="devicebridge-ci-")
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        (self.root / "ci").mkdir()
        shutil.copyfile(SCRIPT, self.root / "ci/check.sh")
        subprocess.run(["git", "init", "-q", str(self.root)], check=True)
        subprocess.run(
            ["git", "-c", "user.name=CI Fixture", "-c", "user.email=ci@example.invalid",
             "commit", "-q", "--allow-empty", "-m", "fixture"],
            cwd=self.root, check=True,
        )

    def invoke(self, *args, env=None):
        return subprocess.run(
            ["bash", str(self.root / "ci/check.sh"), *args],
            cwd=self.root, env=env, capture_output=True, text=True, timeout=30,
        )

    def test_untrusted_selectors_are_rejected_before_output(self):
        for args in [(), ("go", "../../"), ("arm", "amd64"),
                     ("rust", "extra"), ("workflow", "extra")]:
            with self.subTest(args=args):
                result = self.invoke(*args)
                self.assertEqual(result.returncode, 2, result.stdout + result.stderr)
                self.assertFalse((self.root / "dist").exists())

    def test_command_failure_survives_tee_and_is_recorded(self):
        tools = self.root / "tools"
        tools.mkdir()
        fake_go = tools / "go"
        fake_go.write_text("#!/bin/sh\necho deliberate-command-failure >&2\nexit 23\n")
        fake_go.chmod(0o700)
        (self.root / "cli").mkdir()
        artifact = self.root / "dist/ci/go-cli"
        artifact.mkdir(parents=True)
        (artifact / "SHA256SUMS").write_text("stale-success-manifest\n")
        result = self.invoke("go", "cli", env={**os.environ, "PATH": f"{tools}:{os.environ['PATH']}"})
        self.assertEqual(result.returncode, 23, result.stdout + result.stderr)
        log = (artifact / "commands.txt").read_text()
        self.assertIn("deliberate-command-failure", log)
        self.assertIn("[exit 23]", log)
        self.assertTrue((artifact / "SOURCE.txt").is_file())
        self.assertFalse((artifact / "SHA256SUMS").exists())

    def test_unformatted_go_fails_instead_of_only_printing(self):
        module = self.root / "cli"
        module.mkdir()
        (module / "main.go").write_text('package main\nfunc main(){println("fixture")}\n')
        result = self.invoke("go", "cli")
        self.assertEqual(result.returncode, 1, result.stdout + result.stderr)
        self.assertIn("main.go", result.stdout)
        log = (self.root / "dist/ci/go-cli/commands.txt").read_text()
        self.assertIn("[exit 1]", log)
        self.assertNotIn("$ go vet", log)

    def test_gofmt_parse_error_stops_before_vet(self):
        module = self.root / "cli"
        module.mkdir()
        (module / "main.go").write_text("package main\nfunc main( {\n")
        result = self.invoke("go", "cli")
        self.assertEqual(result.returncode, 2, result.stdout + result.stderr)
        log = (self.root / "dist/ci/go-cli/commands.txt").read_text()
        self.assertIn("[exit 2]", log)
        self.assertNotIn("$ go vet", log)


if __name__ == "__main__":
    unittest.main(verbosity=2)
