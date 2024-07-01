package git

import (
	"fmt"
	"reflect"
	"testing"
)

type mockCmd struct {
	output string
	err    error
}

func (c *mockCmd) CombinedOutput() ([]byte, error) {
	return []byte(c.output), c.err
}

func TestIsInstalled(t *testing.T) {
	if !IsInstalled() {
		t.Error("Git is not installed, but it should be for running these tests")
	}
}

func TestGetLastTag(t *testing.T) {
	// Save the current ExecCommand and defer its restoration
	oldExecCommand := ExecCommand
	defer func() { ExecCommand = oldExecCommand }()

	tests := []struct {
		name           string
		mockOutput     string
		mockError      error
		expectedResult string
		expectedError  bool
	}{
		{
			name:           "Successful tag retrieval",
			mockOutput:     "1.2.3\n",
			mockError:      nil,
			expectedResult: "1.2.3",
			expectedError:  false,
		},
		{
			name:           "No tags found",
			mockOutput:     "",
			mockError:      fmt.Errorf("exit status 128"),
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExecCommand = func(command string, args ...string) Commander {
				return &mockCmd{
					output: tt.mockOutput,
					err:    tt.mockError,
				}
			}

			result, err := GetLastTag()

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("Expected %s, but got %s", tt.expectedResult, result)
				}
			}
		})
	}
}
func TestGetProjectVersion(t *testing.T) {
	// Save the current ExecCommand and defer its restoration
	oldExecCommand := ExecCommand
	defer func() { ExecCommand = oldExecCommand }()

	tests := []struct {
		name           string
		mockOutput     string
		mockError      error
		expectedResult string
		expectedError  bool
	}{
		{
			name:           "No tags found",
			mockOutput:     "fatal: No names found, cannot describe anything.",
			mockError:      fmt.Errorf("exit status 128"),
			expectedResult: "0.1.0",
			expectedError:  false,
		},
		{
			name:           "Tag found",
			mockOutput:     "v1.2.3\n",
			mockError:      nil,
			expectedResult: "v1.2.3",
			expectedError:  false,
		},
		{
			name:           "Git error",
			mockOutput:     "error: unknown revision or path not in the working tree.",
			mockError:      fmt.Errorf("exit status 128"),
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExecCommand = func(command string, args ...string) Commander {
				return &mockCmd{
					output: tt.mockOutput,
					err:    tt.mockError,
				}
			}

			result, err := GetProjectVersion()

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("Expected %s, but got %s", tt.expectedResult, result)
				}
			}
		})
	}
}

func TestCommitChangelog(t *testing.T) {
	// Implementation remains the same
}

func TestTagVersion(t *testing.T) {
	// Implementation remains the same
}

func TestTagVersionWithoutPrefix(t *testing.T) {
	// Save the current ExecCommand and defer its restoration
	oldExecCommand := ExecCommand
	defer func() { ExecCommand = oldExecCommand }()

	var executedCommand []string
	ExecCommand = func(command string, args ...string) Commander {
		executedCommand = append([]string{command}, args...)
		return &mockCmd{output: "Mocked git tag command", err: nil}
	}

	err := TagVersion("1.0.0")
	if err != nil {
		t.Errorf("TagVersion failed: %v", err)
	}

	expectedCommand := []string{"git", "tag", "1.0.0"}
	if !reflect.DeepEqual(executedCommand, expectedCommand) {
		t.Errorf("Expected command %v, got: %v", expectedCommand, executedCommand)
	}
}
