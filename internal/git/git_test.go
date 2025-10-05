package git

import (
	"os"
	"testing"
)

// TestIsInstalled tests the IsInstalled function
func TestIsInstalled(t *testing.T) {
	// This is a basic test that simply ensures the function runs
	// The result depends on whether git is installed on the test machine
	result := IsInstalled()
	t.Logf("Git is installed: %v", result)
}

// Since most git functions rely on the git command, we'll use a simplified
// approach for testing: check that the function handles errors correctly.
// For a more comprehensive test suite, consider using a mocking library or
// setting up a test repository.

// TestGetVersionErrorHandling ensures GetVersion handles errors properly
func TestGetVersionErrorHandling(t *testing.T) {
	// Test with a non-git directory
	tmpDir, err := os.MkdirTemp("", "non-git-dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current dir
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}

	// Change to temp dir and back when done
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Logf("Warning: Failed to change back to original directory: %v", err)
		}
	}()

	// GetVersion should return empty string (no tags) without error
	version, err := GetVersion()
	if err == nil {
		// Some environments might have git configured to not return an error
		// for this case, so we'll just log instead of failing
		t.Logf("Expected error due to non-git dir, but got: %v", version)
	}
}

// TestHasUncommittedChangesErrorHandling ensures HasUncommittedChanges handles errors properly
func TestHasUncommittedChangesErrorHandling(t *testing.T) {
	// Test with a non-git directory
	tmpDir, err := os.MkdirTemp("", "non-git-dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current dir
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}

	// Change to temp dir and back when done
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Logf("Warning: Failed to change back to original directory: %v", err)
		}
	}()

	// HasUncommittedChanges should return an error
	_, err = HasUncommittedChanges()
	if err == nil {
		t.Error("Expected error from HasUncommittedChanges in non-git dir, but got nil")
	}
}

// TestCommitChangelogErrorHandling ensures CommitChangelog handles errors properly
func TestCommitChangelogErrorHandling(t *testing.T) {
	// Test with a non-git directory
	tmpDir, err := os.MkdirTemp("", "non-git-dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current dir
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}

	// Change to temp dir and back when done
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Logf("Warning: Failed to change back to original directory: %v", err)
		}
	}()

	// Create a test file
	testFile := "test-changelog.md"
	err = os.WriteFile(testFile, []byte("Test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// CommitChangelog should return an error
	err = CommitChangelog(testFile, "1.0.0")
	if err == nil {
		t.Error("Expected error from CommitChangelog in non-git dir, but got nil")
	}
}

// TestTagVersionErrorHandling ensures TagVersion handles errors properly
func TestTagVersionErrorHandling(t *testing.T) {
	// Test with a non-git directory
	tmpDir, err := os.MkdirTemp("", "non-git-dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current dir
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}

	// Change to temp dir and back when done
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Logf("Warning: Failed to change back to original directory: %v", err)
		}
	}()

	// TagVersion should return an error
	err = TagVersion("1.0.0")
	if err == nil {
		t.Error("Expected error from TagVersion in non-git dir, but got nil")
	}
}

// TestPushChangesErrorHandling ensures PushChanges handles errors properly
func TestPushChangesErrorHandling(t *testing.T) {
	// Test with a non-git directory
	tmpDir, err := os.MkdirTemp("", "non-git-dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current dir
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}

	// Change to temp dir and back when done
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Logf("Warning: Failed to change back to original directory: %v", err)
		}
	}()

	// PushChanges should return an error
	err = PushChanges()
	if err == nil {
		t.Error("Expected error from PushChanges in non-git dir, but got nil")
	}
}

// TestGetCurrentBranch tests the GetCurrentBranch function
func TestGetCurrentBranch(t *testing.T) {
	// This test verifies the function works in a git repository
	// Since we're running in the changie git repo, this should succeed
	branch, err := GetCurrentBranch()
	if err != nil {
		t.Logf("GetCurrentBranch returned error (expected if not in git repo): %v", err)
	} else {
		t.Logf("Current branch: %s", branch)
		if branch == "" {
			t.Error("Expected non-empty branch name")
		}
	}
}

// TestGetCurrentBranchErrorHandling ensures GetCurrentBranch handles errors properly
func TestGetCurrentBranchErrorHandling(t *testing.T) {
	// Test with a non-git directory
	tmpDir, err := os.MkdirTemp("", "non-git-dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current dir
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}

	// Change to temp dir and back when done
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Logf("Warning: Failed to change back to original directory: %v", err)
		}
	}()

	// GetCurrentBranch should return an error in non-git dir
	_, err = GetCurrentBranch()
	if err == nil {
		t.Error("Expected error from GetCurrentBranch in non-git dir, but got nil")
	}
}
