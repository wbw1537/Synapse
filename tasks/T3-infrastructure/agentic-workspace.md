# Establish a recoverable agentic workspace

**ID:** SYN-001
**Status:** done
**Created:** 2026-10-02
**Priority:** high
**Depends-on:** -

## Goal

Provide maintained task contracts, document routing, reusable skills and repeatable workspace checks so Synapse can be resumed and developed autonomously.

## Acceptance

- [x] Tasks have stable IDs/paths, strict states, dependency validation and a reproducible generated index.
- [x] Current implementation, proposals and historical plans are clearly distinguished in the document map.
- [x] Six complete skills are discoverable locally; reusable workflow skills are installed as a versioned global bundle.
- [x] Task viewer works without external libraries and does not expose unrelated repository files.
- [x] Workspace tooling rejects invalid contracts, dependency cycles, metadata drift and conflicting skill installations.
- [x] Feedback records distinguish observations from unproven workflow improvements.

## Scope

Developer environment only. Preserve historical plans and runtime data; do not change application protocol, deploy, or modify i578112-agent.

## Verification

**Evidence:** 2026-10-02, source baseline `36f9ac8` plus the uncommitted workspace change:

- `python3 -m unittest discover -s scripts/tests -v`: 11 tests passed for task
  state/evidence, metadata drift, dependencies, installer conflicts and viewer paths.
- `python3 scripts/agent.py check`: task contracts, complete generated index,
  documentation catalog/maintained links and skill entrypoints passed.
- `python3 scripts/install_skills.py --check` and `doctor`: six local skill links
  and the pinned global workflow bundle `be3a7c145dfb9e5d` verified.
- Skill Creator `quick_validate.py`: all six definitions passed.
- `GOCACHE=/tmp/synapse-go-cache python3 scripts/agent.py verify --scope all`:
  Vue typecheck/build, Go vet/test and integrated binary build passed. Go reports
  no test files; this is build/static evidence, not application behavior coverage.
- Actual viewer HTTP: page/index 200, secret path 404. Cached Chromium exercised
  selection, dependency navigation, search, status filtering and refresh with no
  page errors. Screenshot: `/tmp/synapse-tasks-viewer.png` (ephemeral evidence).
- Temporary-storage launcher started on 16080/11883/18083, HTTP UI/list returned
  200 with an empty service registry; Ctrl+C shut down and removed temporary data.
- `node --check` on viewer JavaScript, `bash -n scripts/tasks.sh`, and
  `git diff --check` passed.

## Execution

Workspace setup complete; seven product/quality tasks remain in the queue.
No application source, deployed runtime, existing database or i578112 files changed.
The default browser CLI expected missing system Chrome; validation used already
cached Chromium instead of installing software. Network-listening validation and
user-level skill installation required sandbox filesystem/network permission.

Next product task: SYN-101. Future workflow observations belong in
`docs/workflow-feedback.md`; no upstream improvement claim has been established.
