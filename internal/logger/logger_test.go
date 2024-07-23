package logger

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testLogSink struct {
	entries []zapcore.Entry
	fields  []map[string]interface{}
	mu      sync.Mutex
}

func (s *testLogSink) Write(p []byte) (int, error) {
	var entry map[string]interface{}
	err := json.Unmarshal(p, &entry)
	if err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fields = append(s.fields, entry)
	return len(p), nil
}

func (s *testLogSink) Sync() error {
	return nil
}

func newTestLogger() (*zap.Logger, *testLogSink) {
	sink := &testLogSink{}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.NanosDurationEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(sink),
		zapcore.DebugLevel,
	)
	return zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel)), sink
}

func TestLogger(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		logFunc  func(string, ...zap.Field)
		message  string
		fields   []zap.Field
		expected map[string]interface{}
	}{
		{
			name:    "Info logging",
			level:   INFO,
			logFunc: Info,
			message: "Test info",
			fields:  []zap.Field{zap.Int("count", 5)},
			expected: map[string]interface{}{
				"level": "INFO",
				"msg":   "Test info",
				"count": float64(5),
			},
		},
		{
			name:    "Error logging with stack trace",
			level:   ERROR,
			logFunc: Error,
			message: "Test error",
			fields:  []zap.Field{zap.Error(assert.AnError)},
			expected: map[string]interface{}{
				"level": "ERROR",
				"msg":   "Test error",
				"error": assert.AnError.Error(),
			},
		},
		{
			name:    "Debug logging",
			level:   DEBUG,
			logFunc: Debug,
			message: "Test debug",
			fields:  nil,
			expected: map[string]interface{}{
				"level": "DEBUG",
				"msg":   "Test debug",
			},
		},
		{
			name:    "Warn logging with complex fields",
			level:   WARN,
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
			testLogger, sink := newTestLogger()
			logger = testLogger

			tt.logFunc(tt.message, tt.fields...)

			if len(sink.fields) == 0 {
				t.Fatalf("No log entries recorded")
			}

			entry := sink.fields[0]

			for key, expectedValue := range tt.expected {
				if actualValue, ok := entry[key]; !ok {
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
			for key := range entry {
				if tt.expected[key] == nil && key != "stacktrace" && key != "ts" {
					t.Errorf("Unexpected key %s found in log entry", key)
				}
			}
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	testLogger, sink := newTestLogger()
	logger = testLogger

	Debug("This is a debug message")
	Info("This is an info message")
	Warn("This is a warning message")
	Error("This is an error message")

	expectedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	expectedMessages := []string{
		"This is a debug message",
		"This is an info message",
		"This is a warning message",
		"This is an error message",
	}

	if len(sink.fields) != len(expectedLevels) {
		t.Errorf("Expected %d log entries, but got %d", len(expectedLevels), len(sink.fields))
	}

	for i, entry := range sink.fields {
		if entry["level"] != expectedLevels[i] {
			t.Errorf("Expected level %v, but got %v", expectedLevels[i], entry["level"])
		}
		if entry["msg"] != expectedMessages[i] {
			t.Errorf("Expected message '%s', but got '%s'", expectedMessages[i], entry["msg"])
		}
	}
}
