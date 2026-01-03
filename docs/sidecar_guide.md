# Synapse Sidecar Implementation Guide

This guide explains how to build a "Sidecar" or "Agent" that reports service status to Synapse.

## Core Concept
Synapse uses a **Push Model**. The server does not scan your network. Instead, your services (or a sidecar process) must periodically publish a JSON payload to the MQTT broker.

## 1. Protocol Basics

*   **Transport**: MQTT (v3.1.1 or v5)
*   **Broker**: Included in Synapse (Default Port: `1883`)
*   **Topic**: `synapse/v1/discovery/{service_id}`
*   **Payload Format**: JSON

## 2. JSON Payload Structure

```json
{
  "api_version": "v1",
  "auth_token": "synapse-secret",
  "id": "my-service-id",
  "name": "My Service",
  "status": "online",
  "ttl": 60,
  "widgets": [
    { 
      "type": "stat", 
      "label": "Uptime", 
      "value": 100,
      "monitors": [
        { "condition": "value < 50", "severity": "error", "message": "Uptime Low!" }
      ]
    }
  ]
}
```

## 3. Implementation Example (Python)

Using `paho-mqtt` library:

```python
import time
import json
import paho.mqtt.client as mqtt

# Configuration
BROKER = "localhost"
PORT = 1883
SERVICE_ID = "python-script-01"
TOKEN = "synapse-secret"

# Payload Generator
def get_payload():
    return {
        "api_version": "v1",
        "auth_token": TOKEN,
        "id": SERVICE_ID,
        "name": "Python Worker",
        "group": "Scripts",
        "status": "online",
        "ttl": 10, # 10 seconds TTL
        "widgets": [
            { "type": "stat", "label": "Time", "value": time.strftime("%H:%M:%S") }
        ]
    }

# Connect
client = mqtt.Client()
client.connect(BROKER, PORT, 60)

# Loop
while True:
    payload = json.dumps(get_payload())
    topic = f"synapse/v1/discovery/{SERVICE_ID}"
    
    print(f"Publishing to {topic}...")
    client.publish(topic, payload)
    
    # Sleep less than TTL (e.g., half of TTL) to ensure heartbeat
    time.sleep(5)
```

## 4. Implementation Example (Bash/Curl)

You can also use the HTTP fallback if MQTT is not available.

```bash
#!/bin/bash

curl -X POST http://localhost:8080/api/v1/discovery \
  -H "Content-Type: application/json" \
  -d '{
    "api_version": "v1",
    "auth_token": "synapse-secret",
    "id": "bash-script-01",
    "name": "Bash Job",
    "status": "online",
    "ttl": 60
  }'
```

## 5. Best Practices

1.  **Unique IDs**: Ensure every service has a globally unique `id`.
2.  **Heartbeat Interval**: Publish your payload at an interval **shorter** than the defined `ttl`. A good rule of thumb is `interval = ttl / 2`.
3.  **Graceful Shutdown**: When your service stops, try to publish a final payload with `"status": "offline"` so the dashboard updates immediately.
