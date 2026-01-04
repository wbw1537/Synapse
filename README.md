# Synapse

**The nervous system of your homelab.**

Synapse is an open-source, self-hosted infrastructure platform designed for homelabs. It acts as a central dashboard where services **self-register** via MQTT, providing a "Push-based" alternative to traditional static dashboards.

![Status](https://img.shields.io/badge/Status-Alpha-orange) ![Go](https://img.shields.io/badge/Go-1.23-blue) ![Vue](https://img.shields.io/badge/Vue-3-green)

## Features

*   **Self-Registration**: Services announce themselves. No more editing YAML config files for the dashboard.
*   **Live Status**: Real-time updates via MQTT WebSockets.
*   **TTL Monitoring**: Automatic "Offline" detection if a service stops reporting.
*   **Single Binary**: Backend (Go), Database (SQLite), Broker (MQTT), and Frontend (Vue) all in one executable.

## Quick Start

### 1. Build from Source

Requirements: Go 1.23+, Node.js 20+

```bash
# 1. Clone the repo
git clone https://github.com/wbw1537/synapse.git
cd synapse

# 2. Build the Frontend
cd web
npm install
npm run build
cd ..

# 3. Build the Backend (embeds frontend)
go build -o synapse cmd/synapse/main.go

# 4. Run
./synapse
```

The dashboard will be available at **http://localhost:8080**.

### 2. Configuration

Synapse looks for a `.env` file in the working directory. You can copy the example to start:

```bash
cp .env.example .env
```

| Variable             | Default          | Description                            |
|:---------------------|:-----------------|:---------------------------------------|
| `SYNAPSE_HTTP_PORT`  | `:8080`          | Port for the Web UI and HTTP API.      |
| `SYNAPSE_MQTT_PORT`  | `:1883`          | TCP Port for the embedded MQTT Broker. |
| `SYNAPSE_WS_PORT`    | `:8083`          | WebSocket Port for MQTT (used by UI).  |
| `SYNAPSE_DB_PATH`    | `synapse.db`     | Path to the SQLite database file.      |
| `SYNAPSE_AUTH_TOKEN` | `synapse-secret` | PSK for service registration.          |

Use this to export the environment variables.

```sh
export $(grep -v '^#' .env | xargs)
```

## Usage

### Registering a Service

You can register a service using **MQTT** (preferred) or **HTTP**.

**Option A: MQTT (Python Example)**
```python
import paho.mqtt.publish as publish
import json

payload = {
    "api_version": "v1",
    "auth_token": "synapse-secret",
    "id": "my-service",
    "name": "My Service",
    "status": "online",
    "ttl": 30
}

publish.single("synapse/v1/discovery/my-service", json.dumps(payload), hostname="localhost")
```

**Option B: HTTP (Curl Example)**
```bash
curl -X POST http://localhost:8080/api/v1/discovery \
  -d '{"id": "my-service", "name": "My Service", "status": "online", "ttl": 30, "auth_token": "synapse-secret"}'
```

See [docs/sidecar_guide.md](docs/sidecar_guide.md) for detailed integration guides.

## Architecture

*   **Core**: Go (Golang)
*   **DB**: SQLite (Embedded, WAL mode)
*   **Broker**: Mochi MQTT (Embedded)
*   **Frontend**: Vue 3 + Tailwind CSS + Pinia

## Troubleshooting

### SMTP / Email Alerts
*   **Port 587**: Synapse assumes `STARTTLS` when using port 587. It connects via plain TCP first, then upgrades.
*   **Port 465**: Implicit SSL/TLS is currently *not* supported by the default `net/smtp` implementation used here.
*   **Auth Failed**: If you see `SASL PLAIN authentication failed`, verify your password is correct in `.env` (ensure no trailing spaces) and that your SMTP server accepts Plain Auth from your IP.

## License

MIT

