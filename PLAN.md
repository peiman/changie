# PLAN.md — Fix REVIEW.md Issues and Rewrite README

## Summary

REVIEW.md identified 1 HIGH and 3 MEDIUM issues. This plan addresses all four required actions plus a full README rewrite. The existing test additions (Steps 1-3 from the previous plan) are already committed at `98ba859`; the formatting issue was introduced by those additions.

---

## Requirements Analysis

From REVIEW.md Required Actions:

| # | Issue | Severity | File:Line | What to Fix |
|---|-------|----------|-----------|-------------|
| 1 | Test formatting failure | HIGH | `internal/ui/ui_test.go:21-23` | `goimports`/`gofmt` reports alignment padding in `autoQuitModel` methods is non-standard — extra spaces used for visual alignment that `gofmt` doesn't produce |
| 2 | README references ckeletin-go | MEDIUM | `README.md:209` | Line says "changie is built on the ckeletin-go framework" — must be removed |
| 3 | Root command help leaks framework | MEDIUM | `cmd/root.go:284` | `Long` description says "Powered by the ckeletin-go framework with Cobra, Viper, and Zerolog." — visible in `changie --help` |
| 4 | Config example uses wrong binary | MEDIUM | `cmd/config.go:39-42` | Example field shows `ckeletin-go config validate` instead of `changie` |
| 5 | Full README rewrite | PLAN Step 4 | `README.md` | Rewrite for standalone changie identity |

### Scope Boundary

**In scope:** The 4 required actions from REVIEW.md + README rewrite.

**Out of scope (acknowledged but deferred):**
- `cmd/dev_progress.go` 185-line violation (MEDIUM, dev-only, not blocking)
- `internal/ping/ping.go:87` hardcoded timestamp (LOW, demo command)
- `cmd/helpers.go:71` panic comment (LOW, already documented)
- `cmd/bump.go`/`cmd/init.go` line count (LOW, mostly config wiring)
- Other `ckeletin-go` references in `.ckeletin/` framework files, `test/integration/`, `internal/check/executor.go`, `internal/check/timing.go`, `internal/dev/config_test.go`, `cmd/config_test.go` — these are internal/framework/test plumbing, not user-facing

---

## Architecture Decisions

### AD-1: Formatting Fix Strategy

The `autoQuitModel` methods at `internal/ui/ui_test.go:21-23` use manual column-alignment:
```go
func (a autoQuitModel) Init() tea.Cmd                           { return tea.Quit }
func (a autoQuitModel) Update(tea.Msg) (tea.Model, tea.Cmd)    { return a, nil }
func (a autoQuitModel) View() string                           { return "" }
```

`gofmt` wants:
```go
func (a autoQuitModel) Init() tea.Cmd                       { return tea.Quit }
func (a autoQuitModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return a, nil }
func (a autoQuitModel) View() string                        { return "" }
```

**Decision:** Run `task format` which invokes `goimports` and resolves this automatically. Do NOT manually edit — let the canonical formatter handle it to avoid introducing new issues.

### AD-2: Root Command Long Description

Current (`cmd/root.go:283-284`):
```go
RootCmd.Long = fmt.Sprintf(`%s is a CLI tool for managing semantic versioning and Keep a Changelog format changelogs.
Powered by the ckeletin-go framework with Cobra, Viper, and Zerolog.`, binaryName)
```

**Decision:** Remove the "Powered by" line entirely. The Long description should be user-facing, not an attribution line. Replace with a concise feature summary that helps users understand what the tool does:

```go
RootCmd.Long = fmt.Sprintf(`%s is a CLI tool for managing semantic versioning and Keep a Changelog format changelogs.

It initializes changelogs, adds entries to the correct section, bumps versions following
semver, and updates the changelog, commits, tags, and optionally pushes — all in one command.`, binaryName)
```

**Trade-off:** Removing framework attribution is appropriate for user-facing help text. The `.ckeletin/` directory and AGENTS.md still document the framework relationship for developers.

### AD-3: Config Example Fix

Current (`cmd/config.go:38-42`):
```go
Example: `  # Validate default config file
  ckeletin-go config validate

  # Validate specific config file
  ckeletin-go config validate --file /path/to/config.yaml`,
```

**Decision:** Use the `binaryName` variable via `fmt.Sprintf` instead of hardcoding. This ensures the example is always correct regardless of binary name. This is the pattern already used in `cmd/completion.go:16-26`.

```go
Example: fmt.Sprintf(`  # Validate default config file
  %s config validate

  # Validate specific config file
  %s config validate --file /path/to/config.yaml`, binaryName, binaryName),
```

