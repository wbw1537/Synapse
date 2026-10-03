# Operations

## Inspect a Compose deployment

```sh
docker compose ps
docker compose logs --tail 100 synapse
```

Check listener addresses, the selected image/revision, mounts and available
storage. Do not publish `.env`, credential-bearing discovery payloads or database
contents when reporting a problem. Authentication and remote exposure are covered
by the [access policy](access.md).

## Upgrade and back up

Back up SQLite before an upgrade and verify the target version's compatibility.
Use a consistent SQLite backup procedure; if copying database files directly,
stop the service first and retain any associated WAL files. Do not delete the
existing database to resolve an upgrade failure. Legacy widget/action-array rows
are not automatically migrated to the current [protocol](protocol.md).

For a source-based Compose deployment, after reviewing and updating the checkout:

```sh
docker compose up -d --build
```

After restarting, verify dashboard login, Axon registration, metrics, declared
callbacks and TTL offline behavior. An accepted command request proves dispatch,
not successful execution inside an Axon.

## Email alerts

Enable and configure SMTP using the [configuration reference](configuration.md#email-alerts).
The SMTP client can upgrade with STARTTLS on servers that advertise it. Implicit
TLS on port 465 is not supported by the current SMTP implementation. For
`SASL PLAIN authentication failed`, verify the credentials, SMTP authentication
policy and shell quoting of the password. Never include the password in a report.

## Client or dashboard state

Use the operator-authenticated service reads in the [protocol](protocol.md) to
check persisted state. MQTT delivery acknowledgment does not prove that a
registration passed validation. Review core validation logs and confirm matching
service ID, token and topic when a client appears absent.

The dashboard consumes full server snapshots over SSE and reconciles state after
reconnecting. Proxy configuration must permit long-lived HTTP event streams.
For Python client heartbeat, reconnect and shutdown behavior, consult the
[SDK lifecycle](python-sdk.md#lifecycle-and-runtime-updates).
