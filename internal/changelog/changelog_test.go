package changelog

import (
	"os"
	"strings"
	"testing"
)

func TestAddChangelogSection(t *testing.T) {
	tests := []struct {
		name            string
		initialContent  string
		section         string
		content         string
		expectedContent string
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

			err = AddChangelogSection(tmpfile.Name(), tt.section, tt.content)
			if err != nil {
				t.Fatalf("AddChangelogSection failed: %v", err)
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
