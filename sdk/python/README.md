# Synapse Python Axon SDK

Python 3.11+; paho-mqtt 2.x. Install from the repository root:

```sh
python -m pip install ./sdk/python
```

A registry release is optional: installation from a checkout or a built wheel
works without PyPI. `synapse-axon` is the distribution name and `synapse_axon` is
the import name. A bare `pip install synapse-axon` requires an actual registry
release; this repository does not claim one.

Set `SYNAPSE_AUTH_TOKEN` to the core's Axon token. Optional connection variables:
`SYNAPSE_MQTT_HOST` and numeric `SYNAPSE_MQTT_PORT`. Run the simulated demo:

```sh
python examples/sdk_memory_axon.py
```

The SDK loads TOML and validates layout references, action IDs, component fields
and initial values before constructing any MQTT client. Gauge updates require
finite values within min/max; status values must have a configured mapping;
log strings append locally and string arrays replace the bounded local history.
Full snapshots, including heartbeat logs, are idempotent at the core.

```python
from synapse_axon import Axon

axon = Axon("axon.toml")

@axon.on_action("recover")
def recover():
    axon.components["memory"].update(45)
    axon.components["logs"].update("Recovered")

with axon:
    axon.components["memory"].update(75)
```

Bind every declared action before `start()`. The network thread subscribes before
initial registration and resubscribes/reports after reconnect. Heartbeats run at
TTL/2. Callback execution uses a bounded queue and separate worker, so callbacks
do not block network processing. Handlers should return promptly; Python cannot
forcefully terminate a running handler. Shutdown waits up to three seconds per
worker, publishes offline when connected, and disconnects; when disconnected the
core's TTL expires instead. Restart is refused while a previous handler remains
active. `start()` proves command subscription, not core persistence: verify via
an operator-authenticated HTTP read. `snapshot()` includes the credential for
serialization; do not log it. Commands are requests, not exactly-once execution.

From the repository root:

```sh
PYTHONPATH=sdk/python/src python -m unittest discover -s sdk/python/tests -v
python scripts/verify_integration.py
```

The integration check builds a disposable core, generates temporary secrets,
binds loopback listeners and verifies registration, metric updates, heartbeats,
log idempotence, forced reconnect/resubscribe, actual recovery callback, monitor
recovery, offline shutdown and credential-free logs. Build frontend assets first.
