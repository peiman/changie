// Package git provides functionality for working with Git repositories.
//
// This package encapsulates Git operations used by the changie tool, including:
// - Version detection from Git tags
// - Checking for uncommitted changes
// - Managing commits and tags for changelog updates
// - Pushing changes to remote repositories
//
// All functions in this package use the git command-line tool and require it to be
// installed and available in the PATH. Functions will return appropriate errors if
// Git operations fail.
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsInstalled checks if git is installed and available in the PATH.
//
// This function attempts to run "git --version" and returns true only if
// the command executes successfully, indicating that Git is properly installed.
//
// Returns:
//   - bool: true if git is installed and available, false otherwise
func IsInstalled() bool {
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	return err == nil
}

// GetVersion returns the current version from git tags.
//
// This function uses "git describe --tags --abbrev=0" to retrieve the most recent tag.
// If no tags exist in the repository, it returns an empty string, which the caller
// should interpret as version 0.0.0.
//
// Returns:
//   - string: The version string (without 'v' prefix) or empty string if no tags exist
//   - error: Any error encountered during the git operation (except for "no tags" errors)
func GetVersion() (string, error) {
	// Use git describe to get the latest tag
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.CombinedOutput()

	// If there are no tags, return empty string which will be treated as 0.0.0
	if err != nil {
		// Check if error is due to no tags
		if strings.Contains(string(output), "No names found") || strings.Contains(string(output), "fatal: No names found") {
			return "", nil
		}
		return "", fmt.Errorf("failed to get latest tag: %w (verify that you are in a git repository and have sufficient permissions)", err)
	}

	// Trim and return the version string
	version := strings.TrimSpace(string(output))
	return version, nil
}

// GetCurrentBranch returns the name of the current git branch.
//
// This function uses "git rev-parse --abbrev-ref HEAD" to get the current branch name.
// This is a reliable way to get the branch name that works in all scenarios.
//
// Returns:
//   - string: The name of the current branch
//   - error: Any error encountered during the git operation
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w (verify you're in a git repository)", err)
	}

	branch := strings.TrimSpace(string(output))
	return branch, nil
}

// HasUncommittedChanges checks if there are any uncommitted changes in the repository.
//
// This function uses "git status --porcelain" to determine if there are any staged
// or unstaged changes that haven't been committed. The porcelain format ensures
// machine-readable output that's stable across git versions.
//
// Returns:
//   - bool: true if there are uncommitted changes, false otherwise
//   - error: Any error encountered during the git operation
func HasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check for uncommitted changes: %w (verify you're in a git repository with proper permissions)", err)
	}

	// If there is any output, there are uncommitted changes
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// CommitChangelog commits changes to the changelog file.
//
// This function:
// 1. Adds the specified changelog file to the staging area using "git add"
// 2. Commits the changes with a standardized message that includes the version
//
// Parameters:
//   - file: Path to the changelog file to commit
//   - version: Version number to include in the commit message
//
// Returns:
//   - error: Any error encountered during the git operations
func CommitChangelog(file, version string) error {
	// Add the file to the staging area
	cmd := exec.Command("git", "add", file)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add changelog to staging area: %w (check if the file '%s' exists and you have write permissions)", err, file)
	}

	// Commit the changes
	commitMsg := fmt.Sprintf("Update changelog for version %s", version)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit changelog: %w (ensure git user.name and user.email are configured correctly with 'git config')", err)
	}

	return nil
}

// TagVersion creates a new git tag for the given version.
//
// This function creates an annotated tag (-a) with a standardized message
// that includes the version number. It respects the user's configured
// preference for using 'v' prefix or not.
//
// Parameters:
//   - version: Version number to tag (with or without "v" prefix)
//
// Returns:
//   - error: Any error encountered during the git operation
func TagVersion(version string) error {
	// The version parameter already has the correct prefix (or no prefix)
	// based on the user's preference, so we use it directly
	tagName := version

	// Create the tag
	tagMsg := fmt.Sprintf("Version %s", version)
	cmd := exec.Command("git", "tag", "-a", tagName, "-m", tagMsg)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create tag: %w (check if tag '%s' already exists, you can delete it with 'git tag -d %s')", err, tagName, tagName)
	}

	return nil
}

// PushChanges pushes commits and tags to the remote repository.
//
// This function performs two git operations:
// 1. "git push" to push commits to the remote repository
// 2. "git push --tags" to push tags to the remote repository
//
// Returns:
//   - error: Any error encountered during the git operations
func PushChanges() error {
	// Push commits
	cmd := exec.Command("git", "push")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push commits: %w (check network connection and remote repository access permissions)", err)
	}

	// Push tags
	cmd = exec.Command("git", "push", "--tags")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push tags: %w (check if you have permission to create tags on the remote repository)", err)
	}

	return nil
}
