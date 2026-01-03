# MVP Backlog: Frontend Dashboard

## Status: Done
## Effort Estimate: High

## Description
A Single Page Application (SPA) served by the Go binary. It displays the grid of services and their live status.

## Requirements
1.  **Tech Stack**: Vue 3 + Vite + TailwindCSS v4.
2.  **Embedding**: Build output (`dist/`) embedded into Go binary via `go:embed`.
3.  **UI Components**:
    *   **Service Grid**: Responsive layout of cards.
    *   **Service Card**: Shows Name, Icon, Status, and simple Widgets.
4.  **Data Sync**: MQTT over WebSocket (Port 8083) for live updates.

## Tasks
- [x] Initialize Vue3 project in `web/` directory.
- [x] Configure Tailwind CSS v4.
- [x] Create `ServiceCard.vue` component.
- [x] Implement State Management (Pinia) with MQTT.js integration.
- [x] Configure `go:embed` in `ui.go` and `internal/api/server.go`.

## Test Cases
- [x] **Build**: `npm run build` generates `dist/` without errors.
- [x] **Serve**: Run Go binary -> Open browser -> UI loads.
- [x] **Real-time**:
    1.  Open Dashboard.
    2.  Publish MQTT status change.
    3.  UI reflects change without full page refresh.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-03
- **Key Artifacts**:
    - `web/`: Vue 3 source code.
    - `ui.go`: Root file for embedding `web/dist`.
    - `internal/api/server.go`: Serving static files and handling SPA routing.
- **Dev Notes**:
    - Used **Tailwind CSS v4** with the new Vite plugin (`@tailwindcss/vite`).
    - Implemented a "Service Store" in Pinia that connects to the embedded MQTT broker over WebSockets.
    - The UI is served at `/` and handles SPA routing by falling back to `index.html`.
- **Verification**:
    - `curl` verified `index.html` is served at `http://localhost:8080/`.
    - Build process is verified via `npm run build`.
