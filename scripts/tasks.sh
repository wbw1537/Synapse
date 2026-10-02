#!/usr/bin/env bash
set -euo pipefail
TASK_SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
python3 "$TASK_SCRIPT_DIR/../tasks/sync_index.py" update
exec python3 "$TASK_SCRIPT_DIR/tasks_server.py" --port "${SYNAPSE_TASKS_PORT:-6060}"
