# Synapse Project Guide

## Overview
Synapse is a self-hosted infrastructure platform acting as a "nervous system" for homelabs. It uses a **Push-based** model where clients ("Axons") self-register via MQTT, eliminating the need for static dashboard configuration.

**Key Technologies:**
- **Backend:** Go (Golang) 1.25+
- **Frontend:** Vue 3 + Tailwind CSS v4 (Embedded in Go binary)
- **Database:** SQLite (Embedded, WAL mode)
- **Broker:** Mochi MQTT (Embedded)

## Directory Structure

### Root
- `cmd/synapse/`: Entry point (`main.go`). Starts DB, Broker, and HTTP server.
- `web/`: The Vue 3 frontend application. Built to `web/dist` and embedded.
- `ui.go`: Go file that embeds the `web/dist` directory.

### Internal (`internal/`)
- `api/`: HTTP Server (Chi), API handlers, and static file serving.
- `broker/`: Embedded MQTT broker implementation.
- `config/`: Environment variable parsing (`.env`).
- `db/`: SQLite connection and GORM models.
- `evaluator/`: Expression engine for server-side monitoring rules.
- `models/`: core data structures (`Service`, `Widget`, `Action`).
- `notification/`: Alert managers (SMTP).
- `service/`: Business logic for Service Registry (CRUD, TTL, Heartbeats).

### Documentation & Planning
- `docs/`: Technical specifications (API, Widgets, Sidecars).
- `plan/`: Project roadmap and task tracking.
  - `mvp/`: Completed initial feature set.
  - `sdk/`: (New) SDK and advanced configuration designs.

## Core Concepts

### 1. Service Registry
Services (Axons) register by publishing a JSON payload to `synapse/v1/discovery/{id}` via MQTT. The `internal/service/manager.go` handles this, validating the token and updating the SQLite DB.

### 2. Heartbeats & TTL
Every service payload includes a `ttl` (Time-to-Live). If an Axon fails to report within this window, the background worker in `internal/service/manager.go` marks it as `offline`.

### 3. Widgets & Actions
- **Widgets**: UI elements (Stats, Gauges, Logs) defined by the Axon.
- **Actions**: Buttons that trigger commands sent back to the Axon via `synapse/v1/command/{id}`.

## Development Workflow

### Prerequisites
- Go 1.25+
- Node.js 20+

### Running Locally
1. **Frontend**:
   ```bash
   cd web && npm install && npm run build && cd ..
   ```
2. **Backend**:
   ```bash
   go run cmd/synapse/main.go
   ```
   The server listens on `:8080` (HTTP) and `:1883` (MQTT).

### Configuration
Copy `.env.example` to `.env`. Key variables:
- `SYNAPSE_AUTH_TOKEN`: Shared secret for Axons.
- `SYNAPSE_DB_PATH`: Location of SQLite DB.

## Key Resources
- **API Spec**: `docs/api_spec.md`
- **Widget Reference**: `docs/widget_reference.md`
- **Future Design**: `plan/sdk/axon_toml_and_sdk.md`
