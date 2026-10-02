# Synapse knowledge map

Current documentation describes implementation or accepted development rules.
Proposed documentation describes future behavior. Historical documents preserve
prior intent and reports; do not use them as proof of present implementation.
The machine-readable [catalog](index.json) tracks these distinctions.

## Current

- [Vision](vision.md): the operator workflow and product scope.
- [Architecture](architecture.md): implementation paths and synchronization gaps.
- [Constraints](constraints.md): development boundaries and Alpha limitations.
- [Protocol](protocol.md): flat identity/properties with layout and components.
- [Development](development.md): task execution, recovery, skills and review.
- [Testing](testing.md): environment, checks and disposable-storage runtime.
- [Operations](operation.md): runtime diagnosis and release authorization.
- [Feature index](feature/README.md): implemented behavior explanations.
- [Workflow decision](decisions/001-development-workflow.md): accepted environment design.
- [Discovery decision](decisions/002-discovery-contract.md): accepted flat v1 and legacy rejection.
- [Feedback](workflow-feedback.md): actual results and candidates for upstream reuse.
- [Tasks](../tasks/README.md): current work contracts and dependency queue.

## Proposed

[Axon TOML](axon_toml_spec.md) and [SDK specification](sdk_specification.md) describe
planned client development. Their mapping now targets the maintained flat
protocol; the configuration parser and SDK are still unimplemented.

## Historical

[Original design](../DESIGN.md), [old API reference](api_spec.md),
[widget reference](widget_reference.md), [sidecar guide](sidecar_guide.md) and
[old plans](../plan/README.md) are preserved. The example and legacy Axon skill
were written against widgets/actions arrays; the maintained skill is now updated,
but the Python example still requires its tracked migration.

## Maintenance

Read current intent/constraints and the task, then relevant implementation.
When they disagree, identify the conflict; code demonstrates current behavior
but does not silently override an accepted design decision. Proposed docs are
not implementation requirements until the task adopts them.

Add catalog entries when adding docs. Mark superseded docs historical and link
the successor. Avoid duplicating behavior across multiple authoritative documents.
Run `python3 scripts/agent.py check` after changes. It checks current links and
catalog coverage, not the truth of prose or historical links. Maintain examples
and executable checks alongside public contracts.
