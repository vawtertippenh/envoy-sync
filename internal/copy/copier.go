package copy

import "sort"

// Options controls Copy behaviour.
type Options struct {
	Keys      []string // if non-empty, only copy these keys
	Overwrite bool     // overwrite existing keys in dst
	Prefix    string   // add prefix to copied keys
}

// Result holds the outcome of a Copy operation.
type Result struct {
	Copied   []string
	Skipped  []string
}

// Copy copies keys from src into dst according to opts.
// dst is not mutated; a new map is returned.
func Copy(src, dst map[string]string, opts Options) (map[string]string, Result) {
	out := copyMap(dst)
	var res Result

	keys := opts.Keys
	if len(keys) == 0 {
		keys = sortedKeys(src)
	}

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			res.Skipped = append(res.Skipped, k)
			continue
		}
		destKey := opts.Prefix + k
		if _, exists := out[destKey]; exists && !opts.Overwrite {
			res.Skipped = append(res.Skipped, k)
			continue
		}
		out[destKey] = v
		res.Copied = append(res.Copied, k)
	}
	return out, res
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
