# Define and enforce runtime access boundaries

**ID:** SYN-107
**Status:** done
**Created:** 2026-10-02
**Priority:** high
**Depends-on:** -

## Goal

Prevent untrusted clients from reading operational data, forging discovery state or invoking declared recovery actions.

## Acceptance

- [x] Document an explicit homelab deployment threat model and chosen access policy.
- [x] HTTP reads/actions and MQTT topics enforce that policy, including browser-originated traffic.
- [x] Listener configuration supports intentional network exposure without breaking the internal core client.
- [x] Tests show permitted interactions succeed and unauthorized publication/action paths fail.

## Scope

Refine authentication and MQTT ACL design before implementation; account for token visibility in raw discovery subscriptions and current all-interface bindings. Do not claim payload token checking secures broker subscriptions or command publication.

## Verification

**Evidence:** `go test -race ./internal/api ./internal/broker ./internal/service -count=1` passed real positive/negative HTTP and MQTT checks, including authenticated commands, forbidden subscriptions/publications, invalid role credentials, Origin, login/logout, session expiry and active-stream revocation. Final browser login, HttpOnly visibility check, foreign-Origin rejection and logout revocation passed against fresh temporary core. Compose configuration validated; native/dev listeners observed bound to loopback. Independent agent review unavailable under active delegation policy; solo cross-layer review plus runtime verification completed, with limits recorded in docs/acceptance.md.

## Execution

Implemented separate operator/Axon realms, bearer or HttpOnly operator sessions, origin checks, broker authentication and service-topic ACLs, ephemeral core credentials, loopback defaults and packet-log redaction. Shared-Axon identity limits and TLS/private-network deployment requirements are documented in docs/access.md.
