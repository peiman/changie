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
		// Edge cases
		{
			name:        "File with unusual path characters",
			filePath:    filepath.Join(tempDir, "CHANGE LOG-v2 (beta).md"),
			preexisting: false,
			wantErr:     false,
		},
		{
			name:        "File in nonexistent subdirectory",
			filePath:    filepath.Join(tempDir, "nonexistent-dir", "CHANGELOG.md"),
			preexisting: false,
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

	// Create a test changelog with existing sections
	testFileWithSections := filepath.Join(tempDir, "CHANGELOG_WITH_SECTIONS.md")
	changelogWithSections := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Existing feature A

### Fixed

- Existing bug fix B

## [1.0.0] - 2023-01-01

### Added

- Initial release
`
	if err := os.WriteFile(testFileWithSections, []byte(changelogWithSections), 0644); err != nil {
		t.Fatalf("Failed to create test changelog file with sections: %v", err)
	}

	// Create a malformed changelog file
	malformedFile := filepath.Join(tempDir, "MALFORMED.md")
	malformedContent := `# Changelog

This is not a properly formatted changelog.
No Unreleased section here.
`
	if err := os.WriteFile(malformedFile, []byte(malformedContent), 0644); err != nil {
		t.Fatalf("Failed to create malformed changelog file: %v", err)
	}

	// Test cases
	testCases := []struct {
		name      string
		file      string
		section   string
		content   string
		wantErr   bool
		duplicate bool
	}{
		{
			name:      "Add to Added section",
			file:      testFile,
			section:   "Added",
			content:   "New feature X",
			wantErr:   false,
			duplicate: false,
		},
		{
			name:      "Add to Fixed section",
			file:      testFile,
			section:   "Fixed",
			content:   "Bug Y",
			wantErr:   false,
			duplicate: false,
		},
		{
			name:      "Invalid section",
			file:      testFile,
			section:   "Invalid",
			content:   "Something",
			wantErr:   true,
			duplicate: false,
		},
		{
			name:      "Duplicate entry",
			file:      testFile,
			section:   "Added",
			content:   "New feature X",
			wantErr:   false,
			duplicate: true,
		},
		// Edge cases
		{
			name:      "Add to existing section",
			file:      testFileWithSections,
			section:   "Added",
			content:   "Another new feature",
			wantErr:   false,
			duplicate: false,
		},
		{
			name:      "Add to nonexistent file",
			file:      filepath.Join(tempDir, "NONEXISTENT.md"),
			section:   "Added",
			content:   "New feature",
			wantErr:   true,
			duplicate: false,
		},
		{
			name:      "Add to malformed file",
			file:      malformedFile,
			section:   "Added",
			content:   "New feature",
			wantErr:   true,
			duplicate: false,
		},
		{
			name:      "Content with special characters",
			file:      testFile,
			section:   "Changed",
			content:   "Update API endpoint from `/api/v1/*` to `/api/v2/*` with **breaking changes**",
			wantErr:   false,
			duplicate: false,
		},
		{
			name:      "Content with very long text",
			file:      testFile,
			section:   "Added",
			content:   strings.Repeat("This is a very long changelog entry that exceeds the typical line length. ", 10),
			wantErr:   false,
			duplicate: false,
		},
		{
			name:      "Empty content",
			file:      testFile,
			section:   "Added",
			content:   "",
			wantErr:   false,
			duplicate: false,
		},
		{
			name:      "Case sensitivity check",
			file:      testFile,
			section:   "added", // lowercase, should fail
			content:   "Test case sensitivity",
			wantErr:   true,
			duplicate: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			isDuplicate, err := AddChangelogSection(tc.file, tc.section, tc.content)

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
				content, err := os.ReadFile(tc.file)
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
- Latest thing

## [1.9.0] - 2023-01-15
### Fixed
- Old bug

## [1.0.0] - 2023-01-01
### Added
- Initial release
`,
			wantVersion: "2.0.0",
			wantErr:     false,
		},
		// Edge cases
		{
			name: "Version with v prefix",
			content: `# Changelog
## [Unreleased]
### Added
- Something

## [v1.2.3] - 2023-01-01
### Added
- Something with v prefix
`,
			wantVersion: "v1.2.3", // Keep the v prefix
			wantErr:     false,
		},
		{
			name: "Version with multiple dots",
			content: `# Changelog
## [Unreleased]
### Added
- Something

## [1.2.3.4] - 2023-01-01
### Added
- Invalid semver but should extract anyway
`,
			wantVersion: "1.2.3.4", // Not valid semver but should extract
			wantErr:     false,
		},
		{
			name: "Version with leading zeros (invalid for semver)",
			content: `# Changelog
## [Unreleased]
### Added
- Something

## [01.02.03] - 2023-01-01
### Added
- Version with leading zeros
`,
			wantVersion: "01.02.03", // Should extract as-is, validation happens elsewhere
			wantErr:     false,
		},
		{
			name:        "Empty content",
			content:     "",
			wantVersion: "",
			wantErr:     false,
		},
		{
			name: "Malformed changelog",
			content: `# This is not a changelog
Just some random text
`,
			wantVersion: "",
			wantErr:     false,
		},
		{
			name: "Version with pre-release suffix",
			content: `# Changelog
## [Unreleased]
### Added
- Something

## [1.0.0-beta.1] - 2023-01-01
### Added
- Beta release
`,
			wantVersion: "1.0.0-beta.1",
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			version, err := GetLatestChangelogVersion(tc.content)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetLatestChangelogVersion() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
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

	// Test cases
	testCases := []struct {
		name               string
		initialContent     string
		version            string
		repositoryProvider string
		wantErr            bool
		expectedLinkCount  int
	}{
		{
			name:               "Simple update",
			initialContent:     changelogTemplate,
			version:            "1.0.0",
			repositoryProvider: "github",
			wantErr:            false,
			expectedLinkCount:  2, // Unreleased and 1.0.0
		},
		{
			name: "Update with existing version",
			initialContent: `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- New feature

## [1.0.0] - 2023-01-01

### Added

- Initial release

[Unreleased]: https://github.com/user/repo/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/user/repo/releases/tag/v1.0.0
`,
			version:            "2.0.0",
			repositoryProvider: "github",
			wantErr:            false,
			expectedLinkCount:  3, // Unreleased, 2.0.0, and 1.0.0
		},
		// Edge cases
		{
			name:               "Version with invalid leading zeros",
			initialContent:     changelogTemplate,
			version:            "01.02.03", // Invalid semver with leading zeros
			repositoryProvider: "github",
			wantErr:            true, // Should error on invalid semver
			expectedLinkCount:  0,
		},
		{
			name:               "Bitbucket provider",
			initialContent:     changelogTemplate,
			version:            "1.0.0",
			repositoryProvider: "bitbucket",
			wantErr:            false,
			expectedLinkCount:  2,
		},
		{
			name:               "Unknown provider (defaults to github)",
			initialContent:     changelogTemplate,
			version:            "1.0.0",
			repositoryProvider: "unknown",
			wantErr:            false,
			expectedLinkCount:  2,
		},
		{
			name: "Malformed changelog (no Unreleased section)",
			initialContent: `# Changelog

This is not a properly formatted changelog.
`,
			version:            "1.0.0",
			repositoryProvider: "github",
			wantErr:            true,
			expectedLinkCount:  0,
		},
		{
			name: "Existing comparison links with different format",
			initialContent: `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- New feature

## [1.0.0] - 2023-01-01

### Added

- Initial release

[Unreleased]: https://custom-domain.com/compare/v1.0.0...HEAD
[1.0.0]: https://custom-domain.com/releases/tag/v1.0.0
`,
			version:            "2.0.0",
			repositoryProvider: "github",
			wantErr:            false,
			expectedLinkCount:  3, // Check that existing links are preserved
		},
		{
			name: "Very long section content",
			initialContent: `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

` + generateLongEntries(20) + `

`,
			version:            "1.0.0",
			repositoryProvider: "github",
			wantErr:            false,
			expectedLinkCount:  2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create the changelog file for this test case
			testFile := filepath.Join(tempDir, fmt.Sprintf("CHANGELOG_%s.md", strings.ReplaceAll(tc.name, " ", "_")))
			if err := os.WriteFile(testFile, []byte(tc.initialContent), 0644); err != nil {
				t.Fatalf("Failed to create test changelog file: %v", err)
			}

			// Call the function
			err := UpdateChangelog(testFile, tc.version, tc.repositoryProvider)
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateChangelog() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				// Read the updated content
				content, err := os.ReadFile(testFile)
				if err != nil {
					t.Errorf("Failed to read updated changelog file: %v", err)
					return
				}

				// Check if the new version header was added
				expectedHeader := fmt.Sprintf("## [%s] -", tc.version)
				if !strings.Contains(string(content), expectedHeader) {
					t.Errorf("Expected version header '%s' not found in updated changelog", expectedHeader)
				}

				// Check if Unreleased section still exists
				if !strings.Contains(string(content), "## [Unreleased]") {
					t.Errorf("Expected 'Unreleased' section not found in updated changelog")
				}

				// Check link count
				linkCount := countLinks(string(content))
				if linkCount != tc.expectedLinkCount {
					t.Errorf("Expected %d links, found %d links in updated changelog", tc.expectedLinkCount, linkCount)
				}
			}
		})
	}
}

// Helper function to generate a large number of changelog entries
func generateLongEntries(count int) string {
	var entries []string
	for i := 0; i < count; i++ {
		entries = append(entries, fmt.Sprintf("- Entry %d with a very long description that goes on and on to test how the system handles very large changelog entries that might cause issues with parsing or updating", i+1))
	}
	return strings.Join(entries, "\n")
}

// Helper function to count comparison links in a changelog
func countLinks(content string) int {
	lines := strings.Split(content, "\n")
	count := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]: ") {
			count++
		}
	}
	return count
}
