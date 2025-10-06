package version

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peiman/changie/internal/output"
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

func TestBumpWithJSONOutput(t *testing.T) {
	// Set up a buffer to capture log output
	var logBuf bytes.Buffer
	log.Logger = zerolog.New(&logBuf).With().Timestamp().Logger()

	tests := []struct {
		name        string
		bumpType    string
		expectedOld string
		expectedNew string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "successful patch bump with JSON",
			bumpType:    "patch",
			expectedOld: "v1.0.0",
			expectedNew: "v1.0.1",
			expectError: false,
		},
		{
			name:        "successful minor bump with JSON",
			bumpType:    "minor",
			expectedOld: "v1.0.0",
			expectedNew: "v1.1.0",
			expectError: false,
		},
		{
			name:        "successful major bump with JSON",
			bumpType:    "major",
			expectedOld: "v1.0.0",
			expectedNew: "v2.0.0",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a separate git repository for each test
			_, cleanup := setupTestGitRepo(t, "v1.0.0", "main", false)
			defer cleanup()
			// Enable JSON output
			viper.Set("app.json_output", true)
			defer viper.Set("app.json_output", false)

			var outputBuf bytes.Buffer

			cfg := BumpConfig{
				BumpType:           tt.bumpType,
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			}

			err := Bump(cfg, &outputBuf)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}

				// Verify JSON error output
				var jsonOutput output.BumpOutput
				jsonErr := json.Unmarshal(outputBuf.Bytes(), &jsonOutput)
				require.NoError(t, jsonErr, "Output should be valid JSON")
				assert.False(t, jsonOutput.Success, "Success should be false")
				assert.NotEmpty(t, jsonOutput.Error, "Error field should be populated")
				assert.Equal(t, tt.bumpType, jsonOutput.BumpType)
			} else {
				require.NoError(t, err)

				// Verify JSON success output
				var jsonOutput output.BumpOutput
				jsonErr := json.Unmarshal(outputBuf.Bytes(), &jsonOutput)
				require.NoError(t, jsonErr, "Output should be valid JSON")

				// Verify JSON fields
				assert.True(t, jsonOutput.Success, "Success should be true")
				assert.Empty(t, jsonOutput.Error, "Error field should be empty")
				assert.Equal(t, tt.bumpType, jsonOutput.BumpType)
				assert.Equal(t, tt.expectedOld, jsonOutput.OldVersion)
				assert.Equal(t, tt.expectedNew, jsonOutput.NewVersion)
				assert.Equal(t, tt.expectedNew, jsonOutput.Tag)
				assert.Equal(t, "CHANGELOG.md", jsonOutput.ChangelogFile)
				assert.False(t, jsonOutput.Pushed, "Pushed should be false when AutoPush is false")
			}
		})
	}
}

func TestBumpWithJSONOutputErrors(t *testing.T) {
	// Set up a buffer to capture log output
	var logBuf bytes.Buffer
	log.Logger = zerolog.New(&logBuf).With().Timestamp().Logger()

	tests := []struct {
		name        string
		setupFunc   func() func()
		cfg         BumpConfig
		errorSubstr string
	}{
		{
			name: "uncommitted changes error with JSON",
			setupFunc: func() func() {
				_, cleanup := setupTestGitRepo(t, "v1.0.0", "main", true)
				return cleanup
			},
			cfg: BumpConfig{
				BumpType:           "patch",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			errorSubstr: "uncommitted changes",
		},
		{
			name: "wrong branch error with JSON",
			setupFunc: func() func() {
				_, cleanup := setupTestGitRepo(t, "v1.0.0", "feature", false)
				return cleanup
			},
			cfg: BumpConfig{
				BumpType:           "patch",
				AllowAnyBranch:     false,
				AutoPush:           false,
				ChangelogFile:      "CHANGELOG.md",
				RepositoryProvider: "github",
				UseVPrefix:         true,
			},
			errorSubstr: "not on main/master branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupFunc()
			defer cleanup()

			// Enable JSON output
			viper.Set("app.json_output", true)
			defer viper.Set("app.json_output", false)

			var outputBuf bytes.Buffer
			err := Bump(tt.cfg, &outputBuf)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorSubstr)

			// Verify JSON error output
			var jsonOutput output.BumpOutput
			jsonErr := json.Unmarshal(outputBuf.Bytes(), &jsonOutput)
			require.NoError(t, jsonErr, "Output should be valid JSON")

			assert.False(t, jsonOutput.Success, "Success should be false")
			assert.NotEmpty(t, jsonOutput.Error, "Error field should be populated")
			assert.Contains(t, jsonOutput.Error, tt.errorSubstr, "Error should contain expected substring")
			assert.Equal(t, tt.cfg.BumpType, jsonOutput.BumpType)
			assert.Equal(t, tt.cfg.ChangelogFile, jsonOutput.ChangelogFile)
		})
	}
}
