package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// Mock implementations
type MockChangelogManager struct {
	initProjectErr         error
	updateChangelogErr     error
	addChangelogSectionErr error
	updateChangelogCalled  int
}

func (m *MockChangelogManager) InitProject(string) error {
	return m.initProjectErr
}
func (m *MockChangelogManager) UpdateChangelog(string, string, string) error {
	m.updateChangelogCalled++
	return m.updateChangelogErr
}
func (m *MockChangelogManager) AddChangelogSection(string, string, string) error {
	return m.addChangelogSectionErr
}

type MockGitManager struct {
	commitChangelogErr      error
	tagVersionErr           error
	getProjectVersionErr    error
	projectVersion          string
	getProjectVersionCalled int
	commitChangelogCalled   int
	tagVersionCalled        int
}

func (m *MockGitManager) CommitChangelog(string, string) error {
	m.commitChangelogCalled++
	return m.commitChangelogErr
}
func (m *MockGitManager) TagVersion(string) error {
	m.tagVersionCalled++
	return m.tagVersionErr
}
func (m *MockGitManager) GetProjectVersion() (string, error) {
	m.getProjectVersionCalled++
	if m.getProjectVersionErr != nil {
		return "", m.getProjectVersionErr
	}
	return m.projectVersion, nil
}

type MockSemverManager struct {
	bumpMajorErr    error
	bumpMinorErr    error
	bumpPatchErr    error
	bumpMinorCalled int
}

func (m *MockSemverManager) BumpMajor(string) (string, error) {
	if m.bumpMajorErr != nil {
		return "", m.bumpMajorErr
	}
	return "2.0.0", nil
}
func (m *MockSemverManager) BumpMinor(string) (string, error) {
	m.bumpMinorCalled++
	if m.bumpMinorErr != nil {
		return "", m.bumpMinorErr
	}
	return "1.1.0", nil
}
func (m *MockSemverManager) BumpPatch(string) (string, error) {
	if m.bumpPatchErr != nil {
		return "", m.bumpPatchErr
	}
	return "1.0.1", nil
}
func captureOutput(f func() error) (string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	err := f()

	w.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestMainPackage(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name             string
		args             []string
		expected         string
		changelogManager ChangelogManager
		gitManager       GitManager
		semverManager    SemverManager
	}{
		{
			name:             "Init Command",
			args:             []string{"changie", "init"},
			expected:         "Project initialized for semver and Keep a Changelog.\n",
			changelogManager: &MockChangelogManager{},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Major Version Bump",
			args:             []string{"changie", "major"},
			expected:         "major release 2.0.0 done.\nDon't forget to git push and git push --tags.\n",
			changelogManager: &MockChangelogManager{},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Minor Version Bump",
			args:             []string{"changie", "minor"},
			expected:         "minor release 1.1.0 done.\nDon't forget to git push and git push --tags.\n",
			changelogManager: &MockChangelogManager{},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Patch Version Bump",
			args:             []string{"changie", "patch"},
			expected:         "patch release 1.0.1 done.\nDon't forget to git push and git push --tags.\n",
			changelogManager: &MockChangelogManager{},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Add Changelog Section",
			args:             []string{"changie", "changelog", "added", "New feature"},
			expected:         "Added Added section to changelog: New feature\n",
			changelogManager: &MockChangelogManager{},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Error Adding Changelog Section",
			args:             []string{"changie", "changelog", "added", "New feature"},
			expected:         "Error adding changelog section: mock error\n",
			changelogManager: &MockChangelogManager{addChangelogSectionErr: fmt.Errorf("mock error")},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Error Getting Project Version",
			args:             []string{"changie", "major"},
			expected:         "Error getting project version: mock error\n",
			changelogManager: &MockChangelogManager{},
			gitManager:       &MockGitManager{getProjectVersionErr: fmt.Errorf("mock error")},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Error Bumping Version",
			args:             []string{"changie", "major"},
			expected:         "Error bumping version: mock error\n",
			changelogManager: &MockChangelogManager{},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{bumpMajorErr: fmt.Errorf("mock error")},
		},
		{
			name:             "Error Updating Changelog",
			args:             []string{"changie", "major"},
			expected:         "Error updating changelog: mock error\n",
			changelogManager: &MockChangelogManager{updateChangelogErr: fmt.Errorf("mock error")},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Error Committing Changelog",
			args:             []string{"changie", "major"},
			expected:         "Error committing changelog: mock error\n",
			changelogManager: &MockChangelogManager{},
			gitManager:       &MockGitManager{projectVersion: "1.0.0", commitChangelogErr: fmt.Errorf("mock error")},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Error Tagging Version",
			args:             []string{"changie", "major"},
			expected:         "Error tagging version: mock error\n",
			changelogManager: &MockChangelogManager{},
			gitManager:       &MockGitManager{projectVersion: "1.0.0", tagVersionErr: fmt.Errorf("mock error")},
			semverManager:    &MockSemverManager{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			output, err := captureOutput(func() error {
				return run(tt.changelogManager, tt.gitManager, tt.semverManager)
			})

			if err != nil {
				output += err.Error() + "\n"
			}

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Test case %s failed.\nExpected output to contain: %q\nGot: %q", tt.name, tt.expected, output)
			} else {
				t.Logf("Test case %s passed", tt.name)
			}

			t.Logf("Test case: %s\nArgs: %v\nExpected: %q\nGot: %q", tt.name, tt.args, tt.expected, output)
		})
	}
}

