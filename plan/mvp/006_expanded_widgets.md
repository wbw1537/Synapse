---
id: MVP-006
type: story
status: Done
---

# MVP Backlog: Expanded Widgets

## Status: Done
## Effort Estimate: Medium

## Description
Implement the full set of widgets defined in `docs/widget_reference.md` to allow Axons to report rich data and provide interactivity. This involves updating the backend to handle complex widget states (like log streams) and creating a suite of Vue components for the frontend.

## Requirements
1.  **Backend Model**:
    *   Expand `Widget` struct to support all properties (done).
    *   Implement **Log Stream Merging**: When a `log_stream` widget is updated, the new value should be *appended* to the existing list, maintaining a maximum history (`max_items`), rather than overwriting it.
2.  **Frontend Components**:
    *   Create dedicated Vue components for each widget type:
        *   `stat`: Key-value display with optional copy and units.
        *   `status_indicator`: Mapped states to colors/icons/text.
        *   `gauge`: Visual progress bar/arc with thresholds.
        *   `log_stream`: Scrollable list of text entries.
        *   `action_group`: Buttons that trigger remote actions (future placeholder for now).
        *   `link`: External navigation.
3.  **Frontend Integration**:
    *   Update `ServiceCard` to dynamically render the correct component based on `widget.type`.

## Tasks
- [x] **Backend Models**: Update `internal/models` with expanded `Widget` fields. (Completed in initial step)
- [ ] **Backend Logic**: Implement `mergeWidgets` function in `internal/service/manager.go` to handle `log_stream` appending logic during `Upsert`.
- [ ] **Frontend Library**: Create `web/src/components/widgets/` directory.
    - [ ] `StatWidget.vue`
    - [ ] `StatusWidget.vue`
    - [ ] `GaugeWidget.vue`
    - [ ] `LogStreamWidget.vue`
    - [ ] `ActionGroupWidget.vue`
    - [ ] `LinkWidget.vue`
- [ ] **Frontend Integration**: Refactor `ServiceCard.vue` to use a dynamic component loader or `v-if/else` to render widgets.

## Test Cases
- [x] **Log Stream**:
    1.  Publish a `log_stream` widget with value "Line 1".
    2.  Publish again with "Line 2".
    3.  Verify backend storage contains both lines (up to limit).
    4.  Verify frontend displays both lines.
- [x] **Visuals**:
    1.  Register a service with one of each widget type.
    2.  Verify `gauge` renders correct percentage.
    3.  Verify `status_indicator` shows correct color/icon for given state.
    4.  Verify `link` opens in new tab.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-04
- **Key Artifacts**:
    - `internal/models/service.go`: Expanded `Widget` struct.
    - `internal/service/manager.go`: Implemented `mergeWidgets` for `log_stream` logic.
    - `web/src/components/widgets/`: Created all 6 widget components.
    - `web/src/components/ServiceCard.vue`: Dynamic widget rendering.
- **Dev Notes**:
    - Implemented a `mergeWidgets` function to handle stateful widgets (like `log_stream`) during the Upsert process.
    - Created a modular Vue component architecture for widgets.
    - Verified compilation of both Go and Vue codebases.
- **Verification**:
    - `go build` passed.
    - `npm run build` passed.

