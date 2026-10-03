# Development constraints

These are design boundaries to preserve. They do not imply the existing Alpha
implements every security or reliability requirement.

- Keep the core deployable with its embedded broker, SQLite, and frontend.
- Axons own measurement and executable action handlers. The core publishes IDs
  of declared actions; it must not send arbitrary executable code.
- HTTP and MQTT registration share business logic. Define wire contracts before
  changing either path; update backend, UI, examples, and documentation together.
- Treat layout/component integrity, version handling, and field shapes as public
  integration behavior. Breaking changes need a documented compatibility decision.
- Monitor evaluation stays server-side unless a deliberate design change says otherwise.
- Preserve bounded log history and distinguish snapshots from appended events.
- Separate runtime validation data from the user's existing database and services.
- Credentials must not appear in docs, task evidence, logs added by development,
  or committed configuration. Use example configuration for variable names.
- Record observable verification; build success alone does not establish end-to-end correctness.

Runtime access follows [the homelab policy](access.md): distinct operator/Axon
credentials, authenticated HTTP reads/actions, MQTT authentication and topic ACLs,
and loopback defaults. Axons share a credential and therefore a trust realm;
per-service identities, roles and native TLS remain future work.
