# Synapse

**The nervous system of your homelab.**

Synapse is an open-source, self-hosted infrastructure platform designed for homelabs. It acts as a central dashboard where **Axons** (clients) self-register via MQTT, providing a "Push-based" alternative to traditional static dashboards.

![Status](https://img.shields.io/badge/Status-Alpha-orange) ![Go](https://img.shields.io/badge/Go-1.25-blue) ![Vue](https://img.shields.io/badge/Vue-3-green)

## Features

*   **Self-Registration**: Axons announce themselves. No more editing YAML config files for the dashboard.
*   **Live Status**: Real-time validated state via same-origin server events.
*   **TTL Monitoring**: Automatic "Offline" detection if an Axon stops reporting.
*   **Single Binary**: Synapse Core (Go), Database (SQLite), Broker (MQTT), and Frontend (Vue) all in one executable.

### Architecture

```mermaid
graph LR
    Axons["Axons (Nerve Endings)"] -- "Self-Register & Report" --> Core["Synapse Core (The Brain)"]
    Core -- "Live Dashboard" --> User((User))
```

The system consists of two main components:

1.  **Synapse Core**: The central brain.
    *   **Backend**: Go (Golang)
    *   **DB**: SQLite (Embedded, WAL mode)
    *   **Broker**: Mochi MQTT (Embedded)
    *   **Frontend**: Vue 3 + Tailwind CSS + Pinia
2.  **Axons**: The nerve endings (clients).
    *   Small scripts or sidecars running alongside your services.
    *   Report status and stats via MQTT.

## Development workspace

Read [AGENTS.md](AGENTS.md) and the [knowledge map](docs/README.md) when resuming development.
Current work lives in [tasks/](tasks/README.md); old `plan/` files are historical.
See [documentation website setup](docs/documentation-site.md) for local preview
and GitHub Pages publishing.

```sh
python3 scripts/install_skills.py
python3 scripts/agent.py doctor
python3 scripts/agent.py check
python3 tasks/sync_index.py next
./scripts/tasks.sh
```

See [testing](docs/testing.md) for builds and temporary-storage development.
The reference [Python SDK](sdk/python/README.md) is available. The Python example uses the maintained
[layout/component protocol](docs/protocol.md).

## Quick Start

### 1. Docker (Recommended)

The easiest way to run Synapse is using Docker Compose.

```bash
# 1. Clone the repo
git clone https://github.com/wbw1537/synapse.git
cd synapse

# 2. Configure distinct random Axon/operator secrets in .env
cp .env.example .env

# 3. Start the application after configuring the secrets
docker compose up -d
```

### 2. Build from Source

Requirements: Go 1.25+, Node.js 20+

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

# 4. Set distinct random secrets, then run
export SYNAPSE_AUTH_TOKEN="<axon-token>"
export SYNAPSE_ADMIN_TOKEN="<different-operator-token>"
./synapse
```

The dashboard will be available at **http://localhost:8080**. Sign in with the
operator token. Native and Compose host listeners default to loopback; see the
[access policy](docs/access.md) before admitting remote clients.

### 2. Configuration

The Go process reads environment variables. It does not load `.env` automatically.
Copy the example and export its values before running:

```bash
cp .env.example .env
```

| Variable                | Default             | Description                                      |
|:------------------------|:--------------------|:-------------------------------------------------|
| **Core**                |                     |                                                  |
| `SYNAPSE_HTTP_PORT`     | `127.0.0.1:8080`             | Port for the Web UI and HTTP API.                |
| `SYNAPSE_MQTT_PORT`     | `127.0.0.1:1883`             | TCP Port for the embedded MQTT Broker.           |
| `SYNAPSE_WS_PORT`       | `127.0.0.1:8083`             | MQTT WebSocket listener; UI uses HTTP SSE.            |
| `SYNAPSE_DB_PATH`       | `synapse.db`        | Path to the SQLite database file.                |
| **Security**            |                     |                                                  |
| `SYNAPSE_AUTH_TOKEN`    | required    | Axon connection and registration secret.                    |
| `SYNAPSE_ADMIN_TOKEN` | required | Distinct operator login/API secret. |
| `SYNAPSE_COOKIE_SECURE` | `false` | Set true behind an HTTPS reverse proxy. |
| **Notifications**       |                     |                                                  |
| `SYNAPSE_ENABLE_ALERTS` | `false`             | Enable SMTP email notifications.                 |
| `SYNAPSE_SMTP_HOST`     |                     | SMTP Server Hostname (e.g., smtp.gmail.com).     |
| `SYNAPSE_SMTP_PORT`     | `587`               | SMTP Port (587 for STARTTLS).                    |
| `SYNAPSE_SMTP_USER`     |                     | SMTP Username.                                   |
| `SYNAPSE_SMTP_PASS`     |                     | SMTP Password.                                   |
| `SYNAPSE_SMTP_FROM`     | `synapse@localhost` | Sender email address.                            |
| `SYNAPSE_SMTP_TO`       |                     | Comma-separated list of recipient emails.        |

Use this to export the environment variables.

```sh
export $(grep -v '^#' .env | xargs)
```

## Usage

### Running an Axon

Axons are the clients that report data to Synapse. We provide a sample Python Axon that monitors memory usage and simulates a critical alert.

**1. Install Dependencies**
```bash
pip install paho-mqtt psutil
```

**2. Run the Axon**
```bash
SYNAPSE_AUTH_TOKEN="<matching-core-token>" python3 examples/memory_axon.py
```

**What happens?**
The script will register a "Memory Monitor" service on the dashboard and cycle through a simulation:
1.  **Normal (10s)**: Reports ~45% memory usage.
2.  **Critical (5s)**: Simulates a spike to 95%, triggering a "Critical" alert on the dashboard.
3.  **Normal (Resumed)**: Returns to normal reporting.

The script reads `SYNAPSE_AUTH_TOKEN`, optional `SYNAPSE_MQTT_HOST`,
`SYNAPSE_MQTT_PORT` (numeric port), and `SYNAPSE_SERVICE_ID` from the environment.
It keeps five log samples and handles declared actions through simulated callbacks.
The built-in scenario uses simulated readings; adapt it to monitor real services.

### Manual Registration

You can also manually register an Axon using **MQTT** (preferred) or **HTTP**.

**Option A: MQTT**
```python
import os
import paho.mqtt.publish as publish
import json

payload = {
    "api_version": "v1",
    "auth_token": os.environ["SYNAPSE_AUTH_TOKEN"],
    "id": "my-service",
    "name": "My Service",
    "status": "online",
    "ttl": 30,
    "layout": {"type": "sections", "root": []},
    "components": {}
}

publish.single("synapse/v1/discovery/my-service", json.dumps(payload),
               hostname="localhost", client_id="my-service",
               auth={"username": "axon", "password": os.environ["SYNAPSE_AUTH_TOKEN"]})
```

**Option B: HTTP (Curl)**
```bash
curl -X POST http://localhost:8080/api/v1/discovery \
  -H "Content-Type: application/json" \
  -d '{"api_version":"v1","id":"my-service","name":"My Service","status":"online","ttl":30,"auth_token":"<axon-token>","layout":{"type":"sections","root":[]},"components":{}}'
```

Use the [Python SDK](sdk/python/README.md), [maintained protocol](docs/protocol.md)
and [access policy](docs/access.md) for new integrations.

## Troubleshooting

### SMTP / Email Alerts
*   **Port 587**: Synapse assumes `STARTTLS` when using port 587. It connects via plain TCP first, then upgrades.
*   **Port 465**: Implicit SSL/TLS is currently *not* supported by the default `net/smtp` implementation used here.
*   **Auth Failed**: If you see `SASL PLAIN authentication failed`, verify your password is correct in `.env` (ensure no trailing spaces) and that your SMTP server accepts Plain Auth from your IP.

## License

Apache-2.0
