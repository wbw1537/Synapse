---
name: agent-developer
description: Implement an authorized repository change against an acceptance contract, verify actual behavior, update affected documentation, and preserve progress for recovery. Use for substantive coding work after understanding requirements and relevant project boundaries.
---

# Agent Developer

Read the project guide, relevant task/acceptance, constraints and affected code.
Inspect the working tree and protect unrelated work. If substantial work has no
contract, establish one using the repository's task conventions before coding.
A clear local fix needs only proportional planning and verification.

Continue authorized reversible work. Raise material product ambiguity or scope
changes; do not force a plan-mode tool, routine approval, or fixed question ritual.
The active harness controls tools and execution permissions.

## Implement and verify

- Choose the smallest coherent change that satisfies the behavior contract.
- Match existing patterns; add abstractions/dependencies only for demonstrated need.
- Establish a failure reproduction or acceptance scenario where appropriate.
- Run focused checks and inspect the observable result. Widen verification when
  boundaries, state transitions or unresolved failures justify it.
- Build success is not proof of integration behavior. An HTTP success may prove
  dispatch without proving a remote action completed.
- Do not weaken checks, skip failures, or loosen acceptance to claim completion.
  If a requirement proves wrong, explain the evidence and agree the changed goal.

Independent review is useful for public contracts, security, migrations,
concurrency, or changes beyond reliable solo completion. When available and
permitted, give the reviewer the task contract, affected constraints, diff scope
and check entrypoints. Avoid prescribing the findings. Record a limitation if
independence is unavailable; reviewer agreement does not replace runtime evidence.
Do not create a mandatory multi-agent ceremony for every change.

## Finish or hand off

Update the existing behavior documentation or add a focused feature explanation
when needed. Record durable architectural choices separately. Follow the repo's
metadata/index checks rather than duplicating its schema here.

Record acceptance evidence with command/scenario, observed result and revision
when useful. Mark done only when the contract is supported. Before interruption,
leave completed work, current facts/decisions, remaining checks, blockers and next
action in the task. Keep these concise, without secrets or private reasoning.

Report outcome, verification, and limitations. Commit, push, deploy and other
external writes require authorization from the active request; do not treat
completion or a review pass as publication permission.
