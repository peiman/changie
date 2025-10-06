package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/peiman/changie/internal/output"
)

// BumpVersionInput defines the input parameters for the bump version tool.
type BumpVersionInput struct {
	Type     string `json:"type" jsonschema:"enum=major,enum=minor,enum=patch,description=Version bump type"`
	AutoPush bool   `json:"auto_push,omitempty" jsonschema:"description=Automatically push changes to remote"`
}

// BumpVersion implements the changie_bump_version MCP tool.
// It executes the changie bump command with the specified type and returns structured output.
func BumpVersion(ctx context.Context, _ *mcpsdk.CallToolRequest, input BumpVersionInput) (
	*mcpsdk.CallToolResult,
	output.BumpOutput,
	error,
) {
	// Build command arguments
	args := []string{"bump", input.Type, "--json"}
	if input.AutoPush {
		args = append(args, "--auto-push")
	}

	// Execute changie command
	cmd := exec.CommandContext(ctx, "changie", args...)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, output.BumpOutput{
			Success: false,
			Error:   fmt.Sprintf("failed to execute changie: %v - %s", err, string(outputBytes)),
		}, fmt.Errorf("changie bump failed: %w", err)
	}

	// Parse JSON output
	var result output.BumpOutput
	if err := json.Unmarshal(outputBytes, &result); err != nil {
		return nil, output.BumpOutput{
			Success: false,
			Error:   fmt.Sprintf("failed to parse output: %v", err),
		}, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	return nil, result, nil
}

// AddChangelogInput defines the input parameters for adding a changelog entry.
type AddChangelogInput struct {
	Section string `json:"section" jsonschema:"enum=added,enum=changed,enum=deprecated,enum=removed,enum=fixed,enum=security,description=Changelog section type"`
	Content string `json:"content" jsonschema:"description=The changelog entry content"`
}

// AddChangelog implements the changie_add_changelog MCP tool.
// It executes the changie changelog command to add an entry.
func AddChangelog(ctx context.Context, _ *mcpsdk.CallToolRequest, input AddChangelogInput) (
	*mcpsdk.CallToolResult,
	output.ChangelogOutput,
	error,
) {
	// Execute changie changelog command
	cmd := exec.CommandContext(ctx, "changie", "changelog", input.Section, input.Content, "--json")
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, output.ChangelogOutput{
			Success: false,
			Error:   fmt.Sprintf("failed to execute changie: %v - %s", err, string(outputBytes)),
		}, fmt.Errorf("changie changelog failed: %w", err)
	}

	// Parse JSON output
	var result output.ChangelogOutput
	if err := json.Unmarshal(outputBytes, &result); err != nil {
		return nil, output.ChangelogOutput{
			Success: false,
			Error:   fmt.Sprintf("failed to parse output: %v", err),
		}, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	return nil, result, nil
}

// InitInput defines the input parameters for initializing a changie project.
type InitInput struct {
	ChangelogFile string `json:"changelog_file,omitempty" jsonschema:"description=Path to changelog file (default: CHANGELOG.md)"`
}

// Init implements the changie_init MCP tool.
// It executes the changie init command to initialize a new project.
func Init(ctx context.Context, _ *mcpsdk.CallToolRequest, input InitInput) (
	*mcpsdk.CallToolResult,
	output.InitOutput,
	error,
) {
	// Build command arguments
	args := []string{"init", "--json"}
	if input.ChangelogFile != "" {
		args = append(args, "--changelog-file", input.ChangelogFile)
	}

	// Execute changie command
	cmd := exec.CommandContext(ctx, "changie", args...)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, output.InitOutput{
			Success: false,
			Error:   fmt.Sprintf("failed to execute changie: %v - %s", err, string(outputBytes)),
		}, fmt.Errorf("changie init failed: %w", err)
	}

	// Parse JSON output
	var result output.InitOutput
	if err := json.Unmarshal(outputBytes, &result); err != nil {
		return nil, output.InitOutput{
			Success: false,
			Error:   fmt.Sprintf("failed to parse output: %v", err),
		}, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	return nil, result, nil
}

// GetVersionOutput represents the output of the get version tool.
type GetVersionOutput struct {
	Success bool   `json:"success"`
	Version string `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`
}

// GetVersion implements the changie_get_version MCP tool.
// It retrieves the current version from git tags.
func GetVersion(ctx context.Context, _ *mcpsdk.CallToolRequest, _ struct{}) (
	*mcpsdk.CallToolResult,
	GetVersionOutput,
	error,
) {
	// Execute git describe to get current version
	cmd := exec.CommandContext(ctx, "git", "describe", "--tags", "--abbrev=0")
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, GetVersionOutput{
			Success: false,
			Error:   fmt.Sprintf("failed to get version: %v - %s", err, string(outputBytes)),
		}, fmt.Errorf("git describe failed: %w", err)
	}

	version := strings.TrimSpace(string(outputBytes))
	return nil, GetVersionOutput{
		Success: true,
		Version: version,
	}, nil
}
