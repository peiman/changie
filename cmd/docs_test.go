// cmd/docs_test.go

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

func TestRunDocsConfig(t *testing.T) {
	// Save original binary name
	originalBinaryName := binaryName
	defer func() { binaryName = originalBinaryName }()
	binaryName = "changie"

	tests := []struct {
		name         string
		format       string
		outputFile   string
		wantErr      bool
		wantContains string
	}{
		{
			name:         "markdown format to stdout",
			format:       "markdown",
			outputFile:   "",
			wantErr:      false,
			wantContains: "# Configuration Options",
		},
		{
			name:         "yaml format to stdout",
			format:       "yaml",
			outputFile:   "",
			wantErr:      false,
			wantContains: "app:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset viper
			viper.Reset()

			// Set up viper values
			viper.Set("app.docs.output_format", tc.format)
			viper.Set("app.docs.output_file", tc.outputFile)

			// Create command with output buffer
			cmd := &cobra.Command{}
			var outBuf bytes.Buffer
			cmd.SetOut(&outBuf)

			// Set flags
			cmd.Flags().String("format", tc.format, "")
			cmd.Flags().String("output", tc.outputFile, "")

			// Execute command
			err := runDocsConfig(cmd, []string{})

			// Check error
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				output := outBuf.String()
				if tc.wantContains != "" {
					assert.Contains(t, output, tc.wantContains)
				}
			}
		})
	}
}

func TestRunDocsConfigWithFile(t *testing.T) {
	// Save original binary name
	originalBinaryName := binaryName
	defer func() { binaryName = originalBinaryName }()
	binaryName = "changie"

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "docs-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	outputFile := filepath.Join(tempDir, "config-docs.md")

	// Reset viper
	viper.Reset()
	viper.Set("app.docs.output_format", "markdown")
	viper.Set("app.docs.output_file", outputFile)

	// Create command
	cmd := &cobra.Command{}
	cmd.Flags().String("format", "markdown", "")
	cmd.Flags().String("output", outputFile, "")

	// Mark output flag as changed
	cmd.Flags().Set("output", outputFile)

	// Execute command
	err = runDocsConfig(cmd, []string{})
	require.NoError(t, err)

	// Verify file was created
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "# Configuration Options")
}

func TestGetConfigValue(t *testing.T) {
	tests := []struct {
		name         string
		viperValue   interface{}
		flagValue    string
		flagChanged  bool
		expectedType string
		want         interface{}
	}{
		{
			name:         "string from viper",
			viperValue:   "viper-value",
			flagValue:    "",
			flagChanged:  false,
			expectedType: "string",
			want:         "viper-value",
		},
		{
			name:         "string from flag",
			viperValue:   "viper-value",
			flagValue:    "flag-value",
			flagChanged:  true,
			expectedType: "string",
			want:         "flag-value",
		},
		{
			name:         "empty viper value",
			viperValue:   nil,
			flagValue:    "",
			flagChanged:  false,
			expectedType: "string",
			want:         "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset viper
			viper.Reset()

			// Set viper value if provided
			if tc.viperValue != nil {
				viper.Set("test.key", tc.viperValue)
			}

			// Create command with flag
			cmd := &cobra.Command{}
			cmd.Flags().String("testflag", "", "test flag")

			// Set flag value if changed
			if tc.flagChanged {
				cmd.Flags().Set("testflag", tc.flagValue)
			}

			// Test string type
			if tc.expectedType == "string" {
				result := getConfigValue[string](cmd, "testflag", "test.key")
				assert.Equal(t, tc.want, result)
			}
		})
	}
}

func TestSetupCommandConfig(t *testing.T) {
	// Create a command
	cmd := &cobra.Command{
		Use: "test",
	}

	// Track if original PreRunE was called
	originalCalled := false
	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		originalCalled = true
		return nil
	}

	// Setup command config
	setupCommandConfig(cmd)

	// Verify PreRunE was replaced
	assert.NotNil(t, cmd.PreRunE)

	// Call the new PreRunE
	err := cmd.PreRunE(cmd, []string{})
	assert.NoError(t, err)
	assert.True(t, originalCalled, "Original PreRunE should be called")
}

func TestSetupCommandConfigNoOriginalPreRunE(t *testing.T) {
	// Create a command without PreRunE
	cmd := &cobra.Command{
		Use: "test",
	}

	// Setup command config
	setupCommandConfig(cmd)

	// Verify PreRunE was set
	assert.NotNil(t, cmd.PreRunE)

	// Call the new PreRunE
	err := cmd.PreRunE(cmd, []string{})
	assert.NoError(t, err)
}
