package changelog

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/peiman/changie/internal/semver"
)

// InitProject initializes the project with a new CHANGELOG.md file
func InitProject(changelogFile string) error {
	content := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
`
	if err := os.WriteFile(changelogFile, []byte(content), 0644); err != nil {
		return err
	}
	return ReformatChangelog(changelogFile)
}

// UpdateChangelog updates the CHANGELOG.md file with the new version
func UpdateChangelog(file, version, provider string) error {
	// Read the changelog file
	changelogContent, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Prepare the new version entry
	newVersionEntry := fmt.Sprintf("## [%s] - %s\n\n### Added\n\n- Feature A\n\n", version, time.Now().Format("2006-01-02"))

	// Replace the placeholder for the "Unreleased" section with the new version entry
	updatedContent := strings.Replace(string(changelogContent), "## [Unreleased]\n\n### Added\n\n- Feature A\n\n", "## [Unreleased]\n\n"+newVersionEntry, 1)

	// Write the updated content back to the changelog file
	return os.WriteFile(file, []byte(updatedContent), 0644)
}

func ReformatChangelog(changelogFile string) error {
	content, err := os.ReadFile(changelogFile)
	if err != nil {
		return fmt.Errorf("error reading changelog: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var reformattedLines []string
	lastLineWasEmpty := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmedLine, "# "):
			reformattedLines = append(reformattedLines, trimmedLine, "")
			lastLineWasEmpty = true
		case trimmedLine == "All notable changes to this project will be documented in this file.":
			reformattedLines = append(reformattedLines, trimmedLine, "")
			lastLineWasEmpty = true
		case strings.Contains(trimmedLine, "Keep a Changelog"):
			reformattedLines = append(reformattedLines, "The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),")
			lastLineWasEmpty = false
		case strings.Contains(trimmedLine, "Semantic Versioning"):
			reformattedLines = append(reformattedLines, "and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).", "")
			lastLineWasEmpty = true
		case strings.HasPrefix(trimmedLine, "## "):
			if !lastLineWasEmpty {
				reformattedLines = append(reformattedLines, "")
			}
			reformattedLines = append(reformattedLines, trimmedLine, "")
			lastLineWasEmpty = true
		case strings.HasPrefix(trimmedLine, "### "):
			if !lastLineWasEmpty {
				reformattedLines = append(reformattedLines, "")
			}
			reformattedLines = append(reformattedLines, trimmedLine, "")
			lastLineWasEmpty = true
		case trimmedLine != "":
			reformattedLines = append(reformattedLines, trimmedLine)
			lastLineWasEmpty = false
		case !lastLineWasEmpty:
			reformattedLines = append(reformattedLines, "")
			lastLineWasEmpty = true
		}
	}

	// Remove trailing newlines
	for len(reformattedLines) > 0 && reformattedLines[len(reformattedLines)-1] == "" {
		reformattedLines = reformattedLines[:len(reformattedLines)-1]
	}

	// Add a single newline at the end of the file
	reformattedLines = append(reformattedLines, "")

	err = os.WriteFile(changelogFile, []byte(strings.Join(reformattedLines, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("error writing changelog: %w", err)
	}

	return nil
}

// AddChangelogSection adds a new section to the Unreleased part of the changelog
func AddChangelogSection(changelogFile, section, content string) (bool, error) {
	// Read the entire file
	existingContent, err := os.ReadFile(changelogFile)
	if err != nil {
		return false, fmt.Errorf("error reading changelog: %w", err)
	}

	lines := strings.Split(string(existingContent), "\n")
	var newLines []string
	unreleasedIndex := -1
	sectionOrder := []string{"Added", "Changed", "Deprecated", "Removed", "Fixed", "Security"}
	sections := make(map[string][]string)

	// Find the [Unreleased] section
	for i, line := range lines {
		if strings.HasPrefix(line, "## [Unreleased]") {
			unreleasedIndex = i
			break
		}
	}

	// If [Unreleased] section doesn't exist, create it
	if unreleasedIndex == -1 {
		unreleasedIndex = 0
		newLines = append(newLines, "## [Unreleased]", "")
	} else {
		// Copy lines before and including [Unreleased]
		newLines = append(newLines, lines[:unreleasedIndex+1]...)
		newLines = append(newLines, "")
	}

	// Process existing entries in [Unreleased]
	currentSection := ""
	for i := unreleasedIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "## [") {
			break
		}
		if strings.HasPrefix(line, "### ") {
			currentSection = strings.TrimPrefix(line, "### ")
			continue
		}
		if line != "" && currentSection != "" {
			sections[currentSection] = append(sections[currentSection], line)
		}
	}

	// Add the new content to the appropriate section, but only if it doesn't already exist
	newEntry := fmt.Sprintf("- %s", content)
	isDuplicate := contains(sections[section], newEntry)
	if !isDuplicate {
		sections[section] = append(sections[section], newEntry)
	}

	// Add sections in the correct order
	for _, s := range sectionOrder {
		if len(sections[s]) > 0 {
			newLines = append(newLines, fmt.Sprintf("### %s", s))
			newLines = append(newLines, sections[s]...)
			newLines = append(newLines, "")
		}
	}

	// Add the rest of the file
	for i := unreleasedIndex + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## [") {
			newLines = append(newLines, lines[i:]...)
			break
		}
	}

	// Remove any trailing empty lines
	for len(newLines) > 0 && strings.TrimSpace(newLines[len(newLines)-1]) == "" {
		newLines = newLines[:len(newLines)-1]
	}

	// Write the updated content back to the file
	err = os.WriteFile(changelogFile, []byte(strings.Join(newLines, "\n")), 0644)
	if err != nil {
		return false, fmt.Errorf("error writing changelog: %w", err)
	}
	// Reformat the entire changelog after adding the new section
	err = ReformatChangelog(changelogFile)
	if err != nil {
		return false, fmt.Errorf("error reformatting changelog: %w", err)
	}

	return isDuplicate, nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func updateDiffLinks(lines []string, newVersion, provider string) []string {
	var updatedLines []string
	var versions []string
	linkLines := map[string]string{}

	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]: ") {
			parts := strings.SplitN(line, "]: ", 2)
			version := strings.Trim(parts[0], "[]")
			linkLines[version] = parts[1]
			if version != "Unreleased" {
				versions = append(versions, version)
			}
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	// Sort versions in descending order
	sort.Slice(versions, func(i, j int) bool {
		result, _ := semver.Compare(versions[i], versions[j])
		return result > 0
	})

	// Insert new version
	versions = append([]string{newVersion}, versions...)

	baseURL := getCompareURL(provider)

	// Update comparison links
	updatedLines = append(updatedLines, fmt.Sprintf("[Unreleased]: %s/compare/%s...HEAD", baseURL, newVersion))
	for i := 0; i < len(versions)-1; i++ {
		updatedLines = append(updatedLines, fmt.Sprintf("[%s]: %s/compare/%s...%s", versions[i], baseURL, versions[i+1], versions[i]))
	}
	lastVersion := versions[len(versions)-1]
	updatedLines = append(updatedLines, fmt.Sprintf("[%s]: %s/releases/tag/%s", lastVersion, baseURL, lastVersion))

	return updatedLines
}

func getCompareURL(provider string) string {
	switch provider {
	case "github":
		return "https://github.com/peiman/changie"
	case "bitbucket":
		return "https://bitbucket.org/peiman/changie"
	default:
		return ""
	}
}
