# 001: Acceptance-based, recoverable development workflow

**Status:** accepted
**Date:** 2026-10-02

## Context

Synapse is being resumed after a protocol refactor. The prior plans contain
conflicting status fields and proposed contracts diverge from implementation.
i578112-agent provides useful repo-local tasks and skills, but fixed approval
rituals, tool-specific instructions and moving task paths add maintenance cost.

## Decision

Keep tasks in Git with stable IDs/paths, observable acceptance, dependency links,
verification evidence and optional resume notes. Generate JSON views from task
Markdown using standard-library tools. Keep T0–T3 as familiar categories, not
priority levels. Maintain a small knowledge map identifying current/proposed/
historical documents. Skills specify outcomes and project procedures while the
active agent harness owns scheduling, permissions and context management.

Use proportional plans/checks and risk-based review. No mandatory plan-mode tool
or automatic publication. Version reusable skills before installing globally.

## Consequences

Agents need fewer routine approvals and can recover from durable task notes.
Current claims remain distinguishable from historical/proposed behavior. Index
and link drift can be checked mechanically. Behavioral quality still requires
real checks and human judgment; metadata validation is not a product test.
The experiment is local to Synapse; evidence should precede upstream adoption.
