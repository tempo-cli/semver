package semver

import (
	"bytes"
	"fmt"
	"github.com/spf13/cast"
	"strings"
)

type Version struct {
	Major, Minor, Patch, Extra int
	PreRelease                 string
	State                      string
	Stability                  string
	Metadata                   string
	Original                   string
	Parsed                     string
	isDate                     bool
	isBranch                   bool
}

func (v *Version) major() int {
	return v.Major
}

func (v *Version) minor() int {
	return v.Minor
}

func (v *Version) patch() int {
	return v.Patch
}

func (v *Version) pre() float32 {
	return cast.ToFloat32(strings.Replace(v.PreRelease, v.Stability, "", -1))
}

func (v *Version) stability() string {
	return v.Stability
}

func (v *Version) String() string {
	var buf bytes.Buffer

	if v.isBranch || v.isDate {
		v.branchString(&buf)
	} else {
		_, _ = fmt.Fprintf(&buf, cast.ToString(v.Major))

		if 9999999 != v.Major {
			fmt.Fprintf(&buf, ".%d", v.Minor)

			fmt.Fprintf(&buf, ".%d", v.Patch)

			fmt.Fprintf(&buf, ".%d", v.Extra)
		}

		if v.Stability != "" {
			fmt.Fprintf(&buf, "-%s", v.Stability)
		}

		if v.PreRelease != "" {
			fmt.Fprintf(&buf, "%s", v.PreRelease)
		}

		if v.State != "" {
			fmt.Fprintf(&buf, "-%s", v.State)
		}
	}

	return buf.String()
}

func (v *Version) branchString(buf *bytes.Buffer) {

	fmt.Fprintf(buf, v.Parsed)

	if v.Stability != "" {
		fmt.Fprintf(buf, "-"+v.Stability)
	}

	if 0 != v.Patch {
		fmt.Fprintf(buf, cast.ToString(v.Patch))
	}

	if v.Metadata != "" {
		fmt.Fprintf(buf, v.Metadata)
	}
}
