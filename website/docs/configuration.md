# Configuration reference

The native Go process reads environment variables and does not automatically
load `.env`. Docker Compose reads `.env` for substitution; its
[Compose file](https://github.com/wbw1537/Synapse/blob/main/docker-compose.yml)
selects which variables reach the container.

Use [.env.example](https://github.com/wbw1537/Synapse/blob/main/.env.example)
as the editable template. For a native run, copy it to `.env`, edit the settings,
then load it in a POSIX-compatible shell:

```sh
set -a
. ./.env
set +a
```

Only source a file you control. Quote values containing spaces or shell special
characters; `.env` here must use shell-compatible assignments. Keep real secrets
out of source control.

## Core and access

| Variable | Native default | Meaning |
| --- | --- | --- |
| `SYNAPSE_HTTP_PORT` | `127.0.0.1:8080` | HTTP/UI listener address (`host:port`) |
| `SYNAPSE_MQTT_PORT` | `127.0.0.1:1883` | Embedded MQTT TCP listener address |
| `SYNAPSE_WS_PORT` | `127.0.0.1:8083` | MQTT WebSocket listener; the dashboard uses HTTP SSE |
| `SYNAPSE_DB_PATH` | `synapse.db` | SQLite database path |
| `SYNAPSE_AUTH_TOKEN` | required | Axon connection and registration secret |
| `SYNAPSE_ADMIN_TOKEN` | required | Distinct operator login/API secret |
| `SYNAPSE_COOKIE_SECURE` | `false` | Enable secure cookies and HTTPS origin checking behind an HTTPS proxy |

The two tokens must be nonempty and different or startup fails. The
[access policy](access.md) defines authentication, topic permissions and the trust
model. Compose binds listeners inside the container and publishes host ports on
loopback; changing native listener variables alone does not change Compose's
published addresses.

## Email alerts

| Variable | Default | Meaning |
| --- | --- | --- |
| `SYNAPSE_ENABLE_ALERTS` | `false` | Enable SMTP notifications |
| `SYNAPSE_SMTP_HOST` | empty | SMTP hostname |
| `SYNAPSE_SMTP_PORT` | `587` | SMTP port |
| `SYNAPSE_SMTP_USER` | empty | SMTP username |
| `SYNAPSE_SMTP_PASS` | empty | SMTP password |
| `SYNAPSE_SMTP_FROM` | `synapse@localhost` | Sender address |
| `SYNAPSE_SMTP_TO` | empty | Comma-separated recipients |

When running through Compose, add the needed SMTP variables to its service
`environment` section; examples are already present as comments. Setting them
only in `.env` does not forward them automatically. See the
[operations guide](operations.md#email-alerts) for SMTP troubleshooting.
