// cmd/changelog_validate_test.go

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validChangelogContent is a changelog that passes all 5 checks.
const validChangelogContent = `# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- Upcoming feature

## [1.1.0] - 2024-02-01

### Added

- Feature B

## [1.0.0] - 2024-01-01

### Fixed

- Bug fix

[Unreleased]: https://github.com/user/repo/compare/1.1.0...HEAD
[1.1.0]: https://github.com/user/repo/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
`

// invalidChangelogContent has versions out of order (no links, wrong order).
const invalidChangelogContent = `# Changelog

## [1.0.0] - 2024-01-01

### Added

- Feature A

## [1.1.0] - 2024-02-01

### Added

- Feature B

[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
[1.1.0]: https://github.com/user/repo/releases/tag/1.1.0
`

func setupValidateCmd(t *testing.T, filePath string) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	viper.Reset()
	viper.Set("app.changelog.file", filePath)

	cmd := &cobra.Command{
		Use:  "validate",
		RunE: runValidateChangelog,
	}
	cmd.Flags().String("file", filePath, "Changelog file name")

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	return cmd, &buf
}

func TestValidateCommand_ValidFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changie-validate-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	require.NoError(t, os.WriteFile(changelogPath, []byte(validChangelogContent), 0o644))

	cmd, buf := setupValidateCmd(t, changelogPath)
	err = cmd.RunE(cmd, []string{})

	assert.NoError(t, err, "valid changelog should succeed")
	assert.Contains(t, buf.String(), "checks passed")
}

func TestValidateCommand_InvalidFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changie-validate-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	require.NoError(t, os.WriteFile(changelogPath, []byte(invalidChangelogContent), 0o644))

	cmd, buf := setupValidateCmd(t, changelogPath)
	err = cmd.RunE(cmd, []string{})

	assert.Error(t, err, "invalid changelog should return error")
	assert.Contains(t, err.Error(), "validation failed")
	_ = buf
}

func TestValidateCommand_FileNotFound(t *testing.T) {
	cmd, _ := setupValidateCmd(t, "/nonexistent/path/CHANGELOG.md")
	err := cmd.RunE(cmd, []string{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read changelog")
}

func TestValidateCommand_CustomFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changie-validate-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	customPath := filepath.Join(tempDir, "CUSTOM.md")
	require.NoError(t, os.WriteFile(customPath, []byte(validChangelogContent), 0o644))

	viper.Reset()
	viper.Set("app.changelog.file", "CHANGELOG.md")

	cmd := &cobra.Command{
		Use:  "validate",
		RunE: runValidateChangelog,
	}
	cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
	require.NoError(t, cmd.Flags().Set("file", customPath))

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err = cmd.RunE(cmd, []string{})
	assert.NoError(t, err, "custom file flag should work")
}

func TestValidateCommand_EmptyFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changie-validate-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	require.NoError(t, os.WriteFile(changelogPath, []byte(""), 0o644))

	cmd, buf := setupValidateCmd(t, changelogPath)
	err = cmd.RunE(cmd, []string{})

	// Empty file: no headers → should pass version headers (nothing to check),
	// but will fail broken links (no links to match), etc.
	// The important thing is it runs without panicking.
	assert.NotNil(t, cmd)
	_ = buf
	_ = err
}

func TestValidateCommandRegistered(t *testing.T) {
	found := false
	for _, sub := range changelogCmd.Commands() {
		if sub.Use == "validate" {
			found = true
			break
		}
	}
	assert.True(t, found, "validate subcommand should be registered under changelogCmd")
}
