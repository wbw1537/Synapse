---
id: SDK-004
type: task
status: Pending
---

# SDK Task: Python SDK Implementation

## Status: Pending
## Effort Estimate: High

## Description
Build the reference Python SDK (`synapse-axon`) that implements the requirements defined in `docs/sdk_specification.md`.

## Requirements

### 1. Project Structure
-   Create `sdk/python/` directory.
-   `pyproject.toml` (Poetry or setuptools).
-   Dependencies: `paho-mqtt`, `toml` (or `tomli` for Py<3.11).

### 2. Core Modules
-   **`config.py`**:
    -   `load_from_file(path)`
    -   **Validation Logic**: Implements the "Validation Rules Matrix" from the spec.
-   **`client.py`**:
    -   `Axon` class (Facade).
    -   MQTT Connection management.
    -   Heartbeat loop (Threaded).
-   **`components.py`**:
    -   Wrapper classes for `stat`, `gauge`, etc., to handle type-safe updates.

### 3. Usage & Examples
-   Create `examples/sdk_demo.py` demonstrating:
    -   Loading `axon.toml`.
    -   Registering an action callback.
    -   Updating a metric in a loop.

## Tasks
- [ ] Setup Python project structure in `sdk/python`.
- [ ] Implement TOML Parser & Validator.
- [ ] Implement Component State Manager.
- [ ] Implement MQTT Client & Heartbeat.
- [ ] Implement Action Dispatcher.
- [ ] Create `examples/sdk_demo.py`.

## Resources
-   `docs/sdk_specification.md`: Generic SDK requirements.
-   `docs/axon_toml_spec.md`: Configuration schema.