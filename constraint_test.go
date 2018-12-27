package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseConstraintsIgnoresStabilityFlag(t *testing.T) {
	constraint, err := NewConstraint("1.0@dev")

	if assert.NoError(t, err) {
		assert.Equal(t, "== 1.0.0.0", constraint.String())
	}
}

func TestParseConstraintsIgnoresReferenceOnDevVersion(t *testing.T) {
	constraint, err := NewConstraint("1.0.x-dev#abcd123")

	if assert.NoError(t, err) {
		assert.Equal(t, "== 1.0.9999999.9999999-dev", constraint.String())
	}

	constraint, err = NewConstraint("1.0.x-dev#trunk/@123")
	if assert.NoError(t, err) {
		assert.Equal(t, "== 1.0.9999999.9999999-dev", constraint.String())
	}
}

func TestParseConstraintsFailsOnBadReference(t *testing.T) {
	_, err := NewConstraint("1.0#abcd123")

	assert.Error(t, err)

	_, err = NewConstraint("1.0#trunk/@123")
	assert.Error(t, err)
}

func TestParseConstraintsNudgesRubyDevsTowardsThePathOfRighteousness(t *testing.T) {
	_, err := NewConstraint("~>1.2")

	assert.EqualError(t, err, "Could not parse version constraint ~>1.2: Invalid operator \"~>\", you probably meant to use the \"~\" operator")
}

func TestParseConstraintsSimple(t *testing.T) {
	cases := []struct {
		name       string
		constraint string
		parsed     string
	}{
		{"match any", "*", "[]"},
		{"match any/2", "*.*", "[]"},
		{"match any/2v", "v*.*", "[]"},
		{"match any/3", "*.x.*", "[]"},
		{"match any/4", "x.X.x.*", "[]"},
		{"not equal", "<>1.0.0", "!= 1.0.0.0"},
		{"not equal/2", "!=1.0.0", "!= 1.0.0.0"},
		{"greater than", ">1.0.0", "> 1.0.0.0"},
		{"lesser than", "<1.2.3.4", "< 1.2.3.4-dev"},
		{"less/eq than", "<=1.2.3", "<= 1.2.3.0"},
		{"great/eq than", ">=1.2.3", ">= 1.2.3.0-dev"},
		{"equals", "=1.2.3", "== 1.2.3.0"},
		{"double equals", "==1.2.3", "== 1.2.3.0"},
		{"no op means eq", "1.2.3", "== 1.2.3.0"},
		{"completes version", "=1.0", "== 1.0.0.0"},
		{"shorthand beta", "1.2.3b5", "== 1.2.3.0-beta5"},
		{"shorthand alpha", "1.2.3a1", "== 1.2.3.0-alpha1"},
		{"shorthand patch", "1.2.3p1234", "== 1.2.3.0-patch1234"},
		{"shorthand patch/2", "1.2.3pl1234", "== 1.2.3.0-patch1234"},
		{"accepts spaces", ">= 1.2.3", ">= 1.2.3.0-dev"},
		{"accepts spaces/2", "< 1.2.3", "< 1.2.3.0-dev"},
		{"accepts spaces/3", "> 1.2.3", "> 1.2.3.0"},
		{"accepts master", ">=dev-master", ">= 9999999-dev"},
		{"accepts master/2", "dev-master", "== 9999999-dev"},
		{"accepts arbitrary", "dev-feature-a", "== dev-feature-a"},
		{"regression #550", "dev-some-fix", "== dev-some-fix"},
		{"regression #935", "dev-CAPS", "== dev-CAPS"},
		{"ignores aliases", "dev-master as 1.0.0", "== 9999999-dev"},
		{"lesser than override", "<1.2.3.4-stable", "< 1.2.3.4"},
		{"great/eq than override", ">=1.2.3.4-stable", ">= 1.2.3.4"},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			constraint, err := NewConstraint(tc.constraint)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.parsed, constraint.String())
			}
		})
	}
}

