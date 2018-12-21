package semver

import (
	"fmt"
	"testing"
)

func TestNewVersion(t *testing.T) {
	cases := []struct {
		name       string
		version    string
		normalized string
	}{
		// Use all test cases from composer to ensure feature parity
		{"none", "1.0.0", "1.0.0.0"},
		{"none/2", "1.2.3.4", "1.2.3.4"},
		{"parses state", "1.0.0RC1dev", "1.0.0.0-RC1-dev"},
		{"CI parsing", "1.0.0-rC15-dev", "1.0.0.0-RC15-dev"},
		{"delimiters", "1.0.0.RC.15-dev", "1.0.0.0-RC15-dev"},
		{"RC uppercase", "1.0.0-rc1", "1.0.0.0-RC1"},
		{"patch replace", "1.0.0.pl3-dev", "1.0.0.0-patch3-dev"},
		{"forces w.x.y.z", "1.0-dev", "1.0.0.0-dev"},
		{"forces w.x.y.z/2", "0", "0.0.0.0"},
		{"parses long", "10.4.13-beta", "10.4.13.0-beta"},
		{"parses long/2", "10.4.13beta2", "10.4.13.0-beta2"},
		{"parses long/semver", "10.4.13beta.2", "10.4.13.0-beta2"},
		{"expand shorthand", "10.4.13-b", "10.4.13.0-beta"},
		{"expand shorthand/2", "10.4.13-b5", "10.4.13.0-beta5"},
		{"strips leading v", "v1.0.0", "1.0.0.0"},
		{"parses dates y-m as classical", "2010.01", "2010.01.0.0"},
		{"parses dates w/ . as classical", "2010.01.02", "2010.01.02.0"},
		{"parses dates y.m.Y as classical", "2010.1.555", "2010.1.555.0"},
		{"parses dates y.m.Y/2 as classical", "2010.10.200", "2010.10.200.0"},
		{"strips v/datetime", "v20100102", "20100102"},
		{"parses dates w/ -", "2010-01-02", "2010.01.02"},
		{"parses numbers", "2010-01-02.5", "2010.01.02.5"},
		{"parses dates y.m.Y", "2010.1.555", "2010.1.555.0"},
		{"parses datetime", "20100102-203040", "20100102.203040"},
		{"parses dt+number", "20100102203040-10", "20100102203040.10"},
		{"parses dt+patch", "20100102-203040-p1", "20100102.203040-patch1"},
		{"parses master", "dev-master", "9999999-dev"},
		{"parses trunk", "dev-trunk", "9999999-dev"},
		{"parses branches", "1.x-dev", "1.9999999.9999999.9999999-dev"},
		{"parses arbitrary", "dev-feature-foo", "dev-feature-foo"},
		{"parses arbitrary/2", "DEV-FOOBAR", "dev-FOOBAR"},
		{"parses arbitrary/3", "dev-feature/foo", "dev-feature/foo"},
		{"parses arbitrary/4", "dev-feature+issue-1", "dev-feature+issue-1"},
		{"ignores aliases", "dev-master as 1.0.0", "9999999-dev"},
		{"semver metadata/2", "1.0.0-beta.5+foo", "1.0.0.0-beta5"},
		{"semver metadata/3", "1.0.0+foo", "1.0.0.0"},
		{"semver metadata/4", "1.0.0-alpha.3.1+foo", "1.0.0.0-alpha3.1"},
		{"semver metadata/5", "1.0.0-alpha2.1+foo", "1.0.0.0-alpha2.1"},
		{"semver metadata/6", "1.0.0-alpha-2.1-3+foo", "1.0.0.0-alpha2.1-3"},
		// not supported for BC {"semver metadata/7", "1.0.0-0.3.7", "1.0.0.0-0.3.7"},
		// not supported for BC {"semver metadata/8", "1.0.0-x.7.z.92", "1.0.0.0-x.7.z.92"},
		{"metadata w/ alias", "1.0.0+foo as 2.0", "1.0.0.0"},
	}

	for _, tc := range cases {
		version, err := NewVersion(tc.version)

		if err != nil {
			t.Fatalf("error for version %s: %s", tc.version, err)
		}

		normalized := version.String()
		if normalized != tc.normalized {
			fmt.Println(version)
			t.Fatalf("error for version %s: %s is not equal to %s", tc.version, normalized, tc.normalized)
		}
	}
}

