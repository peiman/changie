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
