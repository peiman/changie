// cmd/version.go

package cmd

import (
	"fmt"
	"strings"

	"github.com/peiman/changie/internal/changelog"
	"github.com/peiman/changie/internal/git"
	"github.com/peiman/changie/internal/semver"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	majorCmd = &cobra.Command{
		Use:   "major",
		Short: "Bump the major version number",
		Long: `Release a major version by bumping the first version number.

For example, 1.2.3 → 2.0.0

This command will:
1. Check for uncommitted changes
2. Update the changelog
3. Commit the changes
4. Create a new git tag
5. Optionally push changes and tags to remote repository`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersionBump(cmd, "major")
		},
	}

	minorCmd = &cobra.Command{
		Use:   "minor",
		Short: "Bump the minor version number",
		Long: `Release a minor version by bumping the second version number.

For example, 1.2.3 → 1.3.0

This command will:
1. Check for uncommitted changes
2. Update the changelog
3. Commit the changes
4. Create a new git tag
5. Optionally push changes and tags to remote repository`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersionBump(cmd, "minor")
		},
	}

	patchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Bump the patch version number",
		Long: `Release a patch version by bumping the third version number.

For example, 1.2.3 → 1.2.4

This command will:
1. Check for uncommitted changes
2. Update the changelog
3. Commit the changes
4. Create a new git tag
5. Optionally push changes and tags to remote repository`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersionBump(cmd, "patch")
		},
	}
)

func init() {
	// Add common flags to all version commands
	for _, cmd := range []*cobra.Command{majorCmd, minorCmd, patchCmd} {
		cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
		cmd.Flags().String("rrp", "github", "Remote repository provider (github, bitbucket)")
		cmd.Flags().Bool("auto-push", false, "Automatically push changes and tags")

		// Bind flags to Viper
		if err := viper.BindPFlag("app.changelog.file", cmd.Flags().Lookup("file")); err != nil {
			log.Fatal().Err(err).Msg("Failed to bind 'file' flag")
		}
		if err := viper.BindPFlag("app.changelog.repository_provider", cmd.Flags().Lookup("rrp")); err != nil {
			log.Fatal().Err(err).Msg("Failed to bind 'rrp' flag")
		}
		if err := viper.BindPFlag("app.changelog.auto_push", cmd.Flags().Lookup("auto-push")); err != nil {
			log.Fatal().Err(err).Msg("Failed to bind 'auto-push' flag")
		}

		// Add command to RootCmd
		RootCmd.AddCommand(cmd)
	}
}

func initVersionConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("app.changelog.file", "CHANGELOG.md")
	viper.SetDefault("app.changelog.repository_provider", "github")
	viper.SetDefault("app.changelog.auto_push", false)
}

func runVersionBump(cmd *cobra.Command, bumpType string) error {
	log.Debug().Str("type", bumpType).Msg("Starting version bump")
	initVersionConfig()

	// Check if git is installed
	if !git.IsInstalled() {
		err := fmt.Errorf("git is not installed or not available in PATH")
		log.Error().Err(err).Msg("Failed to run git")
		return err
	}

	// Check for uncommitted changes
	hasUncommittedChanges, err := git.HasUncommittedChanges()
	if err != nil {
		log.Error().Err(err).Msg("Failed to check for uncommitted changes")
		return fmt.Errorf("failed to check for uncommitted changes: %w", err)
	}

	if hasUncommittedChanges {
		err := fmt.Errorf("uncommitted changes found, please commit or stash your changes before bumping version")
		log.Error().Err(err).Msg("Failed to bump version")
		return err
	}

	// Get current version from git
	currentVersion, err := git.GetVersion()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current version from git")
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// Log current version
	if currentVersion == "" {
		currentVersion = "0.0.0" // Default if no tag exists
		fmt.Fprintf(cmd.OutOrStdout(), "No version tag found, starting from %s\n", currentVersion)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Current version: %s\n", currentVersion)
	}

	// Bump version according to type
	var newVersion string
	switch bumpType {
	case "major":
		newVersion, err = semver.BumpMajor(currentVersion)
	case "minor":
		newVersion, err = semver.BumpMinor(currentVersion)
	case "patch":
		newVersion, err = semver.BumpPatch(currentVersion)
	default:
		err = fmt.Errorf("invalid bump type: %s", bumpType)
	}

	if err != nil {
		log.Error().Err(err).Str("type", bumpType).Str("current_version", currentVersion).Msg("Failed to bump version")
		return fmt.Errorf("failed to bump version: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "New version: %s\n", newVersion)

	// Get configuration values
	file := viper.GetString("app.changelog.file")
	if cmd.Flags().Changed("file") {
		file, _ = cmd.Flags().GetString("file")
	}

	repositoryProvider := viper.GetString("app.changelog.repository_provider")
	if cmd.Flags().Changed("rrp") {
		repositoryProvider, _ = cmd.Flags().GetString("rrp")
	}

	autoPush := viper.GetBool("app.changelog.auto_push")
	if cmd.Flags().Changed("auto-push") {
		autoPush, _ = cmd.Flags().GetBool("auto-push")
	}

	// Update changelog
	fmt.Fprintf(cmd.OutOrStdout(), "Updating changelog file: %s\n", file)
	err = changelog.UpdateChangelog(file, newVersion, repositoryProvider)
	if err != nil {
		log.Error().Err(err).Str("file", file).Str("version", newVersion).Msg("Failed to update changelog")
		return fmt.Errorf("failed to update changelog: %w", err)
	}

	// Commit changes
	err = git.CommitChangelog(file, newVersion)
	if err != nil {
		log.Error().Err(err).Str("file", file).Str("version", newVersion).Msg("Failed to commit changelog")
		return fmt.Errorf("failed to commit changelog: %w", err)
	}

	// Tag version
	fmt.Fprintf(cmd.OutOrStdout(), "Tagging version: %s\n", newVersion)
	err = git.TagVersion(newVersion)
	if err != nil {
		log.Error().Err(err).Str("version", newVersion).Msg("Failed to tag version")
		return fmt.Errorf("failed to tag version: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s release %s done.\n", bumpType, newVersion)

	// Auto-push if enabled
	if autoPush {
		fmt.Fprintf(cmd.OutOrStdout(), "Pushing changes and tags...\n")
		err = git.PushChanges()
		if err != nil {
			log.Error().Err(err).Msg("Failed to push changes")
			return fmt.Errorf("failed to push changes: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Automatically pushed changes and tags to remote repository.\n")
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Don't forget to git push and git push --tags.\n")
	}

	log.Debug().Str("type", bumpType).Str("version", newVersion).Msg("Version bump completed successfully")
	return nil
}
