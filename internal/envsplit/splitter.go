package envsplit

import (
	"fmt"
	"sort"
	"strings"
)

// Part holds a named subset of env vars.
type Part struct {
	Name string
	Env  map[string]string
}

// Options controls how splitting behaves.
type Options struct {
	Prefixes      []string
	StripPrefix   bool
	KeepRemainder bool
}

// Split partitions env into named parts based on key prefixes.
func Split(env map[string]string, opts Options) ([]Part, error) {
	if len(opts.Prefixes) == 0 {
		return nil, fmt.Errorf("at least one prefix is required")
	}

	groups := make(map[string]map[string]string, len(opts.Prefixes))
	for _, p := range opts.Prefixes {
		groups[p] = make(map[string]string)
	}
	remainder := make(map[string]string)

	for k, v := range env {
		matched := false
		for _, p := range opts.Prefixes {
			if strings.HasPrefix(k, p) {
				key := k
				if opts.StripPrefix {
					if trimmed := strings.TrimPrefix(k, p); trimmed != "" {
						key = trimmed
					}
				}
				groups[p][key] = v
				matched = true
				break
			}
		}
		if !matched {
			remainder[k] = v
		}
	}

	prefixesSorted := make([]string, len(opts.Prefixes))
	copy(prefixesSorted, opts.Prefixes)
	sort.Strings(prefixesSorted)

	var parts []Part
	for _, p := range prefixesSorted {
		parts = append(parts, Part{Name: p, Env: groups[p]})
	}
	if opts.KeepRemainder && len(remainder) > 0 {
		parts = append(parts, Part{Name: "_remainder", Env: remainder})
	}
	return parts, nil
}
