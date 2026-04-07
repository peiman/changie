# PLAN: `changie diff` Command

## Summary

Add a `changie diff v1.0.0 v1.1.0` command that extracts and displays all changelog
entries between two specified versions. This is a **root-level** command (like `bump`),
not a `changelog` subcommand, matching the usage pattern in the task description.

---

## Requirements Analysis

### Functional Requirements

1. **Usage:** `changie diff <from-version> <to-version>`
2. **Behavior:** Extract and display all changelog sections/entries for versions
   greater than `from-version` up to and including `to-version`
3. **Version handling:** Support both `v`-prefixed (`v1.0.0`) and bare (`1.0.0`) versions
4. **Output:** Print the extracted changelog content to stdout
5. **File flag:** Respect `--file` flag for custom changelog paths (like `bump` does)
6. **Error cases:**
   - Changelog file not found → clear error message
   - Version not found in changelog → specific error naming the missing version
   - `from-version` is newer than `to-version` → error with guidance
   - Invalid semver → error with format guidance
   - No content between versions → informational message (not an error)

### Non-Functional Requirements

- Command file ≤30 lines (logic in `internal/changelog/`)
- 85%+ test coverage
- TDD: tests first, then implementation
- Uses `testify/assert` and table-driven tests
- Follows existing patterns (see `changelog_validate.go` + `bump.go`)

---

## Architecture Decisions

### AD-1: Root-level command, not `changelog` subcommand

**Decision:** `changie diff` (root command), not `changie changelog diff`

**Rationale:** The task specifies `changie diff v1.0.0 v1.1.0`. This also follows
the pattern set by `bump` — which is a root command despite operating on changelogs.
`diff` is a high-level workflow command, not a changelog management operation.

