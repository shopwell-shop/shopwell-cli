# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

## Shopwell repository guardrails

- This is an independent Shopwell repository. It must not be placed in the
  upstream GitHub fork network or share upstream Git history.
- Shopwell-owned code uses the Apache License 2.0. Keep the root `LICENSE` as
  the standard Apache-2.0 text and preserve every registered upstream license verbatim in
  the root `NOTICE`.
- Runtime code and workflows must not depend on upstream organization packages,
  Actions, or repositories. Keep real third-party identities unchanged.
- Before committing, pushing, releasing, or completing a sync, run:
  `../sync-upstream/bin/syncctl audit-license shopware-cli`,
  `../sync-upstream/bin/syncctl audit-repository-identity shopware-cli`,
  `../sync-upstream/bin/syncctl audit-dependency-parity shopware-cli`, and
  `../sync-upstream/bin/syncctl audit-upstream-dependencies shopware-cli`.

## Development Commands

### Building and Running
- **Build the CLI**: `go build -o shopwell-cli .`
- **Run directly**: `go run main.go [command]`
- **Run tests**: `go test ./...`
- **Run specific test**: `go test ./[package_path]`
- **Verbose test output**: `go test -v ./...`

### Code Quality
- **Prefer Go 1.24 packages** like slices package
- **Check Go modules**: `go mod tidy`
- **Format code**: `go fmt ./...`
- **Run static analysis**: `go vet ./...`

## Architecture Overview

### Core Command Groups
1. **`account/`** - Shopwell Account management (login, companies, extensions)
2. **`extension/`** - Extension development tools (build, validate, AI assistance)
3. **`project/`** - Shopwell project management (creation, configuration, deployment)

### Key Internal Packages
- **`internal/verifier/`** - Code quality tools (PHPStan, ESLint, Twig linting)
- **`internal/account-api/`** - Shopwell Account API integration
- **`internal/system/`** - System utilities (PHP/Node detection, filesystem)
- **`internal/packagist/`** - Composer/Packagist integration
- **`internal/llm/`** - AI integration (OpenAI, Gemini, OpenRouter)
- **`internal/html/`** - Twig + HTML parsing and formatting
- **`internal/git/`** - Git operations
- **`internal/ci/`** - CI/CD configuration generation

### Extension System
The CLI supports three extension types:
- **Platform Plugins**: Standard Shopwell 6 plugins with `composer.json`
- **Shopwell Apps**: App system extensions with `manifest.xml`
- **Shopwell Bundles**: Custom bundle implementations

Extension detection is automatic based on file presence. Each type has specific build, validation, and packaging rules.

### Configuration Files
- **`.shopwell-cli.yaml`** - Global CLI configuration
- **`.config/shopwell-extension.yml` / `.shopwell-extension.yml` / `.shopwell-extension.yaml`** - Extension-specific settings (schema: `internal/extension/config_schema.json`, embedded and exposed via `shopwell-cli extension config-schema`). Only one file is loaded and the first is preferred / used for new projects, the others are legacy and still supported
- **`.config/shopwell-project.yml` / `.shopwell-project.yaml` / `.shopwell-project.yml`** - Project-specific settings (schema: `internal/shop/config_schema.json`, embedded and exposed via `shopwell-cli project config-schema`). Shop URL and Admin API credentials belong under `environments`; empty `-e` defaults to `environments.local`. Top-level `url` / `admin_api` are deprecated. An environment's `type` is `local`, `docker`, or `ssh`; `ssh` targets point at a remote host via `ssh.host`/`ssh.user`/`ssh.directory` and run commands through a multiplexed SSH connection (database connections are tunneled automatically). Only one file is loaded and the first is preferred / used for new projects, the others are legacy and still supported

## Development Patterns

### Command Structure
Commands follow Cobra CLI patterns with:
- Main command in `cmd/[group]/[group].go`
- Subcommands in `cmd/[group]/[group]_[subcommand].go`
- Service containers for dependency injection

### Testing Strategy
- Unit tests alongside source files (`*_test.go`)
- Use testify assert for test assertions (`github.com/stretchr/testify/assert`)
- Test data in `testdata/` directories
- Integration tests use real extension samples in `testdata/`
- Prefer assert.ElementsMatch on lists to ignore ordering issues
- Use t.Setenv for environment variables
- Use t.Context() for Context creation in tests

### Error Handling
- Use structured logging via `go.uber.org/zap`
- Context-based logging: `logging.FromContext(ctx)`
- Graceful error reporting to users

### Command Conventions
Follow `docs/COMMAND_CONVENTIONS.md` for every command, flag and message. The short version:
- Short texts are imperative, capitalised, no period; a parent command describes the group
- `Use` shows required arguments bare and optional ones in brackets, no angle brackets; every leaf command declares `Args`
- Errors are lowercase `cannot <verb> ...: %w` without Go identifiers, hints in quotes; log lines are capitalised sentences
- Any failure exits non-zero in every output format; never log success after a failed step
- Test the invoked `internal/` function, never the cobra layer, and test behaviour, not message wording

### AI Integration
The CLI includes AI-powered features for:
- Twig template upgrades (`extension ai twig-upgrade`)
- Code quality suggestions
- Automated fixes for common issues

LLM providers are configurable (OpenAI, Gemini, OpenRouter) with API key management.

## Extension Development Workflow

### Building Extensions
```bash
# Build extension (auto-detects type)
shopwell-cli extension build

# Watch mode for development
shopwell-cli extension admin-watch

# Validate extension
shopwell-cli extension validate

# Create distribution package
shopwell-cli extension package
```

### Project Management
```bash
# Create new project
shopwell-cli project create

# Build assets
shopwell-cli project admin-build
shopwell-cli project storefront-build

# Development servers
shopwell-cli project admin-watch
shopwell-cli project storefront-watch
```

## Code Quality Integration

The verifier registers tools through the name-only `Tool` interface. `CheckTool`, `FixTool`, and `FormatTool` add capabilities; commands select the relevant capability before applying `--only` or `--exclude`. An unsupported tool name is an error, and `ToolList[T]` preserves the capability type through filtering.

- **Checkers**: `builtin` (legacy alias: `sw-cli`), PHPStan, ESLint, Stylelint, Storefront Twig
- **Fixers**: Rector, ESLint, Stylelint, Symfony XML conversion
- **Formatters**: PHP-CS-Fixer, Prettier

`extension validate` runs all checkers by default. The deprecated `--full` flag remains accepted but has no effect; use `--only` or `--exclude` to select checkers. The extension commands report whether each tool was invoked or skipped; invocation does not guarantee that files were analyzed or changed.
