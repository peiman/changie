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
    - [`major`, `minor`, `patch` Commands](#major-minor-patch-commands)
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
- **Flexible Output Format**: Generate consistent, well-formatted changelog files

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
   changie minor
   ```

---

## Features

- **Project Initialization**: Generate a properly structured CHANGELOG.md file
- **Entry Management**: Add standardized changelog entries by type (added, changed, fixed, etc.)
- **Version Control**: Bump versions following Semantic Versioning (major, minor, patch)
- **Git Integration**: Commit changes and create version tags automatically
- **Structured Output**: Ensure consistency and readability of changelog files

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
   changie minor
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

### `major`, `minor`, `patch` Commands

Bump the version number according to Semantic Versioning rules.

#### Usage

```bash
changie [major|minor|patch] [flags]
```

#### Flags

- `--file`: Changelog file name (default: "CHANGELOG.md")
- `--rrp`: Remote repository provider (github, bitbucket) (default: "github")
- `--auto-push`: Automatically push changes and tags

#### Examples

```bash
changie major
changie minor --auto-push
changie patch --file HISTORY.md
```

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
