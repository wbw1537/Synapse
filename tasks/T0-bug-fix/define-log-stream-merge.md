# Make log-stream updates consistent and bounded

**ID:** SYN-103
**Status:** done
**Created:** 2026-10-02
**Priority:** medium
**Depends-on:** SYN-101

## Goal

Preserve the same bounded log history across first registration, subsequent updates and page reloads.

## Acceptance

- [x] First registration with a single log string produces the documented initial history.
- [x] A string update appends once and an array snapshot replaces history according to one shared contract.
- [x] Retention respects max_items, and reload produces the same history as the live view.
- [x] Empty/invalid values have explicit behavior covered by focused regression cases.

## Scope

Backend log initialization/merge and frontend state application. Check current early return when an existing component map is nil and existing array replacement discrepancy. Avoid adding long-term log storage.

## Verification

**Evidence:** `go test -race ./internal/service ./internal/api -count=1` passed, including SQLite first-registration/event/snapshot/retention/invalid-value regression cases. Isolated core + memory Axon: browser live log history and reload matched HTTP persisted five-line snapshot (Playwright acceptance 2026-10-03). Empty string/null preserve, empty array clears; invalid types reject unchanged state.

## Execution

Implement server-owned normalization: string events append, string-array snapshots replace, null preserves history, empty arrays clear; reject other values and nonpositive explicit retention. Normalize first registration and changed component types.

Completed implementation; subsequent stage work: implement normalization and persisted regression scenarios, then verify browser consistency with authoritative state work.
