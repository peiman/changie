// Package changelog provides functionality for managing CHANGELOG.md files
// following the Keep a Changelog format.
//
// This package handles all aspects of changelog management according to the
// Keep a Changelog (https://keepachangelog.com) specification, including:
//
// - Initializing new changelog files with proper formatting
// - Adding entries to different changelog sections (Added, Changed, Fixed, etc.)
// - Updating changelogs during version releases
// - Managing comparison links between versions
// - Extracting version information from changelog content
//
// The package maintains proper formatting, spacing, and section organization
// to ensure clean, consistent changelog files that follow established conventions.
package changelog

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Template for a new changelog file following Keep a Changelog format
// with only the Unreleased section and no predefined categories
const changelogTemplate = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

`

// ValidSections represents the valid section types for a changelog
// according to the Keep a Changelog specification.
var ValidSections = map[string]bool{
	"Added":      true, // New features
	"Changed":    true, // Changes to existing functionality
	"Deprecated": true, // Features that will be removed in future versions
	"Removed":    true, // Features that were removed in this version
	"Fixed":      true, // Bug fixes
	"Security":   true, // Vulnerability fixes
}

// InitProject initializes a project with an empty changelog file.
//
// This function creates a new CHANGELOG.md file at the specified path
// using the standard Keep a Changelog template format, with an empty
// Unreleased section ready for new entries.
//
// Parameters:
//   - filePath: Path where the changelog file should be created
//
// Returns:
//   - error: Error if the file already exists or cannot be created
func InitProject(filePath string) error {
	// Check if the file already exists
	_, err := os.Stat(filePath)
	if err == nil {
		return fmt.Errorf("changelog file already exists: %s", filePath)
	}

	// Only proceed if the error is because the file doesn't exist
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if changelog file exists: %w", err)
	}

	// Create the file with the template content
	err = os.WriteFile(filePath, []byte(changelogTemplate), 0644)
	if err != nil {
		return fmt.Errorf("failed to create changelog file: %w", err)
	}

	return nil
}

// AddChangelogSection adds a new entry to a specific section in the changelog.
//
// This function handles several scenarios:
// 1. If the specified section doesn't exist in the Unreleased section, it creates it
// 2. If the section exists, it adds the entry to the section
// 3. If the entry already exists in the section, it doesn't duplicate it
//
// The function maintains proper spacing and formatting in the changelog file.
//
// Parameters:
//   - filePath: Path to the changelog file
//   - section: Section name (must be one of the ValidSections)
//   - content: Text content to add as a new bullet point entry
//
// Returns:
//   - bool: true if the entry already existed (duplicate), false if added successfully
//   - error: Any error encountered during file operations or if section is invalid
func AddChangelogSection(filePath, section, content string) (bool, error) {
	// Validate section
	if !ValidSections[section] {
		return false, fmt.Errorf("invalid section: %s, must be one of: Added, Changed, Deprecated, Removed, Fixed, Security", section)
	}

	// Read the current changelog file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to read changelog file: %w", err)
	}

	fileContent := string(data)

	// Split the content into lines
	lines := strings.Split(fileContent, "\n")

	// Find the Unreleased section
	unreleasedHeader := "## [Unreleased]"
	unreleasedIndex := -1

	for i, line := range lines {
		if strings.TrimSpace(line) == unreleasedHeader {
			unreleasedIndex = i
			break
		}
	}

	if unreleasedIndex == -1 {
		return false, fmt.Errorf("unreleased section not found in changelog")
	}

	// Find the appropriate section header
	sectionHeader := fmt.Sprintf("### %s", section)
	sectionFound := false
	sectionIndex := -1

	// Define the boundary of our search (from Unreleased to the next major section)
	nextMajorSectionIndex := len(lines)
	for i := unreleasedIndex + 1; i < len(lines); i++ {
		trimmedLine := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmedLine, "## ") {
			nextMajorSectionIndex = i
			break
		}

		// Check if this is our section
		if trimmedLine == sectionHeader {
			sectionFound = true
			sectionIndex = i
		}
	}

	// If section is not found, we need to create it
	if !sectionFound {
		// Create a new section
		// First, find where to insert it - right after Unreleased or after the last existing section

		// Find the last section within Unreleased
		lastSectionIndex := unreleasedIndex
		for i := unreleasedIndex + 1; i < nextMajorSectionIndex; i++ {
			trimmedLine := strings.TrimSpace(lines[i])
			if strings.HasPrefix(trimmedLine, "### ") {
				lastSectionIndex = i
			}
		}

		// Create a result array with the new section
		result := make([]string, 0, len(lines)+5) // Allocate extra space

		// Add everything up to and including the last section (or Unreleased if no sections)
		result = append(result, lines[:lastSectionIndex+1]...)

		// If we're adding after an existing section, preserve proper spacing
		if lastSectionIndex > unreleasedIndex {
			// Find the end of the last section
			lastSectionEndIndex := nextMajorSectionIndex
			for i := lastSectionIndex + 1; i < nextMajorSectionIndex; i++ {
				if strings.HasPrefix(strings.TrimSpace(lines[i]), "### ") {
					lastSectionEndIndex = i
					break
				}
			}

			// Include content of the last section
			if lastSectionEndIndex > lastSectionIndex+1 {
				result = append(result, lines[lastSectionIndex+1:lastSectionEndIndex]...)
			}

			// Add one newline, then our new section
			result = append(result, "", sectionHeader)
		} else {
			// Adding right after Unreleased, add a newline then our section
			result = append(result, "", sectionHeader)
		}

		// Add our new entry
		result = append(result, "", fmt.Sprintf("- %s", content))

		// Add the rest of the file
		if nextMajorSectionIndex < len(lines) {
			// Add a blank line before the next section
			result = append(result, "")
			// Then add the rest of the lines
			result = append(result, lines[nextMajorSectionIndex:]...)
		}

		// Write the updated content back to the file
		err = os.WriteFile(filePath, []byte(strings.Join(result, "\n")), 0644)
		if err != nil {
			return false, fmt.Errorf("failed to write updated changelog: %w", err)
		}

		return false, nil
	}

	// Section exists, check if content already exists
	newEntry := fmt.Sprintf("- %s", content)
	isDuplicate := false

	// Find the boundaries of this section
	nextSectionIndex := nextMajorSectionIndex
	for i := sectionIndex + 1; i < nextMajorSectionIndex; i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "### ") {
			nextSectionIndex = i
			break
		}
	}

	// Check entries in this section (until next section or end)
	for i := sectionIndex + 1; i < nextSectionIndex; i++ {
		line := strings.TrimSpace(lines[i])
		if line == newEntry {
			isDuplicate = true
			break
		}
	}

	if isDuplicate {
		return true, nil
	}

	// Create a result with the new entry
	result := make([]string, 0, len(lines)+2) // Allocate space for potential new lines

	// Copy all lines up to and including the section header
	result = append(result, lines[:sectionIndex+1]...)

	// Check if there's any content in this section
	hasContent := false
	for i := sectionIndex + 1; i < nextSectionIndex; i++ {
		if strings.TrimSpace(lines[i]) != "" && !strings.HasPrefix(strings.TrimSpace(lines[i]), "###") {
			hasContent = true
			break
		}
	}

	if hasContent {
		// Existing content found, add our entry after the last content line
		lastContentLine := sectionIndex
		for i := sectionIndex + 1; i < nextSectionIndex; i++ {
			trimmedLine := strings.TrimSpace(lines[i])
			if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "###") && !strings.HasPrefix(trimmedLine, "##") {
				lastContentLine = i
			}
		}

		// Add everything up to the last content line
		result = append(result, lines[sectionIndex+1:lastContentLine+1]...)

		// Add the new entry after the last content
		result = append(result, newEntry)
	} else {
		// No content in this section, add a blank line after header
		result = append(result, "", newEntry)
	}

	// Add the rest of the file with proper spacing
	if nextSectionIndex < len(lines) {
		// Add a blank line before the next section
		result = append(result, "")
		// Then add the rest of the lines
		result = append(result, lines[nextSectionIndex:]...)
	}

	// Write the updated content back to the file
	err = os.WriteFile(filePath, []byte(strings.Join(result, "\n")), 0644)
	if err != nil {
		return false, fmt.Errorf("failed to write updated changelog: %w", err)
	}

	return false, nil
}

// GetLatestChangelogVersion finds the latest version in the changelog.
// Returns an empty string if no version is found.
func GetLatestChangelogVersion(content string) (string, error) {
	// Use regex to find version headers like "## [1.2.3] - 2023-01-01"
	re := regexp.MustCompile(`## \[([0-9]+\.[0-9]+\.[0-9]+)\]`)
	matches := re.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		// No version found
		return "", nil
	}

	// First match is the latest version (assuming the changelog follows the format)
	return matches[0][1], nil
}

