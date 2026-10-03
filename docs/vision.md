# Product intent

Synapse puts a homelab service's state, operational controls, and recovery
instructions in one place. Services describe themselves through Axons, avoiding
manual dashboard entries for every new service.

A useful increment improves a real operator workflow: integrate a service,
understand a fault, inspect supporting data, or invoke a predefined recovery action.
Keep the core self-hosted, lightweight, and easy to deploy as a single binary.

The existing MVP includes discovery, component cards, runbooks, expressions,
SMTP alerts, and command publication. The Python reference SDK and operator-token login support integrations and
controlled operations. Policy overrides, recovery code creation, multi-user login
and additional notification channels remain future capabilities.
They must not be described as implemented without evidence.

`DESIGN.md` preserves the original broader product proposal. This document
summarizes intent; current behavior is recorded separately.
