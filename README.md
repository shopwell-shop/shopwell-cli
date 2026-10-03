# Shopwell CLI

[![Hosted By: Cloudsmith](https://img.shields.io/badge/OSS%20hosting%20by-cloudsmith-blue?logo=cloudsmith&style=flat-square)](https://cloudsmith.com)
[![skills.sh](https://skills.sh/b/shopwell-shop/shopwell-cli)](https://skills.sh/shopwell-shop/shopwell-cli)

Shopwell CLI is a command line companion for common Shopwell account, project, and extension workflows.

## Table of Contents

- [What it helps with](#what-it-helps-with)
- [Highlights](#highlights)
- [Install](#install)
- [Usage](#usage)
- [Repository Layout](#repository-layout)
- [Official Agent Skills](#official-agent-skills)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## What it helps with

Use this CLI to create and set up Shopwell projects, manage Shopwell accounts, build and validate extensions, automate project maintenance, or run everyday developer tasks without leaving the terminal.

## Highlights

- Installation, initialization, and migration wizards for creating new Shopwell projects or moving existing projects into the local development setup
- Integrated Docker-based local development environment with an interactive terminal user interface (TUI)
- Daily development loop support from one place: start environments, run watchers, check logs, access database and mail tools
- Project automation commands for config, cache, admin, storefront, and CI workflows
- Account-related commands under `shopwell-cli account`
- Extension build, validation, formatting, changelog, and packaging helpers
- Non-interactive mode for scripts and CI

## Install

### From source with Go

```bash
go install github.com/shopwell-shop/shopwell-cli@latest
```

### Build locally from this repository

```bash
git clone https://github.com/shopwell-shop/shopwell-cli.git
cd shopwell-cli
go build -o bin/shopwell-cli .
```

## Usage

Show the available commands:

```bash
shopwell-cli --help
```

Common command groups:

```bash
shopwell-cli account --help
shopwell-cli extension --help
shopwell-cli project --help
```

Generate a CycloneDX SBOM for a project without running the full CI build:

```bash
shopwell-cli project sbom
shopwell-cli project sbom ./my-shop --format cyclonedx-json --output sbom.json
```

If you need CI-friendly behavior, disable prompts:

```bash
shopwell-cli --no-interaction <command>
```

Silence update notifications for a single command:

```bash
shopwell-cli --no-update-hint
```

## Configuration

Shop URL and Admin API credentials belong under named `environments` in `.config/shopwell-project.yml`. Project-level keys such as `build`, `dump`, `docker`, and `php_version` stay at the top level.

```yaml
environments:
  local:
    type: local
    url: http://127.0.0.1:8000
    admin_api:
      username: admin
      password: shopwell
  staging:
    url: https://staging.example.com
    admin_api:
      client_id: ...
      client_secret: ...
```

Omit `-e` / `--env` to target `environments.local`. Use `-e staging` (or another name) to target a different environment.

`shopwell-cli project create` and `shopwell-cli project config init` write `environments.local`. Top-level `url` and `admin_api` are deprecated: they are used only when `environments.local` is absent. When both are present, `environments.local` wins.

### Disable update notifications

To disable update notifications for the current shell session, set:

```bash
export SHOPWELL_CLI_NO_UPDATE_NOTIFICATION=true
```

To make the setting **persistent**, add the export command to your shell profile, such as `~/.bashrc` or `~/.zshrc`.

## Repository Layout

- `cmd/`: Cobra command groups for account, extension, and project workflows
- `internal/`: implementation packages for APIs, build steps, validation, TUI, and utilities
- `.github/`: automation and workflow definitions
- `scripts/`: repository helper scripts
- `env-bridge/`: environment bridge helper entrypoint

## Official Agent Skills

Shopwell CLI maintainers provide official Agent Skills in this repository's `skills/` folder. These skills teach compatible AI coding agents how to use Shopwell CLI workflows safely and correctly. 

Install them with:

```bash
npx skills add shopwell-shop/shopwell-cli
```

Available skills:

- **`shopwell-cli`** — Use Shopwell CLI for project, extension, and account workflows. Validate, build, format, and troubleshoot extensions. Reason about command safety and side effects.
- **`shopwell-cli-docker`** — For Docker-backed projects. Run commands in the correct container, manage development environment lifecycle, and troubleshoot Docker execution.
- **`shopwell-cli-extension-store`** — Assess an extension's readiness for Shopwell Store distribution. Read-only: classifies metadata, localization, and asset findings against Store requirements and cites a re-checkable source for each, without modifying files.

## Documentation

- Official docs: <https://developer.shopwell.cn/docs/products/cli/>

## Contributing

Contributions are welcome. If you want to improve commands, docs, or developer workflows, open an issue or send a pull request.
See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

See [LICENSE](LICENSE).
