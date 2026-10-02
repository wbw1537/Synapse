# Align discovery contract, validation and documentation

**ID:** SYN-101
**Status:** done
**Created:** 2026-10-02
**Priority:** high
**Depends-on:** -

## Goal

Make the layout/component contract unambiguous and reject payloads that register while silently losing their capabilities.

## Acceptance

- [x] One maintained contract states identity/property shape, version behavior, component/action IDs and transport paths.
- [x] HTTP and MQTT accept the same valid full snapshot and reject invalid/missing versions or references according to the documented contract.
- [x] Legacy widgets payloads receive a documented failure rather than appearing successfully registered without their capabilities.
- [x] Backend and UI fixtures agree on the supported shape; SDK/TOML docs clearly describe the mapping to that shape.

## Scope

Preserve flat identity/component properties and existing transport namespace as the baseline; explicitly decide supported version values and legacy rejection. Inspect APIVersion shadowing in ServicePayload. Do not implement the SDK or add a compatibility layer by accident.

## Verification

**Evidence:** 2026-10-03, working tree on `main` (baseline `36f9ac8`, uncommitted).

- `python3 scripts/agent.py verify --scope all`: passed all 11 workspace tests,
  document/task checks, Vue typecheck/build, `go vet ./...`, `go test ./...`, and
  the integrated Go binary build. Transport tests required sandbox-external local
  socket permissions; the sandbox failure was rerun successfully, not skipped.
- `TestDiscoveryTransports`: used temporary SQLite, a loopback HTTP server and
  real MQTT broker/clients. Both transports persisted the shared six-component
  snapshot, including api_version, layout and action IDs. Invalid tokens,
  versions, references, IDs, legacy arrays and nested shapes were rejected and
  existing stored state remained unchanged. MQTT topic mismatches were rejected;
  empty full snapshots removed capabilities. Declared action requests delivered
  the expected command topic/envelope to the test Axon.
- `npm --prefix web run test:protocol`: passed the named Vue rendering test using
  the same JSON fixture, showing all six components, CPU value/bar, status,
  log, action label and link. The fixture also compiles against the UI model.
- `git diff --check`: passed. Protocol, SDK/TOML mapping and the accepted decision
  agree on flat v1; documentation/catalog checks passed.

## Execution

Completed shared decoding/validation and MQTT topic identity checking. Removed
APIVersion shadowing so the version is persisted. Kept v1 routes/topics and flat
properties; legacy arrays and nested shapes now fail explicitly. SDK/TOML remains
proposed, with an explicit mapping rather than a conflicting wire definition.
Existing workspace/documentation changes predate this task and were preserved.
Independent agent review was unavailable under the active delegation policy;
solo diff review and behavioral checks were performed. Raw MQTT UI state and log
merge gaps remain SYN-102/SYN-103. No SDK, database migration or deployment added.
