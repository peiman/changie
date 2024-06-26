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
	_, err := exec.LookPath("git")
	return err == nil
}

// GetProjectVersion retrieves the current project version from Git tags
func GetProjectVersion() (string, error) {
	cmd := ExecCommand("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "No names found, cannot describe anything") {
			return "0.0.0", nil // Return 0.0.0 as the initial version
		}
		return "", fmt.Errorf("error getting project version: %v: %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

// CommitChangelog commits the changelog file
func CommitChangelog(file, version string) error {
	addCmd := ExecCommand("git", "add", file)
	_, err := addCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding changelog to git: %v", err)
	}

	commitCmd := ExecCommand("git", "commit", "-m", fmt.Sprintf("Update changelog for version %s", version))
	_, err = commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error committing changelog: %v", err)
	}

	return nil
}

// TagVersion creates a new Git tag for the given version
func TagVersion(version string) error {
	cmd := ExecCommand("git", "tag", version)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error tagging version: %v", err)
	}
	return nil
}
