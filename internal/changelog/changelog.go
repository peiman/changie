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

	"github.com/blang/semver/v4"

	"github.com/peiman/changie/internal/git"
	"github.com/peiman/changie/internal/logger"
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
		return fmt.Errorf("changelog file already exists: %s (use an alternative filename or delete the existing file if you want to recreate it)", filePath)
	}

	// Only proceed if the error is because the file doesn't exist
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if changelog file exists: %w (verify you have read permissions for the directory)", err)
	}

	// Create the file with the template content
	err = os.WriteFile(filePath, []byte(changelogTemplate), 0o644)
	if err != nil {
		return fmt.Errorf("failed to create changelog file: %w (verify you have write permissions for the directory and sufficient disk space)", err)
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
	if !ValidSections[section] {
		return false, fmt.Errorf("invalid section: %s, must be one of: Added, Changed, Deprecated, Removed, Fixed, Security (section names are case-sensitive)", section)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to read changelog file: %w (check if '%s' exists and you have read permissions)", err, filePath)
	}

	lines := strings.Split(string(data), "\n")
	unreleasedIndex := findUnreleasedSection(lines)
	if unreleasedIndex == -1 {
		return false, fmt.Errorf("unreleased section not found in changelog (ensure the file follows the Keep a Changelog format with an '## [Unreleased]' section)")
	}

	sectionHeader := fmt.Sprintf("### %s", section)
	sectionIndex, nextMajorIndex := findSection(lines, unreleasedIndex, sectionHeader)

	if sectionIndex == -1 {
		return false, createNewSection(filePath, lines, unreleasedIndex, nextMajorIndex, sectionHeader, content)
	}

	return addToExistingSection(filePath, lines, sectionIndex, nextMajorIndex, content)
}

func findUnreleasedSection(lines []string) int {
	for i, line := range lines {
		if strings.TrimSpace(line) == "## [Unreleased]" {
			return i
		}
	}
	return -1
}

func findSection(lines []string, unreleasedIndex int, sectionHeader string) (sectionIndex, nextMajorIndex int) {
	sectionIndex = -1
	nextMajorIndex = len(lines)

	for i := unreleasedIndex + 1; i < len(lines); i++ {
		trimmedLine := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmedLine, "## ") {
			nextMajorIndex = i
			break
		}
		if trimmedLine == sectionHeader {
			sectionIndex = i
		}
	}
	return sectionIndex, nextMajorIndex
}

func createNewSection(filePath string, lines []string, unreleasedIndex, nextMajorIndex int, sectionHeader, content string) error {
	lastSectionIndex := unreleasedIndex
	for i := unreleasedIndex + 1; i < nextMajorIndex; i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "### ") {
			lastSectionIndex = i
		}
	}

	result := make([]string, 0, len(lines)+5)
	result = append(result, lines[:lastSectionIndex+1]...)

	if lastSectionIndex > unreleasedIndex {
		lastSectionEndIndex := findSectionEnd(lines, lastSectionIndex, nextMajorIndex)
		if lastSectionEndIndex > lastSectionIndex+1 {
			result = append(result, lines[lastSectionIndex+1:lastSectionEndIndex]...)
		}
		result = append(result, "", sectionHeader)
	} else {
		result = append(result, "", sectionHeader)
	}

	result = append(result, "", fmt.Sprintf("- %s", content))

	if nextMajorIndex < len(lines) {
		result = append(result, "", "")
		result = append(result, lines[nextMajorIndex:]...)
	}

	err := os.WriteFile(filePath, []byte(strings.Join(result, "\n")), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write updated changelog: %w (verify you have write permissions for the file)", err)
	}
	return nil
}

func findSectionEnd(lines []string, sectionIndex, nextMajorIndex int) int {
	for i := sectionIndex + 1; i < nextMajorIndex; i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "### ") {
			return i
		}
	}
	return nextMajorIndex
}

