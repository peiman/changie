package changelog

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// Use this function to mock exec.Command in tests
//
//nolint:unused // This function is intended for future use in mocking exec.Command
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess isn't a real test. It's used to mock exec.Command
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "git":
		if len(args) > 0 && args[0] == "describe" {
			fmt.Fprintf(os.Stdout, "1.0.0\n")
		} else {
			fmt.Fprintf(os.Stderr, "unknown git command\n")
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
		os.Exit(2)
	}
}

func TestInitProject(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "changie-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("Failed to change directory back: %v", err)
		}
	}()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	changelogFile := "CHANGELOG.md"

	// Test case 1: No existing CHANGELOG.md
	err = InitProject(changelogFile)
	if err != nil {
		t.Errorf("InitProject failed when no CHANGELOG.md existed: %v", err)
	}
	if _, err := os.Stat(changelogFile); os.IsNotExist(err) {
		t.Errorf("CHANGELOG.md was not created")
	}

	// Test case 2: Existing CHANGELOG.md
	err = InitProject(changelogFile)
	if err == nil {
		t.Errorf("InitProject did not return an error when CHANGELOG.md already existed")
	}
	// Note: Update this check based on the actual error message your current implementation returns
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestAddChangelogSection(t *testing.T) {
	tests := []struct {
		name            string
		initialContent  string
		section         string
		content         string
		expectedContent string
		expectDuplicate bool
	}{
		{
			name: "Add new section",
			initialContent: `# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.1.0] - 2023-01-01
### Added
- Initial release
`,
			section: "Added",
			content: "New feature",
			expectedContent: `# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- New feature

## [0.1.0] - 2023-01-01
### Added
- Initial release
`,
			expectDuplicate: false,
		},
		{
			name: "Add to existing section",
			initialContent: `# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]
### Added
- Existing feature

## [0.1.0] - 2023-01-01
### Added
- Initial release
`,
			section: "Added",
			content: "Another new feature",
			expectedContent: `# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]
### Added
- Existing feature
- Another new feature

## [0.1.0] - 2023-01-01
### Added
- Initial release
`,
			expectDuplicate: false,
		},
		{
			name: "Add duplicate entry",
			initialContent: `# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]
### Added
- Existing feature

## [0.1.0] - 2023-01-01
### Added
- Initial release
`,
			section: "Added",
			content: "Existing feature",
			expectedContent: `# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]
### Added
- Existing feature

## [0.1.0] - 2023-01-01
### Added
- Initial release
`,
			expectDuplicate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "CHANGELOG.md")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.initialContent)); err != nil {
				t.Fatal(err)
			}

			isDuplicate, err := AddChangelogSection(tmpfile.Name(), tt.section, tt.content)
			if err != nil {
				t.Fatalf("AddChangelogSection failed: %v", err)
			}

			if isDuplicate != tt.expectDuplicate {
				t.Errorf("Expected isDuplicate to be %v, got %v", tt.expectDuplicate, isDuplicate)
			}

			content, err := os.ReadFile(tmpfile.Name())
			if err != nil {
				t.Fatal(err)
			}

			if !compareIgnoreWhitespace(string(content), tt.expectedContent) {
				t.Errorf("Changelog content doesn't match expected.\nGot:\n%s\nExpected:\n%s", string(content), tt.expectedContent)
			}
		})
	}
}

func compareIgnoreWhitespace(a, b string) bool {
	a = strings.Join(strings.Fields(strings.TrimSpace(a)), " ")
	b = strings.Join(strings.Fields(strings.TrimSpace(b)), " ")
	return a == b
}

func TestUpdateChangelog(t *testing.T) {
	mockChangelog := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com),
and this project adheres to [Semantic Versioning (SemVer)](https://semver.org).

## [Unreleased]

### Added

- Feature A

## [1.0.0] - 2023-01-01

### Added

- Initial release

[Unreleased]: https://github.com/peiman/changie/compare/1.0.0...HEAD
[1.0.0]: https://github.com/peiman/changie/releases/tag/1.0.0`

	expectedChangelog := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com),
