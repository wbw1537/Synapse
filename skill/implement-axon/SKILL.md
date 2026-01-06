---
name: implement-axon
description: Guide on how to implement Synapse Axons (clients). Covers Smart Axons (interactive) and Raw Axons (reporting only) using MQTT or HTTP. Use this skill when the user asks how to create a new service, sidecar, or agent for Synapse.
---

# How to Implement a Synapse Axon

An **Axon** is a client (script, sidecar, or service) that reports its status to the Synapse Core. This guide covers how to implement both "Smart" (interactive) and "Normal" (reporting-only) Axons.

## 1. Implement a Smart Axon (No SDK)

A "Smart Axon" reports status **and** can execute actions triggered from the dashboard. Since there is no official SDK yet, you must implement the MQTT protocol directly.

### Prerequisites
*   **Transport:** MQTT (TCP port 1883 or WS port 8083)
*   **Libraries:** Any MQTT client (e.g., `paho-mqtt` for Python, `mqtt.js` for Node)

### Step 1: Define the Configuration
You need to construct a JSON payload that defines your service, its widgets, and its actions.

```json
{
  "api_version": "v1",
  "auth_token": "synapse-secret",
  "id": "my-smart-service",
  "name": "Production Database",
  "status": "online",
  "ttl": 30,
  "actions": [
    {
      "id": "restart",
      "label": "Restart Service",
      "style": "danger",
      "confirm": true
    }
  ],
  "widgets": [
    {
      "type": "stat",
      "id": "cpu_stat",
      "label": "CPU",
      "value": 45,
      "unit": "%"
    }
  ]
}
```

### Step 2: Handle Connections & Reporting
1.  **Connect** to the MQTT Broker (`localhost:1883`).
2.  **Subscribe** to the command topic: `synapse/v1/command/{your_service_id}`.
3.  **Publish** your configuration payload to `synapse/v1/discovery/{your_service_id}`.
4.  **Loop**: Re-publish the payload periodically (e.g., every `TTL / 2` seconds) to act as a heartbeat.

### Step 3: Handle Commands
When you receive a message on `synapse/v1/command/{your_service_id}`, parse the JSON payload to find the `action_id`.

```python
# Python pseudo-code example
def on_message(client, userdata, msg):
    payload = json.loads(msg.payload)
    action = payload.get("action_id")
    
    if action == "restart":
        perform_restart()
```

### Full Example (Python)
See `examples/memory_axon.py` in the codebase for a complete working example using `paho-mqtt`.

---

## 2. Implement a Smart Axon (With SDK)

*Current Status: SDK is not yet implemented.*

---

## 3. Implement a Normal Axon (Reporting Only)

A "Normal" or "Raw" Axon only reports status and metrics. It cannot receive commands.

### Option A: MQTT (Recommended)
Same as the Smart Axon, but you do not need to define `actions` or subscribe to the command topic.

1.  **Connect** to MQTT.
2.  **Publish** payload to `synapse/v1/discovery/{id}`.

### Option B: HTTP (Simple Scripts)
Use this for simple bash scripts or cron jobs where keeping a persistent MQTT connection is difficult.

**Endpoint:** `POST http://core-address:8080/api/v1/discovery`

**Bash Example:**
```bash
curl -X POST http://localhost:8080/api/v1/discovery \
  -H "Content-Type: application/json" \
  -d '{
    "api_version": "v1",
    "auth_token": "synapse-secret",
    "id": "backup-job-01",
    "name": "Nightly Backup",
    "status": "online",
    "ttl": 3600,
    "widgets": [
        { "type": "stat", "label": "Last Run", "value": "Success" }
    ]
  }'
```

### Best Practices for Normal Axons
*   **Heartbeat:** Ensure you send the request more frequently than the `ttl` value.
*   **Offline Status:** If your script is exiting (e.g., error or shutdown), try to send one last payload with `"status": "offline"`.

---

## 4. Reference: Widget Types

When building your payload, you can use these widget types in the `widgets` list:

*   **`stat`**: Key-value display (good for text/numbers).
*   **`gauge`**: 0-100 visual progress bar/arc.
*   **`status_indicator`**: Maps text values (e.g., "connected") to colors/icons.
*   **`log_stream`**: Appends text lines to a scrolling log window.
*   **`action_group`**: Buttons to trigger your actions (Smart Axons only).
*   **`link`**: Static hyperlinks to external dashboards.

See `docs/widget_reference.md` for full JSON schemas.
