// cmd/init.go

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/peiman/changie/internal/changelog"
	"github.com/peiman/changie/internal/git"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// askForVersionPrefix prompts the user for their version prefix preference
func askForVersionPrefix(cmd *cobra.Command) bool {
	fmt.Fprintf(cmd.OutOrStdout(), "No existing version tags found.\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Would you like to use 'v' prefix for version tags? (e.g., v1.0.0 vs 1.0.0) [Y/n]: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Warn().Err(err).Msg("Error reading input, defaulting to use 'v' prefix")
		return true
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "" || input == "y" || input == "yes"
}

func runInit(cmd *cobra.Command, args []string) error {
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
			useVPrefix = askForVersionPrefix(cmd)
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
	log.Debug().Msg("runInit completed successfully")

	return nil
}
