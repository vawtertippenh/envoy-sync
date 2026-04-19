package envchain

import "sort"

// Link represents a single env map in the chain.
type Link struct {
	Name string
	Env  map[string]string
}

// Result holds the resolved env and metadata about each key's origin.
type Result struct {
	Env    map[string]string
	Origin map[string]string // key -> link name
}

// Chain resolves a sequence of env links, where later links override earlier ones.
// If stopOnFirst is true, resolution stops at the first link that defines a key.
func Chain(links []Link, stopOnFirst bool) Result {
	env := make(map[string]string)
	origin := make(map[string]string)

	if stopOnFirst {
		for _, link := range links {
			for k, v := range link.Env {
				if _, exists := env[k]; !exists {
					env[k] = v
					origin[k] = link.Name
				}
			}
		}
	} else {
		for _, link := range links {
			for k, v := range link.Env {
				env[k] = v
				origin[k] = link.Name
			}
		}
	}

	return Result{Env: env, Origin: origin}
}

// Summary returns a human-readable summary of origins.
func Summary(r Result) []string {
	keys := make([]string, 0, len(r.Env))
	for k := range r.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, k+" (from "+r.Origin[k]+")")
	}
	return lines
}
