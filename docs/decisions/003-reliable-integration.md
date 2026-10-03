# Reliable integration Alpha

Accepted under the autonomous next-stage implementation request, 2026-10-03.
Baseline is f6d0daa; work remains in the working tree for human review.

## Decisions

- Browser state comes from same-origin SSE snapshots of persisted, validated
  services. Coalesced change notifications plus an initial snapshot prevent startup
  gaps and slow-client blocking; reconnect reconciles missed changes. Full snapshots
  suit a small homelab and avoid maintaining a second browser log merge engine.
- Log normalization belongs to the core. Strings append; arrays replace; null and
  empty strings preserve; empty arrays clear. Retention always applies. The SDK
  publishes bounded arrays so heartbeats cannot append the same event repeatedly.
- Operator and Axon secrets are separate. Operators use bearer authentication or
  HttpOnly sessions; authenticated Axons have service-topic ACLs. The internal core
  has a startup-generated credential. Axons sharing a token are one trust realm,
  deliberately not multi-tenant/per-service credential isolation.
- Native/dev listeners and Compose host ports default to loopback. Intentional
  remote access requires private networking or TLS protection; native TLS is not
  introduced. No default secret remains. Existing deployments must configure two
  secrets and update MQTT clients to service IDs and Axon authentication.
- The reference SDK is Python 3.11+, TOML plus paho-mqtt 2.x, with full snapshots,
  typed updates and separate bounded callback dispatch. Other languages remain
  future work. MQTT success means delivery/subscription, not business acceptance
  or exactly-once execution.
- Branch/PR/tag verification precedes any version-tag-only publication. Local
  acceptance uses isolated processes, generated credentials and temporary SQLite.

## Evidence and limitations

See [acceptance report](../acceptance.md) for commands and actual runtime outcomes.
Legacy database capability rows are not migrated. Per-Axon keys, users/roles,
command-result history, durable logs and native TLS are future increments.
Independent agent review was unavailable under the active delegation policy;
solo source review, failure-path checks and actual runtime acceptance were used.
