package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Init initializes the logger
func Init(level zapcore.Level, output io.Writer) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.NanosDurationEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(output),
		zap.NewAtomicLevelAt(level),
	)

	logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
