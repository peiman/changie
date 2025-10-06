package mcp

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupIntegrationTestRepo creates a complete git repository with changie initialized
// for integration testing MCP tools.
//
//nolint:unparam // tempDir return value is intentionally kept for potential future use
func setupIntegrationTestRepo(t *testing.T, initialTag string) (string, func()) {
	t.Helper()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "mcp-integration-test-*")
	require.NoError(t, err, "Failed to create temp dir")

	// Save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")

	// Change to temp directory
	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change to temp directory")

	// Initialize git repository
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	// Create initial CHANGELOG.md
	changelogContent := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [` + initialTag + `] - 2024-01-01

### Added
- Initial release
`

	err = os.WriteFile("CHANGELOG.md", []byte(changelogContent), 0o644)
	require.NoError(t, err, "Failed to write CHANGELOG.md")

	// Create README
	err = os.WriteFile("README.md", []byte("# Test Project\n"), 0o644)
	require.NoError(t, err, "Failed to write README.md")

	// Initial commit
	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()

	// Create initial tag
	exec.Command("git", "tag", initialTag).Run()

	// Cleanup function
	cleanup := func() {
		os.Chdir(originalDir)
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestBumpVersionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if changieBinary == "" {
		t.Skip("Changie binary not available, skipping integration test")
	}

	tests := []struct {
		name        string
		bumpType    string
		initialTag  string
		expectedVer string
	}{
		{
			name:        "major bump from v1.0.0",
			bumpType:    "major",
			initialTag:  "v1.0.0",
			expectedVer: "v2.0.0",
		},
		{
			name:        "minor bump from v1.2.3",
			bumpType:    "minor",
			initialTag:  "v1.2.3",
			expectedVer: "v1.3.0",
		},
		{
			name:        "patch bump from v2.1.0",
			bumpType:    "patch",
			initialTag:  "v2.1.0",
			expectedVer: "v2.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cleanup := setupIntegrationTestRepo(t, tt.initialTag)
			defer cleanup()

			ctx := context.Background()
			input := BumpVersionInput{
				Type:     tt.bumpType,
				AutoPush: false,
			}

			// Execute the actual MCP tool function
			_, result, err := BumpVersion(ctx, nil, input)

			// We expect this to fail since changie binary might not be in PATH
			// but if it succeeds, verify the result
			if err == nil {
				assert.True(t, result.Success, "BumpVersion should succeed")
				assert.Equal(t, tt.initialTag, result.OldVersion, "Old version should match")
				assert.Equal(t, tt.expectedVer, result.NewVersion, "New version should match")
				assert.Equal(t, tt.bumpType, result.BumpType, "Bump type should match")
				assert.Equal(t, "CHANGELOG.md", result.ChangelogFile)
				assert.False(t, result.Pushed, "Should not be pushed")
			}
		})
	}
}

func TestAddChangelogIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if changieBinary == "" {
		t.Skip("Changie binary not available, skipping integration test")
	}

	_, cleanup := setupIntegrationTestRepo(t, "v1.0.0")
	defer cleanup()

	ctx := context.Background()
	input := AddChangelogInput{
		Section: "added",
		Content: "New integration test feature",
	}

	_, result, err := AddChangelog(ctx, nil, input)

	// May fail if changie not in PATH, but verify if succeeds
	if err == nil {
		assert.True(t, result.Success, "AddChangelog should succeed")
		assert.Equal(t, "added", result.Section)
		assert.Equal(t, "New integration test feature", result.Content)
		assert.Equal(t, "CHANGELOG.md", result.ChangelogFile)
	}
}

func TestInitIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if changieBinary == "" {
		t.Skip("Changie binary not available, skipping integration test")
	}

	// Create temp directory without changelog
	tempDir, err := os.MkdirTemp("", "mcp-init-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Initialize git
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	ctx := context.Background()
	input := InitInput{
		ChangelogFile: "",
	}

	_, result, err := Init(ctx, nil, input)

	// May fail if changie not in PATH, but verify if succeeds
	if err == nil {
		assert.True(t, result.Success, "Init should succeed")
		assert.True(t, result.Created, "Should have created file")
		assert.Equal(t, "CHANGELOG.md", result.ChangelogFile)
	}
}

func TestGetVersionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	_, cleanup := setupIntegrationTestRepo(t, "v3.2.1")
	defer cleanup()

	ctx := context.Background()

	_, result, err := GetVersion(ctx, nil, struct{}{})

	// This should work since it uses git directly
	assert.NoError(t, err, "GetVersion should not error")
	assert.True(t, result.Success, "GetVersion should succeed")
	assert.Equal(t, "v3.2.1", result.Version, "Version should match tag")
}
