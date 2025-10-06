// Package output provides utilities for formatting command output.
//
// It supports both human-readable text output and machine-readable JSON output,
// allowing commands to easily switch between formats based on user preference.
package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/viper"
)

// IsJSONEnabled returns true if JSON output mode is enabled.
func IsJSONEnabled() bool {
	return viper.GetBool("app.json_output")
}

// WriteJSON writes a value as JSON to the given writer.
// It returns an error if JSON marshaling fails.
func WriteJSON(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

// WriteString writes a string to the given writer.
// It returns an error if writing fails.
func WriteString(w io.Writer, s string) error {
	if _, err := fmt.Fprintln(w, s); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	return nil
}

// Write writes output in the appropriate format (JSON or text) based on configuration.
// If JSON mode is enabled, it marshals the jsonValue.
// Otherwise, it writes the textValue as-is.
func Write(w io.Writer, textValue string, jsonValue interface{}) error {
	if IsJSONEnabled() {
		return WriteJSON(w, jsonValue)
	}
	return WriteString(w, textValue)
}

// BumpOutput represents the JSON output structure for version bump commands.
type BumpOutput struct {
	Success       bool   `json:"success"`
	OldVersion    string `json:"old_version,omitempty"`
	NewVersion    string `json:"new_version,omitempty"`
	Tag           string `json:"tag,omitempty"`
	ChangelogFile string `json:"changelog_file,omitempty"`
	CommitHash    string `json:"commit_hash,omitempty"`
	Pushed        bool   `json:"pushed,omitempty"`
	Error         string `json:"error,omitempty"`
	BumpType      string `json:"bump_type,omitempty"`
}

// InitOutput represents the JSON output structure for init command.
type InitOutput struct {
	Success       bool   `json:"success"`
	ChangelogFile string `json:"changelog_file,omitempty"`
	Created       bool   `json:"created"`
	Error         string `json:"error,omitempty"`
}

// ChangelogOutput represents the JSON output structure for changelog add commands.
type ChangelogOutput struct {
	Success       bool   `json:"success"`
	Section       string `json:"section,omitempty"`
	Content       string `json:"content,omitempty"`
	ChangelogFile string `json:"changelog_file,omitempty"`
	Added         bool   `json:"added"`
	Error         string `json:"error,omitempty"`
}

// DocsOutput represents the JSON output structure for docs generation commands.
type DocsOutput struct {
	Success    bool   `json:"success"`
	OutputFile string `json:"output_file,omitempty"`
	Format     string `json:"format,omitempty"`
	Error      string `json:"error,omitempty"`
}
