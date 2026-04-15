// Package split provides functionality to split a flat env map into
// multiple named groups based on key prefix patterns.
package split

import (
	"fmt"
	"sort"
	"strings"
)

// Group represents a named subset of env vars.
type Group struct {
	Name string
	Env  map[string]string
}

// Options controls splitting behaviour.
type Options struct {
	// StripPrefix removes the matched prefix from keys in the output group.
	StripPrefix bool
	// Prefixes maps group names to their key prefixes (e.g. {"app": "APP_"}).
	Prefixes map[string]string
	// Remainder is the name for keys that matched no prefix.
	// If empty, unmatched keys are discarded.
	Remainder string
}

// Split partitions env into groups according to opts.Prefixes.
// Keys are matched case-sensitively. Each key is assigned to the first
// matching group (in deterministic alphabetical order of group names).
func Split(env map[string]string, opts Options) ([]Group, error) {
	if len(opts.Prefixes) == 0 {
		return nil, fmt.Errorf("split: at least one prefix mapping is required")
	}

	// Stable iteration order over group names.
	groupNames := make([]string, 0, len(opts.Prefixes))
	for name := range opts.Prefixes {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	buckets := make(map[string]map[string]string, len(groupNames))
	for _, name := range groupNames {
		buckets[name] = make(map[string]string)
	}

	var remainderBucket map[string]string
	if opts.Remainder != "" {
		remainderBucket = make(map[string]string)
	}

	for k, v := range env {
		matched := false
		for _, name := range groupNames {
			prefix := opts.Prefixes[name]
			if strings.HasPrefix(k, prefix) {
				outKey := k
				if opts.StripPrefix {
					outKey = strings.TrimPrefix(k, prefix)
				}
				buckets[name][outKey] = v
				matched = true
				break
			}
		}
		if !matched && remainderBucket != nil {
			remainderBucket[k] = v
		}
	}

	result := make([]Group, 0, len(groupNames)+1)
	for _, name := range groupNames {
		result = append(result, Group{Name: name, Env: buckets[name]})
	}
	if opts.Remainder != "" {
		result = append(result, Group{Name: opts.Remainder, Env: remainderBucket})
	}
	return result, nil
}
