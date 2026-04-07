// internal/config/commands/diff_config.go
//
// Diff command configuration: metadata + options

package commands

import "github.com/peiman/changie/.ckeletin/pkg/config"

// DiffMetadata defines all metadata for the diff command.
var DiffMetadata = config.CommandMetadata{
	Use:          "diff FROM TO",
	Short:        "Show changelog entries between two versions",
	ConfigPrefix: "app.diff",
	FlagOverrides: map[string]string{
		"app.changelog.file": "file",
	},
}