// UpdateChangelog updates the changelog by converting Unreleased to a new version.
func UpdateChangelog(filePath, version, repositoryProvider string) error {
	// Read the current changelog file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read changelog file: %w", err)
	}

	content := string(data)

	// Check if Unreleased section exists
	unreleasedRegex := regexp.MustCompile(`## \[Unreleased\]`)
	if !unreleasedRegex.MatchString(content) {
		return fmt.Errorf("unreleased section not found in changelog")
	}

	// Get current date
	today := time.Now().Format("2006-01-02")

	// Create new version section
	versionHeader := fmt.Sprintf("## [%s] - %s", version, today)

	// Replace Unreleased with the version and add a new Unreleased section
	content = unreleasedRegex.ReplaceAllString(content, fmt.Sprintf(`## [Unreleased]

%s`, versionHeader))

	// Update or add comparison links
	linkPrefix := ""
	switch repositoryProvider {
	case "github":
		linkPrefix = "https://github.com/user/repo/compare/"
	case "bitbucket":
		linkPrefix = "https://bitbucket.org/user/repo/compare/"
	default:
		linkPrefix = "https://github.com/user/repo/compare/"
	}

	// Check if version link section exists at end of file
	linkRegex := regexp.MustCompile(`\[Unreleased\]: .*`)
	if linkRegex.MatchString(content) {
		// Update existing links
		content = linkRegex.ReplaceAllString(content, fmt.Sprintf(`[Unreleased]: %sv%s...HEAD
[%s]: %s...v%s`, linkPrefix, version, version, linkPrefix, version))
	} else {
		// Add new links section
		content += fmt.Sprintf(`

[Unreleased]: %sv%s...HEAD
[%s]: %s...v%s`, linkPrefix, version, version, linkPrefix, version)
	}

	// Write the updated content back to the file
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated changelog: %w", err)
	}

	return nil
}
