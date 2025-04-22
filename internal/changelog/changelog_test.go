package changelog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitProject(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "changelog-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test cases
	testCases := []struct {
		name        string
		filePath    string
		preexisting bool
		wantErr     bool
	}{
		{
			name:        "New file creation",
			filePath:    filepath.Join(tempDir, "CHANGELOG.md"),
			preexisting: false,
			wantErr:     false,
		},
		{
			name:        "File already exists",
			filePath:    filepath.Join(tempDir, "EXISTING.md"),
			preexisting: true,
			wantErr:     true,
		},
	}

	// Create the preexisting file
	preexistingFile := filepath.Join(tempDir, "EXISTING.md")
	if err := os.WriteFile(preexistingFile, []byte("existing content"), 0644); err != nil {
		t.Fatalf("Failed to create preexisting file: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			err := InitProject(tc.filePath)

			// Check error result
			if (err != nil) != tc.wantErr {
				t.Errorf("InitProject() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// Check file existence
			if !tc.wantErr {
				_, err := os.Stat(tc.filePath)
				if err != nil {
					t.Errorf("Expected file to exist at %s, but got error: %v", tc.filePath, err)
				}

				// Verify file content
				content, err := os.ReadFile(tc.filePath)
				if err != nil {
					t.Errorf("Failed to read created file: %v", err)
				}
				if string(content) != changelogTemplate {
					t.Errorf("File content does not match template")
				}
			}
		})
	}
}

func TestAddChangelogSection(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "changelog-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test changelog file
	testFile := filepath.Join(tempDir, "CHANGELOG.md")
	if err := os.WriteFile(testFile, []byte(changelogTemplate), 0644); err != nil {
		t.Fatalf("Failed to create test changelog file: %v", err)
	}

	// Test cases
	testCases := []struct {
		name      string
		section   string
		content   string
		wantErr   bool
		duplicate bool
	}{
		{
			name:      "Add to Added section",
			section:   "Added",
			content:   "New feature X",
			wantErr:   false,
			duplicate: false,
		},
		{
			name:      "Add to Fixed section",
			section:   "Fixed",
			content:   "Bug Y",
			wantErr:   false,
			duplicate: false,
		},
		{
			name:      "Invalid section",
			section:   "Invalid",
			content:   "Something",
			wantErr:   true,
			duplicate: false,
		},
		{
			name:      "Duplicate entry",
			section:   "Added",
			content:   "New feature X",
			wantErr:   false,
			duplicate: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			isDuplicate, err := AddChangelogSection(testFile, tc.section, tc.content)

			// Check error result
			if (err != nil) != tc.wantErr {
				t.Errorf("AddChangelogSection() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// Check duplicate flag
			if !tc.wantErr && isDuplicate != tc.duplicate {
				t.Errorf("AddChangelogSection() isDuplicate = %v, want %v", isDuplicate, tc.duplicate)
			}

			// For non-error cases, check if content was added
			if !tc.wantErr && !isDuplicate {
				content, err := os.ReadFile(testFile)
				if err != nil {
					t.Errorf("Failed to read changelog file: %v", err)
					return
				}

				// Check if entry was added
				expectedEntry := fmt.Sprintf("- %s", tc.content)
				if !strings.Contains(string(content), expectedEntry) {
					t.Errorf("Expected entry '%s' not found in changelog", expectedEntry)
				}
			}
		})
	}
}

func TestGetLatestChangelogVersion(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		content     string
		wantVersion string
		wantErr     bool
	}{
		{
			name: "No version",
			content: `# Changelog
## [Unreleased]
### Added
- Something
`,
			wantVersion: "",
			wantErr:     false,
		},
		{
			name: "One version",
			content: `# Changelog
## [Unreleased]
### Added
- Something

## [1.2.3] - 2023-01-01
### Added
- Another thing
`,
			wantVersion: "1.2.3",
			wantErr:     false,
		},
		{
			name: "Multiple versions",
			content: `# Changelog
## [Unreleased]
### Added
- Something

## [2.0.0] - 2023-02-01
### Added
- Major feature

## [1.2.3] - 2023-01-01
### Added
- Another thing
`,
			wantVersion: "2.0.0",
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			version, err := GetLatestChangelogVersion(tc.content)

			// Check error result
			if (err != nil) != tc.wantErr {
				t.Errorf("GetLatestChangelogVersion() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// Check version
			if version != tc.wantVersion {
				t.Errorf("GetLatestChangelogVersion() = %v, want %v", version, tc.wantVersion)
			}
		})
	}
}

func TestUpdateChangelog(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "changelog-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test changelog file
	testFile := filepath.Join(tempDir, "CHANGELOG.md")

	// Test cases
	testCases := []struct {
		name               string
		initialContent     string
		version            string
		repositoryProvider string
		wantErr            bool
		checkContent       func(string) bool
	}{
		{
			name:               "Valid update",
			initialContent:     changelogTemplate,
			version:            "1.0.0",
			repositoryProvider: "github",
			wantErr:            false,
			checkContent: func(content string) bool {
				// Check if version section was added
				versionHeader := "## [1.0.0]"
				unreleasedSection := "## [Unreleased]"
				links := "[Unreleased]: https://github.com/user/repo/compare/v1.0.0...HEAD"

				return strings.Contains(content, versionHeader) &&
					strings.Contains(content, unreleasedSection) &&
					strings.Contains(content, links)
			},
		},
		{
			name:               "No unreleased section",
			initialContent:     "# Changelog\n\nSome content\n",
			version:            "1.0.0",
			repositoryProvider: "github",
			wantErr:            true,
			checkContent:       func(content string) bool { return true },
		},
		{
			name:               "Bitbucket provider",
			initialContent:     changelogTemplate,
			version:            "2.0.0",
			repositoryProvider: "bitbucket",
			wantErr:            false,
			checkContent: func(content string) bool {
				// Check if bitbucket links were used
				links := "[Unreleased]: https://bitbucket.org/user/repo/compare/v2.0.0...HEAD"
				return strings.Contains(content, links)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create initial content
			if err := os.WriteFile(testFile, []byte(tc.initialContent), 0644); err != nil {
				t.Fatalf("Failed to create test changelog file: %v", err)
			}

			// Call the function
			err := UpdateChangelog(testFile, tc.version, tc.repositoryProvider)

			// Check error result
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateChangelog() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// For non-error cases, check content
			if !tc.wantErr {
				content, err := os.ReadFile(testFile)
				if err != nil {
					t.Errorf("Failed to read changelog file: %v", err)
					return
				}

				if !tc.checkContent(string(content)) {
					t.Errorf("UpdateChangelog() produced unexpected content:\n%s", string(content))
				}
			}
		})
	}
}
