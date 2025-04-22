// Package logger provides structured logging functionality for the application.
//
// This package implements a centralized logging system using zerolog, providing:
// - Consistent structured logging patterns across the application
// - Log level configuration via Viper
// - Formatted console output with timestamps
//
// The logger should be initialized once at application startup through the Init function,
// and then accessed through the zerolog.Logger instances throughout the codebase.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Init initializes the logger with options from Viper configuration.
//
// It configures the global zerolog Logger instance with:
// - Log level from viper.GetString("app.log_level")
// - Console writer output formatting
// - Timestamp in RFC3339 format
//
// This function should be called once in rootCmd's PersistentPreRunE or main initialization.
// If an invalid log level is provided, it will default to "info" and log a warning.
//
// Parameters:
//   - out: Writer where logs will be written. If nil, defaults to os.Stderr
//
// Returns:
//   - error: Any error encountered during initialization
func Init(out io.Writer) error {
	if out == nil {
		out = os.Stderr
	}

	logLevelStr := viper.GetString("app.log_level")
	level, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		level = zerolog.InfoLevel
		log.Warn().
			Err(err).
			Str("provided_level", logLevelStr).
			Msg("Invalid log level provided, defaulting to 'info'")
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Logger()

	return nil
}
