#!/usr/bin/env python3
"""Run the SDK lifecycle against a disposable core on loopback listeners."""
import os
import secrets
import signal
import socket
import subprocess
import sys
import tempfile
import time
import urllib.request
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

def free_port():
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        return listener.getsockname()[1]

def main():
    with tempfile.TemporaryDirectory(prefix="synapse-acceptance-") as tmp:
        binary = Path(tmp) / "synapse"
        subprocess.run(["go", "build", "-o", str(binary), "./cmd/synapse"], cwd=ROOT, check=True)
        ports = set()
        while len(ports) < 3:
            ports.add(free_port())
        http_port, mqtt_port, ws_port = ports
        env = dict(os.environ, SYNAPSE_AUTH_TOKEN=secrets.token_urlsafe(32),
                   SYNAPSE_ADMIN_TOKEN=secrets.token_urlsafe(32),
                   SYNAPSE_DB_PATH=str(Path(tmp) / "core.db"),
                   SYNAPSE_HTTP_PORT=f"127.0.0.1:{http_port}",
                   SYNAPSE_MQTT_PORT=f"127.0.0.1:{mqtt_port}",
                   SYNAPSE_WS_PORT=f"127.0.0.1:{ws_port}", SYNAPSE_COOKIE_SECURE="false",
                   SYNAPSE_ENABLE_ALERTS="false")
        with (Path(tmp) / "core.log").open("w+") as output:
            core = subprocess.Popen([str(binary)], cwd=tmp, env=env, stdout=output, stderr=output)
            try:
                base = f"http://127.0.0.1:{http_port}"
                deadline = time.monotonic() + 15
                while True:
                    if core.poll() is not None:
                        raise RuntimeError("isolated core exited before readiness")
                    try:
                        req = urllib.request.Request(base + "/api/v1/session", headers={"Authorization":"Bearer " + env["SYNAPSE_ADMIN_TOKEN"]})
                        with urllib.request.urlopen(req, timeout=1):
                            break
                    except OSError:
                        if time.monotonic() >= deadline:
                            raise TimeoutError("isolated core readiness timed out")
                        time.sleep(.1)
                env.update(PYTHONPATH=str(ROOT / "sdk/python/src"), SYNAPSE_MQTT_PORT=str(mqtt_port), SYNAPSE_HTTP_URL=base)
                subprocess.run([sys.executable, "sdk/python/tests/lifecycle.py"], cwd=ROOT, env=env, check=True, timeout=45)
                output.flush()
                output.seek(0)
                logs = output.read()
                if env["SYNAPSE_AUTH_TOKEN"] in logs or env["SYNAPSE_ADMIN_TOKEN"] in logs:
                    raise AssertionError("runtime logs exposed credentials")
                if "Alert Resolved" not in logs:
                    raise AssertionError("server monitor recovery not observed")
                print("Integrated acceptance passed with temporary storage, loopback listeners, monitor recovery and credential-free logs")
            finally:
                if core.poll() is None:
                    core.send_signal(signal.SIGINT)
                    try:
                        core.wait(timeout=5)
                    except subprocess.TimeoutExpired:
                        core.kill()
                        core.wait()

if __name__ == "__main__":
    main()
