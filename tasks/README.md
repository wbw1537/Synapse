# Task management

Task Markdown is authoritative. `index.json` is generated; the browser view is
read-only. Files have stable paths under categories, not status directories:

- `T0-bug-fix`: broken or inconsistent behavior.
- `T1-feature`: new operator/integrator capabilities.
- `T2-quality`: correctness coverage, maintainability and documentation coherence.
- `T3-infrastructure`: development tooling, CI, deployments and security plumbing.

Categories are types, not priorities. Use `high`, `medium`, or `low` separately.
IDs remain stable when a title, status, priority or implementation plan changes.

## State

`backlog` needs refinement; `ready` has a usable contract and can be selected when
dependencies are done; `in-progress` is active work; `blocked` requires an explicit
external blocker; `done` has acceptance evidence; `obsoleted` explains why work
is no longer needed. Do not silently accept old aliases or unknown status values.
Avoid multiple abandoned in-progress tasks; leave a next action on each active task.

## Commands

```sh
python3 tasks/sync_index.py update
python3 tasks/sync_index.py check
python3 tasks/sync_index.py next
./scripts/tasks.sh
```

`next` lists dependency-satisfied ready tasks by priority and stable ID; it does
not assign work or imply authorization. There is no automatic implementation loop.
Validation rejects missing required fields, duplicate IDs, unknown dependencies,
cycles, and malformed completion/recovery records. A done task must have checked
acceptance and concrete evidence; the tool cannot independently prove its claims.

Use [the template](template.md) as a starting point. Describe a verifiable outcome
rather than prescribing every code edit. Include a short plan when useful; add
resume details only for interrupted/long work. Dependencies use stable IDs in
`Depends-on`, comma separated, or `-` for none.

Historical `plan/` files are retained as evidence, not duplicated wholesale into
the current queue. See [migration notes](migration.md). New work lives here.
