# Agent Skills

Shopwell CLI publishes Agent Skills for AI coding tools.

The canonical skills live in this repository so that changes to Shopwell CLI and
changes to its agent guidance can be reviewed and released together.

## Repository structure

```text
shopwell-cli/
├── AGENTS.md
├── docs/
│   └── AGENT_SKILLS.md
└── skills/
    ├── shopwell-cli/
    │   └── SKILL.md
    ├── shopwell-cli-docker/
    │   └── SKILL.md
    ├── shopwell-cli-extension-store/
    │   └── SKILL.md
    └── shopwell-cli-extension-anonymize/
        └── SKILL.md
```

Each skill is a single `SKILL.md`, optionally with a `scripts/` directory for
helper scripts the agent executes (the Agent Skills layout). Keep reference
material inside `SKILL.md` itself.

Do not maintain separate Claude, Cursor, Codex, Copilot, or other
client-specific copies in this repository.

## Available skills

### `shopwell-cli`

General Shopwell CLI guidance.

It teaches agents how to:

- discover the current CLI command surface;
- prefer Shopwell CLI abstractions over lower-level tooling;
- distinguish project, extension, and account workflows;
- reason about command safety and side effects;
- work non-interactively;
- troubleshoot CLI failures.

### `shopwell-cli-docker`

Guidance for Shopwell CLI projects using Docker.

It teaches agents to let Shopwell CLI resolve the project environment instead
of defaulting to raw Docker commands.

In particular, Symfony Console commands should normally go through:

```bash
shopwell-cli project console <command>
```

and development environment lifecycle through:

```bash
shopwell-cli project dev
shopwell-cli project dev start
shopwell-cli project dev status
shopwell-cli project dev stop
```

### `shopwell-cli-extension-store`

Read-only Shopwell Store readiness assessment for an extension.

It teaches agents to:

- collect evidence by calling the CLI directly — two `extension validate` runs
  (normal and with `validation.store_compliance: true` in the extension
  configuration) plus reading the extension's metadata and icon;
- classify every finding against a fixed table with a re-checkable source;
- keep local file state separate from the remote Store listing, which it never
  inspects;
- never modify files.

### `shopwell-cli-extension-anonymize`

Bootstrap the `anonymize` section of an extension's
`.config/shopwell-extension.yml`.

It teaches agents to:

- edit only the config file the CLI actually loads;
- take `system_config` keys from `config.xml` and SystemConfig calls;
- take table and column names from entity storage names and migrations;
- choose faker expressions for personal data and `''` or `NULL` for secrets;
- omit a `system_config` row, or replace it with a custom value or a faker
  value;
- leave Shopwell core anonymization and unrelated extension validation alone.

## Source of truth

The files under `skills/` are the canonical source.

Do not duplicate their content in generated client-specific files.

The skills should contain durable workflow knowledge rather than an exhaustive
copy of the command reference. Agents are instructed to inspect the current
CLI with `--help`, which reduces the amount of guidance that needs updating
when commands or flags are added.

## Pull request checklist

For user-facing Shopwell CLI changes:

1. Implement the CLI change.
2. Check whether `skills/shopwell-cli/SKILL.md` is affected.
3. Check whether `skills/shopwell-cli-docker/SKILL.md`,
   `skills/shopwell-cli-extension-store/SKILL.md`, or
   `skills/shopwell-cli-extension-anonymize/SKILL.md` is affected.
4. Update the skill in the same PR when required.
5. Validate the skills.
6. Verify that the skills can still be discovered by the skills CLI.

Consider adding the following item to the repository PR template:

```text
- [ ] I reviewed `skills/` for user-facing CLI changes.
```

## Validation

Validate each skill against the Agent Skills format:

```bash
skills-ref validate ./skills/shopwell-cli
skills-ref validate ./skills/shopwell-cli-docker
skills-ref validate ./skills/shopwell-cli-extension-store
skills-ref validate ./skills/shopwell-cli-extension-anonymize
```

Verify repository discovery:

```bash
npx skills add . --list
```

CI should run these checks so an invalid skill cannot be merged.

Where practical, CI should also test commands explicitly referenced by the
skills, such as:

```text
project console
project dev
project dev start
project dev status
project dev stop
project dump
```

This catches obvious drift when commands are renamed or removed.

Semantic changes still require human review.

## Distribution

The canonical Agent Skills are maintained under `skills/` in this repository.

They are distributed through the Agent Skills ecosystem directly from the
`shopwell-shop/shopwell-cli` GitHub repository. Do not maintain client-specific
copies of the skills.

After merging changes to the default branch, verify public discovery:

```bash
npx skills add shopwell-shop/shopwell-cli --list
```

The `skills` CLI is responsible for installing the canonical skills for
supported AI clients. Shopwell CLI does not maintain client-specific copies.

## Installing with Shopwell CLI

Shopwell CLI exposes an `ai` command group as a thin front door for installing
integrations: it runs skills.sh for skills and writes the project's `.mcp.json`
for MCP servers. Installing a skill requires Node.js/npx on PATH (git-delivered
integrations also need git and network access).

