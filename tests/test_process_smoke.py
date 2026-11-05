"""Failure gates for the smoke harness, using real child processes."""
import contextlib
from pathlib import Path
import signal
import subprocess
import sys
import tempfile
import unittest

from process_smoke import command, process, stop, wait_for


class FailureGates(unittest.TestCase):
    def test_wrong_exit_code_fails(self):
        with self.assertRaisesRegex(RuntimeError, "expected 0"):
            command(Path(sys.executable), ["-c", "raise SystemExit(23)"], 0)

    def test_exit_before_readiness_fails(self):
        child = subprocess.Popen([sys.executable, "-c", "pass"])
        child.wait(timeout=5)
        with self.assertRaisesRegex(RuntimeError, "exited before"):
            wait_for(child, lambda: False, "startup")

    def test_unclean_shutdown_fails(self):
        with tempfile.TemporaryDirectory() as directory:
            log = Path(directory) / "process.log"
            log.write_text("started\n")
            child = subprocess.Popen([sys.executable, "-c", "import time; time.sleep(60)"])
            try:
                with self.assertRaisesRegex(RuntimeError, "unclean"):
                    stop(child, signal.SIGTERM, log, "stopped")
            finally:
                if child.poll() is None:
                    child.kill()
                child.wait(timeout=5)

    def test_exception_cleans_up_owned_process(self):
        with tempfile.TemporaryDirectory() as directory:
            fixture = Path(directory) / "idle"
            fixture.write_text("#!/bin/sh\nexec sleep 60\n")
            fixture.chmod(0o700)
            with contextlib.suppress(RuntimeError):
                with process(fixture, Path(directory) / "process.log") as child:
                    raise RuntimeError("simulated probe failure")
            self.assertIsNotNone(child.poll())


if __name__ == "__main__":
    unittest.main(verbosity=2)
