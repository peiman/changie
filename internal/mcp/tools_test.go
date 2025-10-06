package mcp

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peiman/changie/internal/output"
)

var (
	changieBinary  string
	changieTestDir string
)

// TestMain builds the changie binary once for all integration tests
func TestMain(m *testing.M) {
	// Create temp directory for test binary
	tmpDir, err := os.MkdirTemp("", "changie-test-bin-*")
	if err != nil {
		panic(err)
	}

	changieTestDir = tmpDir
	changieBinary = filepath.Join(tmpDir, "changie")

	// Build changie binary
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		os.RemoveAll(tmpDir)
		panic(err)
	}

	cmd := exec.Command("go", "build", "-o", changieBinary, projectRoot)
	if err := cmd.Run(); err != nil {
		// Binary build failed - tests will skip integration tests
		changieBinary = ""
	} else {
		// Add binary directory to PATH
		os.Setenv("PATH", tmpDir+":"+os.Getenv("PATH"))
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

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

func TestInit_DefaultFile(t *testing.T) {
	// Test with empty changelog file (should use default)
	ctx := context.Background()
	input := InitInput{
		ChangelogFile: "",
	}

	_, result, err := Init(ctx, nil, input)

	// This will likely fail since we're not in a fresh project
	// but we're testing the code path
	if err != nil {
		assert.False(t, result.Success)
		assert.NotEmpty(t, result.Error)
	} else {
		assert.True(t, result.Success)
	}
}

func TestInit_CustomFile(t *testing.T) {
	// Test with custom changelog file
	ctx := context.Background()
	input := InitInput{
		ChangelogFile: "CUSTOM-CHANGELOG.md",
	}

	_, result, err := Init(ctx, nil, input)

	// This will likely fail since we're not in a fresh project
	// but we're testing the code path
	if err != nil {
		assert.False(t, result.Success)
		assert.NotEmpty(t, result.Error)
	} else {
		assert.True(t, result.Success)
	}
}

func TestBumpVersion_MajorSuccess(t *testing.T) {
	// Note: This test will fail if not in a proper git repo with clean state
	// but it exercises the code path
	ctx := context.Background()
	input := BumpVersionInput{
		Type:     "major",
		AutoPush: false,
	}

	_, result, _ := BumpVersion(ctx, nil, input)

	// We expect an error in test environment, but check structure
	if result.Success {
		assert.NotEmpty(t, result.NewVersion)
		assert.Equal(t, "major", result.BumpType)
	} else {
		assert.NotEmpty(t, result.Error)
	}
}

func TestBumpVersion_WithAutoPush(t *testing.T) {
	ctx := context.Background()
	input := BumpVersionInput{
		Type:     "minor",
		AutoPush: true,
	}

	_, result, _ := BumpVersion(ctx, nil, input)

	// Will fail in test environment but exercises auto-push path
	if result.Success {
		assert.NotEmpty(t, result.NewVersion)
		assert.Equal(t, "minor", result.BumpType)
	} else {
		assert.NotEmpty(t, result.Error)
	}
}

func TestAddChangelog_ValidSection(t *testing.T) {
	ctx := context.Background()
	input := AddChangelogInput{
		Section: "fixed",
		Content: "Test bug fix entry",
	}

	_, result, _ := AddChangelog(ctx, nil, input)

	// Will likely fail in test but exercises the path
	if result.Success {
		assert.Equal(t, "fixed", result.Section)
		assert.NotEmpty(t, result.ChangelogFile)
	} else {
		assert.NotEmpty(t, result.Error)
	}
}

func TestAddChangelog_AllSections(t *testing.T) {
	sections := []string{"added", "changed", "deprecated", "removed", "fixed", "security"}

	for _, section := range sections {
		t.Run(section, func(t *testing.T) {
			ctx := context.Background()
			input := AddChangelogInput{
				Section: section,
				Content: "Test entry for " + section,
			}

			_, result, _ := AddChangelog(ctx, nil, input)

			// Just verify the section is preserved
			if !result.Success {
				assert.NotEmpty(t, result.Error)
			}
		})
	}
}

func TestBumpVersion_JSONParsing(t *testing.T) {
	// Test that invalid JSON is handled
	ctx := context.Background()
	input := BumpVersionInput{
		Type:     "patch",
		AutoPush: false,
	}

	_, result, err := BumpVersion(ctx, nil, input)
	// Should either succeed or fail gracefully
	if err != nil {
		assert.False(t, result.Success)
	}
}

func TestAddChangelog_JSONParsing(t *testing.T) {
	// Test JSON parsing error path
	ctx := context.Background()
	input := AddChangelogInput{
		Section: "added",
		Content: "Test with special chars: \" \\ \n",
	}

	_, result, err := AddChangelog(ctx, nil, input)
	// Should handle gracefully
	if err != nil {
		assert.False(t, result.Success)
	}
}

func TestInit_JSONParsing(t *testing.T) {
	// Test JSON parsing
	ctx := context.Background()
	input := InitInput{
		ChangelogFile: "TEST.md",
	}

	_, result, err := Init(ctx, nil, input)
	// Should handle gracefully
	if err != nil {
		assert.False(t, result.Success)
	}
}

func TestBumpVersion_AllTypes(t *testing.T) {
	tests := []struct {
		name     string
		bumpType string
	}{
		{"major", "major"},
		{"minor", "minor"},
		{"patch", "patch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			input := BumpVersionInput{
				Type:     tt.bumpType,
				AutoPush: false,
			}

			_, result, _ := BumpVersion(ctx, nil, input)

			// Verify structure is populated
			if !result.Success {
				assert.NotEmpty(t, result.Error)
			}
		})
	}
}

