// internal/config/version_options.go
//
// Version command configuration options
//
// This file contains configuration options specific to version bump commands
// (major, minor, patch).

package config

// VersionOptions returns configuration options for version bump commands
func VersionOptions() []ConfigOption {
	return []ConfigOption{
		{
			Key:          "app.version.tag_prefix",
			DefaultValue: "v",
			Description:  "Prefix for git version tags",
			Type:         "string",
			Required:     false,
			Example:      "",
		},
		{
			Key:          "app.version.auto_push",
			DefaultValue: false,
			Description:  "Automatically push changes and tags after version bump",
			Type:         "bool",
			Required:     false,
			Example:      "true",
		},
		{
			Key:          "app.version.remote_repository_provider",
			DefaultValue: "github",
			Description:  "Remote repository provider (github, bitbucket, gitlab)",
			Type:         "string",
			Required:     false,
			Example:      "gitlab",
		},
	}
}
