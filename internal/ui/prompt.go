// Package ui provides user interface utilities.
package ui

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// AskYesNo prompts the user with a yes/no question and returns their response.
//
// This function displays a prompt to the user and reads their input from stdin.
// It accepts various forms of yes/no responses:
// - Yes: "y", "yes", "" (empty for default)
// - No: "n", "no"
//
// The function is case-insensitive.
//
// Parameters:
//   - prompt: The question to display to the user
//   - defaultYes: If true, empty input is treated as "yes", otherwise as "no"
//   - output: Writer for displaying the prompt (typically os.Stdout or cmd.OutOrStdout())
//
// Returns:
//   - bool: true for yes, false for no
//   - error: Any error encountered while reading input
//
// Example:
//
//	useVPrefix, err := AskYesNo("Use 'v' prefix for version tags?", true, os.Stdout)
//	if err != nil {
//	    return err
//	}
func AskYesNo(prompt string, defaultYes bool, output io.Writer) (bool, error) {
	// Construct the prompt with default indicator
	defaultIndicator := "[Y/n]"
	if !defaultYes {
		defaultIndicator = "[y/N]"
	}

	_, err := output.Write([]byte(prompt + " " + defaultIndicator + ": "))
	if err != nil {
		log.Warn().Err(err).Msg("Error writing prompt")
		return defaultYes, nil // Return nil error since we've logged it and have a default value
	}

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Warn().Err(err).Msgf("Error reading input, defaulting to %v", defaultYes)
		return defaultYes, nil // Return nil error since we've logged it and have a default value
	}

	// Normalize input
	input = strings.TrimSpace(strings.ToLower(input))

	// Handle empty input (use default)
	if input == "" {
		return defaultYes, nil
	}

	// Check for yes/no responses
	switch input {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		// Unrecognized input, use default
		log.Warn().Str("input", input).Msgf("Unrecognized input, defaulting to %v", defaultYes)
		return defaultYes, nil
	}
}
