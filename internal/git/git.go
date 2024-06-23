package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsInstalled checks if Git is installed
func IsInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// GetProjectVersion gets the current project version from Git tags
func GetProjectVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0", "--match", "[0-9]*.[0-9]*.[0-9]*")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting project version: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// CommitChangelog commits the updated changelog
func CommitChangelog(changelogFile, version string) error {
	commitMsg := fmt.Sprintf("Update changelog for version %s", version)
	cmd := exec.Command("git", "commit", "-am", commitMsg)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error committing changelog: %w", err)
	}
	return nil
}

// TagVersion creates a new Git tag for the given version
func TagVersion(version string) error {
	cmd := exec.Command("git", "tag", version)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error tagging version: %w", err)
	}
	return nil
}
