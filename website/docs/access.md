# Runtime access policy

Synapse targets a trusted homelab. The core stores operational data and publishes
predeclared actions; Axons alone execute handlers. Operators and Axons are separate
trust realms. Startup requires distinct, nonempty `SYNAPSE_ADMIN_TOKEN` and
`SYNAPSE_AUTH_TOKEN`; there is no built-in default secret.

## HTTP and browser

Service reads, action requests and `GET /api/v1/events` require an operator bearer
token or session. API clients send `Authorization: Bearer <admin-token>`. Browser
users enter the admin token at the sign-in form; `POST /api/v1/session` exchanges
it for a random server-side session with a twelve-hour lifetime. The cookie is
HttpOnly and SameSite Strict; the token is not stored in browser storage. Sessions
are in memory and expire on restart. `DELETE /api/v1/session` revokes the session;
streams recheck authorization before every snapshot, at most fifteen seconds apart.
Foreign Origin headers are rejected. `SYNAPSE_COOKIE_SECURE=true` enables HTTPS
cookies and HTTPS origin checking when using a TLS reverse proxy.

`POST /api/v1/discovery` checks the Axon token in its body; an Axon token does not
authorize reads or actions. Discovery is limited to 1 MiB. HTTP has a header-read
timeout; slow SSE writes have a deadline. Static assets and the login form are public.

## MQTT

| Client | Authentication | Permitted topics |
| --- | --- | --- |
| Axon | Username `axon`, password Axon token, client ID service ID | Publish `synapse/v1/discovery/{own-id}`; subscribe `synapse/v1/command/{own-id}` |
| Core | Username `synapse-core`, startup-generated secret, client ID `synapse_core` | Subscribe discovery; publish individual command topics |

Anonymous connections, wildcard Axon subscriptions, other-service publications
and Axon command publications are rejected. MQTT v3 forbidden publications may
cause disconnect; PUBACK is never evidence of business acceptance. The discovery
payload must still include the matching Axon token and service ID. The core
resubscribes on reconnect. Raw discovery, including its token, is never sent to
the browser. Broker packet log fields are redacted.

Axons share one credential: an authorized Axon can impersonate another service by
reconnecting under that ID. This is explicitly a trusted-Axon realm, not per-Axon
credential isolation or multi-tenant security. Operators can invoke every declared
action. There are no users, roles, audit history or execution-result acknowledgments.

## Deployment

Native defaults bind all listeners to loopback. The development launcher uses
loopback and temporary storage. Compose listeners bind inside the container while
published host ports default to loopback, and secrets must be provided explicitly.
To admit remote Axons, deliberately select listener/published addresses and a
private network or VPN/firewall policy. For remote browser access use HTTPS through
a trusted reverse proxy and enable secure cookies. MQTT TCP/WS has no built-in TLS;
protect it with a VPN or TLS tunnel. Do not expose plain credentials to an untrusted
network. Rotate both configured secrets and restart after suspected compromise.
