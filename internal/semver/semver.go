// Package semver provides functionality for working with semantic versioning.
package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// BumpMajor increases the major version number and resets minor and patch to 0.
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
func BumpPatch(version string) (string, error) {
	v, err := parseVersion(version)
	if err != nil {
		return "", err
	}
	v[2]++
	return formatVersion(v), nil
}

// Compare compares two version strings.
// It returns -1 if v1 < v2, 0 if v1 == v2, and 1 if v1 > v2.
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
func formatVersion(v [3]int) string {
	return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
}
