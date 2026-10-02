> Proposed reference. See [knowledge map](README.md) and
> [implemented protocol](protocol.md) before relying on this document.

# Axon Configuration & Protocol Specification

This document proposes the `axon.toml` configuration schema. The maintained
wire contract is [Discovery protocol](protocol.md); the SDK is not implemented.

## 1. Overview

The Synapse Protocol separates **Definition** (UI Layout) from **State** (Component Data).

-   **Configuration (`axon.toml`)**: Defines the static structure, layout, and component types.
-   **Wire Protocol (JSON)**: The payload sent to Synapse Core via MQTT.

## 2. Configuration: `axon.toml`

This file is the source of truth for an Axon's UI and capabilities.

### 2.1 Root Structure

```toml
# Configuration schema identifier; maps to wire api_version = "v1"
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
    components = ["wan_ip"]

    [[layout.section]]
    title = "Resources"
    components = ["cpu_usage"]

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

## 3. Mapping TOML to the maintained wire contract

The SDK must emit the flat [discovery snapshot](protocol.md), never nested `meta`
or `props`. This mapping is a planned SDK requirement, not an implemented parser.

| Configuration / runtime source | JSON destination |
| --- | --- |
| `schema = "axon.card.v1"` | `api_version = "v1"`; reject other configuration schemas |
| `[meta]` fields | Copy to the top level (`id`, `name`, `icon`, `ttl`, `group`, `description`, etc.) |
| Runtime credentials | Top-level `auth_token`; do not store real tokens in committed TOML |
| Runtime health | Top-level `status` and `message`; initialize status as `online` |
| `[layout].type` | `layout.type` (`sections`) |
| `[[layout.section]]` | An element of `layout.root`, adding `type = "section"` |
| Section `components` | Section `children` array |
| `[components.<key>]` | Entry in the `components` object with `id = <key>` |
| Component `default` / current state | Component `value`; never emit a `default` field |
| Other component settings | Flat component fields (`min`, `max`, `mapping`, `unit`, etc.) |
| Action item TOML `id` | Wire item `action_id`; do not emit item `id` |

Validate the resulting snapshot against the maintained contract before publication:
all layout references must exist, component IDs must match keys, and action IDs
must be unique. Unreferenced components are allowed by the core; the SDK may warn.
The current configuration schema supports the six component types in the protocol,
including `link` with flat `uri` and `text`. No standalone button component exists.

Runtime updates and heartbeats publish the full snapshot to
`synapse/v1/discovery/{id}` or POST `/api/v1/discovery`. MQTT topics and HTTP routes
retain `v1`. Definition fields may change between snapshots; partial updates are
unsupported. Log-stream semantics must follow their dedicated task before SDK
implementation assumes a merge policy.
