// Package changelog - diff_test.go provides tests for DiffVersions().

package changelog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sampleChangelog is a multi-version changelog for testing.
const sampleChangelog = `# Changelog

## [Unreleased]

### Added

- Upcoming feature

## [1.2.0] - 2024-03-01

### Added

- Feature C

### Fixed

- Bug fix C

## [1.1.0] - 2024-02-01

### Added

- Feature B

## [1.0.0] - 2024-01-01

### Fixed

- Bug fix A

[Unreleased]: https://github.com/user/repo/compare/1.2.0...HEAD
[1.2.0]: https://github.com/user/repo/compare/1.1.0...1.2.0
[1.1.0]: https://github.com/user/repo/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/user/repo/releases/tag/1.0.0
`

// changelogWithVPrefix uses v-prefixed version headers.
const changelogWithVPrefix = `# Changelog

## [v1.1.0] - 2024-02-01

### Added

- Feature B

## [v1.0.0] - 2024-01-01

### Fixed

- Bug fix A

[v1.1.0]: https://github.com/user/repo/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/user/repo/releases/tag/v1.0.0
`

func TestDiffVersions_HappyPath_TwoAdjacentVersions(t *testing.T) {
	result, err := DiffVersions(sampleChangelog, "1.1.0", "1.2.0")

	require.NoError(t, err)
	assert.Contains(t, result, "## [1.2.0]")
	assert.Contains(t, result, "Feature C")
	assert.Contains(t, result, "Bug fix C")
	// The 1.1.0 section should NOT be included
	assert.NotContains(t, result, "## [1.1.0]")
	assert.NotContains(t, result, "Feature B")
}

func TestDiffVersions_MultipleVersionsInRange(t *testing.T) {
	result, err := DiffVersions(sampleChangelog, "1.0.0", "1.2.0")

	require.NoError(t, err)
	assert.Contains(t, result, "## [1.2.0]")
	assert.Contains(t, result, "Feature C")
	assert.Contains(t, result, "## [1.1.0]")
	assert.Contains(t, result, "Feature B")
	// The 1.0.0 section should NOT be included (exclusive lower bound)
	assert.NotContains(t, result, "## [1.0.0]")
	assert.NotContains(t, result, "Bug fix A")
}

func TestDiffVersions_WithVPrefix(t *testing.T) {
	result, err := DiffVersions(sampleChangelog, "v1.1.0", "v1.2.0")

	require.NoError(t, err)
	assert.Contains(t, result, "## [1.2.0]")
	assert.Contains(t, result, "Feature C")
}

func TestDiffVersions_WithoutVPrefix(t *testing.T) {
	result, err := DiffVersions(sampleChangelog, "1.1.0", "1.2.0")

	require.NoError(t, err)
	assert.Contains(t, result, "## [1.2.0]")
	assert.Contains(t, result, "Feature C")
}

func TestDiffVersions_MixedPrefix(t *testing.T) {
	// User passes v1.1.0 but changelog has [1.1.0] — should normalize
	result, err := DiffVersions(sampleChangelog, "v1.1.0", "1.2.0")

	require.NoError(t, err)
	assert.Contains(t, result, "## [1.2.0]")
	assert.Contains(t, result, "Feature C")
}

func TestDiffVersions_ChangelogWithVPrefixedHeaders(t *testing.T) {
	// Changelog uses [v1.0.0] format, user passes bare version
	result, err := DiffVersions(changelogWithVPrefix, "1.0.0", "1.1.0")

	require.NoError(t, err)
	assert.Contains(t, result, "## [v1.1.0]")
	assert.Contains(t, result, "Feature B")
}

func TestDiffVersions_FromVersionNotFound(t *testing.T) {
	_, err := DiffVersions(sampleChangelog, "0.5.0", "1.0.0")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "0.5.0")
	assert.Contains(t, strings.ToLower(err.Error()), "not found")
}

func TestDiffVersions_ToVersionNotFound(t *testing.T) {
	_, err := DiffVersions(sampleChangelog, "1.0.0", "3.0.0")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "3.0.0")
	assert.Contains(t, strings.ToLower(err.Error()), "not found")
}

func TestDiffVersions_InvertedVersions(t *testing.T) {
	_, err := DiffVersions(sampleChangelog, "1.2.0", "1.0.0")

	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "older")
}

func TestDiffVersions_SameVersion(t *testing.T) {
	_, err := DiffVersions(sampleChangelog, "1.0.0", "1.0.0")

	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "different")
}

func TestDiffVersions_InvalidSemverFrom(t *testing.T) {
	_, err := DiffVersions(sampleChangelog, "not-a-version", "1.0.0")

	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "invalid")
}

func TestDiffVersions_InvalidSemverTo(t *testing.T) {
	_, err := DiffVersions(sampleChangelog, "1.0.0", "abc")

	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "invalid")
}

func TestDiffVersions_EmptyContent(t *testing.T) {
	_, err := DiffVersions("", "1.0.0", "2.0.0")

	require.Error(t, err)
}

func TestDiffVersions_OnlyUnreleased(t *testing.T) {
	content := `# Changelog

## [Unreleased]

### Added

- Some feature
`
	_, err := DiffVersions(content, "1.0.0", "2.0.0")

	require.Error(t, err)
}

func TestDiffVersions_MultipleSectionsPerVersion(t *testing.T) {
	// v1.2.0 has Added + Fixed + Changed — all should be in output
	result, err := DiffVersions(sampleChangelog, "1.1.0", "1.2.0")

	require.NoError(t, err)
	assert.Contains(t, result, "### Added")
	assert.Contains(t, result, "Feature C")
	assert.Contains(t, result, "### Fixed")
	assert.Contains(t, result, "Bug fix C")
}

func TestDiffVersions_UnreleasedNotIncluded(t *testing.T) {
	result, err := DiffVersions(sampleChangelog, "1.1.0", "1.2.0")

	require.NoError(t, err)
	assert.NotContains(t, result, "Unreleased")
	assert.NotContains(t, result, "Upcoming feature")
}

func TestDiffVersions_WindowsLineEndings(t *testing.T) {
	content := strings.ReplaceAll(sampleChangelog, "\n", "\r\n")
	result, err := DiffVersions(content, "1.0.0", "1.2.0")

	require.NoError(t, err)
	assert.Contains(t, result, "Feature C")
	assert.Contains(t, result, "Feature B")
}
