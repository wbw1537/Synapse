# Synchronize persisted service state with the open dashboard

**ID:** SYN-102
**Status:** backlog
**Created:** 2026-10-02
**Priority:** high
**Depends-on:** SYN-101

## Goal

Make the running dashboard reflect validated server state for MQTT/HTTP registration, TTL expiration and heartbeat recovery.

## Acceptance

- [ ] HTTP registration appears in an already-open dashboard.
- [ ] Stopped heartbeats become offline in both storage and the open dashboard within TTL plus monitor/delivery tolerance.
- [ ] A new valid heartbeat restores the same service to online without refresh.
- [ ] Rejected discovery payloads never overwrite valid dashboard state; reconnect reconciles missed state.

## Scope

Refine the server-to-UI state transport and reconciliation strategy before marking ready. Account for snapshots/log events, startup races and sensitive auth fields. Avoid independent backend/browser state semantics.

## Verification

Run a disposable Axon through both transports; keep the UI open during expiry/recovery and compare HTTP persisted snapshots with UI. Include rejected input and reconnect. No verification yet.
