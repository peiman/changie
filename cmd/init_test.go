// cmd/init_test.go

package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRunInit(t *testing.T) {
	// Save original binary name
	originalBinaryName := binaryName
	defer func() { binaryName = originalBinaryName }()

	// Set up a buffer to capture output
	var logBuf bytes.Buffer

	// Configure logger to output to buffer for testing
	log.Logger = zerolog.New(&logBuf)

	tests := []struct {
		name       string
		args       []string
		wantMsg    string
		wantErr    bool
		mockGit    bool
		mockTags   bool
		hasVPrefix bool
	}{
		{
			name:    "default init",
			args:    []string{},
			wantMsg: "Project initialized with changelog file: CHANGELOG.md",
			wantErr: false,
		},
		{
			name:    "custom file name",
			args:    []string{"--file", "CUSTOM.md"},
			wantMsg: "Project initialized with changelog file: CUSTOM.md",
			wantErr: false,
		},
		{
			name:    "explicit v prefix",
			args:    []string{"--use-v-prefix", "true"},
			wantMsg: "Version tags will use 'v' prefix",
			wantErr: false,
		},
		{
			name:    "no v prefix",
			args:    []string{"--use-v-prefix=false"},
			wantMsg: "Version tags will not use 'v' prefix",
			wantErr: false,
		},
		{
			name:       "detect existing tags with v prefix",
			args:       []string{},
			mockGit:    true,
			mockTags:   true,
			hasVPrefix: true,
			wantMsg:    "Detected existing version tags with 'v' prefix",
			wantErr:    false,
		},
		{
			name:       "detect existing tags without v prefix",
			args:       []string{},
			mockGit:    true,
			mockTags:   true,
			hasVPrefix: false,
			wantMsg:    "Detected existing version tags without 'v' prefix",
			wantErr:    false,
		},
		{
			name:    "error case",
			args:    []string{"--file", ""},
			wantMsg: "failed to initialize project",
			wantErr: true,
		},
		{
			name:       "creates initial git commit and tag when no tags exist",
			args:       []string{},
			mockGit:    true,
			mockTags:   false, // No existing tags
			hasVPrefix: true,
			wantMsg:    "Project initialized with changelog file: CHANGELOG.md",
			wantErr:    false,
		},
		{
			name:       "creates initial git tag without v prefix when configured",
			args:       []string{"--use-v-prefix=false"},
			mockGit:    true,
			mockTags:   false,
			hasVPrefix: false,
			wantMsg:    "Project initialized with changelog file: CHANGELOG.md",
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clear buffers
			logBuf.Reset()
			viper.Reset()

			// Create a new temporary directory for each test case
			tempDir, err := os.MkdirTemp("", "init-cmd-test-"+tc.name)
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Save current working directory and change to temp dir
			currentDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			if err := os.Chdir(tempDir); err != nil {
				t.Fatalf("Failed to change to temp directory: %v", err)
			}
			defer func() {
				if err := os.Chdir(currentDir); err != nil {
					t.Errorf("Failed to restore working directory: %v", err)
				}
			}()

			// Create a new root command for each test
			rootCmd := &cobra.Command{Use: binaryName}
			cmd := &cobra.Command{
				Use:   "init",
				Short: "Initialize a project with SemVer and Keep a Changelog",
				RunE:  runInit,
			}

			// Set up the command flags
			cmd.Flags().String("file", "CHANGELOG.md", "Changelog file name")
			cmd.Flags().Bool("use-v-prefix", true, "Use 'v' prefix for version tags")
			if err := viper.BindPFlag("app.changelog.file", cmd.Flags().Lookup("file")); err != nil {
				t.Fatalf("Failed to bind file flag: %v", err)
			}
			if err := viper.BindPFlag("app.version.use_v_prefix", cmd.Flags().Lookup("use-v-prefix")); err != nil {
				t.Fatalf("Failed to bind use-v-prefix flag: %v", err)
			}
			rootCmd.AddCommand(cmd)

			// Create output buffer for each test
			var outputBuf bytes.Buffer
			rootCmd.SetArgs(append([]string{"init"}, tc.args...))
			rootCmd.SetOut(&outputBuf)
			rootCmd.SetErr(&outputBuf)
			rootCmd.SilenceUsage = true
			rootCmd.SilenceErrors = true

			// Set up real git repository if needed
			if tc.mockGit {
				// Initialize a real git repository
				initCmd := exec.Command("git", "init")
				initCmd.Dir = tempDir
				if err := initCmd.Run(); err != nil {
					t.Skipf("Skipping test (git not available): %v", err)
				}

				// Configure git user for commits (required for git operations)
				configNameCmd := exec.Command("git", "config", "user.name", "Test User")
				configNameCmd.Dir = tempDir
				if err := configNameCmd.Run(); err != nil {
					t.Skipf("Failed to configure git user.name: %v", err)
				}

				configEmailCmd := exec.Command("git", "config", "user.email", "test@example.com")
				configEmailCmd.Dir = tempDir
				if err := configEmailCmd.Run(); err != nil {
					t.Skipf("Failed to configure git user.email: %v", err)
				}

				// Create git tags if needed
				if tc.mockTags {
					// Create an initial commit (required before creating tags)
					readmeFile := filepath.Join(tempDir, "README.md")
					if err := os.WriteFile(readmeFile, []byte("test"), 0644); err != nil {
						t.Fatalf("Failed to create README: %v", err)
					}

					addCmd := exec.Command("git", "add", "README.md")
					addCmd.Dir = tempDir
					if err := addCmd.Run(); err != nil {
						t.Fatalf("Failed to add README: %v", err)
					}

					commitCmd := exec.Command("git", "commit", "-m", "Initial commit")
					commitCmd.Dir = tempDir
					if err := commitCmd.Run(); err != nil {
						t.Fatalf("Failed to create initial commit: %v", err)
					}

					// Create a tag
					tagName := "1.0.0"
					if tc.hasVPrefix {
						tagName = "v1.0.0"
					}
					tagCmd := exec.Command("git", "tag", tagName)
					tagCmd.Dir = tempDir
					if err := tagCmd.Run(); err != nil {
						t.Fatalf("Failed to create tag: %v", err)
					}
				}
			}

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
			} else if !tc.wantErr && tc.wantMsg != "" && !strings.Contains(gotOutput, tc.wantMsg) {
				t.Errorf("Expected output to contain '%s', but got '%s'", tc.wantMsg, gotOutput)
			}

			// Check file existence if not expecting error
			if !tc.wantErr {
				filename := "CHANGELOG.md"
				if cmd.Flags().Changed("file") {
					filename, _ = cmd.Flags().GetString("file")
				}
				if filename != "" {
					_, err := os.Stat(filename)
					assert.NoError(t, err, "Changelog file should exist")
				}

				// Check if version prefix setting was correctly set in viper
				if strings.Contains(tc.name, "detect existing tags") {
					expectedPrefix := tc.hasVPrefix
					assert.Equal(t, expectedPrefix, viper.GetBool("app.version.use_v_prefix"),
						"Version prefix preference should match detected tags")
				} else if cmd.Flags().Changed("use-v-prefix") {
					expectedPrefix, _ := cmd.Flags().GetBool("use-v-prefix")
					assert.Equal(t, expectedPrefix, viper.GetBool("app.version.use_v_prefix"),
						"Version prefix preference should match flag value")
				}

				// Check for git operations when initializing in a repo with no tags
				if tc.mockGit && !tc.mockTags && strings.Contains(tc.name, "creates initial") {
					// Verify initial commit was created
					logCmd := exec.Command("git", "log", "--oneline")
					logCmd.Dir = tempDir
					logOutput, err := logCmd.CombinedOutput()
					assert.NoError(t, err, "Should be able to get git log")
					assert.Contains(t, string(logOutput), "Update changelog for version 0.0.0",
						"Should have created initial commit with changelog")

					// Verify initial tag was created
					tagCmd := exec.Command("git", "tag", "-l")
					tagCmd.Dir = tempDir
					tagOutput, err := tagCmd.CombinedOutput()
					assert.NoError(t, err, "Should be able to list tags")

					expectedTag := "v0.0.0"
					if !tc.hasVPrefix {
						expectedTag = "0.0.0"
					}
					assert.Contains(t, string(tagOutput), expectedTag,
						"Should have created initial tag %s", expectedTag)
				}
			}
		})
	}
}
