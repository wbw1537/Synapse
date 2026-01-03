---
id: MVP-002
type: story
status: Done
---

# MVP Backlog: Service Registry & Logic

## Status: Done
## Effort Estimate: High

## Description
The "brain" of Synapse. It manages the lifecycle of services, handling registration, updates, and availability monitoring (Heartbeats/TTL).

## Requirements
1.  **Data Model**: Define `Service` struct matching the JSON Schema.
    *   Includes: ID, Name, Group, Tags, Status, Icon, URL, TTL.
    *   Complex Types: Actions (Buttons), Widgets (Graphs/Stats), Markdown Docs.
2.  **Ingestion & Auth**:
    *   Listen on MQTT topic `synapse/v1/discovery/+`.
    *   Validate `auth_token` against `SYNAPSE_AUTH_TOKEN` env var.
3.  **CRUD**: Functions to Upsert (Create/Update) services in SQLite.
4.  **TTL Monitor**: A background worker that periodically checks for expired heartbeats.
    *   Query: `UPDATE services SET status = 'offline' WHERE last_seen < now - ttl`.
    *   This ensures services that crash "silently" are marked as offline.

## Tasks
- [x] **Models**: Update `internal/models` with `ServicePayload`, `Action`, `Widget` structs.
- [x] **Config**: Add `AuthToken` to `internal/config`.
- [x] **Service Manager**: Implement `internal/service` package.
    - [x] `NewManager(db)` constructor.
    - [x] `Upsert(payload)` method.
    - [x] `StartTTLMonitor()` background worker.
- [x] **MQTT Handler**: Integrate `ServiceManager` into `cmd/synapse/main.go` and subscribe to topics.

## Test Cases
- [x] **Registration**: Publish valid JSON to MQTT, verify DB record is created with all fields.
- [x] **Auth Failure**: Publish with wrong token, verify **no** DB record is created.
- [x] **Update**: Publish updated JSON (e.g., status change), verify DB reflects change.
- [x] **TTL Expiry**:
    1.  Register service with `ttl=2`.
    2.  Wait 4 seconds.
    3.  Verify status becomes `offline` automatically.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-03
- **Key Artifacts**:
    - `internal/models/service.go`: Full data model.
    - `internal/service/manager.go`: Business logic for Upsert and TTL.
    - `cmd/synapse/main.go`: Integrated Paho MQTT client to consume local messages.
- **Dev Notes**:
    - Decided to use an internal Paho MQTT client in `main.go` to subscribe to the embedded broker. This decouples the core logic from the specific broker implementation.
    - Used `gorm` with `serializer:json` (implicitly via simple struct tags for now, though SQLite stores them as JSON strings mostly).
    - TTL Monitor runs a raw SQL query for efficiency and simplicity.
- **Verification**:
    - Verified registration via a custom Go publisher (`cmd/test-publisher`).
    - Verified TTL expiry logs showing "Marked 1 services as offline".
