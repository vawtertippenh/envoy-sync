// Package envdiff2 provides a structured diff between two env maps,
// reporting added, removed, and modified keys with optional masking.
package envdiff2

import "sort"

// ChangeKind describes the type of change for a key.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level difference.
type Change struct {
	Key    string
	Kind   ChangeKind
	OldVal string
	NewVal string
}

// Result holds the full diff output.
type Result struct {
	Changes []Change
}

// HasDiff returns true if any non-unchanged entries exist.
func (r Result) HasDiff() bool {
	for _, c := range r.Changes {
		if c.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Diff computes a structured diff between env maps a and b.
// If includeUnchanged is true, unchanged keys are also included.
func Diff(a, b map[string]string, includeUnchanged bool) Result {
	keys := unionKeys(a, b)
	var changes []Change
	for _, k := range keys {
		av, aok := a[k]
		bv, bok := b[k]
		switch {
		case aok && !bok:
			changes = append(changes, Change{Key: k, Kind: Removed, OldVal: av})
		case !aok && bok:
			changes = append(changes, Change{Key: k, Kind: Added, NewVal: bv})
		case av != bv:
			changes = append(changes, Change{Key: k, Kind: Modified, OldVal: av, NewVal: bv})
		default:
			if includeUnchanged {
				changes = append(changes, Change{Key: k, Kind: Unchanged, OldVal: av, NewVal: bv})
			}
		}
	}
	return Result{Changes: changes}
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
