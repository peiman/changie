package semver

import (
	"errors"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
)

func TestBumpMajor(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		preservePrefix bool
		want           string
		shouldFail     bool
	}{
		{
			name:           "bump major from 1.2.3",
			version:        "1.2.3",
			preservePrefix: false,
			want:           "2.0.0",
			shouldFail:     false,
		},
		{
			name:           "bump major with v prefix",
			version:        "v1.2.3",
			preservePrefix: false,
			want:           "2.0.0",
			shouldFail:     false,
		},
		{
			name:           "bump major with empty version",
			version:        "",
			preservePrefix: false,
			want:           "1.0.0",
			shouldFail:     false,
		},
		{
			name:           "bump major with invalid version",
			version:        "invalid",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump major with very large numbers",
			version:        "999999999.999999999.999999999",
			preservePrefix: false,
			want:           "1000000000.0.0",
			shouldFail:     false,
		},
		{
			name:           "bump major with leading zeroes (invalid)",
			version:        "01.02.03",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump major with wrong format (too few segments)",
			version:        "1.2",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump major with wrong format (too many segments)",
			version:        "1.2.3.4",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump major with v prefix and preserve prefix",
			version:        "v1.2.3",
			preservePrefix: true,
			want:           "v2.0.0",
			shouldFail:     false,
		},
		{
			name:           "bump major without v prefix and preserve prefix",
			version:        "1.2.3",
			preservePrefix: true,
			want:           "2.0.0",
			shouldFail:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BumpMajor(tt.version, tt.preservePrefix)

			if (err != nil) != tt.shouldFail {
				t.Errorf("BumpMajor() error = %v, shouldFail = %v", err, tt.shouldFail)
				return
			}
			if !tt.shouldFail && got != tt.want {
				t.Errorf("BumpMajor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBumpMinor(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		preservePrefix bool
		want           string
		shouldFail     bool
	}{
		{
			name:           "bump minor from 1.2.3",
			version:        "1.2.3",
			preservePrefix: false,
			want:           "1.3.0",
			shouldFail:     false,
		},
		{
			name:           "bump minor with v prefix",
			version:        "v1.2.3",
			preservePrefix: false,
			want:           "1.3.0",
			shouldFail:     false,
		},
		{
			name:           "bump minor with empty version",
			version:        "",
			preservePrefix: false,
			want:           "0.1.0",
			shouldFail:     false,
		},
		{
			name:           "bump minor with invalid version",
			version:        "invalid",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump minor with version at zero",
			version:        "1.0.0",
			preservePrefix: false,
			want:           "1.1.0",
			shouldFail:     false,
		},
		{
			name:           "bump minor with negative numbers (invalid)",
			version:        "1.-2.3",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump minor with non-numeric segments",
			version:        "1.a.3",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump minor with leading zeros (invalid)",
			version:        "1.01.3",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump minor with v prefix and preserve prefix",
			version:        "v1.2.3",
			preservePrefix: true,
			want:           "v1.3.0",
			shouldFail:     false,
		},
		{
			name:           "bump minor without v prefix and preserve prefix",
			version:        "1.2.3",
			preservePrefix: true,
			want:           "1.3.0",
			shouldFail:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BumpMinor(tt.version, tt.preservePrefix)

			if (err != nil) != tt.shouldFail {
				t.Errorf("BumpMinor() error = %v, shouldFail = %v", err, tt.shouldFail)
				return
			}
			if !tt.shouldFail && got != tt.want {
				t.Errorf("BumpMinor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBumpPatch(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		preservePrefix bool
		want           string
		shouldFail     bool
	}{
		{
			name:           "bump patch from 1.2.3",
			version:        "1.2.3",
			preservePrefix: false,
			want:           "1.2.4",
			shouldFail:     false,
		},
		{
			name:           "bump patch with v prefix",
			version:        "v1.2.3",
			preservePrefix: false,
			want:           "1.2.4",
			shouldFail:     false,
		},
		{
			name:           "bump patch with empty version",
			version:        "",
			preservePrefix: false,
			want:           "0.0.1",
			shouldFail:     false,
		},
		{
			name:           "bump patch with invalid version",
			version:        "invalid",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump patch on all zeros",
			version:        "0.0.0",
			preservePrefix: false,
			want:           "0.0.1",
			shouldFail:     false,
		},
		{
			name:           "bump patch with decimal numbers (invalid)",
			version:        "1.2.3.4",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump patch with special characters (invalid)",
			version:        "1.2.3-beta",
			preservePrefix: false,
			shouldFail:     false, // Changed to false, as we now handle prerelease identifiers
			want:           "1.2.4-beta",
		},
		{
			name:           "bump patch with leading zeros (invalid)",
			version:        "1.2.03",
			preservePrefix: false,
			shouldFail:     true,
		},
		{
			name:           "bump patch with v prefix and preserve prefix",
			version:        "v1.2.3",
			preservePrefix: true,
			want:           "v1.2.4",
			shouldFail:     false,
		},
		{
			name:           "bump patch without v prefix and preserve prefix",
			version:        "1.2.3",
			preservePrefix: true,
			want:           "1.2.4",
			shouldFail:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BumpPatch(tt.version, tt.preservePrefix)

			if (err != nil) != tt.shouldFail {
				t.Errorf("BumpPatch() error = %v, shouldFail = %v", err, tt.shouldFail)
				return
			}
			if !tt.shouldFail && got != tt.want {
				t.Errorf("BumpPatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name       string
		v1         string
		v2         string
		want       int
		shouldFail bool
	}{
		{
			name:       "equal versions",
			v1:         "1.2.3",
			v2:         "1.2.3",
			want:       0,
			shouldFail: false,
		},
		{
			name:       "v1 > v2 major",
			v1:         "2.0.0",
			v2:         "1.9.9",
			want:       1,
			shouldFail: false,
		},
		{
			name:       "v1 < v2 major",
			v1:         "1.0.0",
			v2:         "2.0.0",
			want:       -1,
			shouldFail: false,
		},
		{
			name:       "v1 > v2 minor",
			v1:         "1.2.0",
			v2:         "1.1.9",
			want:       1,
			shouldFail: false,
		},
		{
			name:       "v1 < v2 minor",
			v1:         "1.1.0",
			v2:         "1.2.0",
			want:       -1,
			shouldFail: false,
		},
		{
			name:       "v1 > v2 patch",
			v1:         "1.1.2",
			v2:         "1.1.1",
			want:       1,
			shouldFail: false,
		},
		{
			name:       "v1 < v2 patch",
			v1:         "1.1.1",
			v2:         "1.1.2",
			want:       -1,
			shouldFail: false,
		},
		{
			name:       "v1 invalid",
			v1:         "invalid",
			v2:         "1.0.0",
			shouldFail: true,
		},
		{
			name:       "v2 invalid",
			v1:         "1.0.0",
			v2:         "invalid",
			shouldFail: true,
		},
		{
			name:       "compare with different v prefix",
			v1:         "v1.0.0",
			v2:         "1.0.0",
			want:       0,
			shouldFail: false,
		},
		{
			name:       "compare with leading zeros (invalid)",
			v1:         "1.01.0",
			v2:         "1.1.0",
			shouldFail: true,
		},
		{
			name:       "compare with very large numbers",
			v1:         "999999999.0.0",
			v2:         "0.999999999.0",
			want:       1,
			shouldFail: false,
		},
		{
			name:       "empty string comparison",
			v1:         "",
			v2:         "",
			want:       0,
			shouldFail: false,
		},
		{
			name:       "compare empty with zero version",
			v1:         "",
			v2:         "0.0.0",
			want:       0,
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compare(tt.v1, tt.v2)
			if (err != nil) != tt.shouldFail {
				t.Errorf("Compare() error = %v, shouldFail = %v", err, tt.shouldFail)
				return
			}
			if !tt.shouldFail && got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		wantVer    semver.Version
		wantPrefix bool
		wantErr    error
	}{
		{
			name:       "standard version",
			version:    "1.2.3",
			wantVer:    semver.MustParse("1.2.3"),
			wantPrefix: false,
		},
		{
			name:       "v prefix",
			version:    "v1.2.3",
			wantVer:    semver.MustParse("1.2.3"),
			wantPrefix: true,
		},
		{
			name:       "version with prerelease",
			version:    "1.2.3-alpha.1",
			wantVer:    semver.MustParse("1.2.3-alpha.1"),
			wantPrefix: false,
		},
		{
			name:       "version with build metadata",
			version:    "1.2.3+build.123",
			wantVer:    semver.MustParse("1.2.3+build.123"),
			wantPrefix: false,
		},
		{
			name:       "version with prerelease and build metadata",
			version:    "1.2.3-alpha.1+build.123",
			wantVer:    semver.MustParse("1.2.3-alpha.1+build.123"),
			wantPrefix: false,
		},
		{
			name:       "zero version",
			version:    "0.0.0",
			wantVer:    semver.MustParse("0.0.0"),
			wantPrefix: false,
		},
		{
			name:    "invalid format",
			version: "not.a.version",
			wantErr: ErrInvalidVersion,
		},
		{
			name:    "invalid format with too many parts",
			version: "1.2.3.4",
			wantErr: ErrInvalidVersion,
		},
		{
			name:    "invalid format with too few parts",
			version: "1.2",
			wantErr: ErrInvalidVersion,
		},
		{
			name:    "leading zero in major",
			version: "01.2.3",
			wantErr: ErrInvalidVersion,
		},
		{
			name:    "leading zero in minor",
			version: "1.02.3",
			wantErr: ErrInvalidVersion,
		},
		{
			name:    "leading zero in patch",
			version: "1.2.03",
			wantErr: ErrInvalidVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ver, hasPrefix, err := ParseVersion(tt.version)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error %v, got %v", tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantVer, ver)
			assert.Equal(t, tt.wantPrefix, hasPrefix)
		})
	}
}

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		name          string
		version       semver.Version
		includePrefix bool
		want          string
	}{
		{
			name:          "no prefix",
			version:       semver.MustParse("1.2.3"),
			includePrefix: false,
			want:          "1.2.3",
		},
		{
			name:          "with prefix",
			version:       semver.MustParse("1.2.3"),
			includePrefix: true,
			want:          "v1.2.3",
		},
		{
			name:          "with prerelease no prefix",
			version:       semver.MustParse("1.2.3-alpha.1"),
			includePrefix: false,
			want:          "1.2.3-alpha.1",
		},
		{
			name:          "with prerelease and prefix",
			version:       semver.MustParse("1.2.3-alpha.1"),
			includePrefix: true,
			want:          "v1.2.3-alpha.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatVersion(tt.version, tt.includePrefix)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBumpVersion(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		bumpType       BumpType
		preservePrefix bool
		want           string
		wantErr        error
	}{
		{
			name:           "bump major",
			version:        "1.2.3",
			bumpType:       Major,
			preservePrefix: false,
			want:           "2.0.0",
		},
		{
			name:           "bump major with v prefix",
			version:        "v1.2.3",
			bumpType:       Major,
			preservePrefix: true,
			want:           "v2.0.0",
		},
		{
			name:           "bump minor",
			version:        "1.2.3",
			bumpType:       Minor,
			preservePrefix: false,
			want:           "1.3.0",
		},
		{
			name:           "bump patch",
			version:        "1.2.3",
			bumpType:       Patch,
			preservePrefix: false,
			want:           "1.2.4",
		},
		{
			name:           "bump major with prerelease",
			version:        "1.2.3-alpha.1",
			bumpType:       Major,
			preservePrefix: false,
			want:           "2.0.0-alpha.1",
		},
		{
			name:           "bump major with build metadata",
			version:        "1.2.3+build.123",
			bumpType:       Major,
			preservePrefix: false,
			want:           "2.0.0+build.123",
		},
		{
			name:           "bump major with prerelease and build metadata",
			version:        "1.2.3-alpha.1+build.123",
			bumpType:       Major,
			preservePrefix: false,
			want:           "2.0.0-alpha.1+build.123",
		},
		{
			name:           "invalid version",
			version:        "not.a.version",
			bumpType:       Major,
			preservePrefix: false,
			wantErr:        ErrInvalidVersion,
		},
		{
			name:           "invalid bump type",
			version:        "1.2.3",
			bumpType:       "invalid",
			preservePrefix: false,
			wantErr:        ErrInvalidBump,
		},
		{
			name:           "bump major with v prefix and preserve prefix false",
			version:        "v1.2.3",
			bumpType:       Major,
			preservePrefix: false,
			want:           "2.0.0",
		},
		{
			name:           "bump minor with v prefix and preserve prefix true",
			version:        "v1.2.3",
			bumpType:       Minor,
			preservePrefix: true,
			want:           "v1.3.0",
		},
		{
			name:           "bump patch with v prefix and preserve prefix true",
			version:        "v1.2.3",
			bumpType:       Patch,
			preservePrefix: true,
			want:           "v1.2.4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BumpVersion(tt.version, tt.bumpType, tt.preservePrefix)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error %v, got %v", tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
