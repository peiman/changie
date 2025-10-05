// Package git provides functionality for working with Git repositories.
//
// This package encapsulates Git operations used by the changie tool, including:
// - Version detection from Git tags
// - Checking for uncommitted changes
// - Managing commits and tags for changelog updates
// - Pushing changes to remote repositories
//
// All functions in this package use the git command-line tool and require it to be
// installed and available in the PATH. Functions will return appropriate errors if
// Git operations fail.
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsInstalled checks if git is installed and available in the PATH.
//
// This function attempts to run "git --version" and returns true only if
// the command executes successfully, indicating that Git is properly installed.
//
// Returns:
//   - bool: true if git is installed and available, false otherwise
func IsInstalled() bool {
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	return err == nil
}

// GetVersion returns the current version from git tags.
//
// This function uses "git describe --tags --abbrev=0" to retrieve the most recent tag.
// If no tags exist in the repository, it returns an empty string, which the caller
// should interpret as version 0.0.0.
//
// Returns:
//   - string: The version string (without 'v' prefix) or empty string if no tags exist
//   - error: Any error encountered during the git operation (except for "no tags" errors)
func GetVersion() (string, error) {
	// Use git describe to get the latest tag
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.CombinedOutput()
	// If there are no tags, return empty string which will be treated as 0.0.0
	if err != nil {
		// Check if error is due to no tags
		if strings.Contains(string(output), "No names found") || strings.Contains(string(output), "fatal: No names found") {
			return "", nil
		}
		return "", fmt.Errorf("failed to get latest tag: %w (verify that you are in a git repository and have sufficient permissions)", err)
	}

	// Trim and return the version string
	version := strings.TrimSpace(string(output))
	return version, nil
}

// GetCurrentBranch returns the name of the current git branch.
//
// This function uses "git rev-parse --abbrev-ref HEAD" to get the current branch name.
// This is a reliable way to get the branch name that works in all scenarios.
//
// Returns:
//   - string: The name of the current branch
//   - error: Any error encountered during the git operation
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w (verify you're in a git repository)", err)
	}

	branch := strings.TrimSpace(string(output))
	return branch, nil
}

// HasUncommittedChanges checks if there are any uncommitted changes in the repository.
//
// This function uses "git status --porcelain" to determine if there are any staged
// or unstaged changes that haven't been committed. The porcelain format ensures
// machine-readable output that's stable across git versions.
//
// Returns:
//   - bool: true if there are uncommitted changes, false otherwise
//   - error: Any error encountered during the git operation
func HasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check for uncommitted changes: %w (verify you're in a git repository with proper permissions)", err)
	}

	// If there is any output, there are uncommitted changes
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// CommitChangelog commits changes to the changelog file.
//
// This function:
// 1. Adds the specified changelog file to the staging area using "git add"
// 2. Commits the changes with a standardized message that includes the version
//
// Parameters:
//   - file: Path to the changelog file to commit
//   - version: Version number to include in the commit message
//
// Returns:
//   - error: Any error encountered during the git operations
func CommitChangelog(file, version string) error {
	// Add the file to the staging area
	cmd := exec.Command("git", "add", file)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add changelog to staging area: %w (check if the file '%s' exists and you have write permissions)", err, file)
	}

	// Commit the changes
	commitMsg := fmt.Sprintf("Release %s", version)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit changelog: %w (ensure git user.name and user.email are configured correctly with 'git config')", err)
	}

	return nil
}

// TagVersion creates a new git tag for the given version.
//
// This function creates an annotated tag (-a) with a standardized message
// that includes the version number. It respects the user's configured
// preference for using 'v' prefix or not.
//
// Parameters:
//   - version: Version number to tag (with or without "v" prefix)
//
// Returns:
//   - error: Any error encountered during the git operation
func TagVersion(version string) error {
	// The version parameter already has the correct prefix (or no prefix)
	// based on the user's preference, so we use it directly
	tagName := version

	// Create the tag
	tagMsg := fmt.Sprintf("Version %s", version)
	cmd := exec.Command("git", "tag", "-a", tagName, "-m", tagMsg)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create tag: %w (check if tag '%s' already exists, you can delete it with 'git tag -d %s')", err, tagName, tagName)
	}

	return nil
}

// PushChanges pushes commits and tags to the remote repository.
//
// This function performs two git operations:
// 1. "git push" to push commits to the remote repository
// 2. "git push --tags" to push tags to the remote repository
//
// Returns:
//   - error: Any error encountered during the git operations
func PushChanges() error {
	// Push commits
	cmd := exec.Command("git", "push")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push commits: %w (check network connection and remote repository access permissions)", err)
	}

	// Push tags
	cmd = exec.Command("git", "push", "--tags")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to push tags: %w (check if you have permission to create tags on the remote repository)", err)
	}

	return nil
}

