// cmd/init.go

package cmd

import (
	"fmt"
	"strings"

	"github.com/peiman/changie/internal/changelog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project with SemVer and Keep a Changelog",
	Long: `Sets up the project directory for Semantic Versioning and Keep a Changelog format.

This command creates a new CHANGELOG.md file in the current directory following the
Keep a Changelog format (https://keepachangelog.com/en/1.1.0/).`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")

	// Bind flags to Viper
	if err := viper.BindPFlag("app.changelog.file", initCmd.Flags().Lookup("file")); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind 'file' flag")
	}

	// Add initCmd to RootCmd
	RootCmd.AddCommand(initCmd)
}

func initInitConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("app.changelog.file", "CHANGELOG.md")
}

func runInit(cmd *cobra.Command, args []string) error {
	log.Debug().Msg("Starting runInit execution")
	initInitConfig()

	file := viper.GetString("app.changelog.file")
	if cmd.Flags().Changed("file") {
		file, _ = cmd.Flags().GetString("file")
	}

	log.Info().Str("file", file).Msg("Initializing project with changelog file")

	if err := changelog.InitProject(file); err != nil {
		log.Error().Err(err).Str("file", file).Msg("Failed to initialize project")
		return fmt.Errorf("failed to initialize project: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Project initialized with changelog file: %s\n", file)
	log.Debug().Msg("runInit completed successfully")

	return nil
}
