"""Black-box M01 process checks. Python is build-host tooling only."""

import contextlib
import http.client
import json
import os
from pathlib import Path
import signal
import socket
import subprocess
import sys
import time


def require(condition, message):
    if not condition:
        raise RuntimeError(message)


def environment(**overrides):
    # Ignore operator control settings and proxy environment; never contact a proxy.
    env = {k: v for k, v in os.environ.items()
           if not k.startswith("DEVICEBRIDGE_CONTROL_")}
    env.update(overrides)
    return env


def command(binary, args, expected, env=None):
    result = subprocess.run([str(binary), *args], env=env or environment(),
                            capture_output=True, text=True, timeout=10)
    require(result.returncode == expected,
            f"{binary.name} {args}: exit {result.returncode}, expected {expected}\n"
            f"stdout: {result.stdout}\nstderr: {result.stderr}")
    print(f"PASS: {binary.name} {args} exit={expected}", flush=True)


@contextlib.contextmanager
def process(binary, log, env=None):
    with log.open("w") as output:
        child = subprocess.Popen([str(binary)], env=env or environment(),
                                 stdout=output, stderr=subprocess.STDOUT)
        try:
            yield child
        finally:
            if child.poll() is None:
                child.terminate()
                try:
                    child.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    child.kill()
                    child.wait(timeout=5)


def wait_for(child, predicate, message):
    deadline = time.monotonic() + 10
    while time.monotonic() < deadline:
        require(child.poll() is None, f"process exited before {message}")
        if predicate():
            return
        time.sleep(0.05)
    raise RuntimeError(f"timeout waiting for {message}")


def probe(port, method, path):
    connection = http.client.HTTPConnection("127.0.0.1", port, timeout=1)
    try:
        connection.request(method, path)
        response = connection.getresponse()
        return response.status, response.read(), dict(response.getheaders())
    finally:
        connection.close()


def ready(port):
    try:
        return probe(port, "GET", "/readyz")[:2] == (200, b'{"status":"ready"}\n')
    except OSError:
        return False


def stop(child, sig, log, expected_log):
    child.send_signal(sig)
    require(child.wait(timeout=5) == 0, "unclean signal shutdown")
    require(expected_log in log.read_text(), "missing clean shutdown log")
    print(f"PASS: {log.stem} {sig.name} exit=0 and stopped log", flush=True)


def main():
    # Run finally blocks on cancellation so owned children are reaped.
    signal.signal(signal.SIGTERM, lambda signum, _frame: sys.exit(128 + signum))
    binaries, logs = (Path(arg).resolve() for arg in sys.argv[1:])
    control = binaries / "control"
    relay = binaries / "devicebridge-relay"
    cli = binaries / "devicebridge"
    agent = binaries / "devicebridge-agent"
    command(control, [], 1, environment(DEVICEBRIDGE_CONTROL_ADDR="0.0.0.0:8080"))
    command(relay, ["--transport", "tls"], 2)
    command(cli, ["connect"], 2)
    command(cli, ["check", "--relay-url", "http://localhost"], 2)
    command(cli, ["check"], 0)
    command(relay, ["--check"], 0)
    command(agent, ["--check-config"], 0)
    command(agent, ["--shutdown-poll-ms", "0"], 2)

    for sig in (signal.SIGTERM, signal.SIGINT):
        # Ephemeral port selection has a small bind race; a conflict fails visibly.
        with socket.socket() as reservation:
            reservation.bind(("127.0.0.1", 0))
            port = reservation.getsockname()[1]
        log = logs / f"control-{sig.name}.txt"
        env = environment(DEVICEBRIDGE_CONTROL_ADDR=f"127.0.0.1:{port}")
        with process(control, log, env) as child:
            wait_for(child, lambda: ready(port), "control readiness")
            for path, status, body in (
                ("/healthz", 200, {"status": "ok"}),
                ("/readyz", 200, {"status": "ready"}),
                ("/debug/pprof/", 404, {"error": "not_found"}),
                ("/devices", 404, {"error": "not_found"}),
            ):
                code, payload, headers = probe(port, "GET", path)
                require(code == status and json.loads(payload) == body,
                        f"unexpected response for {path}")
                require(headers.get("Cache-Control") == "no-store", "probe caching")
                print(f"PASS: GET {path} status={status}", flush=True)
            for path in ("/healthz", "/readyz"):
                require(probe(port, "HEAD", path)[:2] == (200, b""), "HEAD body/status")
                code, _, headers = probe(port, "POST", path)
                require(code == 405 and headers.get("Allow") == "GET, HEAD", "method guard")
                print(f"PASS: HEAD/POST {path} status=200/405", flush=True)
            stop(child, sig, log, "control API stopped")
        with socket.socket() as check:
            check.settimeout(1)
            require(check.connect_ex(("127.0.0.1", port)) != 0, "listener survived shutdown")
        print("PASS: control listener closed", flush=True)

        for binary, started, stopped in (
            (relay, "relay idle", "relay stopped"),
            (agent, '"event":"agent_started"', '"event":"agent_stopped"'),
        ):
            log = logs / f"{binary.name}-{sig.name}.txt"
            with process(binary, log) as child:
                wait_for(child, lambda: started in log.read_text(), f"{binary.name} startup")
                stop(child, sig, log, stopped)
    print("PASS: M01 probes, negative cases and SIGINT/SIGTERM lifecycle", flush=True)


if __name__ == "__main__":
    main()
