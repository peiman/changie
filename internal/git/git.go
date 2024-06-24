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
func CommitChangelog(file, version string) error {
	log.Printf("Attempting to commit changelog for version %s", version)

	// Add the file
	addCmd := exec.Command("git", "add", file)
	log.Printf("Executing command: %v", addCmd.Args)
	output, err := addCmd.CombinedOutput()
	if err != nil {
		log.Printf("Error adding file: %v. Output: %s", err, output)
		return fmt.Errorf("error adding changelog: %w", err)
	}

	// Commit the changes
	commitCmd := exec.Command("git", "commit", "-m", fmt.Sprintf("Update changelog for version %s", version))
	log.Printf("Executing command: %v", commitCmd.Args)
	output, err = commitCmd.CombinedOutput()
	if err != nil {
		log.Printf("Error committing changes: %v. Output: %s", err, output)
		return fmt.Errorf("error committing changelog: %w", err)
	}

	log.Printf("Successfully committed changelog for version %s", version)
	return nil
}

// TagVersion creates a new Git tag for the given version
func TagVersion(version string) error {
	log.Printf("Attempting to tag version %s", version)

	cmd := exec.Command("git", "tag", version)
	log.Printf("Executing command: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error tagging version: %v. Output: %s", err, output)
		return fmt.Errorf("error tagging version: %w", err)
	}

	log.Printf("Successfully tagged version %s", version)
	return nil
}