**Trade-off:** Less organizational grouping under `changelog`, but better UX
(shorter command, matches user's mental model).

### AD-2: Logic placement — `internal/changelog/` (infrastructure layer)

**Decision:** Add `DiffVersions()` to `internal/changelog/diff.go`

**Rationale:** The `internal/changelog/` package is classified as **infrastructure**
in `.go-arch-lint.yml:84`. It already contains all changelog parsing/manipulation
functions. The diff logic is pure changelog content parsing — no orchestration or
external coordination needed. Unlike `bump` (which orchestrates git + changelog +
semver), `diff` only parses changelog content, so it belongs in the infrastructure
layer alongside `validate.go` and `changelog.go`.

**Alternative considered:** `internal/diff/` as a business logic package. Rejected
because it would create an unnecessary package for a single function that naturally
belongs with the other changelog operations.

### AD-3: Semver comparison for version range

**Decision:** Use `internal/semver.Compare()` to validate version ordering and
`internal/semver.ParseVersion()` to normalize version strings.

**Rationale:** The existing `semver` package already handles v-prefix stripping
and comparison. Reusing it avoids duplicating version parsing logic.

### AD-4: Command wiring pattern — direct `cobra.Command` (not `MustNewCommand`)

**Decision:** Use direct `cobra.Command` definition like `bump.go`, not
`MustNewCommand` like `ping.go`.

**Rationale:** `MustNewCommand` requires config registry metadata
(`config.CommandMetadata`) and auto-registers flags from the config registry. The
`diff` command has only the shared `--file` flag (from viper binding) and takes
positional args — no command-specific config options needed. This matches `bump.go`'s
pattern exactly.

---

## Module/File Structure

```
cmd/
  diff.go              # Command wiring (≤30 lines) — NEW
  diff_test.go         # Command-level tests — NEW

internal/changelog/
  diff.go              # DiffVersions() business logic — NEW
  diff_test.go         # Unit tests for diff logic — NEW
```

No changes needed to:
- `.go-arch-lint.yml` — `internal/changelog/` is already in the infrastructure layer
- Config registry — no new config options
- `keys_generated.go` — no new keys

---

## Step-by-Step Implementation Tasks

> **TDD required:** For each step, write the failing test first, then implement.
> Commit test + implementation together atomically.

### Step 1: Create `internal/changelog/diff.go` + `internal/changelog/diff_test.go`

**Tests first** in `internal/changelog/diff_test.go`:

```go
package changelog

// Test cases for DiffVersions():
// 1. Happy path: two versions exist, content between them returned
// 2. Multiple versions in range (v1.0.0 → v1.2.0 with v1.1.0 in between)
// 3. Adjacent versions (nothing between them except the to-version itself)
// 4. from-version not found → error
// 5. to-version not found → error  
// 6. from-version > to-version → error ("from must be older than to")
// 7. Same version for both → error
// 8. Versions with v-prefix work
// 9. Versions without v-prefix work
// 10. Mixed prefix (v1.0.0 and 1.1.0) works
// 11. Empty changelog → error
// 12. Changelog with only Unreleased → error (versions not found)
```

**Implementation** in `internal/changelog/diff.go`:

```go
package changelog

// DiffVersions extracts changelog content between two versions.
//
// It returns all sections and entries for versions strictly greater than
// fromVersion up to and including toVersion. Both versions must exist
// in the changelog content.
//
// Version strings may include a 'v' prefix (e.g., "v1.0.0" or "1.0.0").
//
// Parameters:
//   - content: The full changelog file content as a string
//   - fromVersion: The older version (exclusive — its content is NOT included)
//   - toVersion: The newer version (inclusive — its content IS included)
//
// Returns:
//   - string: The extracted changelog content between the two versions
//   - error: If versions are not found, invalid, or fromVersion >= toVersion
func DiffVersions(content, fromVersion, toVersion string) (string, error)
```

**Algorithm:**
1. Parse `fromVersion` and `toVersion` with `semver.ParseVersion()` → validate both
2. Compare versions with `semver.Compare()` → ensure `fromVersion < toVersion`
3. Scan changelog lines for `## [X.Y.Z]` headers using `reVersionHeader` regex (already compiled in `validate.go`)
4. Find line indices for both version headers (normalize v-prefix for matching)
5. Extract all lines from `toVersion` header (inclusive) to `fromVersion` header (exclusive)
6. Return extracted content as a trimmed string

**Key implementation details:**
- Reuse the existing `reVersionHeader` regex from `validate.go` (it's package-level, already available)
- Strip v-prefix before comparing with changelog headers (changelog may use `[1.0.0]` while user passes `v1.0.0`)
- Handle both `## [1.0.0] - 2024-01-01` and `## [v1.0.0] - 2024-01-01` formats

### Step 2: Create `cmd/diff.go` + `cmd/diff_test.go`

**Tests first** in `cmd/diff_test.go`:

```go
package cmd

// Test cases:
// 1. Happy path: valid versions, output contains expected content
// 2. Missing changelog file → error containing "failed to read changelog"
// 3. Version not found → error containing version string
// 4. Inverted versions → error about ordering
// 5. Wrong number of args (0, 1, 3) → cobra arg validation error
// 6. Custom --file flag works
// 7. Command is registered on RootCmd
```

Follow the test pattern from `cmd/changelog_validate_test.go`:
- Create temp dir with temp changelog file
- Use `viper.Reset()` + `viper.Set("app.changelog.file", path)`
- Create command with `RunE: runDiff`
- Capture output with `bytes.Buffer` via `cmd.SetOut()`

**Implementation** in `cmd/diff.go`:

```go
// cmd/diff.go
package cmd

import (
    "fmt"
    "os"

    "github.com/peiman/changie/internal/changelog"
    "github.com/rs/zerolog/log"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var diffCmd = &cobra.Command{
    Use:   "diff FROM TO",
    Short: "Show changelog entries between two versions",
    Long: `Compare two versions in the changelog and show what changed between them.

Extracts and displays all changelog entries for versions after FROM up to
and including TO. Both versions must exist in the changelog file.

Examples:
  changie diff 1.0.0 1.1.0
  changie diff v1.0.0 v2.0.0
  changie diff 0.9.0 0.9.1 --file HISTORY.md`,
    Args: cobra.ExactArgs(2),
    RunE: runDiff,
}

func init() {
    diffCmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
    if err := viper.BindPFlag("app.changelog.file", diffCmd.Flags().Lookup("file")); err != nil {
        log.Fatal().Err(err).Msg("Failed to bind 'file' flag")
    }
    RootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
    file := getConfigValueWithFlags[string](cmd, "file", "app.changelog.file")
    data, err := os.ReadFile(file)
    if err != nil {
        return fmt.Errorf("failed to read changelog: %w", err)
    }
    result, err := changelog.DiffVersions(string(data), args[0], args[1])
    if err != nil {
        return err
    }
    _, _ = fmt.Fprint(cmd.OutOrStdout(), result)
    return nil
}
```

**Line count of `runDiff`:** ~10 lines ✅ (well under 30-line limit)

### Step 3: Run `task check`

Verify everything passes before committing. Fix any lint, format, or coverage issues.

### Step 4: Commit atomically

Single commit with all 4 files:
- `internal/changelog/diff.go`
- `internal/changelog/diff_test.go`
- `cmd/diff.go`
- `cmd/diff_test.go`

---

## Testing Strategy

### Unit Tests (`internal/changelog/diff_test.go`)

| # | Test Case | Input | Expected |
|---|-----------|-------|----------|
| 1 | Happy path — two adjacent versions | content with v1.0.0 and v1.1.0 | v1.1.0 section content |
| 2 | Multiple versions in range | v1.0.0, v1.0.1, v1.1.0; diff(1.0.0, 1.1.0) | Both v1.0.1 and v1.1.0 sections |
| 3 | Versions with v-prefix | `diff("v1.0.0", "v1.1.0")` | Works, same as bare |
| 4 | Versions without v-prefix | `diff("1.0.0", "1.1.0")` | Works |
| 5 | Mixed prefix user vs changelog | User passes `v1.0.0`, changelog has `[1.0.0]` | Works (normalized) |
| 6 | from-version not found | `diff("0.5.0", "1.0.0")` where 0.5.0 doesn't exist | Error: "version 0.5.0 not found" |
| 7 | to-version not found | `diff("1.0.0", "3.0.0")` where 3.0.0 doesn't exist | Error: "version 3.0.0 not found" |
| 8 | from > to (inverted) | `diff("2.0.0", "1.0.0")` | Error: "from-version must be older" |
| 9 | Same version | `diff("1.0.0", "1.0.0")` | Error: "versions must be different" |
| 10 | Invalid semver (from) | `diff("not-a-version", "1.0.0")` | Error: "invalid version" |
| 11 | Invalid semver (to) | `diff("1.0.0", "abc")` | Error: "invalid version" |
| 12 | Empty content | `diff("1.0.0", "2.0.0")` on `""` | Error |
| 13 | Content with multiple sections per version | v1.1.0 has Added + Fixed + Changed | All sections included in output |
| 14 | Unreleased not included | diff(1.0.0, 1.1.0) with Unreleased present | Unreleased content excluded |

### Command Tests (`cmd/diff_test.go`)

| # | Test Case | Setup | Expected |
|---|-----------|-------|----------|
| 1 | Happy path | Temp file with valid changelog | No error, output contains version content |
| 2 | File not found | Non-existent path | Error: "failed to read changelog" |
| 3 | Version not found | Valid file, bad version | Error containing the version string |
| 4 | Wrong arg count (0) | No args | Cobra error |
| 5 | Wrong arg count (1) | One arg | Cobra error |
| 6 | Custom file flag | `--file custom.md` | Uses custom file |
| 7 | Command registered | Check RootCmd.Commands() | `diff` found |

---

## Edge Cases

1. **Changelog uses v-prefix but user omits it (or vice versa):** Normalize both
   the user input and changelog header by stripping `v` prefix before comparison.

2. **Version exists in link references but not as header:** Only match `## [X.Y.Z]`
   headers, not `[X.Y.Z]: URL` link references. The existing `reVersionHeader`
   regex already anchors to `## [`.

3. **Trailing whitespace in version headers:** Use `strings.TrimSpace()` on lines
   before matching, consistent with `validate.go` patterns.

4. **Windows line endings (`\r\n`):** Normalize with `strings.ReplaceAll(content, "\r\n", "\n")`
   before parsing, matching the pattern in `ValidateChangelog()` at `validate.go:51`.

5. **Changelog with no versions (only Unreleased):** Both versions will fail to
   be found → return clear error for each.

6. **from-version is the very first release:** All content from `to-version` down
   to (but not including) `from-version` header is returned.

7. **to-version is the latest release:** Content starts right after the `[Unreleased]`
   section boundary.

---

## References

| File | Line(s) | Relevance |
|------|---------|-----------|
| `cmd/bump.go:20-42` | — | Pattern for root-level command with `--file` flag |
| `cmd/bump.go:135-145` | — | `runVersionBump` pattern: read config, call internal |
| `cmd/changelog_validate.go:33-51` | — | Pattern for reading changelog + calling internal |
| `cmd/changelog_validate_test.go:64-79` | — | `setupValidateCmd` test helper pattern |
| `internal/changelog/validate.go:40-45` | — | Compiled regex `reVersionHeader` — reuse for diff |
| `internal/changelog/validate.go:49-51` | — | `\r\n` normalization pattern |
| `internal/changelog/changelog.go:262-275` | — | `GetLatestChangelogVersion` — version header parsing |
| `internal/semver/semver.go:34-54` | — | `ParseVersion()` — v-prefix handling |
| `internal/semver/semver.go:120-137` | — | `Compare()` — version ordering |
| `.go-arch-lint.yml:84` | — | `internal/changelog/` is infrastructure layer |
| `.go-arch-lint.yml:104` | — | `internal/version/` is business layer |
| `AGENTS.md:296-307` | — | New Command Checklist |
| `CLAUDE.md` | — | TDD rule, `task check` rule, ≤30 lines rule |
