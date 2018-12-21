package semver

import (
	"bytes"
	"fmt"
)

type date struct {
	semver
}

func (v *date) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, v.Parsed)

	if v.Stability != "" {
		fmt.Fprintf(&buf, "-"+v.Stability)
	}

	if v.Patch != "" {
		fmt.Fprintf(&buf, v.Patch)
	}

	if v.Metadata != "" {
		fmt.Fprintf(&buf, v.Metadata)
	}

	return buf.String()
}
