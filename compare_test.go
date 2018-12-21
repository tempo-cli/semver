package semver

import (
	"fmt"
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
		t.Run(fmt.Sprintf("%s GreaterThan %s", tc.versionA, tc.versionB), func(t *testing.T) {
			version1, _ := NewVersion(tc.versionA)
			version2, _ := NewVersion(tc.versionB)

			assert.Equal(t, tc.result, version1.GreaterThan(version2), "%s is not greater than %s", tc.versionA, tc.versionB)
		})
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
		t.Run(fmt.Sprintf("%s GreaterThanOrEqual %s", tc.versionA, tc.versionB), func(t *testing.T) {
			version1, _ := NewVersion(tc.versionA)
			version2, _ := NewVersion(tc.versionB)

			assert.Equal(t, tc.result, version1.GreaterThanOrEqual(version2), "%s is not greater than or equal to %s", tc.versionA, tc.versionB)
		})
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
		t.Run(fmt.Sprintf("%s LessThan %s", tc.versionA, tc.versionB), func(t *testing.T) {
			version1, _ := NewVersion(tc.versionA)
			version2, _ := NewVersion(tc.versionB)

			assert.Equal(t, tc.result, version1.LessThan(version2), "%s is not less than %s", tc.versionA, tc.versionB)
		})
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
		t.Run(fmt.Sprintf("%s LessThanOrEqual %s", tc.versionA, tc.versionB), func(t *testing.T) {
			version1, _ := NewVersion(tc.versionA)
			version2, _ := NewVersion(tc.versionB)

			assert.Equal(t, tc.result, version1.LessThanOrEqual(version2), "%s is not less than or equal to %s", tc.versionA, tc.versionB)
		})
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
		t.Run(fmt.Sprintf("%s Equal %s", tc.versionA, tc.versionB), func(t *testing.T) {
			version1, _ := NewVersion(tc.versionA)
			version2, _ := NewVersion(tc.versionB)

			assert.Equal(t, tc.result, version1.Equal(version2), "%s is not equal to %s", tc.versionA, tc.versionB)
		})
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
		t.Run(fmt.Sprintf("%s NotEqual %s", tc.versionA, tc.versionB), func(t *testing.T) {
			version1, _ := NewVersion(tc.versionA)
			version2, _ := NewVersion(tc.versionB)

			assert.Equal(t, tc.result, version1.NotEqual(version2), "Expected %s and %s to not be equal", tc.versionA, tc.versionB)
		})
	}
}

func TestCompare(t *testing.T) {
	for _, tc := range getComparisonCases() {
		t.Run(fmt.Sprintf("%s %s %s", tc.versionA, tc.operator, tc.versionB), func(t *testing.T) {
			version1, _ := NewVersion(tc.versionA)
			version2, _ := NewVersion(tc.versionB)

			assert.Equal(t, tc.result, version1.Compare(version2, tc.operator), "failed comparing that %s %s %s", tc.versionA, tc.operator, tc.versionB)
		})
	}
}

func BenchmarkCompare(b *testing.B) {
	for _, tc := range getComparisonCasesForBench() {
		b.Run(fmt.Sprintf("Compare versions: %s %s %s", tc.versionA.String(), tc.operator, tc.versionB.String()), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				tc.versionA.Compare(tc.versionB, tc.operator)
			}
		})
	}
}

func getComparisonCases() []struct {
	versionA string
	operator string
	versionB string
	result   bool
} {
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
	return cases
}

func getComparisonCasesForBench() []struct {
	versionA Version
	operator string
	versionB Version
} {
	cases := []struct {
		versionA Version
		operator string
		versionB Version
	}{
		{&semver{Major: 1, Minor: 25, Patch: 9}, ">", &semver{Major: 1, Minor: 24, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, ">", &semver{Major: 1, Minor: 25, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, ">", &semver{Major: 1, Minor: 26, Patch: 9}},

		{&semver{Major: 1, Minor: 25, Patch: 9}, ">=", &semver{Major: 1, Minor: 24, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, ">=", &semver{Major: 1, Minor: 25, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, ">=", &semver{Major: 1, Minor: 26, Patch: 9}},

		{&semver{Major: 1, Minor: 25, Patch: 9}, "<", &semver{Major: 1, Minor: 24, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "<", &semver{Major: 1, Minor: 25, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "<", &semver{Major: 1, Minor: 26, Patch: 9}},

		{&semver{Major: 1, Minor: 25, Patch: 9, Stability: "beta", PreRelease: "2.1"}, "<", &semver{Major: 1, Minor: 25, Patch: 9, Stability: "beta", PreRelease: "3"}},
		{&semver{Major: 1, Minor: 25, Patch: 9, Stability: "beta", PreRelease: "2.1"}, "<", &semver{Major: 1, Minor: 25, Patch: 9, Stability: "rc"}},

		{&semver{Major: 1, Minor: 25, Patch: 9}, ">=", &semver{Major: 1, Minor: 24, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, ">=", &semver{Major: 1, Minor: 25, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, ">=", &semver{Major: 1, Minor: 26, Patch: 9}},

		{&semver{Major: 1, Minor: 25, Patch: 9}, "==", &semver{Major: 1, Minor: 22, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "==", &semver{Major: 1, Minor: 25, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "==", &semver{Major: 1, Minor: 26, Patch: 9}},

		{&semver{Major: 1, Minor: 25, Patch: 9, Stability: "beta", PreRelease: "2.1"}, "==", &semver{Major: 1, Minor: 25, Patch: 9, Stability: "beta", PreRelease: "2.1"}},

		{&semver{Major: 1, Minor: 25, Patch: 9}, "=", &semver{Major: 1, Minor: 24, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "=", &semver{Major: 1, Minor: 25, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "=", &semver{Major: 1, Minor: 26, Patch: 9}},

		{&semver{Major: 1, Minor: 25, Patch: 9}, "!=", &semver{Major: 1, Minor: 24, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "!=", &semver{Major: 1, Minor: 25, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "!=", &semver{Major: 1, Minor: 26, Patch: 9}},

		{&semver{Major: 1, Minor: 25, Patch: 9}, "<>", &semver{Major: 1, Minor: 24, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "<>", &semver{Major: 1, Minor: 25, Patch: 9}},
		{&semver{Major: 1, Minor: 25, Patch: 9}, "<>", &semver{Major: 1, Minor: 26, Patch: 9}},
	}

	return cases
}
