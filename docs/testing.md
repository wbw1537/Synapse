# Environment and verification

Prerequisites: Python 3.10+ for workspace tools, Go version from `go.mod`, and
Node/npm for `web/`. Workspace tools use the Python standard library.

```sh
python3 scripts/agent.py doctor
python3 scripts/agent.py check
python3 scripts/agent.py verify --scope workspace
python3 scripts/agent.py verify --scope core
python3 scripts/agent.py verify --scope web
python3 scripts/agent.py verify --scope all
```

`check` validates task schema, dependencies/cycles, the generated index, document
catalog and maintained local Markdown links, and skill entrypoints/references.
`verify workspace` also runs the tooling integrity tests.
`verify core` builds frontend assets if absent, then runs `go vet ./...` and
`go test ./...`; `verify web` runs the frontend typecheck/build. `all` runs workspace checks, both application checks,
and builds the integrated binary in a temporary directory. Failed dependencies
or checks remain failures. Discovery transport regression tests use a loopback MQTT broker, HTTP server and
temporary SQLite (`go test ./internal/api -run TestDiscoveryTransports -count=1`).
`npm --prefix web run test:protocol` renders the actual Vue card with the same
six-component fixture; `npm --prefix web run build` checks its UI model types.
Do not call compilation alone behavioral coverage.

## Isolated development runtime

```sh
python3 scripts/dev.py
```

Builds the UI and binary, uses a temporary working directory/database, disables
SMTP, and binds HTTP to loopback. Defaults are the application's
8080/1883/8083; use `--http-port`, `--mqtt-port`, and `--ws-port` to avoid an
existing runtime. The browser uses same-origin SSE, so changing the MQTT/WS ports does not
affect live UI connectivity.
All three listeners bind loopback in the launcher. The core accepts complete
host:port listener addresses and connects to the selected broker address.

Supply distinct `SYNAPSE_AUTH_TOKEN` and `SYNAPSE_ADMIN_TOKEN` in the invoking
environment to integrate an Axon and log in. Otherwise the launcher generates
unprinted ephemeral secrets. The browser uses the operator token, the Axon uses
its token both for MQTT authentication and discovery. See [access policy](../website/docs/access.md).

No `.env` is read automatically by Go; the launcher inherits environment and
supplies temporary DB/token/ports. It never prints credentials. Ctrl+C terminates
the temporary process and removes its working data. It never deletes the existing
`synapse.db` or restarts existing containers.

## Behavior checks

Use disposable Axon IDs and the implemented payload in [public protocol](../website/docs/protocol.md).
Test discovery through each affected transport, section/widget rendering,
monitor state transitions, declared action publication and actual Axon handling,
then stop heartbeats and inspect both persisted state and the still-open UI.
Record expected and observed outcomes, not only HTTP success codes.
Browser work may use an available Playwright tool; do not assume a tool or daemon
from another client's setup exists. Evidence must omit credentials.

## SDK and integrated acceptance

Install the SDK using the [public SDK guide](../website/docs/python-sdk.md#install).
From the repository root, run these development checks:

```sh
PYTHONPATH=sdk/python/src python -m unittest discover -s sdk/python/tests -v
python scripts/verify_integration.py
go test -race ./internal/api ./internal/broker ./internal/service
```

The integration script builds a disposable core, selects loopback ports and
generates temporary secrets. It checks actual SDK registration, heartbeat/log
idempotence, forced reconnect, command callback, monitor recovery and shutdown.
It requires built frontend assets. CI runs these checks after the full build;
HTTP/MQTT tests require permission to listen on local temporary ports.

## Task viewer

`./scripts/tasks.sh` serves the task viewer at `http://127.0.0.1:6060`.
It serves only task data and maintained Markdown documents, not the repo root,
secret files, databases, or source files. The viewer works without CDN libraries.
