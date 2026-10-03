# Public documentation as the instruction authority

**ID:** SYN-111
**Status:** done
**Created:** 2026-10-03
**Priority:** medium
**Depends-on:** SYN-110

## Goal

Make GitHub Pages the single user-facing instruction reference, backed by website/docs/, while docs/ holds internal development records. Remove duplicated setup and SDK guidance.

## Acceptance

- [x] Public setup, configuration, operations, authentication, SDK and protocol guidance has one maintained home in website/docs/.
- [x] Repository and SDK READMEs direct readers to the website without duplicating instructions.
- [x] Internal development/decision/acceptance/history documents are excluded from the published site and refer to public contracts where needed.
- [x] Obsolete staging and duplicated SDK documentation are removed; current local references and document catalogs are consistent.
- [x] Direct MkDocs strict build, workspace checks and actual browser navigation/search pass, including project subpath URLs.
- [x] Documentation ownership and contribution/publication procedures explain the new boundary without claiming remote deployment.

## Scope

Move authoritative public content, preserve existing product behavior and development history, consolidate SDK instructions, simplify site build and update affected links/tooling. No remote publication or runtime changes.

## Verification

Workspace integrity, focused document coverage checks, strict site build and rendered navigation/search; inspect generated files to confirm only public pages are published.

**Evidence:** Workspace checks passed with 12 valid task contracts and 13 behavioral tests, including public/internal link boundaries and task-viewer paths. MkDocs strict build passed with eight public pages. Generated search index contained exactly those eight pages and no internal records. Chromium traversed all seven guide pages from navigation, found the SDK heartbeat result, verified asset/search URLs under `/Synapse/`, and observed no page exceptions. Pages workflow checks confirmed direct build, public path triggers and retained main-only deployment permissions. Root README reduced from 207 to 31 lines; SDK README reduced to an entry point. Protocol/access/TOML references moved, SDK specification merged, and staging script/tests removed. No remote deployment performed.

## Execution

Baseline 167a551. Found full user setup in root README; SDK lifecycle duplicated across SDK README and docs/sdk_specification.md with overlapping TOML explanations; staging publishes every current internal document.

Public source ownership is recorded in docs/decisions/004-documentation-ownership.md. Historical plans and dated acceptance evidence remain internal; maintained references point to their public owners. Final sources are ready for review and publication through the existing Pages workflow.
