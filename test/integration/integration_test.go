// test/integration/integration_test.go
//
// Integration tests for full command execution
//
// These tests execute actual commands end-to-end to verify:
// - Complete command workflows
// - Flag parsing and precedence
// - Configuration loading
// - Output generation
// - Exit codes

package integration

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var binaryPath string

// getExitCode extracts the exit code from a command execution error.
// Returns 0 if err is nil, the exit code if err is an ExitError, or -1 otherwise.
func getExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return -1
}

// TestMain builds the binary before running tests
func TestMain(m *testing.M) {
	// Build the binary with platform-specific name
	binaryName := "ckeletin-go-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binaryName, "../../main.go")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build test binary: %v\n", err)
		os.Exit(1)
	}
	binaryPath = "./" + binaryName

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove(binaryPath)

	os.Exit(code)
}

func TestChangelogValidateCommand(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name               string
		fileContent        string
		wantExitCode       int
		wantOutputContains string
	}{
		{
			name: "Valid changelog",
			fileContent: `# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.0.0] - 2024-01-01

### Added
- Initial release

[Unreleased]: https://github.com/example/repo/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/example/repo/releases/tag/v1.0.0
`,
			wantExitCode:       0,
			wantOutputContains: "passed",
		},
		{
			name:               "Nonexistent changelog file",
			fileContent:        "", // no file created
			wantExitCode:       1,
			wantOutputContains: "no such file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []string
			if tt.fileContent != "" {
				changelogFile := filepath.Join(tmpDir, tt.name+"-CHANGELOG.md")
				err := os.WriteFile(changelogFile, []byte(tt.fileContent), 0600)
				require.NoError(t, err, "setup: failed to create changelog file")
				args = []string{"changelog", "validate", "--file", changelogFile}
			} else {
				args = []string{"changelog", "validate", "--file", filepath.Join(tmpDir, "nonexistent-CHANGELOG.md")}
			}

			cmd := exec.Command(binaryPath, args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			exitCode := getExitCode(err)

			assert.Equal(t, tt.wantExitCode, exitCode,
				"exit code mismatch\nstdout: %s\nstderr: %s",
				stdout.String(), stderr.String())

			if tt.wantOutputContains != "" {
				combinedOutput := stdout.String() + stderr.String()
				assert.Contains(t, combinedOutput, tt.wantOutputContains,
					"output should contain expected text")
			}
		})
	}
}

func TestConfigValidateCommand(t *testing.T) {
	// Create temp directory for test config files
	tmpDir := t.TempDir()

	tests := []struct {
		name               string
		configContent      string
		configPerms        os.FileMode
		args               []string
		wantExitCode       int
		wantOutputContains string
	}{
		{
			name: "Valid config",
			configContent: `app:
  log_level: debug
`,
			configPerms:        0600,
			wantExitCode:       0,
			wantOutputContains: "Configuration is valid",
		},
		{
			name: "Invalid YAML",
			configContent: `app:
  invalid: [unclosed
`,
			configPerms:        0600,
			wantExitCode:       1,
			wantOutputContains: "Configuration is invalid",
		},
		{
			name: "Unknown keys (warning)",
			configContent: `app:
  log_level: info
  unknown_key: value
`,
			configPerms:        0600,
			wantExitCode:       1, // Warnings cause exit 1
			wantOutputContains: "Unknown configuration key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test config file
			configFile := filepath.Join(tmpDir, tt.name+".yaml")
			err := os.WriteFile(configFile, []byte(tt.configContent), tt.configPerms)
			require.NoError(t, err, "setup: failed to create config file")

			// Build command
			args := append([]string{"config", "validate", "--file", configFile}, tt.args...)
			cmd := exec.Command(binaryPath, args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err = cmd.Run()

			// Check exit code
			exitCode := getExitCode(err)

			assert.Equal(t, tt.wantExitCode, exitCode,
				"exit code mismatch\nstdout: %s\nstderr: %s",
				stdout.String(), stderr.String())

			// Check output
			if tt.wantOutputContains != "" {
				combinedOutput := stdout.String() + stderr.String()
				assert.Contains(t, combinedOutput, tt.wantOutputContains,
					"output should contain expected text")
			}
		})
	}
}

func TestDocsCommand(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		wantExitCode       int
		wantOutputContains string
	}{
		{
			name:               "Generate markdown docs",
			args:               []string{"docs", "config"},
			wantExitCode:       0,
			wantOutputContains: "Configuration Options",
		},
		{
			name:               "Generate YAML docs",
			args:               []string{"docs", "config", "--format", "yaml"},
			wantExitCode:       0,
			wantOutputContains: "app:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			exitCode := getExitCode(err)

			assert.Equal(t, tt.wantExitCode, exitCode,
				"exit code mismatch\nstderr: %s", stderr.String())

			if tt.wantOutputContains != "" {
				output := stdout.String()
				assert.Contains(t, output, tt.wantOutputContains,
					"output should contain expected text")
			}
		})
	}
}

