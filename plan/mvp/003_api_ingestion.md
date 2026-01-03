---
id: MVP-003
type: story
status: Pending
---

# MVP Backlog: API & Ingestion

## Status: Done
## Effort Estimate: Medium

## Description
Interfaces for the outside world to talk to Synapse. Primarily listening to MQTT topics for service discovery.

## Requirements
1.  **MQTT Ingestion** (Completed in 002):
    *   Subscribe to `synapse/v1/discovery/#`.
    *   Payload validation (check if JSON matches schema).
    *   Route valid payloads to `ServiceManager`.
2.  **HTTP API (Read/Write)**:
    *   **Framework**: `go-chi/chi`.
    *   `GET /api/v1/services`: List all services (JSON).
    *   `GET /api/v1/services/{id}`: Get details for a specific service.
    *   `POST /api/v1/discovery`: HTTP fallback for service registration.

## Tasks
- [x] Implement `internal/api/mqtt_handler.go`: Subscribe and handle messages. (Handled in `main.go`)
- [x] Implement JSON Schema validation. (Handled in `ServiceManager`)
- [x] Implement `internal/api/server.go`: Chi HTTP server setup.
- [x] Implement Handlers:
    - [x] `ListServices`
    - [x] `GetService`
    - [x] `RegisterServiceHTTP`
- [x] Connect HTTP server to `ServiceManager` and `DB`.

## Test Cases
- [x] **HTTP Discovery**: POST JSON to `/api/v1/discovery` -> Check DB.
- [x] **List API**: GET `/api/v1/services` -> Should return list of services.
- [x] **Detail API**: GET `/api/v1/services/{id}` -> Should return full details.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-03
- **Key Artifacts**:
    - `internal/api/server.go`: Chi Router with CORS and endpoints.
    - `internal/service/manager.go`: Added `List()` and `Get()` methods.
- **Dev Notes**:
    - Used `go-chi/chi` for routing.
    - Enabled CORS for all origins (`*`) for MVP development ease.
    - Exposed `8080` (default) for HTTP API.
- **Verification**:
    - Verified all endpoints using `curl` in `test_api.sh`.
