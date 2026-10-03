---
name: implement-axon
description: Create Synapse reporting clients or interactive sidecars using the implemented discovery and command contract. Use for integrating a service, metric or recovery operation with Synapse; the Python reference SDK is available.
---

# Implement an Axon

Read the public [protocol](../../website/docs/protocol.md),
[access policy](../../website/docs/access.md) and, for Python, the
[SDK guide](../../website/docs/python-sdk.md). These are the integration authority;
this skill records how to develop and verify a client, not a second specification.

Identify measurements, declared actions and the intended service identity before
implementation. Follow the public schema and authentication contract, and use
the reference SDK when it fits. Keep secrets in environment/runtime configuration.
Axons execute predefined local handlers; never execute arbitrary received code.

Use disposable runtime data to verify registration, component rendering, actual
callback execution, monitor normal/triggered/recovered transitions, reconnect,
shutdown and TTL behavior. A successful core response alone does not prove a
handler executed. The examples use simulated metrics and actions; adapt them to
the actual service rather than treating simulation as monitoring.

Update public instructions only in their owner page and record development
acceptance in the task. Use [testing](../../docs/testing.md) for isolated checks.
