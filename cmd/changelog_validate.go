// cmd/changelog_validate.go

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/peiman/changie/internal/changelog"
	"github.com/peiman/changie/internal/ui"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate changelog for common problems",
	Long: `Validate the changelog file for common problems.

Runs five checks:
  1. Version headers — must use ## [X.Y.Z] format with valid semver
  2. Duplicate entries — no identical bullet points within a section
  3. Broken links — every ## [X] header must have a matching [X]: URL reference
  4. Entries without dates — non-Unreleased versions must include YYYY-MM-DD
  5. Semver order — versions must appear in descending order (newest first)`,
	RunE: runValidateChangelog,
}

func init() {
	changelogCmd.AddCommand(validateCmd)
}

func runValidateChangelog(cmd *cobra.Command, _ []string) error {
	file := getConfigValueWithFlags[string](cmd, "file", "app.changelog.file")

	data, err := os.ReadFile(file) //nolint:gosec // G304: file is a user-specified changelog path // nosemgrep: go-path-traversal
	if err != nil {
		return fmt.Errorf("failed to read changelog: %w", err)
	}

	report := changelog.ValidateChangelog(string(data), file)
	if report.Passed {
		return ui.RenderSuccess(cmd.OutOrStdout(),
			fmt.Sprintf("All %d checks passed for %s", report.TotalRules, file), report)
	}

	formatted := changelog.FormatReport(report)
	_, _ = fmt.Fprint(cmd.OutOrStdout(), formatted)
	return fmt.Errorf("changelog validation failed: %d/%d checks failed",
		report.FailCount, report.TotalRules)
}
