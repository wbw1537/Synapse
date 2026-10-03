# Decision: one public instruction source

Date: 2026-10-03. Status: accepted. Task:
[SYN-111](../../tasks/T3-infrastructure/public-documentation-authority.md).

## Context

The first documentation site copied every current catalog entry, including
internal development and acceptance records. Root README, SDK README and two SDK
specifications repeated usage and lifecycle instructions. Editing one copy could
leave the others stale.

## Decision

GitHub Pages is the user-facing instruction authority. `website/docs/` is its sole
versioned Markdown source and MkDocs reads it directly. Each instruction topic has
one owner page. Root/SDK READMEs point to that site and source instead of carrying
another guide. Public instructions cannot link into internal development records.

`docs/` retains development records, decisions, source maps and historical notes;
`tasks/` retains execution contracts and evidence. Catalog status tracks document
maturity across these directories and does not implicitly publish documents.
Remove the catalog-based staging script and merge duplicate SDK prose into the
public SDK guide. Historical evidence remains dated rather than being rewritten
as present instructions.

## Consequences

Behavior changes update the public owner and relevant internal evidence. Preview
requires only MkDocs; staging output and rewrite logic are no longer needed.
Internal documents cannot enter the site by becoming current. New public pages
must be cataloged and included in navigation. External Pages setup and package
publication remain separate operations; local build success is not deployment.
