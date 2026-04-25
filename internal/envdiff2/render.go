package envdiff2

import (
	"fmt"
	"strings"

	"github.com/user/envoy-sync/internal/mask"
)

// RenderOptions controls how the diff is printed.
type RenderOptions struct {
	MaskSensitive   bool
	ShowUnchanged   bool
	Colorize        bool
}

// Render returns a human-readable diff string.
func Render(r Result, opts RenderOptions) string {
	var sb strings.Builder
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			v := maskIf(c.Key, c.NewVal, opts.MaskSensitive)
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", c.Key, v))
		case Removed:
			v := maskIf(c.Key, c.OldVal, opts.MaskSensitive)
			sb.WriteString(fmt.Sprintf("- %s=%s\n", c.Key, v))
		case Modified:
			ov := maskIf(c.Key, c.OldVal, opts.MaskSensitive)
			nv := maskIf(c.Key, c.NewVal, opts.MaskSensitive)
			sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", c.Key, ov, nv))
		case Unchanged:
			if opts.ShowUnchanged {
				v := maskIf(c.Key, c.OldVal, opts.MaskSensitive)
				sb.WriteString(fmt.Sprintf("  %s=%s\n", c.Key, v))
			}
		}
	}
	return sb.String()
}

// Summary returns a one-line summary of the diff.
func Summary(r Result) string {
	var added, removed, modified int
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return fmt.Sprintf("added=%d removed=%d modified=%d", added, removed, modified)
}

func maskIf(key, val string, sensitive bool) string {
	if sensitive && mask.IsSensitive(key, nil) {
		return mask.MaskValue(val)
	}
	return val
}