func TestGitNotInstalled(t *testing.T) {
	oldIsGitInstalled := isGitInstalled
	defer func() { isGitInstalled = oldIsGitInstalled }()

	isGitInstalled = func() bool { return false }

	os.Args = []string{"changie", "major"}

	output, err := captureOutput(func() error {
		return run(&MockChangelogManager{}, &MockGitManager{}, &MockSemverManager{})
	})

	expected := "Error: Git is not installed."
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error %q, but got: %v", expected, err)
	}
	if output != "" {
		t.Errorf("Expected no output, but got: %q", output)
	}
}

func TestInvalidCommand(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"changie", "invalid"}

	output, err := captureOutput(func() error {
		return run(&MockChangelogManager{}, &MockGitManager{}, &MockSemverManager{})
	})

	expectedError := "expected command but got \"invalid\""
	if err == nil {
		t.Errorf("Expected an error for invalid command, but got none")
	} else if err.Error() != expectedError {
		t.Errorf("Expected error %q, but got: %v", expectedError, err)
	}

	if output != "" {
		t.Errorf("Expected no output, but got: %q", output)
	}
}

func TestMinorVersionBump(t *testing.T) {
	// Set up test environment
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"changie", "minor"}

	// Mock dependencies
	mockChangelogManager := &MockChangelogManager{}
	mockGitManager := &MockGitManager{projectVersion: "0.1.0"}
	mockSemverManager := &MockSemverManager{}

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	err := run(mockChangelogManager, mockGitManager, mockSemverManager)

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check results
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedOutput := "minor release 1.1.0 done.\nDon't forget to git push and git push --tags.\n"
	if output != expectedOutput {
		t.Errorf("Expected output '%s', got: '%s'", expectedOutput, output)
	}

	// Check if the correct methods were called
	if mockGitManager.getProjectVersionCalled != 1 {
		t.Errorf("Expected GetProjectVersion to be called once")
	}
	if mockSemverManager.bumpMinorCalled != 1 {
		t.Errorf("Expected BumpMinor to be called once")
	}
	if mockChangelogManager.updateChangelogCalled != 1 {
		t.Errorf("Expected UpdateChangelog to be called once")
	}
	if mockGitManager.commitChangelogCalled != 1 {
		t.Errorf("Expected CommitChangelog to be called once")
	}
	if mockGitManager.tagVersionCalled != 1 {
		t.Errorf("Expected TagVersion to be called once")
	}
}
