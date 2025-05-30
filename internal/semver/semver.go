// Package semver provides functionality for working with semantic versioning.
//
// This package implements the Semantic Versioning 2.0.0 specification (https://semver.org/)
// for working with version numbers in the format MAJOR.MINOR.PATCH. It provides:
// - Functions for parsing and validating SemVer strings
// - Methods for incrementing major, minor, and patch version numbers
// - Utilities for comparing version numbers
//
// The package handles version strings with or without the "v" prefix and provides
// clear error messages for invalid version formats.
package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// BumpMajor increases the major version number and resets minor and patch to 0.
//
// According to SemVer, the major version should be incremented when making
// incompatible API changes. This function also resets the minor and patch
// versions to 0, as per SemVer convention.
//
// Parameters:
//   - version: A semantic version string (e.g., "1.2.3" or "v1.2.3")
//
// Returns:
//   - string: The new version with incremented major version (e.g., "2.0.0")
//   - error: Any error encountered during parsing or formatting
func BumpMajor(version string) (string, error) {
	v, err := parseVersion(version)
	if err != nil {
		return "", err
	}
	v[0]++
	v[1] = 0
	v[2] = 0
	return formatVersion(v), nil
}

// BumpMinor increases the minor version number and resets patch to 0.
//
// According to SemVer, the minor version should be incremented when adding
// functionality in a backward-compatible manner. This function also resets
// the patch version to 0, as per SemVer convention.
//
// Parameters:
//   - version: A semantic version string (e.g., "1.2.3" or "v1.2.3")
//
// Returns:
//   - string: The new version with incremented minor version (e.g., "1.3.0")
//   - error: Any error encountered during parsing or formatting
func BumpMinor(version string) (string, error) {
	v, err := parseVersion(version)
	if err != nil {
		return "", err
	}
	v[1]++
	v[2] = 0
	return formatVersion(v), nil
}

// BumpPatch increases the patch version number.
//
// According to SemVer, the patch version should be incremented when making
// backward-compatible bug fixes.
//
// Parameters:
//   - version: A semantic version string (e.g., "1.2.3" or "v1.2.3")
//
// Returns:
//   - string: The new version with incremented patch version (e.g., "1.2.4")
//   - error: Any error encountered during parsing or formatting
func BumpPatch(version string) (string, error) {
	v, err := parseVersion(version)
	if err != nil {
		return "", err
	}
	v[2]++
	return formatVersion(v), nil
}

// Compare compares two version strings according to semantic versioning rules.
//
// This function implements version precedence as defined in the SemVer specification.
// It compares each component of the version (major, minor, patch) in order of significance.
//
// Parameters:
//   - v1: First version string to compare
//   - v2: Second version string to compare
//
// Returns:
//   - int: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
//   - error: Any error encountered during parsing
func Compare(v1, v2 string) (int, error) {
	ver1, err := parseVersion(v1)
	if err != nil {
		return 0, err
	}
	ver2, err := parseVersion(v2)
	if err != nil {
		return 0, err
	}

	for i := 0; i < 3; i++ {
		if ver1[i] > ver2[i] {
			return 1, nil
		}
		if ver1[i] < ver2[i] {
			return -1, nil
		}
	}
	return 0, nil
}

// parseVersion converts a version string to an array of integers.
//
// This function parses a version string in the format "X.Y.Z" or "vX.Y.Z" into
// its numeric components. It handles the "v" prefix gracefully if present.
//
// Parameters:
//   - version: A version string to parse
//
// Returns:
//   - [3]int: Array containing the major, minor, and patch version numbers
//   - error: Error if the version string is not in a valid format
func parseVersion(version string) ([3]int, error) {
	// Handle empty version string
	if version == "" {
		return [3]int{0, 0, 0}, nil
	}

	// Trim 'v' prefix if present
	trimmedVersion := strings.TrimPrefix(version, "v")

	parts := strings.Split(trimmedVersion, ".")
	if len(parts) != 3 {
		return [3]int{}, fmt.Errorf("invalid version format: %s (expected format: X.Y.Z)", version)
	}

	var v [3]int
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return [3]int{}, fmt.Errorf("invalid version number: %s (must be a valid integer)", part)
		}
		if num < 0 {
			return [3]int{}, fmt.Errorf("invalid version number: %s (must be non-negative)", part)
		}
		v[i] = num
	}
	return v, nil
}

// formatVersion converts an array of integers to a version string.
//
// This function creates a properly formatted semantic version string from
// the provided version components, without the "v" prefix.
//
// Parameters:
//   - v: Array containing the major, minor, and patch version numbers
//
// Returns:
//   - string: Formatted version string in the format "X.Y.Z"
func formatVersion(v [3]int) string {
	return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
}
