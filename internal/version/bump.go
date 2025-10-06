// Package version provides version bump orchestration functionality.
//
// This package handles the workflow for bumping semantic versions,
// including validation, changelog updates, git operations, and pushing changes.
package version

import (
	"fmt"
	"io"

	"github.com/peiman/changie/internal/changelog"
	"github.com/peiman/changie/internal/git"
	"github.com/peiman/changie/internal/logger"
	"github.com/peiman/changie/internal/semver"
)

// BumpConfig holds all configuration needed for a version bump operation.
type BumpConfig struct {
	// BumpType specifies the type of bump: "major", "minor", or "patch"
	BumpType string

	// AllowAnyBranch bypasses the main/master branch check when true
	AllowAnyBranch bool

	// AutoPush automatically pushes changes and tags to remote when true
	AutoPush bool

	// ChangelogFile is the path to the changelog file
	ChangelogFile string

	// RepositoryProvider specifies the git provider (github, gitlab, bitbucket)
	RepositoryProvider string

	// UseVPrefix determines whether to add 'v' prefix to version tags
	UseVPrefix bool
}

// BumpResult holds the result of a version bump operation.
type BumpResult struct {
	// OldVersion is the version before bumping
	OldVersion string

	// NewVersion is the version after bumping
	NewVersion string

	// BumpType is the type of bump performed ("major", "minor", "patch")
	BumpType string

	// ChangelogFile is the path to the changelog file that was updated
	ChangelogFile string

	// Pushed indicates whether changes were automatically pushed to remote
	Pushed bool
}

