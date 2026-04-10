# PLAN: `changie diff` command

## Status: IMPLEMENTED — Needs Review & Hardening

The `changie diff` command is **already implemented** and all tests pass (22/22). This plan documents the architecture, verifies correctness, and identifies hardening work for the Executor.

---

## 1. Requirements

| Requirement | Status | Evidence |
|---|---|---|
| `changie diff v1.0.0 v1.1.0` shows changelog entries between two versions | ✅ Done | `cmd/diff.go:16-29`, `internal/changelog/diff.go:29-91` |
| Extracts entries after FROM (exclusive) up to TO (inclusive) | ✅ Done | `internal/changelog/diff_test.go:81-92` |
| Supports `v`-prefixed and bare version strings | ✅ Done | `internal/changelog/diff_test.go:94-126` |
| Validates semver format | ✅ Done | `internal/changelog/diff_test.go:158-170` |
| Rejects inverted / same versions | ✅ Done | `internal/changelog/diff_test.go:144-156` |
| Reports missing versions clearly | ✅ Done | `internal/changelog/diff_test.go:128-142` |
| Custom `--file` flag for non-default changelog | ✅ Done | `cmd/diff_test.go:107-131` |
| Unit tests | ✅ Done | 16 internal + 6 cmd = 22 tests, all passing |
| Coverage ≥ 85% | ✅ Done | `internal/changelog`: 94.2%, `cmd`: 90.3% |

---

## 2. Architecture Decisions

### 2.1 Top-level command (not subcommand of `changelog`)

The `diff` command is registered directly on `RootCmd` (`cmd/diff.go:37`), yielding `changie diff FROM TO`. This matches the user's specification and keeps the UX short. It is **not** placed under `changie changelog diff` because `diff` is a frequent operation that benefits from a short invocation path.

### 2.2 Ultra-thin command pattern (ADR-001)

`cmd/diff.go` is wiring only — `runDiff` is 12 lines. All business logic is in `internal/changelog/diff.go`. This follows the ≤30-line command rule.

### 2.3 Shared regex from validate.go

`diff.go` reuses `reVersionHeader` (compiled in `internal/changelog/validate.go:40`) rather than defining its own regex. This avoids duplication and ensures consistent version header matching across diff and validate features.

### 2.4 Semver comparison via `internal/semver`

Version parsing and comparison uses the project's `internal/semver.ParseVersion()` → `blang/semver.Compare()` chain (`diff.go:31-43`), keeping semver handling consistent across the codebase.

---

## 3. Module / File Structure

```
cmd/
  diff.go              # Command wiring (52 lines, runDiff=12 lines)  ✅ EXISTS
  diff_test.go         # Command-layer tests (6 tests)                ✅ EXISTS

internal/changelog/
  diff.go              # DiffVersions() business logic (91 lines)     ✅ EXISTS
  diff_test.go         # Business-logic tests (16 tests)              ✅ EXISTS
  validate.go          # Shared reVersionHeader regex                  ✅ EXISTS (no changes)
```

No new files need to be created.

---

## 4. Implementation Tasks for the Executor

All core implementation is complete. The following tasks are **hardening and cleanup** items.

### Task 1: Fix duplicate viper binding for `app.changelog.file`

**Priority: HIGH** — Potential runtime conflict.

`cmd/diff.go:34` binds its own `--file` flag to `app.changelog.file`. But `cmd/changelog.go:24` already binds a *PersistentFlag* with the same viper key. Since both `init()` functions run at package load, the last one to execute wins the binding. Go's `init()` order within a package is alphabetical by filename, so `changelog.go` runs before `diff.go`, meaning `diff.go`'s binding overwrites it.

**Fix:** Remove the `viper.BindPFlag` call from `cmd/diff.go:34` entirely. Instead, rely on `getConfigValueWithFlags` (which checks the flag directly, then falls back to viper) — this is already how `runDiff` retrieves the value at line 41. The binding in `init()` is redundant given the retrieval pattern.

Alternatively, if the binding is kept for env-var support, change to a unique viper key like `app.diff.file` (with its own default), or share via the existing `changelog` binding by making `diff` a subcommand of `changelog`. The simplest fix is removing the redundant `BindPFlag`.

```go
// cmd/diff.go init() — REMOVE these lines:
// if err := viper.BindPFlag("app.changelog.file", diffCmd.Flags().Lookup("file")); err != nil {
//     log.Fatal().Err(err).Msg("Failed to bind 'file' flag")
// }
```

### Task 2: Add `--output json` support

**Priority: MEDIUM** — Consistency with other commands.

The `validate` command (`cmd/changelog_validate.go:43`) renders JSON via `ui.RenderSuccess`. The `diff` command outputs raw text via `fmt.Fprintln`. For agent/CI consumers, add JSON envelope support:

