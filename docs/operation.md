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

See [access policy](access.md) before selecting listener/published addresses. Set
two distinct secrets; native listeners and Compose host ports default to loopback.
Use HTTPS/secure cookies for remote browser access, and a private network or tunnel
for MQTT. Back up SQLite before upgrades; legacy registered rows are not migrated.

`.github/workflows/ci.yml` runs verification on branch pushes, PRs and `v*` tags.
Verification builds frontend before Go embedding checks, runs workspace and
application tests, race checks and SDK lifecycle acceptance. Docker publication
requires a version-tag push and successful verification; ordinary checks do not
publish or deploy. Remote CI execution is separate from local workflow validation.
