// cmd/completion_test.go

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCompletionCmd(t *testing.T) {
	// Save original root command
	origRoot := RootCmd
	defer func() { RootCmd = origRoot }()

	// Create a fresh root command for testing
	RootCmd = &cobra.Command{Use: binaryName}
	RootCmd.AddCommand(completionCmd)

	// Test bash completion (default)
	buf := new(bytes.Buffer)
	completionCmd.SetOut(buf)
	completionCmd.SetErr(buf)

	err := completionCmd.RunE(completionCmd, []string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := buf.String()
	// Bash completion should contain bash-specific content
	if !strings.Contains(output, "bash completion") && !strings.Contains(output, "_"+binaryName) {
		t.Errorf("Expected bash completion script, got: %s", output)
	}
}

func TestCompletionCmd_Exists(t *testing.T) {
	// Verify completion command is registered
	found := false
	for _, cmd := range RootCmd.Commands() {
		if cmd.Name() == "completion" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected completion command to be registered")
	}
}