// Bump performs a complete version bump workflow.
//
// The workflow includes:
// 1. Verifying git is installed and repository state
// 2. Optionally checking branch name (main/master)
// 3. Checking for uncommitted changes
// 4. Getting current version from git tags
// 5. Calculating new version based on bump type
// 6. Updating the changelog file
// 7. Committing changes and creating git tag
// 8. Optionally pushing changes to remote
//
// Parameters:
//   - cfg: Configuration for the bump operation
//   - output: Writer for user-facing output messages
//
// Returns:
//   - *BumpResult: Result of the bump operation (nil if error occurred)
//   - error: Any error encountered during the workflow
func Bump(cfg BumpConfig, output io.Writer) (*BumpResult, error) {
	logger.Version.Debug().Str("type", cfg.BumpType).Msg("Starting version bump")

	// Check if git is installed
	if !git.IsInstalled() {
		err := fmt.Errorf("git is not installed or not available in PATH - please install Git (https://git-scm.com/downloads) and ensure it's in your system PATH")
		logger.Version.Error().Err(err).Msg("Failed to run git")
		return nil, err
	}

	// Check if we're on main/master branch (unless bypassed)
	if !cfg.AllowAnyBranch {
		currentBranch, err := git.GetCurrentBranch()
		if err != nil {
			logger.Version.Error().Err(err).Msg("Failed to get current branch")
			return nil, fmt.Errorf("failed to get current branch: %w", err)
		}

		if currentBranch != "main" && currentBranch != "master" {
			err := fmt.Errorf("not on main/master branch (current: %s) - version bumps should typically be done on the main branch to maintain a clean release history. Use --allow-any-branch to bypass this check if you're working with release branches or have a different workflow", currentBranch)
			logger.Version.Error().Err(err).Str("branch", currentBranch).Msg("Branch check failed")
			return nil, err
		}
		logger.Version.Debug().Str("branch", currentBranch).Msg("Branch check passed")
	} else {
		logger.Version.Debug().Msg("Branch check bypassed with --allow-any-branch flag")
	}

	// Check for uncommitted changes
	hasUncommittedChanges, err := git.HasUncommittedChanges()
	if err != nil {
		logger.Version.Error().Err(err).Msg("Failed to check for uncommitted changes")
		return nil, fmt.Errorf("failed to check for uncommitted changes: %w", err)
	}

	if hasUncommittedChanges {
		err := fmt.Errorf("uncommitted changes found - run 'git status' to see changed files, then either commit changes with 'git commit' or stash them with 'git stash' before bumping version")
		logger.Version.Error().Err(err).Msg("Failed to bump version")
		return nil, err
	}

	// Get current version from git
	currentVersion, err := git.GetVersion()
	if err != nil {
		logger.Version.Error().Err(err).Msg("Failed to get current version from git")
		return nil, fmt.Errorf("failed to get current version: %w - ensure you're in a git repository with at least one tag, or initialize with 'git tag v0.0.0'", err)
	}

	// Log current version
	if currentVersion == "" {
		currentVersion = "0.0.0" // Default if no tag exists
		fmt.Fprintf(output, "No version tag found, starting from %s\n", currentVersion)
	} else {
		fmt.Fprintf(output, "Current version: %s\n", currentVersion)
	}

	// Bump version according to type
	var newVersion string
	switch cfg.BumpType {
	case "major":
		newVersion, err = semver.BumpMajor(currentVersion, cfg.UseVPrefix)
	case "minor":
		newVersion, err = semver.BumpMinor(currentVersion, cfg.UseVPrefix)
	case "patch":
		newVersion, err = semver.BumpPatch(currentVersion, cfg.UseVPrefix)
	default:
		err = fmt.Errorf("invalid bump type: %s - must be one of: major, minor, patch", cfg.BumpType)
	}

	if err != nil {
		logger.Version.Error().Err(err).Str("type", cfg.BumpType).Str("current_version", currentVersion).Msg("Failed to bump version")
		return nil, fmt.Errorf("failed to bump version: %w - check if the current version (%s) is a valid semantic version in the format X.Y.Z", err, currentVersion)
	}

	fmt.Fprintf(output, "New version: %s\n", newVersion)

	// Update changelog
	fmt.Fprintf(output, "Updating changelog file: %s\n", cfg.ChangelogFile)
	err = changelog.UpdateChangelog(cfg.ChangelogFile, newVersion, cfg.RepositoryProvider)
	if err != nil {
		logger.Version.Error().Err(err).Str("file", cfg.ChangelogFile).Str("version", newVersion).Msg("Failed to update changelog")
		return nil, fmt.Errorf("failed to update changelog: %w - verify that '%s' exists and follows the Keep a Changelog format", err, cfg.ChangelogFile)
	}

	// Commit changes
	err = git.CommitChangelog(cfg.ChangelogFile, newVersion)
	if err != nil {
		logger.Version.Error().Err(err).Str("file", cfg.ChangelogFile).Str("version", newVersion).Msg("Failed to commit changelog")
		return nil, fmt.Errorf("failed to commit changelog: %w - ensure git is properly configured and you have permissions to commit changes", err)
	}

	// Tag version
	fmt.Fprintf(output, "Tagging version: %s\n", newVersion)
	err = git.TagVersion(newVersion)
	if err != nil {
		logger.Version.Error().Err(err).Str("version", newVersion).Msg("Failed to tag version")
		return nil, fmt.Errorf("failed to tag version: %w - check if the tag already exists (use 'git tag' to list existing tags)", err)
	}

	fmt.Fprintf(output, "%s release %s done.\n", cfg.BumpType, newVersion)

	// Auto-push if enabled
	pushed := false
	if cfg.AutoPush {
		fmt.Fprintf(output, "Pushing changes and tags...\n")
		err = git.PushChanges()
		if err != nil {
			logger.Version.Error().Err(err).Msg("Failed to push changes")
			return nil, fmt.Errorf("failed to push changes: %w - check network connection and remote repository permissions", err)
		}
		fmt.Fprintf(output, "Automatically pushed changes and tags to remote repository.\n")
		pushed = true
	} else {
		fmt.Fprintf(output, "Don't forget to git push and git push --tags.\n")
	}

	logger.Version.Debug().Str("type", cfg.BumpType).Str("version", newVersion).Msg("Version bump completed successfully")

	result := &BumpResult{
		OldVersion:    currentVersion,
		NewVersion:    newVersion,
		BumpType:      cfg.BumpType,
		ChangelogFile: cfg.ChangelogFile,
		Pushed:        pushed,
	}

	return result, nil
}
