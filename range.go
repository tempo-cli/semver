package semver

import (
	"regexp"
	"strings"
)

type Range struct {
	conjunctive bool
	constraints []*Constraint
	ranges      []*Range
}

var (
	removeStabilityRegex = regexp.MustCompile("(?i)^([^,\\s]*?)@(stable|RC|beta|alpha|dev)$")
	devConstraintRegex   = regexp.MustCompile("(?i)^(dev-[^,\\s@]+?|[^,\\s@]+?\\.x-dev)#.+$")
	orSplitRegex         = regexp.MustCompile("\\s*\\|\\|?\\s*")
	andConstraintRegex   = regexp.MustCompile("\\s*[ ,]\\s*")
)

func NewRange(constraint string) (*Range, error) {
	var version = constraint

	result := removeStabilityRegex.FindStringSubmatch(constraint)

	if nil != result {
		version = result[1]
	}

	result = devConstraintRegex.FindStringSubmatch(constraint)

	if nil != result {
		version = result[1]
	}

	orConstraints := orSplitRegex.Split(version, -1)
	var orGroups []*Range

	var r = &Range{conjunctive: false}
	for _, constraints := range orConstraints {

		andConstraints := parseAndConstraints(constraints)

		if len(andConstraints) > 1 {
			andRange := &Range{conjunctive: true}

			if len(r.constraints) > 0 {
				r.ranges = append(r.ranges, &Range{constraints: r.constraints, conjunctive: true})
				r.constraints = []*Constraint{}
			}

			for _, constraints := range andConstraints {
				c, err := NewConstraint(constraints)

				if nil != err {
					return nil, err
				}

				andRange.constraints = append(andRange.constraints, c...)
			}

			r.ranges = append(r.ranges, andRange)
		} else {
			c, err := NewConstraint(constraints)

			if nil != err {
				return nil, err
			}

			if len(r.ranges) > 0 {
				r.ranges = append(r.ranges, &Range{constraints: c, conjunctive: false})
			} else {
				r.constraints = append(r.constraints, c...)
			}
			orGroups = append(orGroups, &Range{constraints: c, conjunctive: false})
		}
	}

	if
	2 == len(orGroups) &&
	// parse the two OR groups and if they are contiguous we collapse
	// them into one constraint
		2 == len(orGroups[0].constraints) &&
		2 == len(orGroups[1].constraints) &&
		">=" == orGroups[0].constraints[0].operator &&
		"<" == orGroups[0].constraints[1].operator &&
		">=" == orGroups[1].constraints[0].operator &&
		"<" == orGroups[1].constraints[1].operator &&
		orGroups[0].constraints[1].version.String() == orGroups[1].constraints[0].version.String() {

		c := make([]*Constraint, 0, 2)
		c = append(c, &Constraint{">=", orGroups[0].constraints[0].version, false})
		c = append(c, &Constraint{"<", orGroups[1].constraints[1].version, false})

		return &Range{constraints: c, conjunctive: false}, nil
	}

	return r, nil
}

func (c *Range) String() string {

	glue := ""

	if c.conjunctive {
		glue = " "
	} else {
		glue = " || "
	}

	if len(c.ranges) > 0 {
		if 1 == len(c.ranges) {
			return c.ranges[0].String()
		}

		var ranges = make([]string, 0, len(c.ranges))
		for _, r := range c.ranges {
			ranges = append(ranges, r.String())
		}

		return "[" + strings.Join(ranges, glue) + "]"
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

		if "<" == parts[index] || ">" == parts[index] || ">=" == parts[index] || "<=" == parts[index] || "^" == parts[index] {
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
