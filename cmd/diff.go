// cmd/diff.go — wiring only; business logic lives in internal/changelog.

package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/peiman/changie/internal/changelog"
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
	data, err := os.ReadFile(file) //nolint:gosec // G304: file is a user-specified changelog path // nosemgrep: go-path-traversal
	if err != nil {
		return fmt.Errorf("failed to read changelog: %w", err)
	}
	result, err := changelog.DiffVersions(string(data), args[0], args[1])
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), result)
	return nil
}
