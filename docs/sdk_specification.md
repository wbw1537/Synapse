# Reference Python Axon SDK

The implemented SDK lives in [sdk/python](../sdk/python/README.md), targets Python
3.11+ and uses paho-mqtt 2.x. Install with `python -m pip install ./sdk/python`.
The [TOML schema](axon_toml_spec.md) maps to the maintained
[flat discovery contract](protocol.md); other-language SDKs remain future work.

## Lifecycle and API

`Axon(path, token=None, host=None, port=None)` loads and validates TOML before
constructing a network client. Credentials come from arguments or environment,
not TOML. Missing versions, ghost references, duplicate actions, invalid field
shapes and initial values fail immediately; unreferenced components warn.

`axon.components[id].update(value)` validates and updates state, then publishes a
full snapshot when connected. Gauge values must be finite and inside min/max;
status values require a configured mapping. Logs accept string events or string
array replacement, retain max_items locally, and publish arrays so heartbeats
cannot duplicate events. Invalid typed values leave state unchanged.

`@axon.on_action(id)` binds a declared callback; all declared actions must be bound
before start. `start()` connects with service client ID and Axon credentials,
subscribes to its command topic, publishes registration and starts TTL/2 heartbeats.
Reconnect resubscribes and publishes current state. Actions run on a separate
worker with a bounded queue; monitor expressions remain server-side.

`stop()` stops heartbeats, attempts an acknowledged offline publication and
terminates network processing. A disconnected client relies on core TTL. Handlers
must return promptly because Python cannot cancel running callbacks. Context
manager entry/exit wraps start/stop. A successful MQTT subscribe/publication does
not establish persistence or exactly-once execution; use HTTP readback and actual
callback evidence.

See the SDK README for executable examples and verification commands.
