package envdoc

import (
	"fmt"
	"sort"
	"strings"
)

// Entry holds documentation metadata for a single env var.
type Entry struct {
	Key         string
	Value       string
	Description string
	Required    bool
	Sensitive   bool
}

// Result holds all documented entries.
type Result struct {
	Entries []Entry
}

// Options controls documentation generation behaviour.
type Options struct {
	// Descriptions maps key names to human-readable descriptions.
	Descriptions map[string]string
	// RequiredKeys marks certain keys as required.
	RequiredKeys []string
	// SensitiveKeys marks certain keys as sensitive.
	SensitiveKeys []string
}

// Document annotates an env map with descriptions, required and sensitive flags.
func Document(env map[string]string, opts Options) Result {
	required := toSet(opts.RequiredKeys)
	sensitive := toSet(opts.SensitiveKeys)

	keys := sortedKeys(env)
	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		if sensitive[k] {
			v = "***"
		}
		entries = append(entries, Entry{
			Key:         k,
			Value:       v,
			Description: opts.Descriptions[k],
			Required:    required[k],
			Sensitive:   sensitive[k],
		})
	}
	return Result{Entries: entries}
}

// Render formats the Result as a markdown table.
func Render(r Result) string {
	var sb strings.Builder
	sb.WriteString("| Key | Value | Required | Sensitive | Description |\n")
	sb.WriteString("|-----|-------|----------|-----------|-------------|\n")
	for _, e := range r.Entries {
		req := boolMark(e.Required)
		sen := boolMark(e.Sensitive)
		fmt.Fprintf(&sb, "| %s | %s | %s | %s | %s |\n", e.Key, e.Value, req, sen, e.Description)
	}
	return sb.String()
}

func boolMark(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
