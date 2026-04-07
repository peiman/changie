// internal/config/commands/bump_config.go
//
// Bump command configuration: metadata + options

package commands

import "github.com/peiman/changie/.ckeletin/pkg/config"

// BumpMetadata defines all metadata for the bump command
var BumpMetadata = config.CommandMetadata{
	Use:          "bump",
	Short:        "Bump version numbers following semantic versioning",
	ConfigPrefix: "app.bump",
	FlagOverrides: map[string]string{
		"app.changelog.file":                "file",
		"app.changelog.repository_provider": "rrp",
		"app.changelog.auto_push":           "auto-push",
		"app.version.allow_any_branch":      "allow-any-branch",
		"app.version.use_v_prefix":          "use-v-prefix",
	},
}
