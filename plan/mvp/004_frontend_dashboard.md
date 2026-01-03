# MVP Backlog: Frontend Dashboard

## Status: Pending
## Effort Estimate: High

## Description
A Single Page Application (SPA) served by the Go binary. It displays the grid of services and their live status.

## Requirements
1.  **Tech Stack**: Vue 3 + Vite + TailwindCSS (for rapid UI).
2.  **Embedding**: Build output (`dist/`) must be embedded into Go binary via `go:embed`.
3.  **UI Components**:
    *   **Service Grid**: Responsive layout of cards.
    *   **Service Card**: Shows Name, Icon, Status (Online/Offline/Warning), Uptime tags.
    *   **Detail Panel**: Sidebar/Modal showing JSON/Markdown content (read-only for MVP).
4.  **Data Sync**: Poll `/api/v1/services` OR Listen to MQTT over WebSocket (Preferred for "Event Driven" feel). For MVP, Polling is acceptable fallback, but MQTT WS is better aligned with architecture.

## Tasks
- [ ] Initialize Vue3 project in `frontend/` directory.
- [ ] Create `ServiceCard.vue` component.
- [ ] Implement State Management (Pinia) to hold service list.
- [ ] Implement API Client (fetch or MQTT.js).
- [ ] Configure `go:embed` in `internal/api/static.go` to serve the UI.

## Test Cases
- [ ] **Build**: `npm run build` generates `dist/` without errors.
- [ ] **Serve**: Run Go binary -> Open browser -> UI loads.
- [ ] **Real-time**:
    1.  Open Dashboard.
    2.  Manually update DB/Publish MQTT status change.
    3.  UI should reflect change without full page refresh.
