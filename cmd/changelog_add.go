// cmd/changelog_add.go

package cmd

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/peiman/changie/internal/changelog"
)

// createChangelogSectionCmd creates a command for adding entries to a specific changelog section.
func createChangelogSectionCmd(section string) *cobra.Command {
	return &cobra.Command{
		Use:   strings.ToLower(section) + " CONTENT",
		Short: "Add a " + strings.ToLower(section) + " entry to the changelog",
		Long: fmt.Sprintf(`Add an entry to the %s section in your changelog.

This adds a bullet point to the %s section in the [Unreleased] area of your changelog file.
The entry will be formatted as a bullet point according to Keep a Changelog format.`, strings.ToLower(section), section),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddChangelogSection(cmd, args, section)
		},
	}
}

func init() {
	// Create a subcommand for each section type
	for section := range changelog.ValidSections {
		sectionCmd := createChangelogSectionCmd(section)
		changelogCmd.AddCommand(sectionCmd)
	}
}

func runAddChangelogSection(cmd *cobra.Command, args []string, section string) error {
	log.Debug().Str("section", section).Msg("Adding changelog entry")

	file := viper.GetString("app.changelog.file")
	if cmd.Flags().Changed("file") {
		file, _ = cmd.Flags().GetString("file")
	}

	content := args[0]
	log.Info().
		Str("file", file).
		Str("section", section).
		Str("content", content).
		Msg("Adding changelog entry")

	isDuplicate, err := changelog.AddChangelogSection(file, section, content)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", file).
			Str("section", section).
			Str("content", content).
			Msg("Failed to add changelog entry")
		return fmt.Errorf("failed to add changelog entry: %w", err)
	}

	if isDuplicate {
		fmt.Fprintf(cmd.OutOrStdout(), "Entry already exists in the %s section, not added again.\n", section)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Added to %s section: %s\n", section, content)
	}

	log.Debug().Msg("Changelog entry added successfully")
	return nil
}