```bash
# discover integrations
shopwell-cli ai list
shopwell-cli ai info deployment-helper

# install into an agent for the current Shopwell project (pin with @<tag>)
shopwell-cli ai add shopwell-cli --agent claude-code
shopwell-cli ai add deployment-helper@0.1.7 --agent claude-code

# install at the user level, from anywhere
shopwell-cli ai add shopwell-cli --agent claude-code --global

# preview, list what the CLI installed, and remove
shopwell-cli ai add shopwell-cli --agent claude-code --dry-run
shopwell-cli ai list --installed
shopwell-cli ai remove shopwell-cli --agent claude-code
```

A project install resolves the Shopwell project root (so the skill is written
where agents read it). "Project" here means a Shopwell shop (a `bin/console`
plus a `composer.json` requiring `shopwell/core`); a project install refuses to
run anywhere else. An extension (plugin or app) repository is not a shop, so
install the general `shopwell-cli` skill with `--global` there — which is the
natural scope for it anyway, making it available in every repository. The agent
name is whatever skills.sh supports (e.g. `claude-code`, `codex`); the CLI keeps
no record of its own — skills.sh is the source of truth for what is installed.
Commit `.claude/skills/` (or the agent's directory) and `skills-lock.json` to
share an integration with your team, or gitignore them to keep it local.

### MCP integrations

An integration of type `mcp` (for example `meteor`, the Meteor MCP server) is a
remote server rather than a skill. The CLI configures it by writing the
project's `.mcp.json` — the standard `{"mcpServers": {…}}` file MCP clients read
— instead of any agent's own config. Because the CLI owns that one file, `add`
and `remove` target exactly the named server, and no per-agent config paths need
to be known.

An MCP install is project-scoped: `.mcp.json` lives in the Shopwell project
root, and committing it shares the server with your team. It takes no `@<tag>`
and no `--agent` (the file is not agent-specific), and `--global` is not
supported.

```bash
shopwell-cli ai info meteor
shopwell-cli ai add meteor               # writes <project>/.mcp.json
shopwell-cli ai add meteor --dry-run
shopwell-cli ai remove meteor
```

`add` writes a remote `{"type":"http","url":…}` entry under `mcpServers`,
preserving other servers and unrelated keys; `remove` deletes only the named
server. `ai list --installed` reads the project's `.mcp.json` back and reports
configured MCP servers alongside skills, with a `Type` column.

Claude Code reads `.mcp.json` from the project root natively. Clients that use a
different project path (for example Cursor's `.cursor/mcp.json` or VS Code's
`.vscode/mcp.json`) can import it; if broader coverage is needed later, the
server can be distributed with a dedicated MCP installer instead.

### Recommendation

When a supported AI client runs the CLI (detected from the `AI_AGENT`
environment variable), and the `shopwell-cli` skill is not yet installed for it,
the CLI prints a short, non-blocking recommendation to stderr pointing at
`shopwell-cli ai add shopwell-cli --agent <client> --global`. It never installs
anything and asks the agent to confirm with the user first. The detection itself
is the trigger, so it is shown in non-interactive runs too (that is how agents
invoke the CLI). It is not shown again for the same client within 24 hours, and
stays silent in CI runs, during shell completion, while using `ai` commands, and
when `--no-ai-hint` is passed.

### Known limitations

- `ai remove` asks skills.sh whether the skill is installed before deleting it,
  so it reports `Removed` only when there was something to remove. The Codex
  agent can still leave files behind on a project-scope `ai remove` (a skills.sh
  defect); the CLI reports the removal that skills.sh carried out.

## Updating installed skills

The canonical source changes whenever `skills/` changes on the default branch.

Users can check installed skills for available updates with:

```bash
npx skills check
```

Update the Shopwell skills with:

```bash
npx skills update shopwell-cli shopwell-cli-docker shopwell-cli-extension-store shopwell-cli-extension-anonymize
```

or update all installed project skills with:

```bash
npx skills update -y
```

Updating the Shopwell CLI binary does not automatically modify skills installed
by external Agent Skills tooling.

## Keeping the skills current

Agent Skills are part of the user-facing Shopwell CLI contract.

Every pull request that changes user-facing CLI behavior must review `skills/`
for impact.

Review the skills when changing:

- command names or hierarchy;
- important execution modes;
- environment or Docker behavior;
- recommended workflows;
- destructive or state-changing behavior;
- commands explicitly referenced in a skill.

Do not turn the skills into a duplicate command reference. The running CLI and
its `--help` output remain authoritative for exact commands, flags, and
version-specific behavior.

Where possible, CI should verify that commands explicitly referenced by the
skills continue to exist.

## Versioning

Skills do not have an independent Shopwell version.

Because their canonical source lives in the Shopwell CLI repository, every Git
tag and Shopwell CLI release records the exact skill source that existed for
that release.

The default branch contains the latest maintained guidance.

## Future signed distribution

If Shopwell later requires cryptographically verified, pinned, or offline skill
distribution, the canonical `skills/` directory may additionally be published
as a signed OCI artifact.

This must remain a second distribution of the same canonical source, not a
separate copy of the skills.

Until such a requirement exists, the `skills` CLI and skills.sh ecosystem are
the preferred cross-client distribution and update mechanism.
