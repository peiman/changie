// cmd/changelog.go

package cmd

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Manage changelog entries",
	Long: `Manage changelog entries in Keep a Changelog format.

This command group allows you to add entries to different sections of your
CHANGELOG.md file, following the Keep a Changelog format.`,
}

func init() {
	changelogCmd.PersistentFlags().String("file", "CHANGELOG.md", "Changelog file name")

	// Bind flags to Viper
	if err := viper.BindPFlag("app.changelog.file", changelogCmd.PersistentFlags().Lookup("file")); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind 'file' flag")
	}

	// Add changelogCmd to RootCmd
	RootCmd.AddCommand(changelogCmd)
}

func initChangelogConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("app.changelog.file", "CHANGELOG.md")
}
