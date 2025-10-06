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

	"github.com/peiman/changie/internal/output"
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
  # Basic major version bump
  changie bump major

  # With automatic push to remote
  changie bump major --auto-push

  # On a release branch
  changie bump major --allow-any-branch

  # With custom changelog file
  changie bump major --file HISTORY.md

  # JSON output for scripts/automation
  changie bump major --json

COMMON USE CASES:
  - Removing deprecated features
  - Changing API contracts
  - Major architectural changes
  - Database schema migrations requiring data changes
  - Any change that breaks backward compatibility`,
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
  # Basic minor version bump
  changie bump minor

  # With automatic push to remote
  changie bump minor --auto-push

  # On a release branch
  changie bump minor --allow-any-branch

  # With custom changelog file
  changie bump minor --file HISTORY.md

  # JSON output for scripts/automation
  changie bump minor --json

COMMON USE CASES:
  - Adding new API endpoints
  - New command-line options or commands
  - New functionality that doesn't break existing code
  - Performance improvements
  - New optional features`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersionBump(cmd, "minor")
		},
	}

	// patchCmd represents the command to bump the patch version number
	patchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Bump patch version (x.y.Z) for bug fixes",
		Long: `Release a patch version by bumping the third version number.

Use this command when making BUG FIXES in a backward-compatible manner.

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
  # Basic patch version bump
  changie bump patch

  # With automatic push to remote
  changie bump patch --auto-push

  # On a hotfix branch
  changie bump patch --allow-any-branch

  # With custom changelog file
  changie bump patch --file HISTORY.md

  # JSON output for scripts/automation
  changie bump patch --json

COMMON USE CASES:
  - Bug fixes
  - Security patches
  - Documentation corrections
  - Minor refactoring without behavior changes
  - Dependency updates (security fixes)`,
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
	result, err := version.Bump(cfg, cmd.OutOrStdout())

	// If JSON output is requested, output structured result
	if output.IsJSONEnabled() {
		jsonOutput := output.BumpOutput{
			Success:       err == nil,
			BumpType:      bumpType,
			ChangelogFile: file,
		}

		if err != nil {
			jsonOutput.Error = err.Error()
		} else if result != nil {
			jsonOutput.OldVersion = result.OldVersion
			jsonOutput.NewVersion = result.NewVersion
			jsonOutput.Tag = result.NewVersion
			jsonOutput.Pushed = result.Pushed
		}

		return output.WriteJSON(cmd.OutOrStdout(), jsonOutput)
	}

	return err
}
