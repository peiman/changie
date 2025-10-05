// cmd/version_test.go

package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionBumpCommands(t *testing.T) {
	// Save original binary name
	originalBinaryName := binaryName
	defer func() { binaryName = originalBinaryName }()

	// Set up a buffer to capture output
	var logBuf bytes.Buffer

	// Configure logger to output to buffer for testing
	log.Logger = zerolog.New(&logBuf)

	tests := []struct {
		name            string
		bumpType        string
		setupBranch     string // Branch to create and checkout
		uncommitted     bool   // Whether to leave uncommitted changes
		allowAnyBranch  bool   // --allow-any-branch flag
		autoPush        bool   // --auto-push flag
		initialTag      string // Initial tag to create
		expectedVersion string // Expected version after bump
		wantErr         bool   // Whether we expect an error
		errContains     string // String that should be in error message
		wantMsg         string // String that should be in output
	}{
		{
			name:            "major bump on main branch succeeds",
			bumpType:        "major",
			setupBranch:     "main",
			uncommitted:     false,
			allowAnyBranch:  false,
			autoPush:        false,
			initialTag:      "v1.2.3",
			expectedVersion: "v2.0.0",
			wantErr:         false,
			wantMsg:         "major release v2.0.0 done",
		},
		{
			name:            "minor bump on main branch succeeds",
			bumpType:        "minor",
			setupBranch:     "main",
			uncommitted:     false,
			allowAnyBranch:  false,
			autoPush:        false,
			initialTag:      "v1.2.3",
			expectedVersion: "v1.3.0",
			wantErr:         false,
			wantMsg:         "minor release v1.3.0 done",
		},
		{
			name:            "patch bump on main branch succeeds",
			bumpType:        "patch",
			setupBranch:     "main",
			uncommitted:     false,
			allowAnyBranch:  false,
			autoPush:        false,
			initialTag:      "v1.2.3",
			expectedVersion: "v1.2.4",
			wantErr:         false,
			wantMsg:         "patch release v1.2.4 done",
		},
		{
			name:           "bump on feature branch fails without flag",
			bumpType:       "minor",
			setupBranch:    "feature/test",
			uncommitted:    false,
			allowAnyBranch: false,
			autoPush:       false,
			initialTag:     "v1.2.3",
			wantErr:        true,
			errContains:    "not on main/master branch",
		},
		{
			name:            "bump on feature branch succeeds with --allow-any-branch",
			bumpType:        "minor",
			setupBranch:     "feature/test",
			uncommitted:     false,
			allowAnyBranch:  true,
			autoPush:        false,
			initialTag:      "v1.2.3",
			expectedVersion: "v1.3.0",
			wantErr:         false,
			wantMsg:         "minor release v1.3.0 done",
		},
		{
			name:           "bump with uncommitted changes fails",
			bumpType:       "patch",
			setupBranch:    "main",
			uncommitted:    true,
			allowAnyBranch: false,
			autoPush:       false,
			initialTag:     "v1.2.3",
			wantErr:        true,
			errContains:    "uncommitted changes found",
		},
		{
			name:            "bump clears prerelease and build metadata",
			bumpType:        "patch",
			setupBranch:     "main",
			uncommitted:     false,
			allowAnyBranch:  false,
			autoPush:        false,
			initialTag:      "v1.2.3-alpha.1+build.123",
			expectedVersion: "v1.2.4",
			wantErr:         false,
			wantMsg:         "patch release v1.2.4 done",
		},
		{
			name:            "bump on master branch succeeds",
			bumpType:        "minor",
			setupBranch:     "master",
			uncommitted:     false,
			allowAnyBranch:  false,
			autoPush:        false,
			initialTag:      "v0.1.0",
			expectedVersion: "v0.2.0",
			wantErr:         false,
			wantMsg:         "minor release v0.2.0 done",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory for this test
			tempDir, err := os.MkdirTemp("", "changie-version-test-*")
			require.NoError(t, err, "Failed to create temp dir")
			defer os.RemoveAll(tempDir)

			// Save current directory
			originalDir, err := os.Getwd()
			require.NoError(t, err, "Failed to get current dir")

			// Change to temp directory
			err = os.Chdir(tempDir)
			require.NoError(t, err, "Failed to change to temp dir")
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					t.Logf("Warning: Failed to change back to original directory: %v", err)
				}
			}()

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

			// Create initial tag
			exec.Command("git", "tag", tc.initialTag).Run()

			// Checkout the desired branch
			if tc.setupBranch != "main" && tc.setupBranch != "master" {
				// Create and checkout new branch
				exec.Command("git", "checkout", "-b", tc.setupBranch).Run()
			} else if tc.setupBranch == "master" {
				// Rename main to master
				exec.Command("git", "branch", "-m", "main", "master").Run()
			}

			// Add uncommitted changes if needed
			if tc.uncommitted {
				testFile := filepath.Join(tempDir, "test.txt")
				os.WriteFile(testFile, []byte("uncommitted"), 0o644)
			}

			// Reset viper for clean state
			viper.Reset()
			viper.Set("app.version.use_v_prefix", true)
			viper.Set("app.changelog.file", "CHANGELOG.md")
			viper.Set("app.changelog.repository_provider", "github")

			// Create a new root command for this test
			rootCmd := &cobra.Command{Use: binaryName}
			var cmdToTest *cobra.Command

			switch tc.bumpType {
			case "major":
				cmdToTest = &cobra.Command{
					Use: "major",
					RunE: func(cmd *cobra.Command, args []string) error {
						return runVersionBump(cmd, "major")
					},
				}
			case "minor":
				cmdToTest = &cobra.Command{
					Use: "minor",
					RunE: func(cmd *cobra.Command, args []string) error {
						return runVersionBump(cmd, "minor")
					},
				}
			case "patch":
				cmdToTest = &cobra.Command{
					Use: "patch",
					RunE: func(cmd *cobra.Command, args []string) error {
						return runVersionBump(cmd, "patch")
					},
				}
			}

			// Add flags
			cmdToTest.Flags().String("file", "CHANGELOG.md", "Changelog file name")
			cmdToTest.Flags().String("rrp", "github", "Remote repository provider")
			cmdToTest.Flags().Bool("auto-push", false, "Automatically push changes and tags")
			cmdToTest.Flags().Bool("allow-any-branch", false, "Allow version bumping on any branch")

			// Bind flags to viper
			viper.BindPFlag("app.changelog.file", cmdToTest.Flags().Lookup("file"))
			viper.BindPFlag("app.changelog.repository_provider", cmdToTest.Flags().Lookup("rrp"))
			viper.BindPFlag("app.changelog.auto_push", cmdToTest.Flags().Lookup("auto-push"))
			viper.BindPFlag("app.version.allow_any_branch", cmdToTest.Flags().Lookup("allow-any-branch"))

			rootCmd.AddCommand(cmdToTest)

			// Set flags
			if tc.allowAnyBranch {
				cmdToTest.Flags().Set("allow-any-branch", "true")
			}
			if tc.autoPush {
				cmdToTest.Flags().Set("auto-push", "true")
			}

			// Capture output
			var outBuf bytes.Buffer
			cmdToTest.SetOut(&outBuf)
			cmdToTest.SetErr(&outBuf)

			// Execute command directly (not through parent)
			err = runVersionBump(cmdToTest, tc.bumpType)

			// Check results
			output := outBuf.String()

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				} else if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains,
						"Error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Unexpected error: %v\nOutput: %s", err, output)

				// Verify the expected message is in output
				if tc.wantMsg != "" {
					assert.Contains(t, output, tc.wantMsg,
						"Output should contain expected message")
				}

				// Verify the new tag was created
				if tc.expectedVersion != "" {
					tagCmd := exec.Command("git", "tag", "-l", tc.expectedVersion)
					tagCmd.Dir = tempDir
					tagOutput, err := tagCmd.CombinedOutput()
					assert.NoError(t, err, "Should be able to list tags")
					assert.Contains(t, string(tagOutput), tc.expectedVersion,
						"Expected tag %s should exist", tc.expectedVersion)
				}

				// Verify changelog was updated
				changelogContent, err := os.ReadFile(changelogPath)
				assert.NoError(t, err, "Should be able to read changelog")
				if tc.expectedVersion != "" {
					assert.Contains(t, string(changelogContent), tc.expectedVersion,
						"Changelog should contain new version")
				}
			}
		})
	}
}