func TestParseConstraintsWildcard(t *testing.T) {
	cases := []struct {
		name       string
		constraint string
	}{
		{"v2.*", "[>= 2.0.0.0-dev < 3.0.0.0-dev]"},
		{"2.*.*", "[>= 2.0.0.0-dev < 3.0.0.0-dev]"},
		{"20.*", "[>= 20.0.0.0-dev < 21.0.0.0-dev]"},
		{"20.*.*", "[>= 20.0.0.0-dev < 21.0.0.0-dev]"},
		{"2.0.*", "[>= 2.0.0.0-dev < 2.1.0.0-dev]"},
		{"2.x", "[>= 2.0.0.0-dev < 3.0.0.0-dev]"},
		{"2.x.x", "[>= 2.0.0.0-dev < 3.0.0.0-dev]"},
		{"2.2.x", "[>= 2.2.0.0-dev < 2.3.0.0-dev]"},
		{"2.10.X", "[>= 2.10.0.0-dev < 2.11.0.0-dev]"},
		{"2.1.3.*", "[>= 2.1.3.0-dev < 2.1.4.0-dev]"},
		{"0.*", "< 1.0.0.0-dev"},
		{"0.*.*", "< 1.0.0.0-dev"},
		{"0.x", "< 1.0.0.0-dev"},
		{"0.x.x", "< 1.0.0.0-dev"},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			constraint, err := NewConstraint(tc.name)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.constraint, constraint.String())
			}
		})
	}
}

func TestParseTildeWildcard(t *testing.T) {
	cases := []struct {
		name       string
		constraint string
	}{
		{"~v1", "[>= 1.0.0.0-dev < 2.0.0.0-dev]"},
		{"~1.0", "[>= 1.0.0.0-dev < 2.0.0.0-dev]"},
		{"~1.0.0", "[>= 1.0.0.0-dev < 1.1.0.0-dev]"},
		{"~1.2", "[>= 1.2.0.0-dev < 2.0.0.0-dev]"},
		{"~1.2.3", "[>= 1.2.3.0-dev < 1.3.0.0-dev]"},
		{"~1.2.3.4", "[>= 1.2.3.4-dev < 1.2.4.0-dev]"},
		{"~1.2-beta", "[>= 1.2.0.0-beta < 2.0.0.0-dev]"},
		{"~1.2-b2", "[>= 1.2.0.0-beta2 < 2.0.0.0-dev]"},
		{"~1.2-BETA2", "[>= 1.2.0.0-beta2 < 2.0.0.0-dev]"},
		{"~1.2.2-dev", "[>= 1.2.2.0-dev < 1.3.0.0-dev]"},
		{"~1.2.2-stable", "[>= 1.2.2.0 < 1.3.0.0-dev]"},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			constraint, err := NewConstraint(tc.name)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.constraint, constraint.String())
			}
		})
	}
}

func TestParseCaretWildcard(t *testing.T) {
	cases := []struct {
		name       string
		constraint string
	}{
		{"^v1", "[>= 1.0.0.0-dev < 2.0.0.0-dev]"},
		{"^0", "[>= 0.0.0.0-dev < 1.0.0.0-dev]"},
		{"^0.0", "[>= 0.0.0.0-dev < 0.1.0.0-dev]"},
		{"^1.2", "[>= 1.2.0.0-dev < 2.0.0.0-dev]"},
		{"^1.2.3-beta.2", "[>= 1.2.3.0-beta2 < 2.0.0.0-dev]"},
		{"^1.2.3.4", "[>= 1.2.3.4-dev < 2.0.0.0-dev]"},
		{"^1.2.3", "[>= 1.2.3.0-dev < 2.0.0.0-dev]"},
		{"^0.2.3", "[>= 0.2.3.0-dev < 0.3.0.0-dev]"},
		{"^0.2", "[>= 0.2.0.0-dev < 0.3.0.0-dev]"},
		{"^0.2.0", "[>= 0.2.0.0-dev < 0.3.0.0-dev]"},
		{"^0.0.3", "[>= 0.0.3.0-dev < 0.0.4.0-dev]"},
		{"^0.0.3-alpha", "[>= 0.0.3.0-alpha < 0.0.4.0-dev]"},
		{"^0.0.3-dev", "[>= 0.0.3.0-dev < 0.0.4.0-dev]"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			constraint, err := NewConstraint(tc.name)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.constraint, constraint.String())
			}
		})
	}
}

