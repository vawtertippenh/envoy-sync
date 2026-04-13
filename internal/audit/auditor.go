package audit

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// EntryKind represents the type of audit event.
type EntryKind string

const (
	KindAdded   EntryKind = "added"
	KindRemoved EntryKind = "removed"
	KindChanged EntryKind = "changed"
	KindMasked  EntryKind = "masked"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time
	Key       string
	Kind      EntryKind
	Detail    string
}

// Log holds a collection of audit entries.
type Log struct {
	Entries []Entry
}

// Record appends a new entry to the log.
func (l *Log) Record(key string, kind EntryKind, detail string) {
	l.Entries = append(l.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Key:       key,
		Kind:      kind,
		Detail:    detail,
	})
}

// Summary returns a human-readable summary of the audit log.
func (l *Log) Summary() string {
	if len(l.Entries) == 0 {
		return "audit log: no events recorded"
	}

	counts := map[EntryKind]int{}
	for _, e := range l.Entries {
		counts[e.Kind]++
	}

	kinds := make([]string, 0, len(counts))
	for k := range counts {
		kinds = append(kinds, string(k))
	}
	sort.Strings(kinds)

	parts := make([]string, 0, len(kinds))
	for _, k := range kinds {
		parts = append(parts, fmt.Sprintf("%s=%d", k, counts[EntryKind(k)]))
	}

	return fmt.Sprintf("audit log: %d event(s) [%s]", len(l.Entries), strings.Join(parts, ", "))
}

// FilterByKind returns entries matching the given kind.
func (l *Log) FilterByKind(kind EntryKind) []Entry {
	var result []Entry
	for _, e := range l.Entries {
		if e.Kind == kind {
			result = append(result, e)
		}
	}
	return result
}
