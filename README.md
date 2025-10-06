# changie

[![Build Status](https://github.com/peiman/changie/actions/workflows/ci.yml/badge.svg)](https://github.com/peiman/changie/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/codecov/c/github/peiman/changie)](https://codecov.io/gh/peiman/changie)
[![Go Report Card](https://goreportcard.com/badge/github.com/peiman/changie)](https://goreportcard.com/report/github.com/peiman/changie)
[![Version](https://img.shields.io/github/v/release/peiman/changie)](https://github.com/peiman/changie/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/peiman/changie.svg)](https://pkg.go.dev/github.com/peiman/changie)
[![License](https://img.shields.io/github/license/peiman/changie)](LICENSE)
[![CodeQL](https://github.com/peiman/changie/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/peiman/changie/security/code-scanning)
[![Made with Go](https://img.shields.io/badge/made%20with-Go-brightgreen.svg)](https://go.dev)

**A professional Golang CLI tool for managing changelogs following the "Keep a Changelog" format and Semantic Versioning.**

---

## Table of Contents

- [changie](#changie)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Key Highlights](#key-highlights)
  - [Quick Start](#quick-start)
  - [Features](#features)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Using changie](#using-changie)
    - [Important: Single Source of Truth for Names](#important-single-source-of-truth-for-names)
  - [Commands](#commands)
    - [`init` Command](#init-command)
      - [Usage](#usage)
      - [Flags](#flags)
      - [Examples](#examples)
    - [`changelog` Command](#changelog-command)
      - [Usage](#usage-1)
      - [Subcommands](#subcommands)
      - [Flags](#flags-1)
      - [Examples](#examples-1)
    - [`bump` Command](#bump-command)
      - [Usage](#usage-2)
      - [Flags](#flags-2)
      - [Examples](#examples-2)
  - [Configuration](#configuration)
    - [Configuration File](#configuration-file)
    - [Environment Variables](#environment-variables)
    - [Command-Line Flags](#command-line-flags)
  - [Development Workflow](#development-workflow)
    - [Taskfile Tasks](#taskfile-tasks)
    - [Pre-Commit Hooks with Lefthook](#pre-commit-hooks-with-lefthook)
    - [Continuous Integration](#continuous-integration)
  - [Dependency Management](#dependency-management)
    - [Available Tasks](#available-tasks)
    - [Automated Checks](#automated-checks)
    - [Best Practices](#best-practices)
  - [Contributing](#contributing)
  - [License](#license)
  - [Additional Notes](#additional-notes)
  
---

## Introduction

**changie** is a professional Go command-line application designed to help developers manage changelogs according to the [Keep a Changelog](https://keepachangelog.com/) format and [Semantic Versioning](https://semver.org/) principles. It provides a structured workflow for adding, organizing, and releasing changelog entries while integrating seamlessly with Git.

Built on solid engineering principles, changie includes:

- Modular command structure with [Cobra](https://github.com/spf13/cobra)
- Configuration management via [Viper](https://github.com/spf13/viper)
- Structured logging with [Zerolog](https://github.com/rs/zerolog)
- Comprehensive testing and code quality checks

---

## Key Highlights

- **Standardized Changelog Management**: Follow "Keep a Changelog" best practices without manual formatting
- **Semantic Versioning Support**: Automatic version bumping following SemVer principles
- **Git Integration**: Seamless interaction with Git for tagging and committing changes
- **JSON Output**: Machine-readable output with `--json` flag for automation and CI/CD integration
- **MCP Server**: Model Context Protocol server for AI agent integration (Claude Desktop, etc.)
- **AI-Ready**: Comprehensive documentation and tools for AI agents (`llms.txt`, `.ai/` directory)

---

## Quick Start

1. **Install changie**:

   ```bash
   go install github.com/peiman/changie@latest
   ```

2. **Initialize a project**:

   ```bash
   changie init
   ```

3. **Add a changelog entry**:

   ```bash
   changie changelog added "New feature: added user authentication"
   ```

4. **Release a new version**:

   ```bash
   changie bump minor
   ```

---

## Features

- **Project Initialization**: Generate a properly structured CHANGELOG.md file
- **Entry Management**: Add standardized changelog entries by type (added, changed, fixed, etc.)
- **Version Control**: Bump versions following Semantic Versioning (major, minor, patch)
- **Git Integration**: Commit changes and create version tags automatically
- **JSON Output Mode**: Machine-readable output for automation and CI/CD pipelines
- **MCP Server**: AI agent integration via Model Context Protocol
- **AI Agent Resources**: Comprehensive guides, prompts, and workflows in `.ai/` directory
- **LLM-Optimized Docs**: `llms.txt` file for AI assistant integration
- **Usage Examples**: Ready-to-use scripts for common workflows and CI/CD integration

---

## Getting Started

### Prerequisites

- **Go**: 1.20+ recommended
- **Git**: For version control and integration features

### Installation

```bash
go install github.com/peiman/changie@latest
```

Or build from source:

```bash
git clone https://github.com/peiman/changie.git
cd changie
go install
```

### Using changie

1. **Initialize a project**:

   ```bash
   changie init
   ```

   This creates a `CHANGELOG.md` file in your project root.

2. **Add a changelog entry**:

   ```bash
   changie changelog added "New feature: added user authentication"
   ```

3. **Release a new version**:

   ```bash
   changie bump minor
   ```

   This will bump the minor version number and update the changelog.

### Important: Single Source of Truth for Names

This project uses a "single source of truth" approach for configuration:

1. **Binary Name**: Defined only in `Taskfile.yml` as `BINARY_NAME`. This is propagated through the codebase via build flags and the `binaryName` variable in `cmd/root.go`.

2. **Module Path**: Defined only in `go.mod` and referenced in `Taskfile.yml` as `MODULE_PATH`.

When customizing this project:

- Change `BINARY_NAME` in `Taskfile.yml` to your desired binary name
- Change the module path in `go.mod` to your own repository path
- Run `task build` to apply these changes throughout the codebase

---

## Commands

### `init` Command

Initialize a project with a properly formatted CHANGELOG.md file.

#### Usage

```bash
changie init [flags]
```

#### Flags

- `--file`: Changelog file name (default: "CHANGELOG.md")

#### Examples

```bash
changie init
changie init --file HISTORY.md
```

### `changelog` Command

Add entries to different sections of the changelog.

#### Usage

```bash
changie changelog [subcommand] [content]
```

#### Subcommands

- `added`: Add entry to the Added section
- `changed`: Add entry to the Changed section
- `deprecated`: Add entry to the Deprecated section
- `removed`: Add entry to the Removed section
- `fixed`: Add entry to the Fixed section
- `security`: Add entry to the Security section

#### Flags

- `--file`: Changelog file name (default: "CHANGELOG.md")

#### Examples

```bash
changie changelog added "New feature: added user authentication"
changie changelog fixed "Bug in login form"
changie changelog security "Patched XSS vulnerability"
```

### `bump` Command

Bump the version number according to Semantic Versioning rules.

#### Usage

```bash
changie bump [major|minor|patch] [flags]
```

#### Subcommands

- `major`: Bump the major version (breaking changes, X.y.z → X+1.0.0)
- `minor`: Bump the minor version (new features, x.Y.z → x.Y+1.0)
- `patch`: Bump the patch version (bug fixes, x.y.Z → x.y.Z+1)

#### Flags

- `--file`: Changelog file name (default: "CHANGELOG.md")
- `--rrp`: Remote repository provider (github, bitbucket) (default: "github")
- `--auto-push`: Automatically push changes and tags
- `--allow-any-branch`: Allow version bumping on any branch (bypasses main/master branch check)
- `--json`: Output results in JSON format for machine parsing

#### Examples

```bash
changie bump major
changie bump minor --auto-push
changie bump patch --file HISTORY.md
changie bump minor --allow-any-branch  # For release branches or hotfixes

# JSON output for automation
changie bump patch --json
# Returns: {"success":true,"old_version":"1.0.0","new_version":"1.0.1",...}
```

**Note:** By default, version bump commands check that you're on the `main` or `master` branch. This is a best practice to maintain a clean release history. Use `--allow-any-branch` when you need to bump versions on other branches (e.g., release branches, hotfix branches).

---

## Configuration

changie uses Viper for flexible configuration:

### Configuration File

Default config file: `$HOME/.changie.yaml`

Example:

```yaml
app:
  log_level: "info"
  changelog:
    file: "CHANGELOG.md"
  version:
    tag_prefix: "v"
```

### Environment Variables

Override any config via environment variables:

```bash
export APP_LOG_LEVEL="debug"
export APP_CHANGELOG_FILE="HISTORY.md"
```

### Command-Line Flags

Override at runtime:

```bash
changie init --file HISTORY.md
```

---

## JSON Output & Automation

changie supports machine-readable JSON output for all bump commands, making it easy to integrate with CI/CD pipelines, scripts, and automation tools.

### Using JSON Output

Add the `--json` flag to any bump command:

```bash
changie bump patch --json
```

**Example Output:**
```json
{
  "success": true,
  "old_version": "1.2.3",
  "new_version": "1.2.4",
  "tag": "v1.2.4",
  "changelog_file": "CHANGELOG.md",
  "commit_hash": "abc123",
  "pushed": false,
  "bump_type": "patch"
}
```

**Error Output:**
```json
{
  "success": false,
  "error": "uncommitted changes detected - commit or stash changes first",
  "bump_type": "patch"
}
```

### CI/CD Integration Examples

**GitHub Actions:**
```yaml
- name: Bump version
  id: bump
  run: |
    OUTPUT=$(changie bump patch --json)
    echo "version=$(echo $OUTPUT | jq -r '.new_version')" >> $GITHUB_OUTPUT

- name: Use version
  run: echo "Released version ${{ steps.bump.outputs.version }}"
```

**GitLab CI:**
```yaml
bump_version:
  script:
    - changie bump patch --json | tee release.json
    - export VERSION=$(jq -r '.new_version' release.json)
    - echo "VERSION=$VERSION" >> variables.env
  artifacts:
    reports:
      dotenv: variables.env
```

**Shell Scripts:**
```bash
#!/bin/bash
RESULT=$(changie bump patch --json)
if [[ $(echo "$RESULT" | jq -r '.success') == "true" ]]; then
  VERSION=$(echo "$RESULT" | jq -r '.new_version')
  echo "✓ Released version $VERSION"
else
  ERROR=$(echo "$RESULT" | jq -r '.error')
  echo "✗ Release failed: $ERROR"
  exit 1
fi
```

See the `examples/` directory for more complete CI/CD integration examples.

---

## MCP Server (AI Agent Integration)

changie includes an **MCP (Model Context Protocol) server** that exposes changelog operations as tools for AI agents like Claude Desktop.

### What is MCP?

[Model Context Protocol](https://modelcontextprotocol.io/) is an open standard for connecting AI assistants to external tools and data sources. changie's MCP server allows AI agents to manage your changelog autonomously.

### Available Tools

The MCP server exposes 4 tools:

1. **`changie_bump_version`** - Bump semantic version (major/minor/patch)
2. **`changie_add_changelog`** - Add changelog entries
3. **`changie_init`** - Initialize new projects
4. **`changie_get_version`** - Get current version from git tags

### Running the MCP Server

**Option 1: Build from source**
```bash
# Build MCP server binary
task build:mcp

# Run server (stdio mode for MCP)
./changie-mcp-server
```

**Option 2: Docker (recommended)**
```bash
# Build Docker image
task docker:build:mcp

# Run in Docker
task docker:run:mcp
```

**Option 3: With Claude Desktop**

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "changie": {
      "command": "/path/to/changie-mcp-server"
    }
  }
}
```

Then restart Claude Desktop. The changie tools will appear in Claude's available tools.

### MCP Architecture

- **Transport**: stdio (standard for MCP)
- **Protocol**: JSON-RPC 2.0
- **SDK**: Official MCP Go SDK v1.0.0
- **Integration**: Calls `changie` binary with `--json` flag

### Example Usage with AI Agents

Once configured, you can ask Claude:

> "Please bump the patch version and push the changes"

Claude will use the `changie_bump_version` tool to execute the release.

---

## AI Agent Resources

changie includes comprehensive resources for AI agents and LLM integration:

### `.ai/` Directory

Structured resources for AI agents:

- **`context.md`** - Project capabilities, constraints, and best practices
- **`prompts/release-new-version.md`** - Guide for releasing versions
- **`prompts/add-changelog-entry.md`** - Guide for adding changelog entries
- **`workflows/complete-release-workflow.md`** - End-to-end release process

These files help AI agents understand changie's capabilities and execute tasks correctly.

### `llms.txt`

LLM-optimized documentation following the [llms.txt standard](https://llmstxt.org/). This file provides:

- Quick reference for all commands
- Decision trees for choosing bump types
- JSON output examples
- Common workflows and error handling
- Integration patterns

AI assistants can reference this file to understand how to use changie effectively.

### `examples/` Directory

Ready-to-use scripts and comprehensive guides:

- **`basic-workflow.sh`** - Daily development workflow
- **`ci-integration.sh`** - GitHub Actions and GitLab CI examples
- **`release-workflow.sh`** - Complete release automation
- **`README.md`** - Detailed usage guide with troubleshooting

These examples serve as templates for both humans and AI agents.

---

## Development Workflow

### Taskfile Tasks

- `task setup`: Install tools
- `task format`: Format code
- `task lint`: Run linters
- `task test`: Run tests with coverage
- `task build`: Build the binary
- `task run`: Run the binary
- `task check`: All checks (format, lint, deps, tests)

### Pre-Commit Hooks with Lefthook

`task setup` installs hooks that run `format`, `lint`, `test` on commit, ensuring code quality before changes land in the repository.

### Continuous Integration

GitHub Actions runs `task check` on each commit or pull request, maintaining code standards and reliability.

---

## Dependency Management

### Available Tasks

- `task deps:verify`: Verifies that dependencies haven't been modified
- `task deps:outdated`: Checks for outdated dependencies
- `task deps:check`: Runs all dependency checks (verification, outdated, vulnerabilities)

### Automated Checks

Dependency verification is automatically included in:

- Pre-commit hooks via Lefthook
- CI workflow via GitHub Actions
- The comprehensive quality check command: `task check`

### Best Practices

1. Run `task deps:check` before starting a new feature
2. Update dependencies incrementally with `go get -u <package>` followed by `task tidy`
3. Always run tests after dependency updates
4. Document significant dependency changes in commit messages

---

## Contributing

1. Fork & create a new branch
2. Make changes, run `task check`
3. Commit with descriptive messages following the project's commit convention
4. Open a pull request against `main`

---

## License

MIT License. See [LICENSE](LICENSE).

---

## Additional Notes

- Run `task test:coverage-text` to identify uncovered code paths for targeted testing improvements
- Regularly run `task deps:check` to ensure dependencies are up-to-date and secure
- For consistent formatting, run `task format` before committing changes

---
