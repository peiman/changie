package version

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestGitRepo creates a temporary git repository for testing
//
//nolint:unparam // tempDir return value is intentionally kept for future use
func setupTestGitRepo(t *testing.T, initialTag, branch string, addUncommitted bool) (string, func()) {
	t.Helper()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "version-bump-test-*")
	require.NoError(t, err, "Failed to create temp dir")

	// Save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current dir")

	// Change to temp directory
	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change to temp dir")

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err = cmd.Run()
	require.NoError(t, err, "Failed to init git repo")

	// Configure git
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()
	exec.Command("git", "config", "tag.gpgsign", "false").Run()
	exec.Command("git", "config", "commit.gpgsign", "false").Run()

	// Create initial changelog
	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	initialChangelog := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- New feature for testing

`
	err = os.WriteFile(changelogPath, []byte(initialChangelog), 0o644)
	require.NoError(t, err, "Failed to create changelog")

	// Add and commit the changelog
	exec.Command("git", "add", "CHANGELOG.md").Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()

	// Create initial tag if provided
	if initialTag != "" {
		exec.Command("git", "tag", initialTag).Run()
	}

	// Checkout the desired branch
	if branch != "main" && branch != "master" && branch != "" {
		// Create and checkout new branch
		exec.Command("git", "checkout", "-b", branch).Run()
	} else if branch == "master" {
		// Rename main to master
		exec.Command("git", "branch", "-m", "main", "master").Run()
	}

	// Add uncommitted changes if needed
	if addUncommitted {
		testFile := filepath.Join(tempDir, "test.txt")
		os.WriteFile(testFile, []byte("uncommitted"), 0o644)
	}

	// Cleanup function
	cleanup := func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Warning: Failed to change back to original directory: %v", err)
		}
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestBump(t *testing.T) {
	// Set up a buffer to capture log output
	var logBuf bytes.Buffer
	log.Logger = zerolog.New(&logBuf)

	tests := []struct {
		name            string
		cfg             BumpConfig
		setupBranch     string
		initialTag      string
		addUncommitted  bool
		expectedVersion string
		expectError     bool
		errorContains   string
		outputContains  []string
	}{
		{
			name: "major bump with v prefix",
			cfg: BumpConfig{
				BumpType:           "major",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:     "main",
			initialTag:      "v1.2.3",
			addUncommitted:  false,
			expectedVersion: "v2.0.0",
			expectError:     false,
			outputContains: []string{
				"Current version: v1.2.3",
				"New version: v2.0.0",
				"major release v2.0.0 done",
			},
		},
		{
			name: "minor bump with v prefix",
			cfg: BumpConfig{
				BumpType:           "minor",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:     "main",
			initialTag:      "v1.2.3",
			addUncommitted:  false,
			expectedVersion: "v1.3.0",
			expectError:     false,
			outputContains: []string{
				"Current version: v1.2.3",
				"New version: v1.3.0",
				"minor release v1.3.0 done",
			},
		},
		{
			name: "patch bump with v prefix",
			cfg: BumpConfig{
				BumpType:           "patch",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:     "main",
			initialTag:      "v1.2.3",
			addUncommitted:  false,
			expectedVersion: "v1.2.4",
			expectError:     false,
			outputContains: []string{
				"Current version: v1.2.3",
				"New version: v1.2.4",
				"patch release v1.2.4 done",
			},
		},
		{
			name: "major bump without v prefix",
			cfg: BumpConfig{
				BumpType:           "major",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         false,
			},
			setupBranch:     "main",
			initialTag:      "1.2.3",
			addUncommitted:  false,
			expectedVersion: "2.0.0",
			expectError:     false,
			outputContains: []string{
				"Current version: 1.2.3",
				"New version: 2.0.0",
				"major release 2.0.0 done",
			},
		},
		{
			name: "bump on feature branch fails without allow-any-branch",
			cfg: BumpConfig{
				BumpType:           "minor",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:    "feature/test",
			initialTag:     "v1.2.3",
			addUncommitted: false,
			expectError:    true,
			errorContains:  "not on main/master branch",
		},
		{
			name: "bump on feature branch succeeds with allow-any-branch",
			cfg: BumpConfig{
				BumpType:           "minor",
				AllowAnyBranch:     true,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:     "feature/test",
			initialTag:      "v1.2.3",
			addUncommitted:  false,
			expectedVersion: "v1.3.0",
			expectError:     false,
			outputContains: []string{
				"New version: v1.3.0",
			},
		},
		{
			name: "bump on master branch succeeds",
			cfg: BumpConfig{
				BumpType:           "patch",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:     "master",
			initialTag:      "v0.1.0",
			addUncommitted:  false,
			expectedVersion: "v0.1.1",
			expectError:     false,
			outputContains: []string{
				"New version: v0.1.1",
			},
		},
		{
			name: "bump with uncommitted changes fails",
			cfg: BumpConfig{
				BumpType:           "patch",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:    "main",
			initialTag:     "v1.2.3",
			addUncommitted: true,
			expectError:    true,
			errorContains:  "uncommitted changes found",
		},
		{
			name: "bump with prerelease tag clears metadata",
			cfg: BumpConfig{
				BumpType:           "patch",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:     "main",
			initialTag:      "v1.2.3-alpha.1+build.123",
			addUncommitted:  false,
			expectedVersion: "v1.2.4",
			expectError:     false,
			outputContains: []string{
				"New version: v1.2.4",
			},
		},
		{
			name: "bump from no tags starts at 0.0.0",
			cfg: BumpConfig{
				BumpType:           "patch",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:     "main",
			initialTag:      "", // No initial tag
			addUncommitted:  false,
			expectedVersion: "v0.0.1",
			expectError:     false,
			outputContains: []string{
				"No version tag found",
				"New version: v0.0.1",
			},
		},
		{
			name: "invalid bump type returns error",
			cfg: BumpConfig{
				BumpType:           "invalid",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:    "main",
			initialTag:     "v1.2.3",
			addUncommitted: false,
			expectError:    true,
			errorContains:  "invalid bump type",
		},
		{
			name: "bump with auto-push disabled shows reminder",
			cfg: BumpConfig{
				BumpType:           "minor",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			setupBranch:     "main",
			initialTag:      "v1.0.0",
			addUncommitted:  false,
			expectedVersion: "v1.1.0",
			expectError:     false,
			outputContains: []string{
				"Don't forget to git push",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test git repository
			tempDir, cleanup := setupTestGitRepo(t, tt.initialTag, tt.setupBranch, tt.addUncommitted)
			defer cleanup()
			_ = tempDir // We don't use tempDir in tests, but keep it for future use

			// Create output buffer
			var output bytes.Buffer

			// Run the bump function
			err := Bump(tt.cfg, &output)

			// Check error expectations
			if tt.expectError {
				require.Error(t, err, "Expected an error but got none")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains,
						"Error message doesn't contain expected text")
				}
				return
			}

			// No error expected
			require.NoError(t, err, "Unexpected error: %v", err)

			// Check output contains expected strings
			outputStr := output.String()
			for _, expected := range tt.outputContains {
				assert.Contains(t, outputStr, expected,
					"Output doesn't contain expected text: %s", expected)
			}

			// Verify the git tag was created with expected version
			if tt.expectedVersion != "" {
				cmd := exec.Command("git", "tag", "-l", tt.expectedVersion)
				tagOutput, err := cmd.Output()
				require.NoError(t, err, "Failed to list git tags")
				assert.Contains(t, string(tagOutput), tt.expectedVersion,
					"Expected tag %s was not created", tt.expectedVersion)
			}

			// Verify changelog was updated
			changelogContent, err := os.ReadFile("CHANGELOG.md")
			require.NoError(t, err, "Failed to read changelog")
			changelogStr := string(changelogContent)

			// Should have the new version header (with the version as-is, including v prefix if present)
			if tt.expectedVersion != "" {
				assert.Contains(t, changelogStr, "["+tt.expectedVersion+"]",
					"Changelog doesn't contain version header for %s", tt.expectedVersion)
			}

			// Should still have unreleased section
			assert.Contains(t, changelogStr, "[Unreleased]",
				"Changelog doesn't contain Unreleased section")
		})
	}
}

func TestBump_GitNotInstalled(t *testing.T) {
	// This test checks behavior when git is not found
	// We can't actually uninstall git, but we can test the error path
	// by setting up an environment where git command would fail

	cfg := BumpConfig{
		BumpType:           "major",
		AllowAnyBranch:     false,
		AutoPush:           false,
		ChangelogFile:      "CHANGELOG.md",
		RepositoryProvider: "github",
		UseVPrefix:         true,
	}

	var output bytes.Buffer

	// Create a directory that's not a git repo
	tempDir, err := os.MkdirTemp("", "non-git-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Create a minimal changelog
	os.WriteFile("CHANGELOG.md", []byte("# Changelog\n\n## [Unreleased]\n"), 0o644)

	// This should fail because we're not in a git repo
	err = Bump(cfg, &output)
	require.Error(t, err, "Expected error when not in git repo")
}

func TestBump_InvalidChangelogFile(t *testing.T) {
	// Setup test git repository
	tempDir, cleanup := setupTestGitRepo(t, "v1.0.0", "main", false)
	defer cleanup()
	_ = tempDir // We don't use tempDir directly, but keep it for consistency

	cfg := BumpConfig{
		BumpType:           "major",
		AllowAnyBranch:     false,
		AutoPush:           false,
		ChangelogFile:      "NONEXISTENT.md", // File doesn't exist
		RepositoryProvider: "github",
		UseVPrefix:         true,
	}

	var output bytes.Buffer

	err := Bump(cfg, &output)
	require.Error(t, err, "Expected error with nonexistent changelog file")
	assert.Contains(t, err.Error(), "failed to update changelog",
		"Error should mention changelog update failure")
}

func TestBumpConfig(t *testing.T) {
	// Test that BumpConfig struct can be created with all fields
	cfg := BumpConfig{
		BumpType:           "major",
		AllowAnyBranch:     true,
		AutoPush:           true,
		ChangelogFile:      "CHANGELOG.md",
		RepositoryProvider: "gitlab",
		UseVPrefix:         false,
	}

	assert.Equal(t, "major", cfg.BumpType)
	assert.True(t, cfg.AllowAnyBranch)
	assert.True(t, cfg.AutoPush)
	assert.Equal(t, "CHANGELOG.md", cfg.ChangelogFile)
	assert.Equal(t, "gitlab", cfg.RepositoryProvider)
	assert.False(t, cfg.UseVPrefix)
}

// TestBump_Integration tests a complete workflow
func TestBump_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test git repository with initial state
	tempDir, cleanup := setupTestGitRepo(t, "v0.1.0", "main", false)
	defer cleanup()
	_ = tempDir // We don't use tempDir directly, but keep it for consistency

	// Perform a series of bumps
	bumps := []struct {
		bumpType        string
		expectedVersion string
	}{
		{"patch", "v0.1.1"},
		{"minor", "v0.2.0"},
		{"major", "v1.0.0"},
		{"patch", "v1.0.1"},
	}

	var output bytes.Buffer

	for _, bump := range bumps {
		output.Reset()

		cfg := BumpConfig{
			BumpType:           bump.bumpType,
			AllowAnyBranch:     false,
			AutoPush:           false,
			ChangelogFile:      "CHANGELOG.md",
			RepositoryProvider: "github",
			UseVPrefix:         true,
		}

		err := Bump(cfg, &output)
		require.NoError(t, err, "Bump %s failed", bump.bumpType)

		// Verify the version in output
		assert.Contains(t, output.String(), bump.expectedVersion,
			"Output should contain version %s", bump.expectedVersion)

		// Verify git tag exists
		cmd := exec.Command("git", "tag", "-l", bump.expectedVersion)
		tagOutput, err := cmd.Output()
		require.NoError(t, err, "Failed to list git tags")
		assert.Contains(t, string(tagOutput), bump.expectedVersion,
			"Tag %s should exist", bump.expectedVersion)
	}
}

func TestBump_AutoPush(t *testing.T) {
	// Create bare remote repo
	remoteDir, err := os.MkdirTemp("", "version-remote-*")
	require.NoError(t, err)
	defer os.RemoveAll(remoteDir)

	cmd := exec.Command("git", "init", "--bare")
	cmd.Dir = remoteDir
	require.NoError(t, cmd.Run())

	// Setup test repo with tag
	tempDir, cleanup := setupTestGitRepo(t, "v1.0.0", "main", false)
	defer cleanup()
	_ = tempDir

	// Add bare repo as remote
	cmd = exec.Command("git", "remote", "add", "origin", remoteDir)
	require.NoError(t, cmd.Run())

	// Push initial state to remote
	cmd = exec.Command("git", "push", "-u", "origin", "main")
	require.NoError(t, cmd.Run())
	cmd = exec.Command("git", "push", "--tags")
	require.NoError(t, cmd.Run())

	var output bytes.Buffer
	cfg := BumpConfig{
		BumpType:           "patch",
		AllowAnyBranch:     false,
		AutoPush:           true,
		ChangelogFile:      "CHANGELOG.md",
		RepositoryProvider: "github",
		UseVPrefix:         true,
	}

	err = Bump(cfg, &output)
	require.NoError(t, err)

	outputStr := output.String()
	assert.Contains(t, outputStr, "Pushing changes and tags")
	assert.Contains(t, outputStr, "Automatically pushed")
	assert.Contains(t, outputStr, "v1.0.1")

	// Verify the tag exists on the remote
	cmd = exec.Command("git", "ls-remote", "--tags", "origin", "v1.0.1")
	tagOutput, err := cmd.Output()
	require.NoError(t, err)
	assert.Contains(t, string(tagOutput), "v1.0.1")
}

func TestBump_AutoPushFailure(t *testing.T) {
	tempDir, cleanup := setupTestGitRepo(t, "v1.0.0", "main", false)
	defer cleanup()
	_ = tempDir

	// Add a non-existent remote
	cmd := exec.Command("git", "remote", "add", "origin", "/nonexistent/path")
	require.NoError(t, cmd.Run())

	var output bytes.Buffer
	cfg := BumpConfig{
		BumpType:           "patch",
		AllowAnyBranch:     false,
		AutoPush:           true,
		ChangelogFile:      "CHANGELOG.md",
		RepositoryProvider: "github",
		UseVPrefix:         true,
	}

	err := Bump(cfg, &output)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to push")
}

// TestBump_CommitChangelogFailure verifies that Bump returns an error containing
// "failed to commit changelog" when the git commit step cannot complete.
// The .git/index.lock file is used to simulate a git index lock, which causes
// git-add (the first step inside git.CommitChangelog) to fail.
func TestBump_CommitChangelogFailure(t *testing.T) {
	tempDir, cleanup := setupTestGitRepo(t, "v1.0.0", "main", false)
	defer cleanup()

	// Create .git/index.lock to prevent git from acquiring the index.
	// git refuses to run operations that modify the index when this file exists.
	lockFile := filepath.Join(tempDir, ".git", "index.lock")
	err := os.WriteFile(lockFile, []byte("locked"), 0o644)
	require.NoError(t, err)
	defer os.Remove(lockFile)

	cfg := BumpConfig{
		BumpType:           "patch",
		AllowAnyBranch:     false,
		AutoPush:           false,
		ChangelogFile:      "CHANGELOG.md",
		RepositoryProvider: "github",
		UseVPrefix:         true,
	}

	var out bytes.Buffer
	err = Bump(cfg, &out)
	require.Error(t, err, "Expected error when git commit step fails")
	assert.Contains(t, err.Error(), "failed to commit changelog")
}

func TestBump_TagVersionFailure(t *testing.T) {
	// Build a repo where v1.0.1 exists on an older commit and v1.0.0 is on HEAD,
	// so git describe reports v1.0.0 but Bump's attempt to create v1.0.1 fails.
	tempDir, err := os.MkdirTemp("", "version-bump-test-*")
	require.NoError(t, err)

	originalDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Warning: Failed to change back to original directory: %v", err)
		}
		os.RemoveAll(tempDir)
	}()

	// Git init and config
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()
	exec.Command("git", "config", "tag.gpgsign", "false").Run()
	exec.Command("git", "config", "commit.gpgsign", "false").Run()

	// First commit: tag it v1.0.1 (the tag Bump will try to create)
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "placeholder.txt"), []byte("placeholder"), 0o644))
	exec.Command("git", "add", "placeholder.txt").Run()
	exec.Command("git", "commit", "-m", "placeholder").Run()
	exec.Command("git", "tag", "-a", "v1.0.1", "-m", "pre-existing").Run()

	// Second commit: add changelog and tag it v1.0.0 (the "current" version)
	changelog := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- New feature for testing

`
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "CHANGELOG.md"), []byte(changelog), 0o644))
	exec.Command("git", "add", "CHANGELOG.md").Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()
	exec.Command("git", "tag", "-a", "v1.0.0", "-m", "v1.0.0").Run()

	var output bytes.Buffer
	cfg := BumpConfig{
		BumpType:           "patch",
		AllowAnyBranch:     false,
		AutoPush:           false,
		ChangelogFile:      "CHANGELOG.md",
		RepositoryProvider: "github",
		UseVPrefix:         true,
	}

	err = Bump(cfg, &output)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to tag version")
}
