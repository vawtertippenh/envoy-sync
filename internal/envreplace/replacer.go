// Package envreplace provides functionality to find and replace values
// across an env map using string or pattern-based rules.
package envreplace

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a single find-and-replace operation.
type Rule struct {
	Find    string
	Replace string
	Regex   bool
}

// Result holds the outcome of a Replace operation.
type Result struct {
	Env         map[string]string
	ChangedKeys []string
}

// Replace applies the given rules to every value in env.
// Keys are never modified. A new map is returned; the original is not mutated.
func Replace(env map[string]string, rules []Rule) (Result, error) {
	out := copyMap(env)
	changed := []string{}

	for _, rule := range rules {
		if rule.Find == "" {
			return Result{}, fmt.Errorf("replace rule has empty Find field")
		}

		var re *regexp.Regexp
		if rule.Regex {
			var err error
			re, err = regexp.Compile(rule.Find)
			if err != nil {
				return Result{}, fmt.Errorf("invalid regex %q: %w", rule.Find, err)
			}
		}

		for _, k := range sortedKeys(out) {
			old := out[k]
			var updated string
			if rule.Regex {
				updated = re.ReplaceAllString(old, rule.Replace)
			} else {
				updated = strings.ReplaceAll(old, rule.Find, rule.Replace)
			}
			if updated != old {
				out[k] = updated
				changed = appendUnique(changed, k)
			}
		}
	}

	return Result{Env: out, ChangedKeys: changed}, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func appendUnique(s []string, v string) []string {
	for _, x := range s {
		if x == v {
			return s
		}
	}
	return append(s, v)
}
