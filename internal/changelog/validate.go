// Package changelog - validate.go provides validation logic for CHANGELOG.md files.
// It checks for common problems and outputs a structured pass/fail report.

package changelog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/peiman/changie/internal/semver"
)

// ValidationResult represents the outcome of a single validation check.
type ValidationResult struct {
	Name    string   `json:"name"`
	Passed  bool     `json:"passed"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
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

// Regexes compiled once for performance.
var (
	// Matches ## [Something] (with optional date after)
	reVersionHeader = regexp.MustCompile(`(?m)^## \[([^\]]+)\]`)
	// Matches [Something]: URL (reference-style links at bottom)
	reLinkRef = regexp.MustCompile(`(?m)^\[([^\]]+)\]:\s+\S+`)
	// Matches ## [X] with a valid date (any text after date is allowed, e.g. [YANKED])
	reVersionWithDate = regexp.MustCompile(`^## \[([^\]]+)\] - \d{4}-\d{2}-\d{2}`)
)

// ValidateChangelog runs all five validation checks against changelog content.
// The filePath parameter is used only for report metadata (no I/O performed).
func ValidateChangelog(content string, filePath string) *ValidationReport {
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")

	results := []ValidationResult{
		checkVersionHeaders(lines),
		checkDuplicateEntries(lines),
		checkBrokenLinks(lines),
		checkEntriesWithoutDates(lines),
		checkSemverOrder(lines),
	}

	passCount := 0
	failCount := 0
	for _, r := range results {
		if r.Passed {
			passCount++
		} else {
			failCount++
		}
	}

	return &ValidationReport{
		File:       filePath,
		Passed:     failCount == 0,
		TotalRules: len(results),
		PassCount:  passCount,
		FailCount:  failCount,
		Results:    results,
	}
}

// checkVersionHeaders validates that all version-like headers use valid semver.
// Rule: ## [X.Y.Z] - YYYY-MM-DD format; [Unreleased] is always valid.
func checkVersionHeaders(lines []string) ValidationResult {
	name := "Version headers"
	var details []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "## ") {
			continue
		}
		// Must match ## [Something]
		m := reVersionHeader.FindStringSubmatch(line)
		if m == nil {
			// It's a ## header but has no brackets — flag it
			// Strip "## " and use the rest as the problematic token
			token := strings.TrimPrefix(line, "## ")
			details = append(details, fmt.Sprintf("Malformed version header (missing brackets): %q", token))
			continue
		}
		versionStr := m[1]
		if versionStr == "Unreleased" {
			continue
		}
		// Validate the version portion with semver
		_, _, err := semver.ParseVersion(versionStr)
		if err != nil {
			details = append(details, fmt.Sprintf("Invalid semver in header: [%s]", versionStr))
		}
	}

	if len(details) == 0 {
		return ValidationResult{
			Name:    name,
			Passed:  true,
			Message: "All version headers are properly formatted",
		}
	}
	return ValidationResult{
		Name:    name,
		Passed:  false,
		Message: fmt.Sprintf("%d malformed version header(s) found", len(details)),
		Details: details,
	}
}

// checkDuplicateEntries finds duplicate bullet points within the same version+section.
func checkDuplicateEntries(lines []string) ValidationResult {
	name := "Duplicate entries"
	var details []string

	// Track current context
	currentVersion := ""
	currentSection := ""
	// Map of version+section -> set of entries
	seen := map[string]map[string]bool{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect version header
		if strings.HasPrefix(trimmed, "## ") {
			m := reVersionHeader.FindStringSubmatch(trimmed)
			if m != nil {
				currentVersion = m[1]
				currentSection = ""
			}
			continue
		}

		// Detect subsection header
		if strings.HasPrefix(trimmed, "### ") {
			currentSection = strings.TrimPrefix(trimmed, "### ")
			continue
		}

		// Detect bullet entry
		if strings.HasPrefix(trimmed, "- ") && currentVersion != "" {
			entry := strings.TrimPrefix(trimmed, "- ")
			key := currentVersion + "\x00" + currentSection
			if seen[key] == nil {
				seen[key] = map[string]bool{}
			}
			if seen[key][entry] {
				details = append(details, fmt.Sprintf("Duplicate entry in [%s] %s: %q", currentVersion, currentSection, entry))
			}
			seen[key][entry] = true
		}
	}

	if len(details) == 0 {
		return ValidationResult{
			Name:    name,
			Passed:  true,
			Message: "No duplicate entries found",
		}
	}
	return ValidationResult{
		Name:    name,
		Passed:  false,
		Message: fmt.Sprintf("%d duplicate entry/entries found", len(details)),
		Details: details,
	}
}

