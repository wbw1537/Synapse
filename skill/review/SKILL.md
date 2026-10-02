---
name: review
description: Review working-tree, branch, or PR changes against their acceptance contract and actual behavior, identifying concrete regressions, missing requirements and risk. Use for requested reviews or justified independent evaluation of substantive changes.
---

# Review

Read the project guide and linked task/requirements before assessing the change.
Use requirements supplied in the request when no task exists; do not demand a new
task merely to review a small fix. State limits when requirements are missing.

Identify the review surface explicitly: unstaged/staged changes including new
files, branch relative to its actual base, or PR. Do not assume `main` or overlook
untracked implementation files. Read affected code and callers, not only the diff.
For PRs obtain metadata and the linked contract before the diff when practical.

## Evaluation

Inspect these dimensions in proportion to the change:

1. Acceptance: does the observable behavior meet the requested contract and preserve
   important existing behavior? Trace cross-layer paths and failure/recovery cases.
2. Correctness: state, ordering, concurrency, resource cleanup, retries, duplicates,
   input validation, version handling and migrations relevant to the surface.
3. Boundaries: architecture, credential handling, access control and external writes.
4. Verification: do checks exercise the changed result and important regressions?
   Check for skipped tests, weakened assertions or evidence from stale code.
5. Maintainability: unjustified complexity, duplicated state/contracts, unexpected
   dependencies, useful error handling and updated integration examples/docs.

Run focused checks when possible in the permitted environment. Prefer a concrete
reproduction or traced failure path. Distinguish a confirmed defect from a
hypothesis and a subjective preference. Do not invent quotas or style findings.
When reviewing a refactor, check preserved behavior against old callers/contracts.

## Report

Lead with actionable findings ordered by impact. Each should identify source
location, trigger, incorrect result, expected behavior and supporting evidence.
Classify severity by consequence, not a fixed checklist category. Separate missing
acceptance from optional future improvements. State checks run and review limits.
If no defects are found, say so without claiming the system is defect-free.

A reviewer normally reports rather than edits. Review completion does not authorize
merge, commit, deploy or contacting others. A model verdict is advisory; tests,
runtime outcomes and product judgment remain independent sources of evidence.
