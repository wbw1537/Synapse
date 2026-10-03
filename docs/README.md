# Internal development knowledge map

User instructions live in the **[GitHub Pages documentation](https://wbw1537.github.io/Synapse/)**,
whose sole source is [website/docs/](../website/docs/index.md). This directory
records development context and evidence. Do not copy user setup, SDK usage or
configuration instructions into internal records; link to their public owner.

The [catalog](index.json) tracks current, proposed and historical documents.
Current describes implementation or accepted rules, proposed describes future
behavior, and historical preserves earlier intent or evidence. These labels do
not determine what gets published.

## Development records

- [Vision](vision.md): product intent and scope.
- [Architecture](architecture.md): source responsibilities and data paths.
- [Constraints](constraints.md): development boundaries.
- [Development](development.md): task execution, recovery and skills.
- [Testing](testing.md): verification and isolated runtimes.
- [Operations](operation.md): developer release and runtime-change procedures.
- [Feature index](feature/README.md): implementation and evidence pointers.
- [Documentation maintenance](documentation-site.md): ownership, preview and publishing.
- [Acceptance](acceptance.md): integration Alpha outcomes and limits.
- [Feedback](workflow-feedback.md): observed workflow results.
- [Tasks](../tasks/README.md): current contracts and dependency queue.

## Decisions

- [Development workspace](decisions/001-development-workflow.md)
- [Discovery contract](decisions/002-discovery-contract.md)
- [Reliable integration](decisions/003-reliable-integration.md)
- [Documentation ownership](decisions/004-documentation-ownership.md)

## Public contracts for developers

Use the public [protocol](../website/docs/protocol.md),
[access policy](../website/docs/access.md), [SDK guide](../website/docs/python-sdk.md)
and [TOML reference](../website/docs/axon-toml.md). Update those owner pages when
changing public behavior; internal records explain implementation choices and
verification, rather than restating the contract.

## Historical

[Original design](../DESIGN.md), [legacy API](api_spec.md),
[widgets](widget_reference.md), [sidecar guide](sidecar_guide.md) and
[old plans](../plan/README.md) preserve previous designs. They are not current
instructions. Historical reports may name paths that existed at their recorded
revision; use public sources for the present contract.

## Maintenance

Add catalog entries for new documents. Link a successor when superseding a
reference; avoid leaving a second authoritative guide. Run
`python3 scripts/agent.py check` after changes. Code demonstrates behavior, while
accepted decisions and task evidence establish intended changes and verification.