func addToExistingSection(filePath string, lines []string, sectionIndex, nextMajorIndex int, content string) (bool, error) {
	newEntry := fmt.Sprintf("- %s", content)
	nextSectionIndex := findSectionEnd(lines, sectionIndex, nextMajorIndex)

	if isDuplicate(lines, sectionIndex, nextSectionIndex, newEntry) {
		return true, nil
	}

	result := make([]string, 0, len(lines)+2)
	result = append(result, lines[:sectionIndex+1]...)

	if hasContent(lines, sectionIndex, nextSectionIndex) {
		lastContentLine := findLastContentLine(lines, sectionIndex, nextSectionIndex)
		result = append(result, lines[sectionIndex+1:lastContentLine+1]...)
		result = append(result, newEntry)
	} else {
		result = append(result, "", newEntry)
	}

	if nextSectionIndex < len(lines) {
		result = append(result, "")
		result = append(result, lines[nextSectionIndex:]...)
	}

	err := os.WriteFile(filePath, []byte(strings.Join(result, "\n")), 0o644)
	if err != nil {
		return false, fmt.Errorf("failed to write updated changelog: %w (verify you have write permissions for the file)", err)
	}
	return false, nil
}

func isDuplicate(lines []string, sectionIndex, nextSectionIndex int, newEntry string) bool {
	for i := sectionIndex + 1; i < nextSectionIndex; i++ {
		if strings.TrimSpace(lines[i]) == newEntry {
			return true
		}
	}
	return false
}

func hasContent(lines []string, sectionIndex, nextSectionIndex int) bool {
	for i := sectionIndex + 1; i < nextSectionIndex; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "" && !strings.HasPrefix(trimmed, "###") {
			return true
		}
	}
	return false
}

func findLastContentLine(lines []string, sectionIndex, nextSectionIndex int) int {
	lastContentLine := sectionIndex
	for i := sectionIndex + 1; i < nextSectionIndex; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "" && !strings.HasPrefix(trimmed, "###") && !strings.HasPrefix(trimmed, "##") {
			lastContentLine = i
		}
	}
	return lastContentLine
}

