#!/usr/bin/env python3
"""Loopback, allowlisted, read-only task viewer; never serves the repository root."""
import argparse
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[1]


class Handler(BaseHTTPRequestHandler):
    def do_HEAD(self):
        self.serve(False)

    def do_GET(self):
        self.serve(True)

    def serve(self, body):
        route = unquote(urlsplit(self.path).path)
        relative = "scripts/tasks-viewer.html" if route in ("/", "/scripts/tasks-viewer.html") else route.lstrip("/")
        file = (ROOT / relative).resolve()
        allowed = file.is_relative_to(ROOT / "tasks") and file.suffix in (".md", ".json")
        allowed |= file.is_relative_to(ROOT / "docs") and file.suffix == ".md"
        allowed |= file == ROOT / "scripts/tasks-viewer.html"
        if not file.is_relative_to(ROOT) or not allowed or not file.is_file():
            self.send_error(404)
            return
        data = file.read_bytes()
        mime = "text/html" if file.suffix == ".html" else "application/json" if file.suffix == ".json" else "text/plain"
        self.send_response(200)
        self.send_header("Content-Type", f"{mime}; charset=utf-8")
        self.send_header("Content-Length", str(len(data)))
        self.send_header("Cache-Control", "no-store")
        self.send_header("X-Content-Type-Options", "nosniff")
        self.end_headers()
        if body:
            self.wfile.write(data)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--port", type=int, default=6060)
    args = parser.parse_args()
    with ThreadingHTTPServer(("127.0.0.1", args.port), Handler) as server:
        print(f"Task viewer: http://127.0.0.1:{args.port}", flush=True)
        try:
            server.serve_forever()
        except KeyboardInterrupt:
            pass


if __name__ == "__main__":
    main()
