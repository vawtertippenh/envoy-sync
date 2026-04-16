// Package diff2 provides line-by-line diff of two env maps with change classification.
package diff2

import "sort"

// ChangeType represents the kind of change detected.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Entry represents a single diff entry.
type Entry struct {
	Key    string
	OldVal string
	NewVal string
	Change ChangeType
}

// Result holds all diff entries and metadata.
type Result struct {
	Entries []Entry
}

// Added returns entries that were added.
func (r Result) Added() []Entry { return r.filter(Added) }

// Removed returns entries that were removed.
func (r Result) Removed() []Entry { return r.filter(Removed) }

// Modified returns entries that were modified.
func (r Result) Modified() []Entry { return r.filter(Modified) }

// HasChanges returns true if any non-unchanged entries exist.
func (r Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Change != Unchanged {
			return true
		}
	}
	return false
}

func (r Result) filter(ct ChangeType) []Entry {
	var out []Entry
	for _, e := range r.Entries {
		if e.Change == ct {
			out = append(out, e)
		}
	}
	return out
}

// Diff computes the difference between env maps a (old) and b (new).
func Diff(a, b map[string]string) Result {
	seen := map[string]bool{}
	var entries []Entry

	for k, av := range a {
		seen[k] = true
		if bv, ok := b[k]; ok {
			if av == bv {
				entries = append(entries, Entry{Key: k, OldVal: av, NewVal: bv, Change: Unchanged})
			} else {
				entries = append(entries, Entry{Key: k, OldVal: av, NewVal: bv, Change: Modified})
			}
		} else {
			entries = append(entries, Entry{Key: k, OldVal: av, Change: Removed})
		}
	}

	for k, bv := range b {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, NewVal: bv, Change: Added})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Result{Entries: entries}
}
