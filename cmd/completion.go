// cmd/completion.go
// ckeletin:allow-custom-command

package cmd

import (
	"github.com/spf13/cobra"
)

// completionCmd generates shell completion scripts.
// Note: Long is set in root.go's init() after binaryName is resolved,
// because Go evaluates var declarations before init() runs.
var completionCmd = &cobra.Command{
	Use:                   "completion",
	Short:                 "Generate the autocompletion script for the specified shell",
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default to bash if no args provided:
		return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
