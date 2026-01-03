# MVP Backlog: Core Infrastructure

## Status: Pending
## Effort Estimate: Medium

## Description
Initialize the Go module, set up the project structure, and integrate the core persistent storage (SQLite) and the message bus (Embedded MQTT Broker).

## Requirements
1.  **Project Structure**: Standard Go layout (`cmd/`, `internal/`, `pkg/`).
2.  **Database**: Embedded SQLite (using `mattn/go-sqlite3` or CGO-free `modernc.org/sqlite` if cross-compilation is priority).
    *   Must enable WAL mode.
    *   Schema migration system (simple SQL files or Go-based migration).
3.  **MQTT Broker**: Embed a pure Go MQTT broker (e.g., `github.com/mochi-mqtt/server`).
    *   Must listen on TCP (default 1883) and WS (WebSockets) for UI.

## Tasks
- [ ] `go mod init github.com/wbw1537/synapse`
- [ ] Create folder structure.
- [ ] Implement `internal/db`: SQLite connection & Schema initialization.
- [ ] Implement `internal/broker`: Mochi MQTT server embedding.
- [ ] Create `cmd/synapse/main.go` entry point to start DB and Broker.

## Test Cases
- [ ] **DB Init**: Run binary, check if `synapse.db` is created.
- [ ] **MQTT Connect**: Run binary, use `mosquitto_pub` to publish a message to `test/topic` and ensure connection succeeds.
- [ ] **Concurrency**: Ensure DB and Broker run concurrently without blocking.
