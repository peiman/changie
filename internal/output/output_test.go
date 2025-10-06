package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsJSONEnabled(t *testing.T) {
	// Save original value
	originalValue := viper.Get("app.json_output")
	defer func() {
		if originalValue != nil {
			viper.Set("app.json_output", originalValue)
		}
	}()

	tests := []struct {
		name     string
		setValue bool
		want     bool
	}{
		{
			name:     "JSON disabled",
			setValue: false,
			want:     false,
		},
		{
			name:     "JSON enabled",
			setValue: true,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("app.json_output", tt.setValue)
			got := IsJSONEnabled()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name: "simple struct",
			value: BumpOutput{
				Success:    true,
				OldVersion: "1.0.0",
				NewVersion: "1.1.0",
				BumpType:   "minor",
			},
			wantErr: false,
			check: func(t *testing.T, output string) {
				var result BumpOutput
				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err)
				assert.True(t, result.Success)
				assert.Equal(t, "1.0.0", result.OldVersion)
				assert.Equal(t, "1.1.0", result.NewVersion)
				assert.Equal(t, "minor", result.BumpType)
			},
		},
		{
			name: "with error message",
			value: BumpOutput{
				Success: false,
				Error:   "something went wrong",
			},
			wantErr: false,
			check: func(t *testing.T, output string) {
				var result BumpOutput
				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err)
				assert.False(t, result.Success)
				assert.Equal(t, "something went wrong", result.Error)
			},
		},
		{
			name: "omitempty fields",
			value: InitOutput{
				Success: true,
				Created: true,
			},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// ChangelogFile should be omitted
				assert.NotContains(t, output, "changelog_file")
				assert.Contains(t, output, "success")
				assert.Contains(t, output, "created")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteJSON(&buf, tt.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.check != nil {
					tt.check(t, buf.String())
				}
			}
		})
	}
}

func TestWriteString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "simple string",
			input:   "Hello, World!",
			want:    "Hello, World!\n",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "\n",
			wantErr: false,
		},
		{
			name:    "multiline string",
			input:   "Line 1\nLine 2",
			want:    "Line 1\nLine 2\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteString(&buf, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, buf.String())
			}
		})
	}
}

func TestWrite(t *testing.T) {
	// Save original value
	originalValue := viper.Get("app.json_output")
	defer func() {
		if originalValue != nil {
			viper.Set("app.json_output", originalValue)
		}
	}()

	tests := []struct {
		name      string
		jsonMode  bool
		textValue string
		jsonValue interface{}
		checkText func(t *testing.T, output string)
		checkJSON func(t *testing.T, output string)
	}{
		{
			name:      "text mode",
			jsonMode:  false,
			textValue: "Success message",
			jsonValue: BumpOutput{Success: true},
			checkText: func(t *testing.T, output string) {
				assert.Equal(t, "Success message\n", output)
			},
		},
		{
			name:      "JSON mode",
			jsonMode:  true,
			textValue: "Success message",
			jsonValue: BumpOutput{Success: true, NewVersion: "1.2.3"},
			checkJSON: func(t *testing.T, output string) {
				var result BumpOutput
				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err)
				assert.True(t, result.Success)
				assert.Equal(t, "1.2.3", result.NewVersion)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("app.json_output", tt.jsonMode)

			var buf bytes.Buffer
			err := Write(&buf, tt.textValue, tt.jsonValue)
			assert.NoError(t, err)

			if tt.jsonMode && tt.checkJSON != nil {
				tt.checkJSON(t, buf.String())
			} else if !tt.jsonMode && tt.checkText != nil {
				tt.checkText(t, buf.String())
			}
		})
	}
}

func TestBumpOutput(t *testing.T) {
	output := BumpOutput{
		Success:       true,
		OldVersion:    "1.0.0",
		NewVersion:    "2.0.0",
		Tag:           "v2.0.0",
		ChangelogFile: "CHANGELOG.md",
		CommitHash:    "abc123",
		Pushed:        true,
		BumpType:      "major",
	}

	var buf bytes.Buffer
	err := WriteJSON(&buf, output)
	require.NoError(t, err)

	var result BumpOutput
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, output.Success, result.Success)
	assert.Equal(t, output.OldVersion, result.OldVersion)
	assert.Equal(t, output.NewVersion, result.NewVersion)
	assert.Equal(t, output.Tag, result.Tag)
	assert.Equal(t, output.ChangelogFile, result.ChangelogFile)
	assert.Equal(t, output.Pushed, result.Pushed)
	assert.Equal(t, output.BumpType, result.BumpType)
}

func TestInitOutput(t *testing.T) {
	output := InitOutput{
		Success:       true,
		ChangelogFile: "CHANGELOG.md",
		Created:       true,
	}

	var buf bytes.Buffer
	err := WriteJSON(&buf, output)
	require.NoError(t, err)

	var result InitOutput
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, output.Success, result.Success)
	assert.Equal(t, output.ChangelogFile, result.ChangelogFile)
	assert.Equal(t, output.Created, result.Created)
}

func TestChangelogOutput(t *testing.T) {
	output := ChangelogOutput{
		Success:       true,
		Section:       "Added",
		Content:       "New feature",
		ChangelogFile: "CHANGELOG.md",
		Added:         true,
	}

	var buf bytes.Buffer
	err := WriteJSON(&buf, output)
	require.NoError(t, err)

	var result ChangelogOutput
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, output.Success, result.Success)
	assert.Equal(t, output.Section, result.Section)
	assert.Equal(t, output.Content, result.Content)
	assert.Equal(t, output.ChangelogFile, result.ChangelogFile)
	assert.Equal(t, output.Added, result.Added)
}

func TestDocsOutput(t *testing.T) {
	output := DocsOutput{
		Success:    true,
		OutputFile: "config.md",
		Format:     "markdown",
	}

	var buf bytes.Buffer
	err := WriteJSON(&buf, output)
	require.NoError(t, err)

	var result DocsOutput
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, output.Success, result.Success)
	assert.Equal(t, output.OutputFile, result.OutputFile)
	assert.Equal(t, output.Format, result.Format)
}

func TestJSONIndentation(t *testing.T) {
	output := BumpOutput{
		Success:    true,
		NewVersion: "1.0.0",
	}

	var buf bytes.Buffer
	err := WriteJSON(&buf, output)
	require.NoError(t, err)

	// Check that output is indented (has newlines and spaces)
	result := buf.String()
	assert.True(t, strings.Contains(result, "\n"))
	assert.True(t, strings.Contains(result, "  ")) // 2-space indent
}
