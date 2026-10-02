# Synapse Agent Guide

Synapse is a push-based homelab operations dashboard. Axons register capabilities,
report state, and receive predefined actions. The core embeds Go, SQLite, MQTT,
and a Vue UI. Keep deployment and integration lightweight.

## Start here

- Read `docs/README.md` for the knowledge map and known documentation conflicts.
- Product intent: `docs/vision.md`; implementation: `docs/architecture.md`.
- Boundaries: `docs/constraints.md`; current wire shape: `docs/protocol.md`.
- Work queue: `tasks/README.md`, `tasks/index.json`, and the relevant task file.
- Environment and checks: `docs/testing.md`; workflow: `docs/development.md`.
- Historical plans are in `plan/`; they are not the current work queue.

Read relevant context, not the entire knowledge base. Report conflicts and record
important decisions rather than silently treating proposed behavior as implemented.

## Working rules

- Check `git status --short` and preserve unrelated changes.
- For substantial work, establish a task with observable acceptance criteria
  before implementation. Small, clear fixes need only proportional verification.
- Continue authorized, reversible work autonomously. Ask when a missing product
  decision changes the outcome or an external operation lacks authorization.
- Prefer a complete behavior increment over splitting work by file or layer.
- Verify the result using the affected behavior, not just compilation or a second
  model's opinion. Do not weaken acceptance criteria or checks to obtain a pass.
- Before interruption or handoff, record the current state, evidence, remaining
  work, and next action in the task. Use `blocked` for a concrete external blocker.
- Update only affected docs. Architecture decisions go in `docs/decisions/`;
  implementation behavior goes in existing feature/protocol docs.
- Run `python3 tasks/sync_index.py update` after task metadata changes and
  `python3 scripts/agent.py check` before handing off workflow/doc changes.
- Commit, push, deploy, and mutate live services only within the user's request.
  Implementation completion alone is not deployment authorization.

## Skills

Versioned definitions live in `skill/`. Install with `python3 scripts/install_skills.py`.
The installer exposes them under `.agents/skills/`; workflow skills also have a
versioned user-level installation under `~/.agents/skills/`.

| Folder / skill name | Use |
| --- | --- |
| `planning` / `project-planner` | Scope substantial work and write an acceptance contract |
| `development` / `agent-developer` | Implement, verify, document, and leave recoverable progress |
| `review` / `review` | Assess actual behavior and diffs against requirements |
| `dev-testing` / `dev-testing` | Synapse builds, isolated runtime checks, and UI verification |
| `operation` / `operation` | Synapse runtime/CI diagnosis and authorized operations |
| `implement-axon` / `implement-axon` | Integrate an Axon with the implemented protocol |

Read the applicable skill before using it. Skills describe project work; the
active harness owns tools, permissions, context management, and agent scheduling.
Do not require tools from a different client or spawn agents merely to follow a ritual.

## Source layout

`cmd/synapse/` wires the process; `internal/models/` defines the data contract;
`internal/service/` owns persistence, TTL, monitors, and commands; `internal/api/`
serves HTTP; `internal/broker/` serves MQTT; `web/src/` renders the dashboard.
`examples/` contains integrations. The Python SDK is planned, not implemented.

## Credentials and local state

Use `.env.example` to discover configuration. Do not print secret files or include
credentials in task evidence. The Go process reads environment variables, not
`.env` automatically. Keep test databases outside the repository and do not use
`synapse.db` for disposable validation. Build output is not source evidence.
