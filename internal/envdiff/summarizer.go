package envdiff

import "sort"

// ChangeKind represents the type of change between two env maps.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
)

// Change describes a single key-level difference.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Summary holds the full diff result between two env maps.
type Summary struct {
	Changes []Change
}

// HasDrift returns true if there are any changes.
func (s Summary) HasDrift() bool {
	return len(s.Changes) > 0
}

// Counts returns added, removed, modified counts.
func (s Summary) Counts() (added, removed, modified int) {
	for _, c := range s.Changes {
		switch c.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return
}

// Summarize computes the diff between base and target env maps.
func Summarize(base, target map[string]string) Summary {
	changes := []Change{}

	for k, bv := range base {
		if tv, ok := target[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: bv})
		} else if tv != bv {
			changes = append(changes, Change{Key: k, Kind: Modified, OldValue: bv, NewValue: tv})
		}
	}

	for k, tv := range target {
		if _, ok := base[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: tv})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return Summary{Changes: changes}
}