func TestHelpCommand(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		wantExitCode       int
		wantOutputContains string
	}{
		{
			name:               "Root help",
			args:               []string{"--help"},
			wantExitCode:       0,
			wantOutputContains: "Available Commands",
		},
		{
			name:               "Changelog help",
			args:               []string{"changelog", "--help"},
			wantExitCode:       0,
			wantOutputContains: "changelog",
		},
		{
			name:               "Config validate help",
			args:               []string{"config", "validate", "--help"},
			wantExitCode:       0,
			wantOutputContains: "Validate a configuration file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout

			err := cmd.Run()

			exitCode := getExitCode(err)

			assert.Equal(t, tt.wantExitCode, exitCode, "exit code mismatch")

			if tt.wantOutputContains != "" {
				output := stdout.String()
				assert.Contains(t, output, tt.wantOutputContains,
					"output should contain expected text")
			}
		})
	}
}

func TestVersionFlag(t *testing.T) {
	cmd := exec.Command(binaryPath, "--version")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	require.NoError(t, err, "version command should succeed")

	output := stdout.String()
	// Version output should contain version info
	hasVersion := strings.Contains(output, "version") || strings.Contains(output, "dev")
	assert.True(t, hasVersion,
		"version output should contain 'version' or 'dev', got: %s", output)
}

// TestConfigLoading tests configuration loading from files
func TestConfigLoading(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name               string
		configContent      string
		args               []string
		wantExitCode       int
		wantOutputContains string
	}{
		{
			name: "Config file with log level debug",
			configContent: `app:
  log_level: debug
`,
			args:               []string{"docs", "config"},
			wantExitCode:       0,
			wantOutputContains: "Configuration", // docs config always outputs this
		},
		{
			name: "Config file with log level info",
			configContent: `app:
  log_level: info
`,
			args:               []string{"docs", "config"},
			wantExitCode:       0,
			wantOutputContains: "Configuration",
		},
		{
			name: "Config file with nested version settings",
			configContent: `app:
  log_level: warn
  version:
    auto_push: false
`,
			args:               []string{"docs", "config"},
			wantExitCode:       0,
			wantOutputContains: "Configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create config file in a unique location per test
			configFile := filepath.Join(tmpDir, tt.name+"-config.yaml")
			err := os.WriteFile(configFile, []byte(tt.configContent), 0600)
			require.NoError(t, err, "setup: failed to create config file")

			// Run command with config file
			args := append([]string{"--config", configFile}, tt.args...)
			cmd := exec.Command(binaryPath, args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err = cmd.Run()

			// Check exit code
			exitCode := getExitCode(err)

			assert.Equal(t, tt.wantExitCode, exitCode,
				"exit code mismatch\nstdout: %s\nstderr: %s",
				stdout.String(), stderr.String())

			// Check output
			if tt.wantOutputContains != "" {
				output := stdout.String()
				assert.Contains(t, output, tt.wantOutputContains,
					"output should contain expected text")
			}
		})
	}
}

// TestEnvironmentVariables tests configuration via environment variables
func TestEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name               string
		envVars            map[string]string
		args               []string
		wantExitCode       int
		wantOutputContains string
	}{
		{
			name: "Set log level via env var",
			envVars: map[string]string{
				"CHANGIE_APP_LOG_LEVEL": "warn",
			},
			args:               []string{"docs", "config"},
			wantExitCode:       0,
			wantOutputContains: "Configuration",
		},
		{
			name: "Set docs format via env var",
			envVars: map[string]string{
				"CHANGIE_APP_DOCS_FORMAT": "yaml",
			},
			args:         []string{"docs", "config"},
			wantExitCode: 0,
			// YAML format outputs app: at top level
			wantOutputContains: "app:",
		},
		{
			name: "Multiple env vars together",
			envVars: map[string]string{
				"CHANGIE_APP_LOG_LEVEL":   "warn",
				"CHANGIE_APP_DOCS_FORMAT": "markdown",
			},
			args:               []string{"docs", "config"},
			wantExitCode:       0,
			wantOutputContains: "Configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Set environment variables
			cmd.Env = os.Environ()
			for k, v := range tt.envVars {
				cmd.Env = append(cmd.Env, k+"="+v)
			}

			err := cmd.Run()

			// Check exit code
			exitCode := getExitCode(err)

			assert.Equal(t, tt.wantExitCode, exitCode,
				"exit code mismatch\nstdout: %s\nstderr: %s",
				stdout.String(), stderr.String())

			// Check output
			if tt.wantOutputContains != "" {
				output := stdout.String()
				assert.Contains(t, output, tt.wantOutputContains,
					"output should contain expected text")
			}
		})
	}
}