func TestParseHyphen(t *testing.T) {
	cases := []struct {
		name       string
		constraint string
	}{
		{"v1 - v2", "[>= 1.0.0.0-dev < 3.0.0.0-dev]"},
		{"1.2.3 - 2.3.4.5", "[>= 1.2.3.0-dev <= 2.3.4.5]"},
		{"1.2-beta - 2.3", "[>= 1.2.0.0-beta < 2.4.0.0-dev]"},
		{"1.2-beta - 2.3-dev", "[>= 1.2.0.0-beta <= 2.3.0.0-dev]"},
		{"1.2-RC - 2.3.1", "[>= 1.2.0.0-RC <= 2.3.1.0]"},
		{"1.2.3-alpha - 2.3-RC", "[>= 1.2.3.0-alpha <= 2.3.0.0-RC]"},
		{"1 - 2.0", "[>= 1.0.0.0-dev < 2.1.0.0-dev]"},
		{"1 - 2.1", "[>= 1.0.0.0-dev < 2.2.0.0-dev]"},
		{"1.2 - 2.1.0", "[>= 1.2.0.0-dev <= 2.1.0.0]"},
		{"1.3 - 2.1.3", "[>= 1.3.0.0-dev <= 2.1.3.0]"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			constraint, err := NewConstraint(tc.name)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.constraint, constraint.String())
			}
		})
	}
}

func TestParseConstraintsMulti(t *testing.T) {
	cases := []struct {
		constraint string
	}{
		{">2.0,<=3.0"},
		{">2.0 <=3.0"},
		{">2.0  <=3.0"},
		{">2.0, <=3.0"},
		{">2.0 ,<=3.0"},
		{">2.0 , <=3.0"},
		{">2.0   , <=3.0"},
		{"> 2.0   <=  3.0"},
		{"> 2.0  ,  <=  3.0"},
		{"  > 2.0  ,  <=  3.0 "},
	}

	for _, tc := range cases {
		t.Run(tc.constraint, func(t *testing.T) {
			constraint, err := NewConstraint(tc.constraint)
			if assert.NoError(t, err) {
				assert.Equal(t, "[> 2.0.0.0 <= 3.0.0.0]", constraint.String())
			}
		})
	}
}

func TestParseConstraintsMultiCollapsesContiguous(t *testing.T) {
	constraint, err := NewConstraint("^2.5 || ^3.0")
	if assert.NoError(t, err) {
		assert.Equal(t, "[>= 2.5.0.0-dev || < 4.0.0.0-dev]", constraint.String())
	}
}

func TestParseCaretConstraintsMultiDoesNotCollapseNonContiguousRange(t *testing.T) {
	constraint, err := NewConstraint("^0.2 || ^1.0")
	if assert.NoError(t, err) {
		assert.Equal(t, "[[>= 0.2.0.0-dev < 0.3.0.0-dev] || [>= 1.0.0.0-dev < 2.0.0.0-dev]]", constraint.String())
	}
}

func TestDoNotCollapseContiguousRangeIfOtherConstraintsAlsoApply(t *testing.T) {
	constraint, err := NewConstraint("~0.1 || ~1.0 !=1.0.1")
	if assert.NoError(t, err) {
		assert.Equal(t, "[[>= 0.1.0.0-dev < 1.0.0.0-dev] || [>= 1.0.0.0-dev < 2.0.0.0-dev != 1.0.1.0]]", constraint.String())
	}
}

