// Package version provides version bump orchestration functionality.
//
// This package handles the workflow for bumping semantic versions,
// including validation, changelog updates, git operations, and pushing changes.
package version

import (
	"fmt"
	"io"
	"os"

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

// Bump performs a complete version bump workflow.
//
// The workflow is atomic: if any step fails after mutations begin,
// all changes are rolled back (changelog restored, commit reverted, tag deleted).
//
// Phase 1 — Preflight (no mutations):
//  1. Verify git is installed
//  2. Check branch (main/master unless bypassed)
//  3. Check for uncommitted changes
//  4. Get current version from git tags
//  5. Calculate new version
//  6. Read current changelog content (for rollback)
//
// Phase 2 — Mutate (with rollback on failure):
//  7. Update changelog file
//  8. Commit changelog
//  9. Create git tag
//  10. Optionally push
func Bump(cfg BumpConfig, output io.Writer) error {
	// ── Phase 1: Preflight ──────────────────────────────────────────────
	logger.Version.Debug().Str("type", cfg.BumpType).Msg("Starting version bump")

	if !git.IsInstalled() {
		return fmt.Errorf("git is not installed or not available in PATH - please install Git (https://git-scm.com/downloads) and ensure it's in your system PATH")
	}

	if !cfg.AllowAnyBranch {
		currentBranch, err := git.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		if currentBranch != "main" && currentBranch != "master" {
			return fmt.Errorf("not on main/master branch (current: %s) - use --allow-any-branch to bypass", currentBranch)
		}
		logger.Version.Debug().Str("branch", currentBranch).Msg("Branch check passed")
	}

	hasUncommitted, err := git.HasUncommittedChanges()
	if err != nil {
		return fmt.Errorf("failed to check for uncommitted changes: %w", err)
	}
	if hasUncommitted {
		return fmt.Errorf("uncommitted changes found - commit or stash them before bumping version")
	}

	currentVersion, err := git.GetVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w - ensure you have at least one tag, or run 'git tag v0.0.0'", err)
	}
	if currentVersion == "" {
		currentVersion = "0.0.0"
		_, _ = fmt.Fprintf(output, "No version tag found, starting from %s\n", currentVersion)
	} else {
		_, _ = fmt.Fprintf(output, "Current version: %s\n", currentVersion)
	}

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
		return fmt.Errorf("failed to bump version: %w", err)
	}

	_, _ = fmt.Fprintf(output, "New version: %s\n", newVersion)

	// Save changelog content for rollback
	originalChangelog, err := os.ReadFile(cfg.ChangelogFile) //nolint:gosec // G304: user-specified changelog path
	if err != nil {
		return fmt.Errorf("failed to read changelog for backup: %w", err)
	}

	// ── Phase 2: Mutate (with rollback) ─────────────────────────────────
	rollback := &rollbackState{
		changelogFile:   cfg.ChangelogFile,
		originalContent: originalChangelog,
		committed:       false,
		tagged:          false,
		tagName:         newVersion,
	}

	_, _ = fmt.Fprintf(output, "Updating changelog file: %s\n", cfg.ChangelogFile)
	if err := changelog.UpdateChangelog(cfg.ChangelogFile, newVersion, cfg.RepositoryProvider); err != nil {
		rollback.execute(output)
		return fmt.Errorf("failed to update changelog: %w", err)
	}

	if err := git.CommitChangelog(cfg.ChangelogFile, newVersion); err != nil {
		rollback.execute(output)
		return fmt.Errorf("failed to commit changelog: %w", err)
	}
	rollback.committed = true

	_, _ = fmt.Fprintf(output, "Tagging version: %s\n", newVersion)
	if err := git.TagVersion(newVersion); err != nil {
		rollback.execute(output)
		return fmt.Errorf("failed to tag version: %w", err)
	}
	rollback.tagged = true

	_, _ = fmt.Fprintf(output, "%s release %s done.\n", cfg.BumpType, newVersion)

	if cfg.AutoPush {
		_, _ = fmt.Fprintf(output, "Pushing changes and tags...\n")
		if err := git.PushChanges(); err != nil {
			rollback.execute(output)
			return fmt.Errorf("failed to push changes: %w", err)
		}
		_, _ = fmt.Fprintf(output, "Pushed changes and tags to remote.\n")
	} else {
		_, _ = fmt.Fprintf(output, "Don't forget to git push and git push --tags.\n")
	}

	logger.Version.Debug().Str("type", cfg.BumpType).Str("version", newVersion).Msg("Version bump completed")
	return nil
}

// rollbackState tracks which mutations have been applied so they can be undone.
type rollbackState struct {
	changelogFile   string
	originalContent []byte
	committed       bool
	tagged          bool
	tagName         string
}

// execute undoes all completed mutations in reverse order.
func (r *rollbackState) execute(output io.Writer) {
	_, _ = fmt.Fprintf(output, "Rolling back changes...\n")

	if r.tagged {
		if err := git.DeleteTag(r.tagName); err != nil {
			logger.Version.Error().Err(err).Str("tag", r.tagName).Msg("Rollback: failed to delete tag")
		}
	}

	if r.committed {
		if err := git.UndoLastCommit(); err != nil {
			logger.Version.Error().Err(err).Msg("Rollback: failed to undo commit")
		}
	}

	// Always restore the changelog file to its original content
	if err := os.WriteFile(r.changelogFile, r.originalContent, 0o644); err != nil { //nolint:gosec // G306: changelog needs 0644
		logger.Version.Error().Err(err).Str("file", r.changelogFile).Msg("Rollback: failed to restore changelog")
	}

	_, _ = fmt.Fprintf(output, "Rollback complete.\n")
}
