---
id: SDK-003
type: task
status: Pending
---

# SDK Task: Core Refactor (New Protocol)

## Status: Pending
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
    -   Rationale: Reduce runner usage during rapid dev, only validate releases/PRs that are release candidates? (Clarify: User said "only ci when new version tag is marked for the PR" - likely means on `push: tags: [ 'v*' ]` or similar).

## Tasks
- [ ] **CI**: Update `.github/workflows/ci.yml` trigger rules.
- [ ] **Models**: Refactor `internal/models/service.go` (Breaking).
- [ ] **Manager**: Update `internal/service/manager.go`.
- [ ] **Frontend Types**: Update TypeScript interfaces.
- [ ] **Frontend UI**: Rewrite `ServiceCard.vue` for Layout/Component rendering.
- [ ] **Cleanup**: Remove any v1 legacy code.

## Resources
-   `docs/axon_toml_spec.md`: The definitive protocol source.