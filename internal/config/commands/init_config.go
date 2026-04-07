// internal/config/commands/init_config.go
//
// Init command configuration: metadata + options

package commands

import "github.com/peiman/changie/.ckeletin/pkg/config"

// InitMetadata defines all metadata for the init command
var InitMetadata = config.CommandMetadata{
	Use:          "init",
	Short:        "Initialize a project with SemVer and Keep a Changelog",
	ConfigPrefix: "app.init",
	FlagOverrides: map[string]string{
		"app.changelog.file":       "file",
		"app.version.use_v_prefix": "use-v-prefix",
	},
}
