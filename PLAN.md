# PLAN: Port Changie to ckeletin-go Scaffold

## Overview

Port the changie changelog management CLI from its current standalone structure into the ckeletin-go scaffold framework. The source changie code lives at `/Users/peiman/dev/cli/changie/changie/` and the ckeletin-go scaffold at `/Users/peiman/dev/cli/ckeletin-go/`.

**Target**: `/Users/peiman/dev/workhorse/repos/changie/` (this repo, currently empty).

**Goal**: A fully working changie CLI built on ckeletin-go's framework, passing `task check`, with all existing functionality preserved.

---

## Phase 1: Scaffold Setup

### Step 1.1: Copy ckeletin-go scaffold
Copy the entire ckeletin-go repo contents into this repo:
```
cp -r /Users/peiman/dev/cli/ckeletin-go/* /Users/peiman/dev/workhorse/repos/changie/
cp -r /Users/peiman/dev/cli/ckeletin-go/.ckeletin /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.gitignore /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.golangci.yml /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.goreleaser.yml /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.go-version /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.go-arch-lint.yml /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.semgrep.yml /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.lefthook.yml /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.gitleaks.toml /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.lichen.yaml /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/.gitattributes /Users/peiman/dev/workhorse/repos/changie/
cp /Users/peiman/dev/cli/ckeletin-go/codecov.yml /Users/peiman/dev/workhorse/repos/changie/
cp -r /Users/peiman/dev/cli/ckeletin-go/.github /Users/peiman/dev/workhorse/repos/changie/
cp -r /Users/peiman/dev/cli/ckeletin-go/.vscode /Users/peiman/dev/workhorse/repos/changie/
cp -r /Users/peiman/dev/cli/ckeletin-go/.devcontainer /Users/peiman/dev/workhorse/repos/changie/
cp -r /Users/peiman/dev/cli/ckeletin-go/test /Users/peiman/dev/workhorse/repos/changie/
cp -r /Users/peiman/dev/cli/ckeletin-go/testdata /Users/peiman/dev/workhorse/repos/changie/
```

Also copy:
- `main.go`, `main_test.go`
- `cmd/root.go`, `cmd/helpers.go`, `cmd/flags.go`
- `cmd/config.go`, `cmd/check.go`, `cmd/docs.go` (framework commands)
- `cmd/dev*.go` (dev commands, build-tagged)
- `internal/ui/`, `internal/xdg/`, `internal/progress/`, `internal/dev/`, `internal/check/`, `internal/docs/` (framework internals)
- `internal/config/commands/` (framework config definitions)
- `pkg/checkmate/` (public package)
- `Taskfile.yml`

