// internal/ui/prompt_test.go

package ui

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failWriter is a writer that always returns an error.
type failWriter struct{}

func (f failWriter) Write(_ []byte) (int, error) { return 0, errors.New("write failed") }

func TestAskYesNo(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultYes bool
		want       bool
	}{
		{"yes input default yes", "y\n", true, true},
		{"yes full input", "yes\n", true, true},
		{"no input", "n\n", true, false},
		{"no full input", "no\n", true, false},
		{"empty input default yes", "\n", true, true},
		{"empty input default no", "\n", false, false},
		{"unrecognized input default yes", "maybe\n", true, true},
		{"unrecognized input default no", "maybe\n", false, false},
		{"uppercase YES", "YES\n", true, true},
		{"uppercase NO", "NO\n", false, false},
		{"mixed case Yes", "Yes\n", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			require.NoError(t, err)

			oldStdin := os.Stdin
			os.Stdin = r
			defer func() { os.Stdin = oldStdin }()

			_, err = w.Write([]byte(tt.input))
			require.NoError(t, err)
			w.Close()

			var buf bytes.Buffer
			result, err := AskYesNo("Test?", tt.defaultYes, &buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestAskYesNo_WriteError(t *testing.T) {
	result, err := AskYesNo("Test?", true, failWriter{})
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestAskYesNo_WriteError_DefaultNo(t *testing.T) {
	result, err := AskYesNo("Test?", false, failWriter{})
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestAskYesNo_ReadError(t *testing.T) {
	// Provide a pipe with only write end closed (EOF on read) so ReadString returns error
	r, w, err := os.Pipe()
	require.NoError(t, err)

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// Close write end without writing — reading will return io.EOF
	w.Close()

	var buf bytes.Buffer
	result, err := AskYesNo("Test?", true, &buf)
	assert.NoError(t, err)
	assert.True(t, result) // defaults to defaultYes on read error
}

func TestAskYesNo_PromptContents(t *testing.T) {
	tests := []struct {
		name          string
		defaultYes    bool
		wantIndicator string
	}{
		{"default yes shows Y/n", true, "[Y/n]"},
		{"default no shows y/N", false, "[y/N]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			require.NoError(t, err)

			oldStdin := os.Stdin
			os.Stdin = r
			defer func() { os.Stdin = oldStdin }()

			_, err = w.Write([]byte("\n"))
			require.NoError(t, err)
			w.Close()

			var buf bytes.Buffer
			_, err = AskYesNo("Question?", tt.defaultYes, &buf)
			assert.NoError(t, err)
			assert.Contains(t, buf.String(), tt.wantIndicator)
		})
	}
}

// Ensure failWriter satisfies io.Writer at compile time.
var _ io.Writer = failWriter{}
