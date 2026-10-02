---
name: operation
description: Diagnose Synapse local or deployed runtime, Docker Compose and tag-triggered CI, and carry out explicitly authorized deployment or recovery. Use for status/log inspection, runtime failures, releases or operational changes.
---

# Synapse Operation

Read [operations](../../docs/operation.md). Identify the actual target environment;
this repository does not define Beanpilot or other projects' hosts/runbooks.

Inspect status, listener availability, Compose state, relevant logs, selected image,
Git revision and DB volume before mutation. Never print secret files or
credential-bearing discovery payloads. Explain whether a fault is build, startup,
transport, persistence, UI state or external notification delivery.

Use existing Docker/CI configuration. Deployment uses the current working tree,
and pushing a `v*` tag triggers artifact publication. Neither is implied by
finishing local implementation or review. Honor deployment/restart authorization
already present in the session; do not repeat approval rituals.

Before destructive migration/recovery, establish data ownership, backup and rollback
approach. Old Alpha advice to reset the database is not authorization to delete
current data. After an authorized change verify both HTTP and MQTT paths and the
operator-visible result. Report the environment, revision and observed outcome.
