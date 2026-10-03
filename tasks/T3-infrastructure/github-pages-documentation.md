# Publishable GitHub Pages documentation

**ID:** SYN-110
**Status:** done
**Created:** 2026-10-03
**Priority:** medium
**Depends-on:** -

## Goal

Provide a searchable documentation site from maintained repository Markdown, ready for GitHub Pages without duplicating SDK or product documentation.

## Acceptance

- [x] Home, integration/protocol, Python SDK, operations and development documents are navigable and searchable.
- [x] Current Markdown is reused; historical docs/task/code references resolve to GitHub source rather than broken site URLs.
- [x] Only selected public Markdown is staged, with no environment files, databases or repository source dump.
- [x] Strict local build and actual browser navigation/search/subpath behavior pass.
- [x] PR checks build docs without deploying; main pushes/manual main runs can deploy a verified Pages artifact with scoped permissions.
- [x] Setup explains optional SDK distribution, local preview and the one-time Pages setting; remote publication status is accurately reported.

## Scope

Use MkDocs/Material and a standard-library staging/link-rewrite script over docs/index.json plus root README. Add dedicated Pages workflow and documentation build instructions. This task does not include SDK package publication, commit/push or live Pages setting changes.

## Verification

Focused staging regressions, strict MkDocs build, local rendered navigation/search and workflow event/permission checks. Read-only Pages API returns 404; repository metadata reports `has_pages: false`. The site is not live.

**Evidence:** Staged 20 current Markdown pages; MkDocs 1.6.1 / Material 9.7.7 strict build passed. Three staging regression tests passed, including excluded private files, source/history link mapping, stale-output removal and invalid input preserving an earlier build. Chromium navigated the Python SDK page and source link, found SDK heartbeat results in 15 search links, verified `/Synapse/` resource/search URLs and historical links, and visually confirmed the rendered Mermaid architecture graph. No page JavaScript exceptions occurred; optional GitHub repository-stat requests received API 403 responses without affecting navigation/search. Workflow structural checks passed for PR exclusion, main-only artifact/deploy, strict build, build dependency and scoped permissions. Read-only GitHub repository metadata confirmed `has_pages: false`; no workflow was run remotely and no package/site was published.

## Execution

Existing SDK is a setuptools package with paho-mqtt and Python 3.11+; wheel installation was verified in SYN-104. Source docs are distributed across docs/, root README, sdk/python/README.md and tasks/. Existing repository modifications are preserved.