func TestParseConstraintsMultiWithStabilitySuffix(t *testing.T) {
	constraint, err := NewConstraint(">=1.1.0-alpha4,<1.2.x-dev")
	if assert.NoError(t, err) {
		assert.Equal(t, "[>= 1.1.0.0-alpha4 < 1.2.9999999.9999999-dev]", constraint.String())
	}

	constraint, err = NewConstraint(">=1.1.0-alpha4,<1.2-beta2")
	if assert.NoError(t, err) {
		assert.Equal(t, "[>= 1.1.0.0-alpha4 < 1.2.0.0-beta2]", constraint.String())
	}
}

func TestParseConstraintsMultiWithStabilities(t *testing.T) {
	constraint, err := NewConstraint(">2.0@stable,<=3.0@dev")
	if assert.NoError(t, err) {
		assert.Equal(t, "[> 2.0.0.0 <= 3.0.0.0-dev]", constraint.String())
	}
}

func TestParseConstraintsMultiDisjunctiveHasPrioOverConjuctive(t *testing.T) {
	cases := []struct {
		constraint string
	}{
		{">2.0,<2.0.5 | >2.0.6"},
		{">2.0,<2.0.5 || >2.0.6"},
		{"> 2.0 , <2.0.5 | >  2.0.6"},
	}

	for _, tc := range cases {
		t.Run(tc.constraint, func(t *testing.T) {
			constraint, err := NewConstraint(tc.constraint)
			if assert.NoError(t, err) {
				assert.Equal(t, "[[> 2.0.0.0 < 2.0.5.0-dev] || > 2.0.6.0]", constraint.String())
			}
		})
	}
}

func TestParseConstraintsFails(t *testing.T) {
	cases := []struct {
		name       string
		constraint string
	}{
		{"empty", ""},
		{"invalid version", "1.0.0-meh"},
		{"operator abuse", ">2.0,,<=3.0"},
		{"operator abuse/2", ">2.0 ,, <=3.0"},
		{"operator abuse/3", ">2.0 ||| <=3.0"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConstraint(tc.name)
			assert.Error(t, err)
		})
	}
}

