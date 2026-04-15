// Package resolve provides functionality to resolve an env map against
// a set of override sources, applying them in priority order.
package resolve

import "sort"

// Source represents a named set of key-value overrides.
type Source struct {
	Name   string
	Values map[string]string
}

// Result holds the resolved environment and metadata about applied overrides.
type Result struct {
	Env       map[string]string
	Overrides map[string]string // key -> source name that provided the final value
}

// Resolve merges base with the provided sources in order (last wins).
// Keys present in later sources override earlier ones and the base.
func Resolve(base map[string]string, sources []Source) Result {
	env := copyMap(base)
	overrides := make(map[string]string)

	for _, src := range sources {
		for k, v := range src.Values {
			env[k] = v
			overrides[k] = src.Name
		}
	}

	// Remove override tracking for keys whose final value matches base
	for k, srcName := range overrides {
		if bv, ok := base[k]; ok && bv == env[k] {
			delete(overrides, k)
			_ = srcName
		}
	}

	return Result{Env: env, Overrides: overrides}
}

// Summary returns a human-readable list of keys that were overridden and by which source.
func Summary(r Result) []string {
	keys := make([]string, 0, len(r.Overrides))
	for k := range r.Overrides {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, k+" <- "+r.Overrides[k])
	}
	return lines
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