1. In `cmd/diff.go`, check `output.IsJSONMode()`.
2. If JSON, wrap the result in `output.JSONEnvelope{Status: "success", Command: "diff", Data: map[string]string{"diff": result}}`.
3. Use `output.RenderJSON(cmd.OutOrStdout(), envelope)`.

This is optional — the current text output is functional. Prioritize only if CI/agent integration is planned.

### Task 3: Add Windows line-ending test

**Priority: LOW** — `DiffVersions` already normalizes `\r\n` at `internal/changelog/diff.go:51`, but there is no test verifying this path.

```go
func TestDiffVersions_WindowsLineEndings(t *testing.T) {
    content := strings.ReplaceAll(sampleChangelog, "\n", "\r\n")
    result, err := DiffVersions(content, "1.0.0", "1.2.0")

    require.NoError(t, err)
    assert.Contains(t, result, "Feature C")
    assert.Contains(t, result, "Feature B")
}
```

### Task 4: Run `task check` to verify full quality gate

**Priority: HIGH** — Mandatory before any commit.

```bash
task check
```

This validates formatting, linting, architecture rules, security scanning, license compliance, and full test suite with race detection.

---

## 5. Testing Strategy

### Existing Coverage (all passing)

**`internal/changelog/diff_test.go` (16 tests):**

| Test | Category |
|---|---|
| `HappyPath_TwoAdjacentVersions` | Happy path |
| `MultipleVersionsInRange` | Happy path — multi-version span |
| `WithVPrefix` | v-prefix normalization |
| `WithoutVPrefix` | Bare version |
| `MixedPrefix` | Mixed v-prefix (user vs header) |
| `ChangelogWithVPrefixedHeaders` | Changelog uses `[v1.0.0]` format |
| `FromVersionNotFound` | Error — missing FROM |
| `ToVersionNotFound` | Error — missing TO |
| `InvertedVersions` | Error — FROM > TO |
| `SameVersion` | Error — FROM == TO |
| `InvalidSemverFrom` | Error — bad FROM format |
| `InvalidSemverTo` | Error — bad TO format |
| `EmptyContent` | Error — empty changelog |
| `OnlyUnreleased` | Error — no versioned sections |
| `MultipleSectionsPerVersion` | Multiple subsections (Added, Fixed) |
| `UnreleasedNotIncluded` | Unreleased section excluded |

**`cmd/diff_test.go` (6 tests):**

| Test | Category |
|---|---|
| `HappyPath` | End-to-end via command |
| `FileNotFound` | Error — missing file |
| `VersionNotFound` | Error — invalid version |
| `InvertedVersions` | Error — wrong order |
| `CustomFileFlag` | `--file HISTORY.md` |
| `CommandRegistered` | Command wiring verification |

### Recommended Additional Tests (Task 3)

- Windows line endings (`\r\n`)
- Large changelog (50+ versions) — performance sanity
- Pre-release versions (e.g., `1.0.0-beta.1`)

---

## 6. Edge Cases

| Edge Case | Handled | Location |
|---|---|---|
| `v`-prefix on user input | ✅ | `diff.go:55-56` — strips prefix for comparison |
| `v`-prefix in changelog headers | ✅ | `diff.go:68` — normalizes header prefix |
| Mixed prefix (user=`v1.0.0`, header=`1.0.0`) | ✅ | `diff_test.go:110-117` |
| FROM == TO (same version) | ✅ | `diff.go:44-45` — explicit error |
| FROM > TO (inverted order) | ✅ | `diff.go:47-48` — explicit error |
| Invalid semver (not a version) | ✅ | `diff.go:31-38` — ParseVersion errors |
| Version not in changelog | ✅ | `diff.go:77-82` — not-found errors |
| Empty changelog content | ✅ | Falls through to "not found" |
| Only `[Unreleased]` section | ✅ | Falls through to "not found" |
| Windows line endings (`\r\n`) | ✅ (code) | `diff.go:51` — normalized, but **untested** |
| Multiple sections per version | ✅ | `diff_test.go:192-201` |
| Unreleased content excluded | ✅ | `diff_test.go:203-209` |
| Custom changelog filename | ✅ | `cmd/diff_test.go:107-131` |
| Pre-release versions (e.g., `1.0.0-rc.1`) | ⚠️ Untested | `blang/semver` supports it, but no test |

---

## 7. Summary for Executor

The implementation is complete and passing. Priority work:

1. **Fix the duplicate viper binding** in `cmd/diff.go:34` (HIGH — remove the `BindPFlag` call or use a unique key)
2. **Run `task check`** to verify all quality gates pass (HIGH — mandatory before commit)
3. **Add CRLF test** for `DiffVersions` (LOW — covers existing `\r\n` normalization code)
4. **Add `--output json` support** (MEDIUM — only if CI/agent integration is a priority)
