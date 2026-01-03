# Product Requirement Document: Synapse

| Document Details |  |
| --- | --- |
| **Project Name** | Synapse |
| **Version** | 1.0 (Draft) |
| **Status** | In Planning |
| **Slogan** | *The nervous system of your homelab.* |
| **Repository** | `github.com/wbw1537/synapse` |

---

## 1. Executive Summary

**Synapse** is an open-source, self-hosted infrastructure platform designed for homelabs and small-scale private clouds. Unlike traditional static dashboards (like Homepage) or external monitors (like Uptime Kuma), Synapse utilizes an **Event-Driven Architecture (EDA)** via MQTT.

It acts as a central nervous system where services **self-register**, report their status, provide their own documentation, and expose actionable operations (e.g., "Restart," "Clear Cache"). It combines the functions of a Service Catalog, a Status Dashboard, and a Notification Gateway into a single, lightweight binary.

---

## 2. Problem Statement

Current homelab tooling suffers from three major disconnects:

1. **The "Static Config" Fatigue:** Adding a new service requires manually editing YAML files on the dashboard. This leads to friction and outdated dashboards.
2. **Dead Documentation:** Recovery instructions (Runbooks) exist in a separate wiki or text file, often disconnected from the actual monitoring alert.
3. **Read-Only Operations:** Dashboards show that a service is down, but require the user to SSH into a server to fix it. There is no bridge between "Monitoring" and "Action."
4. **Fragmented Notifications:** Every service (Radarr, Sonarr, PVE) needs its own SMTP/Telegram configuration.

---

## 3. User Personas

### **The "Architect" (Primary Persona)**

* **Profile:** A software developer or sysadmin running a PVE/Docker cluster at home.
* **Needs:** Loves automation. Wants services to "just appear" on the dashboard when deployed. Values low resource usage (Go/Rust preference).
* **Pain Point:** Hates context switching between the dashboard, the terminal, and the wiki when things break.

---

## 4. Product Principles

1. **Push over Pull:** Synapse does not guess where services are; services tell Synapse they exist.
2. **Documentation as Code:** Documentation and operational actions are part of the service payload, not stored in the dashboard database.
3. **Batteries Included, But Swappable:** Comes with a built-in MQTT broker and SQLite DB, but supports external ones.
4. **Single Binary Distribution:** No complex dependencies. One binary (or Docker container) contains the Backend, Frontend, and DB logic.

---

## 5. System Architecture

### 5.1 High-Level Diagram

### 5.2 Technology Stack

* **Core Backend:** **Go (Golang)**. Chosen for concurrency, type safety, and single-binary compilation.
* **Database:** **SQLite** (Embedded, with WAL mode). Stores snapshots and event history.
* **Transport:** **MQTT** (v3.1.1/v5). Used for service registration, heartbeats, and real-time UI updates.
* **Frontend:** **TypeScript + Vue 3 (or React)**. Compiled and embedded into the Go binary via `go:embed`.
* **Notifications:** **Shoutrrr** (Library) integrated natively, plus generic Webhook support.

---

## 6. Functional Requirements

### 6.1 Service Discovery (The "Nervous System")

* **MQTT Ingestion:** The Core must listen to `synapse/v1/discovery/#`.
* **Auto-Registration:** Upon receiving a valid payload, the system creates or updates the service record in the DB.
* **Heartbeat Monitoring (TTL):**
* If a service stops sending heartbeats for > `ttl` seconds, status automatically changes to `offline`.
* Triggers a "Service Lost" notification.


* **Docker Auto-Discovery Agent:**
* A sidecar container that listens to `docker.sock`.
* Parses container labels (e.g., `synapse.enable=true`) and automatically publishes MQTT payloads to the Core.



### 6.2 The Dashboard (Server-Driven UI)

* **Service Cards:** Display status, uptime, and real-time metrics (defined by the `widgets` payload).
* **Context Panel:** Clicking a service opens a side panel displaying:
* Rendered Markdown (`markdown_docs` payload).
* Connection details (IP/Port).
* Tag metadata.


* **Action Buttons:** UI renders buttons based on the `actions` payload.
* *Example:* If payload contains `action: { label: "Restart", id: "restart" }`, UI shows a "Restart" button.



### 6.3 Notification Gateway

* **Unified Routing:** Synapse exposes a single API endpoint for notifications.
* **Deduplication:** Logic to prevent notification storms (e.g., "If > 5 alerts in 10s, send summary").
* **Multi-Channel:** Support Telegram, Discord, Email (via SMTP), and Gotify out of the box.

### 6.4 Action Execution (Remote Ops)

* **Command Flow:** User clicks button -> Core publishes to `synapse/v1/command/{id}` -> Agent receives -> Agent verifies whitelist -> Agent executes local command.
* **Security:** Core must enforce Action ID verification (no raw shell commands allowed).

### 6.5 AI Agent (The "Cortex")

* **Analysis:** When status turns `error`, Synapse sends the service's logs (from payload) and docs to an LLM (OpenAI/Ollama).
* **Suggestion:** The LLM output is displayed on the dashboard (e.g., "Suggestion: The logs indicate OOM. Try increasing memory limit.").

---

## 7. Data Contract (JSON Schema)

The system relies on a strict schema. (Refer to the technical specification for full details).

**Key Fields:**

* `id`: Unique identifier.
* `auth_token`: Security token.
* `status`: `online` | `warning` | `error`.
* `markdown_docs`: Contextual help.
* `actions`: Array of executable operations.
* `widgets`: Array of UI component definitions.

---

## 8. Security Requirements

1. **Token Authentication:** Agents must provide a pre-shared key (PSK) or token in the MQTT payload.
2. **Action Whitelisting:** Agents must strictly define valid Action IDs. Arbitrary command execution is strictly forbidden.
3. **UI Authentication:** The web dashboard must be protected by a login screen (OIDC or simple username/password).

---

## 9. Roadmap

### Phase 1: MVP (Minimum Viable Product)

* [ ] Go Core with embedded SQLite & MQTT Broker.
* [ ] Basic Vue3 Dashboard (Read-only status).
* [ ] API for JSON ingestion via HTTP & MQTT.
* [ ] Simple notification support (Shoutrrr).

### Phase 2: Interactivity

* [ ] Implementation of "Actions" (Command topic).
* [ ] "Docker Watcher" Agent (Go binary).
* [ ] Markdown rendering in UI sidebar.

### Phase 3: Intelligence & Polish

* [ ] AI Agent integration (Log analysis).
* [ ] Server-driven UI Widgets (Gauges, Graphs).
* [ ] Public Release (Documentation, Logo, Docker Hub).

---

## 10. Success Metrics

1. **Deployment Time:** A new user should be able to run `docker run synapse` and see the dashboard in < 60 seconds.
2. **Integration Speed:** Adding an existing Docker container to Synapse (via Labels) should take < 30 seconds.
3. **Performance:** Idle memory usage of the Core container should be < 50MB.