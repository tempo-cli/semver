package semver

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"math"
	"regexp"
	"strings"
)

type Constraint struct {
	operator    string
	version     *Version
	constraints []*Constraint
	conjunctive bool
	isEmpty     bool
}

var (
	versionReg             = `v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?` + stabilityRegex + `?([.-]?dev)?(?:\+[^\s]+)?`
	operatorMap            = map[string]string{"=": "==", "==": "==", "<>": "!=", "!=": "!=", ">": ">", "<": "<", "<=": "<=", ">=": ">="}
	stabilityModifierRegex = regexp.MustCompile(`(?i)^([^,\s]*?)@(stable|RC|beta|alpha|dev)$`)
	simpleComparisonRegex  = regexp.MustCompile(`^(<>|!=|>=?|<=?|==?)?\s*(.*)`)
	xRangeRegex            = regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.[xX*])+$`)
	tildeRegex             = regexp.MustCompile(`(?i)^~>?` + versionReg + `$`)
	caretRegex             = regexp.MustCompile(`(?i)^\^` + versionReg + `$`)
	hyphenRegex            = regexp.MustCompile(`(?i)^(` + versionReg + `) +- +(` + versionReg + `)($)`)
	devConstraintRegex     = regexp.MustCompile(`(?i)^(dev-[^,\s@]+?|[^,\s@]+?\.x-dev)#.+$`)
	orSplitRegex           = regexp.MustCompile(`\s*\|\|?\s*`)
	andConstraintRegex     = regexp.MustCompile(`\s*[ ,]\s*`)
)

func NewConstraint(constraint string) (*Constraint, error) {
	var version = constraint

	result := stabilityModifierRegex.FindStringSubmatch(constraint)

	if nil != result {
		version = result[1]
	}

	result = devConstraintRegex.FindStringSubmatch(constraint)

	if nil != result {
		version = result[1]
	}

	orConstraints := orSplitRegex.Split(version, -1)
	var orGroups []*Constraint

	for _, constraints := range orConstraints {

		andConstraints := parseAndConstraints(constraints)

		if len(andConstraints) > 1 {
			andRange := &Constraint{conjunctive: true}

			for _, constraints := range andConstraints {
				c, err := parseConstraint(constraints)

				if nil != err {
					return nil, err
				}

				if len(c.constraints) > 0 {
					andRange.constraints = append(andRange.constraints, c.constraints...)
				} else {
					andRange.constraints = append(andRange.constraints, c)
				}
			}

			orGroups = append(orGroups, andRange)

		} else {
			c, err := parseConstraint(constraints)

			if nil != err {
				return nil, err
			}

			orGroups = append(orGroups, c)
		}

	}

	if 1 == len(orGroups) {
		return orGroups[0], nil
	} else if 2 == len(orGroups) &&
		// parse the two OR groups and if they are contiguous we collapse
		// them into one constraint
		2 == len(orGroups[0].constraints) &&
		2 == len(orGroups[1].constraints) &&
		">=" == orGroups[0].constraints[0].operator &&
		"<" == orGroups[0].constraints[1].operator &&
		">=" == orGroups[1].constraints[0].operator &&
		"<" == orGroups[1].constraints[1].operator &&
		orGroups[0].constraints[1].version.String() == orGroups[1].constraints[0].version.String() {

		return &Constraint{constraints: []*Constraint{orGroups[0].constraints[0], orGroups[1].constraints[1]}, conjunctive: false}, nil
	}

	return &Constraint{conjunctive: false, constraints: orGroups}, nil
}

func (c *Constraint) Matches(version *Version) bool {
	if len(c.constraints) > 0 {
		if false == c.conjunctive {
			for _, c := range c.constraints {
				if c.Matches(version) {
					return true
				}

			}

			return false
		}

		for _, c := range c.constraints {
			if !c.Matches(version) {
				return false
			}
		}

		return true
	}

	if c.isEmpty {
		return true
	}

	return version.Compare(c.version, c.operator)
}

func (c *Constraint) String() string {
	if c.isEmpty {
		return "[]"
	}

	if 0 == len(c.constraints) {
		result := fmt.Sprintf("%s %s", c.operator, c.version.String())
		if "-stable" == result[len(result)-7:] {
			result = result[0 : len(result)-7]
		}
		return result
	}

	var glue string

	if true == c.conjunctive {
		glue = " "
	} else {
		glue = " || "
	}

	totalConstraints := len(c.constraints)

	if 1 == totalConstraints {
		constraint := c.constraints[0]

		if constraint.isEmpty {
			return "[]"
		}

		return constraint.String()
	}

	var constraints = make([]string, 0, totalConstraints)
	for _, constraint := range c.constraints {
		constraints = append(constraints, constraint.String())
	}

	return "[" + strings.Join(constraints, glue) + "]"
}

