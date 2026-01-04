# Product Requirement Document: Synapse

| Document Details |  |
| --- | --- |
| **Project Name** | Synapse |
| **Version** | 1.1 (Revised) |
| **Status** | In Development |
| **Slogan** | *The nervous system of your homelab.* |
| **Repository** | `github.com/wbw1537/synapse` |

---

## 1. Executive Summary

**Synapse** is an open-source, self-hosted infrastructure platform designed for homelabs and small-scale private clouds. Unlike traditional static dashboards or external monitors, Synapse utilizes an **Event-Driven Architecture (EDA)** via MQTT.

It acts as a central nervous system where services (clients) **self-register**, report their status, provide their own documentation, and expose actionable operations.

The system consists of two main components:

1. **Synapse Core:** The central brain that handles logic, storage, and the dashboard UI.
2. **Axons:** The nerve endings (clients) that run alongside your services to report data and execute commands.

---

## 2. Problem Statement

Current homelab tooling suffers from major disconnects:

1. **The "Static Config" Fatigue:** Adding a new service requires manually editing YAML files on the dashboard.
2. **Dead Documentation:** Recovery instructions exist in separate wikis, disconnected from the actual monitoring alert.
3. **Read-Only Operations:** Dashboards show a service is down, but require SSH access to fix it.
4. **Configuration Drift:** Scripts and agents are hard to manage centrally; changing a threshold usually means editing code on the server.
5. **Recovery Anxiety:** If a monitoring script is lost, the logic and configuration are lost with it.

---

## 3. User Personas

### **The "Architect" (Primary Persona)**

* **Profile:** A software developer or sysadmin running a PVE/Docker cluster.
* **Needs:** Automation, low resource usage (Go), and "Infrastructure as Code" principles.
* **Pain Point:** Wants to write a script once, have it appear on the dashboard, but allow for tuning (alerts/names) via UI without redeploying the script.

---

## 4. Product Principles

1. **Push over Pull:** Synapse does not guess; Axons tell Synapse they exist.
2. **Code is Definition, Core is Policy:** The Axon defines *what* it can do (Actions) and *what* it measures (Widgets); the Core defines *how* it behaves (Alert Thresholds, Display Names).
3. **Bi-directional Source of Truth:** Axons provide the default config; Synapse Core saves the runtime overrides.
4. **Single Binary Core:** The server must remain lightweight and easy to deploy.

---

## 5. System Architecture

### 5.1 Terminology

* **Synapse Core:** The server application (Go + SQLite + MQTT Broker + Vue Frontend).
* **Smart Axon:** A client using the official SDK (Python/Go/Node). It creates a persistent connection, supports bi-directional communication, and can execute **Actions**.
* **Raw Axon:** A simple script (Bash/Curl) that sends uni-directional status updates via HTTP/MQTT. It cannot execute actions.

### 5.2 Technology Stack

* **Core Backend:** **Go (Golang)**.
* **Database:** **SQLite**. Stores snapshots, event history, and user override configurations.
* **Transport:** **MQTT** (v3.1.1/v5) & **HTTP**.
* **Frontend:** **TypeScript + Vue 3**.
* **Notifications:** **Shoutrrr**.

---

## 6. Functional Requirements

### 6.1 Service Discovery (The "Nervous System")

* **MQTT Ingestion:** Core listens to `synapse/register/{id}` and `synapse/report/{id}`.
* **Heartbeat Monitoring (TTL):** Services define their own TTL. If silent for > TTL, status becomes `offline`.

### 6.2 Configuration Management (The "Override Strategy")

* **Definition Layer (Read-Only):** Fields that define the code logic (ID, Action IDs, Widget Data Types) are **immutable** in the UI. They must be updated by the Axon.
* **Policy Layer (Writable):** Fields that define operation preference (Display Name, Icon, Documentation, Alert Rules) can be **overridden** in the Synapse UI.
* **Merge Logic:** On every heartbeat, Core merges the *Axon Payload* with the *Database Overrides* to generate the final state.

### 6.3 Server-Side Monitoring (Expression Engine)

* **Monitors:** Widgets include monitor rules defined as string expressions.
* **Logic:** Uses a Go expression engine (e.g., `expr`) to evaluate logic like `"cpu_value > 90"` or `"status == 'backup_failed'"`.
* **Evaluation:** Occurs on the Server side upon receiving a payload.

### 6.4 Action Execution (Smart Axons Only)

* **Command Flow:** User clicks button -> Core publishes to `synapse/command/{id}` -> Smart Axon receives -> SDK maps Action ID to local function -> Function executes.
* **Security:** Core sends only an Action ID. The Axon SDK strictly maps IDs to pre-defined functions. **No arbitrary shell execution.**

### 6.5 Disaster Recovery (Axon Rehydration)

* **Snapshotting:** Core saves the raw JSON configuration provided by the Axon upon registration.
* **Recovery Kit:** If an Axon is lost (e.g., script deleted), the user can download a "Recovery Kit" from Synapse.
* **Code Generation:** Synapse generates a Python/Node boilerplate file based on the snapshot, including function stubs for all registered Actions and Widgets, allowing the user to fill in the logic gaps.

### 6.6 Notification Gateway

* **Unified Routing:** Single API endpoint for all notifications.
* **Deduplication:** Storm prevention logic.
* **Channels:** SMTP, Telegram, Discord, etc.

---

## 7. Data Contract (Axon Configuration)

Smart Axons must register with a strict JSON schema.

```json
{
  "id": "my-service-01",
  "name": "Production DB",
  "auth_token": "secret-psk",
  "docs": "# Markdown content...",
  "actions": [
    { "id": "restart", "label": "Restart Service", "style": "danger" }
  ],
  "widgets": [
    {
      "id": "cpu",
      "type": "stat",
      "label": "CPU Usage",
      "monitors": [
        { "condition": "value > 90", "severity": "error", "message": "High CPU" }
      ]
    }
  ]
}

```

---

## 8. Security Requirements

1. **Token Authentication:** All Axons must authenticate via PSK.
2. **Action Whitelisting:** Core never sends executable code, only IDs.
3. **UI Authentication:** Dashboard protected by login.
4. **Override Precedence:** User overrides in DB always take precedence over metadata (Name/Docs/Rules) sent by Axon, preventing "config thrashing."

---

## 9. Roadmap

### Phase 1: MVP (The Core)

* [ ] Go Core with embedded SQLite & MQTT Broker.
* [ ] HTTP/MQTT Ingestion API.
* [ ] Basic Dashboard with "Raw Axon" support (Uni-directional).
* [ ] Notification Gateway (Shoutrrr).

### Phase 2: The Smart Axon (Interactivity)

* [ ] **Python SDK** & **Go SDK** development.
* [ ] Implementation of `synapse/command` topic logic.
* [ ] Server-side Expression Engine for alerts.
* [ ] Configuration "Override Strategy" implementation in DB.

### Phase 3: Resilience & Intelligence

* [ ] **Axon Rehydration** (Code generator for lost scripts).
* [ ] AI Agent integration for log analysis.
* [ ] Community Widget Library.

---

## 10. Success Metrics

1. **Time-to-First-Action:** A user should be able to write a 10-line Python Axon that restarts a service via the dashboard in < 5 minutes.
2. **Recovery Confidence:** A user should be able to restore a deleted Axon script using the Synapse UI generator in < 2 minutes.
3. **Performance:** Core runs on < 50MB RAM.