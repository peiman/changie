package git

import (
	"fmt"
	"strings"
	"testing"
)

type mockCmd struct {
	output []byte
	err    error
}

func (c *mockCmd) CombinedOutput() ([]byte, error) {
	return c.output, c.err
}

func TestIsInstalled(t *testing.T) {
	oldExecCommand := ExecCommand
	defer func() { ExecCommand = oldExecCommand }()

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte("git version 2.30.1"), err: nil}
	}

	if !IsInstalled() {
		t.Error("IsInstalled returned false, expected true")
	}

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte(""), err: fmt.Errorf("git not found")}
	}

	if IsInstalled() {
		t.Error("IsInstalled returned true, expected false")
	}
}

func TestGetVersion(t *testing.T) {
	oldExecCommand := ExecCommand
	defer func() { ExecCommand = oldExecCommand }()

	tests := []struct {
		name        string
		mockOutputs map[string][]byte
		mockErrors  map[string]error
		expected    string
		expectError bool
	}{
		{
			name: "Tagged version",
			mockOutputs: map[string][]byte{
				"git describe --tags --abbrev=0":         []byte("v1.2.3"),
				"git describe --exact-match --tags HEAD": []byte("v1.2.3"),
			},
			expected:    "v1.2.3",
			expectError: false,
		},
		{
			name: "Dev version",
			mockOutputs: map[string][]byte{
				"git describe --tags --abbrev=0":    []byte("v1.2.3"),
				"git rev-parse --short HEAD":        []byte("abc1234"),
				"git rev-list v1.2.3..HEAD --count": []byte("5"),
			},
			mockErrors: map[string]error{
				"git describe --exact-match --tags HEAD": fmt.Errorf("not a tag"),
			},
			expected:    "v1.2.3-dev.5+abc1234",
			expectError: false,
		},
		{
			name: "No tags",
			mockOutputs: map[string][]byte{
				"git describe --tags --abbrev=0": []byte("fatal: No names found, cannot describe anything"),
			},
			mockErrors: map[string]error{
				"git describe --tags --abbrev=0": fmt.Errorf("exit status 128"),
			},
			expected:    "dev",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExecCommand = func(command string, args ...string) Commander {
				cmdString := command + " " + strings.Join(args, " ")
				return &mockCmd{
					output: tt.mockOutputs[cmdString],
					err:    tt.mockErrors[cmdString],
				}
			}

			version, err := GetVersion()

			t.Logf("Test case: %s", tt.name)
			t.Logf("Version: %s", version)
			t.Logf("Error: %v", err)
			if err != nil {
				t.Logf("Error type: %T", err)
				t.Logf("Error contains 'No names found': %v", strings.Contains(err.Error(), "No names found"))
			}

			if tt.expectError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if version != tt.expected {
				t.Errorf("Expected version %s, but got %s", tt.expected, version)
			}
		})
	}
}

func TestCommitChangelog(t *testing.T) {
	oldExecCommand := ExecCommand
	defer func() { ExecCommand = oldExecCommand }()

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte(""), err: nil}
	}

	err := CommitChangelog("CHANGELOG.md", "v1.0.0")
	if err != nil {
		t.Errorf("CommitChangelog failed: %v", err)
	}

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte(""), err: fmt.Errorf("git error")}
	}

	err = CommitChangelog("CHANGELOG.md", "v1.0.0")
	if err == nil {
		t.Error("CommitChangelog should have failed, but didn't")
	}
}

func TestTagVersion(t *testing.T) {
	oldExecCommand := ExecCommand
	defer func() { ExecCommand = oldExecCommand }()

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte(""), err: nil}
	}

	err := TagVersion("v1.0.0")
	if err != nil {
		t.Errorf("TagVersion failed: %v", err)
	}

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte(""), err: fmt.Errorf("git error")}
	}

	err = TagVersion("v1.0.0")
	if err == nil {
		t.Error("TagVersion should have failed, but didn't")
	}
}

func TestHasUncommittedChanges(t *testing.T) {
	oldExecCommand := ExecCommand
	defer func() { ExecCommand = oldExecCommand }()

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte(" M file.txt"), err: nil}
	}

	hasChanges, err := HasUncommittedChanges()
	if err != nil {
		t.Errorf("HasUncommittedChanges failed: %v", err)
	}
	if !hasChanges {
		t.Error("HasUncommittedChanges should have returned true, but returned false")
	}

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte(""), err: nil}
	}

	hasChanges, err = HasUncommittedChanges()
	if err != nil {
		t.Errorf("HasUncommittedChanges failed: %v", err)
	}
	if hasChanges {
		t.Error("HasUncommittedChanges should have returned false, but returned true")
	}
}

func TestPushChanges(t *testing.T) {
	oldExecCommand := ExecCommand
	defer func() { ExecCommand = oldExecCommand }()

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte(""), err: nil}
	}

	err := PushChanges()
	if err != nil {
		t.Errorf("PushChanges failed: %v", err)
	}

	ExecCommand = func(command string, args ...string) Commander {
		return &mockCmd{output: []byte(""), err: fmt.Errorf("git error")}
	}

	err = PushChanges()
	if err == nil {
		t.Error("PushChanges should have failed, but didn't")
	}
}