func TestConstraintMatches(t *testing.T) {

	cases := []struct {
		constraint string
		version    string
		matches    bool
	}{
		{"*", "1.4.0", true},
		{"^1.0", "1.4.0", true},
		{"^1.0 || ^2.0", "1.4.0", true},
		{"^1.0 || ^2.0", "2.9.0", true},
		{"1.x", "1.1", true},
		{">= 1.2.3 < 1.2.6 != 1.2.5", "1.2.5", false},
		{">= 1.2.3 < 1.2.6 != 1.2.5", "1.2.4", true},
		{">=1.1.0-alpha4,<1.2.x-dev", "1.1.0-alpha3", false},
		{">=1.1.0-alpha4,<1.2.x-dev", "1.1.0-beta", true},
		{"^0.4", "0.5", false},
		{"^0.4", "0.4.1", true},

		// Satifies test suit from Composer

		{"1.0.0 - 2.0.0", "1.2.3", true},
		{"^1.2.3+build", "1.2.3", true},
		{"^1.2.3+build", "1.3.0", true},
		{"1.2.3+asdf - 2.4.3+asdf", "2.4.3-alpha", true},
		{">1.2", "1.3.0-beta", true},
		{"<=1.2.3", "1.2.3-beta", true},
		{"^1.2.3", "1.2.3-beta", true},
		{"1.2.3+asdf - 2.4.3+asdf", "1.2.3", true},
		{"1.0.0", "1.0.0", true},
		{"*", "1.2.3", true},
		{"*", "v1.2.3", true},
		{">=1.0.0", "1.0.0", true},
		{">=1.0.0", "1.0.1", true},
		{">=1.0.0", "1.1.0", true},
		{">1.0.0", "1.0.1", true},
		{">1.0.0", "1.1.0", true},
		{"<=2.0.0", "2.0.0", true},
		{"<=2.0.0", "1.9999.9999", true},
		{"<=2.0.0", "0.2.9", true},
		{"<2.0.0", "1.9999.9999", true},
		{"<2.0.0", "0.2.9", true},
		{">= 1.0.0", "1.0.0", true},
		{">=  1.0.0", "1.0.1", true},
		{">=   1.0.0", "1.1.0", true},
		{"> 1.0.0", "1.0.1", true},
		{">  1.0.0", "1.1.0", true},
		{"<=   2.0.0", "2.0.0", true},
		{"<= 2.0.0", "1.9999.9999", true},
		{"<=  2.0.0", "0.2.9", true},
		{"<    2.0.0", "1.9999.9999", true},
		{"<\t2.0.0", "0.2.9", true},
		{">=0.1.97", "v0.1.97", true},
		{">=0.1.97", "0.1.97", true},
		{"0.1.20 || 1.2.4", "1.2.4", true},
		{">=0.2.3 || <0.0.1", "0.0.0", true},
		{">=0.2.3 || <0.0.1", "0.2.3", true},
		{">=0.2.3 || <0.0.1", "0.2.4", true},
		{"2.x.x", "2.1.3", true},
		{"1.2.x", "1.2.3", true},
		{"1.2.x || 2.x", "2.1.3", true},
		{"1.2.x || 2.x", "1.2.3", true},
		{"x", "1.2.3", true},
		{"2.*.*", "2.1.3", true},
		{"1.2.*", "1.2.3", true},
		{"1.2.* || 2.*", "2.1.3", true},
		{"1.2.* || 2.*", "1.2.3", true},
		{"*", "1.2.3", true},
		{"~2.4", "2.9.0", true}, // >=2.4.0 <3.0.0
		{"~2.4", "2.4.5", true},
		{"~1", "1.2.3", true},   //  >=1.0.0 <2.0.0
		{"~1.0", "1.4.7", true}, // >=1.0.0 <2.0.0
		{">=1'),", "1.0.0", true},
		{">= 1", "1.0.0", true},
		{">1.2", "1.2.8", true}, // >1.2.0
		{"<1.2", "1.1.1", true},
		{"< 1.2", "1.1.1", true},
		{"~1.2.1 >=1.2.3", "1.2.3", true},
		{"~1.2.1 =1.2.3", "1.2.3", true},
		{"~1.2.1 1.2.3", "1.2.3", true},
		{"~1.2.1 >=1.2.3 1.2.3", "1.2.3", true},
		{"~1.2.1 1.2.3 >=1.2.3", "1.2.3", true},
		{"~1.2.1 1.2.3", "1.2.3", true},
		{">=1.2.1 1.2.3", "1.2.3", true},
		{"1.2.3 >=1.2.1", "1.2.3", true},
		{">=1.2.3 >=1.2.1", "1.2.3", true},
		{">=1.2.1 >=1.2.3", "1.2.3", true},
		{">=1.2", "1.2.8", true},
		{"^1.2.3", "1.8.1", true},
		{"^0.1.2", "0.1.2", true},
		{"^0.1", "0.1.2", true},
		{"^1.2", "1.4.2", true},
		{"^1.2 ^1", "1.4.2", true},
		{"^0.0.1-alpha", "0.0.1-beta", true},

		{"1.0.0 - 2.0.0", "2.2.3", false},
		{"^1.2.3+build", "2.0.0", false},
		{"^1.2.3+build", "1.2.0", false},
		{"1", "1.0.0beta", false},
		{"<1", "1.0.0beta", false},
		{"< 1", "1.0.0beta", false},
		{"1.0.0", "1.0.1", false},
		{">=1.0.0", "0.0.0", false},
		{">=1.0.0", "0.0.1", false},
		{">=1.0.0", "0.1.0", false},
		{">1.0.0", "0.0.1", false},
		{">1.0.0", "0.1.0", false},
		{"<=2.0.0", "3.0.0", false},
		{"<=2.0.0", "2.9999.9999", false},
		{"<=2.0.0", "2.2.9", false},
		{"<2.0.0", "2.9999.9999", false},
		{"<2.0.0", "2.2.9", false},
		{">=0.1.97", "v0.1.93", false},
		{">=0.1.97", "0.1.93", false},
		{"0.1.20 || 1.2.4", "1.2.3", false},
		{">=0.2.3 || <0.0.1", "0.0.3", false},
		{">=0.2.3 || <0.0.1", "0.2.2", false},
		{"2.x.x", "1.1.3", false},
		{"2.x.x", "3.1.3", false},
		{"1.2.x", "1.3.3", false},
		{"1.2.x || 2.x", "3.1.3", false},
		{"1.2.x || 2.x", "1.1.3", false},
		{"2.*.*", "1.1.3", false},
		{"2.*.*", "3.1.3", false},
		{"1.2.*", "1.3.3", false},
		{"1.2.* || 2.*", "3.1.3", false},
		{"1.2.* || 2.*", "1.1.3", false},
		{"2", "1.1.2", false},
		{"2.3", "2.4.1", false},
		{"~2.4", "3.0.0", false}, // >=2.4.0 <3
		{"~2.4", "2.3.9", false},
		{"~1", "0.2.3", false}, // >=1.0.0 <2
		{"<1", "1.0.0", false},
		{">=1.2", "1.1.1", false},
		{"1", "2.0.0beta", false},
		{"~v0.5.4-beta", "0.5.4-alpha", false},
		{"<1.2.3", "1.2.3-beta", false},
		{"^1.2.3", "2.0.0-alpha", false},
		{"^1.2.3", "1.2.2", false},
		{"^1.2", "1.1.9", false},
	}

	for _, tc := range cases {
		t.Run(tc.constraint, func(t *testing.T) {
			c, err := NewConstraint(tc.constraint)
			if assert.NoError(t, err) {
				version, err := NewVersion(tc.version)
				if assert.NoError(t, err) {
					assert.Equal(t, tc.matches, c.Matches(version))
				}
			}
		})
	}

}

