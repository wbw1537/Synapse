---
id: MVP-008
type: story
status: Done
---

# MVP Backlog: Action Triggering

## Status: Done
## Effort Estimate: Medium

## Description
Implement the mechanism for users to trigger actions defined by Axons. This involves an HTTP endpoint that publishes a command to the MQTT broker, which the Axon then consumes to execute the local function.

## Requirements
1.  **Backend API**:
    *   `POST /api/v1/services/{service_id}/actions/{action_id}`: Endpoint to trigger an action.
    *   **Validation**: Ensure the `service_id` exists and the `action_id` is valid for that service (though the Axon is the ultimate authority, the backend should do a basic check against the registered actions).
2.  **MQTT Command Topic**:
    *   Publish to `synapse/v1/command/{service_id}`.
    *   Payload: `{"action_id": "restart", "issued_by": "user"}` (or similar).
3.  **Frontend Integration**:
    *   Update `ActionGroupWidget.vue` and `ServiceDrawer.vue` to call the API instead of `alert()`.
    *   Handle success/failure notifications (mocked or simple toast for now).
4.  **Axon Implementation (Python)**:
    *   Update `memory_axon.py` to subscribe to `synapse/v1/command/+`.
    *   Parse the payload and execute the corresponding logic (mocked logic is fine).
    *   Validate: If an Axon receives a command for an action it didn't register (or doesn't support), it should log a warning.

## Tasks
- [ ] **Backend**: Implement `ExecuteAction` in `ServiceManager` (publish to MQTT).
- [ ] **Backend**: Add API route `POST /api/v1/services/{id}/actions/{action_id}`.
- [ ] **Frontend**: Implement `executeAction` in `useServiceStore`.
- [ ] **Frontend**: Connect UI buttons to the store action.
- [ ] **Axon**: Implement MQTT subscription and command handling loop in `memory_axon.py`.

## Test Cases
- [x] **Trigger**: Click "Restart" button on UI.
- [x] **Network**: Verify backend publishes to `synapse/v1/command/memory-sidecar-python`.
- [x] **Axon**: Verify `memory_axon.py` prints "Executing action: restart".
- [x] **Validation**: Try to trigger a non-existent action ID via curl; should return 400 or 404.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-04
- **Key Artifacts**:
    - `internal/api/server.go`: Added `executeAction` endpoint.
    - `internal/service/manager.go`: Added `ExecuteAction` logic with MQTT publishing.
    - `web/src/stores/services.ts`: Added frontend action handler.
    - `examples/memory_axon.py`: Added command subscription and handling.
- **Dev Notes**:
    - Implemented a secure action triggering flow where the backend validates the action ID before publishing the MQTT command.
    - The Axon now listens for commands and executes logic locally.
- **Verification**:
    - Verified full flow: UI Button -> API -> MQTT -> Python Script.