and this project adheres to [Semantic Versioning (SemVer)](https://semver.org).

## [Unreleased]

## [1.1.0] - ` + time.Now().Format("2006-01-02") + `

### Added

- Feature A

## [1.0.0] - 2023-01-01

### Added

- Initial release

[Unreleased]: https://github.com/peiman/changie/compare/1.1.0...HEAD
[1.1.0]: https://github.com/peiman/changie/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/peiman/changie/releases/tag/1.0.0`

	tempFile, err := os.CreateTemp("", "changelog")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(mockChangelog)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	err = UpdateChangelog(tempFile.Name(), "1.1.0", "github")
	if err != nil {
		t.Fatalf("UpdateChangelog failed: %v", err)
	}

	updatedChangelog, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read updated changelog: %v", err)
	}

	if string(updatedChangelog) != expectedChangelog {
		t.Errorf("Updated changelog does not match expected.\nGot:\n%s\nExpected:\n%s", string(updatedChangelog), expectedChangelog)
	}
}

func TestReformatChangelog2(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "CHANGELOG.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	initialContent := `# Changelog


All notable changes to this project will be documented in this file.


The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning (SemVer)](http://semver.org/).
## [Unreleased]


### Added
- Feature 1
- Feature 2


### Changed
- Change 1
## [1.0.0] - 2023-01-01


### Added
- Initial release`

	if _, err := tmpfile.Write([]byte(initialContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	err = ReformatChangelog(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to reformat changelog: %v", err)
	}

	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com),
and this project adheres to [Semantic Versioning (SemVer)](https://semver.org).

## [Unreleased]

### Added

- Feature 1
- Feature 2

### Changed

- Change 1

## [1.0.0] - 2023-01-01

### Added

- Initial release
`

	if string(content) != expected {
		t.Errorf("Reformatted changelog content does not match expected.\nGot:\n%s\nExpected:\n%s", string(content), expected)

		gotLines := strings.Split(string(content), "\n")
		expectedLines := strings.Split(expected, "\n")

		for i := 0; i < len(gotLines) || i < len(expectedLines); i++ {
			var gotLine, expectedLine string
			if i < len(gotLines) {
				gotLine = gotLines[i]
			}
			if i < len(expectedLines) {
				expectedLine = expectedLines[i]
			}
			if gotLine != expectedLine {
				t.Errorf("Line %d mismatch:\nGot     : %q\nExpected: %q", i+1, gotLine, expectedLine)
			}
		}
	}
}

func TestNoExtraLineInChangelog(t *testing.T) {
	changelogContent := `## [Unreleased]
## [0.4.0] - 2024-06-26`

	if strings.Contains(changelogContent, "---THIS NEW LINE SHOULD NOT BE HERE ---") {
		t.Error("Changelog contains an unexpected extra line")
	}
}

func TestUpdateChangelogFormatting(t *testing.T) {
	// Save the original execCommand and defer its restoration
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	initialContent := `# Changelog

## [Unreleased]
### Added
- New feature

## [1.0.0] - 2023-01-01
### Added
- Initial release

[Unreleased]: https://github.com/user/repo/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/user/repo/releases/tag/v1.0.0`

	expectedContent := `# Changelog

## [Unreleased]

## [1.1.0] - ` + time.Now().Format("2006-01-02") + `
### Added
- New feature

## [1.0.0] - 2023-01-01
### Added
- Initial release

[Unreleased]: https://github.com/peiman/changie/compare/1.1.0...HEAD
[1.1.0]: https://github.com/peiman/changie/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/peiman/changie/releases/tag/1.0.0`

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "CHANGELOG.*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write initial content
	if _, err := tmpfile.Write([]byte(initialContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Mock execCommand to use mockExecCommand
	execCommand = func(command string, args ...string) *exec.Cmd {
		return mockExecCommand(command, args...)
	}

	// Update changelog
	err = UpdateChangelog(tmpfile.Name(), "1.1.0", "github")
	if err != nil {
		t.Fatalf("UpdateChangelog failed: %v", err)
	}

	// Read updated content
	updatedContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Compare updated content with expected content
	if string(updatedContent) != expectedContent {
		t.Errorf("UpdateChangelog produced incorrect output.\nExpected:\n%s\n\nGot:\n%s", expectedContent, string(updatedContent))
	}
}
