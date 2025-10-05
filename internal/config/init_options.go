// internal/config/init_options.go
//
// Init command configuration options
//
// This file contains configuration options specific to the 'init' command.

package config

// InitOptions returns configuration options for the init command
func InitOptions() []ConfigOption {
	return []ConfigOption{
		{
			Key:          "app.changelog.file",
			DefaultValue: "CHANGELOG.md",
			Description:  "Name of the changelog file to create and manage",
			Type:         "string",
			Required:     false,
			Example:      "HISTORY.md",
		},
		{
			Key:          "app.version.use_v_prefix",
			DefaultValue: true,
			Description:  "Use 'v' prefix for version tags (e.g., v1.0.0 vs 1.0.0)",
			Type:         "bool",
			Required:     false,
			Example:      "false",
		},
	}
}
