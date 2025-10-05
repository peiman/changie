// Package cmd implements command line interface commands for the application.
//
// This file contains the implementation of the version-related commands:
// - major: Bump the major version number (X.y.z -> X+1.0.0)
// - minor: Bump the minor version number (x.Y.z -> x.Y+1.0)
// - patch: Bump the patch version number (x.y.Z -> x.y.Z+1)
//
// These commands manage semantic versioning operations including checking for
// uncommitted changes, changelog updates, and git tagging.
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/peiman/changie/internal/version"
)

var (
	// majorCmd represents the command to bump the major version number
	majorCmd = &cobra.Command{
		Use:   "major",
		Short: "Bump the major version number",
		Long: `Release a major version by bumping the first version number.

For example, 1.2.3 → 2.0.0

This command will:
1. Check that you're on main/master branch (use --allow-any-branch to bypass)
2. Check for uncommitted changes
3. Update the changelog
4. Commit the changes
5. Create a new git tag
6. Optionally push changes and tags to remote repository`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersionBump(cmd, "major")
		},
	}

	// minorCmd represents the command to bump the minor version number
	minorCmd = &cobra.Command{
		Use:   "minor",
		Short: "Bump the minor version number",
		Long: `Release a minor version by bumping the second version number.

For example, 1.2.3 → 1.3.0

This command will:
1. Check that you're on main/master branch (use --allow-any-branch to bypass)
2. Check for uncommitted changes
3. Update the changelog
4. Commit the changes
5. Create a new git tag
6. Optionally push changes and tags to remote repository`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersionBump(cmd, "minor")
		},
	}

	// patchCmd represents the command to bump the patch version number
	patchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Bump the patch version number",
		Long: `Release a patch version by bumping the third version number.

For example, 1.2.3 → 1.2.4

This command will:
1. Check that you're on main/master branch (use --allow-any-branch to bypass)
2. Check for uncommitted changes
3. Update the changelog
4. Commit the changes
5. Create a new git tag
6. Optionally push changes and tags to remote repository`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersionBump(cmd, "patch")
		},
	}
)

// init registers the version commands with the root command and
// defines and binds their flags to viper configuration values.
func init() {
	// Add common flags to all version commands
	for _, cmd := range []*cobra.Command{majorCmd, minorCmd, patchCmd} {
		cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
		cmd.Flags().String("rrp", "github", "Remote repository provider (github, bitbucket)")
		cmd.Flags().Bool("auto-push", false, "Automatically push changes and tags")
		cmd.Flags().Bool("allow-any-branch", false, "Allow version bumping on any branch (bypasses main/master branch check)")

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
		if err := viper.BindPFlag("app.version.allow_any_branch", cmd.Flags().Lookup("allow-any-branch")); err != nil {
			log.Fatal().Err(err).Msg("Failed to bind 'allow-any-branch' flag")
		}

		// Add command to RootCmd
		RootCmd.AddCommand(cmd)
	}
}

// runVersionBump is a thin wrapper that constructs configuration and delegates
// to the version.Bump function in the internal/version package.
//
// This function is responsible only for:
// 1. Reading configuration from viper and flags
// 2. Constructing the BumpConfig struct
// 3. Calling version.Bump with the appropriate output writer
//
// All business logic resides in internal/version/bump.go
//
// Parameters:
//   - cmd: The cobra command being executed
//   - bumpType: Type of version bump ("major", "minor", or "patch")
//
// Returns:
//   - error: Any error that occurred during execution
func runVersionBump(cmd *cobra.Command, bumpType string) error {
	// Get configuration values from viper/flags
	allowAnyBranch := viper.GetBool("app.version.allow_any_branch")
	if cmd.Flags().Changed("allow-any-branch") {
		allowAnyBranch, _ = cmd.Flags().GetBool("allow-any-branch")
	}

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

	useVPrefix := viper.GetBool("app.version.use_v_prefix")

	// Construct configuration
	cfg := version.BumpConfig{
		BumpType:           bumpType,
		AllowAnyBranch:     allowAnyBranch,
		AutoPush:           autoPush,
		ChangelogFile:      file,
		RepositoryProvider: repositoryProvider,
		UseVPrefix:         useVPrefix,
	}

	// Delegate to the version package
	return version.Bump(cfg, cmd.OutOrStdout())
}
