---
id: SDK-001
type: task
status: Done
---

# SDK Task: TOML Structure & Protocol Design

## Status: Done
## Effort Estimate: Medium

## Description
Design the `axon.toml` configuration structure and the corresponding JSON wire protocol updates. This task involves bridging the gap between the current flat `widgets` list and the proposed hierarchical `layout` + `components` model.

**Goal:** Create a definitive specification that allows Axons to define *structure* (TOML) separately from *runtime data* (SDK Code).

## Requirements

### 1. TOML Schema Design
-   Define the exact structure of `axon.toml`.
-   **Root Sections**:
    -   `[meta]`: ID, Name, Icon, TTL.
    -   `[layout]`: Hierarchy definition (Sections, potentially grids in future).
    -   `[components]`: Registry of widget definitions (Type, Static Props).
-   **Constraint**: Must support all existing widget types (`stat`, `gauge`, `log_stream`, `action_group`, `status_indicator`).

### 2. Wire Protocol Update (JSON)
-   The current JSON payload places `widgets` (def + value) directly in the root.
-   New Protocol must split this:
    -   `layout`: A JSON tree representing the UI structure.
    -   `components`: A map of `id -> { type, props, value }`.
-   **Backward Compatibility**: Decide if we support the old "flat list" format or require a version bump (e.g., `api_version: "v2"`).

### 3. Gap Analysis & Specs
-   **Backend (`internal/models/service.go`)**:
    -   Identify necessary struct changes to `Service`. The current `Widgets []Widget` is likely insufficient for a nested layout.
    -   Proposed Model: `Layout json.RawMessage` (for flexibility) or a recursive struct.
-   **Frontend (`ServiceCard.vue`)**:
    -   Describe how the frontend will render the new layout.
    -   Currently: `v-for="w in widgets"`.
    -   Target: Recursive component or Section iterator looking up data in `components` map.

## Deliverables
1.  **Document**: `docs/axon_toml_spec.md`
    -   Full TOML example.
    -   Full JSON Payload example.
    -   Field-level descriptions.
2.  **Document**: `docs/protocol_v2_migration.md`
    -   List of changes required in `internal/models`.
    -   List of changes required in `web/`.

## Resources
-   `plan/sdk/axon_toml_and_sdk.md`: High-level design.
-   `internal/models/service.go`: Current data model.
-   `docs/widget_reference.md`: Current widget properties.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-31
- **Key Artifacts**:
    - `docs/axon_toml_spec.md`: Detailed TOML schema and v2 Wire Protocol.
    - `docs/protocol_v2_migration.md`: Gap analysis for Backend/Frontend implementation.
- **Design Decisions**:
    - **Protocol v2**: Split `layout` (tree) and `components` (map) to support hierarchical UIs.
    - **Backward Compatibility**: `Service` model will hold both `Widgets` (v1) and `Layout`/`Components` (v2) fields. Frontend will switch rendering based on `api_version`.
