# Reliable integration Alpha stage

**ID:** SYN-108
**Status:** done
**Created:** 2026-10-03
**Priority:** high
**Depends-on:** -

## Goal

Advance Synapse to a reliable integration Alpha: a runnable Axon and packaged Python SDK, authoritative live dashboard state, consistent bounded logs, explicit enforced runtime access boundaries, and reproducible checks. User authorized autonomous design, implementation and acceptance; final changes remain available for review in the working tree.

## Acceptance

- [x] Reference Axon displays components, exercises server monitors and handles simulated declared actions.
- [x] HTTP/MQTT updates, rejected input, TTL expiry, recovery and reconnect agree between persisted state and the open UI.
- [x] Log snapshots/events are bounded and agree after reload; SDK heartbeats do not duplicate log events.
- [x] Python SDK validates TOML before connecting and supports typed updates, heartbeats, reconnect/resubscribe, callbacks and shutdown.
- [x] HTTP reads/actions and MQTT publication/subscription enforce a documented homelab access policy with positive and negative runtime checks.
- [x] Branch/PR CI checks workspace, Go, Vue and integrated builds; publication remains tag-only.
- [x] Updated current documentation and an acceptance report describe actual behavior, checks and remaining limitations.

## Scope

Complete SYN-102 through SYN-107 as independently verifiable increments. Keep single-binary deployment and declared-action execution in Axons. Choose a lightweight same-origin server state transport and explicit authentication policy; no deployment or mutation of existing services is needed for acceptance.

## Verification

**Evidence:** [acceptance report](../../docs/acceptance.md) records requirement-level command and runtime outcomes. All SYN-102 through SYN-107 and the dependency-refresh SYN-109 are complete. Clean export passed npm ci, integrated Go/Vue build and tests, race checks, SDK source/wheel regressions and real lifecycle acceptance. Final actual-browser assertions passed SDK recovery callback, HTTP registration, snapshot retention, TTL offline/recovery, rejected input, reload, mobile logout, disconnect reconciliation and session revocation. Compose config and static integrated binary build passed. Architecture/protocol/SDK/access/setup/operations docs match implemented behavior.

## Execution

Completed reliable integration Alpha from baseline f6d0daa. Temporary runtimes/storage and test credentials were used; existing deployments/data were untouched. Final source review and browser acceptance completed. Independent agent review unavailable under delegation policy; solo review and behavioral evidence recorded. Remaining product capabilities are future increments, not unfinished stage acceptance: per-Axon credentials/roles, native TLS, command result acknowledgments, durable history and legacy capability-row migration. Changes remain uncommitted for user review.
