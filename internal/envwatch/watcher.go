package envwatch

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// WatchResult holds the outcome of a watch check.
type WatchResult struct {
	Changed  bool
	OldHash  string
	NewHash  string
	DiffKeys []string
}

// Watch compares two env maps and reports whether anything changed.
// It returns a WatchResult with the affected keys and hashes.
func Watch(previous, current map[string]string) WatchResult {
	oldHash := hashEnv(previous)
	newHash := hashEnv(current)

	if oldHash == newHash {
		return WatchResult{Changed: false, OldHash: oldHash, NewHash: newHash}
	}

	diffKeys := changedKeys(previous, current)
	return WatchResult{
		Changed:  true,
		OldHash:  oldHash,
		NewHash:  newHash,
		DiffKeys: diffKeys,
	}
}

// Summary returns a human-readable summary of the WatchResult.
func (r WatchResult) Summary() string {
	if !r.Changed {
		return "no changes detected"
	}
	return fmt.Sprintf("%d key(s) changed: %s", len(r.DiffKeys), strings.Join(r.DiffKeys, ", "))
}

func hashEnv(env map[string]string) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s;", k, env[k])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func changedKeys(prev, curr map[string]string) []string {
	seen := map[string]bool{}
	var keys []string

	for k, v := range curr {
		if prev[k] != v {
			keys = append(keys, k)
		}
		seen[k] = true
	}
	for k := range prev {
		if !seen[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}
