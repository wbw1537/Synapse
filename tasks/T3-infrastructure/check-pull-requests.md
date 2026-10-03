# Run reproducible checks on normal branches and PRs

**ID:** SYN-106
**Status:** done
**Created:** 2026-10-02
**Priority:** medium
**Depends-on:** -

## Goal

Catch task/document drift and backend/frontend build regressions before version-tag publication.

## Acceptance

- [x] PRs and relevant branch changes run workspace and Go/Vue verification.
- [x] Go embedding checks use a built frontend asset rather than assuming web/dist exists.
- [x] Ordinary PR checks do not publish images or deploy.
- [x] Tag publication remains explicit and runs required checks first.

## Scope

CI workflow only; inspect current tag-only trigger and missing tracked frontend placeholder. Use scripts/agent.py for workspace checks and the existing build toolchain.

## Verification

**Evidence:** YAML structure assertions passed branch/PR/tag triggers, verification dependency and exact tag-only publication condition. Clean export with no web/dist or node_modules passed npm ci, scripts/agent.py verify --scope all, Vue fixture rendering, Go race tests, SDK regressions and scripts/verify_integration.py. CGO-disabled integrated binary built. Ordinary verification jobs have read-only contents permissions and no publication steps. Remote workflow was not triggered; local equivalent commands and workflow configuration were verified.

## Execution

Workflow runs verification for all branch pushes, PRs and version tags; image publication requires a version-tag push and successful verification. Build frontend before Go embedding; SDK install/regressions/lifecycle and race checks are included.
