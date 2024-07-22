package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Commander is an interface for command execution
type Commander interface {
	CombinedOutput() ([]byte, error)
}

// ExecCommand is a variable that holds the function to execute commands
var ExecCommand = func(command string, args ...string) Commander {
	return exec.Command(command, args...)
}

// IsInstalled checks if Git is installed
func IsInstalled() bool {
	cmd := ExecCommand("git", "--version")
	_, err := cmd.CombinedOutput()
	return err == nil
}

// GetVersion retrieves the current version based on git tags and commits
func GetVersion() (string, error) {
	cmd := ExecCommand("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "No names found") {
			return "dev", nil // Return "dev" without an error when no tags are found
		}
		return "", fmt.Errorf("error getting latest tag: %w", err)
	}
	tag := strings.TrimSpace(string(output))

	// Check if the current commit is tagged
	cmd = ExecCommand("git", "describe", "--exact-match", "--tags", "HEAD")
	if _, err := cmd.CombinedOutput(); err == nil {
		// Current commit is tagged, return the tag
		return tag, nil
	}

	// Get the current commit hash
	cmd = ExecCommand("git", "rev-parse", "--short", "HEAD")
	commitHashOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error getting commit hash: %w", err)
	}
	commitHash := strings.TrimSpace(string(commitHashOutput))

	// Count commits since the latest tag
	cmd = ExecCommand("git", "rev-list", tag+"..HEAD", "--count")
	revListOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error counting commits since last tag: %w", err)
	}
	commitCount := strings.TrimSpace(string(revListOutput))

	return fmt.Sprintf("%s-dev.%s+%s", tag, commitCount, commitHash), nil
}

// CommitChangelog commits the changelog file
func CommitChangelog(file, version string) error {
	addCmd := ExecCommand("git", "add", file)
	_, err := addCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding changelog to git: %w", err)
	}

	commitCmd := ExecCommand("git", "commit", "-m", fmt.Sprintf("Update changelog for version %s", version))
	_, err = commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error committing changelog: %w", err)
	}

	return nil
}

// TagVersion creates a new Git tag for the given version
func TagVersion(version string) error {
	cmd := ExecCommand("git", "tag", version)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error tagging version: %w", err)
	}
	return nil
}

// HasUncommittedChanges checks if there are any uncommitted changes in the repository
func HasUncommittedChanges() (bool, error) {
	cmd := ExecCommand("git", "status", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}
	return len(output) > 0, nil
}

// PushChanges pushes the changes and tags to the remote repository
func PushChanges() error {
	cmd := ExecCommand("git", "push", "--follow-tags")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push changes: %w\nCommand output: %s", err, string(output))
	}
	return nil
}