func parseAndConstraints(constraint string) []string {
	var (
		index          = 0
		constraintPart = 0
		constraints    []string
	)

	split := andConstraintRegex.Split(constraint, -1)

	if 1 == len(split) {
		return []string{constraint}
	}

	var parts []string
	for _, str := range split {
		if str != "" {
			parts = append(parts, str)
		}
	}

	partsLen := len(parts)

	for {

		if index >= partsLen {
			break
		}

		constraints = append(constraints, "")

		if "<" == parts[index] || ">" == parts[index] || ">=" == parts[index] || "<=" == parts[index] || "^" == parts[index] || "!=" == parts[index] {
			constraints[constraintPart] += parts[index]
			index++

			if index >= partsLen {
				break
			}
		}

		constraints[constraintPart] += parts[index]

		if index+1 >= partsLen {
			break
		}

		if "as" == parts[index+1] || "-" == parts[index+1] {
			index++
			constraints[constraintPart] += " " + parts[index]

			index++
			constraints[constraintPart] += " " + parts[index]

		}

		index++
		constraintPart++
	}

	return constraints
}

func parseConstraint(constraint string) (*Constraint, error) {
	b := []byte(constraint)

	if match, _ := regexp.Match(`^v?[xX*](\.[xX*])*$`, b); match {
		return &Constraint{isEmpty: true}, nil
	}

	if tildeRegex.Match(b) {
		return parseTilde(constraint)
	}

	if caretRegex.Match(b) {
		return caretRange(constraint)
	}

	if xRangeRegex.Match(b) {
		return xRange(constraint)
	}

	if hyphenRegex.Match(b) {
		return hyphenRange(constraint)
	}

	if simpleComparisonRegex.Match(b) {
		return basicRange(constraint)
	}

	return nil, fmt.Errorf("unable to parse constraint %s", constraint)
}

func basicRange(constraint string) (*Constraint, error) {
	result := stabilityModifierRegex.FindStringSubmatch(constraint)

	var stability = ""

	if nil != result {
		constraint = result[1]
		if "stable" != result[2] {
			stability = result[2]
		}
	}

	matches := simpleComparisonRegex.FindStringSubmatch(constraint)

	version := matches[2]

	if "" != stability && "stable" == ParseStability(version) {
		version += "-" + stability
	} else if "<" == matches[1] || ">=" == matches[1] {
		if !stabilityRegexC.Match([]byte(version)) {
			if len(version) < 4 || "dev-" != version[0:4] {
				version += "-dev"
			}
		}
	}

	v, err := NewVersion(version)

	if nil != err {
		return nil, err
	}

	operator := matches[1]

	if "" == operator {
		operator = "="
	}

	return &Constraint{operator: operatorMap[operator], version: v, conjunctive: true}, nil
}

/*
 Hyphen Range

 Specifies an inclusive set. If a partial version is provided as the first version in the inclusive range,
 then the missing pieces are replaced with zeroes. If a partial version is provided as the second version in
 the inclusive range, then all versions that start with the supplied parts of the tuple are accepted, but
 nothing that would be greater than the provided tuple parts.
*/
func hyphenRange(constraint string) (*Constraint, error) {
	matches := hyphenRegex.FindStringSubmatch(constraint)

	c := &Constraint{conjunctive: true}

	// Calculate the stability suffix
	var lowStabilitySuffix = ""
	if "" == matches[6] && "" == matches[8] {
		lowStabilitySuffix = "-dev"
	}
	lowVersion, err := NewVersion(matches[1] + lowStabilitySuffix)

	if nil != err {
		return nil, err
	}

	c.constraints = append(c.constraints, &Constraint{operator: ">=", version: lowVersion, conjunctive: true})
	var isEmpty = func(x string) bool {
		if "0" == x {
			return false
		}

		return "" == x
	}

	if (!isEmpty(matches[11]) && !isEmpty(matches[12])) || "" != matches[14] || "" != matches[16] {
		highVersion, err := NewVersion(matches[9])
		if nil != err {
			return nil, err
		}

		c.constraints = append(c.constraints, &Constraint{operator: "<=", version: highVersion, conjunctive: true})
	} else {
		var position = 0
		highMatch := []string{"", matches[10], matches[11], matches[12], matches[13]}
		if isEmpty(matches[11]) {
			position = 1
		} else {
			position = 2
		}
		highVersion, err := expandVersion(highMatch, position, 1, "0", "-dev")

		if nil != err {
			return nil, err
		}

		c.constraints = append(c.constraints, &Constraint{operator: "<", version: highVersion, conjunctive: true})
	}

	return c, nil
}

