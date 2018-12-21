package semver

import (
	"bytes"
	"fmt"
	"github.com/spf13/cast"
)

type branch struct {
	semver
}

func (v *branch) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, v.Parsed)

	if v.Stability != "" {
		fmt.Fprintf(&buf, "-"+v.Stability)
	}

	if 0 != v.Patch {
		fmt.Fprintf(&buf, cast.ToString(v.Patch))
	}

	if v.Metadata != "" {
		fmt.Fprintf(&buf, v.Metadata)
	}

	return buf.String()
}
