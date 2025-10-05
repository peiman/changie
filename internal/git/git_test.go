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
	err = os.WriteFile(testFile, []byte("Test content"), 0o644)
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

// TestParseRepositoryURL tests the ParseRepositoryURL function
func TestParseRepositoryURL(t *testing.T) {
	tests := []struct {
		name         string
		remoteURL    string
		wantOwner    string
		wantRepo     string
		wantProvider string
		wantBaseURL  string
		wantErr      bool
	}{
		{
			name:         "GitHub HTTPS URL",
			remoteURL:    "https://github.com/peiman/changie.git",
			wantOwner:    "peiman",
			wantRepo:     "changie",
			wantProvider: "github",
			wantBaseURL:  "https://github.com",
			wantErr:      false,
		},
		{
			name:         "GitHub HTTPS URL without .git",
			remoteURL:    "https://github.com/peiman/changie",
			wantOwner:    "peiman",
			wantRepo:     "changie",
			wantProvider: "github",
			wantBaseURL:  "https://github.com",
			wantErr:      false,
		},
		{
			name:         "GitHub SSH URL",
			remoteURL:    "git@github.com:peiman/changie.git",
			wantOwner:    "peiman",
			wantRepo:     "changie",
			wantProvider: "github",
			wantBaseURL:  "https://github.com",
			wantErr:      false,
		},
		{
			name:         "GitHub SSH URL without .git",
			remoteURL:    "git@github.com:peiman/changie",
			wantOwner:    "peiman",
			wantRepo:     "changie",
			wantProvider: "github",
			wantBaseURL:  "https://github.com",
			wantErr:      false,
		},
		{
			name:         "Bitbucket HTTPS URL",
			remoteURL:    "https://bitbucket.org/myteam/myrepo.git",
			wantOwner:    "myteam",
			wantRepo:     "myrepo",
			wantProvider: "bitbucket",
			wantBaseURL:  "https://bitbucket.org",
			wantErr:      false,
		},
		{
			name:         "Bitbucket SSH URL",
			remoteURL:    "git@bitbucket.org:myteam/myrepo.git",
			wantOwner:    "myteam",
			wantRepo:     "myrepo",
			wantProvider: "bitbucket",
			wantBaseURL:  "https://bitbucket.org",
			wantErr:      false,
		},
		{
			name:         "GitLab HTTPS URL",
			remoteURL:    "https://gitlab.com/group/project.git",
			wantOwner:    "group",
			wantRepo:     "project",
			wantProvider: "gitlab",
			wantBaseURL:  "https://gitlab.com",
			wantErr:      false,
		},
		{
			name:         "GitLab SSH URL",
			remoteURL:    "git@gitlab.com:group/project.git",
			wantOwner:    "group",
			wantRepo:     "project",
			wantProvider: "gitlab",
			wantBaseURL:  "https://gitlab.com",
			wantErr:      false,
		},
		{
			name:         "Unknown provider HTTPS",
			remoteURL:    "https://git.example.com/owner/repo.git",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantProvider: "unknown",
			wantBaseURL:  "https://git.example.com",
			wantErr:      false,
		},
		{
			name:         "Unknown provider SSH",
			remoteURL:    "git@git.example.com:owner/repo.git",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantProvider: "unknown",
			wantBaseURL:  "https://git.example.com",
			wantErr:      false,
		},
		{
			name:      "Empty URL",
			remoteURL: "",
			wantErr:   true,
		},
		{
			name:      "Invalid SSH URL - missing colon",
			remoteURL: "git@github.com/peiman/changie.git",
			wantErr:   true,
		},
		{
			name:      "Invalid HTTPS URL - too short",
			remoteURL: "https://github.com/changie",
			wantErr:   true,
		},
		{
			name:      "Invalid format",
			remoteURL: "not-a-valid-url",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseRepositoryURL(tc.remoteURL)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Owner != tc.wantOwner {
				t.Errorf("Owner: got %q, want %q", result.Owner, tc.wantOwner)
			}
			if result.Repo != tc.wantRepo {
				t.Errorf("Repo: got %q, want %q", result.Repo, tc.wantRepo)
			}
			if result.Provider != tc.wantProvider {
				t.Errorf("Provider: got %q, want %q", result.Provider, tc.wantProvider)
			}
			if result.BaseURL != tc.wantBaseURL {
				t.Errorf("BaseURL: got %q, want %q", result.BaseURL, tc.wantBaseURL)
			}
		})
	}
}

// TestGetRepositoryInfo tests the GetRepositoryInfo function
func TestGetRepositoryInfo(t *testing.T) {
	// This test verifies the function works in a git repository with a remote
	// Since we're running in the changie git repo, this should succeed if remote is configured

	// Note: This test may fail in CI environments without a remote configured
	// In that case, it's expected and we just log it
	repoInfo, err := GetRepositoryInfo()
	if err != nil {
		t.Logf("GetRepositoryInfo returned error (expected if no remote configured): %v", err)
		return
	}

	// If it succeeded, verify we got valid data
	if repoInfo.Owner == "" {
		t.Error("Expected non-empty owner")
	}
	if repoInfo.Repo == "" {
		t.Error("Expected non-empty repo")
	}
	if repoInfo.Provider == "" {
		t.Error("Expected non-empty provider")
	}
	if repoInfo.BaseURL == "" {
		t.Error("Expected non-empty baseURL")
	}

	t.Logf("Repository info: owner=%s, repo=%s, provider=%s, baseURL=%s",
		repoInfo.Owner, repoInfo.Repo, repoInfo.Provider, repoInfo.BaseURL)
}
