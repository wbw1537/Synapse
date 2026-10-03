# Synapse

**The nervous system of your homelab.**

Synapse is a self-hosted homelab dashboard. Axons register over MQTT, report
service health and metrics, and execute declared actions. The Go binary embeds
the Vue interface, SQLite and the MQTT broker.

**[Read the documentation →](https://wbw1537.github.io/Synapse/)**

GitHub Pages is the single reference for installing, configuring and using Synapse:

- [Get started](https://wbw1537.github.io/Synapse/getting-started/)
- [Configuration](https://wbw1537.github.io/Synapse/configuration/)
- [Access and deployment](https://wbw1537.github.io/Synapse/access/)
- [Python SDK](https://wbw1537.github.io/Synapse/python-sdk/)
- [Discovery protocol](https://wbw1537.github.io/Synapse/protocol/)

The website source is [website/docs/](website/docs/index.md). If Pages has not yet
been enabled, read those files directly; they are the same instructions.

## Development

Start with [AGENTS.md](AGENTS.md) and the [internal knowledge map](docs/README.md).
Development procedures and verification are in [docs/development.md](docs/development.md)
and [docs/testing.md](docs/testing.md); work is tracked in [tasks/](tasks/README.md).
See [documentation maintenance](docs/documentation-site.md) for preview and publishing.

## License

Apache-2.0
