// Package envsorter provides functionality to sort environment variable maps
// by key name, value, or length, with optional grouping by prefix.
package envsorter

import (
	"sort"
	"strings"
)

// Strategy defines how keys should be sorted.
type Strategy string

const (
	StrategyAlpha   Strategy = "alpha"
	StrategyValue   Strategy = "value"
	StrategyLength  Strategy = "length"
	StrategyPrefix  Strategy = "prefix"
)

// Options controls sorting behaviour.
type Options struct {
	Strategy   Strategy
	Descending bool
	// PrefixSep is the separator used to detect prefix groups (e.g. "_").
	PrefixSep string
}

// Result holds the sorted key order and the original map.
type Result struct {
	Env  map[string]string
	Keys []string
}

// Sort returns a Result with keys ordered according to opts.
func Sort(env map[string]string, opts Options) Result {
	if opts.PrefixSep == "" {
		opts.PrefixSep = "_"
	}
	if opts.Strategy == "" {
		opts.Strategy = StrategyAlpha
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch opts.Strategy {
	case StrategyValue:
		sort.Slice(keys, func(i, j int) bool {
			return less(env[keys[i]], env[keys[j]], opts.Descending)
		})
	case StrategyLength:
		sort.Slice(keys, func(i, j int) bool {
			li, lj := len(keys[i]), len(keys[j])
			if li == lj {
				return less(keys[i], keys[j], opts.Descending)
			}
			if opts.Descending {
				return li > lj
			}
			return li < lj
		})
	case StrategyPrefix:
		sort.Slice(keys, func(i, j int) bool {
			pi := extractPrefix(keys[i], opts.PrefixSep)
			pj := extractPrefix(keys[j], opts.PrefixSep)
			if pi != pj {
				return less(pi, pj, opts.Descending)
			}
			return less(keys[i], keys[j], opts.Descending)
		})
	default: // StrategyAlpha
		sort.Slice(keys, func(i, j int) bool {
			return less(keys[i], keys[j], opts.Descending)
		})
	}

	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return Result{Env: copy, Keys: keys}
}

func less(a, b string, desc bool) bool {
	if desc {
		return strings.ToLower(a) > strings.ToLower(b)
	}
	return strings.ToLower(a) < strings.ToLower(b)
}

func extractPrefix(key, sep string) string {
	if idx := strings.Index(key, sep); idx > 0 {
		return key[:idx]
	}
	return key
}
