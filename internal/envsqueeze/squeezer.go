// Package envсqueeze removes redundant or duplicate entries from an env map,
// collapsing keys that share the same value into a canonical form.
package envsqueeze

import "sort"

// Options controls the behaviour of Squeeze.
type Options struct {
	// DedupeValues collapses keys with identical values, keeping only the first
	// (alphabetically) key per unique value.
	DedupeValues bool
	// RemoveEmpty drops keys whose value is the empty string.
	RemoveEmpty bool
	// RemovePlaceholders drops keys whose value matches a placeholder pattern
	// such as "CHANGE_ME", "TODO", or "<placeholder>".
	RemovePlaceholders bool
}

// Result holds the squeezed env map and metadata about what was removed.
type Result struct {
	Env     map[string]string
	Dropped []string // keys that were removed
}

var placeholders = []string{"CHANGE_ME", "TODO", "<placeholder>", "PLACEHOLDER", "YOUR_VALUE_HERE"}

// Squeeze returns a new env map with redundant entries removed according to opts.
func Squeeze(env map[string]string, opts Options) Result {
	out := copyMap(env)
	var dropped []string

	if opts.RemoveEmpty {
		for k, v := range out {
			if v == "" {
				delete(out, k)
				dropped = append(dropped, k)
			}
		}
	}

	if opts.RemovePlaceholders {
		for k, v := range out {
			if isPlaceholder(v) {
				delete(out, k)
				dropped = append(dropped, k)
			}
		}
	}

	if opts.DedupeValues {
		seen := map[string]string{} // value -> first key
		for _, k := range sortedKeys(out) {
			v := out[k]
			if prev, exists := seen[v]; exists {
				_ = prev
				delete(out, k)
				dropped = append(dropped, k)
			} else {
				seen[v] = k
			}
		}
	}

	sort.Strings(dropped)
	return Result{Env: out, Dropped: dropped}
}

func isPlaceholder(v string) bool {
	for _, p := range placeholders {
		if v == p {
			return true
		}
	}
	return false
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
