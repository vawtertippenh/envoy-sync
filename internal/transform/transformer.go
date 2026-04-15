package transform

import (
	"fmt"
	"sort"
	"strings"
)

// Rule defines a single transformation to apply to env values.
type Rule struct {
	Key   string // exact key to target, or "*" for all
	Op    string // "upper", "lower", "trim", "replace"
	From  string // used by "replace"
	To    string // used by "replace"
}

// Result holds the outcome of a Transform call.
type Result struct {
	Env     map[string]string
	Changed []string // keys whose values were modified
}

// Transform applies a list of Rules to the given env map.
// It returns a new map and a Result describing what changed.
func Transform(env map[string]string, rules []Rule) (Result, error) {
	out := copyMap(env)
	changed := []string{}

	for _, rule := range rules {
		keys := targetKeys(out, rule.Key)
		for _, k := range keys {
			original := out[k]
			newVal, err := applyOp(original, rule)
			if err != nil {
				return Result{}, fmt.Errorf("rule op %q on key %q: %w", rule.Op, k, err)
			}
			if newVal != original {
				out[k] = newVal
				changed = append(changed, k)
			}
		}
	}

	sort.Strings(changed)
	return Result{Env: out, Changed: changed}, nil
}

func applyOp(val string, rule Rule) (string, error) {
	switch rule.Op {
	case "upper":
		return strings.ToUpper(val), nil
	case "lower":
		return strings.ToLower(val), nil
	case "trim":
		return strings.TrimSpace(val), nil
	case "replace":
		return strings.ReplaceAll(val, rule.From, rule.To), nil
	default:
		return "", fmt.Errorf("unknown op %q", rule.Op)
	}
}

func targetKeys(env map[string]string, key string) []string {
	if key == "*" {
		keys := make([]string, 0, len(env))
		for k := range env {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}
	if _, ok := env[key]; ok {
		return []string{key}
	}
	return nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