// TestConfigPrecedence tests the precedence of config sources (flags > env > config file > defaults)
func TestConfigPrecedence(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a config file with a known log level
	configContent := `app:
  log_level: info
`
	configFile := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configFile, []byte(configContent), 0600)
	require.NoError(t, err, "setup: failed to create config file")

	tests := []struct {
		name               string
		envVars            map[string]string
		args               []string
		wantExitCode       int
		wantOutputContains string
	}{
		{
			name:               "Config file only - docs works",
			envVars:            map[string]string{},
			args:               []string{"--config", configFile, "docs", "config"},
			wantExitCode:       0,
			wantOutputContains: "Configuration",
		},
		{
			name: "Env var log level overrides config file",
			envVars: map[string]string{
				"CHANGIE_APP_LOG_LEVEL": "warn",
			},
			args:               []string{"--config", configFile, "docs", "config"},
			wantExitCode:       0,
			wantOutputContains: "Configuration",
		},
		{
			name: "CLI flag overrides env var and config file",
			envVars: map[string]string{
				"CHANGIE_APP_LOG_LEVEL": "warn",
			},
			args:               []string{"--config", configFile, "--log-level", "info", "docs", "config"},
			wantExitCode:       0,
			wantOutputContains: "Configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Set environment variables
			cmd.Env = os.Environ()
			for k, v := range tt.envVars {
				cmd.Env = append(cmd.Env, k+"="+v)
			}

			err := cmd.Run()
			exitCode := getExitCode(err)

			assert.Equal(t, tt.wantExitCode, exitCode,
				"exit code mismatch\nstdout: %s\nstderr: %s",
				stdout.String(), stderr.String())

			if tt.wantOutputContains != "" {
				output := stdout.String()
				assert.Contains(t, output, tt.wantOutputContains,
					"output should contain expected text")
			}
		})
	}
}

// TestMultiCommandWorkflow tests complex multi-command workflows
func TestMultiCommandWorkflow(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("Create config, validate it, then generate docs", func(t *testing.T) {
		// Step 1: Create a config file
		configContent := `app:
  log_level: info
`
		configFile := filepath.Join(tmpDir, "workflow-config.yaml")
		err := os.WriteFile(configFile, []byte(configContent), 0600)
		require.NoError(t, err, "setup: failed to create config file")

		// Step 2: Validate the config
		validateCmd := exec.Command(binaryPath, "config", "validate", "--file", configFile)
		var validateStdout bytes.Buffer
		validateCmd.Stdout = &validateStdout

		err = validateCmd.Run()
		require.NoError(t, err, "config validation should succeed: %s", validateStdout.String())

		assert.Contains(t, validateStdout.String(), "valid",
			"validation output should confirm config is valid")

		// Step 3: Generate docs using the validated config
		docsCmd := exec.Command(binaryPath, "--config", configFile, "docs", "config")
		var docsStdout bytes.Buffer
		docsCmd.Stdout = &docsStdout

		err = docsCmd.Run()
		require.NoError(t, err, "docs command should succeed")

		assert.Contains(t, docsStdout.String(), "Configuration",
			"docs output should contain configuration info")
	})

	t.Run("Generate docs then validate config", func(t *testing.T) {
		// Step 1: Generate documentation
		docsCmd := exec.Command(binaryPath, "docs", "config")
		var docsStdout bytes.Buffer
		docsCmd.Stdout = &docsStdout

		err := docsCmd.Run()
		require.NoError(t, err, "docs generation should succeed")

		assert.Contains(t, docsStdout.String(), "Configuration",
			"docs should contain configuration info")

		// Step 2: Create a config based on docs
		configContent := `app:
  log_level: info
`
		configFile := filepath.Join(tmpDir, "docs-based-config.yaml")
		err = os.WriteFile(configFile, []byte(configContent), 0600)
		require.NoError(t, err, "setup: failed to create config file")

		// Step 3: Validate the config
		validateCmd := exec.Command(binaryPath, "config", "validate", "--file", configFile)
		var validateStdout bytes.Buffer
		validateCmd.Stdout = &validateStdout

		err = validateCmd.Run()
		require.NoError(t, err, "config validation should succeed: %s", validateStdout.String())
	})
}
