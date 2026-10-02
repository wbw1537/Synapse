---
name: implement-axon
description: Create Synapse reporting clients or interactive sidecars using the implemented discovery and command contract. Use for integrating a service, metric or recovery operation with Synapse; the Python SDK remains planned.
---

# Implement an Axon

Read [implemented protocol](../../docs/protocol.md) first. The layout/component
schema is current, while the old `widgets`/`actions` arrays and proposed nested
`meta`/`props` are not supported shapes. `api_version` and topics still use `v1`
despite the refactor being called v2. The core requires exactly `api_version: "v1"`
and rejects legacy arrays, nested fields, mismatched IDs and ghost references.

Construct a full discovery snapshot with top-level identity/token, status, positive
TTL, section layout and a component map. Flatten component properties as the
current model expects. Use IDs referenced by section children consistently.

Reporting clients publish to `synapse/v1/discovery/{id}` or POST
`/api/v1/discovery`. Resend before TTL expires, typically TTL/2. Do not expose
secrets in logs or hardcode real credentials into committed examples.

Interactive clients also subscribe to `synapse/v1/command/{id}` and bind declared
`action_group.items[].action_id` values to predefined local handlers. Reject
unknown IDs and untrusted arguments; never execute arbitrary received shell code.
Handle reconnect/subscription restoration and cleanup proportionally to the client.
A successful core response only establishes dispatch, not completed execution.

Monitor expressions run in the core against `value`, with notifications on state
transitions. Test a normal/triggered/recovered sequence. Log-stream merging differs
between backend and UI; define whether the client sends a snapshot or appended event
and verify behavior rather than assuming lists and strings are interchangeable.

Verify registration, rendering, commands where supported, stop/TTL and recovery
using disposable runtime data. `examples/memory_axon.py` is a legacy example until
its tracked migration is complete. Do not use it as a proven new-protocol template.
