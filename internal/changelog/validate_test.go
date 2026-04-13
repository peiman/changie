package changelog

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- ValidateChangelog integration tests ---

func TestValidateChangelog_AllPass(t *testing.T) {
	content := `# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- Some upcoming feature

## [1.1.0] - 2024-02-01

### Added

- Another feature

## [1.0.0] - 2024-01-01

### Fixed

- A bug fix

[Unreleased]: https://github.com/user/repo/compare/1.1.0...HEAD
[1.1.0]: https://github.com/user/repo/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
`
	report := ValidateChangelog(content, "CHANGELOG.md")
	assert.Equal(t, "CHANGELOG.md", report.File)
	assert.True(t, report.Passed)
	assert.Greater(t, report.TotalRules, 0)
	assert.Equal(t, report.TotalRules, report.PassCount)
	assert.Equal(t, 0, report.FailCount)
	assert.Len(t, report.Results, report.TotalRules)
}

func TestValidateChangelog_WithFailures(t *testing.T) {
	// Changelog with broken links and out-of-order versions
	content := `# Changelog

## [1.0.0] - 2024-01-01

- Some entry

## [1.1.0] - 2024-02-01

- Another entry

[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
[1.1.0]: https://github.com/user/repo/releases/tag/1.1.0
`
	report := ValidateChangelog(content, "CHANGELOG.md")
	assert.False(t, report.Passed)
	assert.Greater(t, report.TotalRules, 0)
	assert.Greater(t, report.FailCount, 0)
	assert.Equal(t, report.PassCount+report.FailCount, report.TotalRules)
}

func TestValidateChangelog_EmptyContent(t *testing.T) {
	report := ValidateChangelog("", "CHANGELOG.md")
	assert.NotNil(t, report)
	assert.Greater(t, report.TotalRules, 0)
}

// --- ValidationReport JSON tests ---

