// Package envmerge provides deep merging of multiple env maps with
// conflict resolution strategies and priority ordering.
package envmerge

import "fmt"

// Strategy defines how key conflicts are resolved.
type Strategy string

const (
	StrategyFirst Strategy = "first" // keep first occurrence
	StrategyLast  Strategy = "last"  // keep last occurrence
	StrategyError Strategy = "error" // return error on conflict
)

// Options configures the merge behaviour.
type Options struct {
	Strategy Strategy
	Prefix   string // optional prefix to apply to all result keys
}

// Result holds the merged map and metadata.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Conflict records a key that had differing values across sources.
type Conflict struct {
	Key    string
	Values []string
}

// Merge combines sources according to opts.
func Merge(sources []map[string]string, opts Options) (Result, error) {
	if opts.Strategy == "" {
		opts.Strategy = StrategyLast
	}

	env := make(map[string]string)
	seen := make(map[string]string) // key -> first value
	var conflicts []Conflict
	conflictIndex := make(map[string][]string)

	for _, src := range sources {
		for k, v := range src {
			key := opts.Prefix + k
			if prev, exists := seen[key]; exists && prev != v {
				if opts.Strategy == StrategyError {
					return Result{}, fmt.Errorf("conflict on key %q: %q vs %q", key, prev, v)
				}
				conflictIndex[key] = append(conflictIndex[key], v)
			}
			if _, exists := seen[key]; !exists {
				seen[key] = v
				conflictIndex[key] = []string{v}
			} else {
				conflictIndex[key] = append(conflictIndex[key], v)
			}
			switch opts.Strategy {
			case StrategyFirst:
				if _, ok := env[key]; !ok {
					env[key] = v
				}
			default: // last
				env[key] = v
			}
		}
	}

	for key, vals := range conflictIndex {
		if len(vals) > 1 {
			uniq := dedupe(vals)
			if len(uniq) > 1 {
				conflicts = append(conflicts, Conflict{Key: key, Values: uniq})
			}
		}
	}

	return Result{Env: env, Conflicts: conflicts}, nil
}

func dedupe(ss []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
