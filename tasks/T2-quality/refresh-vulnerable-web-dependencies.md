# Refresh vulnerable web dependencies

**ID:** SYN-109
**Status:** done
**Created:** 2026-10-03
**Priority:** high
**Depends-on:** -

## Goal

Remove the seven security advisories reported by npm after removing unused browser MQTT, preserving dashboard behavior and the current build toolchain.

## Acceptance

- [x] Compatible dependency refresh removes all current npm audit findings.
- [x] Frontend typecheck/build, actual Vue rendering and dashboard runtime remain correct.
- [x] Lockfile is reproducible with npm ci in a clean export.

## Scope

Use nonbreaking updates within existing ranges. Retain the pinned Vite alias unless evidence requires a change. Audit currently reports markdown-it/linkify-it and build-tool transitive packages; no force upgrades.

## Verification

**Evidence:** `npm --prefix web audit fix --ignore-scripts` reported zero vulnerabilities. Fresh clean-export npm ci also reported zero findings, and Vue typecheck/build/rendering passed. Fresh final browser exercised login, live discovery, SDK callback, logs, TTL/recovery, reconnection and mobile sign-out using updated assets.

## Execution

Removed unused MQTT browser dependency and refreshed eleven compatible packages without changing the pinned Vite alias.
