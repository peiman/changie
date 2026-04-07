VERDICT: REQUEST CHANGES

# Code Review — Reviewer Fix Pass (coverage artifact, binaryName in config)

**Reviewed commits:** `2b18413..50d6e9d` (2 commits: reviewer fix pass + executor implementation pass)
**Files changed:** 5 files (`cmd/config.go`, `cmd/config_test.go`, `.gitignore`, `coverage.tmp` removed, `REVIEW.md`)
**Date:** 2026-04-07
**Previous review:** REVIEW.md at `4f2b585` (REQUEST CHANGES — 1 HIGH, 2 MEDIUM, 3 LOW)

---

## Summary

| Severity | Count |
|----------|-------|
| CRITICAL | 0 |
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 2 |

---

## Stage 1: Spec Compliance (PLAN.md)

### PLAN.md Required Actions Checklist

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 1 | Remove `coverage.tmp` from git tracking | ✅ DONE | `git ls-files coverage.tmp` returns nothing; file in `.gitignore` |
| 2 | Add `coverage.tmp` to `.gitignore` | ✅ DONE | Added at line 63 alongside `coverage.txt` |
| 3 | Fix `cmd/config.go:38-42` to use `binaryName` variable | ⚠️ BROKEN | Moved to `init()` per plan, but **init ordering bug** — see HIGH issue below |
| 4 | Fix `cmd/config_test.go:153` to use `"changie"` | ✅ DONE | Line 153 now reads `binaryName = "changie"` |
| 5 | Run `task check` | ✅ DONE | All 23 checks pass, 88.5% coverage |

**Spec deviation:** PLAN.md AD-3 correctly identified the risk: "binaryName empty at var-declaration time — HIGH." The plan prescribed moving the `Example` assignment to `init()`. The implementation followed this advice, but `config.go`'s `init()` runs BEFORE `root.go`'s `init()` (Go processes init functions in filename alphabetical order within a package: "config" < "root"). The result is that `binaryName` is still `""` when `config.go`'s `init()` executes `fmt.Sprintf`.

---

## Stage 2: Issues

### [HIGH] Init ordering bug — `binaryName` is empty in `config.go` init()

**File:** `cmd/config.go:44-49`
**Issue:** The `configValidateCmd.Example` is set in `config.go`'s `init()` using `fmt.Sprintf(..., binaryName, binaryName)`. However, `binaryName` is initialized to `"changie"` in `root.go`'s `init()` (line 272-273). In Go, `init()` functions within a package execute in **source file name alphabetical order**. Since `config.go` sorts before `root.go`, `config.go`'s `init()` runs first — when `binaryName` is still `""`.

**Verified by building and running:**
```
$ ./changie config validate --help
Examples:
  # Validate default config file
   config validate          ← MISSING binary name

  # Validate specific config file
   config validate --file /path/to/config.yaml   ← MISSING binary name
```

**Impact:** Users see broken example text with missing binary name in `changie config validate --help`.

**Fix:** Move the Example assignment to `root.go`'s `init()` function AFTER `binaryName` is resolved (around line 286), alongside the other `binaryName`-dependent assignments that are already there:

```go
// In cmd/root.go init(), after line 286:
configValidateCmd.Example = fmt.Sprintf(`  # Validate default config file
  %s config validate

  # Validate specific config file
  %s config validate --file /path/to/config.yaml`, binaryName, binaryName)
```

And remove the Example assignment from `cmd/config.go`'s `init()`.

**Alternative fix:** Keep it in `config.go` but use a `cobra.OnInitialize` callback which runs after all `init()` functions complete. However, the simpler approach is consolidating all `binaryName`-dependent assignments in `root.go`'s `init()`, which already has the pattern at lines 282-286.

**Note:** The pre-existing `cmd/completion.go:16-26` has the same bug — `binaryName` is captured at package-level `var` declaration time (even earlier than `init()`). The `completion --help` output also shows empty binary names. This was NOT introduced by the reviewed commits but is the same class of defect.

---

### [LOW] `cmd/dev_progress.go` exceeds 30-line command guidance (185 lines)

**File:** `cmd/dev_progress.go` (185 lines)
**Issue:** AGENTS.md mandates "Commands ≤30 lines." Contains demo business logic directly in cmd layer.
**Note:** Carried forward from previous reviews. Dev-only build-tagged file, not a production risk.
**Fix:** Extract demo logic into `internal/dev/progress_demo.go`. Lower priority.

---

### [LOW] Hardcoded `"now"` timestamp in ping response

**File:** `internal/ping/ping.go:86`
**Issue:** `Timestamp: "now"` instead of `time.Now().Format(time.RFC3339)`.
**Note:** Carried forward. Demo/example command.
**Fix:** Replace with actual timestamp or remove the field.

---

## Verification Evidence

| Check | Result |
|-------|--------|
| `task test` (1794 tests, 12 skipped) | ✅ PASS |
| `task lint` | ✅ PASS (no issues) |
| `task check` (all 23 gates) | ✅ PASS |
| Coverage | 88.5% (above 85% minimum) |
| `goimports -l internal/ui/ui_test.go` | ✅ No output (clean) |
| `git ls-files coverage.tmp` | ✅ Not tracked |
| `.gitignore` includes `coverage.tmp` | ✅ Line 63 |
| `grep ckeletin-go README.md` | ✅ No matches |
| `./changie --help \| grep ckeletin` | ✅ No matches |
| `./changie config validate --help` | ❌ Missing binary name in Example (see HIGH issue) |
| `./changie completion --help` | ❌ Missing binary name (pre-existing, not from these commits) |
| `cmd/config_test.go:153` uses `"changie"` | ✅ Confirmed |
| Hardcoded secrets scan | ✅ PASS |
| SAST scan | ✅ PASS |
| Vulnerability scan | ✅ PASS |

---

## Positive Observations

- **Previous HIGH issue fully resolved:** `coverage.tmp` is removed from git tracking and properly added to `.gitignore`. Clean execution.
- **Previous MEDIUM (config_test.go) resolved:** Test now uses `"changie"` matching production binary name.
- **Correct approach for config.go:** The *strategy* of using `fmt.Sprintf` with `binaryName` in `init()` is exactly right per PLAN.md AD-3. The issue is solely the init ordering between files, not the approach.
- **All quality gates green:** 1794 tests pass, lint clean, 88.5% coverage, security scans clear.
- **Clean working tree:** `git status` shows no uncommitted changes.

---

## Recommendation

**VERDICT: REQUEST CHANGES** — 1 HIGH issue: the `binaryName` init ordering bug in `cmd/config.go:44-49` produces broken help text for `changie config validate --help`. The binary name is missing from the Example output.

### Required Actions

1. **[HIGH]** Move `configValidateCmd.Example` assignment from `cmd/config.go` init() to `cmd/root.go` init() (after line 286 where `binaryName` is resolved). Verify with: `go build -o ./changie . && ./changie config validate --help` — should show `changie config validate` in examples.
2. **(Recommended)** Also fix `cmd/completion.go:16-26` — move `completionCmd.Long` assignment to `root.go` init() to fix the same pre-existing bug class.
3. Run `task check` to verify all issues are resolved.
4. Commit fixes.
