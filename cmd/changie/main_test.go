package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/peiman/changie/internal/changelog"
)

// Mock implementations
type MockChangelogManager struct {
	initProjectErr         error
	updateChangelogErr     error
	addChangelogSectionErr error
	updateChangelogCalled  int
	isDuplicate            bool
	changelogContent       string
}

func (m *MockChangelogManager) GetChangelogContent() (string, error) {
	if m.changelogContent == "" {
		m.changelogContent = `# Changelog

## [Unreleased]

## [1.0.0] - 2023-01-01

[Unreleased]: https://github.com/peiman/changie/compare/1.0.0...HEAD
[1.0.0]: https://github.com/peiman/changie/releases/tag/1.0.0`
	}
	return m.changelogContent, nil
}

func (m *MockChangelogManager) InitProject(string) error {
	return m.initProjectErr
}
func (m *MockChangelogManager) UpdateChangelog(string, string, string) error {
	m.updateChangelogCalled++
	return m.updateChangelogErr
}
func (m *MockChangelogManager) AddChangelogSection(string, string, string) (bool, error) {
	return m.isDuplicate, m.addChangelogSectionErr
}

type MockGitManager struct {
	commitChangelogErr      error
	tagVersionErr           error
	getProjectVersionErr    error
	projectVersion          string
	getProjectVersionCalled int
	commitChangelogCalled   int
	tagVersionCalled        int
	hasUncommittedChanges   bool
	pushChangesCalled       int
	pushChangesErr          error
}

func (m *MockGitManager) CommitChangelog(string, string) error {
	m.commitChangelogCalled++
	return m.commitChangelogErr
}
func (m *MockGitManager) TagVersion(version string) error {
	m.tagVersionCalled++
	return m.tagVersionErr
}

// Reset the tagVersionCalled counter
func (m *MockGitManager) ResetTagVersionCalled() {
	m.tagVersionCalled = 0
}
func (m *MockGitManager) GetProjectVersion() (string, error) {
	m.getProjectVersionCalled++
	if m.getProjectVersionErr != nil {
		return "", m.getProjectVersionErr
	}
	return m.projectVersion, nil
}
func (m *MockGitManager) HasUncommittedChanges() (bool, error) {
	return m.hasUncommittedChanges, nil
}
func (m *MockGitManager) PushChanges() error {
	m.pushChangesCalled++
	return m.pushChangesErr
}

type MockSemverManager struct {
	bumpMajorErr    error
	bumpMinorErr    error
	bumpPatchErr    error
	bumpMinorCalled int
}

func (m *MockSemverManager) BumpMajor(version string) (string, error) {
	if m.bumpMajorErr != nil {
		return "", m.bumpMajorErr
	}
	parts := strings.Split(version, ".")
	major, _ := strconv.Atoi(parts[0])
	return fmt.Sprintf("%d.0.0", major+1), nil
}

func (m *MockSemverManager) BumpMinor(version string) (string, error) {
	m.bumpMinorCalled++
	if m.bumpMinorErr != nil {
		return "", m.bumpMinorErr
	}
	parts := strings.Split(version, ".")
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	return fmt.Sprintf("%d.%d.0", major, minor+1), nil
}

func (m *MockSemverManager) BumpPatch(version string) (string, error) {
	if m.bumpPatchErr != nil {
		return "", m.bumpPatchErr
	}
	parts := strings.Split(version, ".")
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])
	return fmt.Sprintf("%d.%d.%d", major, minor, patch+1), nil
}

func captureOutput(t *testing.T, f func() error) (string, error) {
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
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to copy output: %v", err)
	}
	return buf.String(), err
}

