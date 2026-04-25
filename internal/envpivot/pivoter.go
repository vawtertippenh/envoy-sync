// Package envpivot provides functionality to pivot env maps by value,
// grouping keys that share the same value together.
package envpivot

import (
	"fmt"
	"sort"
)

// Group holds a set of keys that share a common value.
type Group struct {
	Value string
	Keys  []string
}

// Result is the output of a Pivot operation.
type Result struct {
	Groups     []Group
	Singletons int // groups with only one key
	Shared     int // groups with more than one key
}

// Pivot inverts an env map, grouping keys by their value.
// Keys within each group are sorted alphabetically.
// Groups are sorted by value.
func Pivot(env map[string]string) Result {
	index := make(map[string][]string)
	for k, v := range env {
		index[v] = append(index[v], k)
	}

	groups := make([]Group, 0, len(index))
	for val, keys := range index {
		sort.Strings(keys)
		groups = append(groups, Group{Value: val, Keys: keys})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Value < groups[j].Value
	})

	var singletons, shared int
	for _, g := range groups {
		if len(g.Keys) == 1 {
			singletons++
		} else {
			shared++
		}
	}

	return Result{
		Groups:     groups,
		Singletons: singletons,
		Shared:     shared,
	}
}

// Summary returns a human-readable summary of the pivot result.
func Summary(r Result) string {
	if len(r.Groups) == 0 {
		return "no entries"
	}
	total := r.Singletons + r.Shared
	if r.Shared == 0 {
		return fmt.Sprintf("%d unique values, no shared values", total)
	}
	return fmt.Sprintf("%d value groups (%d shared, %d unique)", total, r.Shared, r.Singletons)
}
