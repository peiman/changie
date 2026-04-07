// internal/config/commands/changelog_config.go
//
// Changelog command configuration: metadata + options

package commands

import "github.com/peiman/changie/.ckeletin/pkg/config"

// ChangelogMetadata defines all metadata for the changelog command
var ChangelogMetadata = config.CommandMetadata{
	Use:          "changelog",
	Short:        "Manage changelog entries",
	ConfigPrefix: "app.changelog",
	FlagOverrides: map[string]string{
		"app.changelog.file": "file",
	},
}
