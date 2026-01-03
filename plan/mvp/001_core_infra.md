---
id: MVP-001
type: story
status: Done
---

# MVP Backlog: Core Infrastructure

## Status: Done
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
- [x] `go mod init github.com/wbw1537/synapse`
- [x] Create folder structure.
- [x] Implement `internal/db`: SQLite connection & Schema initialization.
- [x] Implement `internal/broker`: Mochi MQTT server embedding.
- [x] Create `cmd/synapse/main.go` entry point to start DB and Broker.

## Test Cases
- [x] **DB Init**: Run binary, check if `synapse.db` is created.
- [x] **MQTT Connect**: Run binary, use `mosquitto_pub` to publish a message to `test/topic` and ensure connection succeeds.
- [x] **Concurrency**: Ensure DB and Broker run concurrently without blocking.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-03
- **Key Artifacts**:
    - `cmd/synapse/main.go`: Entry point.
    - `internal/db/sqlite.go`: SQLite connection (modernc.org/sqlite) with WAL mode.
    - `internal/broker/broker.go`: Embedded MQTT broker (mochi-mqtt).
    - `internal/config/config.go`: Environment variable configuration.
- **Dev Notes**:
    - Chose `modernc.org/sqlite` for CGO-free compilation to simplify cross-platform builds.
    - Added `internal/config` to handle environment variables (`SYNAPSE_DB_PATH`, `SYNAPSE_MQTT_PORT`, etc.) using `github.com/caarlos0/env`.
    - **Refactor**: Switched to **GORM** (with `github.com/glebarez/sqlite` for CGO-free support) to simplify DB interactions and schema management, replacing raw SQL.
    - Implemented a basic `InitSchema` in `db` to create the `services` table.
- **Verification**:
    - `synapse.db` is created upon startup.
    - Ports 1883 (TCP) and 8083 (WS) are listening.
    - Application shuts down gracefully on SIGINT.
