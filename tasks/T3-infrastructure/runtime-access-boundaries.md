# Define and enforce runtime access boundaries

**ID:** SYN-107
**Status:** backlog
**Created:** 2026-10-02
**Priority:** high
**Depends-on:** -

## Goal

Prevent untrusted clients from reading operational data, forging discovery state or invoking declared recovery actions.

## Acceptance

- [ ] Document an explicit homelab deployment threat model and chosen access policy.
- [ ] HTTP reads/actions and MQTT topics enforce that policy, including browser-originated traffic.
- [ ] Listener configuration supports intentional network exposure without breaking the internal core client.
- [ ] Tests show permitted interactions succeed and unauthorized publication/action paths fail.

## Scope

Refine authentication and MQTT ACL design before implementation; account for token visibility in raw discovery subscriptions and current all-interface bindings. Do not claim payload token checking secures broker subscriptions or command publication.

## Verification

Independent review plus positive/negative HTTP and MQTT integration checks against disposable storage. Current behavior is Alpha; no hardening results yet.