func TestMainPackage(t *testing.T) {
	isTestMode = true
	defer func() { isTestMode = false }()
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
			expected:         "Project initialized for SemVer and Keep a Changelog.\n",
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
			expected:         "Added section: New feature\n",
			changelogManager: &MockChangelogManager{isDuplicate: false},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Add Duplicate Changelog Section",
			args:             []string{"changie", "changelog", "added", "New feature"},
			expected:         "Added section: New feature (duplicate entry, not added)\n",
			changelogManager: &MockChangelogManager{isDuplicate: true},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Add Changelog Changed Section",
			args:             []string{"changie", "changelog", "changed", "Updated feature"},
			expected:         "Changed section: Updated feature\n",
			changelogManager: &MockChangelogManager{isDuplicate: false},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Add Changelog Deprecated Section",
			args:             []string{"changie", "changelog", "deprecated", "Old feature"},
			expected:         "Deprecated section: Old feature\n",
			changelogManager: &MockChangelogManager{isDuplicate: false},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Add Changelog Removed Section",
			args:             []string{"changie", "changelog", "removed", "Obsolete feature"},
			expected:         "Removed section: Obsolete feature\n",
			changelogManager: &MockChangelogManager{isDuplicate: false},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Add Changelog Fixed Section",
			args:             []string{"changie", "changelog", "fixed", "Bug fix"},
			expected:         "Fixed section: Bug fix\n",
			changelogManager: &MockChangelogManager{isDuplicate: false},
			gitManager:       &MockGitManager{projectVersion: "1.0.0"},
			semverManager:    &MockSemverManager{},
		},
		{
			name:             "Add Changelog Security Section",
			args:             []string{"changie", "changelog", "security", "Security patch"},
			expected:         "Security section: Security patch\n",
			changelogManager: &MockChangelogManager{isDuplicate: false},
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

			output, err := captureOutput(t, func() error {
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

	output, err := captureOutput(t, func() error {
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

func TestGetVersion(t *testing.T) {
	// Save the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(pwd); err != nil {
			t.Errorf("Failed to change directory back: %v", err)
		}
	}()
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "changie-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to the temporary directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize a git repository
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("Failed to initialize git repo: %v", err)
	}
	if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}
	if err := exec.Command("git", "config", "user.name", "Test User").Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	// Test case 1: No tags
	version, err := getVersion()
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	if version != "dev" {
		t.Errorf("Expected 'dev', got '%s'", version)
	}

	// Create a commit
	if err := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit").Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Test case 2: With a tag
	if err := exec.Command("git", "tag", "v1.0.0").Run(); err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}
	version, err = getVersion()
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	if version != "v1.0.0" {
		t.Errorf("Expected 'v1.0.0', got '%s'", version)
	}

	// Test case 3: Commit after tag
	if err := exec.Command("git", "commit", "--allow-empty", "-m", "Another commit").Run(); err != nil {
		t.Fatalf("Failed to create another commit: %v", err)
	}
	version, err = getVersion()
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	if !strings.HasPrefix(version, "v1.0.0-dev.") {
		t.Errorf("Expected version starting with 'v1.0.0-dev.', got '%s'", version)
	}
}

func TestInvalidCommand(t *testing.T) {
	originalArgs := os.Args
	originalExit := exitFunction
	defer func() {
		os.Args = originalArgs
		exitFunction = originalExit
	}()

	os.Args = []string{"changie", "invalid"}

	var exitCode int
	exitFunction = func(code int) {
		exitCode = code
	}

	output, err := captureOutput(t, func() error {
		return run(&MockChangelogManager{}, &MockGitManager{}, &MockSemverManager{})
	})

	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	expectedErrors := []string{
		"Unknown command: invalid",
		"expected command but got \"invalid\"",
	}

	found := false
	for _, expectedError := range expectedErrors {
		if strings.Contains(output, expectedError) || (err != nil && strings.Contains(err.Error(), expectedError)) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected output or error to contain one of %v, but got output: %q and error: %v", expectedErrors, output, err)
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, but got: %d", exitCode)
	}
}

func TestMinorVersionBump(t *testing.T) {
	isTestMode = true
	defer func() { isTestMode = false }()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"changie", "minor"}

	initialVersion := "0.1.0"
	expectedNewVersion := "0.2.0"

	mockGitManager := &MockGitManager{
		projectVersion:        initialVersion,
		hasUncommittedChanges: false,
	}
	mockChangelogManager := &MockChangelogManager{
		changelogContent: fmt.Sprintf(`# Changelog

## [Unreleased]

## [%s] - 2023-01-01

[Unreleased]: https://github.com/peiman/changie/compare/%s...HEAD
[%s]: https://github.com/peiman/changie/releases/tag/%s`, initialVersion, initialVersion, initialVersion, initialVersion),
	}
	mockSemverManager := &MockSemverManager{}

	// Reset the tagVersionCalled counter before running the test
	mockGitManager.ResetTagVersionCalled()

	output, err := captureOutput(t, func() error {
		return run(mockChangelogManager, mockGitManager, mockSemverManager)
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedOutput := fmt.Sprintf("minor release %s done.\nDon't forget to git push and git push --tags.\n", expectedNewVersion)
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain '%s', got: '%s'", expectedOutput, output)
	}

	if mockGitManager.tagVersionCalled != 1 {
		t.Errorf("Expected TagVersion to be called once, got: %d", mockGitManager.tagVersionCalled)
	}
}

