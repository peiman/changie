package logger

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name     string
		level    zapcore.Level
		logFunc  func(string, ...zap.Field)
		message  string
		fields   []zap.Field
		expected map[string]interface{}
	}{
		{
			name:    "Info logging",
			level:   zapcore.InfoLevel,
			logFunc: Info,
			message: "Test info",
			fields:  []zap.Field{zap.Int("count", 5)},
			expected: map[string]interface{}{
				"level": "INFO",
				"msg":   "Test info",
				"count": float64(5), // JSON numbers are floats
			},
		},
		{
			name:    "Error logging with stack trace",
			level:   zapcore.ErrorLevel,
			logFunc: Error,
			message: "Test error",
			fields:  []zap.Field{zap.Error(assert.AnError)},
			expected: map[string]interface{}{
				"level":      "ERROR",
				"msg":        "Test error",
				"error":      assert.AnError.Error(),
				"stacktrace": "", // We now expect a stacktrace, but don't check its content
			},
		},
		{
			name:     "Debug logging below threshold",
			level:    zapcore.InfoLevel,
			logFunc:  Debug,
			message:  "Test debug",
			fields:   nil,
			expected: nil, // Expect no output
		},
		{
			name:    "Warn logging with complex fields",
			level:   zapcore.WarnLevel,
			logFunc: Warn,
			message: "Test warn",
			fields: []zap.Field{
				zap.String("user", "test_user"),
				zap.Time("timestamp", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				zap.Duration("elapsed", 5*time.Second),
			},
			expected: map[string]interface{}{
				"level":     "WARN",
				"msg":       "Test warn",
				"user":      "test_user",
				"timestamp": "2023-01-01T00:00:00.000Z",
				"elapsed":   float64(5000000000), // nanoseconds
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(tt.level, &buf)

			tt.logFunc(tt.message, tt.fields...)

			if tt.expected == nil {
				if buf.Len() > 0 {
					t.Errorf("Expected no output, but got: %s", buf.String())
				}
				return
			}

			var logEntry map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &logEntry)
			if err != nil {
				t.Fatalf("Failed to parse log output: %v", err)
			}

			for key, expectedValue := range tt.expected {
				if actualValue, ok := logEntry[key]; !ok {
					t.Errorf("Expected key %s not found in log entry", key)
				} else if key == "stacktrace" {
					if actualValue == "" {
						t.Errorf("Expected non-empty stacktrace, but got empty string")
					}
				} else if actualValue != expectedValue {
					t.Errorf("For key %s, expected %v, but got %v", key, expectedValue, actualValue)
				}
			}

			// Check for unexpected fields
			for key := range logEntry {
				if key != "ts" && tt.expected[key] == nil && key != "stacktrace" {
					t.Errorf("Unexpected key %s found in log entry", key)
				}
			}
		})
	}
}

// ... rest of the file remains the same
