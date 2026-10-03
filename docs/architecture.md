# Current architecture

Baseline: working tree based on `f6d0daa`, updated 2026-10-03. This is a source-level map,
not evidence that the deployed system or every runtime path has been verified.

| Path | Responsibility |
| --- | --- |
| `cmd/synapse/main.go` | Configuration, DB/schema, HTTP, embedded broker, internal MQTT client, shutdown |
| `internal/models/service.go` | Service, section layout, component map, monitor definitions |
| `internal/service/manager.go` | Token validation, snapshot upsert, log merge, TTL, expressions, command publication |
| `internal/api/server.go` | Service reads, HTTP discovery, action requests, embedded frontend |
| `internal/broker/broker.go` | TCP/WS MQTT listeners; authenticated core/Axon clients with topic ACLs |
| `internal/db/sqlite.go` | SQLite and GORM schema initialization |
| `internal/evaluator/` | Server-side expressions against component `value` |
| `internal/notification/` | In-memory alert transitions and SMTP delivery |
| `web/src/stores/services.ts` | Authoritative SSE snapshots, UI state and actions |
| `web/src/components/` | Section/component cards, detail drawer, runbooks and widgets |

The HTTP listener opens after initial MQTT discovery subscription and command
publisher initialization, so API readiness includes ingestion and dispatch.

## Data paths

1. An Axon publishes discovery; the internal MQTT client passes it to `UpsertMQTT`, checking topic identity.
2. HTTP discovery calls the same manager. The manager validates the v1 snapshot, IDs/references and token, merges
   log state, writes SQLite, then evaluates monitors.
3. The UI opens same-origin `GET /api/v1/events`. The server subscribes to change
   notifications before reading SQLite and sends complete service snapshots.
4. Accepted registration and TTL updates invalidate streams; notifications are
   coalesced per subscriber. A slow browser does not block mutations. Reconnect
   starts with a fresh snapshot; periodic snapshots keep idle connections active.
5. An action HTTP request checks the stored capability and publishes an action ID
   through MQTT. HTTP success means publication succeeded, not execution completed.

Log normalization is server-owned: string events append, array snapshots replace,
and retention applies from first registration. The UI replaces state from server
snapshots without merging raw events. Mutations are serialized so concurrent
read-modify-write updates cannot lose log entries. Browser rendering and runtime lifecycle acceptance passed SYN-102/SYN-103.
Operator sessions protect reads/actions/state streams; broker access and listener
defaults follow [the access policy](access.md). The Python reference SDK performs
TOML validation, typed updates and client lifecycle management.

The frontend is built before Go compilation because `ui.go` embeds `web/dist/*`.
There is one Git repository and no nested application repository.
