package semver

type Version interface {
	String() string

	major() int
	minor() int
	patch() int
	pre() float32
	stability() string

	GreaterThan(b Version) bool
	GreaterThanOrEqual(b Version) bool
	LessThan(b Version) bool
	LessThanOrEqual(b Version) bool
	Equal(b Version) bool
	NotEqual(b Version) bool
	Compare(b Version, operator string) bool
}