// GetLatestChangelogVersion finds the latest version in the changelog.
// Returns an empty string if no version is found.
func GetLatestChangelogVersion(content string) (string, error) {
	// Use regex to find version headers like "## [1.2.3] - 2023-01-01" or "## [v1.2.3-alpha] - 2023-01-01"
	// Match any version format, including those with 'v' prefix, multiple segments, and prerelease/build metadata
	re := regexp.MustCompile(`## \[([v]?[0-9]+(?:\.[0-9]+)*(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?)\]`)
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
	// Check if the version starts with 'v' and remove it for semver validation
	semverStr := version
	if strings.HasPrefix(version, "v") {
		semverStr = strings.TrimPrefix(version, "v")
	}

	// Validate the version using strict semver
	_, err := semver.Parse(semverStr)
	if err != nil {
		return fmt.Errorf("invalid version format: %w (verify that the version follows semantic versioning rules)", err)
	}

	// Read the current changelog file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read changelog file: %w (verify that '%s' exists and you have read permissions)", err, filePath)
	}

	content := string(data)

	// Check if Unreleased section exists
	unreleasedRegex := regexp.MustCompile(`## \[Unreleased\]`)
	if !unreleasedRegex.MatchString(content) {
		return fmt.Errorf("unreleased section not found in changelog (ensure your changelog follows the Keep a Changelog format with an '## [Unreleased]' section)")
	}

	// Get current date
	today := time.Now().Format("2006-01-02")

	// Create new version section - use the original version format for display
	versionHeader := fmt.Sprintf("## [%s] - %s", version, today)

	// Replace Unreleased with the version and add a new Unreleased section
	content = unreleasedRegex.ReplaceAllString(content, fmt.Sprintf(`## [Unreleased]

%s`, versionHeader))

	// Update or add comparison links
	// Try to get repository info from git remote
	linkPrefix := ""
	repoInfo, err := git.GetRepositoryInfo()
	if err != nil {
		// If we can't get repo info from git, use provided repository provider with placeholder
		logger.Changelog.Warn().Err(err).Msg("Failed to get repository info from git, using placeholder URL")

		switch repositoryProvider {
		case "github":
			linkPrefix = "https://github.com/user/repo/compare/"
		case "bitbucket":
			linkPrefix = "https://bitbucket.org/user/repo/compare/"
		case "gitlab":
			linkPrefix = "https://gitlab.com/user/repo/-/compare/"
		default:
			linkPrefix = "https://github.com/user/repo/compare/"
		}
	} else {
		// Use actual repository info
		logger.Changelog.Info().
			Str("owner", repoInfo.Owner).
			Str("repo", repoInfo.Repo).
			Str("provider", repoInfo.Provider).
			Msg("Using repository info from git remote")

		// Override the provider with the detected one from git URL
		repositoryProvider = repoInfo.Provider

		switch repositoryProvider {
		case "github":
			linkPrefix = fmt.Sprintf("%s/%s/%s/compare/", repoInfo.BaseURL, repoInfo.Owner, repoInfo.Repo)
		case "bitbucket":
			linkPrefix = fmt.Sprintf("%s/%s/%s/compare/", repoInfo.BaseURL, repoInfo.Owner, repoInfo.Repo)
		case "gitlab":
			linkPrefix = fmt.Sprintf("%s/%s/%s/-/compare/", repoInfo.BaseURL, repoInfo.Owner, repoInfo.Repo)
		default:
			// For unknown providers, use GitHub-style format
			linkPrefix = fmt.Sprintf("%s/%s/%s/compare/", repoInfo.BaseURL, repoInfo.Owner, repoInfo.Repo)
		}
	}

	// Use the version as-is for tag links (respects the user's v-prefix preference)
	tagVersion := version

	// Find the previous version by looking at version headers in the changelog
	// The previous version is the first version header after our new version
	versionHeaderRegex := regexp.MustCompile(`## \[([^\]]+)\] - `)
	versionMatches := versionHeaderRegex.FindAllStringSubmatch(content, -1)
	var previousVersion string

	for i, match := range versionMatches {
		if len(match) >= 2 && match[1] == version {
			// Found our new version, check if there's a next one (which is the previous release)
			if i+1 < len(versionMatches) {
				nextMatch := versionMatches[i+1]
				if nextMatch[1] != "Unreleased" {
					previousVersion = nextMatch[1]
				}
			}
			break
		}
	}

	// Find where the links section starts
	linksRegex := regexp.MustCompile(`(?m)^\[.+?\]: .+$`)
	linksMatches := linksRegex.FindAllStringIndex(content, -1)

	// Build the new version link
	var newVersionLink string
	if previousVersion != "" {
		// Link to comparison with previous version
		newVersionLink = fmt.Sprintf("[%s]: %s%s...%s", version, linkPrefix, previousVersion, tagVersion)
	} else {
		// First release - link to releases/tag instead of comparison
		baseURL := strings.TrimSuffix(linkPrefix, "/compare/")
		newVersionLink = fmt.Sprintf("[%s]: %s/releases/tag/%s", version, baseURL, tagVersion)
	}

	// Build new links section
	newUnreleasedLink := fmt.Sprintf("[Unreleased]: %s%s...HEAD", linkPrefix, tagVersion)

	if len(linksMatches) > 0 {
		// Extract existing links section and preserve non-Unreleased links
		existingLinksStart := linksMatches[0][0]
		linksSection := content[existingLinksStart:]

		// Split into individual link lines and filter
		linkLines := strings.Split(linksSection, "\n")
		var otherLinks []string
		for _, line := range linkLines {
			trimmed := strings.TrimSpace(line)
			// Keep links that are not Unreleased and not the new version
			if trimmed != "" && !strings.HasPrefix(trimmed, "[Unreleased]:") && !strings.HasPrefix(trimmed, "["+version+"]:") {
				otherLinks = append(otherLinks, trimmed)
			}
		}

		// Build new links section: Unreleased, new version, then all existing links
		var linksBuilder strings.Builder
		linksBuilder.WriteString(newUnreleasedLink)
		linksBuilder.WriteString("\n")
		linksBuilder.WriteString(newVersionLink)
		if len(otherLinks) > 0 {
			linksBuilder.WriteString("\n")
			linksBuilder.WriteString(strings.Join(otherLinks, "\n"))
		}

		// Replace the entire links section
		content = strings.TrimRight(content[:existingLinksStart], "\n") + "\n\n" + linksBuilder.String()
	} else {
		// No existing links, add new section
		content = strings.TrimRight(content, "\n") + "\n\n" + newUnreleasedLink + "\n" + newVersionLink
	}

	// Write the updated content back to the file
	err = os.WriteFile(filePath, []byte(content), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write updated changelog: %w (check if you have write permissions for the file and sufficient disk space)", err)
	}

	return nil
}
