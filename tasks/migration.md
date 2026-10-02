# Historical plan reconciliation

Inspected 2026-10-02 against source at `36f9ac8` and local commit history.
Historical completion reports are preserved; their claimed tests were not rerun.

| Old plan | Observed current state / successor |
| --- | --- |
| MVP-001–005 | Core, registry, HTTP, UI and SMTP/expression code exist |
| MVP-006–008 | Expanded widgets, runbook drawer and command publication exist; some original checkboxes/status headers remain stale |
| SDK-001–002 | Design documents exist, but completion reports reference missing migration/API documents; design and implementation conflict |
| SDK-003 | Layout/component refactor exists; compilation did not establish protocol/example compatibility |
| SDK-004 | No `sdk/python/`; replaced by SYN-104 after contract alignment |

Outstanding behavior is captured in the current queue. Old plans are neither
retroactively declared tested nor used to select the next task. Original SDK
compatibility ideas were superseded by the breaking refactor; do not reintroduce
legacy compatibility merely because an old design calls for it.
