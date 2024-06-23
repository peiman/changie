package semver

import (
	"testing"
)

func TestBumpMajor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.0.0", "2.0.0"},
		{"0.1.2", "1.0.0"},
		{"1.2.3", "2.0.0"},
	}

	for _, test := range tests {
		result, err := BumpMajor(test.input)
		if err != nil {
			t.Errorf("BumpMajor(%s) returned an error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("BumpMajor(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestBumpMinor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.0.0", "1.1.0"},
		{"0.1.2", "0.2.0"},
		{"1.2.3", "1.3.0"},
	}

	for _, test := range tests {
		result, err := BumpMinor(test.input)
		if err != nil {
			t.Errorf("BumpMinor(%s) returned an error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("BumpMinor(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestBumpPatch(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.0.0", "1.0.1"},
		{"0.1.2", "0.1.3"},
		{"1.2.3", "1.2.4"},
	}

	for _, test := range tests {
		result, err := BumpPatch(test.input)
		if err != nil {
			t.Errorf("BumpPatch(%s) returned an error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("BumpPatch(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.1.0", "1.0.0", 1},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.1.0", -1},
		{"1.0.0", "1.0.1", -1},
	}

	for _, test := range tests {
		result, err := Compare(test.v1, test.v2)
		if err != nil {
			t.Errorf("Compare(%s, %s) returned an error: %v", test.v1, test.v2, err)
		}
		if result != test.expected {
			t.Errorf("Compare(%s, %s) = %d, expected %d", test.v1, test.v2, result, test.expected)
		}
	}
}
