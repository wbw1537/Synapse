# Documentation website

The website uses MkDocs with Material, navigation, full-text search and Mermaid
rendering. Markdown remains authoritative in the root README, docs/ and the
Python SDK README. `docs/index.json` selects current documents. Historical plans,
tasks and code links point back to GitHub source instead of becoming duplicate
website pages. Build output contains only explicitly selected Markdown and the
static site's generated assets.

## Preview and verify locally

From the repository root:

```sh
python -m venv .venv-docs
.venv-docs/bin/python -m pip install -r requirements-docs.txt
python scripts/docs_site.py
.venv-docs/bin/python -m mkdocs build --strict
.venv-docs/bin/python -m mkdocs serve
```

MkDocs prints the local URL, including the `/Synapse/` project subpath. After
editing source Markdown rerun `scripts/docs_site.py`; preview then reloads the
staged changes. `.build/` and `.venv-docs/` are generated and ignored. Build tools
run only on the developer/CI machine; deployed Pages is a static site.

Add new current pages to the catalog and `mkdocs.yml` navigation. Relative links
among published documents remain site links. Links to non-published files and
historical documents become GitHub source links. A missing/outside-repository
relative target fails staging; fenced examples remain unchanged. Use the source
Markdown files, not generated copies, when editing.

## GitHub Pages setup

After reviewing and merging the changes, open the repository's **Settings →
Pages → Build and deployment → Source** and select **GitHub Actions**. Then run
the Documentation workflow from main, or push a relevant documentation change.
The expected default address is `https://wbw1537.github.io/Synapse/`; the workflow's
`github-pages` environment reports the actual published URL.

The workflow strictly builds docs on PRs without publishing. Main pushes and
manual runs on main upload a verified Pages artifact and deploy it. Deployment
alone receives `pages: write` and `id-token: write`; build jobs use read-only
repository permissions. GitHub Pages settings/environment protections can still
require a repository administrator's setup. A successful local build does not
mean the site is already live. See [GitHub's workflow documentation](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages).

## Python SDK distribution

SDK publication is separate from publishing this website. Python 3.11+ clients
can install `python -m pip install ./sdk/python` from a checkout now. Once the
SDK changes exist on GitHub, another option is:

```sh
python -m pip install 'git+https://github.com/wbw1537/Synapse.git#subdirectory=sdk/python'
```

PyPI publication is optional. It would provide installation by package name,
`pip install synapse-axon`, after a verified version is actually released. No
PyPI publication or package-name availability is claimed here. See the
[SDK guide](../sdk/python/README.md) for its implementation and lifecycle.
