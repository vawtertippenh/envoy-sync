// Package envindex builds a searchable index over an env map,
// allowing fast key lookup by prefix, suffix, or substring.
package envindex

import "sort"

// Entry holds a single indexed key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Index is a searchable collection of env entries.
type Index struct {
	entries []Entry
}

// Build creates a new Index from the given env map.
func Build(env map[string]string) *Index {
	entries := make([]Entry, 0, len(env))
	for k, v := range env {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return &Index{entries: entries}
}

// ByPrefix returns all entries whose key starts with the given prefix.
func (idx *Index) ByPrefix(prefix string) []Entry {
	var out []Entry
	for _, e := range idx.entries {
		if len(e.Key) >= len(prefix) && e.Key[:len(prefix)] == prefix {
			out = append(out, e)
		}
	}
	return out
}

// BySuffix returns all entries whose key ends with the given suffix.
func (idx *Index) BySuffix(suffix string) []Entry {
	var out []Entry
	for _, e := range idx.entries {
		if len(e.Key) >= len(suffix) && e.Key[len(e.Key)-len(suffix):] == suffix {
			out = append(out, e)
		}
	}
	return out
}

// BySubstring returns all entries whose key contains the given substring.
func (idx *Index) BySubstring(sub string) []Entry {
	var out []Entry
	for _, e := range idx.entries {
		if contains(e.Key, sub) {
			out = append(out, e)
		}
	}
	return out
}

// All returns every entry in the index in sorted key order.
func (idx *Index) All() []Entry {
	return idx.entries
}

func contains(s, sub string) bool {
	if sub == "" {
		return true
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
