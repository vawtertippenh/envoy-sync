// Package scope provides filtering of env maps by key prefix or suffix patterns.
package scope

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the filtered environment and metadata.
type Result struct {
	Matched   map[string]string
	Unmatched map[string]string
}

// Summary returns a human-readable summary of the scope result.
func (r Result) Summary() string {
	return fmt.Sprintf("%d matched, %d unmatched", len(r.Matched), len(r.Unmatched))
}

// Options controls how scoping is applied.
type Options struct {
	Prefixes []string
	Suffixes []string
	Strip    bool // strip the matched prefix from the key
}

// Scope filters env by the given prefix or suffix patterns.
// If no patterns are provided, all keys are returned as matched.
func Scope(env map[string]string, opts Options) Result {
	matched := make(map[string]string)
	unmatched := make(map[string]string)

	for _, k := range sortedKeys(env) {
		v := env[k]
		if matches(k, opts) {
			outKey := k
			if opts.Strip {
				outKey = stripPrefix(k, opts.Prefixes)
			}
			matched[outKey] = v
		} else {
			unmatched[k] = v
		}
	}

	return Result{Matched: matched, Unmatched: unmatched}
}

func matches(key string, opts Options) bool {
	if len(opts.Prefixes) == 0 && len(opts.Suffixes) == 0 {
		return true
	}
	for _, p := range opts.Prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	for _, s := range opts.Suffixes {
		if strings.HasSuffix(key, s) {
			return true
		}
	}
	return false
}

func stripPrefix(key string, prefixes []string) string {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return strings.TrimPrefix(key, p)
		}
	}
	return key
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
