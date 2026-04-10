// Package changelog - diff.go provides DiffVersions for extracting changelog
// content between two semantic versions.

package changelog

import (
	"fmt"
	"strings"

	"github.com/peiman/changie/internal/semver"
)

// DiffVersions extracts changelog content between two versions.
//
// It returns all sections and entries for versions strictly greater than
// fromVersion up to and including toVersion. Both versions must exist
// in the changelog content.
//
// Version strings may include a 'v' prefix (e.g., "v1.0.0" or "1.0.0").
//
// Parameters:
//   - content: The full changelog file content as a string
//   - fromVersion: The older version (exclusive — its content is NOT included)
//   - toVersion: The newer version (inclusive — its content IS included)
//
// Returns:
//   - string: The extracted changelog content between the two versions
//   - error: If versions are not found, invalid, or fromVersion >= toVersion
func DiffVersions(content, fromVersion, toVersion string) (string, error) {
	// Validate and parse both versions
	fromParsed, _, err := semver.ParseVersion(fromVersion)
	if err != nil {
		return "", fmt.Errorf("invalid version %q: %w", fromVersion, err)
	}

	toParsed, _, err := semver.ParseVersion(toVersion)
	if err != nil {
		return "", fmt.Errorf("invalid version %q: %w", toVersion, err)
	}

	// Enforce ordering: fromVersion must be strictly older than toVersion
	cmp := fromParsed.Compare(toParsed)
	if cmp == 0 {
		return "", fmt.Errorf("versions must be different: both are %q", fromVersion)
	}
	if cmp > 0 {
		return "", fmt.Errorf("from-version must be older than to-version: %q >= %q", fromVersion, toVersion)
	}

	// Normalize line endings and split into lines
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")

	// Strip v-prefix from user-supplied versions for comparison against headers
	fromBare := strings.TrimPrefix(fromVersion, "v")
	toBare := strings.TrimPrefix(toVersion, "v")

	// Find line indices for the two version headers.
	// Changelog headers may appear as ## [1.0.0] or ## [v1.0.0]; normalize both.
	fromIdx := -1
	toIdx := -1

	for i, line := range lines {
		m := reVersionHeader.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		headerVer := strings.TrimPrefix(m[1], "v") // normalize header's v-prefix
		if headerVer == fromBare {
			fromIdx = i
		}
		if headerVer == toBare {
			toIdx = i
		}
	}

	if toIdx == -1 {
		return "", fmt.Errorf("version %q not found in changelog", toBare)
	}
	if fromIdx == -1 {
		return "", fmt.Errorf("version %q not found in changelog", fromBare)
	}

	// Extract lines from toVersion header (inclusive) up to fromVersion header (exclusive)
	extracted := lines[toIdx:fromIdx]

	// Trim trailing blank lines
	result := strings.TrimRight(strings.Join(extracted, "\n"), "\n ")

	return result, nil
}
