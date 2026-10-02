# Run reproducible checks on normal branches and PRs

**ID:** SYN-106
**Status:** ready
**Created:** 2026-10-02
**Priority:** medium
**Depends-on:** -

## Goal

Catch task/document drift and backend/frontend build regressions before version-tag publication.

## Acceptance

- [ ] PRs and relevant branch changes run workspace and Go/Vue verification.
- [ ] Go embedding checks use a built frontend asset rather than assuming web/dist exists.
- [ ] Ordinary PR checks do not publish images or deploy.
- [ ] Tag publication remains explicit and runs required checks first.

## Scope

CI workflow only; inspect current tag-only trigger and missing tracked frontend placeholder. Use scripts/agent.py for workspace checks and the existing build toolchain.

## Verification

Validate workflow trigger/job structure and run corresponding commands from a clean temporary export. Confirm publication conditions are restricted to intended events. No checks yet.
