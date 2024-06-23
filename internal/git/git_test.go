package git

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var executedCmd []string

func mockExecCommand(command string, args ...string) *exec.Cmd {
	executedCmd = append([]string{command}, args...)
	return exec.Command("echo", "mocked")
}
func TestIsInstalled(t *testing.T) {
	if !IsInstalled() {
		t.Error("Git is not installed, but it should be for running these tests")
	}
}

func TestGetProjectVersion(t *testing.T) {
	// Set up a temporary git repository
	tempDir, err := os.MkdirTemp("", "changie-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	os.Chdir(tempDir)
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	// Create an initial commit
	exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit").Run()

	// Create a tag
	exec.Command("git", "tag", "1.0.0").Run()

	version, err := GetProjectVersion()
	if err != nil {
		t.Fatalf("GetProjectVersion() returned an error: %v", err)
	}
	if version != "1.0.0" {
		t.Errorf("GetProjectVersion() = %s, expected 1.0.0", version)
	}
}

func TestCommitChangelog(t *testing.T) {
	// Set up a temporary git repository
	tempDir, err := os.MkdirTemp("", "changie-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	os.Chdir(tempDir)
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	// Create a dummy changelog file
	os.WriteFile("CHANGELOG.md", []byte("# Changelog"), 0644)
	exec.Command("git", "add", "CHANGELOG.md").Run()

	err = CommitChangelog("CHANGELOG.md", "1.0.0")
	if err != nil {
		t.Fatalf("CommitChangelog() returned an error: %v", err)
	}

	// Check if the commit was created
	out, err := exec.Command("git", "log", "--oneline").Output()
	if err != nil {
		t.Fatalf("Failed to get git log: %v", err)
	}
	if !strings.Contains(string(out), "Update changelog for version 1.0.0") {
		t.Error("Commit message not found in git log")
	}
}

func TestTagVersion(t *testing.T) {
	// Set up a temporary git repository
	tempDir, err := os.MkdirTemp("", "changie-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	os.Chdir(tempDir)
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()
	exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit").Run()

	err = TagVersion("1.0.0")
	if err != nil {
		t.Fatalf("TagVersion() returned an error: %v", err)
	}

	// Check if the tag was created
	out, err := exec.Command("git", "tag").Output()
	if err != nil {
		t.Fatalf("Failed to get git tags: %v", err)
	}
	if !strings.Contains(string(out), "1.0.0") {
		t.Error("Tag 1.0.0 not found in git tags")
	}
}
func TestTagVersionWithoutPrefix(t *testing.T) {
	// Mock exec.Command
	oldExecCommand := execCommand
	execCommand = mockExecCommand
	defer func() { execCommand = oldExecCommand }()

	// Reset executedCmd before the test
	executedCmd = nil

	// Test tagging version without 'v' prefix
	err := TagVersion("1.0.0")
	if err != nil {
		t.Errorf("TagVersion failed: %v", err)
	}

	// Check if the correct command was executed
	expectedCmd := []string{"git", "tag", "1.0.0"}
	if len(executedCmd) != len(expectedCmd) {
		t.Errorf("Expected command %v, got: %v", expectedCmd, executedCmd)
	} else {
		for i := range expectedCmd {
			if executedCmd[i] != expectedCmd[i] {
				t.Errorf("Expected command %v, got: %v", expectedCmd, executedCmd)
				break
			}
		}
	}
}
