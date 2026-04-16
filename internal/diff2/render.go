package diff2

import (
	"fmt"
	"strings"
)

// RenderOptions controls how the diff is rendered.
type RenderOptions struct {
	ShowUnchanged bool
	MaskValues   bool
}

const maskPlaceholder = "***"

func maskIf(val string, mask bool) string {
	if mask {
		return maskPlaceholder
	}
	return val
}

// Render formats the diff result as a human-readable string.
func Render(r Result, opts RenderOptions) string {
	var sb strings.Builder
	for _, e := range r.Entries {
		switch e.Change {
		case Added:
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, maskIf(e.NewVal, opts.MaskValues))
		case Removed:
			fmt.Fprintf(&sb, "- %s=%s\n", e.Key, maskIf(e.OldVal, opts.MaskValues))
		case Modified:
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", e.Key, maskIf(e.OldVal, opts.MaskValues), maskIf(e.NewVal, opts.MaskValues))
		case Unchanged:
			if opts.ShowUnchanged {
				fmt.Fprintf(&sb, "  %s=%s\n", e.Key, maskIf(e.NewVal, opts.MaskValues))
			}
		}
	}
	return sb.String()
}

// Summary returns a short statistics line.
func Summary(r Result) string {
	return fmt.Sprintf("added=%d removed=%d modified=%d unchanged=%d",
		len(r.Added()), len(r.Removed()), len(r.Modified()), len(r.filter(Unchanged)))
}
