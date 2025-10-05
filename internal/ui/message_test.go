// internal/ui/message_test.go

package ui

import (
	"bytes"
	"errors"
	"testing"
)

func TestPrintColoredMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	err := PrintColoredMessage(buf, "Test Message", "green")
	if err != nil {
		t.Fatalf("PrintColoredMessage returned an error: %v", err)
	}

	output := buf.String()
	expected := "Test Message"
	if !bytes.Contains([]byte(output), []byte(expected)) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}

func TestPrintColoredMessageInvalidColor(t *testing.T) {
	buf := new(bytes.Buffer)
	err := PrintColoredMessage(buf, "Test Message", "invalid-color")
	if err == nil {
		t.Errorf("Expected error for invalid color, got nil")
	}
}

func TestPrintColoredMessageAllColors(t *testing.T) {
	colors := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}

	for _, color := range colors {
		t.Run(color, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := PrintColoredMessage(buf, "Test Message", color)
			if err != nil {
				t.Errorf("PrintColoredMessage with color %q returned error: %v", color, err)
			}

			output := buf.String()
			if !bytes.Contains([]byte(output), []byte("Test Message")) {
				t.Errorf("Expected output to contain message for color %q", color)
			}
		})
	}
}

// errorWriter is a writer that always returns an error
type errorWriter struct{}

func (e *errorWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("write error")
}

func TestPrintColoredMessageWriteError(t *testing.T) {
	writer := &errorWriter{}
	err := PrintColoredMessage(writer, "Test Message", "green")
	if err == nil {
		t.Errorf("Expected error when writer fails, got nil")
	}
}
