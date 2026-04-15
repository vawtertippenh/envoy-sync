// Package dedupe provides functionality to detect and remove duplicate
// keys from a parsed .env map, reporting any conflicts found.
package dedupe

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the deduplicated environment map and any duplicate keys found.
type Result struct {
	Env        map[string]string
	Duplicates []string
}

// Summary returns a human-readable summary of the deduplication result.
func (r Result) Summary() string {
	if len(r.Duplicates) == 0 {
		return "No duplicate keys found."
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d duplicate key(s) removed:\n", len(r.Duplicates))
	for _, k := range r.Duplicates {
		fmt.Fprintf(&sb, "  - %s\n", k)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Dedupe scans rawLines (each element being a raw "KEY=VALUE" line or comment)
// and returns a Result where only the first occurrence of each key is kept.
// If keepLast is true, the last occurrence wins instead.
func Dedupe(rawLines []string, keepLast bool) Result {
	seen := make(map[string]int)   // key -> index of first/last kept line
	dupeSet := make(map[string]bool)
	var ordered []string // keys in insertion order
	values := make(map[string]string)

	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if _, exists := seen[key]; exists {
			dupeSet[key] = true
			if keepLast {
				values[key] = val
			}
		} else {
			seen[key] = len(ordered)
			ordered = append(ordered, key)
			values[key] = val
		}
	}

	env := make(map[string]string, len(ordered))
	for _, k := range ordered {
		env[k] = values[k]
	}

	duplicates := make([]string, 0, len(dupeSet))
	for k := range dupeSet {
		duplicates = append(duplicates, k)
	}
	sort.Strings(duplicates)

	return Result{Env: env, Duplicates: duplicates}
}
