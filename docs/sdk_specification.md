# Synapse Axon SDK Specification

This document defines the standard behavior and architecture for Synapse Axon SDKs, regardless of the implementation language (Python, Go, Rust, etc.).

All official SDKs must adhere to these requirements to ensure consistent behavior across the ecosystem.

## 1. Design Principles
*   **Fail Fast**: Configuration errors (invalid TOML, schema violations) must cause the application to crash/exit *immediately* during initialization, before network connections are attempted.
*   **Type Safety**: The SDK must validate data types at runtime (e.g., ensuring a string isn't sent to a Gauge widget).
*   **Definition vs. State**: The SDK loads the static definition from `axon.toml` and manages the dynamic state in memory.

## 2. SDK Lifecycle

An Axon SDK must implement the following state machine:

### Phase 1: Initialization
*Input*: Path to `axon.toml`.

1.  **Load File**: Read and parse the TOML configuration.
2.  **Validate Schema**:
    *   Ensure `schema` matches the supported version.
    *   **Integrity Check**:
        *   **Orphans**: Warn or Error if a component is defined in `[components]` but not used in `[layout]`.
        *   **Ghosts**: **MUST Error** if an ID is used in `[layout]` but not defined in `[components]`.
    *   **Type Check**: Validate required properties for each component type (e.g., `gauge` needs `min`/`max`).
3.  **Initialize State**: Create an in-memory registry of components with their `default` values.

### Phase 2: Registration
*Trigger*: User calls `start()`.

1.  **Connect**: Establish connection to the MQTT Broker.
2.  **Publish Discovery**: Serialize the configuration + current state into the JSON Wire Protocol and publish to `synapse/v1/discovery/{id}`.
3.  **Start Heartbeat**: Spawn a background task to re-publish the Discovery payload every `TTL / 2` seconds.

### Phase 3: Runtime Loop
*Trigger*: Application Logic.

1.  **State Updates**:
    *   User updates a component (e.g., `components["cpu"].set(50)`).
    *   SDK updates internal state.
    *   SDK publishes the updated payload to MQTT immediately (or debounced).
2.  **Command Handling**:
    *   Subscribe to `synapse/v1/command/{id}`.
    *   On message: Parse `action_id`.
    *   Invoke the registered callback/handler for that action.

---

## 3. Standard API Surface

While syntax varies by language, the concepts should map 1:1.

### 3.1 Constructor
Should accept the config path.
*   Python: `Axon("axon.toml")`
*   Go: `NewAxon("axon.toml")`

### 3.2 Component Access
Users should access components by ID to update them.
*   Python: `axon.components["cpu"].update(50)`
*   Go: `axon.Component("cpu").Update(50)`

### 3.3 Action Binding
A mechanism to bind code to `action_group` buttons.
*   Python: Decorator `@axon.on_action("restart")`
*   Go: `axon.OnAction("restart", func() { ... })`

---

## 4. Validation Rules Matrix

| Rule | Severity | Description |
| :--- | :--- | :--- |
| **File Missing** | **Fatal** | Config file does not exist. |
| **Invalid TOML** | **Fatal** | Syntax error in file. |
| **Ghost ID** | **Fatal** | Layout references ID not in Components. |
| **Missing Props** | **Fatal** | Component missing mandatory fields (e.g. `gauge` without `max`). |
| **Orphan ID** | Warning | Component defined but not used in Layout. |
| **Type Mismatch** | Runtime Error | Sending `string` to numeric widget. |

---

## 5. Security Recommendations

*   **Sanitize Inputs**: SDKs should assume MQTT payloads are untrusted.
*   **Monitor Safety**: If the SDK evaluates local monitoring rules, avoid using unsafe `eval()` functions. Use parsed expression engines.