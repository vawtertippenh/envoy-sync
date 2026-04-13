// Package promote handles promoting env variables from one environment to another,
// applying optional filters and transformations in the process.
package promote

import (
	"fmt"
	"sort"
)

// Options configures how promotion behaves.
type Options struct {
	// Keys restricts promotion to only these keys. Empty means all keys.
	Keys []string
	// Overwrite controls whether existing keys in the target are overwritten.
	Overwrite bool
	// DryRun reports what would change without modifying the target.
	DryRun bool
}

// Result describes the outcome of a promotion.
type Result struct {
	Promoted  []string
	Skipped   []string
	Overwrite []string
}

// Summary returns a human-readable summary of the promotion result.
func (r Result) Summary() string {
	return fmt.Sprintf(
		"promoted=%d skipped=%d overwritten=%d",
		len(r.Promoted), len(r.Skipped), len(r.Overwrite),
	)
}

// Promote copies keys from src into dst according to opts.
// It returns the modified dst map and a Result describing what happened.
func Promote(src, dst map[string]string, opts Options) (map[string]string, Result, error) {
	if src == nil {
		return dst, Result{}, fmt.Errorf("promote: source env must not be nil")
	}
	if dst == nil {
		dst = make(map[string]string)
	}

	allow := toSet(opts.Keys)
	out := copyMap(dst)
	var res Result

	for _, k := range sortedKeys(src) {
		if len(allow) > 0 && !allow[k] {
			continue
		}
		_, exists := out[k]
		if exists && !opts.Overwrite {
			res.Skipped = append(res.Skipped, k)
			continue
		}
		if !opts.DryRun {
			out[k] = src[k]
		}
		if exists {
			res.Overwrite = append(res.Overwrite, k)
		} else {
			res.Promoted = append(res.Promoted, k)
		}
	}

	if opts.DryRun {
		return dst, res, nil
	}
	return out, res, nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
