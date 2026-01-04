---
id: MVP-007
type: story
status: Pending
---

# MVP Backlog: Service Detail Side-Panel

## Status: Pending
## Effort Estimate: Medium

## Description
Implement a slide-over side-panel (drawer) to display detailed information about a selected service. This panel addresses the need to view runbooks (documentation), logs, and perform critical actions without cluttering the main dashboard card.

## Requirements
1.  **UX/UI**:
    *   **Slide-over**: A drawer sliding in from the right when a service card is clicked.
    *   **Header**: Service Icon, Name, IP/URL, Status, Edit button (placeholder), Close button.
    *   **Tabs**: Navigation for "Runbook" (Docs), "Logs" (Log Streams), and "Health" (Stats/Widgets).
    *   **Action Area**: Sticky bottom footer for "Critical" actions (those requiring confirmation or context).
2.  **Components**:
    *   `ServiceDrawer.vue`: The main container component.
    *   Markdown Renderer: To render `service.markdown_docs`.
    *   Widget Reuse: Reuse existing widget components (`LogStreamWidget`, `StatWidget`, etc.) inside the panel.
3.  **Interaction**:
    *   Clicking a `ServiceCard` opens the drawer.
    *   Clicking "X" or outside the drawer closes it.
    *   Critical actions in `action_group` widgets should be moved to or accessible via the panel (or specifically designated for the panel, though for now we might just render all widgets in the "Health" or "Logs" tab and keep specific actions in the footer if we parse them out).
    *   *Refinement*: The requirement mentions separating "Quick" vs "Critical" actions.
        *   **Quick**: Stay on Card.
        *   **Critical**: Move to Panel Footer.
        *   *Implementation*: We need a way to distinguish. For MVP, we can treat `style="danger"` or `confirm=true` as candidates for the panel, OR just replicate them in the panel for now. The prompt says "separate them", implying a logical split. We will assume for this task we render `action_group` widgets in the panel's footer or specific tab if they are deemed critical. However, since the backend model just has `actions` (top level) and `widgets` (which include `action_group`), we might need to parse the top-level `actions` list for the footer.
        *   *Clarification*: The `Service` model has a top-level `Actions []Action`. It also has `Widgets` which can be `action_group`. The `action_group` widget references `action_id`s.
        *   *Decision*: We will render the top-level `Actions` in the sticky footer if they exist. We will display the `markdown_docs` in the "Runbook" tab. We will display `log_stream` widgets in the "Logs" tab. We will display other widgets in the "Health" tab.

## Tasks
- [ ] **Dependencies**: Install `markdown-it` (or similar) for rendering docs.
- [ ] **Store Update**: Add `selectedServiceId` state to `useServiceStore`.
- [ ] **Component: ServiceDrawer**:
    -   Implement the sliding transition.
    -   Implement the Tab interface.
    -   Implement Markdown rendering for Runbook.
    -   Implement Widget filtering (Logs vs Stats).
    -   Implement Sticky Footer for Actions.
- [ ] **Component: ServiceCard**: Update to toggle the drawer on click (excluding the URL link).
- [ ] **Integration**: Add `<ServiceDrawer />` to `App.vue`.

## Test Cases
- [x] **Open/Close**: Click card -> Panel opens. Click X -> Panel closes.
- [x] **Data Sync**: Open panel for Service A. Publish update to Service A. Panel should reflect changes (e.g. new log line).
- [x] **Markdown**: Verify `markdown_docs` renders correctly as HTML.
- [x] **Tabs**: Verify switching between Runbook/Health/Logs works.
- [x] **Actions**: Verify buttons in the footer trigger the (mock) action.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-04
- **Key Artifacts**:
    - `web/src/components/ServiceDrawer.vue`: The slide-over panel component.
    - `web/src/components/widgets/LogStreamWidget.vue`: Enhanced with `autoHeight` support.
    - `examples/memory_axon.py`: Added sample runbook documentation.
- **Dev Notes**:
    - Integrated `markdown-it` for safe documentation rendering.
    - Reordered tabs to prioritize Runbook -> Health -> Logs based on user feedback.
    - Implemented a sticky footer for critical actions in the drawer.
- **Verification**:
    - Build successful.
    - Manual verification of panel transitions and markdown rendering.

