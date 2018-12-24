package semver

import (
	"fmt"
	"github.com/spf13/cast"
	"regexp"
	"strings"
)

var (
	aliasRegex = regexp.MustCompile(`^([^,\s]+)\s+as\s+([^,\s]+)$`)

	stabilityRegex = `(?i)[._-]?(?:(dev|stable|beta|b|RC|alpha|a|patch|pl|p)((?:[.-]?\d+)*)?)`

	branchRegex = regexp.MustCompile(`^v?(\d+)(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?$`)

	versionRegex =
	// Match normal version string (1.2.3)
		`^v?([0-9]{1,5})(\.[0-9]+)?(\.[0-9]+)?(\.[0-9]+)?` +

		// Match pre-release info (-beta.2). This supports dot, underscore, dash or nothing as a prefix to match Composers rules
			stabilityRegex + "?([.-]?dev)?"

	// Match metadata (E.G +build.1234)
	versionRegexC = regexp.MustCompile(versionRegex + `(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

	dateTimeRegex = regexp.MustCompile(`^v?(\d{4}(?:[.:-]?\d{2}){1,6}(?:[.:-]?\d{1,3})?)` + stabilityRegex + `?$`)

	stabilityRegexC = regexp.MustCompile(stabilityRegex)
	branchMatcher   = regexp.MustCompile(`(?i)(.*?)[.-]?dev$`)
	replaceRegex    = regexp.MustCompile(`([^0-9]+)`)
)

func NewVersion(version string) (*Version, error) {
	originalVersion := version
	alias := aliasRegex.FindStringSubmatch(version)

	if alias != nil {
		version = alias[1]
	}

	if match, _ := regexp.Match("(?i)^(?:dev-)?(?:master|trunk|default)$", []byte(version)); match {
		return &Version{
			Major:     9999999,
			Stability: "dev",
			State:     "",
			Original:  originalVersion,
		}, nil
	}

	if len(version) > 4 && "dev-" == strings.ToLower(version[0:4]) {
		return &Version{
			Parsed:   fmt.Sprintf("dev-%s", version[4:]),
			Original: originalVersion,
			isBranch: true,
		}, nil
	}

	versionMatch := versionRegexC.FindStringSubmatch(version)

	if versionMatch != nil {
		stability := expandStability(versionMatch[5])
		return &Version{
			Major:      cast.ToInt(versionMatch[1]),
			Minor:      parseVersionNumber(versionMatch[2]),
			Patch:      parseVersionNumber(versionMatch[3]),
			Extra:      parseVersionNumber(versionMatch[4]),
			PreRelease: strings.TrimLeft(versionMatch[6], ".-"),
			Stability:  stability,
			State:      strings.TrimLeft(versionMatch[7], "-"),
			Metadata:   versionMatch[9],
			Original:   originalVersion,
		}, nil
	}

	dateTimeMatch := dateTimeRegex.FindStringSubmatch(version)

	if dateTimeMatch != nil {

		versionString := replaceRegex.ReplaceAllString(dateTimeMatch[1], `.`)

		return &Version{
			Stability: expandStability(dateTimeMatch[2]),
			Patch:     cast.ToInt(dateTimeMatch[3]),
			Parsed:    versionString,
			Original:  originalVersion,
			isDate:    true,
		}, nil
	}

	branchMatches := branchMatcher.FindStringSubmatch(version)
	if nil != branchMatches {
		return NormalizeBranch(branchMatches[1])
	}

	return nil, fmt.Errorf("unable to parse version %s", version)
}

func NormalizeBranch(branch string) (*Version, error) {

	valid := map[string]bool{"master": true, "trunk": true, "default": true}

	if valid[branch] {
		return NewVersion(branch)
	}

	branchMatches := branchRegex.FindStringSubmatch(branch)

	if nil != branchMatches {
		versionString := ""
		matchesLength := len(branchMatches)

		for i := 1; i < 5; i++ {
			if i < matchesLength && "" != branchMatches[i] {
				versionString += strings.Replace(strings.Replace(branchMatches[i], "X", "x", -1), "*", "x", -1)
			} else {
				versionString += ".x"
			}
		}

		return NewVersion(strings.Replace(versionString, "x", "9999999", -1) + "-dev")
	}

	return NewVersion("dev-" + branch)
}

func expandStability(stability string) string {
	switch strings.ToLower(stability) {
	case "alpha", "a":
		return "alpha"
	case "beta", "b":
		return "beta"
	case "p", "pl":
		return "patch"
	case "rc":
		return "RC"
	}

	return stability
}

func ParseStability(stability string) string {
	if "" == stability {
		return stability
	}

	if len(stability) >= 4 && ("dev-" == strings.ToLower(stability[0:4]) || "-dev" == strings.ToLower(stability[len(stability)-4:])) {
		return "dev"
	}

	stabilityMatch := stabilityRegexC.FindStringSubmatch(stability)

	if nil != stabilityMatch {
		switch strings.ToLower(stabilityMatch[1]) {
		case "alpha", "a":
			return "alpha"
		case "beta", "b":
			return "beta"
		case "rc":
			return "RC"
		case "dev":
			return "dev"
		}
	}

	return "stable"
}

func parseVersionNumber(version string) int {
	if "" == version {
		return 0
	}

	return cast.ToInt(strings.TrimPrefix(version, "."))
}
