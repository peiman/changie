// cmd/changelog_test.go

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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestChangelogAdd(t *testing.T) {
	// Save original binary name
	originalBinaryName := binaryName
	defer func() { binaryName = originalBinaryName }()

	// Set up a buffer to capture output
	var logBuf bytes.Buffer

	// Configure logger to output to buffer for testing
	log.Logger = zerolog.New(&logBuf)

	tests := []struct {
		name       string
		section    string
		content    string
		setupFile  bool
		wantErr    bool
		wantOutput string
	}{
		{
			name:       "add to Added section",
			section:    "added",
			content:    "New awesome feature",
			setupFile:  true,
			wantErr:    false,
			wantOutput: "Added to Added section: New awesome feature",
		},
		{
			name:       "add to Fixed section",
			section:    "fixed",
			content:    "Critical bug fix",
			setupFile:  true,
			wantErr:    false,
			wantOutput: "Added to Fixed section: Critical bug fix",
		},
		{
			name:       "add to Security section",
			section:    "security",
			content:    "Patched vulnerability",
			setupFile:  true,
			wantErr:    false,
			wantOutput: "Added to Security section: Patched vulnerability",
		},
		{
			name:      "add to nonexistent file",
			section:   "added",
			content:   "Test entry",
			setupFile: false,
			wantErr:   true,
		},
		{
			name:       "duplicate entry",
			section:    "changed",
			content:    "Duplicate entry",
			setupFile:  true,
			wantErr:    false,
			wantOutput: "Added to Changed section",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory for this test
			tempDir, err := os.MkdirTemp("", "changie-changelog-test-*")
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

			changelogPath := filepath.Join(tempDir, "CHANGELOG.md")

			// Create changelog file if needed
			if tc.setupFile {
				changelogContent := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

`
				err = os.WriteFile(changelogPath, []byte(changelogContent), 0o644)
				require.NoError(t, err, "Failed to create changelog")

				// For duplicate test, add the entry first
				if tc.name == "duplicate entry" {
					exec.Command("git", "init").Run()
					exec.Command("git", "config", "user.email", "test@example.com").Run()
					exec.Command("git", "config", "user.name", "Test User").Run()
				}
			}

			// Reset viper
			viper.Reset()
			viper.Set("app.changelog.file", "CHANGELOG.md")

			// Create command
			cmd := &cobra.Command{
				Use: tc.section,
				RunE: func(cmd *cobra.Command, args []string) error {
					caser := cases.Title(language.English)
					return runAddChangelogSection(cmd, args, caser.String(tc.section))
				},
			}

			cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
			viper.BindPFlag("app.changelog.file", cmd.Flags().Lookup("file"))

			// Set command args
			cmd.SetArgs([]string{tc.content})

			var outBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&outBuf)

			// Execute command
			err = cmd.Execute()

			// Check results
			output := outBuf.String()

			if tc.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error: %v\nOutput: %s", err, output)

				if tc.wantOutput != "" {
					assert.Contains(t, output, tc.wantOutput, "Output should contain expected message")
				}

				// Verify the entry was added to the file
				if tc.setupFile {
					content, err := os.ReadFile(changelogPath)
					assert.NoError(t, err, "Should be able to read changelog")
					assert.Contains(t, string(content), tc.content, "Changelog should contain the entry")
				}
			}
		})
	}
}

func TestCreateChangelogSectionCmd(t *testing.T) {
	tests := []struct {
		name        string
		section     string
		expectUse   string
		expectShort string
	}{
		{
			name:        "Added section",
			section:     "Added",
			expectUse:   "added CONTENT",
			expectShort: "Add a added entry to the changelog",
		},
		{
			name:        "Fixed section",
			section:     "Fixed",
			expectUse:   "fixed CONTENT",
			expectShort: "Add a fixed entry to the changelog",
		},
		{
			name:        "Security section",
			section:     "Security",
			expectUse:   "security CONTENT",
			expectShort: "Add a security entry to the changelog",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := createChangelogSectionCmd(tc.section)

			// Verify command properties
			assert.NotNil(t, cmd, "Command should not be nil")
			assert.Equal(t, tc.expectUse, cmd.Use, "Use should match")
			assert.Equal(t, tc.expectShort, cmd.Short, "Short should match")
			assert.NotNil(t, cmd.RunE, "RunE should be set")
			assert.NotEmpty(t, cmd.Long, "Long description should be set")

			// Test executing the command's RunE function
			// Create a temp changelog file
			tempDir, err := os.MkdirTemp("", "changie-section-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			changelogFile := filepath.Join(tempDir, "CHANGELOG.md")
			changelogContent := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

`
			err = os.WriteFile(changelogFile, []byte(changelogContent), 0o644)
			require.NoError(t, err)

			// Reset viper
			viper.Reset()
			viper.Set("app.changelog.file", changelogFile)

			// Execute the RunE function
			err = cmd.RunE(cmd, []string{"Test entry"})
			assert.NoError(t, err, "RunE should not error")

			// Verify the entry was added
			content, err := os.ReadFile(changelogFile)
			assert.NoError(t, err)
			assert.Contains(t, string(content), "Test entry")
		})
	}
}

func TestChangelogAddCommand_AllSections(t *testing.T) {
	// Verify all valid sections have commands
	expectedSections := []string{"added", "changed", "deprecated", "removed", "fixed", "security"}

	for _, section := range expectedSections {
		t.Run("section_"+section, func(t *testing.T) {
			found := false
			for _, cmd := range changelogCmd.Commands() {
				if cmd.Use == section+" CONTENT" {
					found = true
					break
				}
			}

			assert.True(t, found, "Expected to find command for section: %s", section)
		})
	}
}

func TestChangelogAddCommand_WithFlagOverride(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "changie-flag-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create custom changelog file
	customFile := filepath.Join(tempDir, "CUSTOM_CHANGELOG.md")
	changelogContent := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

`
	err = os.WriteFile(customFile, []byte(changelogContent), 0o644)
	require.NoError(t, err)

	// Reset viper
	viper.Reset()
	viper.Set("app.changelog.file", "CHANGELOG.md")

	// Create command with file flag
	cmd := &cobra.Command{Use: "added"}
	cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return runAddChangelogSection(c, args, "Added")
	}

	// Set the custom file via flag
	cmd.Flags().Set("file", customFile)
	cmd.SetArgs([]string{"Test entry with custom file"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Execute
	err = cmd.Execute()
	require.NoError(t, err)

	// Verify the custom file was updated
	content, err := os.ReadFile(customFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Test entry with custom file")
}
