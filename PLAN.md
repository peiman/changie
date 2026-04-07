# PLAN: `changelog validate` Command

## Overview

Add a `changelog validate` subcommand that checks CHANGELOG.md for common problems and outputs a pass/fail report. The command follows the established ultra-thin command pattern: wiring in `cmd/`, business logic in `internal/changelog/`, configuration in `internal/config/commands/`.

## Requirements Analysis

### Functional Requirements

1. **Five validation rules**, each producing pass/fail:
   - **Missing version headers**: Version lines must match `## [X.Y.Z] - YYYY-MM-DD` (the Keep a Changelog format)
   - **Duplicate entries**: No identical bullet point text within the same section of the same version
   - **Broken links**: Every `## [X.Y.Z]` header must have a matching `[X.Y.Z]: URL` reference link at the bottom, and vice versa (orphan links)
   - **Entries without dates**: Version headers (non-Unreleased) must include a date in `YYYY-MM-DD` format
   - **Versions not in semver order**: Versions must appear in descending semver order top-to-bottom

2. **Report output**: Structured report showing each check name + pass/fail + details of failures
3. **JSON mode**: Must work with `--output json` (inherited from root command)
4. **Exit code**: Non-zero when any check fails
5. **File flag**: Reuse existing `--file` flag from `changelogCmd` (defaults to `CHANGELOG.md`)

### Non-Functional Requirements

- Unit tests for every validation rule (TDD — tests first)
- 85%+ test coverage
- Command file ≤30 lines
- Uses `testify/assert` and table-driven tests
- Uses existing `internal/semver` package for version comparison
- Uses `ui.RenderSuccess()` for output (gets JSON for free)

## Architecture Decisions

### AD-1: Subcommand of `changelog`, not top-level

**Decision**: `changie changelog validate` (subcommand under existing `changelogCmd`)

**Rationale**: The `changelog` command group already exists (`cmd/changelog.go:11-17`) with subcommands for each section type (`cmd/changelog_add.go`). Validation is a changelog operation, so it belongs here. This also gives us the `--file` flag for free since it's a persistent flag on `changelogCmd` (`cmd/changelog.go:21`).

**Trade-off**: Slightly longer command (`changelog validate` vs `validate`), but consistent grouping.

### AD-2: Validation logic in existing `internal/changelog/` package

**Decision**: Add `validate.go` to the existing `internal/changelog/` package rather than creating a new `internal/validate/` package.

**Rationale**: The validation logic operates on changelog content and reuses concepts already in this package (version header parsing at `changelog.go:263-275`, section detection at `changelog.go:129-136`). Adding to the existing package avoids circular dependencies and enables reuse of existing helper functions.

**Trade-off**: Makes the `internal/changelog/` package larger, but the validation code is cohesive with existing changelog operations.

### AD-3: Pure function validation with `io.Reader` input

**Decision**: Core validation function takes `string` content (not file path) to enable pure unit testing without filesystem.

**Rationale**: Following the dependency injection pattern (ADR-003). The cmd layer handles file reading; the business logic validates content. This makes tests fast and deterministic.

### AD-4: Use `ui.RenderSuccess()` for output, not checkmate

**Decision**: Use `ui.RenderSuccess()` with a structured `ValidationReport` type that implements `output.JSONResponder`.

**Rationale**: The `check` command uses checkmate because it runs long-running shell commands with progress tracking. Our validation is instantaneous (parse-only, no I/O). Using `ui.RenderSuccess()` gives us JSON mode for free and keeps the implementation simple. Text-mode output will format the report with pass/fail indicators.

**Trade-off**: Less visually rich than checkmate output, but appropriate for the use case and significantly simpler to implement.

## Module/File Structure

```
cmd/
  changelog_validate.go          # Thin command wiring (≤30 lines) — NEW
  changelog_validate_test.go     # Command integration tests — NEW

internal/
  changelog/
    validate.go                  # Core validation logic (5 rules) — NEW
    validate_test.go             # Unit tests for all validation rules — NEW

internal/config/commands/
    changelog_validate_config.go # Command metadata — NEW (minimal, no custom options)
```

### No new config keys needed

The only flag is `--file` which is already a persistent flag on `changelogCmd` (`cmd/changelog.go:21`). No new entries in the config registry, no regeneration of `keys_generated.go`.

## Detailed Type Design

### `internal/changelog/validate.go`

