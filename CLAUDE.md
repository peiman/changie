# Claude Code - Project Guide

This document provides guidance for AI agents (Claude Code) and developers working on this project. It outlines the architectural principles, best practices, and project-specific conventions being followed.

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture Principles](#architecture-principles)
- [Go Best Practices](#go-best-practices)
- [ckeletin-go Scaffold](#ckeletin-go-scaffold)
- [Project Structure](#project-structure)
- [Key Conventions](#key-conventions)
- [Testing](#testing)
- [Common Tasks](#common-tasks)
- [Where to Find Things](#where-to-find-things)

## Project Overview

**changie** is a professional changelog management CLI tool for SemVer projects that follows the Keep a Changelog format.

- **Language**: Go
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Config Management**: [Viper](https://github.com/spf13/viper)
- **Logging**: [zerolog](https://github.com/rs/zerolog/log)
- **Testing**: Standard Go testing with [gotestsum](https://github.com/gotestyourself/gotestsum)
- **Build Automation**: [Task](https://taskfile.dev/)

## Architecture Principles

### 1. **No "Service Layer" Pattern**

This project follows Go idioms and **avoids** traditional OOP patterns like service layers, repository patterns, or dependency injection containers.

**Instead, we use:**
- Plain functions grouped by domain/feature in `internal/` packages
- Direct function calls (no unnecessary interfaces)
- Dependency injection via function parameters (e.g., `io.Writer` for testability)

### 2. **Thin cmd/ Layer**

Command files in `cmd/` should be **thin wrappers** that:
- Define Cobra commands and flags
- Bind flags to Viper configuration
- Map configuration values to structs
- Call business logic functions from `internal/`
- Handle I/O (stdout, stderr)

**cmd/ files should NOT contain:**
- Complex business logic
- Workflow orchestration
- Data transformation
- Validation logic (beyond flag validation)

### 3. **Business Logic in internal/**

All business logic resides in `internal/` packages as **plain, testable functions**.

**Example:**
```go
// Good: internal/version/bump.go
func Bump(cfg BumpConfig, output io.Writer) error {
    // Business logic here
}

// Bad: cmd/version.go
func runVersionBump(cmd *cobra.Command, bumpType string) error {
    // 200 lines of business logic mixed with CLI concerns
}
```

## Go Best Practices

This project follows established Go community standards and idioms.

### Primary Sources

1. **[golang-standards/project-layout](https://github.com/golang-standards/project-layout)**
   - The de facto standard for Go project structure
   - Defines `/cmd`, `/internal`, `/pkg` conventions
   - Enforced by Go compiler (internal packages)

2. **[Effective Go](https://go.dev/doc/effective_go)**
   - Official Go documentation on writing clear, idiomatic Go code
   - Covers naming, formatting, control structures, functions, data structures

3. **[Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)**
   - Common comments made during reviews of Go code
   - Best practices from the Go team

4. **[Standard Go Project Layout](https://go.dev/doc/modules/layout)**
   - Official Go module organization guidance
   - How to structure packages and modules

### Key Articles

- **[Common Anti-Patterns in Go Web Applications](https://threedots.tech/post/common-anti-patterns-in-go-web-applications/)** by Three Dots Labs
  - Avoid single model with multiple tags (JSON, DB, validation)
  - Don't start with database schema - start with domain
  - Separate models for different layers (HTTP, storage, application)

- **[Why Clean Architecture Struggles in Golang](https://dev.to/lucasdeataides/why-clean-architecture-struggles-in-golang-and-what-works-better-m4g)**
  - Go favors practicality over rigid architectural patterns
  - Avoid over-engineering with excessive layers
  - Focus on packages and modularity, not strict layering

- **[The Fat Service Pattern](https://www.alexedwards.net/blog/the-fat-service-pattern)** by Alex Edwards
  - Split code into application layer (HTTP/CLI) and logic layer
  - Keep handlers lightweight, focused on I/O
  - Business logic in separate packages

### Anti-Patterns to Avoid

❌ **Service classes with methods**
```go
// Don't do this
type VersionService struct {
    gitRepo GitRepository
    changelog ChangelogRepository
}
func (s *VersionService) BumpVersion(...) error { }
```

✅ **Plain functions grouped by domain**
```go
// Do this instead
package version

func Bump(cfg BumpConfig, output io.Writer) error { }
```

❌ **Business logic in cmd/ files**
```go
// Don't do this in cmd/version.go
func runVersionBump(...) error {
    // 200 lines of git operations, validation, changelog updates
}
```

✅ **cmd/ as thin wrappers**
```go
// Do this instead in cmd/version.go
func runVersionBump(cmd *cobra.Command, bumpType string) error {
    cfg := version.BumpConfig{ /* map from viper */ }
    return version.Bump(cfg, cmd.OutOrStdout())
}
```

## ckeletin-go Scaffold

This project is built upon the **[ckeletin-go](https://github.com/peiman/ckeletin-go)** scaffold, a professional Go CLI application template.

### Scaffold Repository
**https://github.com/peiman/ckeletin-go**

### Key Scaffold Principles

1. **Modular Configuration Management**
   - Single source of truth: `internal/config/registry.go`
   - Command-specific config files: `internal/config/*_options.go`
   - **NEVER** use `viper.SetDefault()` outside of registry
   - Configuration precedence: CLI flags > Env vars > Config file > Defaults

2. **Options Pattern**
   - Each command has its own `*_options.go` file in `internal/config/`
   - Registry aggregates all options in one place
   - Type-safe configuration handling

3. **Structured Logging**
   - Use `zerolog` for all logging
   - Log at appropriate levels (Debug, Info, Warn, Error)
   - Include context with `.Str()`, `.Err()`, etc.

4. **Documentation Generation**
   - Auto-generate config documentation with `changie docs config`
   - Keep config options well-documented in code
   - Support multiple output formats (markdown, yaml)

### Scaffold Structure

```
.
├── cmd/                    # Command definitions (thin wrappers)
│   ├── root.go            # Root command and config initialization
│   ├── version.go         # Version bump commands
│   └── *.go               # Other commands
├── internal/              # Private application code
│   ├── config/            # Configuration management
│   │   ├── registry.go    # Single source of truth for config
│   │   ├── core_options.go
│   │   ├── *_options.go   # Command-specific configs
│   ├── <domain>/          # Domain-specific packages
│   │   └── *.go           # Business logic as plain functions
│   ├── ui/                # User interface utilities
│   └── logger/            # Logging configuration
├── scripts/               # Build and automation scripts
├── Taskfile.yml           # Task automation
└── .golangci.yml          # Linter configuration
```

## Project Structure

### `/cmd` - Command Layer

Entry points for the CLI. Each file defines cobra commands.

**Files:**
- `root.go` - Root command, config initialization, global flags (includes `--json` flag)
- `version.go` - Version bump commands (major, minor, patch)
- `init.go` - Project initialization command
- `changelog.go` - Changelog command group
- `changelog_add.go` - Subcommands for adding changelog entries
- `docs.go` - Documentation generation commands
- `completion.go` - Shell completion
- `mcp-server/main.go` - MCP (Model Context Protocol) server entry point

**Responsibilities:**
- Define cobra commands and flags
- Bind flags to viper
- Map config to structs
- Call internal/ functions
- Handle stdout/stderr

### `/internal` - Business Logic

Private packages that cannot be imported by external projects.

#### `internal/version/`
Version bump orchestration logic.

**Files:**
- `bump.go` - `Bump()` function for complete version bump workflow

**Key Functions:**
- `Bump(cfg BumpConfig, output io.Writer) error` - Orchestrates version bumping

#### `internal/config/`
Configuration management following the registry pattern.

**Files:**
- `registry.go` - Single source of truth, aggregates all config options
- `core_options.go` - Application-wide settings
- `version_options.go` - Version command settings
- `init_options.go` - Init command settings
- `docs_options.go` - Docs command settings
- `environment.go` - Environment variable utilities
- `paths.go` - Config file path utilities
- `options.go` - ConfigOption type definition

**Key Functions:**
- `Registry() []ConfigOption` - Returns all config options
- `SetDefaults()` - Applies defaults to viper
- `EnvPrefix(binaryName string) string` - Sanitizes env var prefix
- `DefaultPaths(binaryName string) PathsConfig` - Returns config paths

#### `internal/changelog/`
Changelog file management following Keep a Changelog format.

**Files:**
- `changelog.go` - Changelog operations

**Key Functions:**
- `InitProject(filePath string) error` - Creates new changelog
- `AddChangelogSection(filePath, section, content string) (bool, error)` - Adds entry
- `UpdateChangelog(filePath, version, provider string) error` - Updates for release
- `GetLatestChangelogVersion(content string) (string, error)` - Extracts version

#### `internal/git/`
Git operations wrapper.

**Files:**
- `git.go` - Git command wrappers

**Key Functions:**
- `IsInstalled() bool` - Checks if git is available
- `GetVersion() (string, error)` - Gets current version from tags
- `GetCurrentBranch() (string, error)` - Gets current branch name
- `HasUncommittedChanges() (bool, error)` - Checks for uncommitted changes
- `CommitChangelog(file, version string) error` - Commits changelog
- `TagVersion(version string) error` - Creates version tag
- `PushChanges() error` - Pushes commits and tags
- `GetRepositoryInfo() (*RepositoryInfo, error)` - Parses remote URL

#### `internal/semver/`
Semantic versioning operations.

**Files:**
- `semver.go` - SemVer parsing and bumping

**Key Functions:**
- `ParseVersion(version string) (semver.Version, bool, error)` - Parses version
- `BumpMajor/Minor/Patch(version string, useVPrefix bool) (string, error)` - Bumps version
- `Compare(v1, v2 string) (int, error)` - Compares versions

#### `internal/ui/`
User interface utilities and prompts.

**Files:**
- `ui.go` - UI interface definition
- `message.go` - Message formatting
- `colors.go` - Color definitions
- `prompt.go` - User interaction functions
- `mock.go` - Mock UI for testing

**Key Functions:**
- `AskYesNo(prompt string, defaultYes bool, output io.Writer) (bool, error)` - Y/N prompt

#### `internal/logger/`
Logging configuration and initialization.

**Files:**
- `logger.go` - Logger setup

**Key Functions:**
- `Init(options *Options) error` - Initializes zerolog

#### `internal/docs/`
Documentation generation.

**Files:**
- `config.go` - Config struct
- `generator.go` - Documentation generator
- `markdown.go` - Markdown formatter
- `yaml.go` - YAML formatter

**Key Functions:**
- `NewGenerator(cfg *Config) *Generator` - Creates generator
- `Generate() error` - Generates documentation

#### `internal/output/`
JSON output utilities for machine-readable command results.

**Files:**
- `output.go` - JSON output formatting and structures

**Key Functions:**
- `IsJSONEnabled() bool` - Checks if JSON output mode is enabled
- `WriteJSON(w io.Writer, v interface{}) error` - Writes JSON to writer
- `Write(w io.Writer, textValue string, jsonValue interface{}) error` - Writes in appropriate format

**Output Structures:**
- `BumpOutput` - Version bump results (success, old_version, new_version, tag, etc.)
- `InitOutput` - Project initialization results
- `ChangelogOutput` - Changelog operation results
- `DocsOutput` - Documentation generation results

#### `internal/mcp/`
MCP (Model Context Protocol) server tool implementations.

**Files:**
- `tools.go` - MCP tool handler functions
- `tools_test.go` - Comprehensive integration tests

**Key Functions:**
- `BumpVersion(ctx, req, input) (result, output, error)` - Bump version MCP tool
- `AddChangelog(ctx, req, input) (result, output, error)` - Add changelog entry tool
- `Init(ctx, req, input) (result, output, error)` - Initialize project tool
- `GetVersion(ctx, req, input) (result, output, error)` - Get current version tool

**Note:** Uses official MCP Go SDK v1.0.0 (`github.com/modelcontextprotocol/go-sdk`). Implements 4 tools that call `changie` with `--json` flag for structured output.

## Key Conventions

### Configuration Management

**✅ DO:**
```go
// In internal/config/version_options.go
func VersionOptions() []ConfigOption {
    return []ConfigOption{
        {
            Key:          "app.version.use_v_prefix",
            DefaultValue: true,
            Description:  "Use 'v' prefix for version tags",
        },
    }
}

// In internal/config/registry.go
func Registry() []ConfigOption {
    options := CoreOptions()
    options = append(options, VersionOptions()...)
    return options
}
```

**❌ DON'T:**
```go
// In cmd/version.go - NEVER do this!
viper.SetDefault("some.config", "value")
```

### Error Handling

**✅ DO:**
```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to update changelog: %w - verify that '%s' exists", err, file)
}

// Log errors before returning
log.Error().Err(err).Str("file", file).Msg("Failed to update changelog")
return fmt.Errorf("failed to update changelog: %w", err)
```

**❌ DON'T:**
```go
// Don't return generic errors without context
if err != nil {
    return err
}

// Don't lose error information
if err != nil {
    return fmt.Errorf("something failed")
}
```

### Function Signatures

**For testability, pass `io.Writer` for output:**
```go
// Good - testable with any writer
func Bump(cfg BumpConfig, output io.Writer) error {
    fmt.Fprintf(output, "Current version: %s\n", version)
}

// Bad - hardcoded to stdout
func Bump(cfg BumpConfig) error {
    fmt.Printf("Current version: %s\n", version)
}
```

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `version`, `changelog`, not `versionBump`)
- **Files**: lowercase with underscores (e.g., `version_options.go`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase
- **Constants**: PascalCase or SCREAMING_SNAKE_CASE for package-level

## Testing

### Test Organization

- Test files alongside source: `*_test.go`
- Use table-driven tests for multiple cases
- Use test helpers to reduce duplication
- Mock external dependencies (git, filesystem) using interfaces or test fixtures

### Running Tests

```bash
# Run all tests
task test

# Run tests with coverage
go test -v -coverprofile=coverage.txt ./...

# Run specific package tests
go test -v ./internal/version/

# Run with race detector
go test -race ./...
```

### Test Coverage Requirements

Current coverage (as of latest commit):
- Overall: **83.8%**
- `internal/changelog`: 86.3%
- `internal/config`: 51.4%
- `internal/docs`: 89.3%
- `internal/git`: 81.0%
- `internal/logger`: 100%
- `internal/semver`: 100%
- `internal/ui`: 54.1%
- `internal/version`: **0%** (needs tests!)
- `cmd`: 75.1%

**Priority**: Add tests for `internal/version/bump.go`

### Example Test Pattern

```go
func TestBump(t *testing.T) {
    tests := []struct {
        name    string
        cfg     version.BumpConfig
        want    string
        wantErr bool
    }{
        {
            name: "major bump",
            cfg: version.BumpConfig{
                BumpType: "major",
                // ...
            },
            want:    "v2.0.0",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var buf bytes.Buffer
            err := version.Bump(tt.cfg, &buf)
            // assertions...
        })
    }
}
```

## Common Tasks

### Adding a New Command

1. **Create command file**: `cmd/mycommand.go`
   ```go
   package cmd

   import (
       "github.com/spf13/cobra"
       "github.com/peiman/changie/internal/myfeature"
   )

   var myCmd = &cobra.Command{
       Use:   "mycommand",
       Short: "Description",
       RunE:  runMyCommand,
   }

   func init() {
       RootCmd.AddCommand(myCmd)
       // Add flags and bind to viper
   }

   func runMyCommand(cmd *cobra.Command, args []string) error {
       // Map config and call internal function
       return myfeature.DoSomething(cfg, cmd.OutOrStdout())
   }
   ```

2. **Create config options**: `internal/config/mycommand_options.go`
   ```go
   package config

   func MyCommandOptions() []ConfigOption {
       return []ConfigOption{
           {
               Key:          "app.mycommand.setting",
               DefaultValue: "default",
               Description:  "Description of setting",
           },
       }
   }
   ```

3. **Update registry**: `internal/config/registry.go`
   ```go
   func Registry() []ConfigOption {
       options := CoreOptions()
       options = append(options, MyCommandOptions()...)
       return options
   }
   ```

4. **Create business logic**: `internal/myfeature/myfeature.go`
   ```go
   package myfeature

   import "io"

   type Config struct {
       Setting string
   }

   func DoSomething(cfg Config, output io.Writer) error {
       // Business logic here
       return nil
   }
   ```

5. **Add tests**: `internal/myfeature/myfeature_test.go`

### Adding Configuration Option

1. Add to appropriate `*_options.go` file in `internal/config/`
2. Ensure it's included in `Registry()` in `registry.go`
3. Never use `viper.SetDefault()` outside of registry
4. Run `./scripts/check-defaults.sh` to verify

### Code Quality Checks

```bash
# Run all checks (format, lint, test, vuln scan)
task check

# Individual checks
task format      # Format code with goimports and gofmt
task lint        # Run golangci-lint
task test        # Run tests with coverage
task vuln        # Check for vulnerabilities
task deps-verify # Verify dependencies
```

### Pre-commit Hooks

Managed by [Lefthook](https://github.com/evilmartians/lefthook). Configured in `.lefthook.yml`.

**Runs automatically on commit:**
- check-defaults (no unauthorized viper.SetDefault calls)
- format (goimports, gofmt)
- deps-verify (go mod verify)
- lint (go vet, golangci-lint)
- test (full test suite with coverage)

## Where to Find Things

### "I need to add validation logic"
→ Create a function in the appropriate `internal/` package, not in `cmd/`

### "I need to add a new flag to a command"
→ Add flag in `cmd/*_command.go` `init()` function, bind to viper, add to config options

### "I need to change default configuration"
→ Update the appropriate `*_options.go` file in `internal/config/`

### "I need to add user interaction (prompts)"
→ Add function to `internal/ui/prompt.go`

### "I need to call git commands"
→ Use or extend functions in `internal/git/git.go`

### "I need to parse/manipulate versions"
→ Use functions in `internal/semver/semver.go`

### "I need to read/write changelog files"
→ Use functions in `internal/changelog/changelog.go`

### "I need to add logging"
→ Use `log` from `github.com/rs/zerolog/log` (already initialized)

### "I need to format output for users"
→ Use `fmt.Fprintf(output, ...)` where `output` is `io.Writer` (typically `cmd.OutOrStdout()`)

### "I need to add JSON output support to a command"
→ Use functions in `internal/output/output.go` - `Write()`, `WriteJSON()`, or create new output struct

### "I need to implement MCP tools for AI agents"
→ Add tool handler in `internal/mcp/tools.go`, register in `cmd/mcp-server/main.go`

### "I need AI agent documentation/prompts"
→ Check `.ai/` directory for context, prompts, and workflows; see `llms.txt` for LLM docs

### "I need CI/CD integration examples"
→ Look at `examples/ci-integration.sh` for GitHub Actions and GitLab CI patterns

### "I need to add environment variable support"
→ Use viper automatic env binding (already configured with prefix in `cmd/root.go`)

### "I need to understand the config system"
→ Read `internal/config/registry.go` and look at examples in `*_options.go` files

### "I need to see how tests are structured"
→ Look at `internal/semver/semver_test.go` or `cmd/version_test.go` for patterns

### "Command runs but doesn't respect config file"
→ Check viper binding in command's `init()` function

### "Linter is complaining about complexity"
→ Check if exception is needed in `.golangci.yml` exclude-rules section

### "Tests are failing in CI but pass locally"
→ Check if git operations are involved (they need real git repos in tests)

## Best Practices Checklist

When adding/modifying code, ensure:

- [ ] Business logic is in `internal/`, not `cmd/`
- [ ] cmd/ files are thin wrappers (< 200 lines ideally)
- [ ] Functions accept `io.Writer` for output (testability)
- [ ] Configuration defaults are in `internal/config/*_options.go`, not scattered
- [ ] No `viper.SetDefault()` calls outside of registry
- [ ] Errors are wrapped with context using `%w`
- [ ] Important operations are logged with zerolog
- [ ] New functions have documentation comments
- [ ] Tests are added for new functionality
- [ ] `task check` passes (format, lint, test, vuln)
- [ ] No service classes - use plain functions
- [ ] Packages are organized by domain/feature, not by layer

## Resources

### Project Resources
- **Main Repository**: https://github.com/peiman/changie
- **Scaffold Repository**: https://github.com/peiman/ckeletin-go
- **Keep a Changelog**: https://keepachangelog.com/
- **Semantic Versioning**: https://semver.org/
- **Model Context Protocol**: https://modelcontextprotocol.io/

### Go Standards & Best Practices
- **[golang-standards/project-layout](https://github.com/golang-standards/project-layout)** - Standard project structure
- **[Effective Go](https://go.dev/doc/effective_go)** - Official Go best practices
- **[Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)** - Review guidelines
- **[Common Anti-Patterns](https://threedots.tech/post/common-anti-patterns-in-go-web-applications/)** - What to avoid

### Tools Documentation
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework
- **[Viper](https://github.com/spf13/viper)** - Configuration management
- **[zerolog](https://github.com/rs/zerolog)** - Structured logging
- **[Task](https://taskfile.dev/)** - Build automation
- **[golangci-lint](https://golangci-lint.run/)** - Linter aggregator
- **[MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)** - Official MCP SDK

---

**Last Updated**: 2025-10-06
**Project Version**: Based on `ai` branch (pre-v1.2.0)
**New Features**: JSON output, MCP server, AI agent integration

For questions or clarifications, refer to the source code and tests as the ultimate source of truth.
