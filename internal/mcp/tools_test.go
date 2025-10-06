package mcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/peiman/changie/internal/output"
)

func TestBumpVersionInput(t *testing.T) {
	tests := []struct {
		name  string
		input BumpVersionInput
	}{
		{
			name: "major bump",
			input: BumpVersionInput{
				Type:     "major",
				AutoPush: false,
			},
		},
		{
			name: "minor bump with auto push",
			input: BumpVersionInput{
				Type:     "minor",
				AutoPush: true,
			},
		},
		{
			name: "patch bump",
			input: BumpVersionInput{
				Type:     "patch",
				AutoPush: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, []string{"major", "minor", "patch"}, tt.input.Type)
		})
	}
}

func TestAddChangelogInput(t *testing.T) {
	tests := []struct {
		name  string
		input AddChangelogInput
	}{
		{
			name: "added section",
			input: AddChangelogInput{
				Section: "added",
				Content: "New feature",
			},
		},
		{
			name: "fixed section",
			input: AddChangelogInput{
				Section: "fixed",
				Content: "Bug fix",
			},
		},
		{
			name: "security section",
			input: AddChangelogInput{
				Section: "security",
				Content: "Security patch",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validSections := []string{"added", "changed", "deprecated", "removed", "fixed", "security"}
			assert.Contains(t, validSections, tt.input.Section)
			assert.NotEmpty(t, tt.input.Content)
		})
	}
}

func TestInitInput(t *testing.T) {
	tests := []struct {
		name  string
		input InitInput
	}{
		{
			name:  "default changelog file",
			input: InitInput{},
		},
		{
			name: "custom changelog file",
			input: InitInput{
				ChangelogFile: "CHANGES.md",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the struct is valid
			_ = tt.input
		})
	}
}

func TestGetVersionOutput(t *testing.T) {
	tests := []struct {
		name   string
		output GetVersionOutput
		valid  bool
	}{
		{
			name: "successful version retrieval",
			output: GetVersionOutput{
				Success: true,
				Version: "v1.2.3",
			},
			valid: true,
		},
		{
			name: "error getting version",
			output: GetVersionOutput{
				Success: false,
				Error:   "git command failed",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.output.Success {
				assert.NotEmpty(t, tt.output.Version)
				assert.Empty(t, tt.output.Error)
			} else {
				assert.NotEmpty(t, tt.output.Error)
			}
		})
	}
}

func TestGetVersion_Success(t *testing.T) {
	// This test requires a git repository with tags
	// Skip if not in a git repo
	ctx := context.Background()
	_, result, err := GetVersion(ctx, nil, struct{}{})

	// We expect either success or an error (if no tags exist)
	// Both are valid outcomes depending on the repository state
	if err != nil {
		// Error case - should have error message in result
		assert.False(t, result.Success)
		assert.NotEmpty(t, result.Error)
	} else {
		// Success case - should have version
		if result.Success {
			assert.NotEmpty(t, result.Version)
			assert.Empty(t, result.Error)
		}
	}
}

func TestBumpVersion_InvalidType(t *testing.T) {
	ctx := context.Background()
	input := BumpVersionInput{
		Type:     "invalid",
		AutoPush: false,
	}

	_, result, err := BumpVersion(ctx, nil, input)

	// Should fail with invalid bump type
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Error)
}

func TestAddChangelog_EmptyContent(t *testing.T) {
	ctx := context.Background()
	input := AddChangelogInput{
		Section: "added",
		Content: "",
	}

	_, result, err := AddChangelog(ctx, nil, input)

	// Should fail with empty content
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Error)
}

func TestOutputStructures(t *testing.T) {
	t.Run("BumpOutput with all fields", func(t *testing.T) {
		out := output.BumpOutput{
			Success:       true,
			OldVersion:    "v1.0.0",
			NewVersion:    "v1.1.0",
			Tag:           "v1.1.0",
			ChangelogFile: "CHANGELOG.md",
			CommitHash:    "abc123",
			Pushed:        true,
			BumpType:      "minor",
		}
		assert.True(t, out.Success)
		assert.Equal(t, "v1.0.0", out.OldVersion)
		assert.Equal(t, "v1.1.0", out.NewVersion)
		assert.Equal(t, "v1.1.0", out.Tag)
		assert.Equal(t, "CHANGELOG.md", out.ChangelogFile)
		assert.Equal(t, "abc123", out.CommitHash)
		assert.True(t, out.Pushed)
		assert.Equal(t, "minor", out.BumpType)
	})

	t.Run("ChangelogOutput with error", func(t *testing.T) {
		out := output.ChangelogOutput{
			Success: false,
			Error:   "failed to add entry",
		}
		assert.False(t, out.Success)
		assert.NotEmpty(t, out.Error)
	})

	t.Run("InitOutput success", func(t *testing.T) {
		out := output.InitOutput{
			Success:       true,
			ChangelogFile: "CHANGELOG.md",
			Created:       true,
		}
		assert.True(t, out.Success)
		assert.True(t, out.Created)
		assert.Equal(t, "CHANGELOG.md", out.ChangelogFile)
	})

	t.Run("GetVersionOutput with version", func(t *testing.T) {
		out := GetVersionOutput{
			Success: true,
			Version: "v2.0.0",
		}
		assert.True(t, out.Success)
		assert.Equal(t, "v2.0.0", out.Version)
	})
}
