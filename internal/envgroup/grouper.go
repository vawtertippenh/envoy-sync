package envgroup

import "sort"

// Group holds a named collection of env entries.
type Group struct {
	Name string
	Keys []string
	Env  map[string]string
}

// Options controls how grouping is performed.
type Options struct {
	// Prefixes maps a group name to a key prefix.
	Prefixes map[string]string
	// Remainder is the name for keys that don't match any prefix.
	// If empty, unmatched keys are dropped.
	Remainder string
}

// GroupBy splits env into named groups based on key prefixes.
func GroupBy(env map[string]string, opts Options) []Group {
	groupMap := make(map[string]map[string]string)

	for key, val := range env {
		matched := false
		for groupName, prefix := range opts.Prefixes {
			if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
				if groupMap[groupName] == nil {
					groupMap[groupName] = make(map[string]string)
				}
				groupMap[groupName][key] = val
				matched = true
				break
			}
		}
		if !matched && opts.Remainder != "" {
			if groupMap[opts.Remainder] == nil {
				groupMap[opts.Remainder] = make(map[string]string)
			}
			groupMap[opts.Remainder][key] = val
		}
	}

	names := make([]string, 0, len(groupMap))
	for n := range groupMap {
		names = append(names, n)
	}
	sort.Strings(names)

	result := make([]Group, 0, len(names))
	for _, n := range names {
		keys := sortedKeys(groupMap[n])
		result = append(result, Group{Name: n, Keys: keys, Env: groupMap[n]})
	}
	return result
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
