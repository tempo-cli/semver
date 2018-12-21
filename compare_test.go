package semver

import (
	"testing"
)

func TestCompareGreaterThan(t *testing.T) {
	cases := []struct {
		versionA string
		versionB string
		result   bool
	}{
		{"1.25.0", "1.24.0", true},
		{"1.25.0", "1.25.0", false},
		{"1.25.0", "1.26.0", false},
	}

	for _, tc := range cases {
		version1, _ := NewVersion(tc.versionA)
		version2, _ := NewVersion(tc.versionB)

		if tc.result != version1.GreaterThan(version2) {
			t.Fatalf("error for comparison: %s is not greater than %s", tc.versionA, tc.versionB)
		}
	}
}

func TestCompareGreaterThanOrEqual(t *testing.T) {
	cases := []struct {
		versionA string
		versionB string
		result   bool
	}{
		{"1.25.0", "1.24.0", true},
		{"1.25.0", "1.25.0", true},
		{"1.25.0", "1.26.0", false},
	}

	for _, tc := range cases {
		version1, _ := NewVersion(tc.versionA)
		version2, _ := NewVersion(tc.versionB)

		if tc.result != version1.GreaterThanOrEqual(version2) {
			t.Fatalf("error for comparison: %s is not greater than or equal to %s", tc.versionA, tc.versionB)
		}
	}
}

func TestCompareLessThan(t *testing.T) {
	cases := []struct {
		versionA string
		versionB string
		result   bool
	}{
		{"1.25.0", "1.24.0", false},
		{"1.25.0", "1.25.0", false},
		{"1.25.0", "1.26.0", true},
		{"1.0.0", "1.2-dev", true},
	}

	for _, tc := range cases {
		version1, _ := NewVersion(tc.versionA)
		version2, _ := NewVersion(tc.versionB)

		if tc.result != version1.LessThan(version2) {
			t.Fatalf("error for comparison: %s is not less than %s", tc.versionA, tc.versionB)
		}
	}
}

func TestCompareLessThanOrEqual(t *testing.T) {
	cases := []struct {
		versionA string
		versionB string
		result   bool
	}{
		{"1.25.0", "1.24.0", false},
		{"1.25.0", "1.25.0", true},
		{"1.25.0", "1.26.0", true},
	}

	for _, tc := range cases {
		version1, _ := NewVersion(tc.versionA)
		version2, _ := NewVersion(tc.versionB)

		if tc.result != version1.LessThanOrEqual(version2) {
			t.Fatalf("error for comparison: %s is not less than or equal to %s", tc.versionA, tc.versionB)
		}
	}
}

func TestCompareEqual(t *testing.T) {
	cases := []struct {
		versionA string
		versionB string
		result   bool
	}{
		{"1.25.0", "1.24.0", false},
		{"1.25.0", "1.25.0", true},
		{"1.25.0", "1.26.0", false},
	}

	for _, tc := range cases {
		version1, _ := NewVersion(tc.versionA)
		version2, _ := NewVersion(tc.versionB)

		if tc.result != version1.Equal(version2) {
			t.Fatalf("error for comparison: %s is not equal to %s", tc.versionA, tc.versionB)
		}
	}
}

func TestCompareNotEqual(t *testing.T) {
	cases := []struct {
		versionA string
		versionB string
		result   bool
	}{
		{"1.25.0", "1.24.0", true},
		{"1.25.0", "1.25.0", false},
		{"1.25.0", "1.26.0", true},
	}

	for _, tc := range cases {
		version1, _ := NewVersion(tc.versionA)
		version2, _ := NewVersion(tc.versionB)

		if tc.result != version1.NotEqual(version2) {
			t.Fatalf("error for comparison: %s is equal to %s", tc.versionA, tc.versionB)
		}
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		versionA   string
		operator string
		versionB   string
		result     bool
	}{
		{"1.25.0", ">", "1.24.0", true},
		{"1.25.0", ">", "1.25.0", false},
		{"1.25.0", ">", "1.26.0", false},
		{"1.25.0", ">=", "1.24.0", true},
		{"1.25.0", ">=", "1.25.0", true},
		{"1.25.0", ">=", "1.26.0", false},
		{"1.25.0", "<", "1.24.0", false},
		{"1.25.0", "<", "1.25.0", false},
		{"1.25.0", "<", "1.26.0", true},
		{"1.25.0-beta2.1", "<", "1.25.0-b.3", true},
		{"1.25.0-b2.1", "<", "1.25.0beta.3", true},
		{"1.25.0-b-2.1", "<", "1.25.0-rc", true},
		{"1.25.0", "<=", "1.24.0", false},
		{"1.25.0", "<=", "1.25.0", true},
		{"1.25.0", "<=", "1.26.0", true},
		{"1.25.0", "==", "1.24.0", false},
		{"1.25.0", "==", "1.25.0", true},
		{"1.25.0", "==", "1.26.0", false},
		{"1.25.0-beta2.1", "==", "1.25.0-b.2.1", true},
		{"1.25.0beta2.1", "==", "1.25.0-b2.1", true},
		{"1.25.0", "=", "1.24.0", false},
		{"1.25.0", "=", "1.25.0", true},
		{"1.25.0", "=", "1.26.0", false},
		{"1.25.0", "!=", "1.24.0", true},
		{"1.25.0", "!=", "1.25.0", false},
		{"1.25.0", "!=", "1.26.0", true},
		{"1.25.0", "<>", "1.24.0", true},
		{"1.25.0", "<>", "1.25.0", false},
		{"1.25.0", "<>", "1.26.0", true},
	}

	for _, tc := range cases {
		version1, _ := NewVersion(tc.versionA)
		version2, _ := NewVersion(tc.versionB)

		if tc.result != version1.Compare(version2, tc.operator) {
			t.Fatalf("error for comparison: %s does not match constraint %s %s", tc.versionA, tc.operator, tc.versionB)
		}
	}
}
