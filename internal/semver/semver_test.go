package semver

import (
	"testing"
)

func TestBumpMajor(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		want       string
		shouldFail bool
	}{
		{
			name:       "bump major from 1.2.3",
			version:    "1.2.3",
			want:       "2.0.0",
			shouldFail: false,
		},
		{
			name:       "bump major with v prefix",
			version:    "v1.2.3",
			want:       "2.0.0",
			shouldFail: false,
		},
		{
			name:       "bump major with empty version",
			version:    "",
			want:       "1.0.0",
			shouldFail: false,
		},
		{
			name:       "bump major with invalid version",
			version:    "invalid",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BumpMajor(tt.version)
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
		name       string
		version    string
		want       string
		shouldFail bool
	}{
		{
			name:       "bump minor from 1.2.3",
			version:    "1.2.3",
			want:       "1.3.0",
			shouldFail: false,
		},
		{
			name:       "bump minor with v prefix",
			version:    "v1.2.3",
			want:       "1.3.0",
			shouldFail: false,
		},
		{
			name:       "bump minor with empty version",
			version:    "",
			want:       "0.1.0",
			shouldFail: false,
		},
		{
			name:       "bump minor with invalid version",
			version:    "invalid",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BumpMinor(tt.version)
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
		name       string
		version    string
		want       string
		shouldFail bool
	}{
		{
			name:       "bump patch from 1.2.3",
			version:    "1.2.3",
			want:       "1.2.4",
			shouldFail: false,
		},
		{
			name:       "bump patch with v prefix",
			version:    "v1.2.3",
			want:       "1.2.4",
			shouldFail: false,
		},
		{
			name:       "bump patch with empty version",
			version:    "",
			want:       "0.0.1",
			shouldFail: false,
		},
		{
			name:       "bump patch with invalid version",
			version:    "invalid",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BumpPatch(tt.version)
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
		want       [3]int
		shouldFail bool
	}{
		{
			name:       "valid version",
			version:    "1.2.3",
			want:       [3]int{1, 2, 3},
			shouldFail: false,
		},
		{
			name:       "with v prefix",
			version:    "v1.2.3",
			want:       [3]int{1, 2, 3},
			shouldFail: false,
		},
		{
			name:       "empty version",
			version:    "",
			want:       [3]int{0, 0, 0},
			shouldFail: false,
		},
		{
			name:       "invalid format",
			version:    "1.2",
			shouldFail: true,
		},
		{
			name:       "non-numeric part",
			version:    "1.2.a",
			shouldFail: true,
		},
		{
			name:       "negative number",
			version:    "1.-2.3",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersion(tt.version)
			if (err != nil) != tt.shouldFail {
				t.Errorf("parseVersion() error = %v, shouldFail = %v", err, tt.shouldFail)
				return
			}
			if !tt.shouldFail && got != tt.want {
				t.Errorf("parseVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		name    string
		version [3]int
		want    string
	}{
		{
			name:    "format 1.2.3",
			version: [3]int{1, 2, 3},
			want:    "1.2.3",
		},
		{
			name:    "format 0.0.0",
			version: [3]int{0, 0, 0},
			want:    "0.0.0",
		},
		{
			name:    "format 10.20.30",
			version: [3]int{10, 20, 30},
			want:    "10.20.30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatVersion(tt.version)
			if got != tt.want {
				t.Errorf("formatVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
