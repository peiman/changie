// Package git provides functionality for working with Git repositories.
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsInstalled checks if git is installed and available.
func IsInstalled() bool {
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	return err == nil
}

// GetVersion returns the current version from git tags.
// If no tags exist, returns an empty string.
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
		return "", fmt.Errorf("failed to get latest tag: %w", err)
	}

	// Trim and return the version string
	version := strings.TrimSpace(string(output))
	return version, nil
}

// HasUncommittedChanges checks if there are any uncommitted changes in the repository.
func HasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check for uncommitted changes: %w", err)
	}

	// If there is any output, there are uncommitted changes
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// CommitChangelog commits changes to the changelog file.
func CommitChangelog(file, version string) error {
	// Add the file to the staging area
	cmd := exec.Command("git", "add", file)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add changelog to staging area: %w", err)
	}

	// Commit the changes
	commitMsg := fmt.Sprintf("Update changelog for version %s", version)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit changelog: %w", err)
	}

	return nil
}

// TagVersion creates a new git tag for the given version.
func TagVersion(version string) error {
	// Ensure version has v prefix
	tagName := version
	if !strings.HasPrefix(version, "v") {
		tagName = "v" + version
	}

	// Create the tag
	tagMsg := fmt.Sprintf("Version %s", version)
	cmd := exec.Command("git", "tag", "-a", tagName, "-m", tagMsg)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

// PushChanges pushes commits and tags to the remote repository.
func PushChanges() error {
	// Push commits
	cmd := exec.Command("git", "push")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push commits: %w", err)
	}

	// Push tags
	cmd = exec.Command("git", "push", "--tags")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push tags: %w", err)
	}

	return nil
}
