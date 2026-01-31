---
id: SDK-003
type: task
status: Done
---

# SDK Task: Core Refactor (New Protocol)

## Status: Done
## Effort Estimate: High

## Description
Refactor Synapse Core and Frontend to exclusively support the new **Axon Protocol** (`layout` + `components`). This is a **breaking change** that replaces the legacy v1 "flat widget list" architecture.

**Goal:** Simplify the codebase by removing legacy compatibility layers and establishing a clean, hierarchical data model.

## Requirements

### 1. Backend (`internal/`)
-   **Models**: Update `Service` struct in `internal/models/service.go`.
    -   **Remove**: `Widgets []Widget` (Legacy).
    -   **Add**: `Layout LayoutSchema` and `Components map[string]Component`.
-   **Service Manager**:
    -   Update `Upsert` logic to validate and store the new structure.
    -   Update `TTL` monitor (no logic change, just ensure it persists).
    -   **Data Migration**: *None*. Drop existing DB tables or advise users to reset `synapse.db` since this is Alpha.

### 2. Frontend (`web/`)
-   **Store**: Update `web/src/stores/services.ts` to match the new model.
-   **ServiceCard**: Rewrite `ServiceCard.vue` to render based on `layout.root` (Sections) -> `components` (Map).
-   **Components**: Update widget components (`StatWidget`, etc.) to receive data from the map structure.

### 3. CI/CD (`.github/workflows/ci.yml`)
-   **Optimization**: Update the workflow triggers.
-   **Rule**: Only run CI checks when a **new version tag** is pushed (e.g., `v*`).

## Tasks
- [x] **CI**: Update `.github/workflows/ci.yml` trigger rules.
- [x] **Models**: Refactor `internal/models/service.go` (Breaking).
- [x] **Manager**: Update `internal/service/manager.go`.
- [x] **Frontend Types**: Update TypeScript interfaces.
- [x] **Frontend UI**: Rewrite `ServiceCard.vue` for Layout/Component rendering.
- [x] **Cleanup**: Remove any v1 legacy code.

## Resources
-   `docs/axon_toml_spec.md`: The definitive protocol source.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-31
- **Key Artifacts**:
    - `internal/models/service.go`: Updated `Service` struct (removed Widgets, added Layout/Components).
    - `internal/service/manager.go`: Updated logic to handle Component maps and Layout merging.
    - `web/src/stores/services.ts`: Frontend type defs updated.
    - `web/src/components/ServiceCard.vue`: UI rendering updated to use `layout` and `components`.
    - `web/src/components/ServiceDrawer.vue`: UI panel updated to look up widgets in `components` map.
- **Verification**:
    - CI config updated to strict version tagging.
    - Code compiles (verified via `go build` and `npm run build`).