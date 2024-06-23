package changelog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// InitProject initializes the project with a new CHANGELOG.md file
func InitProject(changelogFile string) error {
	content := `# Changelog
All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
## [Unreleased]
`
	return os.WriteFile(changelogFile, []byte(content), 0644)
}

// UpdateChangelog updates the CHANGELOG.md file with the new version
func UpdateChangelog(file, version, provider string) error {
	log.Printf("Updating changelog file: %s", file)
	absPath, err := filepath.Abs(file)
	if err != nil {
		log.Printf("Error getting absolute path: %v", err)
		return fmt.Errorf("error getting absolute path: %v", err)
	}
	log.Printf("Absolute path of changelog: %s", absPath)
	content, err := os.ReadFile(absPath)
	if err != nil {
		log.Printf("Error reading changelog: %v", err)
		return fmt.Errorf("error reading changelog: %v", err)
	}
	lines := strings.Split(string(content), "\n")
	unreleasedIndex := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "## [Unreleased]") {
			unreleasedIndex = i
			break
		}
	}
	if unreleasedIndex == -1 {
		return fmt.Errorf("couldn't find Unreleased section in changelog")
	}
	newLines := append([]string{}, lines[:unreleasedIndex+1]...)
	newLines = append(newLines, "", fmt.Sprintf("## [%s] - %s", version, time.Now().Format("2006-01-02")))
	newLines = append(newLines, lines[unreleasedIndex+1:]...)

	// Update existing diff links and add new ones
	newLines = updateDiffLinks(newLines, version, provider)

	return os.WriteFile(file, []byte(strings.Join(newLines, "\n")), 0644)
}

// AddChangelogSection adds a new section to the Unreleased part of the changelog
func AddChangelogSection(changelogFile, section string) error {
	file, err := os.OpenFile(changelogFile, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening changelog: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	unreleasedFound := false
	sectionAdded := false
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if strings.HasPrefix(line, "## [Unreleased]") {
			unreleasedFound = true
		} else if unreleasedFound && !sectionAdded && line == "" {
			lines = append(lines, fmt.Sprintf("### %s", section), "")
			sectionAdded = true
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning changelog: %w", err)
	}
	if !unreleasedFound {
		return fmt.Errorf("couldn't find Unreleased section in changelog")
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("error seeking to start of file: %w", err)
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error writing to changelog: %w", err)
	}
	return nil
}

func updateDiffLinks(lines []string, newVersion, provider string) []string {
	var updatedLines []string
	var versions []string
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]: ") {
			parts := strings.Split(line, "]:")
			version := strings.Trim(parts[0], "[]")
			if version != "Unreleased" {
				versions = append(versions, version)
			}
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	versions = append([]string{newVersion}, versions...)
	baseURL := getCompareURL(provider)

	updatedLines = append(updatedLines, fmt.Sprintf("[Unreleased]: %s/compare/%s...HEAD", baseURL, versions[0]))
	for i := 0; i < len(versions)-1; i++ {
		updatedLines = append(updatedLines, fmt.Sprintf("[%s]: %s/compare/%s...%s", versions[i], baseURL, versions[i+1], versions[i]))
	}
	updatedLines = append(updatedLines, fmt.Sprintf("[%s]: %s/releases/tag/%s", versions[len(versions)-1], baseURL, versions[len(versions)-1]))

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
