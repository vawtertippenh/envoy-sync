package envtag

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a label attached to an env key.
type Tag struct {
	Key   string
	Value string
}

// Result holds the tagged output for a single env key.
type Result struct {
	EnvKey string
	Tags   []Tag
}

// Options controls tagging behaviour.
type Options struct {
	// Tags maps tag name -> list of env key patterns (prefix* or exact)
	Tags map[string][]string
	// DefaultTag is applied to keys that match no explicit tag.
	DefaultTag string
}

// Tag applies tags to env keys based on pattern rules.
func Tag(env map[string]string, opts Options) []Result {
	keys := sortedKeys(env)
	results := make([]Result, 0, len(keys))

	for _, k := range keys {
		matched := []Tag{}
		for tagName, patterns := range opts.Tags {
			for _, pat := range patterns {
				if matchPattern(k, pat) {
					matched = append(matched, Tag{Key: tagName, Value: k})
					break
				}
			}
		}
		if len(matched) == 0 && opts.DefaultTag != "" {
			matched = append(matched, Tag{Key: opts.DefaultTag, Value: k})
		}
		results = append(results, Result{EnvKey: k, Tags: matched})
	}
	return results
}

// Render returns a human-readable summary of tagged results.
func Render(results []Result) string {
	var sb strings.Builder
	for _, r := range results {
		if len(r.Tags) == 0 {
			fmt.Fprintf(&sb, "%s: (untagged)\n", r.EnvKey)
			continue
		}
		tagStrs := make([]string, 0, len(r.Tags))
		for _, t := range r.Tags {
			tagStrs = append(tagStrs, t.Key)
		}
		fmt.Fprintf(&sb, "%s: [%s]\n", r.EnvKey, strings.Join(tagStrs, ", "))
	}
	return sb.String()
}

func matchPattern(key, pattern string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}
	return key == pattern
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