**Do NOT copy**: `cmd/ping.go`, `internal/ping/`, `internal/config/commands/ping_config.go` (example command we don't need).

### Step 1.2: Update module identity
In `go.mod`:
- Change module to `github.com/peiman/changie`
- Keep Go version from ckeletin-go (1.26.1)

In `Taskfile.yml`:
- Change `BINARY_NAME` to `changie`
- Change `MODULE_PATH` to `github.com/peiman/changie`

### Step 1.3: Update all import paths
Replace all occurrences of `github.com/peiman/ckeletin-go` with `github.com/peiman/changie` across all `.go` files.

### Step 1.4: Remove example command references
- Delete `cmd/ping.go` if copied
- Delete `internal/ping/` if copied
- Delete `internal/config/commands/ping_config.go` if copied
- Remove ping-related entries from any config files

### Step 1.5: Verify scaffold compiles
```
cd /Users/peiman/dev/workhorse/repos/changie
go mod tidy
go build ./...
```

**Commit**: "feat: initialize changie from ckeletin-go scaffold"

---

## Phase 2: Port Dependencies

### Step 2.1: Add changie-specific dependencies
```
go get github.com/blang/semver/v4
```
Then run `task check:license:source` to verify license compliance.

**Note**: The existing ckeletin-go deps (cobra, viper, zerolog, bubbletea, lipgloss, testify) are already present.

**Commit**: "feat: add semver dependency"

---

## Phase 3: Port Business Logic (internal packages)

Port packages in dependency order (leaf packages first, then packages that depend on them).

### Step 3.1: Port `internal/semver/`
**Source**: `/Users/peiman/dev/cli/changie/changie/internal/semver/semver.go`
**Target**: `/Users/peiman/dev/workhorse/repos/changie/internal/semver/semver.go`

Create `internal/semver/` with:
- `semver.go` - BumpType type, ParseVersion, FormatVersion, BumpVersion, BumpMajor, BumpMinor, BumpPatch, Compare
- `semver_test.go` - Port all tests

This package has no internal dependencies (only `blang/semver/v4`).

**Key types to port**:
```go
type BumpType string
const (Major, Minor, Patch BumpType)
```

**Key functions**:
- `ParseVersion(version string) (semver.Version, bool, error)` - parse with v-prefix detection
- `FormatVersion(ver semver.Version, includePrefix bool) string`
- `BumpVersion(version, bumpType, useVPrefix) (string, error)`
- `BumpMajor/BumpMinor/BumpPatch` convenience functions
- `Compare(v1, v2 string) (int, error)`

### Step 3.2: Port `internal/git/`
**Source**: `/Users/peiman/dev/cli/changie/changie/internal/git/git.go`
**Target**: `/Users/peiman/dev/workhorse/repos/changie/internal/git/git.go`

Create `internal/git/` with:
- `git.go` - All git operations
- `git_test.go` - Port all tests

**Key functions**:
- `IsInstalled() bool`
- `GetVersion() (string, error)`
- `GetCurrentBranch() (string, error)`
- `HasUncommittedChanges() (bool, error)`
- `CommitChangelog(file, version string) error`
- `TagVersion(version string) error`
- `PushChanges() error`
- `GetRemoteURL() (string, error)`
- `ParseRepositoryURL(remoteURL string) (*RepositoryInfo, error)`

**Adapt logging**: Use ckeletin-go's `zerolog` (already global via `log` package). Use `log.Debug()` for operation-level logging.

### Step 3.3: Port `internal/changelog/`
**Source**: `/Users/peiman/dev/cli/changie/changie/internal/changelog/changelog.go`
**Target**: `/Users/peiman/dev/workhorse/repos/changie/internal/changelog/changelog.go`

Create `internal/changelog/` with:
- `changelog.go` - Core changelog operations
- `changelog_test.go` - Port all tests

**Key exports**:
- `ValidSections` map
- `InitProject(filePath string) error`
- `AddChangelogSection(filePath, section, content string) (bool, error)`
- `UpdateChangelog(filePath, version, repositoryProvider string) error`
- `GetLatestChangelogVersion(content string) (string, error)`

### Step 3.4: Port `internal/version/`
**Source**: `/Users/peiman/dev/cli/changie/changie/internal/version/bump.go`
**Target**: `/Users/peiman/dev/workhorse/repos/changie/internal/version/bump.go`

Create `internal/version/` with:
- `bump.go` - Version bump orchestrator (the main workflow)
- `bump_test.go` - Write tests (source has 0% coverage, we need to write them)

**BumpConfig struct**:
```go
type BumpConfig struct {
    BumpType           string
    AllowAnyBranch     bool
    AutoPush           bool
    ChangelogFile      string
    RepositoryProvider string
    UseVPrefix         bool
}
```

**Adapt to Executor pattern** (ckeletin-go convention):
```go
type Executor struct {
    cfg    BumpConfig
    writer io.Writer
}

func NewExecutor(cfg BumpConfig, w io.Writer) *Executor
func (e *Executor) Execute() error
```

The `Execute()` method implements the bump workflow:
1. Verify git installed
2. Check branch (main/master)
3. Check uncommitted changes
4. Get current version
5. Bump version
6. Update changelog
7. Commit
8. Tag
9. Optionally push
10. Output result

**Adapt output**: Use ckeletin-go's `.ckeletin/pkg/output/` JSON envelope pattern instead of changie's custom output package.

### Step 3.5: Port output types
**Source**: `/Users/peiman/dev/cli/changie/changie/internal/output/output.go`

We do NOT port changie's output package. Instead, define result structs in the relevant internal packages and use ckeletin-go's `output.JSONEnvelope` / `ui.RenderSuccess()` pattern.

Define these structs in their respective packages:
- `internal/version/bump.go`: `BumpResult` struct
- `internal/changelog/changelog.go`: `InitResult`, `AddSectionResult` structs
- `internal/docs/`: result structs as needed

### Step 3.6: Port `internal/docs/` (changie docs, not ckeletin-go docs)
**Source**: `/Users/peiman/dev/cli/changie/changie/internal/docs/`

**Note**: ckeletin-go already has its own `internal/docs/` for framework documentation. Changie's docs package generates config documentation. We need to either:
- Merge changie's doc generation into the existing ckeletin-go docs system, OR
- Keep changie's docs as a separate subpackage

**Recommendation**: Since ckeletin-go already has `cmd/docs.go` that generates config documentation, check if its existing functionality covers changie's needs. If so, no porting needed. If changie has special doc generation (changelog-format docs), create it as a separate concern.

### Step 3.7: Port UI utilities
**Source**: `/Users/peiman/dev/cli/changie/changie/internal/ui/prompt.go`

Port changie's `AskYesNo()` function. ckeletin-go already has `internal/ui/` with rendering utilities. Add the prompt functionality:

**Target**: Add to `/Users/peiman/dev/workhorse/repos/changie/internal/ui/prompt.go`
- `AskYesNo(prompt string, defaultYes bool, output io.Writer) (bool, error)`

**Commit after each package**: Atomic commits with test + implementation together (TDD).

---

## Phase 4: Port Config Definitions

### Step 4.1: Create changie command config files

Create `internal/config/commands/` files for each changie command following ckeletin-go pattern.

#### `internal/config/commands/init_config.go`:
```go
var InitMetadata = config.CommandMetadata{
    Use:          "init",
    Short:        "Initialize a project with a CHANGELOG.md",
    ConfigPrefix: "app.init",
    FlagOverrides: map[string]string{
        "app.changelog.file":       "file",
        "app.version.use_v_prefix": "v-prefix",
    },
}

func InitOptions() []config.ConfigOption {
    return []config.ConfigOption{
        {
            Key:          "app.changelog.file",
            DefaultValue: "CHANGELOG.md",
            Description:  "Changelog file name",
            Type:         "string",
            ShortFlag:    "f",
        },
        {
            Key:          "app.version.use_v_prefix",
            DefaultValue: true,
            Description:  "Use 'v' prefix for version tags",
            Type:         "bool",
        },
    }
}

func init() {
    config.RegisterOptionsProvider(InitOptions)
}
```

#### `internal/config/commands/bump_config.go`:
```go
var BumpMetadata = config.CommandMetadata{
    Use:          "bump",
    Short:        "Bump the project version",
    ConfigPrefix: "app.bump",
}

// BumpMajorMetadata, BumpMinorMetadata, BumpPatchMetadata for subcommands

func BumpOptions() []config.ConfigOption {
    return []config.ConfigOption{
        {Key: "app.bump.auto_push", DefaultValue: false, Type: "bool", Description: "Auto-push after bump"},
        {Key: "app.bump.allow_any_branch", DefaultValue: false, Type: "bool", Description: "Allow bump on any branch"},
        {Key: "app.changelog.file", DefaultValue: "CHANGELOG.md", Type: "string", ShortFlag: "f"},
        {Key: "app.changelog.repository_provider", DefaultValue: "github", Type: "string"},
        {Key: "app.version.use_v_prefix", DefaultValue: true, Type: "bool"},
    }
}
```

#### `internal/config/commands/changelog_config.go`:
```go
var ChangelogMetadata = config.CommandMetadata{
    Use:          "changelog",
    Short:        "Manage changelog entries",
    ConfigPrefix: "app.changelog",
    FlagOverrides: map[string]string{
        "app.changelog.file": "file",
    },
}
```

**Note**: Some config keys are shared between commands (e.g., `app.changelog.file`). Only register each key once. Use the first command that defines it, or create a shared options provider.

### Step 4.2: Generate key constants
```
task generate:config:key-constants
```

This creates/updates `.ckeletin/pkg/config/keys_generated.go` with constants like:
- `KeyAppChangelogFile`
- `KeyAppVersionUseVPrefix`
- `KeyAppBumpAutoPush`
- `KeyAppBumpAllowAnyBranch`
- `KeyAppChangelogRepositoryProvider`

### Step 4.3: Update config.schema.json
Generate the config schema to include changie's configuration keys.

**Commit**: "feat: add changie configuration definitions"

---

## Phase 5: Port Commands

### Step 5.1: Port `cmd/init.go`
**Target**: `/Users/peiman/dev/workhorse/repos/changie/cmd/init.go`

Follow ultra-thin command pattern:
```go
var initCmd = MustNewCommand(commands.InitMetadata, runInit)

func init() {
    MustAddToRoot(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
    // Create init executor config from flags/viper
    // Call internal/init executor
    // Return error
}
```

**Note**: The init command has interactive prompts (AskYesNo). Pass `cmd.InOrStdin()` and `cmd.OutOrStdout()` for testability.

Create `internal/init/` package with executor:
```go
type Config struct {
    ChangelogFile string
    UseVPrefix    bool
}

type Executor struct {
    cfg    Config
    writer io.Writer
    reader io.Reader // for prompts
}
```

### Step 5.2: Port `cmd/bump.go`
**Target**: `/Users/peiman/dev/workhorse/repos/changie/cmd/bump.go`

This has a parent command (`bump`) with three subcommands (`major`, `minor`, `patch`).

```go
var bumpCmd = &cobra.Command{
    Use:   "bump",
    Short: "Bump the project version",
}

var bumpMajorCmd = MustNewCommand(commands.BumpMajorMetadata, runBumpMajor)
var bumpMinorCmd = MustNewCommand(commands.BumpMinorMetadata, runBumpMinor)
var bumpPatchCmd = MustNewCommand(commands.BumpPatchMetadata, runBumpPatch)

func init() {
    MustAddToRoot(bumpCmd)
    bumpCmd.AddCommand(bumpMajorCmd)
    bumpCmd.AddCommand(bumpMinorCmd)
    bumpCmd.AddCommand(bumpPatchCmd)
}

func runBumpMajor(cmd *cobra.Command, args []string) error {
    return runBump(cmd, "major")
}

func runBump(cmd *cobra.Command, bumpType string) error {
    cfg := version.BumpConfig{
        BumpType:           bumpType,
        AllowAnyBranch:     getConfigValueWithFlags[bool](cmd, "allow-any-branch", config.KeyAppBumpAllowAnyBranch),
        AutoPush:           getConfigValueWithFlags[bool](cmd, "auto-push", config.KeyAppBumpAutoPush),
        ChangelogFile:      getConfigValueWithFlags[string](cmd, "file", config.KeyAppChangelogFile),
        RepositoryProvider: getConfigValueWithFlags[string](cmd, "rrp", config.KeyAppChangelogRepositoryProvider),
        UseVPrefix:         getConfigValueWithFlags[bool](cmd, "v-prefix", config.KeyAppVersionUseVPrefix),
    }
    return version.NewExecutor(cfg, cmd.OutOrStdout()).Execute()
}
```

### Step 5.3: Port `cmd/changelog.go` and `cmd/changelog_add.go`
**Target**: `/Users/peiman/dev/workhorse/repos/changie/cmd/changelog.go`

Parent command with dynamic subcommands for each section type.

```go
var changelogCmd = &cobra.Command{
    Use:   "changelog",
    Short: "Manage changelog entries",
}

func init() {
    MustAddToRoot(changelogCmd)
    // Register subcommands for each section
    for section := range changelog.ValidSections {
        cmd := createSectionCmd(section)
        changelogCmd.AddCommand(cmd)
    }
}

func createSectionCmd(section string) *cobra.Command {
    return &cobra.Command{
        Use:   strings.ToLower(section) + " CONTENT",
        Short: fmt.Sprintf("Add a %s entry to the changelog", section),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            file := getConfigValueWithFlags[string](cmd, "file", config.KeyAppChangelogFile)
            _, err := changelog.AddChangelogSection(file, section, args[0])
            return err
        },
    }
}
```

### Step 5.4: Port `cmd/completion.go`
**Target**: `/Users/peiman/dev/workhorse/repos/changie/cmd/completion.go`

Simple port using Cobra's built-in completion generation.

### Step 5.5: MCP Server decision
The MCP server is a separate binary (`cmd/mcp-server/main.go`). This can be ported as-is but should use the changie module path. Consider whether to include it in this initial port or defer.

**Recommendation**: Defer MCP server to a follow-up PR. Focus on the core CLI first.

**Commit after each command**: Atomic commits (TDD: test + implementation).

---

## Phase 6: Update Framework Integration

### Step 6.1: Update AGENTS.md
Write a new AGENTS.md for changie that includes:
- Project description (changelog management CLI)
- Commands reference (init, bump, changelog, completion, docs)
- Configuration keys reference
- Architecture overview
- All ckeletin-go conventions

### Step 6.2: Update config.schema.json
Regenerate with changie's config keys.

### Step 6.3: Update .go-arch-lint.yml
Add changie-specific packages to the architecture rules:
- `internal/semver/` - no imports from other internal packages
- `internal/git/` - no imports from other internal packages
- `internal/changelog/` - may import `internal/semver/`
- `internal/version/` - may import `internal/changelog/`, `internal/semver/`, `internal/git/`
- `internal/init/` - may import `internal/changelog/`, `internal/git/`, `internal/ui/`

### Step 6.4: Update README.md
Write changie-specific README with:
- What changie does
- Installation
- Quick start
- Commands reference
- Configuration
- Contributing

**Commit**: "docs: add changie documentation"

---

## Phase 7: Quality Gate

### Step 7.1: Run `task format`
Fix all formatting issues.

### Step 7.2: Run `task lint`
Fix all linter warnings.

### Step 7.3: Run `task test`
Ensure all tests pass with 85%+ coverage.

### Step 7.4: Run `task check`
Full quality gate - must pass completely.

### Step 7.5: Build and smoke test
```
task build
./changie --help
./changie init
./changie changelog added "Test entry"
./changie bump patch
```

**Commit**: "fix: resolve quality gate issues" (if needed)

---

## Phase 8: Create PR

Create a PR from this repo back to `https://github.com/peiman/changie`.

---

## File Mapping (Source → Target)

| Source (changie) | Target (ckeletin-go scaffold) | Notes |
|---|---|---|
| `internal/semver/semver.go` | `internal/semver/semver.go` | Direct port |
| `internal/git/git.go` | `internal/git/git.go` | Direct port, adapt logging |
| `internal/changelog/changelog.go` | `internal/changelog/changelog.go` | Direct port |
| `internal/version/bump.go` | `internal/version/bump.go` | Convert to Executor pattern |
| `internal/output/output.go` | *NOT PORTED* | Use ckeletin-go's `output` pkg |
| `internal/logger/` | *NOT PORTED* | Use ckeletin-go's `logger` pkg |
| `internal/config/registry.go` | `internal/config/commands/*.go` | Convert to ckeletin-go pattern |
| `internal/ui/prompt.go` | `internal/ui/prompt.go` | Add to existing UI pkg |
| `internal/docs/` | Evaluate if needed | ckeletin-go has docs already |
| `internal/mcp/` | Deferred | Separate PR |
| `cmd/root.go` | Use ckeletin-go's | Already has root |
| `cmd/init.go` | `cmd/init.go` | Rewrite for scaffold |
| `cmd/bump.go` | `cmd/bump.go` | Rewrite for scaffold |
| `cmd/changelog.go` | `cmd/changelog.go` | Rewrite for scaffold |
| `cmd/changelog_add.go` | Merged into `cmd/changelog.go` | Dynamic subcommands |
| `cmd/completion.go` | `cmd/completion.go` | Simple port |
| `cmd/docs.go` | Use ckeletin-go's | Already has docs |

## Key Architectural Decisions

1. **Use ckeletin-go's logger, output, and config systems** - Don't port changie's custom versions
2. **Convert version/bump to Executor pattern** - Matches ckeletin-go convention
3. **Keep internal packages pure** - `semver/`, `git/`, `changelog/` have no framework dependencies
4. **Defer MCP server** - Separate concern, separate PR
5. **Share config keys between commands** - Register each key only once via the first provider
6. **Dynamic changelog subcommands** - Generate from ValidSections, not hardcode each one

## Risk Areas

1. **Shared config keys**: Multiple commands use `app.changelog.file`. Need careful registration to avoid duplicates.
2. **Interactive prompts**: The `init` command uses `AskYesNo`. Ensure this works with ckeletin-go's command execution flow.
3. **Architecture validation**: `.go-arch-lint.yml` may need updates for changie's package dependencies.
4. **Coverage for version/**: Source has 0% coverage. Must write comprehensive tests.
5. **ckeletin-go framework commands**: `check`, `config`, `dev`, `docs` commands from ckeletin-go should be kept as-is (they're framework features).