func TestFailedNewVersion(t *testing.T) {
	cases := []struct {
		name    string
		version string
	}{
		// Use all test cases from composer to ensure feature parity
		{"empty", ""},
		{"invalid chars", "a"},
		{"invalid type", "1.0.0-meh"},
		{"too many bits", "1.0.0.0.0"},
		{"non-dev arbitrary", "feature-foo"},
		{"metadata w/ space", "1.0.0+foo bar"},
		{"maven style release", "1.0.1-SNAPSHOT"},
	}

	for _, tc := range cases {
		_, err := NewVersion(tc.version)

		if err == nil {
			t.Fatalf("error for version %s", tc.version)
		}
	}
}

func TestNormalizeBranch(t *testing.T) {
	cases := []struct {
		name       string
		version    string
		normalized string
	}{
		// Branches
		{"parses x", "v1.x", "1.9999999.9999999.9999999-dev"},
		{"parses *", "v1.*", "1.9999999.9999999.9999999-dev"},
		{"parses digits", "v1.0", "1.0.9999999.9999999-dev"},
		{"parses digits/2", "2.0", "2.0.9999999.9999999-dev"},
		{"parses long x", "v1.0.x", "1.0.9999999.9999999-dev"},
		{"parses long *", "v1.0.3.*", "1.0.3.9999999-dev"},
		{"parses long digits", "v2.4.0", "2.4.0.9999999-dev"},
		{"parses long digits/2", "2.4.4", "2.4.4.9999999-dev"},
		{"parses master", "master", "9999999-dev"},
		{"parses trunk", "trunk", "9999999-dev"},
		{"parses arbitrary", "feature-a", "dev-feature-a"},
		{"parses arbitrary/2", "FOOBAR", "dev-FOOBAR"},
		{"parses arbitrary/3", "feature+issue-1", "dev-feature+issue-1"},
	}

	for _, tc := range cases {
		version, err := NormalizeBranch(tc.version)

		if err != nil {
			t.Fatalf("error for version %s: %s", tc.version, err)
		}

		normalized := version.String()
		if normalized != tc.normalized {
			fmt.Println(version)
			t.Fatalf("error for version %s: %s is not equal to %s", tc.version, normalized, tc.normalized)
		}
	}
}

func TestParseStability(t *testing.T) {
	cases := []struct {
		stablility string
		version    string
	}{
		// Stablility
		{"stable", "1"},
		{"stable", "1.0"},
		{"stable", "3.2.1"},
		{"stable", "v3.2.1"},
		{"dev", "v2.0.x-dev"},
		{"dev", "v2.0.x-dev#abc123"},
		{"dev", "v2.0.x-dev#trunk/@123"},
		{"RC", "3.0-RC2"},
		{"dev", "dev-master"},
		{"dev", "3.1.2-dev"},
		{"dev", "dev-feature+issue-1"},
		{"stable", "3.1.2-p1"},
		{"stable", "3.1.2-pl2"},
		{"stable", "3.1.2-patch"},
		{"alpha", "3.1.2-alpha5"},
		{"beta", "3.1.2-beta"},
		{"beta", "2.0B1"},
		{"alpha", "1.2.0a1"},
		{"alpha", "1.2_a1"},
		{"RC", "2.0.0rc1"},
		{"alpha", "1.0.0-alpha11+cs-1.1.0"},
	}

	for _, tc := range cases {
		stablility := ParseStability(tc.version)

		if stablility != tc.stablility {
			t.Fatalf("error for stablility %s: %s is not equal to %s", tc.version, stablility, tc.stablility)
		}
	}
}
