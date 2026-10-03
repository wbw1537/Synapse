# Developer operations and releases

For operator procedures use the public [operations guide](../website/docs/operations.md)
and [access policy](../website/docs/access.md). This record governs development
and release work, rather than duplicating deployment instructions.

Identify the actual target environment before changing runtime state. No
production target is configured by the development workflow. Do not reuse
Beanpilot hosts, credentials or scripts. Deployment, restarts and database
migrations require authorization for the intended target. Plan compatibility and
backups before making changes; use disposable runtimes for development checks.

`.github/workflows/ci.yml` verifies branch pushes, PRs and `v*` tags. Verification
builds frontend assets before Go embedding checks, runs workspace and application
tests, race checks and SDK lifecycle acceptance. Docker publication requires a
version-tag push and successful verification. Ordinary checks do not publish or
deploy. Remote CI execution is separate from local workflow validation.

Documentation publication is maintained in [documentation-site.md](documentation-site.md).
SDK distribution metadata is in `sdk/python/pyproject.toml`; build/install evidence
belongs in the task and [testing records](testing.md), while public installation
steps belong in the [SDK guide](../website/docs/python-sdk.md).
