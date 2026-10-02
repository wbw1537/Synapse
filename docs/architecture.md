# Current architecture

Baseline: source at `36f9ac8`, inspected 2026-10-02. This is a source-level map,
not evidence that the deployed system or every runtime path has been verified.

| Path | Responsibility |
| --- | --- |
| `cmd/synapse/main.go` | Configuration, DB/schema, HTTP, embedded broker, internal MQTT client, shutdown |
| `internal/models/service.go` | Service, section layout, component map, monitor definitions |
| `internal/service/manager.go` | Token validation, snapshot upsert, log merge, TTL, expressions, command publication |
| `internal/api/server.go` | Service reads, HTTP discovery, action requests, embedded frontend |
| `internal/broker/broker.go` | TCP/WS MQTT listeners; currently allows all clients |
| `internal/db/sqlite.go` | SQLite and GORM schema initialization |
| `internal/evaluator/` | Server-side expressions against component `value` |
| `internal/notification/` | In-memory alert transitions and SMTP delivery |
| `web/src/stores/services.ts` | Initial HTTP snapshot, raw MQTT discovery subscription, UI state and actions |
| `web/src/components/` | Section/component cards, detail drawer, runbooks and widgets |

## Data paths

1. An Axon publishes discovery; the internal MQTT client passes it to `UpsertMQTT`, checking topic identity.
2. HTTP discovery calls the same manager. The manager validates the v1 snapshot, IDs/references and token, merges
   log state, writes SQLite, then evaluates monitors.
3. The UI fetches persisted services once and separately consumes raw MQTT
   discovery messages. These paths are not a canonical server event stream.
4. TTL updates SQLite every ten seconds; it currently does not notify the UI.
5. An action HTTP request checks the stored capability and publishes an action ID
   through MQTT. HTTP success means publication succeeded, not execution completed.

The backend and browser both merge log streams, with differences in behavior.
The browser sees raw payloads before server validation and does not automatically
receive HTTP registration or TTL changes. These are tracked gaps, not intended
architectural invariants.

The frontend is built before Go compilation because `ui.go` embeds `web/dist/*`.
There is one Git repository and no nested application repository.
