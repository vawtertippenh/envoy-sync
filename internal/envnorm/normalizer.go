package envnorm

import (
	"sort"
	"strings"
)

// Options controls normalization behavior.
type Options struct {
	UppercaseKeys   bool
	TrimValues      bool
	RemoveEmpty     bool
	SortKeys        bool
}

// Result holds the normalized env map and a log of changes.
type Result struct {
	Env     map[string]string
	Changes []Change
}

// Change describes a single normalization action.
type Change struct {
	Key    string
	Action string // "uppercase_key", "trim_value", "removed_empty"
	OldKey string // populated when key was renamed
}

// Normalize applies normalization rules to env according to opts.
func Normalize(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	var changes []Change

	for k, v := range env {
		newKey := k
		newVal := v

		if opts.TrimValues {
			trimmed := strings.TrimSpace(v)
			if trimmed != v {
				changes = append(changes, Change{Key: k, Action: "trim_value"})
				newVal = trimmed
			}
		}

		if opts.RemoveEmpty && newVal == "" {
			changes = append(changes, Change{Key: k, Action: "removed_empty"})
			continue
		}

		if opts.UppercaseKeys {
			up := strings.ToUpper(k)
			if up != k {
				changes = append(changes, Change{Key: up, OldKey: k, Action: "uppercase_key"})
				newKey = up
			}
		}

		out[newKey] = newVal
	}

	if opts.SortKeys {
		sorted := make(map[string]string, len(out))
		keys := make([]string, 0, len(out))
		for k := range out {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sorted[k] = out[k]
		}
		out = sorted
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return Result{Env: out, Changes: changes}
}
