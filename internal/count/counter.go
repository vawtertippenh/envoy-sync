package count

import (
	"fmt"
	"sort"
)

// Result holds the counts derived from an env map.
type Result struct {
	Total     int
	Empty     int
	NonEmpty  int
	Sensitive int
	Prefixes  map[string]int // count of keys per detected prefix (e.g. "DB", "AWS")
}

// Options controls how counting is performed.
type Options struct {
	// SensitivePatterns are additional key patterns to consider sensitive.
	SensitivePatterns []string
	// PrefixSep is the separator used to detect key prefixes (default "_").
	PrefixSep string
}

var defaultSensitivePatterns = []string{
	"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL", "AUTH",
}

// Count analyses env and returns a Result.
func Count(env map[string]string, opts Options) Result {
	sep := opts.PrefixSep
	if sep == "" {
		sep = "_"
	}

	sensitiveSet := toSet(defaultSensitivePatterns)
	for _, p := range opts.SensitivePatterns {
		sensitiveSet[p] = struct{}{}
	}

	prefixes := map[string]int{}
	r := Result{Prefixes: prefixes}

	for k, v := range env {
		r.Total++
		if v == "" {
			r.Empty++
		} else {
			r.NonEmpty++
		}
		if isSensitive(k, sensitiveSet) {
			r.Sensitive++
		}
		if idx := indexOf(k, sep); idx > 0 {
			prefix := k[:idx]
			prefixes[prefix]++
		}
	}
	return r
}

// Summary returns a human-readable summary of the Result.
func Summary(r Result) string {
	lines := []string{
		fmt.Sprintf("Total keys   : %d", r.Total),
		fmt.Sprintf("Non-empty    : %d", r.NonEmpty),
		fmt.Sprintf("Empty        : %d", r.Empty),
		fmt.Sprintf("Sensitive    : %d", r.Sensitive),
	}
	if len(r.Prefixes) > 0 {
		lines = append(lines, "Prefixes:")
		keys := make([]string, 0, len(r.Prefixes))
		for k := range r.Prefixes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			lines = append(lines, fmt.Sprintf("  %s: %d", k, r.Prefixes[k]))
		}
	}
	out := ""
	for _, l := range lines {
		out += l + "\n"
	}
	return out
}

func isSensitive(key string, patterns map[string]struct{}) bool {
	for p := range patterns {
		if contains(key, p) {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub ||
		len(s) > 0 && (indexOf(s, sub) >= 0))
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func toSet(ss []string) map[string]struct{} {
	m := make(map[string]struct{}, len(ss))
	for _, s := range ss {
		m[s] = struct{}{}
	}
	return m
}