```go
// ValidationResult represents the outcome of a single validation check.
type ValidationResult struct {
    Name    string   `json:"name"`
    Passed  bool     `json:"passed"`
    Message string   `json:"message"`
    Details []string `json:"details,omitempty"` // Individual failure descriptions
}

// ValidationReport represents the complete validation report for a changelog file.
type ValidationReport struct {
    File       string             `json:"file"`
    Passed     bool               `json:"passed"`
    TotalRules int                `json:"total_rules"`
    PassCount  int                `json:"pass_count"`
    FailCount  int                `json:"fail_count"`
    Results    []ValidationResult `json:"results"`
}

// JSONResponse implements output.JSONResponder for clean JSON output.
func (r *ValidationReport) JSONResponse() interface{} {
    return r
}

// ValidateChangelog runs all validation checks against changelog content.
// The filePath parameter is used for report metadata only (not for file I/O).
func ValidateChangelog(content string, filePath string) *ValidationReport { ... }
```

### Five Validation Functions (internal, unexported)

```go
func checkVersionHeaders(lines []string) ValidationResult       // Rule 1: malformed version headers
func checkDuplicateEntries(lines []string) ValidationResult     // Rule 2: duplicate bullet points
func checkBrokenLinks(lines []string) ValidationResult          // Rule 3: header↔link mismatches
func checkEntriesWithoutDates(lines []string) ValidationResult  // Rule 4: versions missing dates
func checkSemverOrder(lines []string) ValidationResult          // Rule 5: descending semver order
```

### `cmd/changelog_validate.go`

```go
var validateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate changelog for common problems",
    ...
    RunE: runValidateChangelog,
}

func init() {
    changelogCmd.AddCommand(validateCmd)
}

func runValidateChangelog(cmd *cobra.Command, args []string) error {
    file := getConfigValueWithFlags[string](cmd, "file", "app.changelog.file")
    data, err := os.ReadFile(file)
    // ... error handling ...
    report := changelog.ValidateChangelog(string(data), file)
    // ... render with ui.RenderSuccess or return error ...
}
```

## Step-by-Step Implementation Tasks

> **TDD is mandatory**: For each step, write failing tests FIRST, then implement to make them pass. Commit test + implementation together as one atomic unit.

### Step 1: Create validation types and stub (`internal/changelog/validate.go`)

**Files**: `internal/changelog/validate.go`, `internal/changelog/validate_test.go`

1. Write tests for `ValidateChangelog` with a valid changelog (all checks pass)
2. Write tests for `ValidationReport` struct (verify JSON serialization, `Passed` calculation)
3. Implement the types (`ValidationResult`, `ValidationReport`) and `ValidateChangelog` stub that calls all five check functions (initially returning all-pass)
4. Verify tests pass

### Step 2: Implement Rule 1 — `checkVersionHeaders`

**File**: `internal/changelog/validate.go`, `internal/changelog/validate_test.go`

