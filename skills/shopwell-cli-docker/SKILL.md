---
name: shopwell-cli-docker
description: Use when working on Docker-backed Shopwell projects and needing to run Shopwell or Symfony CLI commands, interact with project services or databases, start or stop the development environment, troubleshoot Docker execution, or determine whether a command should run on the host or inside the project container.
---

# Shopwell CLI with Docker

In a Shopwell CLI-managed Docker project, prefer asking Shopwell CLI to perform the operation instead of manually deciding how to enter containers.

Shopwell CLI knows the project configuration and selected execution environment and can route supported operations through the correct executor.

## Core rule

Do not default to:

```bash
docker compose exec ...
```

when Shopwell CLI already exposes the operation.

Prefer:

```bash
shopwell-cli project ...
```

The CLI should decide when Docker is required.

## Symfony and Shopwell Console commands

When a task requires `bin/console`, use:

```bash
shopwell-cli project console <command> [arguments]
```

Instead of manually doing something like:

```bash
docker compose exec web php bin/console <command>
```

prefer:

```bash
shopwell-cli project console <command>
```

`project console` resolves the current Shopwell project and its configured executor. In a Docker-backed environment, the CLI routes the Symfony command or Composer script through the project container.

For the full Composer CLI, use:

```bash
shopwell-cli project composer <arguments>
```

Instead of manually doing something like:

```bash
docker compose exec web composer <arguments>
```

In a Docker-backed environment, the CLI routes the Composer command through the project container.

This applies to Shopwell and Symfony Console workflows such as:

- cache commands;
- plugin and app lifecycle;
- indexing;
- migrations;
- scheduled tasks;
- system configuration;
- other commands exposed by `bin/console`;
- custom Composer scripts from the project `composer.json`.

Before executing an unfamiliar or state-changing Symfony command:

```bash
shopwell-cli project console <command> --help
```

Do not assume a Symfony command is safe merely because it exists.

Discover Composer scripts from the project `composer.json` or `shopwell-cli project console list`. Do not pass `--help` to a Composer script to inspect it; that argument is forwarded to the script and can still run it.

## Development environment

Use Shopwell CLI to manage the development environment:

```bash
shopwell-cli project dev
shopwell-cli project dev start
shopwell-cli project dev status
shopwell-cli project dev stop
```

Do not replace these with `docker compose up` or `docker compose down` unless there is a specific reason to bypass the CLI.

The CLI knows the configured environment and can prepare and manage the Docker setup required by the project.

## Other project operations

Prefer the corresponding Shopwell CLI command for supported tasks such as builds, watchers, logs, validation, upgrades, project maintenance, and extension management.

Do not manually wrap a Shopwell CLI command in `docker compose exec`.

If Shopwell CLI already exposes the operation, let the CLI handle Docker.

## Database work

First determine what kind of database operation the user actually needs.

For database exports, inspect:

```bash
shopwell-cli project dump --help
```

For database-related operations exposed through Symfony or Shopwell Console, prefer:

```bash
shopwell-cli project console <command>
```

For arbitrary SQL, first inspect the available Shopwell CLI and Symfony commands. Do not invent a CLI command.

If no suitable Shopwell CLI abstraction exists, direct database access may be necessary. Before doing that:

1. identify the configured database service and target environment;
2. prefer discovered configuration over hard-coded container names or credentials;
3. use read-only queries where possible;
4. treat writes, deletes, schema changes, imports, and bulk updates as destructive operations.

Never execute database-changing commands against an ambiguous environment.

## Composer, PHP, and npm

Prefer a higher-level Shopwell CLI operation when one exists.

Do not immediately run:

```bash
docker compose exec web composer ...
docker compose exec web php ...
docker compose exec web npm ...
```

First check whether Shopwell CLI already exposes the intended workflow.

For Composer itself, prefer:

```bash
shopwell-cli project composer ...
```

If no CLI abstraction exists:

1. confirm this using the relevant command group's `--help`;
2. determine the configured Docker service and working directory;
3. run the low-level tool in the correct project environment;
4. avoid hard-coding paths and service names when they can be discovered.

Direct Docker execution is the fallback, not the default.

## Environment selection

Do not assume that an environment named `local`, `dev`, `staging`, or similar is Docker-backed or safe to modify.

Inspect the project config using the CLI's discovery precedence (`.config/shopwell-project.yml`, `.shopwell-project.yaml`, then `.shopwell-project.yml`) and use the CLI's environment selection mechanisms. Empty `-e` / `--env` targets `environments.local`. Store shop URL and Admin API credentials under `environments`, not at the top level.

Before destructive commands, establish exactly which environment will be affected.

## Starting from a user request

When a user asks for something like:

- "clear the Shopwell cache";
- "run this Symfony command";
- "update the plugin";
- "run migrations";
- "check the database";
- "start the shop";
- "show me the logs";

do not translate the request directly into raw Docker commands.

First ask: **Does Shopwell CLI already know how to do this?**

For Symfony commands, the answer should usually begin with:

```bash
shopwell-cli project console ...
```

For development environment lifecycle, begin with:

```bash
shopwell-cli project dev ...
```

For other workflows, inspect:

```bash
shopwell-cli project --help
```

## Troubleshooting

When Docker-backed execution fails:

1. Run `shopwell-cli project dev status`.
2. Inspect the relevant Shopwell CLI command's `--help`.
3. For Shopwell CLI debug output, place `--verbose` before the command group, for example `shopwell-cli --verbose project console <command>`.
4. Verify the selected environment in project configuration.
5. Inspect Docker directly only after establishing how Shopwell CLI attempted to execute the operation.

Do not move a command from the container to the host just to make it work unless the project is explicitly configured for local execution.

## Safety

Docker does not make a command safe.

Commands executed through `project console`, Composer, the database, or project-management commands can still modify files, dependencies, application state, or data.

Before a destructive action:

- identify the environment;
- inspect the command;
- understand the expected changes;
- prefer previews or read-only checks when available;
- preserve unrelated user work.
