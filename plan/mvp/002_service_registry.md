# MVP Backlog: Service Registry & Logic

## Status: Pending
## Effort Estimate: High

## Description
The "brain" of Synapse. It manages the lifecycle of services, handling registration, updates, and availability monitoring (Heartbeats/TTL).

## Requirements
1.  **Data Model**: Define `Service` struct matching JSON Schema (ID, Name, Group, Tags, Status, etc.).
2.  **CRUD**: Functions to Create, Read, Update, Delete services in SQLite.
3.  **TTL Monitor**: A background worker that periodically checks for expired heartbeats.
    *   If `now() - last_seen > ttl`, update status to `offline`.
    *   Fire an event/notification when status changes.

## Tasks
- [ ] Define `Service` struct in `internal/models`.
- [ ] Implement `ServiceRepository` interface in `internal/db`.
- [ ] Implement `ServiceManager` logic (Upsert logic).
- [ ] Implement `HeartbeatMonitor` (Ticker based loop).

## Test Cases
- [ ] **Registration**: Save a service, query DB to verify it exists.
- [ ] **Update**: Update an existing service (e.g., change status), verify change.
- [ ] **TTL Expiry**:
    1.  Register service with `ttl=2` seconds.
    2.  Wait 3 seconds.
    3.  Check DB: Status should be `offline`.