func TestVersionBumpWithoutGit(t *testing.T) {
	// Save original binary name
	originalBinaryName := binaryName
	defer func() { binaryName = originalBinaryName }()

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "changie-nogit-test-*")
	require.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tempDir)

	// Save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current dir")

	// Change to temp directory (no git init)
	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change to temp dir")
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Warning: Failed to change back to original directory: %v", err)
		}
	}()

	// Reset viper
	viper.Reset()
	viper.Set("app.version.use_v_prefix", true)

	// Create command
	cmd := &cobra.Command{
		Use: "minor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersionBump(cmd, "minor")
		},
	}

	cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
	cmd.Flags().Bool("allow-any-branch", false, "Allow version bumping on any branch")

	var outBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&outBuf)

	// Execute - should fail because we're not in a git repo
	err = runVersionBump(cmd, "minor")
	assert.Error(t, err, "Should fail when not in a git repository")
	assert.Contains(t, err.Error(), "failed to get current branch",
		"Error should mention git issue")
}

func TestVersionBumpFlagPrecedence(t *testing.T) {
	// Test that command-line flags take precedence over viper config

	// Save original binary name
	originalBinaryName := binaryName
	defer func() { binaryName = originalBinaryName }()

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "changie-flag-test-*")
	require.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tempDir)

	// Save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current dir")

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change to temp dir")
	defer os.Chdir(originalDir)

	// Initialize git repository
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	// Create changelog
	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	changelog := `# Changelog

## [Unreleased]

### Added
- Test feature
`
	os.WriteFile(changelogPath, []byte(changelog), 0o644)
	exec.Command("git", "add", "CHANGELOG.md").Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()
	exec.Command("git", "tag", "v1.0.0").Run()

	// Create feature branch
	exec.Command("git", "checkout", "-b", "feature/test").Run()

	// Set viper config to disallow any branch
	viper.Reset()
	viper.Set("app.version.allow_any_branch", false)
	viper.Set("app.version.use_v_prefix", true)
	viper.Set("app.changelog.file", "CHANGELOG.md")

	// Create command
	cmd := &cobra.Command{
		Use: "minor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersionBump(cmd, "minor")
		},
	}

	cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
	cmd.Flags().Bool("allow-any-branch", false, "Allow version bumping on any branch")
	viper.BindPFlag("app.version.allow_any_branch", cmd.Flags().Lookup("allow-any-branch"))

	var outBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&outBuf)

	// Test 1: Without flag, should fail on feature branch
	err = runVersionBump(cmd, "minor")
	assert.Error(t, err, "Should fail without --allow-any-branch flag")
	assert.Contains(t, err.Error(), "not on main/master branch")

	// Reset command
	cmd = &cobra.Command{
		Use: "minor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersionBump(cmd, "minor")
		},
	}
	cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
	cmd.Flags().Bool("allow-any-branch", false, "Allow version bumping on any branch")
	viper.BindPFlag("app.version.allow_any_branch", cmd.Flags().Lookup("allow-any-branch"))

	outBuf.Reset()
	cmd.SetOut(&outBuf)
	cmd.SetErr(&outBuf)

	// Test 2: With flag, should succeed
	cmd.Flags().Set("allow-any-branch", "true")
	err = runVersionBump(cmd, "minor")
	assert.NoError(t, err, "Should succeed with --allow-any-branch flag")

	output := outBuf.String()
	assert.Contains(t, output, "v1.1.0", "Should bump to v1.1.0")
}
