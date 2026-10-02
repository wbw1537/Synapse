# Synapse operations

Inspect the intended environment first. No production host or deployment target
is configured by this workflow. Do not reuse Beanpilot hosts, credentials, or scripts.

Local Compose: `docker compose ps`, `docker compose logs --tail 100 synapse`.
Inspect relevant logs without printing `.env`, raw credential-bearing discovery
payloads, or DB contents. Check listener ports, volume mounts, selected image and
Git revision when diagnosing a deployment.

`docker compose up -d --build` deploys the current working tree and changes the
runtime. Run only when requested/authorized. Plan DB backup and compatibility
before migrations; do not delete `synapse.db` merely because an old Alpha task
suggested resetting it. Verify HTTP, MQTT ingestion and UI after operations.

`.github/workflows/ci.yml` currently runs on `v*` tags only and publishes GHCR
images. Pushing a version tag therefore publishes artifacts. Routine PR checks
are a tracked improvement, not a currently available guarantee.
