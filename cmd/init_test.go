// cmd/init_test.go

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestRunInit(t *testing.T) {
	// Save original binary name
	originalBinaryName := binaryName
	defer func() { binaryName = originalBinaryName }()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "init-cmd-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up a buffer to capture output
	var outputBuf bytes.Buffer
	var logBuf bytes.Buffer

	// Configure logger to output to buffer for testing
	log.Logger = zerolog.New(&logBuf)

	// Test cases
	testCases := []struct {
		name      string
		args      []string
		tempFile  string
		createDir bool
		wantErr   bool
		wantMsg   string
	}{
		{
			name:     "Default file",
			args:     []string{},
			tempFile: "CHANGELOG.md",
			wantErr:  false,
			wantMsg:  "Project initialized with changelog file: CHANGELOG.md",
		},
		{
			name:     "Custom file name",
			args:     []string{"--file", "CUSTOM.md"},
			tempFile: "CUSTOM.md",
			wantErr:  false,
			wantMsg:  "Project initialized with changelog file: CUSTOM.md",
		},
		{
			name:     "File exists",
			args:     []string{"--file", "EXISTING.md"},
			tempFile: "EXISTING.md",
			wantErr:  true,
			wantMsg:  "failed to initialize project",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear buffers
			outputBuf.Reset()
			logBuf.Reset()
			viper.Reset()

			// Change to temp dir
			originalDir, _ := os.Getwd()
			err := os.Chdir(tempDir)
			if err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					t.Logf("Warning: Failed to change back to original directory: %v", err)
				}
			}()

			// Create existing file if needed
			if tc.name == "File exists" {
				fullPath := filepath.Join(tempDir, tc.tempFile)
				err := os.WriteFile(fullPath, []byte("existing content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create existing file: %v", err)
				}
			}

			// Create a new root command for each test
			rootCmd := &cobra.Command{Use: binaryName}
			cmd := &cobra.Command{
				Use:   "init",
				Short: "Initialize a project with SemVer and Keep a Changelog",
				RunE:  runInit,
			}

			// Set up the command flags
			cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
			rootCmd.AddCommand(cmd)
			rootCmd.SetArgs(append([]string{"init"}, tc.args...))
			rootCmd.SetOut(&outputBuf)
			rootCmd.SetErr(&outputBuf)
			rootCmd.SilenceUsage = true
			rootCmd.SilenceErrors = true

			// Execute the command
			err = rootCmd.Execute()

			// Check error result
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// Check command output
			gotOutput := outputBuf.String()
			if !tc.wantErr && tc.wantMsg != "" && len(gotOutput) == 0 {
				t.Errorf("Expected output to contain '%s', but got empty output", tc.wantMsg)
			} else if !tc.wantErr && tc.wantMsg != "" && !bytes.Contains(outputBuf.Bytes(), []byte(tc.wantMsg)) {
				t.Errorf("Expected output to contain '%s', but got '%s'", tc.wantMsg, gotOutput)
			}

			// Check file existence if not expecting error
			if !tc.wantErr {
				fullPath := filepath.Join(tempDir, tc.tempFile)
				_, err := os.Stat(fullPath)
				if err != nil {
					t.Errorf("Expected file to exist at %s, but got error: %v", fullPath, err)
				}
			}
		})
	}
}
