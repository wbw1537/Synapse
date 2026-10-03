# Python Axon SDK

The reference SDK requires Python 3.11+ and paho-mqtt 2.x. The distribution name
is `synapse-axon`; the import name is `synapse_axon`.

## Install

From a repository checkout:

```sh
python -m pip install ./sdk/python
```

A PyPI release is optional. Local checkout or wheel installation works without
publishing. `pip install synapse-axon` requires an actual registry release; no
PyPI publication or package-name availability is claimed here.

Once the SDK source is available on GitHub, installation from Git is also possible:

```sh
python -m pip install 'git+https://github.com/wbw1537/Synapse.git#subdirectory=sdk/python'
```

## Configure and run

Create `axon.toml` using the [TOML reference](axon-toml.md). Set `SYNAPSE_AUTH_TOKEN`
to the core's Axon token. Optional connection variables are `SYNAPSE_MQTT_HOST`
and numeric `SYNAPSE_MQTT_PORT`. Authentication and listener exposure are explained
in the [access policy](access.md).

From the repository root, run the simulated memory example:

```sh
python examples/sdk_memory_axon.py
```

The [example configuration](https://github.com/wbw1537/Synapse/blob/main/examples/sdk_memory_axon.toml)
and [script](https://github.com/wbw1537/Synapse/blob/main/examples/sdk_memory_axon.py)
report a simulated gauge and logs, and handle a simulated recovery action. Adapt
the measurement and handler to your service before treating this as monitoring.

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

This snippet requires components `memory` and `logs`, and action `recover`,
as declared in the example configuration. Keep the context alive for the
monitoring lifetime; exiting it shuts down the client.

## Lifecycle and runtime updates

`Axon(path, token=None, host=None, port=None)` loads and validates TOML before
creating any MQTT client. Credentials come from arguments or environment, not
TOML. Invalid versions, references, action IDs, fields and initial values fail
early; unreferenced components warn. For schema fields and JSON mapping, use the
[TOML reference](axon-toml.md).

`axon.components[id].update(value)` validates and changes local state, then
publishes a full snapshot when connected. Gauge values must be finite and within
min/max; status values must have a configured mapping. Log strings append locally
and string arrays replace the bounded local history. The SDK publishes log arrays,
so heartbeat snapshots do not append the same events again. Invalid values and
oversized snapshots leave local state unchanged.

`@axon.on_action(id)` binds a declared callback. Bind every declared action before
`start()`. The network thread subscribes to its command topic before initial
registration, and resubscribes/reports after reconnect. Heartbeats run at TTL/2.
Callbacks use a bounded queue and a separate worker, so they do not block network
processing. Monitor expressions are evaluated by the core.

`stop()` stops heartbeats, publishes offline when connected and disconnects.
When disconnected, the core's TTL expires instead. Handlers should return promptly:
Python cannot forcefully terminate a running handler. Shutdown waits up to three
seconds per worker. Restart is refused while a previous handler remains active.
Context manager entry/exit wraps start/stop.

`start()` proves command subscription, not core persistence. Verify registration
with an operator-authenticated HTTP read as described in the [protocol](protocol.md).
`snapshot()` includes the credential for serialization; do not log it. Commands
are requests, with no exactly-once execution or execution-result acknowledgment.
