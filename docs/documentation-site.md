# Documentation ownership and publishing

## Ownership

GitHub Pages is the single user-facing instruction reference. Its only Markdown
source is `website/docs/`. Edit that source for installation, configuration,
operations, authentication, SDK usage and protocol changes. Each topic has one
owner page; related pages link to it rather than repeat its instructions.

`docs/` holds internal development workflow, source maps, constraints, decisions,
acceptance evidence and historical records. It is not a second user guide and is
not published by MkDocs. Tasks remain in `tasks/`. Root and SDK READMEs are short
entry points that direct users to the website, with links to the same source for
readers before Pages is enabled.

The [catalog](index.json) records current, proposed and historical documents in
both source directories. Status describes document maturity, not publication.
MkDocs publishes only `website/docs/`, independently of catalog status. Add a
new public page to the catalog and `mkdocs.yml` navigation. Keep code/sample links
as explicit GitHub links; public relative links should stay inside website/docs/.

## Preview and verify

From the repository root:

```sh
python -m venv .venv-docs
.venv-docs/bin/python -m pip install -r requirements-docs.txt
.venv-docs/bin/python -m mkdocs build --strict
.venv-docs/bin/python -m mkdocs serve
python3 scripts/agent.py verify --scope workspace
```

MkDocs reads public sources directly and reloads when they change. There is no
staging or copying step. The local URL includes `/Synapse/`. `.build/site/` and
`.venv-docs/` are generated and ignored. Strict builds check navigation and links;
workspace checks cover both source directories and documentation boundaries.

## GitHub Pages

After merging, select **Settings → Pages → Build and deployment → Source →
GitHub Actions**, then run the Documentation workflow on main or push a relevant
change. The default project URL is `https://wbw1537.github.io/Synapse/`; the
`github-pages` environment reports the actual deployment URL.

PRs only build and verify. Main pushes and manual main runs can upload and deploy
a verified artifact. Only deployment has `pages: write` and `id-token: write`.
Repository/environment settings may require administrator setup. Configuration
and local acceptance do not imply that the site has been published. See
[GitHub's Pages workflow reference](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages).

SDK package publication is separate; its user-facing installation instructions
live only in the [public SDK guide](../website/docs/python-sdk.md).
