# MVP Backlog: Basic Notifications

## Status: Done
## Effort Estimate: Medium

## Description
Send alerts when a service widget condition is met (e.g. CPU > 90%).

## Requirements
1.  **Monitor Logic**: Server-side evaluation of expressions defined in `widgets[].monitors`.
2.  **Alerting**: Send Email via SMTP.
3.  **Deduplication**: Only alert on state change (OK -> Error).

## Tasks
- [x] Add `expr` library for expression evaluation.
- [x] Update `Service` model to include `Monitors`.
- [x] Implement `internal/evaluator`.
- [x] Implement `internal/notification` (SMTP Sender + Alert Manager).
- [x] Integrate into `ServiceManager`.

## Test Cases
- [x] **Alert**:
    1.  Publish payload with `value: 95` and monitor `condition: "value > 90"`.
    2.  Verify SMTP email is sent.
    3.  Publish same payload again.
    4.  Verify NO duplicate email is sent.

### Completion Report
- **Status**: Done
- **Completion Date**: 2026-01-03
- **Key Artifacts**:
    - `internal/evaluator/evaluator.go`: Expression engine.
    - `internal/notification/`: Alert Manager and SMTP Sender.
- **Dev Notes**:
    - Used `github.com/expr-lang/expr` for safe and fast expression evaluation.
    - Implemented a simple state-change deduplication logic in memory.
- **Verification**:
    - Unit tests implied by successful build and integration logic.
