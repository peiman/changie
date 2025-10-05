// Package logger provides structured logging functionality for the application.
package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Component-specific sub-loggers provide pre-configured loggers with component context.
// This allows easy filtering and tracing of logs by component/package.
//
// Usage:
//   logger.Version.Debug().Str("type", "major").Msg("Starting version bump")
//
// Benefits:
// - Easy filtering: grep for "component":"version" in JSON logs
// - Better traceability in complex operations
// - Structured organization of logs

var (
	// Version logger for version-related operations
	Version zerolog.Logger

	// Changelog logger for changelog operations
	Changelog zerolog.Logger

	// Git logger for git operations
	Git zerolog.Logger

	// Config logger for configuration operations
	Config zerolog.Logger

	// UI logger for user interface operations
	UI zerolog.Logger

	// Docs logger for documentation generation
	Docs zerolog.Logger
)

// InitComponentLoggers initializes all component-specific loggers.
// This should be called after Init() has configured the global logger.
func InitComponentLoggers() {
	Version = log.With().Str("component", "version").Logger()
	Changelog = log.With().Str("component", "changelog").Logger()
	Git = log.With().Str("component", "git").Logger()
	Config = log.With().Str("component", "config").Logger()
	UI = log.With().Str("component", "ui").Logger()
	Docs = log.With().Str("component", "docs").Logger()
}
