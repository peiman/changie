// Package semver provides functions for parsing, validating, and bumping
// semantic version numbers according to the Semantic Versioning 2.0.0 specification.
// See https://semver.org/ for the full specification.
//
// This package handles versions with or without a leading 'v' prefix
// and provides options to preserve or remove the prefix when formatting.
package semver

import (
	"errors"
	"strings"

	"github.com/blang/semver/v4"
)

// Common errors returned by semver operations
var (
	ErrInvalidVersion = errors.New("invalid version format")
	ErrInvalidBump    = errors.New("invalid bump type")
)

// BumpType represents the type of version bump to perform
type BumpType string

// Constants for the different types of version bumps
const (
	Major BumpType = "major"
	Minor BumpType = "minor"
	Patch BumpType = "patch"
)

// ParseVersion parses a version string into its components
// Returns the parsed semver.Version and whether it had a v prefix
func ParseVersion(version string) (semver.Version, bool, error) {
	if version == "" {
		// Default to 0.0.0 if empty
		return semver.MustParse("0.0.0"), false, nil
	}

	// Check for v prefix
	hasPrefix := strings.HasPrefix(version, "v")
	versionWithoutPrefix := version
	if hasPrefix {
		versionWithoutPrefix = version[1:]
	}

	// Parse using the blang/semver library
	parsedVersion, err := semver.Parse(versionWithoutPrefix)
	if err != nil {
		return semver.Version{}, false, ErrInvalidVersion
	}

	return parsedVersion, hasPrefix, nil
}

// FormatVersion formats a semver.Version object as a string
// If includePrefix is true, the string will be prefixed with 'v'
func FormatVersion(ver semver.Version, includePrefix bool) string {
	if includePrefix {
		return "v" + ver.String()
	}
	return ver.String()
}

// BumpVersion increments a version number based on the bump type
// Returns the new version string and an error if the version is invalid
func BumpVersion(version string, bumpType BumpType, preservePrefix bool) (string, error) {
	parsedVersion, hasPrefix, err := ParseVersion(version)
	if err != nil {
		return "", err
	}

	// Apply the bump
	switch bumpType {
	case Major:
		parsedVersion.Major++
		parsedVersion.Minor = 0
		parsedVersion.Patch = 0
	case Minor:
		parsedVersion.Minor++
		parsedVersion.Patch = 0
	case Patch:
		parsedVersion.Patch++
	default:
		return "", ErrInvalidBump
	}

	// Format the result based on preservePrefix and original prefix
	return FormatVersion(parsedVersion, preservePrefix && hasPrefix), nil
}

// BumpMajor increments the major version number
// This resets the minor and patch numbers to 0
func BumpMajor(version string, preservePrefix bool) (string, error) {
	return BumpVersion(version, Major, preservePrefix)
}

// BumpMinor increments the minor version number
// This resets the patch number to 0
func BumpMinor(version string, preservePrefix bool) (string, error) {
	return BumpVersion(version, Minor, preservePrefix)
}

// BumpPatch increments the patch version number
func BumpPatch(version string, preservePrefix bool) (string, error) {
	return BumpVersion(version, Patch, preservePrefix)
}

// Compare compares two version strings and returns:
//
//	-1 if v1 < v2
//	 0 if v1 == v2
//	+1 if v1 > v2
func Compare(v1, v2 string) (int, error) {
	// Parse both versions
	parsedV1, _, err := ParseVersion(v1)
	if err != nil {
		return 0, err
	}

	parsedV2, _, err := ParseVersion(v2)
	if err != nil {
		return 0, err
	}

	// Compare using the semver library
	return parsedV1.Compare(parsedV2), nil
}
