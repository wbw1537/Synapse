#!/usr/bin/env python3
"""Build/run Synapse with disposable storage, without touching existing runtime data."""
import argparse
import os
from pathlib import Path
import secrets
import signal
import socket
import subprocess
import tempfile
import time

ROOT = Path(__file__).resolve().parents[1]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    for name, port in (("http", 8080), ("mqtt", 1883), ("ws", 8083)):
        parser.add_argument(f"--{name}-port", type=int, default=port)
    args = parser.parse_args()
    ports = [args.http_port, args.mqtt_port, args.ws_port]
    if len(set(ports)) != 3 or any(not 1 <= port <= 65535 for port in ports):
        parser.error("choose three distinct valid ports")
    for port in ports:
        with socket.socket() as sock:
            try:
                sock.bind(("0.0.0.0", port))
            except OSError:
                parser.error(f"port {port} unavailable; existing processes are left untouched")
    subprocess.run(["npm", "run", "build"], cwd=ROOT / "web", check=True)
    with tempfile.TemporaryDirectory(prefix="synapse-dev-") as tmp:
        binary = Path(tmp) / "synapse"
        subprocess.run(["go", "build", "-o", str(binary), "./cmd/synapse"], cwd=ROOT, check=True)
        env = dict(os.environ, SYNAPSE_DB_PATH=str(Path(tmp) / "synapse.db"),
                   SYNAPSE_HTTP_PORT=f"127.0.0.1:{args.http_port}",
                   SYNAPSE_MQTT_PORT=f"127.0.0.1:{args.mqtt_port}", SYNAPSE_WS_PORT=f"127.0.0.1:{args.ws_port}",
                   SYNAPSE_ENABLE_ALERTS="false")
        env["SYNAPSE_AUTH_TOKEN"] = os.environ.get("SYNAPSE_AUTH_TOKEN") or secrets.token_urlsafe(24)
        env["SYNAPSE_ADMIN_TOKEN"] = os.environ.get("SYNAPSE_ADMIN_TOKEN") or secrets.token_urlsafe(24)
        print(f"Dev UI: http://127.0.0.1:{args.http_port}; temporary DB: {tmp}", flush=True)
        print("Listeners bind loopback; set SYNAPSE_ADMIN_TOKEN before launching to log in", flush=True)
        process = subprocess.Popen([str(binary)], cwd=tmp, env=env, start_new_session=True)
        old_handlers = {}
        def stop(signum, frame):
            if process.poll() is None:
                process.send_signal(signal.SIGINT)
        for sig in (signal.SIGINT, signal.SIGTERM):
            old_handlers[sig] = signal.signal(sig, stop)
        try:
            while process.poll() is None:
                time.sleep(0.2)
            return process.returncode
        finally:
            if process.poll() is None:
                process.send_signal(signal.SIGINT)
                try:
                    process.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    process.kill()
                    process.wait()
            for sig, handler in old_handlers.items():
                signal.signal(sig, handler)


if __name__ == "__main__":
    raise SystemExit(main())
