// Package config provides configuration management utilities.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// PathsConfig holds standard paths and filenames for config files.
type PathsConfig struct {
	// Default config name with dot prefix (e.g. ".changie")
	DefaultName string

	// Config file extension
	Extension string

	// Default full config name (e.g. ".changie.yaml")
	DefaultFullName string

	// Default config file with home directory (e.g. "$HOME/.changie.yaml")
	DefaultPath string

	// Default ignore pattern for gitignore (e.g. "changie.yaml")
	IgnorePattern string
}

// DefaultPaths returns standard paths and filenames for config files based on the binary name.
//
// This function generates standardized configuration file paths and names following
// common conventions for CLI applications:
// - Config files are prefixed with a dot (e.g., .changie.yaml)
// - Default location is in the user's home directory
// - Provides gitignore patterns without the leading dot
//
// Parameters:
//   - binaryName: The name of the binary/application
//
// Returns:
//   - PathsConfig: A struct containing all standard config paths and names
//
// Example:
//
//	paths := DefaultPaths("changie")
//	// paths.DefaultName = ".changie"
//	// paths.Extension = "yaml"
//	// paths.DefaultFullName = ".changie.yaml"
//	// paths.DefaultPath = "$HOME/.changie.yaml"
//	// paths.IgnorePattern = "changie.yaml"
func DefaultPaths(binaryName string) PathsConfig {
	ext := "yaml"
	defaultName := fmt.Sprintf(".%s", binaryName)
	defaultFullName := fmt.Sprintf("%s.%s", defaultName, ext)

	home, err := os.UserHomeDir()
	defaultPath := defaultFullName // Fallback if home dir not available
	if err == nil {
		defaultPath = filepath.Join(home, defaultFullName)
	}

	// Used for .gitignore - without leading dot
	ignorePattern := fmt.Sprintf("%s.%s", binaryName, ext)

	return PathsConfig{
		DefaultName:     defaultName,
		Extension:       ext,
		DefaultFullName: defaultFullName,
		DefaultPath:     defaultPath,
		IgnorePattern:   ignorePattern,
	}
}
