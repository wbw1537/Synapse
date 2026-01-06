# Synapse Widget Reference

Widgets are the UI building blocks used by Axons to report data and provide interactivity on the Synapse Dashboard. These are defined in the `widgets` array within the Axon configuration payload.

## 1. Common Properties

All widget types share the following base properties:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | `string` | **Yes** | Unique identifier for the widget within the Axon. |
| `type` | `string` | **Yes** | The component type (e.g., `stat`, `gauge`). |
| `label` | `string` | No | The display title of the widget. |
| `visible` | `boolean` | No | Default `true`. Set to `false` for "headless" monitoring. |
| `monitors` | `array` | No | Array of alert rules (see *Server-Side Monitoring*). |

---

## 2. Widget Types

### 2.1 `stat`

The most versatile widget for displaying key-value pairs (text or numbers).

**Specific Properties:**

* `value` (any): The data to display.
* `unit` (string, optional): Suffix unit (e.g., "MB", "%", "°C").
* `copyable` (boolean, optional): If `true`, clicking the value copies it to the clipboard.

**JSON Sample:**

```json
{
  "type": "stat",
  "id": "wan_ip",
  "label": "Public IP",
  "value": "203.0.113.15",
  "copyable": true,
  "monitors": [
    { 
      "condition": "value != '203.0.113.15'", 
      "severity": "info", 
      "message": "WAN IP Changed" 
    }
  ]
}

```

---

### 2.2 `status_indicator`

Used for boolean or enum states (Online/Offline, Healthy/Degraded). It maps raw values to visual states (colors, icons).

**Specific Properties:**

* `value` (string | boolean): The raw state key.
* `mapping` (object): A dictionary mapping raw values to UI styles.
* `text` (string): Display text.
* `color` (string): `success` (green), `warning` (yellow), `danger` (red), `primary` (blue), `neutral` (gray).
* `icon` (string, optional): MDI icon name (e.g., `check`, `alert`).
* `animate` (boolean, optional): If `true`, adds a pulse effect.



**JSON Sample:**

```json
{
  "type": "status_indicator",
  "id": "login_state",
  "label": "Bot Status",
  "value": "connected",
  "mapping": {
    "connected": {
      "text": "Online",
      "color": "success",
      "icon": "wifi",
      "animate": true
    },
    "disconnected": {
      "text": "Offline",
      "color": "neutral",
      "icon": "wifi-off"
    },
    "banned": {
      "text": "Account Banned",
      "color": "danger",
      "icon": "gavel"
    }
  }
}

```

---

### 2.3 `gauge`

Visualizes a numerical metric against a min/max scale.

**Specific Properties:**

* `value` (number): Current numerical value.
* `min` (number): Minimum value (default 0).
* `max` (number): Maximum value (default 100).
* `unit` (string): Unit label.
* `thresholds` (object, optional): Map of value thresholds to colors.

**JSON Sample:**

```json
{
  "type": "gauge",
  "id": "memory_usage",
  "label": "RAM Usage",
  "value": 82,
  "unit": "%",
  "min": 0,
  "max": 100,
  "thresholds": {
    "70": "warning",
    "90": "danger"
  }
}

```

---

### 2.4 `log_stream`

Displays a chronological list of events. Unlike other widgets that replace the value, sending a value to this widget **appends** it to the stream.

**Specific Properties:**

* `value` (string): The latest log line to append.
* `max_items` (number): Maximum number of lines to keep in history (default 10).

**JSON Sample:**

```json
{
  "type": "log_stream",
  "id": "backup_log",
  "label": "Backup History",
  "value": "[2026-01-04 10:00] Backup successful (4.2GB)",
  "max_items": 5
}

```

---

### 2.5 `action_group`

A container for buttons that trigger remote actions defined in the Axon's `actions` array.

**Specific Properties:**

* `items` (array): List of buttons to render.
* `action_id` (string): Must match an ID in the Axon's top-level `actions` list.
* `label` (string): Button text.
* `style` (string): `primary`, `danger`, `default`.
* `confirm` (boolean): If `true`, requires user confirmation before execution.



**JSON Sample:**

```json
{
  "type": "action_group",
  "id": "service_controls",
  "label": "Operations",
  "items": [
    {
      "action_id": "reload_nginx",
      "label": "Reload Config",
      "style": "primary",
      "confirm": false
    },
    {
      "action_id": "stop_nginx",
      "label": "Stop Service",
      "style": "danger",
      "confirm": true
    }
  ]
}

```

---

---

### 2.6 `link`

A static navigation shortcut, useful for linking to external Web UIs.

**Specific Properties:**

* `uri` (string): The destination URL.
* `text` (string, optional): Link text (defaults to "Open").
* `icon` (string, optional): MDI icon name (e.g., `open-in-new`).

**JSON Sample:**

```json
{
  "type": "link",
  "id": "webui_link",
  "label": "Management Portal",
  "uri": "http://192.168.1.50:8123",
  "text": "Open Dashboard",
  "icon": "view-dashboard"
}

```

---

## 3. Server-Side Monitoring (Monitors)

Monitors allow Synapse to evaluate incoming data and trigger notifications (e.g., email via SMTP) when specific conditions are met. Notifications are only sent when the state changes (e.g., from `ok` to `error`).

### 3.1 Monitor Object Structure

| Field | Type | Description |
| --- | --- | --- |
| `condition` | `string` | An expression evaluated against the widget's `value`. |
| `severity` | `string` | `warning` or `critical` (used in notification subject). |
| `message` | `string` | The alert message to include in the notification. |

### 3.2 Writing Conditions by Widget Type

The `condition` string uses a Go-based expression engine where the current widget data is available as the variable `value`.

#### Numeric Widgets (`stat`, `gauge`)
`value` is a number. Use standard comparison operators (`>`, `<`, `==`, `!=`).
```json
"monitors": [
  { "condition": "value > 90", "severity": "critical", "message": "High usage!" }
]
```

#### State-based Widgets (`status_indicator`)
`value` is a string or boolean. Wrap strings in single quotes.
```json
"monitors": [
  { "condition": "value == 'offline'", "severity": "critical", "message": "Service is down!" }
]
```

#### List-based Widgets (`log_stream`)
`value` is an **array of strings** (the log history). 
```json
"monitors": [
  { 
    "condition": "len(value) > 0 && value[len(value)-1] contains 'ERROR'", 
    "severity": "critical", 
    "message": "Error detected in logs!" 
  }
]
```
*Note: For `log_stream`, it's often more reliable to use a separate hidden `stat` widget for specific error flags.*