func TestVersionMatchBetweenChangelogAndGitTags(t *testing.T) {
	mockGitManager := &MockGitManager{projectVersion: "0.4.0"}
	mockChangelogManager := &MockChangelogManager{
		changelogContent: `## [Unreleased]

## [0.4.0] - 2024-06-26`,
	}

	changelogContent, _ := mockChangelogManager.GetChangelogContent()
	latestChangelogVersion, err := changelog.GetLatestChangelogVersion(changelogContent)
	if err != nil {
		t.Fatalf("Failed to get latest changelog version: %v", err)
	}

	gitVersion, err := mockGitManager.GetProjectVersion()
	if err != nil {
		t.Fatalf("Failed to get git version: %v", err)
	}

	if latestChangelogVersion != gitVersion {
		t.Errorf("Version mismatch: Changelog version %s does not match git tag version %s", latestChangelogVersion, gitVersion)
	}
}

func TestPreventChangelogUpdateOnVersionMismatch(t *testing.T) {
	mockGitManager := &MockGitManager{projectVersion: "0.3.0"}
	mockChangelogManager := &MockChangelogManager{}
	mockSemverManager := &MockSemverManager{}

	os.Args = []string{"changie", "minor"}

	err := run(mockChangelogManager, mockGitManager, mockSemverManager)

	if err == nil {
		t.Error("Expected an error due to version mismatch, but got none")
	}

	if mockChangelogManager.updateChangelogCalled > 0 {
		t.Error("Changelog was updated despite version mismatch")
	}
}

func TestRejectBumpWithUncommittedChanges(t *testing.T) {
	isTestMode = true
	defer func() { isTestMode = false }()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"changie", "minor"}

	initialVersion := "0.1.0"

	mockGitManager := &MockGitManager{
		projectVersion:        initialVersion,
		hasUncommittedChanges: true,
	}
	mockChangelogManager := &MockChangelogManager{
		changelogContent: fmt.Sprintf(`# Changelog

## [Unreleased]

## [%s] - 2023-01-01

[Unreleased]: https://github.com/peiman/changie/compare/%s...HEAD
[%s]: https://github.com/peiman/changie/releases/tag/%s`, initialVersion, initialVersion, initialVersion, initialVersion),
	}
	mockSemverManager := &MockSemverManager{}

	output, err := captureOutput(t, func() error {
		return run(mockChangelogManager, mockGitManager, mockSemverManager)
	})

	if err == nil {
		t.Error("Expected an error due to uncommitted changes, but got none")
	}

	expectedError := "Error: Uncommitted changes found. Please commit or stash your changes before bumping the version."
	if err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message '%s', got: '%v'", expectedError, err)
	}

	if strings.Contains(output, "release") {
		t.Errorf("Unexpected output indicating version bump: %s", output)
	}

	if mockChangelogManager.updateChangelogCalled > 0 {
		t.Error("Changelog was updated despite uncommitted changes")
	}

	if mockGitManager.tagVersionCalled > 0 {
		t.Error("Version was tagged despite uncommitted changes")
	}
}

func TestAllowBumpWithNoUncommittedChanges(t *testing.T) {
	isTestMode = true
	defer func() { isTestMode = false }()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"changie", "minor"}

	initialVersion := "0.1.0"
	expectedNewVersion := "0.2.0"

	mockGitManager := &MockGitManager{
		projectVersion:        initialVersion,
		hasUncommittedChanges: false,
	}
	mockChangelogManager := &MockChangelogManager{
		changelogContent: fmt.Sprintf(`# Changelog

## [Unreleased]

## [%s] - 2023-01-01

[Unreleased]: https://github.com/peiman/changie/compare/%s...HEAD
[%s]: https://github.com/peiman/changie/releases/tag/%s`, initialVersion, initialVersion, initialVersion, initialVersion),
	}
	mockSemverManager := &MockSemverManager{}

	output, err := captureOutput(t, func() error {
		return run(mockChangelogManager, mockGitManager, mockSemverManager)
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedOutput := fmt.Sprintf("minor release %s done.\nDon't forget to git push and git push --tags.\n", expectedNewVersion)
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain '%s', got: '%s'", expectedOutput, output)
	}

	if mockChangelogManager.updateChangelogCalled != 1 {
		t.Errorf("Expected UpdateChangelog to be called once, got: %d", mockChangelogManager.updateChangelogCalled)
	}

	if mockGitManager.tagVersionCalled != 1 {
		t.Errorf("Expected TagVersion to be called once, got: %d", mockGitManager.tagVersionCalled)
	}
}

