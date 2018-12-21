package semver

import (
	"github.com/stretchr/testify/assert"
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

		assert.Equal(t, tc.result, version1.GreaterThan(version2), "%s is not greater than %s", tc.versionA, tc.versionB)
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

		assert.Equal(t, tc.result, version1.GreaterThanOrEqual(version2), "%s is not greater than or equal to %s", tc.versionA, tc.versionB)
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

		assert.Equal(t, tc.result, version1.LessThan(version2), "%s is not less than %s", tc.versionA, tc.versionB)
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

		assert.Equal(t, tc.result, version1.LessThanOrEqual(version2), "%s is not less than or equal to %s", tc.versionA, tc.versionB)
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

		assert.Equal(t, tc.result, version1.Equal(version2), "%s is not equal to %s", tc.versionA, tc.versionB)
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

		assert.Equal(t, tc.result, version1.NotEqual(version2), "Expected %s and %s to not be equal", tc.versionA, tc.versionB)
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		versionA string
		operator string
		versionB string
		result   bool
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

		assert.Equal(t, tc.result, version1.Compare(version2, tc.operator), "failed comparing that %s %s %s", tc.versionA, tc.operator, tc.versionB)
	}
}
