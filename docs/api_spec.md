# Synapse API Specification

This document defines the data models and interfaces for Synapse v1.

## 1. Data Models

### Service Object
The core entity representing a monitored resource.

```json
{
  "id": "my-service-01",
  "name": "My Service",
  "group": "Production",
  "tags": ["database", "critical"],
  "icon": "server", 
  "url": "http://myservice.local",
  "status": "online", 
  "message": "All systems operational",
  "ttl": 60,
  "description": "Primary database server",
  "markdown_docs": "# Runbook\n\nIf down, restart...",
  "actions": [
    {
      "id": "restart",
      "label": "Restart",
      "style": "danger",
      "require_confirmation": true
    }
  ],
  "widgets": [
    {
      "type": "stat",
      "label": "CPU",
      "value": 95,
      "monitors": [
        {
          "condition": "value > 90",
          "severity": "error",
          "message": "CPU Critical"
        }
      ]
    }
  ],
  "last_seen": "2026-01-03T21:00:00Z"
}
```

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | string | **Required**. Unique identifier. |
| `name` | string | Display name. |
| `status` | string | `online` \| `warning` \| `error` \| `offline`. |
| `ttl` | int | Time-to-live in seconds. If no heartbeat received, status becomes `offline`. |

### Widget Object

| Field | Type | Description |
| :--- | :--- | :--- |
| `type` | string | `stat`, `progress`, etc. |
| `value` | any | The data value. |
| `monitors` | array | List of evaluation rules. |

### Monitor Object

| Field | Type | Description |
| :--- | :--- | :--- |
| `condition` | string | Expression to evaluate (e.g. `value > 80`). |
| `severity` | string | `warning` \| `error`. |
| `message` | string | Alert message if true. |

### Discovery Payload
Used for registration via MQTT or HTTP.

```json
{
  "api_version": "v1",
  "auth_token": "your-secret-token",
  ... (All Service Object fields)
}
```

---

## 2. MQTT Interface

**Broker Port**: `1883` (TCP), `8083` (WS)

### Topics

#### Discovery / Heartbeat
*   **Topic**: `synapse/v1/discovery/{service_id}`
*   **Payload**: `Discovery Payload` (JSON)
*   **Description**: Publish to this topic to register or update a service. The `{service_id}` in the topic should match the `id` in the JSON.

---

## 3. HTTP API

**Base URL**: `http://localhost:8080/api/v1`

### Endpoints

#### List Services
*   **GET** `/services`
*   **Response**: `200 OK` (Array of Service Objects)

#### Get Service Detail
*   **GET** `/services/{id}`
*   **Response**: `200 OK` (Service Object) or `404 Not Found`

#### Register Service
*   **POST** `/discovery`
*   **Body**: `Discovery Payload`
*   **Headers**: `Content-Type: application/json`
*   **Response**: `200 OK`

```