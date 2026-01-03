# MVP Backlog: API & Ingestion

## Status: Pending
## Effort Estimate: Medium

## Description
Interfaces for the outside world to talk to Synapse. Primarily listening to MQTT topics for service discovery.

## Requirements
1.  **MQTT Ingestion**:
    *   Subscribe to `opshub/v1/discovery/#`.
    *   Payload validation (check if JSON matches schema).
    *   Route valid payloads to `ServiceManager`.
2.  **HTTP API (Read)**:
    *   `GET /api/v1/services`: List all services (for Frontend).
    *   `GET /api/v1/services/{id}`: Details.

## Tasks
- [ ] Implement `internal/api/mqtt_handler.go`: Subscribe and handle messages.
- [ ] Implement JSON Schema validation (basic checks for required fields).
- [ ] Implement `internal/api/http_server.go`: Gin or Standard Lib HTTP server.
- [ ] Connect HTTP server to `ServiceRepository` for read operations.

## Test Cases
- [ ] **Discovery**:
    1.  Publish JSON to `opshub/v1/discovery/my-service`.
    2.  Check HTTP API `/api/v1/services` to see if `my-service` appears.
- [ ] **Invalid Data**: Publish malformed JSON -> Should be ignored/logged (not crash).
