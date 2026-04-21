// Package envlookup provides key-based lookup and inspection of env maps.
package envlookup

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the outcome of a single key lookup.
type Result struct {
	Key     string
	Value   string
	Found   bool
	Masked  bool
}

// Options configures lookup behaviour.
type Options struct {
	// Keys is the list of keys to look up. If empty, all keys are returned.
	Keys []string
	// MaskSensitive replaces values for sensitive keys with "***".
	MaskSensitive bool
	// SensitivePatterns extends the default sensitive-key heuristics.
	SensitivePatterns []string
	// CaseFold performs case-insensitive key matching.
	CaseFold bool
}

// Lookup searches env for the requested keys and returns ordered results.
func Lookup(env map[string]string, opts Options) []Result {
	keys := opts.Keys
	if len(keys) == 0 {
		keys = sortedKeys(env)
	}

	results := make([]Result, 0, len(keys))
	for _, k := range keys {
		result := resolve(env, k, opts)
		results = append(results, result)
	}
	return results
}

func resolve(env map[string]string, key string, opts Options) Result {
	v, ok := env[key]
	if !ok && opts.CaseFold {
		for ek, ev := range env {
			if strings.EqualFold(ek, key) {
				v, ok = ev, true
				key = ek
				break
			}
		}
	}

	if !ok {
		return Result{Key: key, Found: false}
	}

	masked := false
	if opts.MaskSensitive && isSensitive(key, opts.SensitivePatterns) {
		v = "***"
		masked = true
	}
	return Result{Key: key, Value: v, Found: true, Masked: masked}
}

func isSensitive(key string, extra []string) bool {
	lower := strings.ToLower(key)
	defaults := []string{"password", "secret", "token", "api_key", "apikey", "private_key"}
	for _, p := range append(defaults, extra...) {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// Render formats results as human-readable lines.
func Render(results []Result) string {
	var sb strings.Builder
	for _, r := range results {
		if !r.Found {
			sb.WriteString(fmt.Sprintf("%-30s  (not found)\n", r.Key))
			continue
		}
		suffix := ""
		if r.Masked {
			suffix = "  [masked]"
		}
		sb.WriteString(fmt.Sprintf("%-30s = %s%s\n", r.Key, r.Value, suffix))
	}
	return sb.String()
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
