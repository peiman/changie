// cmd/output_json_test.go

package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/peiman/changie/.ckeletin/pkg/logger"
	"github.com/peiman/changie/.ckeletin/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validChangelogFixture is minimal valid changelog content used by output_json tests.
const validChangelogFixture = `# Changelog

## [Unreleased]

## [1.0.0] - 2024-01-01

### Added

- Initial release

[Unreleased]: https://example.com/compare/v1.0.0...HEAD
[1.0.0]: https://example.com/releases/tag/v1.0.0
`

// createTempChangelog writes validChangelogFixture to a temp file and returns its path.
func createTempChangelog(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "CHANGELOG.md")
	require.NoError(t, os.WriteFile(path, []byte(validChangelogFixture), 0o644))
	return path
}

func findSubcommand(root *cobra.Command, name string) *cobra.Command {
	for _, cmd := range root.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

// resetOutputJSONTestState resets all global state that integration tests modify.
// Must be called via defer at the start of each test.
func resetOutputJSONTestState(origCfgFile string, origStatus string, origUsed string) {
	output.SetOutputMode("")
	output.SetCommandName("")
	viper.Reset()
	cfgFile = origCfgFile
	configFileStatus = origStatus
	configFileUsed = origUsed
	RootCmd.SetArgs(nil)
	RootCmd.SetOut(nil)
	RootCmd.SetErr(nil)

	// Reset persistent flags to their default values.
	// Without this, Cobra retains flag values from a previous Execute() call,
	// causing subsequent tests to inherit prior test settings.
	resetFlags := map[string]string{"output": "text", "log-level": "info"}
	for name, def := range resetFlags {
		if f := RootCmd.PersistentFlags().Lookup(name); f != nil {
			f.Value.Set(def) //nolint:errcheck // resetting to known-good default
			f.Changed = false
		}
	}
}

func TestOutputJSON_ValidateCommand(t *testing.T) {
	savedLogger, savedLevel := logger.SaveLoggerState()
	defer logger.RestoreLoggerState(savedLogger, savedLevel)

	origCfgFile := cfgFile
	origStatus := configFileStatus
	origUsed := configFileUsed
	defer resetOutputJSONTestState(origCfgFile, origStatus, origUsed)

	tempFile := createTempChangelog(t)

	var stdout bytes.Buffer
	RootCmd.SetOut(&stdout)
	RootCmd.SetErr(&bytes.Buffer{})
	RootCmd.SetArgs([]string{"changelog", "validate", "--file", tempFile, "--output", "json"})

	err := RootCmd.Execute()
	require.NoError(t, err)

	var envelope output.JSONEnvelope
	err = json.Unmarshal(stdout.Bytes(), &envelope)
	require.NoError(t, err, "stdout should contain valid JSON, got: %s", stdout.String())

	assert.Equal(t, "success", envelope.Status)
	assert.Equal(t, "validate", envelope.Command)
	assert.Nil(t, envelope.Error)
	assert.NotNil(t, envelope.Data)
}

func TestOutputJSON_DefaultIsText(t *testing.T) {
	savedLogger, savedLevel := logger.SaveLoggerState()
	defer logger.RestoreLoggerState(savedLogger, savedLevel)

	origCfgFile := cfgFile
	origStatus := configFileStatus
	origUsed := configFileUsed
	defer resetOutputJSONTestState(origCfgFile, origStatus, origUsed)

	tempFile := createTempChangelog(t)

	var stdout bytes.Buffer
	RootCmd.SetOut(&stdout)
	RootCmd.SetErr(&bytes.Buffer{})
	RootCmd.SetArgs([]string{"changelog", "validate", "--file", tempFile})

	err := RootCmd.Execute()
	require.NoError(t, err)

	textOutput := stdout.String()
	assert.Contains(t, textOutput, "✔", "text mode should contain checkmark")

	var envelope output.JSONEnvelope
	err = json.Unmarshal(stdout.Bytes(), &envelope)
	assert.Error(t, err, "text mode output should not be valid JSON")
}

func TestOutputJSON_InvalidFormat(t *testing.T) {
	savedLogger, savedLevel := logger.SaveLoggerState()
	defer logger.RestoreLoggerState(savedLogger, savedLevel)

	origCfgFile := cfgFile
	origStatus := configFileStatus
	origUsed := configFileUsed
	defer resetOutputJSONTestState(origCfgFile, origStatus, origUsed)

	tempFile := createTempChangelog(t)

	var stdout, stderr bytes.Buffer
	RootCmd.SetOut(&stdout)
	RootCmd.SetErr(&stderr)
	RootCmd.SetArgs([]string{"changelog", "validate", "--file", tempFile, "--output", "xml"})

	err := RootCmd.Execute()
	assert.Error(t, err, "invalid output format should cause an error")
}

func TestOutputJSON_StderrSilent(t *testing.T) {
	savedLogger, savedLevel := logger.SaveLoggerState()
	defer logger.RestoreLoggerState(savedLogger, savedLevel)

	origCfgFile := cfgFile
	origStatus := configFileStatus
	origUsed := configFileUsed
	defer resetOutputJSONTestState(origCfgFile, origStatus, origUsed)

	tempFile := createTempChangelog(t)

	var stdout, stderr bytes.Buffer
	RootCmd.SetOut(&stdout)
	RootCmd.SetErr(&stderr)
	RootCmd.SetArgs([]string{"changelog", "validate", "--file", tempFile, "--output", "json"})

	err := RootCmd.Execute()
	require.NoError(t, err)

	// Stderr should be empty in JSON mode (zerolog disabled)
	assert.Empty(t, stderr.String(), "stderr should be empty in JSON mode")

	// Stdout should have valid JSON
	var envelope output.JSONEnvelope
	err = json.Unmarshal(stdout.Bytes(), &envelope)
	assert.NoError(t, err, "stdout should be valid JSON")
}

func TestOutputJSON_ErrorCommand(t *testing.T) {
	savedLogger, savedLevel := logger.SaveLoggerState()
	defer logger.RestoreLoggerState(savedLogger, savedLevel)

	origCfgFile := cfgFile
	origStatus := configFileStatus
	origUsed := configFileUsed
	defer resetOutputJSONTestState(origCfgFile, origStatus, origUsed)

	var stdout bytes.Buffer
	RootCmd.SetOut(&stdout)
	RootCmd.SetErr(&bytes.Buffer{})
	// Use a nonexistent file path to trigger an error from the validate command
	RootCmd.SetArgs([]string{"changelog", "validate", "--file", "/nonexistent/path/CHANGELOG.md", "--output", "json"})

	err := RootCmd.Execute()
	assert.Error(t, err, "nonexistent file should cause an error")

	// JSON mode should have been activated early (before config validation),
	// so main.go's error handler will emit the JSON error envelope.
	// The actual JSON envelope emission is tested in main_test.go (TestRun_JSONMode_Error).
	assert.True(t, output.IsJSONMode(), "JSON mode should be active even when command fails")
}

func TestOutputJSON_EnvelopeStructure(t *testing.T) {
	savedLogger, savedLevel := logger.SaveLoggerState()
	defer logger.RestoreLoggerState(savedLogger, savedLevel)

	origCfgFile := cfgFile
	origStatus := configFileStatus
	origUsed := configFileUsed
	defer resetOutputJSONTestState(origCfgFile, origStatus, origUsed)

	tempFile := createTempChangelog(t)

	var stdout bytes.Buffer
	RootCmd.SetOut(&stdout)
	RootCmd.SetErr(&bytes.Buffer{})
	RootCmd.SetArgs([]string{"changelog", "validate", "--file", tempFile, "--output", "json"})

	err := RootCmd.Execute()
	require.NoError(t, err)

	// Parse as raw JSON to check exact field presence
	var raw map[string]json.RawMessage
	err = json.Unmarshal(stdout.Bytes(), &raw)
	require.NoError(t, err)

	// All four fields should be present
	assert.Contains(t, raw, "status")
	assert.Contains(t, raw, "command")
	assert.Contains(t, raw, "data")
	assert.Contains(t, raw, "error")

	// error should be null on success
	assert.Equal(t, "null", string(raw["error"]))
}
