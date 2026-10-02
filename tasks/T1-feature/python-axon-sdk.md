# Deliver a reference Python Axon SDK

**ID:** SYN-104
**Status:** backlog
**Created:** 2026-10-02
**Priority:** medium
**Depends-on:** SYN-101,SYN-105

## Goal

Let an integrator load an axon.toml, report metrics and bind predefined actions without implementing MQTT lifecycle boilerplate.

## Acceptance

- [ ] Invalid config/layout references fail before network connection.
- [ ] A runnable demo registers, updates metrics and dispatches declared action callbacks.
- [ ] Heartbeat, reconnect/resubscribe and shutdown behavior pass observable lifecycle checks.
- [ ] The serialized payload matches the canonical contract and never logs credentials.

## Scope

Successor to SDK-004. Refine packaging/runtime support and component type rules after contract alignment and reference demo. Split additional tasks only where independently verifiable deliverables warrant it.

## Verification

Parser/type regression tests plus a disposable core/client lifecycle scenario. Existing SDK design docs are proposals, not proof of implementation.
