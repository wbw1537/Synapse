# Reliable integration Alpha acceptance

Date: 2026-10-03. Source: working tree based on f6d0daa. The stage contract is
[SYN-108](../tasks/T1-feature/reliable-integration-alpha.md). No commit, push,
publication, deployment or mutation of existing runtime data was performed.

## Verified outcomes

| Outcome | Evidence |
| --- | --- |
| Reference Axon | Actual isolated core persisted and browser rendered six component types. 45/95/50 simulation exercised server monitor recovery. HTTP/browser actions reached simulated callbacks. Duplicate-client takeover forced reconnect/resubscribe; subsequent commands succeeded. Ctrl+C sent offline and exited. |
| Authoritative dashboard | HTTP and MQTT accepted state feed SSE; rejected registrations cannot mutate it. Service/API regressions include notification coalescing, TTL, recovery and reconnect snapshots. Final browser run observed HTTP registration without refresh, TTL offline, recovery online, rejected input absent, and missed updates reconciled after network restoration. |
| Bounded logs | SQLite regressions verify first string/array registration, append/replacement, clearing/preservation, invalid input rejection and retention. Thirty concurrent events were retained without loss. Actual browser array replacement removed old logs, retained newest entries and matched after reload. |
| Python SDK | Parser/type tests reject invalid config before client construction, retain state on invalid updates and validate flat mapping. Real core/client lifecycle verifies registration, heartbeats without log duplication, forced reconnect/resubscribe, action execution while heartbeats continue, monitor recovery, offline shutdown and restart. Built wheel installed into a separate temporary target and passed regression tests. |
| Access boundaries | Real MQTT broker tests reject anonymous/wrong/core-spoofing credentials, foreign/wildcard subscriptions and forged publications, while permitting discovery and commands. HTTP tests cover reads/actions/events, wrong/role credentials, foreign Origin, login/logout, cookie attributes, expiry and active-stream revocation. Final browser login, actual action, foreign-Origin rejection, mobile logout visibility and session revocation passed. |
| Reproducible checks | Clean export excluded dependency/build artifacts. npm ci, full workspace/Go/Vue checks, actual Vue fixture rendering, race tests and isolated SDK acceptance passed. CGO-disabled integrated binary built. SDK wheel built/installed; Compose configuration validated without deployment. |
| CI and dependencies | Workflow structure verifies all branch pushes, PRs and v* tags; image publication is gated on successful verification and version-tag push. Compatible dependency refresh reports zero npm audit findings. Remote GitHub checks were not triggered. |

## Repeatable commands

```sh
npm ci --prefix web
python -m pip install ./sdk/python
python3 scripts/agent.py verify --scope all
npm --prefix web run test:protocol
go test -race ./internal/api ./internal/broker ./internal/service -count=1
PYTHONPATH=sdk/python/src python3 -m unittest discover -s sdk/python/tests -v
python3 scripts/verify_integration.py
npm --prefix web audit --json
git diff --check
```

Workspace integrity: eleven tests passed. UI fixture rendering: one test passed.
SDK regression suites: two tests with invalid-config/type/retention cases passed.
Go tests include real HTTP/MQTT and temporary SQLite. Integration uses generated
secrets and checks that core logs do not contain them. Local runtime used Go from
go.mod, Node 22 and Python 3.14; CI selects Python 3.11. Sandbox-denied socket tests
were rerun with local-listener permission, never skipped.

Final browser ran Chromium against a freshly built temporary core on loopback.
Its assertions covered SDK recovery callback, HTTP live state, snapshot retention,
TTL UI, rejected state, recovery, reload logs and mobile sign-out. A separate
network-offline/state-change/network-online scenario verified reconnection and
logout revocation. Browser HTTP 401/400 entries from deliberate auth/rejection
checks are expected; there were no observed application exceptions.

## Review and limits

Solo review traced state ordering, read-modify-write concurrency, credentials,
ACLs, sessions, retry/shutdown paths and integration callers, including untracked
source. It fixed stale-stream callbacks after a new connection, mobile-hidden
logout and invalid oversized SDK updates retaining state. HTTP startup now waits
for broker subscription and publisher readiness; isolated lifecycle acceptance
was rerun successfully after this change. Independent agent review
was unavailable under the active delegation policy; the final human review remains
with the user.

This is a trusted-homelab Alpha. Axons share one credential and can impersonate
another ID if that credential is compromised. No native TLS, per-user roles,
per-Axon keys, exactly-once action execution, result acknowledgments, durable
history or legacy database migration is claimed. Use the [access policy](access.md)
and updated setup instructions before upgrading or admitting remote traffic.
