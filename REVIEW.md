VERDICT: REQUEST CHANGES

# Code Review — changie Production-Ready Changes

**Reviewed commits:** `c6ff157..98ba859` (3 commits: ckeletin-go port, quality gate fixes, production-ready push)
**Files changed:** 310 files, +52,694 / -2,838 lines
**Date:** 2026-04-07

---

## Summary

| Severity | Count |
|----------|-------|
| CRITICAL | 0 |
| HIGH | 1 |
| MEDIUM | 4 |
| LOW | 3 |

---

## Stage 1: Spec Compliance (PLAN.md)

### PLAN.md Checklist

| Requirement | Status | Notes |
|-------------|--------|-------|
| Step 1: Add tests to `internal/ui` (87.6% → 92%+) | ✅ DONE | `TestRunUI_SuccessPathWithAutoQuitProgram`, `TestNewDefaultUIRunner_LambdaIsCallable` added |
| Step 2: Add tests to `internal/version` (86.6% → 90%+) | ✅ DONE | `TestBump_CommitChangelogFailure` added |
| Step 3: Harden tests in `internal/check` (85.5% → 88%+) | ✅ DONE | JSON mode + executor tests present |
| Step 4: Rewrite README.md | ❌ PARTIAL | ckeletin-go reference remains on line 209 |
| Step 5: Run quality checks | ❌ FAIL | `internal/ui/ui_test.go` has formatting issues; integration test fails |
| Step 6: Commit and push | ✅ DONE | Committed as `98ba859` |

---

## Stage 2: Issues

### [HIGH] Test failure — `TestCheckCommand_QualityCategory` fails due to unformatted file

**File:** `internal/ui/ui_test.go` (formatting) → `test/integration/check_command_test.go:93`
**Issue:** `internal/ui/ui_test.go` needs formatting (`goimports` reports it as needing changes). This causes the integration test `TestCheckCommand_QualityCategory` to fail because the `check --category quality` command detects the formatting violation and exits non-zero.
**Impact:** `go test ./...` exits with FAIL. CI will reject this commit.
**Fix:** Run `task format` to fix `internal/ui/ui_test.go`, then re-run `task check` and commit.

---

### [MEDIUM] Spec non-compliance — README.md still references ckeletin-go framework

**File:** `README.md:209`
**Issue:** PLAN.md Step 4 explicitly states: "Remove all ckeletin-go/scaffold references." The README still contains:
```
changie is built on the [ckeletin-go](https://github.com/peiman/ckeletin-go) framework.
```
**Fix:** Remove or replace line 209 with contributor-focused text that does not reference the scaffold framework, e.g., replace the Development section intro with a standalone description.

---

### [MEDIUM] Spec non-compliance — RootCmd.Long references ckeletin-go

**File:** `cmd/root.go:284`
**Issue:** The Cobra root command's `Long` description says:
```go
"Powered by the ckeletin-go framework with Cobra, Viper, and Zerolog."
```
This appears in `changie --help` output, leaking internal framework details to end users.
**Fix:** Replace with a user-facing description, e.g.:
```go
"%s is a CLI tool for managing semantic versioning and Keep a Changelog format changelogs."
```

---

### [MEDIUM] Example text references wrong binary name

**File:** `cmd/config.go:39-42`
**Issue:** The `Example` field in `configValidateCmd` shows:
```
ckeletin-go config validate
ckeletin-go config validate --file /path/to/config.yaml
```
Users will see `ckeletin-go` in `changie config validate --help`.
**Fix:** Replace `ckeletin-go` with the `binaryName` variable or the literal `changie`.

---

### [MEDIUM] cmd/dev_progress.go contains 185 lines of business logic

**File:** `cmd/dev_progress.go`
**Issue:** AGENTS.md rule: "Commands ≤30 lines — `cmd/*.go` files wire things together; logic goes in `internal/`." The `demoSpinner`, `demoProgressBar`, and `demoMultiPhase` functions (lines 114-185) contain demo business logic directly in the cmd layer instead of delegating to `internal/progress/` or a dedicated `internal/dev/` package.
**Fix:** Extract demo logic into `internal/dev/progress_demo.go` and have `cmd/dev_progress.go` call it. Since this is a `dev`-only build-tag file, this is a guideline violation rather than a production risk—lower urgency.

---

### [LOW] Hardcoded "now" timestamp in ping response

**File:** `internal/ping/ping.go:87`
**Issue:** `Timestamp: "now"` is hardcoded instead of using `time.Now().Format(time.RFC3339)`. The inline comment acknowledges this: `// In a real app, use time.Now().Format(time.RFC3339)`.
**Fix:** Replace with actual timestamp or remove the `Timestamp` field if this is truly a demo command.

---

### [LOW] panic() in MustNewCommand

**File:** `cmd/helpers.go:71`
**Issue:** `MustNewCommand` uses `panic(err)` which is appropriate for `init()` functions but could be surprising if ever called at runtime.
**Fix:** No change needed — the function's doc comment already explains the `init()`-only usage. Consider adding a `// SAFETY: ...` comment to make the panic rationale even clearer.

---

### [LOW] cmd/bump.go and cmd/init.go exceed 30-line guidance

**File:** `cmd/bump.go` (145 lines), `cmd/init.go` (146 lines)
**Issue:** These files exceed the ≤30-line guidance for cmd files. However, much of the content is Cobra command struct definitions, flags, and thin wiring — the actual `RunE` bodies delegate to `internal/`.
**Fix:** Consider moving the lengthy flag definitions and validation to `internal/config/commands/` for consistency, but this is not blocking.

---

## Positive Observations

- **Excellent test architecture**: Table-driven tests, dependency injection (UIRunner interface, programFactory), temp git repo helpers, error writers — all best practices.
- **Thorough bump_test.go**: 12+ test cases covering major/minor/patch, branch restrictions, uncommitted changes, auto-push, tag conflicts, commit failures, and integration workflows. Well done.
- **Security-conscious config handling**: `ValidateConfigFileSecurity()` checks both file size (DoS prevention) and permissions (world-writable detection). This is above-average for CLI tools.
- **Clean error messages**: Every error in `internal/version/bump.go` includes actionable guidance (e.g., "ensure git is properly configured", "use 'git tag' to list existing tags").
- **Good framework/app separation**: `.ckeletin/` framework layer is cleanly separated from application code in `cmd/`, `internal/`, `pkg/`.
- **Structured logging**: Proper use of zerolog with component tags, debug-level for returnable errors, and info for important events.
- **JSON output mode**: The `--output json` flag with proper envelope format is excellent for CI/CD integration.
- **`go vet` passes cleanly**: No type safety issues.
- **No hardcoded secrets**: Grep found no API keys, passwords, or tokens in Go source files.

---

## Recommendation

**VERDICT: REQUEST CHANGES** — 1 HIGH issue (test failure) must be fixed before merge. The 4 MEDIUM issues (ckeletin-go references in README, root command help, config example, and dev_progress.go line count) should be addressed in the same commit.

### Required Actions
1. Run `task format` to fix `internal/ui/ui_test.go` formatting
2. Remove ckeletin-go reference from `README.md:209`
3. Fix `cmd/root.go:284` — remove ckeletin-go from Long description
4. Fix `cmd/config.go:39-42` — replace `ckeletin-go` with `changie` in examples
5. Run `task check` to verify all issues are resolved
6. Commit fixes
