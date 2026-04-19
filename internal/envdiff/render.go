package envdiff

import "fmt"

// RenderOptions controls how diffs are rendered.
type RenderOptions struct {
	ShowUnchanged bool
	MaskValues    bool
}

const masked = "***"

// Render returns a human-readable slice of lines for the diff result.
func Render(r Result, opts RenderOptions) []string {
	var lines []string
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			v := maskIf(c.NewValue, opts.MaskValues)
			lines = append(lines, fmt.Sprintf("+ %s=%s", c.Key, v))
		case Removed:
			v := maskIf(c.OldValue, opts.MaskValues)
			lines = append(lines, fmt.Sprintf("- %s=%s", c.Key, v))
		case Modified:
			old := maskIf(c.OldValue, opts.MaskValues)
			new := maskIf(c.NewValue, opts.MaskValues)
			lines = append(lines, fmt.Sprintf("~ %s: %s -> %s", c.Key, old, new))
		case Unchanged:
			if opts.ShowUnchanged {
				v := maskIf(c.OldValue, opts.MaskValues)
				lines = append(lines, fmt.Sprintf("  %s=%s", c.Key, v))
			}
		}
	}
	return lines
}

// Summary returns a one-line summary string.
func Summary(r Result) string {
	return fmt.Sprintf("+%d -%d ~%d",
		len(r.Added()), len(r.Removed()), len(r.Modified()))
}

func maskIf(v string, mask bool) string {
	if mask {
		return masked
	}
	return v
}
