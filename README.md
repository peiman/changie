# changie

A CLI tool for managing semantic versioning and [Keep a Changelog](https://keepachangelog.com) format changelogs.

[![Build Status](https://github.com/peiman/changie/actions/workflows/ci.yml/badge.svg)](https://github.com/peiman/changie/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/peiman/changie)](https://goreportcard.com/report/github.com/peiman/changie)
[![License](https://img.shields.io/github/license/peiman/changie)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/peiman/changie)](go.mod)

---

## What it does

changie automates the tedious parts of releasing software:

- **Initializes** a project with a properly formatted CHANGELOG.md
- **Adds** changelog entries to the correct section (Added, Changed, Fixed, etc.)
- **Bumps** versions following [Semantic Versioning](https://semver.org) (major, minor, patch)
- **Updates** the changelog, commits, tags, and optionally pushes -- all in one command

No more hand-editing changelogs or forgetting to tag releases.

## Install

### From source

```bash
go install github.com/peiman/changie@latest
```

### Build from source

```bash
git clone https://github.com/peiman/changie.git
cd changie
task build
```

Requires [Go 1.26+](https://go.dev/dl/) and [Task](https://taskfile.dev).

## Quick start

```bash
# Initialize a new project (creates CHANGELOG.md and v0.0.0 tag)
changie init

# Add changelog entries as you work
changie changelog added "User authentication via OAuth2"
changie changelog fixed "Race condition in session handler"

# Release a new version
changie bump patch    # 0.0.0 -> 0.0.1
changie bump minor    # 0.0.1 -> 0.1.0
changie bump major    # 0.1.0 -> 1.0.0
```

That's it. changie updates the changelog, commits the change, creates a git tag, and reminds you to push.

## Commands

### `changie init`

Creates a CHANGELOG.md following the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format. If you're in a git repo with no tags, it also creates an initial `v0.0.0` tag.

```bash
changie init
changie init --file HISTORY.md         # Use a different filename
changie init --use-v-prefix=false      # Tags without 'v' prefix (1.0.0 instead of v1.0.0)
```

If existing git tags are found, changie adopts their naming convention automatically.

### `changie changelog <section> <content>`

Adds an entry to a section in the `[Unreleased]` area of your changelog.

Valid sections: `added`, `changed`, `deprecated`, `removed`, `fixed`, `security`

```bash
changie changelog added "New REST API endpoints"
changie changelog changed "Upgraded database driver to v3"
changie changelog deprecated "Legacy XML export format"
changie changelog removed "Python 2 support"
changie changelog fixed "Memory leak in connection pool"
changie changelog security "Patched XSS vulnerability in search"
```

Duplicate entries are detected and skipped.

### `changie bump <major|minor|patch>`

Performs a complete version release:

1. Verifies you're on main/master (configurable)
2. Checks for uncommitted changes
3. Reads the current version from git tags
4. Calculates the new version
5. Updates the changelog (moves `[Unreleased]` content to the new version)
6. Commits the changelog update
7. Creates a git tag
8. Optionally pushes to remote

```bash
changie bump patch                     # Bug fix release
changie bump minor                     # Feature release
changie bump major                     # Breaking change release

changie bump patch --auto-push         # Push automatically after bumping
changie bump minor --allow-any-branch  # Bump from any branch, not just main/master
changie bump patch --use-v-prefix=false  # Tag as 1.2.4 instead of v1.2.4
changie bump major --rrp gitlab        # Use GitLab-style comparison links
```

### `changie config validate`

Validates the current configuration.

```bash
changie config validate
```

## Configuration

changie uses [Viper](https://github.com/spf13/viper) for configuration. Precedence (highest to lowest):

1. Command-line flags
2. Environment variables
3. Config file
4. Defaults

### Config file

Place a `config.yaml` in one of these locations:

- `./config.yaml` (project directory)
- `~/.config/changie/config.yaml` (XDG config directory)

```yaml
app:
  log_level: info
  changelog:
    file: CHANGELOG.md
    repository_provider: github     # github, gitlab, bitbucket
  version:
    use_v_prefix: true
    auto_push: false
```

### Environment variables

All config keys map to environment variables with a `CHANGIE_` prefix:

```bash
export CHANGIE_APP_LOG_LEVEL=debug
export CHANGIE_APP_CHANGELOG_FILE=HISTORY.md
export CHANGIE_APP_VERSION_USE_V_PREFIX=false
export CHANGIE_APP_VERSION_AUTO_PUSH=true
```

### JSON output

For CI/CD pipelines and automation, use `--output json` to get machine-readable output:

```bash
changie bump patch --output json
```

## Changelog format

changie follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and generates files like:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- New feature description

## [1.1.0] - 2025-03-15

### Added

- OAuth2 authentication
- Rate limiting

### Fixed

- Memory leak in connection pool

## [1.0.0] - 2025-01-10

### Added

- Initial release

[Unreleased]: https://github.com/user/repo/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/user/repo/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/user/repo/releases/tag/v1.0.0
```

Comparison links are generated automatically based on your git remote (supports GitHub, GitLab, and Bitbucket).

## Development

changie is built on the [ckeletin-go](https://github.com/peiman/ckeletin-go) framework.

```bash
task setup     # Install development tools
task test      # Run tests with coverage
task check     # Run all quality checks (lint, vet, security, architecture)
task build     # Build the binary
task format    # Format code
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## License

MIT License. See [LICENSE](LICENSE).
