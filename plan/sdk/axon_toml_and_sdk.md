# Synapse Axon Design Document

---

## 1. Executive Summary

This document defines the architecture for the **Synapse Axon SDK (v1)**. The goal is to standardize how clients (Axons) define their UI and report data to the Synapse Core.

The design adopts a **"Definition vs. Runtime" separation**:

* **Definition (TOML):** The UI layout, component types, static labels, and monitoring rules are defined in a strict configuration file (`axon.toml`).
* **Runtime (SDK):** The business logic resides in code (Python/Go/etc.), which binds to the configuration via IDs.

This approach eliminates the need for complex code generation while ensuring type safety through rigorous startup validation ("Fail Fast").

---

## 2. Architecture Decisions (ADR)

### 2.1 Configuration Format: TOML

* **Decision:** Use TOML for configuration.
* **Rationale:** Superior readability for configuration files compared to JSON; standard in the target ecosystem (Rust/Go); supports comments.
* **Constraint:** SDKs must implement a compliant TOML parser with schema validation.

### 2.2 Wire Protocol: Split Layout/Data

* **Decision:** The MQTT payload will separate `layout` (presentation tree) from `components` (state dictionary).
* **Rationale:** Enables complex UI structures (grids, tabs, sections) without bloating the component state definitions. Matches modern declarative UI frameworks.

### 2.3 Runtime Model: Validated Runtime SDK

* **Decision:** Use a library-based approach (Bean Lifecycle) rather than code generation.
* **Rationale:** Lowers barrier to entry. User writes a config file + standard code. Safety is guaranteed by an aggressive validation phase at SDK startup (`init`), not at compile time.

---

## 3. Specification

### 3.1 The Configuration Schema (`axon.toml`)

The source of truth for the Axon.

```toml
# MUST be present. Defines parser rules.
schema = "axon.card.v1"

[meta]
id = "ddns-updater"
name = "DDNS Updater"
icon = "cloud-sync"
ttl = 30

# The Layout Tree. Defines ORDER and HIERARCHY.
[layout]
type = "sections"  # Future: "grid", "tabs"

  [[layout.section]]
  title = "Network Status"
  # References IDs defined in [components]
  components = ["current_ip", "sync_status"]

  [[layout.section]]
  title = "Logs"
  components = ["sync_log"]

# The Component Registry. Defines TYPE and BEHAVIOR.
# Keys must match IDs used in [layout].

[components.current_ip]
type = "stat"
label = "Current IP"
default = "Initializing..."

[components.sync_status]
type = "status_indicator"
label = "Status"
default = "idle"

  # Structured Mapping
  [components.sync_status.mapping]
  idle    = { text = "Idle", color = "neutral" }
  syncing = { text = "Syncing", color = "primary", animate = true }
  error   = { text = "Error", color = "danger" }

[components.error_count]
type = "stat"
label = "Errors"
unit = "errs"

  # Namespaced Monitors
  [[components.error_count.monitors]]
  type = "threshold"
  condition = "value >= 10"  # Simple expression DSL
  severity = "critical"
  message = "High error rate detected"

[components.sync_now]
type = "action"
label = "Sync Now"
style = "primary"

```

### 3.2 The Wire Protocol (JSON Payload)

The compiled output sent to Synapse Core via MQTT.

```json
{
  "spec_version": "axon.card.v1",
  "meta": { ... },
  "layout": {
    "type": "sections",
    "root": [
      {
        "type": "section",
        "title": "Network Status",
        "children": ["current_ip", "sync_status"]
      }
    ]
  },
  "components": {
    "current_ip": {
      "type": "stat",
      "label": "Current IP",
      "value": "Initializing..."
    },
    "sync_status": {
      "type": "status_indicator",
      "label": "Status",
      "value": "idle",
      "mapping": { ... }
    }
  }
}

```

---

## 4. The SDK Lifecycle Requirements

The SDK must implement a strict state machine to ensure safety.

### Phase 1: Initialization (The Validator)

*Trigger:* `sdk = SynapseAxon("axon.toml")`

1. **File Load:** Read TOML file.
2. **Schema Check:** Verify `schema` version string is supported.
3. **Integrity Validation:**
* **Orphan Check:** Ensure every ID in `[layout]` exists in `[components]`.
* **Ghost Check:** Warn if an ID in `[components]` is not used in `[layout]`.


4. **Type Validation:** iterate through all components and validate required fields per type (e.g., `status_indicator` *must* have `mapping`).
5. **Monitor DSL Check:** Parse `condition` strings to ensure no unsafe syntax (e.g., ensure no `exec()` or function calls).

### Phase 2: Registration (The Handshake)

*Trigger:* `sdk.start()`

1. **Connect:** Establish MQTT connection.
2. **Compile:** Convert validated internal state to Wire Protocol JSON.
3. **Publish:** Send to `synapse/v1/discovery/{id}`.
4. **Subscribe:** Listen to `synapse/v1/command/{id}`.
5. **Heartbeat Loop:** Start background thread to republish payload every `TTL / 2` seconds.

### Phase 3: Operation (The Runtime)

*Trigger:* `sdk.components.update("id", value)`

1. **Lookup:** Check if "id" exists in the internal component map.
* *If missing:* Throw `ComponentNotFoundError`.


2. **Type Check (Optional but recommended):** Check if `value` matches component expectations (e.g., don't send a String to a Gauge).
3. **State Update:** Update internal state.
4. **Flush:** Push update to MQTT.

---

## 5. Mandatory Tasks (Implementation Plan)

### Task Group A: Synapse Core (The Server)

* [ ] **A1. Protocol Upgrade:** Update MQTT handler to parse the new `axon.card.v1` schema (splitting `layout` and `components`).
* [ ] **A2. Rename Existing Widget to Component**
* [ ] **A3. Layout Renderer:** Update Vue frontend to render based on the `layout` tree object instead of a flat list.
* [ ] **A4. Component Registry:** Refactor frontend component loader to handle the new explicit mapping structure.

### Task Group B: SDK Core (The Library)

* [ ] **B1. TOML Parser:** Implement `ConfigLoader` class that reads TOML.
* [ ] **B2. Validation Engine:** Implement `Validator` class with:
* `validate_schema_version()`
* `validate_integrity()` (Orphans/Ghosts)
* `validate_component_props()`


* [ ] **B3. Safe Evaluator:** Implement a tiny parser for monitor conditions (e.g., using Python `ast` or regex) to validate `value >= 10`. **Do not use `eval()`.**

### Task Group C: Runtime & Glue

* [ ] **C1. State Manager:** Implement the internal dictionary to hold component states.
* [ ] **C2. MQTT Client Wrapper:** Encapsulate Paho-MQTT (or equivalent) logic, handling auto-reconnect and heartbeats.
* [ ] **C3. Action Dispatcher:** Implement the `@on_action("id")` decorator pattern to route incoming MQTT commands to user functions.

---

## 6. Future Considerations (Not in v1)

* **Grid Layouts:** Adding `row`/`col` definitions to the layout block.
* **Input Components:** Text inputs or sliders for user interaction (beyond simple buttons).
* **Multi-Page:** Defining `[[layout.tab]]` for complex Axons.