// checkBrokenLinks verifies every ## [X] header has a matching [X]: URL reference, and vice versa.
func checkBrokenLinks(lines []string) ValidationResult {
	name := "Broken links"
	var details []string

	// Collect all version header labels
	headersSet := map[string]bool{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		m := reVersionHeader.FindStringSubmatch(line)
		if m != nil {
			headersSet[m[1]] = true
		}
	}

	// Collect all link reference labels
	linksSet := map[string]bool{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		m := reLinkRef.FindStringSubmatch(line)
		if m != nil {
			linksSet[m[1]] = true
		}
	}

	// If no headers at all, nothing to check
	if len(headersSet) == 0 {
		return ValidationResult{Name: name, Passed: true, Message: "No version headers found to validate"}
	}

	// Find headers missing links
	for v := range headersSet {
		if !linksSet[v] {
			details = append(details, fmt.Sprintf("Version [%s] has no matching reference link", v))
		}
	}

	// Find orphan links (link with no matching header)
	for v := range linksSet {
		if !headersSet[v] {
			details = append(details, fmt.Sprintf("Orphan link [%s] has no matching version header", v))
		}
	}

	if len(details) == 0 {
		return ValidationResult{
			Name:    name,
			Passed:  true,
			Message: "All links are valid",
		}
	}
	return ValidationResult{
		Name:    name,
		Passed:  false,
		Message: fmt.Sprintf("%d link issue(s) found", len(details)),
		Details: details,
	}
}

// checkEntriesWithoutDates ensures all non-Unreleased version headers include a YYYY-MM-DD date.
func checkEntriesWithoutDates(lines []string) ValidationResult {
	name := "Entries without dates"
	var details []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "## ") {
			continue
		}
		m := reVersionHeader.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		versionStr := m[1]
		if versionStr == "Unreleased" {
			continue
		}
		// Check the full line for a valid date
		if !reVersionWithDate.MatchString(line) {
			details = append(details, fmt.Sprintf("Version [%s] is missing a valid YYYY-MM-DD date", versionStr))
		}
	}

	if len(details) == 0 {
		return ValidationResult{
			Name:    name,
			Passed:  true,
			Message: "All versions have dates",
		}
	}
	return ValidationResult{
		Name:    name,
		Passed:  false,
		Message: fmt.Sprintf("%d version(s) missing dates", len(details)),
		Details: details,
	}
}

// checkSemverOrder verifies versions appear in strictly descending semver order (top = newest).
func checkSemverOrder(lines []string) ValidationResult {
	name := "Semver order"

	// Collect version strings in document order (skipping Unreleased)
	var versions []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		m := reVersionHeader.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		v := m[1]
		if v == "Unreleased" {
			continue
		}
		versions = append(versions, v)
	}

	if len(versions) <= 1 {
		return ValidationResult{
			Name:    name,
			Passed:  true,
			Message: "Versions are in correct descending order",
		}
	}

	var details []string
	for i := 0; i < len(versions)-1; i++ {
		cmp, err := semver.Compare(versions[i], versions[i+1])
		if err != nil {
			// If we can't parse, skip ordering check for this pair
			continue
		}
		if cmp <= 0 {
			details = append(details, fmt.Sprintf(
				"Version [%s] should be greater than [%s] (expected descending order)",
				versions[i], versions[i+1],
			))
		}
	}

	if len(details) == 0 {
		return ValidationResult{
			Name:    name,
			Passed:  true,
			Message: "Versions are in correct descending order",
		}
	}
	return ValidationResult{
		Name:    name,
		Passed:  false,
		Message: fmt.Sprintf("%d ordering violation(s) found", len(details)),
		Details: details,
	}
}

// FormatReport renders a ValidationReport as human-readable text.
func FormatReport(report *ValidationReport) string {
	var sb strings.Builder

	title := fmt.Sprintf("Changelog Validation: %s", report.File)
	sb.WriteString(title + "\n")
	sb.WriteString(strings.Repeat("=", len(title)) + "\n\n")

	for _, r := range report.Results {
		icon := "✅"
		if !r.Passed {
			icon = "❌"
		}
		sb.WriteString(fmt.Sprintf("  %s %-30s %s\n", icon, r.Name, r.Message))
		for _, d := range r.Details {
			sb.WriteString(fmt.Sprintf("     • %s\n", d))
		}
	}

	sb.WriteString("\n")
	if report.Passed {
		sb.WriteString(fmt.Sprintf("Result: %d/%d checks passed\n", report.PassCount, report.TotalRules))
	} else {
		sb.WriteString(fmt.Sprintf("Result: %d/%d checks passed, %d failed\n",
			report.PassCount, report.TotalRules, report.FailCount))
	}

	return sb.String()
}
