// cmd/init.go

package cmd

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/peiman/changie/internal/changelog"
	"github.com/peiman/changie/internal/git"
	"github.com/peiman/changie/internal/ui"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project with SemVer and Keep a Changelog",
	Long: `Sets up the project directory for Semantic Versioning and Keep a Changelog format.

This command creates a new CHANGELOG.md file in the current directory following the
Keep a Changelog format (https://keepachangelog.com/en/1.1.0/).

It also detects or configures whether to use a 'v' prefix for version tags (e.g., v1.0.0 vs 1.0.0).
If git tags already exist, it will adopt their convention. Otherwise, it will ask for your preference.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
	initCmd.Flags().Bool("use-v-prefix", true, "Use 'v' prefix for version tags (e.g., v1.0.0)")

	// Bind flags to Viper
	if err := viper.BindPFlag("app.changelog.file", initCmd.Flags().Lookup("file")); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind 'file' flag")
	}
	if err := viper.BindPFlag("app.version.use_v_prefix", initCmd.Flags().Lookup("use-v-prefix")); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind 'use-v-prefix' flag")
	}

	// Add initCmd to RootCmd
	RootCmd.AddCommand(initCmd)

	// Setup command configuration inheritance
	setupCommandConfig(initCmd)
}

func runInit(cmd *cobra.Command, _ []string) error {
	log.Debug().Msg("Starting runInit execution")

	file := viper.GetString("app.changelog.file")
	if cmd.Flags().Changed("file") {
		file, _ = cmd.Flags().GetString("file")
	}

	// Determine if we should use 'v' prefix for versions
	useVPrefix := viper.GetBool("app.version.use_v_prefix")
	explicitPrefixSet := cmd.Flags().Changed("use-v-prefix")

	// Check if git is installed
	if git.IsInstalled() {
		// Get current version from git if available
		currentVersion, err := git.GetVersion()
		if err == nil && currentVersion != "" {
			// Existing tags found, detect if they use 'v' prefix
			hasPrefix := strings.HasPrefix(currentVersion, "v")

			if explicitPrefixSet {
				// User explicitly set the preference, inform them if it differs from existing tags
				if hasPrefix != useVPrefix {
					fmt.Fprintf(cmd.OutOrStdout(), "Note: Your specified version prefix setting (%v) differs from existing tags (%v).\n",
						useVPrefix, hasPrefix)
					fmt.Fprintf(cmd.OutOrStdout(), "Using your specified preference: %v\n", useVPrefix)
				}
			} else {
				// No explicit preference, adopt existing convention
				useVPrefix = hasPrefix
				withoutText := ""
				if !hasPrefix {
					withoutText = "out"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Detected existing version tags with%s 'v' prefix. Using this convention.\n",
					withoutText)
			}
		} else if currentVersion == "" && !explicitPrefixSet {
			// No tags found and no explicit preference, ask user
			fmt.Fprintf(cmd.OutOrStdout(), "No existing version tags found.\n")
			result, err := ui.AskYesNo("Would you like to use 'v' prefix for version tags? (e.g., v1.0.0 vs 1.0.0)", true, cmd.OutOrStdout())
			if err != nil {
				log.Warn().Err(err).Msg("Error reading input, defaulting to use 'v' prefix")
			}
			useVPrefix = result
		}
	}

	// Save the version prefix preference to viper
	viper.Set("app.version.use_v_prefix", useVPrefix)

	// Log the decision
	log.Info().Bool("use_v_prefix", useVPrefix).Msg("Version prefix configuration set")
	prefixText := ""
	if !useVPrefix {
		prefixText = "not "
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Version tags will %suse 'v' prefix.\n", prefixText)

	// Initialize the changelog
	log.Info().Str("file", file).Msg("Initializing project with changelog file")

	if err := changelog.InitProject(file); err != nil {
		log.Error().Err(err).Str("file", file).Msg("Failed to initialize project")
		return fmt.Errorf("failed to initialize project: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Project initialized with changelog file: %s\n", file)

	// If git is available and there are no tags, create initial commit and tag
	if git.IsInstalled() {
		currentVersion, err := git.GetVersion()
		if err == nil && currentVersion == "" {
			// No tags exist, create initial commit and tag
			log.Info().Msg("No version tags found, creating initial commit and tag")

			// Add changelog file to git
			if err := git.CommitChangelog(file, "0.0.0"); err != nil {
				log.Warn().Err(err).Msg("Failed to create initial commit, continuing anyway")
			} else {
				// Create initial tag
				initialTag := "0.0.0"
				if useVPrefix {
					initialTag = "v0.0.0"
				}

				if err := git.TagVersion(initialTag); err != nil {
					log.Warn().Err(err).Str("tag", initialTag).Msg("Failed to create initial tag, continuing anyway")
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Created initial git tag: %s\n", initialTag)
					log.Info().Str("tag", initialTag).Msg("Created initial git tag")
				}
			}
		}
	}

	log.Debug().Msg("runInit completed successfully")

	return nil
}
