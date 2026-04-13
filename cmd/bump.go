// Package cmd implements command line interface commands for the application.
//
// This file contains the implementation of the bump command and its subcommands:
// - bump major: Bump the major version number (X.y.z -> X+1.0.0)
// - bump minor: Bump the minor version number (x.Y.z -> x.Y+1.0)
// - bump patch: Bump the patch version number (x.y.Z -> x.y.Z+1)
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
	// bumpCmd represents the parent bump command
	bumpCmd = &cobra.Command{
		Use:   "bump",
		Short: "Bump version numbers following semantic versioning",
		Long: `Bump version numbers following semantic versioning (SemVer) principles.

Use subcommands to specify the type of version bump:
  - major: Breaking changes (X.y.z -> X+1.0.0)
  - minor: New features, backward compatible (x.Y.z -> x.Y+1.0)
  - patch: Bug fixes, backward compatible (x.y.Z -> x.y.Z+1)

Each bump command will:
1. Check that you're on main/master branch (use --allow-any-branch to bypass)
2. Check for uncommitted changes
3. Update the changelog
4. Commit the changes
5. Create a new git tag
6. Optionally push changes and tags to remote repository`,
	}

	// majorCmd represents the command to bump the major version number
	majorCmd = &cobra.Command{
		Use:   "major",
		Short: "Bump major version (X.0.0) for breaking changes",
		Long: `Release a major version by bumping the first version number.

Use this command when making BREAKING CHANGES or incompatible API changes.

Version change: 1.2.3 → 2.0.0

WHAT IT DOES:
1. Validates you're on main/master branch (bypass with --allow-any-branch)
2. Checks for uncommitted changes (must have clean working directory)
3. Gets current version from git tags
4. Bumps major version and resets minor/patch to 0
5. Updates CHANGELOG.md (moves Unreleased → new version section)
6. Commits changelog with message "Release vX.0.0"
7. Creates git tag (e.g., v2.0.0)
8. Optionally pushes to remote (with --auto-push)

EXAMPLES:
  changie bump major
  changie bump major --auto-push
  changie bump major --allow-any-branch
  changie bump major --output json

COMMON USE CASES:
  - Removing or renaming public API endpoints
  - Changing configuration file format
  - Dropping support for a platform or runtime
  - Any change that requires users to update their code`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersionBump(cmd, "major")
		},
	}

	// minorCmd represents the command to bump the minor version number
	minorCmd = &cobra.Command{
		Use:   "minor",
		Short: "Bump minor version (x.Y.0) for new features",
		Long: `Release a minor version by bumping the second version number.

Use this command when adding NEW FEATURES in a backward-compatible manner.

Version change: 1.2.3 → 1.3.0

WHAT IT DOES:
1. Validates you're on main/master branch (bypass with --allow-any-branch)
2. Checks for uncommitted changes (must have clean working directory)
3. Gets current version from git tags
4. Bumps minor version and resets patch to 0
5. Updates CHANGELOG.md (moves Unreleased → new version section)
6. Commits changelog with message "Release vX.Y.0"
7. Creates git tag (e.g., v1.3.0)
8. Optionally pushes to remote (with --auto-push)

EXAMPLES:
  changie bump minor
  changie bump minor --auto-push
  changie bump minor --allow-any-branch
  changie bump minor --output json

COMMON USE CASES:
  - Adding new API endpoints
  - New command-line options or commands
  - New functionality that doesn't break existing code
  - Performance improvements`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersionBump(cmd, "minor")
		},
	}

	// patchCmd represents the command to bump the patch version number
	patchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Bump patch version (x.y.Z) for bug fixes",
		Long: `Release a patch version by bumping the third version number.

Use this command for BUG FIXES and backward-compatible corrections.

Version change: 1.2.3 → 1.2.4

WHAT IT DOES:
1. Validates you're on main/master branch (bypass with --allow-any-branch)
2. Checks for uncommitted changes (must have clean working directory)
3. Gets current version from git tags
4. Bumps patch version
5. Updates CHANGELOG.md (moves Unreleased → new version section)
6. Commits changelog with message "Release vX.Y.Z"
7. Creates git tag (e.g., v1.2.4)
8. Optionally pushes to remote (with --auto-push)

EXAMPLES:
  changie bump patch
  changie bump patch --auto-push
  changie bump patch --allow-any-branch
  changie bump patch --output json

COMMON USE CASES:
  - Fixing bugs without adding new features
  - Correcting typos in output
  - Updating dependencies for security patches
  - Minor documentation fixes that affect behavior`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersionBump(cmd, "patch")
		},
	}
)

// init registers the bump command with the root command and
// defines and binds flags to viper configuration values.
func init() {
	// Add common flags to all bump subcommands
	for _, cmd := range []*cobra.Command{majorCmd, minorCmd, patchCmd} {
		cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
		cmd.Flags().String("rrp", "github", "Remote repository provider (github, bitbucket)")
		cmd.Flags().Bool("auto-push", false, "Automatically push changes and tags")
		cmd.Flags().Bool("allow-any-branch", false, "Allow version bumping on any branch (bypasses main/master branch check)")
		cmd.Flags().Bool("use-v-prefix", true, "Use 'v' prefix for version tags (e.g., v1.0.0)")

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

		// Add as subcommand of bump
		bumpCmd.AddCommand(cmd)
	}

	// Add bump command to root
	RootCmd.AddCommand(bumpCmd)
}

func runVersionBump(cmd *cobra.Command, bumpType string) error {
	cfg := version.BumpConfig{
		BumpType:           bumpType,
		AllowAnyBranch:     getConfigValueWithFlags[bool](cmd, "allow-any-branch", "app.version.allow_any_branch"),
		AutoPush:           getConfigValueWithFlags[bool](cmd, "auto-push", "app.changelog.auto_push"),
		ChangelogFile:      getConfigValueWithFlags[string](cmd, "file", "app.changelog.file"),
		RepositoryProvider: getConfigValueWithFlags[string](cmd, "rrp", "app.changelog.repository_provider"),
		UseVPrefix:         getConfigValueWithFlags[bool](cmd, "use-v-prefix", "app.version.use_v_prefix"),
	}
	return version.Bump(cfg, cmd.OutOrStdout())
}