/*
 X Range

 Any of X, x, or * may be used to "stand in" for one of the numeric values in the [major, minor, patch] tuple.
 A partial version range is treated as an X-Range, so the special character is in fact optional.
*/
func xRange(constraint string) (*Constraint, error) {
	matches := xRangeRegex.FindStringSubmatch(constraint)
	position := 0

	for i := 3; i > 0; i-- {
		if "" != matches[i] {
			position = i
			break
		}
	}

	lowVersion, err := expandVersion(matches, position, 0, "0", "-dev")

	if nil != err {
		return nil, err
	}

	highVersion, err := expandVersion(matches, position, 1, "0", "-dev")

	if nil != err {
		return nil, err
	}

	if lowVersion.String() == "0.0.0.0-dev" {
		return &Constraint{operator: "<", version: highVersion, conjunctive: true}, nil
	}

	return &Constraint{constraints: []*Constraint{{operator: ">=", version: lowVersion, conjunctive: true}, {operator: "<", version: highVersion, conjunctive: true}}, conjunctive: true}, nil
}

/*
 Caret Range

 Allows changes that do not modify the left-most non-zero digit in the [major, minor, patch] tuple.
 In other words, this allows patch and minor updates for versions 1.0.0 and above, patch updates for
 versions 0.X >=0.1.0, and no updates for versions 0.0.X
*/
func caretRange(constraint string) (*Constraint, error) {
	matches := caretRegex.FindStringSubmatch(constraint)
	stabilitySuffix := ""
	position := 0

	if "0" != matches[1] || "" == matches[2] {
		position = 1
	} else if "0" != matches[2] || "" == matches[3] {
		position = 2
	} else {
		position = 3
	}

	if "" == matches[5] && "" == matches[7] {
		stabilitySuffix = "-dev"
	}

	lowVersion, err := NewVersion(constraint[1:] + stabilitySuffix)

	if nil != err {
		return nil, err
	}
	// For upper bound, we increment the position of one more significance,
	// but highPosition = 0 would be illegal
	highVersion, err := expandVersion(matches, position, 1, "0", "-dev")

	if nil != err {
		return nil, err
	}

	return &Constraint{constraints: []*Constraint{{operator: ">=", version: lowVersion, conjunctive: true}, {operator: "<", version: highVersion, conjunctive: true}}, conjunctive: true}, nil
}

/*
 Tilde Range

 Like wildcard constraints, unsuffixed tilde constraints say that they must be greater than the previous
 version, to ensure that unstable instances of the current version are allowed. However, if a stability
 suffix is added to the constraint, then a >= match on the current version is used instead.
*/
func parseTilde(constraint string) (*Constraint, error) {
	matches := tildeRegex.FindStringSubmatch(constraint)

	if "~>" == constraint[0:2] {
		return nil, errors.New(fmt.Sprintf(`Could not parse version constraint %s: `+
			`Invalid operator "~>", you probably meant to use the "~" operator`, constraint))
	}
	position := 0

	for i := 4; i > 0; i-- {
		if "" != matches[i] {
			position = i
			break
		}
	}

	stabilitySuffix := ""
	if "" != matches[5] {
		stabilitySuffix = "-" + expandStability(matches[5]) + matches[6]
	}

	if "" != matches[7] || "" == stabilitySuffix {
		stabilitySuffix = "-dev"
	}

	lowVersion, err := expandVersion(matches, position, 0, "0", stabilitySuffix)

	if nil != err {
		return nil, err
	}

	// For upper bound, we increment the position of one more significance,
	// but highPosition = 0 would be illegal
	highPosition := math.Max(1, cast.ToFloat64(position-1))
	highVersion, err := expandVersion(matches, cast.ToInt(highPosition), 1, "0", "-dev")

	if nil != err {
		return nil, err
	}

	return &Constraint{constraints: []*Constraint{{operator: ">=", version: lowVersion, conjunctive: true}, {operator: "<", version: highVersion, conjunctive: true}}, conjunctive: true}, nil
}

func expandVersion(matches []string, position int, increment int, pad string, append string) (*Version, error) {
	var (
		i      = 4
		result = make([]interface{}, 5, 5)
	)

	for i > 0 {
		if i > position {
			result[i-1] = pad
		} else if i == position && increment > 0 {
			currentValue := cast.ToInt(matches[i])
			result[i-1] = cast.ToString(currentValue + increment)
			if currentValue < 0 {
				result[i-1] = pad
				position--
				if i == 1 {
					return nil, fmt.Errorf("carry overflow error")
				}
			}
		} else {
			result[i-1] = matches[i]
		}
		i--
	}

	result[4] = append

	version, err := NewVersion(fmt.Sprintf("%s.%s.%s.%s%s", result...))

	if nil != err {
		return nil, err
	}

	return version, nil
}
