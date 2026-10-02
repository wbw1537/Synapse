# Restore a working layout/component reference Axon

**ID:** SYN-105
**Status:** ready
**Created:** 2026-10-02
**Priority:** high
**Depends-on:** SYN-101

## Goal

Give developers one truthful runnable example that demonstrates discovery, widgets, monitor transitions and simulated recovery actions.

## Acceptance

- [ ] The memory example uses the maintained layout/component wire shape.
- [ ] Running it displays expected widgets and invokes a simulated declared action.
- [ ] Normal/critical/recovered metric transitions exercise server monitor behavior.
- [ ] Setup, credential configuration, expected results and known TTL/UI limitations are documented.

## Scope

Update examples/memory_axon.py and relevant README/sidecar reference. Retain simulated actions; do not mutate host memory/kernel settings.

## Verification

Run against temporary storage and a matching environment token; observe widgets, command callback and monitor transitions. Assert client shutdown/reconnect behavior appropriate to the example. No verification yet.
