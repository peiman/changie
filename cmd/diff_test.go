// cmd/diff_test.go

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

	"github.com/peiman/changie/.ckeletin/pkg/output"
)

// diffTestChangelog is a minimal valid changelog for diff command tests.
const diffTestChangelog = `# Changelog

## [1.1.0] - 2024-02-01

### Added

- Feature B

## [1.0.0] - 2024-01-01

### Fixed

- Bug fix A

[1.1.0]: https://github.com/user/repo/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
`

func setupDiffCmd(t *testing.T, filePath string) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	viper.Reset()
	viper.Set("app.changelog.file", filePath)

	cmd := &cobra.Command{
		Use:  "diff",
		RunE: runDiff,
	}
	cmd.Flags().String("file", filePath, "Changelog file name")

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	return cmd, &buf
}

func TestDiffCommand_HappyPath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changie-diff-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	require.NoError(t, os.WriteFile(changelogPath, []byte(diffTestChangelog), 0o644))

	cmd, buf := setupDiffCmd(t, changelogPath)
	err = cmd.RunE(cmd, []string{"1.0.0", "1.1.0"})

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "1.1.0")
	assert.Contains(t, buf.String(), "Feature B")
}

func TestDiffCommand_FileNotFound(t *testing.T) {
	cmd, _ := setupDiffCmd(t, "/nonexistent/path/CHANGELOG.md")
	err := cmd.RunE(cmd, []string{"1.0.0", "1.1.0"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read changelog")
}

func TestDiffCommand_VersionNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changie-diff-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	require.NoError(t, os.WriteFile(changelogPath, []byte(diffTestChangelog), 0o644))

	cmd, _ := setupDiffCmd(t, changelogPath)
	err = cmd.RunE(cmd, []string{"1.0.0", "9.9.9"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "9.9.9")
}

func TestDiffCommand_InvertedVersions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changie-diff-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	require.NoError(t, os.WriteFile(changelogPath, []byte(diffTestChangelog), 0o644))

	cmd, _ := setupDiffCmd(t, changelogPath)
	err = cmd.RunE(cmd, []string{"1.1.0", "1.0.0"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "older")
}

func TestDiffCommand_CustomFileFlag(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changie-diff-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	customPath := filepath.Join(tempDir, "HISTORY.md")
	require.NoError(t, os.WriteFile(customPath, []byte(diffTestChangelog), 0o644))

	viper.Reset()
	viper.Set("app.changelog.file", "CHANGELOG.md")

	cmd := &cobra.Command{
		Use:  "diff",
		RunE: runDiff,
	}
	cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
	require.NoError(t, cmd.Flags().Set("file", customPath))

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err = cmd.RunE(cmd, []string{"1.0.0", "1.1.0"})
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Feature B")
}

func TestDiffCommandRegistered(t *testing.T) {
	found := false
	for _, sub := range RootCmd.Commands() {
		if sub.Use == "diff FROM TO" {
			found = true
			break
		}
	}
	assert.True(t, found, "diff command should be registered on RootCmd")
}

func TestDiffCommand_JSONOutput(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "changie-diff-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	require.NoError(t, os.WriteFile(changelogPath, []byte(diffTestChangelog), 0o644))

	output.SetOutputMode("json")
	output.SetCommandName("diff")
	defer output.SetOutputMode("")

	cmd, buf := setupDiffCmd(t, changelogPath)
	err = cmd.RunE(cmd, []string{"1.0.0", "1.1.0"})

	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, `"status": "success"`)
	assert.Contains(t, out, `"command": "diff"`)
	assert.Contains(t, out, "Feature B")
}
