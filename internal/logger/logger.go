// Package logger provides structured logging functionality for the application.
//
// This package implements a centralized logging system using zerolog, providing:
// - Consistent structured logging patterns across the application
// - Log level configuration via Viper
// - Conditional JSON or console output based on environment
// - Optional caller information for debugging
// - Formatted console output with timestamps
//
// The logger should be initialized once at application startup through the Init function,
// and then accessed through the zerolog.Logger instances throughout the codebase.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Init initializes the logger with options from Viper configuration.
//
// It configures the global zerolog Logger instance with:
// - Log level from viper.GetString("app.log_level")
// - Log format (JSON or console) from viper.GetString("app.log_format")
// - Optional caller information from viper.GetBool("app.log_caller")
// - Timestamp in RFC3339 format
//
// Log format behavior:
// - "console": Always use human-readable console output
// - "json": Always use JSON structured logging
// - "auto": Use console if output is a TTY, JSON otherwise (default)
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

	// Configure log level
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

	// Determine output format
	logFormat := viper.GetString("app.log_format")
	useConsole := shouldUseConsoleWriter(logFormat, out)

	// Create base logger with appropriate output format
	var logger zerolog.Logger
	if useConsole {
		// Human-readable console output for development
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: time.RFC3339,
		})
	} else {
		// JSON structured logging for production
		logger = zerolog.New(out)
	}

	// Add timestamp
	logger = logger.With().Timestamp().Logger()

	// Optionally add caller information
	if viper.GetBool("app.log_caller") {
		logger = logger.With().Caller().Logger()
	}

	log.Logger = logger

	// Initialize component-specific sub-loggers
	InitComponentLoggers()

	return nil
}

// shouldUseConsoleWriter determines whether to use ConsoleWriter based on format setting and output type.
//
// Parameters:
//   - format: The log format setting ("console", "json", or "auto")
//   - out: The output writer
//
// Returns:
//   - bool: true if ConsoleWriter should be used, false for JSON output
func shouldUseConsoleWriter(format string, out io.Writer) bool {
	switch format {
	case "console":
		return true
	case "json":
		return false
	case "auto":
		// Auto mode: use console if output is a terminal (TTY), JSON otherwise
		if file, ok := out.(*os.File); ok {
			return isatty.IsTerminal(file.Fd()) || isatty.IsCygwinTerminal(file.Fd())
		}
		// For non-file outputs (like test buffers), default to JSON for structured parsing
		return false
	default:
		// Unknown format, default to console for backward compatibility
		return true
	}
}
