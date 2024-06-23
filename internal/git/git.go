package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var execCommand = exec.Command

// IsInstalled checks if Git is installed
func IsInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// GetProjectVersion gets the current project version from Git tags
func GetProjectVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	log.Printf("Executing command: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command output: %s", string(output))
		log.Printf("Error executing git command: %v", err)
		return "", fmt.Errorf("error getting project version: %v", err)
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
	cmd := execCommand("git", "tag", version)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error tagging version: %v\nCommand output: %s", err, string(out))
	}
	return nil
}