### AD-4: README Rewrite Strategy

The current README is 80% correct for changie. The rewrite should:

1. **Keep** the excellent structure (What it does, Install, Quick start, Commands, Configuration, Changelog format)
2. **Remove** the single ckeletin-go reference on line 209
3. **Add** a "Why changie?" section for quick value proposition
4. **Add** `check`, `completion`, and `ping` command docs (they exist but aren't documented)
5. **Add** CI/CD integration examples (the JSON output mode deserves highlighting)
6. **Rewrite** Development section to be contributor-focused without framework references
7. **NOT add** `check`/`ping`/`dev`/`docs` since those are dev-build-tag commands — users won't see them

**Trade-off:** A minimal fix (just delete line 209) would address the REVIEW.md item, but the task says "rewrite README fully for changie," so a comprehensive rewrite is warranted.

---

## Step-by-Step Implementation Tasks

### Step 1: Fix Test Formatting (HIGH — Must Be First)

**Action:** Run `task format`

**Why first:** The integration test `TestCheckCommand_QualityCategory` detects formatting violations and will fail until this is fixed. All subsequent `task check` runs depend on this being resolved.

**Files affected:** `internal/ui/ui_test.go`

**What changes:** Lines 21-23 — the `autoQuitModel` method alignment will be reformatted by `goimports`:
```
BEFORE: func (a autoQuitModel) Init() tea.Cmd                           { return tea.Quit }
AFTER:  func (a autoQuitModel) Init() tea.Cmd                       { return tea.Quit }
```
(3 lines change, spacing adjustment only)

**Verification:** `goimports -l internal/ui/ui_test.go` should output nothing after formatting.

---

### Step 2: Fix `cmd/root.go:283-284` — Remove ckeletin-go from Help Output

**File:** `cmd/root.go`
**Line:** 283-284

**Change:** Replace the current `Long` description:

```go
// BEFORE (line 283-284):
RootCmd.Long = fmt.Sprintf(`%s is a CLI tool for managing semantic versioning and Keep a Changelog format changelogs.
Powered by the ckeletin-go framework with Cobra, Viper, and Zerolog.`, binaryName)

// AFTER:
RootCmd.Long = fmt.Sprintf(`%s is a CLI tool for managing semantic versioning and Keep a Changelog format changelogs.

It initializes changelogs, adds entries to the correct section, bumps versions following
semver, and updates the changelog, commits, tags, and optionally pushes — all in one command.`, binaryName)
```

**Verification:** `go build ./... && ./changie --help` — confirm no "ckeletin" appears in output.

---

### Step 3: Fix `cmd/config.go:38-42` — Replace ckeletin-go in Example

**File:** `cmd/config.go`
**Lines:** 22 (the `configValidateCmd` struct), 38-42 (the `Example` field)

**Change:** Replace the hardcoded `Example` string with a `fmt.Sprintf` call:

```go
// BEFORE (line 38-42):
	Example: `  # Validate default config file
  ckeletin-go config validate

  # Validate specific config file
  ckeletin-go config validate --file /path/to/config.yaml`,

// AFTER:
	Example: fmt.Sprintf(`  # Validate default config file
  %s config validate

  # Validate specific config file
  %s config validate --file /path/to/config.yaml`, binaryName, binaryName),
```

**Note:** The `configValidateCmd` is defined as a package-level `var`. Since `binaryName` is also package-level and initialized in `init()`, and Go evaluates `var` declarations before `init()`, we need to move the `Example` assignment into `init()`. Check whether `binaryName` is available at `var` declaration time.

**Actually:** Looking at `cmd/root.go:272-274`, `binaryName` defaults to `""` and is set to `"changie"` inside `init()`. So `configValidateCmd` at var-declaration time would see `binaryName = ""`. The `Example` assignment must happen in `init()`:

```go
func init() {
    // Set Example after binaryName is resolved
    configValidateCmd.Example = fmt.Sprintf(`  # Validate default config file
  %s config validate

  # Validate specific config file
  %s config validate --file /path/to/config.yaml`, binaryName, binaryName)

    configCmd.AddCommand(configValidateCmd)
    // ...
}
```

**Verification:** `go build ./... && ./changie config validate --help` — confirm "changie" appears in examples, not "ckeletin-go".

---

### Step 4: Rewrite README.md

**File:** `README.md`

**Full rewrite.** Preserve accurate content, remove all ckeletin-go references, add value proposition. Structure:

```
# changie
[tagline + badges]

## Why changie?
- 5 bullet points on value proposition

## Install
- go install
- Build from source

## Quick start
- 6-line example (init → add entries → bump)

## Commands
### changie init
### changie changelog <section> <content>
### changie bump <major|minor|patch>
### changie config validate
### changie completion

## Configuration
### Config file
### Environment variables
### JSON output

## Changelog format
[Keep the existing excellent example]

## Development
[Contributor-focused, no ckeletin-go mention]

## License
```

**Key content decisions:**
- Do NOT document `check`, `ping`, `dev`, `docs` — they require dev build tag and aren't user-facing
- DO document `completion` — it's user-facing and available in production builds
- Remove "changie is built on the ckeletin-go framework" (line 209)
- Keep the changelog format example — it's genuinely useful
- Add a note about `--output json` for CI/CD in the Quick Start or Configuration section

**Content for the README (full text):**

The README should contain:

1. **Header:** "changie" with a one-line tagline: "Changelog management and semantic versioning, automated."
2. **Badges:** Keep the 4 existing badges (CI, Go Report Card, License, Go Version)
3. **Why changie?** section with 5 value bullets
4. **Install:** `go install` + build from source
5. **Quick start:** Same example as current (it's good)
6. **Commands:** `init`, `changelog`, `bump`, `config validate`, `completion`
7. **Configuration:** Config file, env vars, JSON output
8. **Changelog format:** Keep the current example
9. **Development:** task commands, link to CONTRIBUTING.md, no framework reference
10. **License:** MIT

---

### Step 5: Run Quality Checks

```bash
task format    # Should be a no-op after Step 1
task check     # Must pass completely
```

**Expected outcome:** All checks pass. If any fail, fix before proceeding.

---

### Step 6: Stage and Verify

```bash
# Verify no ckeletin-go in user-facing files
grep -n "ckeletin-go" README.md          # Should return nothing
grep -n "ckeletin-go" cmd/root.go        # Should only be in import paths, not strings
grep -n "ckeletin-go" cmd/config.go      # Should only be in comment line 3, not Example

# Build and verify help output
go build ./...
./changie --help | grep -i ckeletin      # Should return nothing
./changie config validate --help | grep -i ckeletin  # Should return nothing
```

---

## Testing Strategy

### What Needs Testing

This change is primarily cosmetic (string changes + formatting). The existing test suite validates all behavior. No new tests are needed.

### Verification Checklist

| Check | Command | Expected |
|-------|---------|----------|
| Formatting passes | `goimports -l ./...` | No output |
| All tests pass | `task test` | PASS |
| Full quality gate | `task check` | PASS |
| No ckeletin-go in help | `./changie --help` | No "ckeletin" |
| Config example correct | `./changie config validate --help` | Shows "changie" |
| README has no framework ref | `grep ckeletin-go README.md` | No output |

### Existing Tests That Validate These Changes

- `cmd/root_test.go` — Tests root command behavior (will exercise new `Long` text)
- `cmd/config_test.go` — Tests config command (will exercise new `Example`)
- `test/integration/check_command_test.go:93` — Tests quality checks (will pass once formatting is fixed)

---

## Edge Cases

| Edge Case | Risk | Mitigation |
|-----------|------|------------|
| `binaryName` empty at var-declaration time | HIGH — `fmt.Sprintf` with empty string produces wrong Example | Move Example assignment to `init()` where `binaryName` is resolved |
| `configValidateCmd.Example` set in `init()` runs after command registration | LOW — Go `init()` in same file runs top-to-bottom | `init()` in `config.go` sets Example before `MustAddToRoot` |
| README links broken after rewrite | LOW | Verify all `[text](url)` links point to existing files |
| Test still fails after `task format` | LOW | Verify with `goimports -l` before committing |
| Other ckeletin-go string references in non-user-facing code | NONE | Intentionally out of scope — `.ckeletin/` framework files, integration test plumbing, and internal comments are not user-facing |

---

## Execution Order

```
1. task format                           ~1 min  (fixes HIGH issue)
2. Edit cmd/root.go:283-284             ~2 min  (MEDIUM #1)
3. Edit cmd/config.go:38-42 + init()    ~3 min  (MEDIUM #2)
4. Rewrite README.md                    ~15 min (full rewrite)
5. task check                           ~3 min  (verify everything)
6. Manual verification (help output)    ~2 min
7. git add + commit                     ~1 min
```

Total estimated time: ~27 minutes

---

## Files Modified (Complete List)

| File | Change Type | Lines Changed |
|------|-------------|---------------|
| `internal/ui/ui_test.go` | Format fix | 3 lines (whitespace only) |
| `cmd/root.go` | String edit | 2 lines |
| `cmd/config.go` | String edit + move to init() | ~8 lines |
| `README.md` | Full rewrite | ~180 lines |
