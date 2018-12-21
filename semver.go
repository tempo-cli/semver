package semver

import (
	"bytes"
	"fmt"
	"github.com/spf13/cast"
)

type semver struct {
	Major, Minor, Patch, Extra string
	PreRelease                 string
	State                      string
	Stability                  string
	Metadata                   string
	Original                   string
	Parsed                     string
}

func (v *semver) major() int {
	return cast.ToInt(v.Major)
}

func (v *semver) minor() int {
	return cast.ToInt(v.Major)
}

func (v *semver) patch() int {
	return cast.ToInt(v.Major)
}

func (v *semver) pre() int {
	return cast.ToInt(v.Major)
}

func (v *semver) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, v.Major)

	if "" != v.Minor {
		fmt.Fprintf(&buf, ".%s", v.Minor)
	}

	if "" != v.Patch {
		fmt.Fprintf(&buf, ".%s", v.Patch)
	}

	if "" != v.Extra {
		fmt.Fprintf(&buf, ".%s", v.Extra)
	}

	if v.PreRelease != "" {
		fmt.Fprintf(&buf, "-%s", v.PreRelease)
	}

	if v.State != "" {
		fmt.Fprintf(&buf, "-%s", v.State)
	}

	return buf.String()
}