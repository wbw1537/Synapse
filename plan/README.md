# Synapse Planning & Roadmap

This directory houses the engineering plans, feature breakdowns, and progress tracking for the Synapse project.

## Directory Structure

- `mvp/`: Backlog items for the Minimum Viable Product.
- `phase_2/`: (Planned) Advanced interactivity and agents.

## Working with Backlog Items

Each file (e.g., `001_core_infra.md`) follows a strict structure:
1. **Metadata**: Status (Pending/In Progress/Done) and Effort Estimation.
2. **Requirements**: High-level goals.
3. **Tasks**: Granular checklist for implementation.
4. **Test Cases**: Specific scenarios to verify the feature.

## Feature Completion Workflow

When a feature is finished, update the corresponding markdown file with the following "Completion Report" section at the bottom:

### Post-Implementation Audit
- **Status**: Updated to `Done`.
- **Completion Date**: [YYYY-MM-DD]
- **Key Artifacts**: List of primary files created (e.g., `internal/db/sqlite.go`).
- **Dev Notes**: Any architectural decisions made during implementation that deviated from the plan.
- **Verification**: Evidence that all **Test Cases** passed (logs, screenshots, or command output).

## Principles
- **Keep it Lightweight**: Avoid over-tooling. Markdown is the source of truth.
- **Test-Driven**: Features are not "Done" until the Test Cases section is verified.
- **Iterative**: It is okay to update a plan file mid-development if technical constraints are discovered.
