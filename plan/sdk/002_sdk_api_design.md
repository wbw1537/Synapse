---
id: SDK-002
type: task
status: Done
---

# SDK Task: SDK API & Lifecycle Design

## Status: Done
## Effort Estimate: Medium

## Description
Design the public API for the Synapse Axon SDK. The SDK will be responsible for loading the `axon.toml` configuration, validating it, connecting to the MQTT broker, and managing the runtime state of components.

**Goal:** define the "Developer Experience" (DX) for writing an Axon.

## Requirements

### 1. SDK Lifecycle
Define the behavior for each phase:
1.  **Initialization**:
    -   Loading `axon.toml`.
    -   Validation (Schema check, Orphan/Ghost check).
2.  **Registration**:
    -   Connecting to MQTT.
    -   Constructing and publishing the initial "Discovery" payload (merging TOML static data with default values).
3.  **Runtime Loop**:
    -   Heartbeat management (automatic re-publishing).
    -   Action handling (callbacks).

### 2. API Surface (Python/Go)
Design the classes and methods users will interact with.

*   **Constructor**: How to load the config? `app = Axon("axon.toml")`?
*   **State Updates**: How to update a widget?
    -   `app.components["cpu"].update(50)`?
    -   `app.set_value("cpu", 50)`?
*   **Action Handling**: Decorators vs Callbacks.
    -   `@app.on_action("restart")`?

### 3. Safety & Validation rules
-   Define what checks the SDK *must* perform before connecting.
    -   e.g. "Do not allow `update("unknown_id")`".
    -   e.g. "Validate value type against widget type (no string for gauge)".

## Deliverables
1.  **Document**: `docs/sdk_design_python.md` (or general SDK spec).
    -   Pseudo-code for the main `Axon` class.
    -   Example user code (how the `memory_axon.py` would look using the SDK).
    -   Sequence diagram of the startup/registration flow.
2.  **Validation Matrix**:
    -   A list of validation rules the SDK must enforce (e.g. "Warning if defined component is unused in layout").

## Resources
-   `plan/sdk/axon_toml_and_sdk.md`: High-level design.
-   `examples/memory_axon.py`: Current "manual" implementation.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-31
- **Key Artifacts**:
    - `docs/sdk_design_python.md`: Python SDK API surface, lifecycle state machine, and validation rules.
- **Design Decisions**:
    - **Fail Fast**: The SDK will perform rigorous validation of the TOML config upon initialization, raising errors for orphans, ghosts, or type mismatches before connecting.
    - **Decorator Pattern**: Actions will be handled via `@app.on_action("id")` decorators for clean, idiomatic Python code.
