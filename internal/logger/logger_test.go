// internal/logger/logger_test.go
package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func TestInit(t *testing.T) {
	buf := new(bytes.Buffer)
	viper.Set("app.log_level", "info")
	if err := Init(buf); err != nil {
		t.Fatalf("Init returned an error: %v", err)
	}
	log.Info().Msg("Test message")

	if !bytes.Contains(buf.Bytes(), []byte("Test message")) {
		t.Errorf("Expected 'Test message' in log output")
	}

	// Test with invalid log level
	viper.Set("app.log_level", "invalid")
	buf.Reset()
	if err := Init(buf); err != nil {
		t.Fatalf("Init returned an error: %v", err)
	}
	log.Info().Msg("Test message with invalid level")

	if !bytes.Contains(buf.Bytes(), []byte("Test message with invalid level")) {
		t.Errorf("Expected 'Test message with invalid level' in log output")
	}

	// Test with 'debug' log level
	viper.Set("app.log_level", "debug")
	buf.Reset()
	if err := Init(buf); err != nil {
		t.Fatalf("Init returned an error: %v", err)
	}
	log.Debug().Msg("Debug message")

	if !bytes.Contains(buf.Bytes(), []byte("Debug message")) {
		t.Errorf("Expected 'Debug message' in log output")
	}
}

func TestInit_ValidLogLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	viper.Set("app.log_level", "debug")

	err := Init(buf)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	log.Debug().Msg("Debug message")
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Debug message")) {
		t.Errorf("Expected 'Debug message' in log output")
	}
}

func TestInit_InvalidLogLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	viper.Set("app.log_level", "invalid")

	err := Init(buf)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	log.Info().Msg("Info message")
	log.Debug().Msg("Debug message")

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Info message")) {
		t.Errorf("Expected 'Info message' in log output")
	}
	if bytes.Contains([]byte(output), []byte("Debug message")) {
		t.Errorf("Did not expect 'Debug message' in log output")
	}
}

func TestInit_NilOutput(t *testing.T) {
	// Save the original os.Stderr
	oldStderr := os.Stderr

	// Create a pipe to capture os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Redirect os.Stderr to the write end of the pipe
	os.Stderr = w

	// Initialize the logger with nil output
	if err := Init(nil); err != nil {
		t.Fatalf("Failed to initialize logger with nil output: %v", err)
	}

	// Log a message to test the output
	log.Info().Msg("Test message to stderr")

	// Close the write end of the pipe and restore os.Stderr
	w.Close()
	os.Stderr = oldStderr

	// Read the captured output from the read end of the pipe
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	// Close the read end of the pipe
	r.Close()

	// Verify that the output contains the test message
	if !bytes.Contains(buf.Bytes(), []byte("Test message to stderr")) {
		t.Errorf("Expected 'Test message to stderr' in output, got '%s'", buf.String())
	}
}

func TestInit_LogFormat_JSON(t *testing.T) {
	buf := new(bytes.Buffer)
	viper.Set("app.log_level", "info")
	viper.Set("app.log_format", "json")
	viper.Set("app.log_caller", false)

	err := Init(buf)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	log.Info().Str("test_key", "test_value").Msg("JSON test message")
	output := buf.String()

	// JSON output should contain structured fields
	if !strings.Contains(output, `"test_key":"test_value"`) {
		t.Errorf("Expected JSON format with structured fields, got: %s", output)
	}
	if !strings.Contains(output, `"message":"JSON test message"`) {
		t.Errorf("Expected JSON format with message field, got: %s", output)
	}
}

func TestInit_LogFormat_Console(t *testing.T) {
	buf := new(bytes.Buffer)
	viper.Set("app.log_level", "info")
	viper.Set("app.log_format", "console")
	viper.Set("app.log_caller", false)

	err := Init(buf)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	log.Info().Msg("Console test message")
	output := buf.String()

	// Console output should be human-readable, not JSON
	if strings.Contains(output, `"message"`) {
		t.Errorf("Expected console format (not JSON), got: %s", output)
	}
	if !strings.Contains(output, "Console test message") {
		t.Errorf("Expected 'Console test message' in output, got: %s", output)
	}
}

func TestInit_LogFormat_Auto(t *testing.T) {
	buf := new(bytes.Buffer)
	viper.Set("app.log_level", "info")
	viper.Set("app.log_format", "auto")
	viper.Set("app.log_caller", false)

	err := Init(buf)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	log.Info().Str("test_key", "test_value").Msg("Auto test message")
	output := buf.String()

	// Auto mode with bytes.Buffer (not a TTY) should use JSON
	if !strings.Contains(output, `"test_key":"test_value"`) {
		t.Errorf("Expected JSON format for non-TTY output, got: %s", output)
	}
}

func TestInit_WithCaller(t *testing.T) {
	buf := new(bytes.Buffer)
	viper.Set("app.log_level", "info")
	viper.Set("app.log_format", "json")
	viper.Set("app.log_caller", true)

	err := Init(buf)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	log.Info().Msg("Test with caller")
	output := buf.String()

	// Should include caller information
	if !strings.Contains(output, `"caller"`) {
		t.Errorf("Expected caller information in output, got: %s", output)
	}
	if !strings.Contains(output, "logger_test.go") {
		t.Errorf("Expected file name in caller info, got: %s", output)
	}
}

func TestShouldUseConsoleWriter(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		isFile   bool
		expected bool
	}{
		{
			name:     "console format",
			format:   "console",
			isFile:   false,
			expected: true,
		},
		{
			name:     "json format",
			format:   "json",
			isFile:   false,
			expected: false,
		},
		{
			name:     "auto format with buffer",
			format:   "auto",
			isFile:   false,
			expected: false,
		},
		{
			name:     "unknown format",
			format:   "unknown",
			isFile:   false,
			expected: true, // defaults to console
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out io.Writer
			if tt.isFile {
				out = os.Stderr
			} else {
				out = new(bytes.Buffer)
			}

			result := shouldUseConsoleWriter(tt.format, out)
			if result != tt.expected {
				t.Errorf("shouldUseConsoleWriter(%q, %T) = %v, want %v",
					tt.format, out, result, tt.expected)
			}
		})
	}
}

func TestInitComponentLoggers(t *testing.T) {
	// Initialize the main logger first
	buf := new(bytes.Buffer)
	viper.Set("app.log_level", "info")
	viper.Set("app.log_format", "json")

	err := Init(buf)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// Test that component loggers are initialized
	Version.Info().Msg("Version test")
	output := buf.String()

	if !strings.Contains(output, `"component":"version"`) {
		t.Errorf("Expected component field in Version logger output, got: %s", output)
	}

	// Reset buffer
	buf.Reset()

	// Test another component logger
	Changelog.Info().Msg("Changelog test")
	output = buf.String()

	if !strings.Contains(output, `"component":"changelog"`) {
		t.Errorf("Expected component field in Changelog logger output, got: %s", output)
	}
}
