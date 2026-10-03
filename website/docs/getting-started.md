# Get started

## Prepare a checkout

```sh
git clone https://github.com/wbw1537/Synapse.git
cd Synapse
cp .env.example .env
```

Choose Docker Compose or a native build below.

## Docker Compose

Install Docker with Compose.

Edit `.env` and set two distinct, nonempty random secrets:
`SYNAPSE_AUTH_TOKEN` for Axons and `SYNAPSE_ADMIN_TOKEN` for operators. Consult the
[configuration reference](configuration.md) for all settings.

```sh
docker compose up -d --build
```

Open **http://localhost:8080** and sign in with the operator token. Compose
publishes the HTTP, MQTT and MQTT WebSocket ports on host loopback and persists
SQLite in the `synapse-data` volume. See [access and deployment](access.md) before
admitting remote browsers or Axons.

## Build from source

Requirements: Go 1.25+ and Node.js 20+ with npm.

```sh
npm ci --prefix web
npm --prefix web run build
go build -o synapse ./cmd/synapse
```

Build the frontend first: Go embeds its generated assets. Set the two distinct
secrets and any optional settings in the process environment using the
[configuration reference](configuration.md), then run:

```sh
./synapse
```

Native listeners default to loopback. Open **http://localhost:8080** and sign in
with the operator token. The default database is `synapse.db` in the process's
working directory.

## Add an Axon

Follow the [Python SDK guide](python-sdk.md) for the reference client, or implement
the [discovery protocol](protocol.md) in another language. Axons measure their own
services and execute their own declared actions; the core stores and displays
reported state and dispatches declared action IDs.