// RepositoryInfo holds parsed repository information
type RepositoryInfo struct {
	Owner    string // Repository owner/organization
	Repo     string // Repository name
	Provider string // Provider name (github, bitbucket, gitlab, etc.)
	BaseURL  string // Base URL for the provider
}

// GetRemoteURL returns the remote origin URL from git config.
//
// This function uses "git config --get remote.origin.url" to retrieve the
// remote repository URL.
//
// Returns:
//   - string: The remote origin URL
//   - error: Any error encountered during the git operation
func GetRemoteURL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w (verify you have a remote repository configured with 'git remote -v')", err)
	}

	url := strings.TrimSpace(string(output))
	if url == "" {
		return "", fmt.Errorf("remote URL is empty (configure a remote repository with 'git remote add origin <url>')")
	}

	return url, nil
}

// ParseRepositoryURL parses a git remote URL and extracts repository information.
//
// This function handles both HTTPS and SSH URL formats:
//   - HTTPS: https://github.com/owner/repo.git
//   - SSH: git@github.com:owner/repo.git
//
// Parameters:
//   - remoteURL: The git remote URL to parse
//
// Returns:
//   - *RepositoryInfo: Parsed repository information
//   - error: Any error encountered during parsing
func ParseRepositoryURL(remoteURL string) (*RepositoryInfo, error) {
	remoteURL = strings.TrimSpace(remoteURL)
	if remoteURL == "" {
		return nil, fmt.Errorf("remote URL is empty")
	}

	var owner, repo, provider, baseURL string

	// Handle SSH format: git@github.com:owner/repo.git
	if strings.HasPrefix(remoteURL, "git@") {
		parts := strings.Split(remoteURL, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid SSH URL format: %s (expected format: git@host:owner/repo.git)", remoteURL)
		}

		// Extract host
		host := strings.TrimPrefix(parts[0], "git@")

		// Extract owner/repo
		path := strings.TrimSuffix(parts[1], ".git")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			return nil, fmt.Errorf("invalid repository path: %s (expected format: owner/repo)", path)
		}

		owner = pathParts[len(pathParts)-2]
		repo = pathParts[len(pathParts)-1]

		// Determine provider from host
		switch {
		case strings.Contains(host, "github.com"):
			provider = "github"
			baseURL = "https://github.com"
		case strings.Contains(host, "bitbucket.org"):
			provider = "bitbucket"
			baseURL = "https://bitbucket.org"
		case strings.Contains(host, "gitlab.com"):
			provider = "gitlab"
			baseURL = "https://gitlab.com"
		default:
			provider = "unknown"
			baseURL = "https://" + host
		}
	} else if strings.HasPrefix(remoteURL, "http://") || strings.HasPrefix(remoteURL, "https://") {
		// Handle HTTPS format: https://github.com/owner/repo.git
		// Remove scheme
		url := strings.TrimPrefix(remoteURL, "https://")
		url = strings.TrimPrefix(url, "http://")

		// Remove .git suffix
		url = strings.TrimSuffix(url, ".git")

		// Split into parts
		parts := strings.Split(url, "/")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid HTTPS URL format: %s (expected format: https://host/owner/repo)", remoteURL)
		}

		host := parts[0]
		owner = parts[len(parts)-2]
		repo = parts[len(parts)-1]

		// Determine provider from host
		switch {
		case strings.Contains(host, "github.com"):
			provider = "github"
			baseURL = "https://github.com"
		case strings.Contains(host, "bitbucket.org"):
			provider = "bitbucket"
			baseURL = "https://bitbucket.org"
		case strings.Contains(host, "gitlab.com"):
			provider = "gitlab"
			baseURL = "https://gitlab.com"
		default:
			provider = "unknown"
			baseURL = "https://" + host
		}
	} else {
		return nil, fmt.Errorf("unsupported URL format: %s (expected SSH or HTTPS format)", remoteURL)
	}

	return &RepositoryInfo{
		Owner:    owner,
		Repo:     repo,
		Provider: provider,
		BaseURL:  baseURL,
	}, nil
}

// GetRepositoryInfo gets repository information from git remote configuration.
//
// This is a convenience function that combines GetRemoteURL and ParseRepositoryURL.
//
// Returns:
//   - *RepositoryInfo: Parsed repository information
//   - error: Any error encountered during the operation
func GetRepositoryInfo() (*RepositoryInfo, error) {
	url, err := GetRemoteURL()
	if err != nil {
		return nil, err
	}

	return ParseRepositoryURL(url)
}
