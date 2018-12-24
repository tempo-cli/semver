package semver

const (
	LessThan = iota - 1
	Equal
	GreaterThan
)

func (a *Version) LessThan(b *Version) bool {
	return LessThan == compare(a, b)
}

func (a *Version) GreaterThan(b *Version) bool {
	return GreaterThan == compare(a, b)
}

func (a *Version) LessThanOrEqual(b *Version) bool {
	comparison := compare(a, b)
	return LessThan == comparison || Equal == comparison
}

func (a *Version) GreaterThanOrEqual(b *Version) bool {
	comparison := compare(a, b)
	return GreaterThan == comparison || Equal == comparison
}

func (a *Version) Equal(b *Version) bool {
	return Equal == compare(a, b)
}

func (a *Version) NotEqual(b *Version) bool {
	return Equal != compare(a, b)
}

func (a *Version) Compare(b *Version, operator string) bool {

	switch operator {
	case ">":
		return GreaterThan == compare(a, b)
	case ">=":
		comparison := compare(a, b)
		return GreaterThan == comparison || Equal == comparison
	case "<":
		return LessThan == compare(a, b)
	case "<=":
		comparison := compare(a, b)
		return LessThan == comparison || Equal == comparison
	case "==", "=":
		return Equal == compare(a, b)
	case "!=", "<>":
		return Equal != compare(a, b)
	}

	return false
}

func compare(a *Version, b *Version) int {
	if d := comparePart(a.major(), b.major()); d != Equal {
		return d
	}

	if d := comparePart(a.minor(), b.minor()); d != Equal {
		return d
	}

	if d := comparePart(a.patch(), b.patch()); d != Equal {
		return d
	}

	if d := compareStability(a.stability(), b.stability()); d != Equal {
		return d
	}

	aPre := a.pre()
	bPre := b.pre()

	if aPre > bPre {
		return GreaterThan
	}

	if aPre < bPre {
		return LessThan
	}

	return Equal
}

func compareStability(a string, b string) int {
	matches := map[string]int{"dev": 1, "alpha": 2, "beta": 3, "RC": 4, "stable": 5}

	if matches[a] > matches[b] {
		return GreaterThan
	}

	if matches[a] < matches[b] {
		return LessThan
	}

	return Equal
}

func comparePart(a int, b int) int {
	if a > b {
		return GreaterThan
	}

	if b > a {
		return LessThan
	}

	return Equal
}

