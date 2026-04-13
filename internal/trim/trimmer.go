// Package trim provides utilities for removing unused or duplicate keys
// from .env files, helping keep environments clean and minimal.
package trim

import "sort"

// Result holds the output of a Trim operation.
type Result struct {
	Kept    map[string]string
	Removed []string
}

// Options controls Trim behaviour.
type Options struct {
	// AllowList, when non-empty, keeps only keys present in this set.
	AllowList []string
	// DenyList removes any keys present in this set.
	DenyList []string
	// RemoveEmpty removes keys whose value is an empty string.
	RemoveEmpty bool
}

// Trim filters env according to opts and returns a Result describing what was
// kept and what was removed.
func Trim(env map[string]string, opts Options) Result {
	allowSet := toSet(opts.AllowList)
	denySet := toSet(opts.DenyList)

	kept := make(map[string]string)
	var removed []string

	for _, k := range sortedKeys(env) {
		v := env[k]

		if len(allowSet) > 0 && !allowSet[k] {
			removed = append(removed, k)
			continue
		}
		if denySet[k] {
			removed = append(removed, k)
			continue
		}
		if opts.RemoveEmpty && v == "" {
			removed = append(removed, k)
			continue
		}

		kept[k] = v
	}

	return Result{Kept: kept, Removed: removed}
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
