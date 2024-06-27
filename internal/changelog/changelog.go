package changelog

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/peiman/changie/internal/semver"
)

// InitProject initializes the project with a new CHANGELOG.md file
func InitProject(changelogFile string) error {
	// Check if CHANGELOG.md already exists
	if _, err := os.Stat(changelogFile); err == nil {
		return fmt.Errorf("CHANGELOG.md already exists. Please rename or remove the existing file before running changie init.\n\n" +
			"After initializing with changie, you can manually transfer the content from your old changelog to the new one, " +
			"following the Keep a Changelog format: https://keepachangelog.com/")
	}

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

func GetLatestChangelogVersion(content string) (string, error) {
	re := regexp.MustCompile(`## \[(\d+\.\d+\.\d+)\]`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("no version found in changelog")
	}
	return matches[1], nil
}

// UpdateChangelog updates the CHANGELOG.md file with the new version
func UpdateChangelog(file string, version string, provider string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	unreleasedAdded := false
	versionAdded := false

	for _, line := range lines {
		if strings.HasPrefix(line, "## [Unreleased]") && !unreleasedAdded {
			newLines = append(newLines, "## [Unreleased]", "")
			newLines = append(newLines, fmt.Sprintf("## [%s] - %s", version, time.Now().Format("2006-01-02")))
			unreleasedAdded = true
			versionAdded = true
		} else if strings.HasPrefix(line, "## [") && !versionAdded {
			newLines = append(newLines, fmt.Sprintf("## [%s] - %s", version, time.Now().Format("2006-01-02")))
			newLines = append(newLines, line)
			versionAdded = true
		} else {
			newLines = append(newLines, line)
		}
	}

	// Update comparison links
	updatedLines := updateDiffLinks(newLines, version, provider)

	return os.WriteFile(file, []byte(strings.Join(updatedLines, "\n")), 0644)
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

	// Preserve existing lines and collect version information
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

	// Add new version if it doesn't exist
	if !contains(versions, newVersion) {
		versions = append(versions, newVersion)
	}

	// Sort versions in descending order
	sort.Slice(versions, func(i, j int) bool {
		result, _ := semver.Compare(versions[i], versions[j])
		return result > 0
	})

	baseURL := getCompareURL(provider)

	// Update comparison links
	newLinkLines := []string{fmt.Sprintf("[Unreleased]: %s/compare/%s...HEAD", baseURL, versions[0])}
	for i := 0; i < len(versions)-1; i++ {
		newLinkLines = append(newLinkLines, fmt.Sprintf("[%s]: %s/compare/%s...%s", versions[i], baseURL, versions[i+1], versions[i]))
	}
	lastVersion := versions[len(versions)-1]
	newLinkLines = append(newLinkLines, fmt.Sprintf("[%s]: %s/releases/tag/%s", lastVersion, baseURL, lastVersion))

	// Append updated link lines
	updatedLines = append(updatedLines, newLinkLines...)

	return updatedLines
}

func getCompareURL(provider string) string {
	switch provider {
	case "github":
		return "https://github.com/peiman/changie"
	case "bitbucket":
		return "https://bitbucket.org/peiman/changie"
	default:
		return "https://github.com/peiman/changie" // Default to GitHub
	}
}
