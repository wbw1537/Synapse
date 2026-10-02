# Discovery contract alignment (SYN-101)

Accepted 2026-10-03 for the implementation of SYN-101.

Keep the existing flat identity/component representation and `v1` transport
namespace. Require exactly `api_version: "v1"`; the historical “v2” refactor name
does not introduce a second wire version. Store the version through the embedded
Service field, removing the shadowing payload field.

Reject legacy arrays and unknown/nested fields instead of providing implicit
conversion: silent capability loss is worse than a visible registration failure.
Validate layout references, component key/ID equality, supported component types
and service-wide action ID uniqueness before persistence. IDs are restricted to
ASCII letters, digits, underscores and hyphens so HTTP/MQTT paths are unambiguous.
Actions are declared in action groups; the speculative standalone action field
is unsupported. Empty snapshots and unreferenced components remain supported.

The [protocol](../protocol.md) is the single maintained wire reference. Proposed
TOML/SDK documents map to it rather than defining another wire format. No SDK,
legacy compatibility layer or existing-row migration is introduced. Raw MQTT UI
state and log-stream merge semantics remain separate tasks.
