package semver

import (
	"bytes"
	"fmt"
	"github.com/spf13/cast"
	"strings"
)

type semver struct {
	Major, Minor, Patch, Extra int
	PreRelease                 string
	State                      string
	Stability                  string
	Metadata                   string
	Original                   string
	Parsed                     string
}

func (v *semver) major() int {
	return v.Major
}

func (v *semver) minor() int {
	return v.Minor
}

func (v *semver) patch() int {
	return v.Patch
}

func (v *semver) pre() float32 {
	return cast.ToFloat32(strings.Replace(v.PreRelease, v.Stability, "", -1))
}

func (v *semver) stability() string {
	return v.Stability
}

func (v *semver) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, cast.ToString(v.Major))

	if 9999999 != v.Major {
		fmt.Fprintf(&buf, ".%d", v.Minor)

		fmt.Fprintf(&buf, ".%d", v.Patch)

		fmt.Fprintf(&buf, ".%d", v.Extra)
	}

	if v.PreRelease != "" {
		fmt.Fprintf(&buf, "-%s", v.PreRelease)
	}

	if v.State != "" {
		fmt.Fprintf(&buf, "-%s", v.State)
	}

	return buf.String()
}
