// Package wrap provides functionality to wrap long env values at a specified
// column width, producing multi-line shell-compatible output.
package wrap

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls wrapping behaviour.
type Options struct {
	// Width is the maximum line length before wrapping. Defaults to 80.
	Width int
	// Indent is the string prepended to continuation lines. Defaults to two spaces.
	Indent string
	// ContinuationChar is the character appended to indicate line continuation.
	// Defaults to backslash.
	ContinuationChar string
}

func defaults(o Options) Options {
	if o.Width <= 0 {
		o.Width = 80
	}
	if o.Indent == "" {
		o.Indent = "  "
	}
	if o.ContinuationChar == "" {
		o.ContinuationChar = "\\"
	}
	return o
}

// Wrap takes a map of env vars and returns a slice of lines where any
// KEY=VALUE pair whose total length exceeds opts.Width is split across
// multiple lines using shell continuation syntax.
func Wrap(env map[string]string, opts Options) []string {
	opts = defaults(opts)

	keys := sortedKeys(env)
	var lines []string

	for _, k := range keys {
		line := fmt.Sprintf("%s=%s", k, env[k])
		if len(line) <= opts.Width {
			lines = append(lines, line)
			continue
		}
		lines = append(lines, wrapLine(k, env[k], opts)...)
	}
	return lines
}

// wrapLine splits a single KEY=VALUE into continuation lines.
func wrapLine(key, value string, opts Options) []string {
	prefix := key + "="
	chunk := opts.Width - len(prefix) - len(opts.ContinuationChar)
	if chunk <= 0 {
		// Key itself is too long; emit as-is.
		return []string{prefix + value}
	}

	var lines []string
	remaining := value
	first := true

	for len(remaining) > 0 {
		var linePrefix string
		var available int
		if first {
			linePrefix = prefix
			available = opts.Width - len(prefix) - len(opts.ContinuationChar)
			first = false
		} else {
			linePrefix = opts.Indent
			available = opts.Width - len(opts.Indent) - len(opts.ContinuationChar)
		}
		if available <= 0 {
			available = 1
		}

		if len(remaining) <= available {
			lines = append(lines, linePrefix+remaining)
			break
		}
		lines = append(lines, linePrefix+remaining[:available]+opts.ContinuationChar)
		remaining = remaining[available:]
	}
	return lines
}

// Render joins wrapped lines into a single string.
func Render(lines []string) string {
	return strings.Join(lines, "\n")
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
