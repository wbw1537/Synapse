# Axon Configuration & Protocol Specification

This document defines the `axon.toml` configuration schema and the Wire Protocol for Synapse.

## 1. Overview

The Synapse Protocol separates **Definition** (UI Layout) from **State** (Component Data).

-   **Configuration (`axon.toml`)**: Defines the static structure, layout, and component types.
-   **Wire Protocol (JSON)**: The payload sent to Synapse Core via MQTT.

## 2. Configuration: `axon.toml`

This file is the source of truth for an Axon's UI and capabilities.

### 2.1 Root Structure

```toml
# Schema version matches the wire protocol version
schema = "axon.card.v1"

[meta]
id = "my-service-id"      # Unique Identifier
name = "My Service"       # Display Name
icon = "server"           # MDI Icon Name
ttl = 30                  # Heartbeat TTL in seconds
group = "Production"      # (Optional) Grouping
description = "A short description."

[layout]
# Defines the visual hierarchy.
type = "sections"         # Currently only "sections" is supported.

    [[layout.section]]
    title = "Network"
    # References IDs defined in [components]
    components = ["wan_ip", "latency"]

    [[layout.section]]
    title = "Resources"
    components = ["cpu_usage", "ram_usage"]

[components]
# Registry of all widgets. Keys are the IDs used in [layout].

    [components.wan_ip]
    type = "stat"
    label = "WAN IP"
    default = "Unknown"   # Initial value before runtime updates

    [components.cpu_usage]
    type = "gauge"
    label = "CPU"
    min = 0
    max = 100
    unit = "%"
```

### 2.2 Component Types

#### `stat`
Displays a key-value pair.
```toml
[components.my_stat]
type = "stat"
label = "Uptime"
unit = "hrs" (optional)
copyable = true (optional)
```

#### `status_indicator`
Maps raw values to visual states.
```toml
[components.my_status]
type = "status_indicator"
label = "Health"
default = "ok"

    [components.my_status.mapping.ok]
    text = "Healthy"
    color = "success"
    icon = "check"

    [components.my_status.mapping.error]
    text = "Critical"
    color = "danger"
    icon = "alert"
    animate = true
```

#### `gauge`
Visual progress bar/arc.
```toml
[components.my_gauge]
type = "gauge"
label = "Memory"
min = 0
max = 1024
unit = "MB"
    
    [components.my_gauge.thresholds]
    800 = "warning"
    1000 = "danger"
```

#### `log_stream`
A list of events.
```toml
[components.app_logs]
type = "log_stream"
label = "Application Logs"
max_items = 50
```

#### `action_group`
Buttons that trigger remote actions.
```toml
[components.controls]
type = "action_group"
label = "Service Control"

    [[components.controls.items]]
    id = "restart_svc"    # Action ID sent to SDK
    label = "Restart"
    style = "danger"
    confirm = true
```

---

## 3. Wire Protocol (JSON)

When the SDK registers with Synapse Core (via MQTT or HTTP), it sends this payload.

**Topic:** `synapse/v1/discovery/{id}`

### 3.1 Structure

```json
{
  "api_version": "v1", 
  "auth_token": "...",
  
  "meta": {
    "id": "my-service-id",
    "name": "My Service",
    "icon": "server",
    "ttl": 30,
    "group": "Production",
    "description": "..."
  },

  "layout": {
    "type": "sections",
    "root": [
      {
        "type": "section",
        "title": "Network",
        "children": ["wan_ip", "latency"]
      },
      {
        "type": "section",
        "title": "Resources",
        "children": ["cpu_usage", "ram_usage"]
      }
    ]
  },

  "components": {
    "wan_ip": {
      "type": "stat",
      "label": "WAN IP",
      "value": "Unknown",
      "props": {
        "copyable": true
      }
    },
    "cpu_usage": {
      "type": "gauge",
      "label": "CPU",
      "value": 0,
      "props": {
        "min": 0,
        "max": 100,
        "unit": "%"
      }
    }
  }
}
```

### 3.2 Key Concepts
1.  **Split Structure**: Instead of a flat `widgets` array, we use `layout` (tree) and `components` (map).
2.  **Stateful Updates**: The `value` field in `components` is the only part that changes at runtime.

### 3.3 Runtime Updates
After registration, the Axon pushes updates by re-publishing the full payload to the **Discovery** topic.
*Optimization Note: Future versions may support partial updates.*