func TestAutoPushAfterBump(t *testing.T) {
	isTestMode = true
	defer func() { isTestMode = false }()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"changie", "minor", "--auto-push"}

	initialVersion := "0.1.0"
	expectedNewVersion := "0.2.0"

	mockGitManager := &MockGitManager{
		projectVersion:        initialVersion,
		hasUncommittedChanges: false,
	}
	mockChangelogManager := &MockChangelogManager{
		changelogContent: fmt.Sprintf(`# Changelog

## [Unreleased]

## [%s] - 2023-01-01

[Unreleased]: https://github.com/peiman/changie/compare/%s...HEAD
[%s]: https://github.com/peiman/changie/releases/tag/%s`, initialVersion, initialVersion, initialVersion, initialVersion),
	}
	mockSemverManager := &MockSemverManager{}

	output, err := captureOutput(t, func() error {
		return run(mockChangelogManager, mockGitManager, mockSemverManager)
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedOutputs := []string{
		fmt.Sprintf("Current version from git tags: %s", initialVersion),
		fmt.Sprintf("New version: %s", expectedNewVersion),
		"Updating changelog file: CHANGELOG.md",
		fmt.Sprintf("Tagging version: %s", expectedNewVersion),
		fmt.Sprintf("minor release %s done.", expectedNewVersion),
		"Pushing changes and tags...",
		"Automatically pushed changes and tags to remote repository.",
	}

	for _, expectedOutput := range expectedOutputs {
		if !strings.Contains(output, expectedOutput) {
			t.Errorf("Expected output to contain '%s', got: '%s'", expectedOutput, output)
		}
	}

	if mockGitManager.pushChangesCalled != 1 {
		t.Errorf("Expected PushChanges to be called once, got: %d", mockGitManager.pushChangesCalled)
	}
}
func TestVersionBumpWithExistingTag(t *testing.T) {
	isTestMode = true
	defer func() { isTestMode = false }()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"changie", "minor"}

	initialVersion := "0.3.0"
	expectedNewVersion := "0.4.0"

	mockGitManager := &MockGitManager{
		projectVersion:        initialVersion,
		hasUncommittedChanges: false,
		tagVersionErr:         fmt.Errorf("tag %s already exists", expectedNewVersion),
	}
	mockChangelogManager := &MockChangelogManager{
		changelogContent: fmt.Sprintf(`# Changelog

## [Unreleased]

## [%s] - 2023-01-01

[Unreleased]: https://github.com/peiman/changie/compare/%s...HEAD
[%s]: https://github.com/peiman/changie/releases/tag/%s`, initialVersion, initialVersion, initialVersion, initialVersion),
	}
	mockSemverManager := &MockSemverManager{}

	output, err := captureOutput(t, func() error {
		return run(mockChangelogManager, mockGitManager, mockSemverManager)
	})

	if err == nil {
		t.Errorf("Expected an error, got nil")
	} else if !strings.Contains(err.Error(), fmt.Sprintf("Error tagging version: tag %s already exists", expectedNewVersion)) {
		t.Errorf("Expected error to contain 'Error tagging version: tag %s already exists', got: %v", expectedNewVersion, err)
	}

	expectedOutputs := []string{
		fmt.Sprintf("Current version from git tags: %s", initialVersion),
		fmt.Sprintf("New version: %s", expectedNewVersion),
		"Updating changelog file: CHANGELOG.md",
		fmt.Sprintf("Tagging version: %s", expectedNewVersion),
	}

	for _, expectedOutput := range expectedOutputs {
		if !strings.Contains(output, expectedOutput) {
			t.Errorf("Expected output to contain '%s', got: '%s'", expectedOutput, output)
		}
	}
}

func TestVersionMismatchWarning(t *testing.T) {
	mockGitManager := &MockGitManager{projectVersion: "1.0.0"}
	mockChangelogManager := &MockChangelogManager{
		changelogContent: "## [2.0.0]",
	}

	// Test with printWarning = true
	output, err := captureOutput(t, func() error {
		return checkVersionMismatch(mockGitManager, mockChangelogManager, true)
	})

	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	expectedWarning := "Warning: Version mismatch: Git tag version 1.0.0 does not match changelog version 2.0.0"
	if !strings.Contains(output, expectedWarning) {
		t.Errorf("Expected warning %q, but got output: %q", expectedWarning, output)
	}

	// Test with printWarning = false
	output, err = captureOutput(t, func() error {
		return checkVersionMismatch(mockGitManager, mockChangelogManager, false)
	})

	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	if output != "" {
		t.Errorf("Expected no output, but got: %q", output)
	}
}