func BenchmarkParseConstraintsSimple(b *testing.B) {
	cases := []struct {
		name       string
		constraint string
	}{
		{"match any", "*"},
		{"match any/2", "*.*"},
		{"match any/2v", "v*.*"},
		{"match any/3", "*.x.*"},
		{"match any/4", "x.X.x.*"},
		{"not equal", "<>1.0.0"},
		{"not equal/2", "!=1.0.0"},
		{"greater than", ">1.0.0"},
		{"lesser than", "<1.2.3.4"},
		{"less/eq than", "<=1.2.3"},
		{"great/eq than", ">=1.2.3"},
		{"equals", "=1.2.3"},
		{"double equals", "==1.2.3"},
		{"no op means eq", "1.2.3"},
		{"completes version", "=1.0"},
		{"shorthand beta", "1.2.3b5"},
		{"shorthand alpha", "1.2.3a1"},
		{"shorthand patch", "1.2.3p1234"},
		{"shorthand patch/2", "1.2.3pl1234"},
		{"accepts spaces", ">= 1.2.3"},
		{"accepts spaces/2", "< 1.2.3"},
		{"accepts spaces/3", "> 1.2.3"},
		{"accepts master", ">=dev-master"},
		{"accepts master/2", "dev-master"},
		{"accepts arbitrary", "dev-feature-a"},
		{"regression #550", "dev-some-fix"},
		{"regression #935", "dev-CAPS"},
		{"ignores aliases", "dev-master as 1.0.0"},
		{"lesser than override", "<1.2.3.4-stable"},
		{"great/eq than override", ">=1.2.3.4-stable"},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, _ = NewConstraint(tc.constraint)
			}
		})
	}
}
