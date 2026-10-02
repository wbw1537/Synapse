---
name: dev-testing
description: Verify Synapse changes with focused Go/Vue checks and isolated-storage runtime scenarios for discovery, UI state, monitors, TTL and Axon command handling. Use for local validation or troubleshooting a regression; distinguish build checks from behavior evidence.
---

# Synapse Testing

Read [testing reference](../../docs/testing.md) and the task acceptance before
choosing checks. Use [implemented protocol](../../docs/protocol.md), not legacy
widget payloads or proposed nested meta/props, for current runtime integrations.

Start with `python3 scripts/agent.py doctor`. Run workspace checks for tooling/docs,
`verify --scope core` for backend and `verify --scope web` for UI; use `all` when
integrated embedding/build changes. Existing absence of test cases must be reported
rather than described as functional coverage.

Use `python3 scripts/dev.py` for disposable storage when a runtime is needed. Check
ports first, preserve existing services, and follow its networking limitation.
Provide credentials through the environment if an Axon needs the same token.
Do not print environment secrets or test on the existing root database.

For integration changes exercise affected transport and visible behavior:
registration, widgets/layout, monitor transition, command receipt at the Axon,
TTL expiration in the still-open UI, and recovery. Compare persisted state with
UI state where synchronization is relevant. Use an available browser tool if it
helps; do not assume another client's tool names or daemon paths exist.

Capture the acceptance scenario, observed outcome and code state. A successful
publish does not prove execution; compilation does not prove UI synchronization.
Use meaningful regression coverage for behavior fixes, not tests mirroring prose.
