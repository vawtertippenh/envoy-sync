// Package defaults provides functionality to apply default values
// to missing or empty keys in an env map.
package defaults

import "sort"

// Rule defines a default value rule for a single key.
type Rule struct {
	Key      string
	Value    string
	Override bool // if true, overwrite even if key exists
}

// Result holds the outcome of applying defaults.
type Result struct {
	Env     map[string]string
	Applied []string // keys that received a default value
	Skipped []string // keys that already had a value
}

// Apply applies the given default rules to the env map.
// It returns a new map and a Result describing what changed.
func Apply(env map[string]string, rules []Rule) Result {
	out := copyMap(env)
	var applied, skipped []string

	for _, r := range rules {
		existing, exists := out[r.Key]
		if !exists || existing == "" || r.Override {
			out[r.Key] = r.Value
			applied = append(applied, r.Key)
		} else {
			skipped = append(skipped, r.Key)
		}
	}

	sort.Strings(applied)
	sort.Strings(skipped)

	return Result{
		Env:     out,
		Applied: applied,
		Skipped: skipped,
	}
}

// Summary returns a human-readable summary of the result.
func (r Result) Summary() string {
	if len(r.Applied) == 0 {
		return "no defaults applied"
	}
	msg := "applied defaults for: "
	for i, k := range r.Applied {
		if i > 0 {
			msg += ", "
		}
		msg += k
	}
	return msg
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