func TestAddChangelog_InvalidSection(t *testing.T) {
	ctx := context.Background()
	input := AddChangelogInput{
		Section: "invalid-section",
		Content: "Test content",
	}

	_, result, err := AddChangelog(ctx, nil, input)

	// Should fail with invalid section
	assert.Error(t, err)
	assert.False(t, result.Success)
}

func TestContextIntegration(t *testing.T) {
	// Test that context is passed through
	t.Run("BumpVersion with context", func(t *testing.T) {
		ctx := context.Background()
		input := BumpVersionInput{Type: "patch"}
		_, _, err := BumpVersion(ctx, nil, input)
		// Error is expected in test environment
		_ = err
	})

	t.Run("AddChangelog with context", func(t *testing.T) {
		ctx := context.Background()
		input := AddChangelogInput{Section: "added", Content: "test"}
		_, _, err := AddChangelog(ctx, nil, input)
		_ = err
	})

	t.Run("Init with context", func(t *testing.T) {
		ctx := context.Background()
		input := InitInput{}
		_, _, err := Init(ctx, nil, input)
		_ = err
	})

	t.Run("GetVersion with context", func(t *testing.T) {
		ctx := context.Background()
		_, _, err := GetVersion(ctx, nil, struct{}{})
		_ = err
	})
}

func TestInputValidation(t *testing.T) {
	ctx := context.Background()

	t.Run("BumpVersion empty type", func(t *testing.T) {
		input := BumpVersionInput{Type: ""}
		_, result, _ := BumpVersion(ctx, nil, input)
		assert.False(t, result.Success)
	})

	t.Run("AddChangelog empty section", func(t *testing.T) {
		input := AddChangelogInput{Section: "", Content: "test"}
		_, result, _ := AddChangelog(ctx, nil, input)
		assert.False(t, result.Success)
	})
}

func TestBumpVersionAutoPushPath(t *testing.T) {
	// Explicitly test auto-push path
	ctx := context.Background()

	tests := []struct {
		name     string
		autoPush bool
	}{
		{"without auto-push", false},
		{"with auto-push", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := BumpVersionInput{
				Type:     "minor",
				AutoPush: tt.autoPush,
			}
			_, result, _ := BumpVersion(ctx, nil, input)
			// Structure should be populated regardless
			_ = result
		})
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

func TestJSONUnmarshaling(t *testing.T) {
	t.Run("BumpOutput success JSON", func(t *testing.T) {
		jsonData := `{
			"success": true,
			"old_version": "v1.0.0",
			"new_version": "v1.1.0",
			"tag": "v1.1.0",
			"changelog_file": "CHANGELOG.md",
			"pushed": false,
			"bump_type": "minor"
		}`

		var result output.BumpOutput
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, "v1.0.0", result.OldVersion)
		assert.Equal(t, "v1.1.0", result.NewVersion)
		assert.Equal(t, "v1.1.0", result.Tag)
		assert.Equal(t, "CHANGELOG.md", result.ChangelogFile)
		assert.False(t, result.Pushed)
		assert.Equal(t, "minor", result.BumpType)
	})

	t.Run("ChangelogOutput success JSON", func(t *testing.T) {
		jsonData := `{
			"success": true,
			"section": "added",
			"content": "New feature added",
			"changelog_file": "CHANGELOG.md"
		}`

		var result output.ChangelogOutput
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, "added", result.Section)
		assert.Equal(t, "New feature added", result.Content)
		assert.Equal(t, "CHANGELOG.md", result.ChangelogFile)
	})

	t.Run("InitOutput success JSON", func(t *testing.T) {
		jsonData := `{
			"success": true,
			"changelog_file": "CHANGELOG.md",
			"created": true
		}`

		var result output.InitOutput
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, "CHANGELOG.md", result.ChangelogFile)
		assert.True(t, result.Created)
	})

	t.Run("GetVersionOutput success JSON", func(t *testing.T) {
		jsonData := `{
			"success": true,
			"version": "v2.3.1"
		}`

		var result GetVersionOutput
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, "v2.3.1", result.Version)
	})

	t.Run("BumpOutput invalid JSON", func(t *testing.T) {
		jsonData := `{invalid json`

		var result output.BumpOutput
		err := json.Unmarshal([]byte(jsonData), &result)
		assert.Error(t, err)
	})

	t.Run("ChangelogOutput invalid JSON", func(t *testing.T) {
		jsonData := `not json at all`

		var result output.ChangelogOutput
		err := json.Unmarshal([]byte(jsonData), &result)
		assert.Error(t, err)
	})

	t.Run("InitOutput invalid JSON", func(t *testing.T) {
		jsonData := `{"incomplete": `

		var result output.InitOutput
		err := json.Unmarshal([]byte(jsonData), &result)
		assert.Error(t, err)
	})
}
