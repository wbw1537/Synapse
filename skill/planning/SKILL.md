---
name: project-planner
description: Scope substantial repository changes and create acceptance-based task contracts, including dependency ordering and architecture decisions. Use for new features, cross-layer fixes, refactors, or focused technical discovery; small clear edits need only proportional planning.
---

# Project Planner

Read the repository's agent guide and task conventions, then the relevant product,
architecture, constraints, decisions, task and implementation. Narrow context to
what affects the requested outcome. Separate implemented facts from proposals.

Clarify the observable goal, scope, preserved behavior, uncertain decisions and
blast radius. Resolve routine technical choices autonomously. Ask only when an
unknown changes the product outcome, risk or authorization; record assumptions.

Create or refine one task per independently verifiable behavior increment, using
the repo's template and stable identifiers. A task may span layers. State:

- The operator/user outcome and why it matters.
- Observable acceptance criteria, including important failure/recovery behavior.
- Scope, dependencies and meaningful checks with expected results.
- A short implementation strategy when complexity warrants it.

Discovery tasks may deliver a concrete decision, evidence and follow-up contracts.
Do not label unresolved implementation work ready. Do not pre-specify every file
edit or invent extra product scope to fill a template. Prefer existing patterns.

Plan depth follows uncertainty. A complex or long task needs durable decisions and
resume notes; a local fix does not need a separate design document. Present key
tradeoffs concisely. Existing authorization to implement remains valid; a written
plan is not an automatic new approval gate.

Use the repo's index command after task changes. Planning is complete when the next
step is clear and the acceptance contract can guide implementation and evaluation.
Do not mark acceptance passed until actual verification supports it.
