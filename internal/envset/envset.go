// Package envset provides set operations (union, intersection, difference) on env maps.
package envset

import "sort"

// Result holds the output of a set operation.
type Result struct {
	Env  map[string]string
	Keys []string
}

// Union returns all keys from both a and b. Values from b override a on conflict.
func Union(a, b map[string]string) Result {
	out := copyMap(a)
	for k, v := range b {
		out[k] = v
	}
	return Result{Env: out, Keys: sortedKeys(out)}
}

// Intersect returns keys present in both a and b, with values from a.
func Intersect(a, b map[string]string) Result {
	out := make(map[string]string)
	for k, v := range a {
		if _, ok := b[k]; ok {
			out[k] = v
		}
	}
	return Result{Env: out, Keys: sortedKeys(out)}
}

// Difference returns keys in a that are not in b.
func Difference(a, b map[string]string) Result {
	out := make(map[string]string)
	for k, v := range a {
		if _, ok := b[k]; !ok {
			out[k] = v
		}
	}
	return Result{Env: out, Keys: sortedKeys(out)}
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
