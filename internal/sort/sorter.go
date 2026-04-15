// Package sort provides utilities for sorting .env file keys
// by various strategies: alphabetical, length, or grouped by prefix.
package sort

import (
	"sort"
	"strings"
)

// Strategy defines how keys should be sorted.
type Strategy string

const (
	Alpha  Strategy = "alpha"  // alphabetical order
	Length Strategy = "length" // shortest key first
	Prefix Strategy = "prefix" // grouped by prefix (e.g. DB_, APP_)
)

// Options controls Sort behaviour.
type Options struct {
	Strategy  Strategy
	Descending bool
}

// Sort returns a new map identical to env but provides a deterministic
// key order via the returned slice. The map itself is unchanged.
func Sort(env map[string]string, opts Options) (map[string]string, []string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch opts.Strategy {
	case Length:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) != len(keys[j]) {
				return len(keys[i]) < len(keys[j])
			}
			return keys[i] < keys[j]
		})
	case Prefix:
		sort.Slice(keys, func(i, j int) bool {
			pi := extractPrefix(keys[i])
			pj := extractPrefix(keys[j])
			if pi != pj {
				return pi < pj
			}
			return keys[i] < keys[j]
		})
	default: // Alpha
		sort.Strings(keys)
	}

	if opts.Descending {
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	return out, keys
}

// extractPrefix returns the portion of a key before the first underscore,
// or the full key if no underscore is present.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
