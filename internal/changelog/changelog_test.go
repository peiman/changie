package changelog

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestInitProject(t *testing.T) {
	tempFile, err := os.CreateTemp("", "CHANGELOG.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	err = InitProject(tempFile.Name())
	if err != nil {
		t.Fatalf("InitProject() returned an error: %v", err)
	}

	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read changelog file: %v", err)
	}

	expectedContent := "# Changelog"
	if !strings.Contains(string(content), expectedContent) {
		t.Errorf("InitProject() did not create the expected content, got: %s", string(content))
	}
}

func TestUpdateChangelog(t *testing.T) {
	tempFile, err := os.CreateTemp("", "CHANGELOG.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	initialContent := `# Changelog
## [Unreleased]
### Added
- New feature

## [1.0.0] - 2023-01-01
### Added
- Initial release
`
	err = os.WriteFile(tempFile.Name(), []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial content: %v", err)
	}

	err = UpdateChangelog(tempFile.Name(), "1.1.0", "github")
	if err != nil {
		t.Fatalf("UpdateChangelog() returned an error: %v", err)
	}

	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read changelog file: %v", err)
	}

	expectedContent := "## [1.1.0]"
	if !strings.Contains(string(content), expectedContent) {
		t.Errorf("UpdateChangelog() did not update the content as expected, got: %s", string(content))
	}
}

func TestAddChangelogSection(t *testing.T) {
	tempFile, err := os.CreateTemp("", "CHANGELOG.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	initialContent := `# Changelog
## [Unreleased]

## [1.0.0] - 2023-01-01
### Added
- Initial release
`
	err = os.WriteFile(tempFile.Name(), []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial content: %v", err)
	}

	err = AddChangelogSection(tempFile.Name(), "Added")
	if err != nil {
		t.Fatalf("AddChangelogSection() returned an error: %v", err)
	}

	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read changelog file: %v", err)
	}

	expectedContent := "### Added"
	if !strings.Contains(string(content), expectedContent) {
		t.Errorf("AddChangelogSection() did not add the section as expected, got: %s", string(content))
	}
}
func TestUpdateChangelogFormat(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "CHANGELOG.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	initialContent := `# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]
### Added
- Feature A

## [0.1.0] - 2023-01-01
### Added
- Initial release

[Unreleased]: https://github.com/peiman/changie/compare/0.1.0...HEAD
[0.1.0]: https://github.com/peiman/changie/releases/tag/0.1.0
`
	if _, err := tmpfile.Write([]byte(initialContent)); err != nil {
		t.Fatal(err)
	}

	err = UpdateChangelog(tmpfile.Name(), "0.2.0", "github")
	if err != nil {
		t.Fatalf("UpdateChangelog failed: %v", err)
	}

	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := `# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.2.0] - ` + time.Now().Format("2006-01-02") + `
### Added
- Feature A

## [0.1.0] - 2023-01-01
### Added
- Initial release

[Unreleased]: https://github.com/peiman/changie/compare/0.2.0...HEAD
[0.2.0]: https://github.com/peiman/changie/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/peiman/changie/releases/tag/0.1.0
`

	if !compareIgnoreWhitespace(string(content), expected) {
		t.Errorf("Changelog format doesn't match expected.\nGot:\n%s\nExpected:\n%s", string(content), expected)
	}
}

func compareIgnoreWhitespace(a, b string) bool {
	a = strings.Join(strings.Fields(a), " ")
	b = strings.Join(strings.Fields(b), " ")
	return a == b
}
