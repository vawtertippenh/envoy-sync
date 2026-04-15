// Package inject provides functionality to inject key-value pairs into an
// existing env map, with support for overwrite control and prefix namespacing.
package inject

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls injection behaviour.
type Options struct {
	// Overwrite allows existing keys to be replaced. Default: false.
	Overwrite bool
	// Prefix is prepended to every injected key before insertion.
	Prefix string
}

// Result holds the outcome of an Inject call.
type Result struct {
	Injected  []string
	Skipped   []string
	Overwrite []string
}

// Summary returns a human-readable description of the result.
func (r Result) Summary() string {
	var sb strings.Builder
	if len(r.Injected) > 0 {
		fmt.Fprintf(&sb, "injected: %s\n", strings.Join(r.Injected, ", "))
	}
	if len(r.Overwrite) > 0 {
		fmt.Fprintf(&sb, "overwritten: %s\n", strings.Join(r.Overwrite, ", "))
	}
	if len(r.Skipped) > 0 {
		fmt.Fprintf(&sb, "skipped (already exists): %s\n", strings.Join(r.Skipped, ", "))
	}
	if sb.Len() == 0 {
		return "nothing to inject"
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Inject merges pairs into base according to opts. base is not mutated; a new
// map is returned alongside a Result describing what changed.
func Inject(base, pairs map[string]string, opts Options) (map[string]string, Result) {
	out := copyMap(base)
	var res Result

	for _, k := range sortedKeys(pairs) {
		v := pairs[k]
		destKey := opts.Prefix + k
		if _, exists := out[destKey]; exists {
			if opts.Overwrite {
				out[destKey] = v
				res.Overwrite = append(res.Overwrite, destKey)
			} else {
				res.Skipped = append(res.Skipped, destKey)
			}
		} else {
			out[destKey] = v
			res.Injected = append(res.Injected, destKey)
		}
	}
	return out, res
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