Tests to write first:
- Valid changelog → pass
- Version header missing brackets `## 1.0.0 - 2024-01-01` → fail with detail
- Version with invalid semver `## [abc] - 2024-01-01` → fail with detail
- `## [Unreleased]` is NOT flagged (it's valid without version/date)
- Version with prerelease `## [1.0.0-beta.1] - 2024-01-01` → pass
- Version with `v` prefix `## [v1.0.0] - 2024-01-01` → pass
- Multiple valid versions → pass
- Mix of valid and invalid → fail listing only the invalid ones

Implementation:
- Regex: `^## \[([^\]]+)\]` to find all version-like headers
- Skip `[Unreleased]`
- For remaining, validate the version portion parses as semver (use `internal/semver.ParseVersion`)

### Step 3: Implement Rule 2 — `checkDuplicateEntries`

**File**: `internal/changelog/validate.go`, `internal/changelog/validate_test.go`

Tests to write first:
- No duplicates → pass
- Duplicate in same section → fail with detail showing the duplicate text and section
- Same text in different sections → pass (not a duplicate)
- Same text in different versions → pass (not a duplicate)
- Case-sensitive comparison (not case-insensitive)
- Whitespace-normalized comparison (leading/trailing spaces stripped)

Implementation:
- Walk through lines tracking current version + section
- For each section, collect bullet entries (`- ` prefix) in a set
- Flag any entry seen twice in the same version+section

### Step 4: Implement Rule 3 — `checkBrokenLinks`

**File**: `internal/changelog/validate.go`, `internal/changelog/validate_test.go`

Tests to write first:
- All version headers have matching links → pass
- Version header without matching link → fail with detail
- Link without matching version header (orphan link) → fail with detail
- `[Unreleased]` header with `[Unreleased]: URL` link → pass
- `[Unreleased]` header without link → fail
- Extra link not matching any header → fail (orphan)
- No links section at all → fail (if there are version headers)

Implementation:
- Collect all versions from `## [X]` headers → set A
- Collect all versions from `[X]: URL` link references → set B
- Report A - B as "missing links" and B - A as "orphan links"

### Step 5: Implement Rule 4 — `checkEntriesWithoutDates`

**File**: `internal/changelog/validate.go`, `internal/changelog/validate_test.go`

Tests to write first:
- All versions have dates → pass
- Version without date `## [1.0.0]` → fail with detail
- `## [Unreleased]` without date → pass (Unreleased doesn't need a date)
- Invalid date format `## [1.0.0] - 01-01-2024` → fail
- Valid date format `## [1.0.0] - 2024-01-01` → pass
- Date with extra text `## [1.0.0] - 2024-01-01 [YANKED]` → pass (date is present)

Implementation:
- Find all version headers (non-Unreleased)
- Check each has ` - YYYY-MM-DD` following the closing bracket
- Use regex: `## \[[^\]]+\] - (\d{4}-\d{2}-\d{2})`

### Step 6: Implement Rule 5 — `checkSemverOrder`

**File**: `internal/changelog/validate.go`, `internal/changelog/validate_test.go`

Tests to write first:
- Descending order `1.2.0, 1.1.0, 1.0.0` → pass
- Out of order `1.0.0, 1.2.0, 1.1.0` → fail with detail showing which pair is wrong
- Single version → pass
- No versions (only Unreleased) → pass
- Versions with `v` prefix `v2.0.0, v1.0.0` → pass
- Mixed `v` prefix and no prefix → handled correctly
- Pre-release ordering `2.0.0, 1.1.0-beta.1, 1.0.0` → pass (1.1.0-beta.1 < 1.1.0 per semver)

Implementation:
- Extract version strings from headers (skip Unreleased)
- Parse each with `internal/semver.ParseVersion`
- Verify each version is strictly greater than the next (descending order)
- Use `internal/semver.Compare` for comparison

### Step 7: Implement text-mode report rendering

**File**: `internal/changelog/validate.go`

Add a `FormatReport` function to render the report as human-readable text:

```
Changelog Validation: CHANGELOG.md
===================================

  ✅ Version headers         All version headers are properly formatted
  ✅ Duplicate entries        No duplicate entries found
  ❌ Broken links            2 issues found
     • Version [1.0.0] has no matching reference link
     • Orphan link [0.5.0] has no matching version header
  ✅ Entries without dates    All versions have dates
  ✅ Semver order            Versions are in correct descending order

Result: 4/5 checks passed, 1 failed
```

Tests:
- All-pass report formats correctly
- Mix of pass/fail formats correctly
- Details are indented under failures
- Summary line is accurate

### Step 8: Create command configuration metadata

**File**: `internal/config/commands/changelog_validate_config.go`

Minimal metadata — no custom config options needed since we reuse the parent `--file` flag:

```go
var ChangelogValidateMetadata = config.CommandMetadata{
    Use:          "validate",
    Short:        "Validate changelog for common problems",
    ConfigPrefix: "app.changelog.validate",
}
```

Note: This step is optional. Since the validate command has no custom flags beyond the inherited `--file`, we may not need a separate config file. The command can be created directly as a `cobra.Command` in `cmd/changelog_validate.go` following the same pattern as `changelog_add.go` (which also doesn't use `MustNewCommand`).

**Decision**: Skip the config file. Create the command directly like `changelog_add.go` does, since there are no command-specific flags.

### Step 9: Create the command wiring

**File**: `cmd/changelog_validate.go`, `cmd/changelog_validate_test.go`

Command wiring (must be ≤30 lines in `runValidateChangelog`):

```go
func runValidateChangelog(cmd *cobra.Command, args []string) error {
    file := getConfigValueWithFlags[string](cmd, "file", "app.changelog.file")
    data, err := os.ReadFile(file)
    if err != nil {
        return fmt.Errorf("failed to read changelog: %w", err)
    }
    report := changelog.ValidateChangelog(string(data), file)
    if report.Passed {
        return ui.RenderSuccess(cmd.OutOrStdout(), fmt.Sprintf("All %d checks passed", report.TotalRules), report)
    }
    // Print the report even on failure
    formatted := changelog.FormatReport(report)
    fmt.Fprint(cmd.OutOrStdout(), formatted)
    return fmt.Errorf("changelog validation failed: %d/%d checks failed", report.FailCount, report.TotalRules)
}
```

Command-level tests:
- Valid changelog file → exit 0, success message
- Invalid changelog file → exit non-zero, error details
- Missing file → appropriate error message
- `--output json` → JSON envelope with ValidationReport as data
- `--file custom.md` → reads the custom file

### Step 10: Run `task check` and fix any issues

- Run `task format` to fix formatting
- Run `task lint` to fix lint issues
- Run `task test` to verify all tests pass
- Run `task check` to run the full quality gate
- Fix any issues found

## Testing Strategy

### Unit Tests (`internal/changelog/validate_test.go`)

Table-driven tests for each validation rule. Tests operate on string content (no filesystem):

| Test Group | Count | Description |
|-----------|-------|-------------|
| `TestCheckVersionHeaders` | ~8 cases | Valid/invalid version header formats |
| `TestCheckDuplicateEntries` | ~6 cases | Duplicate detection across sections/versions |
| `TestCheckBrokenLinks` | ~7 cases | Header↔link matching, orphans |
| `TestCheckEntriesWithoutDates` | ~6 cases | Date presence and format |
| `TestCheckSemverOrder` | ~7 cases | Descending order, edge cases |
| `TestValidateChangelog` | ~4 cases | Integration: all rules together |
| `TestFormatReport` | ~3 cases | Text rendering of reports |
| `TestValidationReportJSONResponse` | ~2 cases | JSON serialization |

**Total**: ~43 test cases

### Command Tests (`cmd/changelog_validate_test.go`)

Test the command wiring with temp files:

| Test | Description |
|------|-------------|
| `TestValidateCommand_ValidFile` | Valid changelog → success output |
| `TestValidateCommand_InvalidFile` | Changelog with errors → failure output with details |
| `TestValidateCommand_FileNotFound` | Missing file → error message |
| `TestValidateCommand_CustomFile` | `--file` flag works |
| `TestValidateCommand_EmptyFile` | Empty file → appropriate failures |

### Coverage Target

- `internal/changelog/validate.go`: 95%+ (pure logic, fully testable)
- `cmd/changelog_validate.go`: 80%+ (command wiring)
- Overall project coverage: maintained above 85%

## Edge Cases

1. **Empty changelog file**: Should fail version headers check (no headers), pass duplicates (nothing to duplicate), etc.
2. **Changelog with only `[Unreleased]`**: All checks should pass (no versions to validate ordering on, no links needed for Unreleased if no versions exist — actually, `[Unreleased]` should have a link if there's a repo URL, but we'll be lenient here)
3. **Yanked versions**: `## [1.0.0] - 2024-01-01 [YANKED]` — the date check should still pass
4. **Pre-release versions**: `## [1.0.0-alpha.1] - 2024-01-01` — should be valid
5. **`v`-prefixed versions**: `## [v1.0.0] - 2024-01-01` — should be valid
6. **Windows line endings (`\r\n`)**: Normalize before processing
7. **No link reference section**: If there are version headers, broken links check should fail
8. **Changelog with content before first header**: Preamble text (title, description) is ignored
9. **Duplicate entries across different versions**: NOT considered duplicates (only within same version+section)
10. **Non-standard sections**: e.g., `### Custom` — not flagged (we only validate format, not section names)

## Dependencies

### Existing (no new dependencies)
- `github.com/blang/semver/v4` — already used via `internal/semver`
- `github.com/stretchr/testify` — test assertions
- `internal/semver` — version parsing and comparison
- `internal/ui` — output rendering
- `.ckeletin/pkg/output` — JSON mode support

### No new external dependencies needed

All validation is string parsing that can be done with the standard library (`regexp`, `strings`) plus the existing `internal/semver` package.

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Regex too strict (rejects valid changelogs) | Medium | High | Test against the project's own CHANGELOG.md as a golden test |
| Regex too lenient (misses problems) | Low | Medium | Comprehensive edge case tests |
| Performance on large files | Low | Low | Changelog files are typically <1000 lines |
| Breaking existing `changelog` command group | Low | High | Only adding a new subcommand, not modifying existing ones |

## Validation Against Project CHANGELOG.md

The project's own `CHANGELOG.md` (484 lines, 12 versions) should pass all 5 checks. Use it as a smoke test during development.
