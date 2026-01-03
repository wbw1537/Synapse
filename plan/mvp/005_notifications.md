# MVP Backlog: Basic Notifications

## Status: Pending
## Effort Estimate: Low

## Description
Simple integration to send alerts when a service goes down or recovers.

## Requirements
1.  **Library**: Use `github.com/containrrr/shoutrrr`.
2.  **Configuration**: Load notification URLs from environment variables or config file (e.g., `SYNAPSE_NOTIFICATIONS=discord://...,telegram://...`).
3.  **Trigger**: Hook into the `ServiceManager` status change logic.

## Tasks
- [ ] Implement `internal/notification/sender.go`.
- [ ] Add config parsing for notification URLs.
- [ ] Call notification sender when Status transitions (Online -> Offline, Offline -> Online).

## Test Cases
- [ ] **Alert**:
    1.  Configure a mock Shoutrrr URL (or a real Discord webhook).
    2.  Force a service to timeout (TTL).
    3.  Verify message is received in channel.
