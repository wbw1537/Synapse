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
but `value` is arbitrary JSON; per-widget value/range validation and SDK fail-fast
checks remain future work. Valid structural discovery does not guarantee sensible
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
Omitted capabilities are removed; partial updates are unsupported. The existing
log-stream merge is a special case and has known snapshot/event inconsistencies
tracked in [SYN-103](../tasks/T0-bug-fix/define-log-stream-merge.md).

The browser initially reads stored services, then consumes **raw** MQTT discovery
messages. It can see messages rejected by the core and misses HTTP/TTL updates.
This task aligns structure and persistence validation; authoritative browser state
is tracked in [SYN-102](../tasks/T0-bug-fix/synchronize-authoritative-service-state.md).

The shared [fixture](../testdata/discovery-v1.json) exercises all six components in
backend transport tests and actual Vue card rendering. [TOML](axon_toml_spec.md)
and [SDK](sdk_specification.md) documents specify planned mapping to this contract;
neither describes an implemented SDK. The legacy Python example has its own
[migration task](../tasks/T1-feature/restore-reference-axon.md).
