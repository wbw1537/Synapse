# Synchronize persisted service state with the open dashboard

**ID:** SYN-102
**Status:** done
**Created:** 2026-10-02
**Priority:** high
**Depends-on:** SYN-101

## Goal

Make the running dashboard reflect validated server state for MQTT/HTTP registration, TTL expiration and heartbeat recovery.

## Acceptance

- [x] HTTP registration appears in an already-open dashboard.
- [x] Stopped heartbeats become offline in both storage and the open dashboard within TTL plus monitor/delivery tolerance.
- [x] A new valid heartbeat restores the same service to online without refresh.
- [x] Rejected discovery payloads never overwrite valid dashboard state; reconnect reconciles missed state.

## Scope

Refine the server-to-UI state transport and reconciliation strategy before marking ready. Account for snapshots/log events, startup races and sensitive auth fields. Avoid independent backend/browser state semantics.

## Verification

**Evidence:** Race-enabled service/API tests passed: HTTP/MQTT accepted changes, no invalid notification, TTL persistence/notification, heartbeat recovery, coalescing and SSE reconnect initial snapshot. Playwright isolated core at port 18080 observed HTTP registration without refresh, rejected input absent, TTL icon offline, heartbeat icon online, and reload reconciliation; no console errors. Actual Vue fixture rendering and frontend typecheck/build passed.

## Execution

Use same-origin SSE full snapshots from persisted server state. Subscribe before initial read; bounded coalesced change notifications avoid slow clients blocking registration. Reconnect starts with a fresh snapshot. Serialize upsert/TTL mutations to protect log read-modify-write. No discovery token is exposed to the browser.

Completed implementation; subsequent stage work: implement stream and browser replacement, verify HTTP/MQTT rejection, TTL/recovery and reconnect scenarios.
