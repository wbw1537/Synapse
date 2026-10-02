# Development workflow

## Task size and autonomy

Clear, reversible small fixes can be implemented directly with proportional
checks. Substantial work gets a persistent task contract first. Long-running or
uncertain work adds a short implementation plan and resume notes to that task.
Ask for missing product decisions, not routine implementation choices. Preserve
existing authorization; do not demand plan approval for every edit.

Use one task per independently verifiable behavior increment. Cross-layer changes
may belong together. A discovery task may end with a concrete design decision
instead of production code; label that outcome explicitly.

## Execution loop

1. Inspect the task, relevant docs/code, working tree, and environment.
2. Establish acceptance and identify meaningful verification before implementation.
3. Mark the task `in-progress`, recording the branch and next action if needed.
4. Implement, run focused checks, and inspect actual behavior.
5. Add independent review for security, migrations, public contracts, concurrency,
   or unusually complex changes when delegation is available and authorized.
   Otherwise record the review limitation. Reviewer verdicts do not replace checks.
6. Update affected docs and task evidence. Mark `done` only when acceptance is
   supported; record partial work and remaining checks honestly.
7. Sync/check the index and present the outcome. Commit/publish only when requested.

## Recovery

Before a handoff, update `## Execution` with facts/decisions, completed work,
remaining work, recent evidence and the next action. Use `blocked` only for an
identified external blocker and include a condition for resuming. On resume,
compare the recorded revision and working tree before trusting old evidence.
Keep decisions and results, not full private conversations or chain-of-thought.

## Skills and portability

Canonical definitions live in `skill/` and travel with Git. Three portable
workflow skills are installed into a user-level, content-versioned Synapse bundle
under `~/.agents/skills/`; local `.agents/skills/` symlinks select that bundle.
Project skills link directly to versioned repo sources. This prevents one global
edit from silently changing all repositories. Reinstall after intentional skill
updates. Other repositories may reuse or adapt the workflow bundle deliberately.

Installation: `python3 scripts/install_skills.py`. It refuses to overwrite a
non-symlink local skill directory and unrelated existing link targets.

## Learning from use

Record observed failures or friction in `docs/workflow-feedback.md`: task ID,
model/harness when known, extra human intervention, verification result, and a
proposed correction. Do not infer improvement from one successful run.
Review several real tasks before moving a convention back to i578112-agent.
Only promote repeatable methods; keep Synapse-specific commands local.
