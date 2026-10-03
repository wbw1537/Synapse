# Deliver a reference Python Axon SDK

**ID:** SYN-104
**Status:** done
**Created:** 2026-10-02
**Priority:** medium
**Depends-on:** SYN-101,SYN-105

## Goal

Let an integrator load an axon.toml, report metrics and bind predefined actions without implementing MQTT lifecycle boilerplate.

## Acceptance

- [x] Invalid config/layout references fail before network connection.
- [x] A runnable demo registers, updates metrics and dispatches declared action callbacks.
- [x] Heartbeat, reconnect/resubscribe and shutdown behavior pass observable lifecycle checks.
- [x] The serialized payload matches the canonical contract and never logs credentials.

## Scope

Successor to SDK-004. Refine packaging/runtime support and component type rules after contract alignment and reference demo. Split additional tasks only where independently verifiable deliverables warrant it.

## Verification

**Evidence:** `PYTHONPATH=sdk/python/src python3 -m unittest discover -s sdk/python/tests -v` passed invalid schema/reference/field/action/type/retention/oversize-state checks. `python3 scripts/verify_integration.py` passed real registration, heartbeat/log idempotence, forced reconnect/resubscribe, callback while heartbeat continued, monitor recovery, offline shutdown/restart and credential-free runtime logs. Clean-export wheel build and separate wheel installation passed the same regressions. Actual browser SDK demo recovery button reached the handler and rendered its recovery log.

## Execution

Packaged Python reference SDK implemented with strict TOML/flat mapping, typed updates, idempotent log snapshots, heartbeat, bounded callback worker, reconnect subscription restoration and shutdown/restart.
