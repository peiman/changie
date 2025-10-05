// Package config provides configuration management utilities.
package config

import (
	"regexp"
	"strings"
)

// EnvPrefix returns a sanitized environment variable prefix based on the binary name.
//
// This function converts the binary name to a valid environment variable prefix by:
// 1. Converting to uppercase
// 2. Replacing non-alphanumeric characters with underscores
// 3. Ensuring it doesn't start with a number
// 4. Handling edge cases where all characters were special
//
// Parameters:
//   - binaryName: The name of the binary/application
//
// Returns:
//   - string: A sanitized environment variable prefix (e.g., "CHANGIE" from "changie")
//
// Example:
//   - "changie" -> "CHANGIE"
//   - "my-app" -> "MY_APP"
//   - "123app" -> "_123APP"
func EnvPrefix(binaryName string) string {
	// Convert to uppercase and replace non-alphanumeric characters with underscore
	prefix := strings.ToUpper(binaryName)
	re := regexp.MustCompile(`[^A-Z0-9]`)
	prefix = re.ReplaceAllString(prefix, "_")

	// Ensure it doesn't start with a number (invalid for env vars)
	if prefix != "" && prefix[0] >= '0' && prefix[0] <= '9' {
		prefix = "_" + prefix
	}

	// Handle case where all characters were special and got replaced
	re = regexp.MustCompile(`^_+$`)
	if re.MatchString(prefix) {
		prefix = "_"
	}

	return prefix
}
