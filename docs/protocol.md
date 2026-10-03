# Discovery protocol

This is the maintained wire contract for Synapse discovery, implemented by
`internal/models/discovery.go` and the shared service manager. The layout/component
refactor was historically called “v2”; the supported `api_version` is **exactly
`"v1"`**. Missing, null, numeric and all other version values are rejected.

## Transports and failure behavior

| Operation | Path |
| --- | --- |
| HTTP full snapshot | `POST /api/v1/discovery` |
| MQTT full snapshot | `synapse/v1/discovery/{id}` |
| Read stored services | `GET /api/v1/services` or `GET /api/v1/services/{id}` |
| Request declared action | `POST /api/v1/services/{id}/actions/{action_id}` |
| Axon command subscription | `synapse/v1/command/{id}` |

HTTP discovery authenticates with the payload token. MQTT additionally requires
username `axon`, Axon token password and a service-matching client ID, with topic
ACLs. Reads, actions and SSE require operator authentication; see [access policy](access.md).

HTTP and MQTT call the same decoder, validator and persistence logic. MQTT also
requires the topic suffix to equal the top-level service `id`, with no extra path
segments. HTTP success is `200 OK` with `OK`; invalid registration is `400` with
an explanatory text body. MQTT invalid registration is logged by the core and
leaves storage unchanged. There is no application acknowledgment topic; MQTT
PUBACK acknowledges delivery, not successful registration. Verify storage via HTTP.

A snapshot is one JSON object. Legacy top-level `widgets` or `actions` fields are
rejected even if empty/null or mixed with valid components. Unknown fields,
including nested `meta`, component `props`, and action item `id`, are rejected.
They are never translated automatically. This intentionally breaks legacy clients
that previously appeared registered while losing their capabilities. Migrate them
to the shape below. Existing database rows are not migrated by this task.

## Snapshot shape

Identity and component properties are **flat**. Required fields are `api_version`,
`auth_token` (matching the configured token), `id`, `status`, positive integer
`ttl`, `layout`, and `components`. Optional service fields are `name`, `group`,
`tags` (string array), `icon`, `url`, `message`, `description`, `markdown_docs`.
Model metadata `last_seen`, `created_at`, `updated_at` can be decoded; clients
should omit them and let the core/GORM supply timestamps. `last_seen` is set at every accepted
registration. IDs use ASCII letters, digits, `_`, `-` and must be nonempty.
Status is `online`, `warning`, `error` or `offline`.

```json
{
  "api_version": "v1",
  "auth_token": "<configured-token>",
  "id": "demo-service",
  "name": "Demo Service",
  "group": "Examples",
  "tags": [],
  "status": "online",
  "ttl": 30,
  "layout": {
    "type": "sections",
    "root": [{"type": "section", "title": "Health", "children": ["cpu", "controls"]}]
  },
  "components": {
    "cpu": {"id": "cpu", "type": "gauge", "label": "CPU", "value": 25,
            "min": 0, "max": 100, "unit": "%", "monitors": []},
    "controls": {"id": "controls", "type": "action_group", "label": "Controls",
                 "items": [{"action_id": "restart", "label": "Restart", "style": "danger", "confirm": true}],
                 "monitors": []}
  }
}
```

`layout.type` must be `sections`, `root` an array, each section's `type` must be
`section`, and `children` an array of component IDs. Every child must exist in
`components`. `components` is an object keyed by ID, and each entry must contain
an identical `id`. Unsupported component types are rejected. Empty `root`,
`children`, component maps and action item arrays are allowed; missing/null
containers are rejected. Unreferenced components are allowed and retained (the
drawer can still display them). Repeated layout references are allowed.

| Component `type` | Flat properties used by the UI |
| --- | --- |
| `stat` | `label`, `value`, `unit`, `copyable` |
| `gauge` | `label`, numeric `value`, `min`, `max`, `unit`, `thresholds` |
| `status_indicator` | `label`, `value`, `mapping` keyed by stringified value |
| `log_stream` | `label`, `value` (string event or array), `max_items` |
| `action_group` | `label`, `items` array |
| `link` | `label`, `uri`, `text` |

Mapping entries contain `text`, `color`, `icon`, optional `animate`. All components
may carry `monitors` with `condition`, `severity`, `message`; the core evaluates
conditions against `value`. The typed Go decoder checks declared field types,
but `value` is arbitrary JSON except validated log_stream forms. The Python SDK
adds typed initial/update checks; core per-widget range/expression validation
remains future work. Valid structural discovery does not guarantee sensible
measurement values or a valid monitor expression.

Actions are supported only through `action_group.items[].action_id`, using the
same ID character rules and unique IDs within a service. Item fields are
`action_id`, `label`, optional `style` and `confirm`. Standalone component
`action_id` and items on other component types are rejected. Commands contain
`action_id`, `issued_by` (`synapse-ui`) and an RFC3339 `timestamp`. The core checks
the action against stored declarations before publication; dispatch success does
not prove that the Axon executed it.

## Updates and known synchronization gaps

Heartbeats and updates resend a full snapshot, typically every TTL/2 seconds.
Omitted capabilities are removed; partial updates are unsupported. Log streams are normalized by the server on every registration: a nonempty string
appends one event, a string array replaces the history (an empty array clears it),
and null/missing or an empty string preserves history. Other value types and
non-string array entries are rejected without changing storage. `max_items` is
nonnegative, with zero/omitted meaning ten; positive values keep only the newest
entries. First registration follows the same normalization and retention rules.
Each repeated string publication is a new event; clients should send bounded
array snapshots for heartbeats to avoid duplicating events. Browser consistency
is being verified in [SYN-103](../tasks/T0-bug-fix/define-log-stream-merge.md).

The browser consumes same-origin `GET /api/v1/events` SSE `services` events.
Each event contains a complete array of persisted services, with no registration
credentials. Accepted HTTP/MQTT discovery and TTL transitions notify connected
browsers. Every connection begins with a full snapshot; clients replace their
state and reconnect to reconcile missed updates. Periodic snapshots are sent
at fifteen-second intervals. Raw discovery messages are not browser state.

The shared [fixture](../testdata/discovery-v1.json) exercises all six components in
backend transport tests and actual Vue card rendering. [TOML](axon_toml_spec.md)
and [SDK](sdk_specification.md) documents describe the implemented Python mapping. The Python reference example uses the maintained shape; runtime acceptance is
tracked in its [migration task](../tasks/T1-feature/restore-reference-axon.md).
