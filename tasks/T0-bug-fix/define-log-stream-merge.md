# Make log-stream updates consistent and bounded

**ID:** SYN-103
**Status:** ready
**Created:** 2026-10-02
**Priority:** medium
**Depends-on:** SYN-101

## Goal

Preserve the same bounded log history across first registration, subsequent updates and page reloads.

## Acceptance

- [ ] First registration with a single log string produces the documented initial history.
- [ ] A string update appends once and an array snapshot replaces history according to one shared contract.
- [ ] Retention respects max_items, and reload produces the same history as the live view.
- [ ] Empty/invalid values have explicit behavior covered by focused regression cases.

## Scope

Backend log initialization/merge and frontend state application. Check current early return when an existing component map is nil and existing array replacement discrepancy. Avoid adding long-term log storage.

## Verification

Use two successive payloads per input form, query persisted services, reload the browser and compare histories; test retention boundaries. No verification yet.
