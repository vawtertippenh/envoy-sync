package envdiff

import "sort"

// ChangeKind describes the type of change between two env maps.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds the full diff summary.
type Result struct {
	Changes []Change
}

// Summarize compares two env maps and returns a Result.
func Summarize(base, target map[string]string) Result {
	seen := map[string]bool{}
	var changes []Change

	for k, bv := range base {
		seen[k] = true
		if tv, ok := target[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: bv})
		} else if bv != tv {
			changes = append(changes, Change{Key: k, Kind: Modified, OldValue: bv, NewValue: tv})
		} else {
			changes = append(changes, Change{Key: k, Kind: Unchanged, OldValue: bv, NewValue: tv})
		}
	}

	for k, tv := range target {
		if !seen[k] {
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: tv})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return Result{Changes: changes}
}

// Added returns only added changes.
func (r Result) Added() []Change { return r.filter(Added) }

// Removed returns only removed changes.
func (r Result) Removed() []Change { return r.filter(Removed) }

// Modified returns only modified changes.
func (r Result) Modified() []Change { return r.filter(Modified) }

// HasDrift returns true if there are any non-unchanged entries.
func (r Result) HasDrift() bool {
	for _, c := range r.Changes {
		if c.Kind != Unchanged {
			return true
		}
	}
	return false
}

func (r Result) filter(kind ChangeKind) []Change {
	var out []Change
	for _, c := range r.Changes {
		if c.Kind == kind {
			out = append(out, c)
		}
	}
	return out
}
