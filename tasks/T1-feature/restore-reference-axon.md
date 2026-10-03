# Restore a working layout/component reference Axon

**ID:** SYN-105
**Status:** done
**Created:** 2026-10-02
**Priority:** high
**Depends-on:** SYN-101

## Goal

Give developers one truthful runnable example that demonstrates discovery, widgets, monitor transitions and simulated recovery actions.

## Acceptance

- [x] The memory example uses the maintained layout/component wire shape.
- [x] Running it displays expected widgets and invokes a simulated declared action.
- [x] Normal/critical/recovered metric transitions exercise server monitor behavior.
- [x] Setup, credential configuration, expected results and known TTL/UI limitations are documented.

## Scope

Update examples/memory_axon.py and relevant README/sidecar reference. Retain simulated actions; do not mutate host memory/kernel settings.

## Verification

**Evidence:** isolated core/reference Axon run on temporary SQLite and loopback ports. HTTP persisted all six widgets; actual browser rendered them and bounded logs. 45 → 95 → 50 simulated readings exercised monitor recovery (core printed Alert Resolved). HTTP and browser declared-action requests reached simulated Axon callback. A disposable duplicate-client takeover disconnected the example; it automatically reconnected, resubscribed and received a subsequent command. Ctrl+C published offline and exited successfully. Missing token fails before connection. README documents environment settings and simulation; TTL/UI synchronization is now fixed by SYN-102.

## Execution

Migrated six components to flat v1 layout/components, environment connection/token settings, bounded array log snapshots, QoS1 publication with shutdown wait and reconnect resubscription. Actions remain simulated; host memory/kernel settings are not mutated. Reference-client and browser acceptance completed.