func TestValidationReportJSONResponse(t *testing.T) {
	report := &ValidationReport{
		File:       "CHANGELOG.md",
		Passed:     true,
		TotalRules: 3,
		PassCount:  3,
		FailCount:  0,
		Results:    []ValidationResult{{Name: "Test", Passed: true, Message: "ok"}},
	}
	data := report.JSONResponse()
	assert.Equal(t, report, data)

	// Verify JSON serialization
	b, err := json.Marshal(report)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"file":"CHANGELOG.md"`)
	assert.Contains(t, string(b), `"passed":true`)
	assert.Contains(t, string(b), `"total_rules":3`)
}

func TestValidationResultJSONOmitEmptyDetails(t *testing.T) {
	result := ValidationResult{
		Name:    "Version headers",
		Passed:  true,
		Message: "All good",
	}
	b, err := json.Marshal(result)
	require.NoError(t, err)
	assert.NotContains(t, string(b), "details")
}

// --- checkVersionHeaders tests ---

func TestCheckVersionHeaders(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantPassed bool
		wantDetail string
	}{
		{
			name: "valid versions pass",
			content: `## [Unreleased]

## [1.0.0] - 2024-01-01
## [1.1.0] - 2024-02-01
`,
			wantPassed: true,
		},
		{
			name: "missing brackets fails",
			content: `## 1.0.0 - 2024-01-01
`,
			wantPassed: false,
			wantDetail: "1.0.0",
		},
		{
			name: "invalid semver fails",
			content: `## [abc] - 2024-01-01
`,
			wantPassed: false,
			wantDetail: "abc",
		},
		{
			name: "Unreleased is not flagged",
			content: `## [Unreleased]
`,
			wantPassed: true,
		},
		{
			name: "prerelease version passes",
			content: `## [1.0.0-beta.1] - 2024-01-01
`,
			wantPassed: true,
		},
		{
			name: "v-prefixed version passes",
			content: `## [v1.0.0] - 2024-01-01
`,
			wantPassed: true,
		},
		{
			name: "mix of valid and invalid",
			content: `## [1.0.0] - 2024-01-01
## [bad-version] - 2024-02-01
## [2.0.0] - 2024-03-01
`,
			wantPassed: false,
			wantDetail: "bad-version",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines := strings.Split(tc.content, "\n")
			result := checkVersionHeaders(lines)
			assert.Equal(t, tc.wantPassed, result.Passed, "result.Passed mismatch")
			if tc.wantDetail != "" {
				found := false
				for _, d := range result.Details {
					if strings.Contains(d, tc.wantDetail) {
						found = true
						break
					}
				}
				assert.True(t, found, "expected detail containing %q, got %v", tc.wantDetail, result.Details)
			}
		})
	}
}

// --- checkDuplicateEntries tests ---

func TestCheckDuplicateEntries(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantPassed bool
		wantDetail string
	}{
		{
			name: "no duplicates passes",
			content: `## [1.0.0] - 2024-01-01

### Added

- Feature A
- Feature B
`,
			wantPassed: true,
		},
		{
			name: "duplicate in same section fails",
			content: `## [1.0.0] - 2024-01-01

### Added

- Feature A
- Feature A
`,
			wantPassed: false,
			wantDetail: "Feature A",
		},
		{
			name: "same text in different sections passes",
			content: `## [1.0.0] - 2024-01-01

### Added

- Feature A

### Fixed

- Feature A
`,
			wantPassed: true,
		},
		{
			name: "same text in different versions passes",
			content: `## [2.0.0] - 2024-02-01

### Added

- Feature A

## [1.0.0] - 2024-01-01

### Added

- Feature A
`,
			wantPassed: true,
		},
		{
			name: "case sensitive - different case is not duplicate",
			content: `## [1.0.0] - 2024-01-01

### Added

- Feature A
- feature a
`,
			wantPassed: true,
		},
		{
			name: "whitespace normalized - trailing space is duplicate",
			content: `## [1.0.0] - 2024-01-01

### Added

- Feature A
- Feature A
`,
			wantPassed: false,
			wantDetail: "Feature A",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines := strings.Split(tc.content, "\n")
			result := checkDuplicateEntries(lines)
			assert.Equal(t, tc.wantPassed, result.Passed, "result.Passed mismatch")
			if tc.wantDetail != "" {
				found := false
				for _, d := range result.Details {
					if strings.Contains(d, tc.wantDetail) {
						found = true
						break
					}
				}
				assert.True(t, found, "expected detail containing %q, got %v", tc.wantDetail, result.Details)
			}
		})
	}
}

// --- checkBrokenLinks tests ---

func TestCheckBrokenLinks(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantPassed bool
		wantDetail string
	}{
		{
			name: "all headers have links passes",
			content: `## [Unreleased]

## [1.0.0] - 2024-01-01

[Unreleased]: https://github.com/user/repo/compare/1.0.0...HEAD
[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
`,
			wantPassed: true,
		},
		{
			name: "version header without link fails",
			content: `## [1.0.0] - 2024-01-01

[Unreleased]: https://github.com/user/repo/compare/1.0.0...HEAD
`,
			wantPassed: false,
			wantDetail: "1.0.0",
		},
		{
			name: "orphan link fails",
			content: `## [1.0.0] - 2024-01-01

[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
[0.5.0]: https://github.com/user/repo/releases/tag/0.5.0
`,
			wantPassed: false,
			wantDetail: "0.5.0",
		},
		{
			name: "Unreleased with link passes",
			content: `## [Unreleased]

## [1.0.0] - 2024-01-01

[Unreleased]: https://github.com/user/repo/compare/1.0.0...HEAD
[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
`,
			wantPassed: true,
		},
		{
			name: "Unreleased header without link fails",
			content: `## [Unreleased]

## [1.0.0] - 2024-01-01

[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
`,
			wantPassed: false,
			wantDetail: "Unreleased",
		},
		{
			name: "no links section with version headers fails",
			content: `## [1.0.0] - 2024-01-01

- Some entry
`,
			wantPassed: false,
		},
		{
			name: "no headers and no links passes",
			content: `# Changelog

Just some text.
`,
			wantPassed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines := strings.Split(tc.content, "\n")
			result := checkBrokenLinks(lines)
			assert.Equal(t, tc.wantPassed, result.Passed, "result.Passed mismatch")
			if tc.wantDetail != "" {
				found := false
				for _, d := range result.Details {
					if strings.Contains(d, tc.wantDetail) {
						found = true
						break
					}
				}
				assert.True(t, found, "expected detail containing %q, got %v", tc.wantDetail, result.Details)
			}
		})
	}
}

// --- checkEntriesWithoutDates tests ---

func TestCheckEntriesWithoutDates(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantPassed bool
		wantDetail string
	}{
		{
			name: "all versions have dates passes",
			content: `## [1.0.0] - 2024-01-01
## [2.0.0] - 2024-02-01
`,
			wantPassed: true,
		},
		{
			name: "version without date fails",
			content: `## [1.0.0]
`,
			wantPassed: false,
			wantDetail: "1.0.0",
		},
		{
			name: "Unreleased without date passes",
			content: `## [Unreleased]
`,
			wantPassed: true,
		},
		{
			name: "invalid date format fails",
			content: `## [1.0.0] - 01-01-2024
`,
			wantPassed: false,
			wantDetail: "1.0.0",
		},
		{
			name: "valid date format passes",
			content: `## [1.0.0] - 2024-01-01
`,
			wantPassed: true,
		},
		{
			name: "date with YANKED marker passes",
			content: `## [1.0.0] - 2024-01-01 [YANKED]
`,
			wantPassed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines := strings.Split(tc.content, "\n")
			result := checkEntriesWithoutDates(lines)
			assert.Equal(t, tc.wantPassed, result.Passed, "result.Passed mismatch")
			if tc.wantDetail != "" {
				found := false
				for _, d := range result.Details {
					if strings.Contains(d, tc.wantDetail) {
						found = true
						break
					}
				}
				assert.True(t, found, "expected detail containing %q, got %v", tc.wantDetail, result.Details)
			}
		})
	}
}

// --- checkSemverOrder tests ---

func TestCheckSemverOrder(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantPassed bool
		wantDetail string
	}{
		{
			name: "descending order passes",
			content: `## [1.2.0] - 2024-03-01
## [1.1.0] - 2024-02-01
## [1.0.0] - 2024-01-01
`,
			wantPassed: true,
		},
		{
			name: "out of order fails",
			content: `## [1.0.0] - 2024-01-01
## [1.2.0] - 2024-03-01
## [1.1.0] - 2024-02-01
`,
			wantPassed: false,
			wantDetail: "1.0.0",
		},
		{
			name: "single version passes",
			content: `## [1.0.0] - 2024-01-01
`,
			wantPassed: true,
		},
		{
			name: "no versions only Unreleased passes",
			content: `## [Unreleased]
`,
			wantPassed: true,
		},
		{
			name: "v-prefixed versions pass",
			content: `## [v2.0.0] - 2024-02-01
## [v1.0.0] - 2024-01-01
`,
			wantPassed: true,
		},
		{
			name: "Unreleased before versions passes",
			content: `## [Unreleased]
## [2.0.0] - 2024-02-01
## [1.0.0] - 2024-01-01
`,
			wantPassed: true,
		},
		{
			name: "pre-release ordering passes",
			content: `## [2.0.0] - 2024-03-01
## [1.1.0-beta.1] - 2024-02-01
## [1.0.0] - 2024-01-01
`,
			wantPassed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines := strings.Split(tc.content, "\n")
			result := checkSemverOrder(lines)
			assert.Equal(t, tc.wantPassed, result.Passed, "result.Passed mismatch")
			if tc.wantDetail != "" {
				found := false
				for _, d := range result.Details {
					if strings.Contains(d, tc.wantDetail) {
						found = true
						break
					}
				}
				assert.True(t, found, "expected detail containing %q, got %v", tc.wantDetail, result.Details)
			}
		})
	}
}

// --- FormatReport tests ---

func TestFormatReport(t *testing.T) {
	t.Run("all pass report", func(t *testing.T) {
		report := &ValidationReport{
			File:       "CHANGELOG.md",
			Passed:     true,
			TotalRules: 5,
			PassCount:  5,
			FailCount:  0,
			Results: []ValidationResult{
				{Name: "Version headers", Passed: true, Message: "All version headers are properly formatted"},
				{Name: "Duplicate entries", Passed: true, Message: "No duplicate entries found"},
				{Name: "Broken links", Passed: true, Message: "All links are valid"},
				{Name: "Entries without dates", Passed: true, Message: "All versions have dates"},
				{Name: "Semver order", Passed: true, Message: "Versions are in correct descending order"},
			},
		}
		out := FormatReport(report)
		assert.Contains(t, out, "CHANGELOG.md")
		assert.Contains(t, out, "5/5")
		assert.Contains(t, out, "Version headers")
	})

	t.Run("mix of pass and fail", func(t *testing.T) {
		report := &ValidationReport{
			File:       "CHANGELOG.md",
			Passed:     false,
			TotalRules: 5,
			PassCount:  4,
			FailCount:  1,
			Results: []ValidationResult{
				{Name: "Version headers", Passed: true, Message: "All version headers are properly formatted"},
				{Name: "Broken links", Passed: false, Message: "2 issues found", Details: []string{"Version [1.0.0] has no matching reference link"}},
				{Name: "Duplicate entries", Passed: true, Message: "No duplicate entries found"},
				{Name: "Entries without dates", Passed: true, Message: "All versions have dates"},
				{Name: "Semver order", Passed: true, Message: "Versions are in correct descending order"},
			},
		}
		out := FormatReport(report)
		assert.Contains(t, out, "4/5")
		assert.Contains(t, out, "1 failed")
		assert.Contains(t, out, "1.0.0")
	})

	t.Run("details are indented under failures", func(t *testing.T) {
		report := &ValidationReport{
			File:       "CHANGELOG.md",
			Passed:     false,
			TotalRules: 1,
			PassCount:  0,
			FailCount:  1,
			Results: []ValidationResult{
				{Name: "Broken links", Passed: false, Message: "1 issue found", Details: []string{"Version [2.0.0] has no matching reference link"}},
			},
		}
		out := FormatReport(report)
		// Detail should appear indented (after the rule line)
		lines := strings.Split(out, "\n")
		detailFound := false
		for _, l := range lines {
			if strings.Contains(l, "2.0.0") {
				detailFound = true
				// Should be indented
				assert.True(t, strings.HasPrefix(l, "  ") || strings.HasPrefix(l, "\t"),
					"detail line should be indented, got: %q", l)
			}
		}
		assert.True(t, detailFound, "detail should appear in output")
	})
}
