package changelog

import (
	"os"
	"strings"
	"testing"
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
