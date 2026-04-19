package envfreeze

import (
	"errors"
	"sort"
)

// FreezeResult holds the outcome of a freeze operation.
type FreezeResult struct {
	Frozen map[string]string
	Skipped []string
}

// Options configures the Freeze operation.
type Options struct {
	// AllowKeys, if set, only freezes these keys.
	AllowKeys []string
	// DenyKeys are keys to exclude from freezing.
	DenyKeys []string
	// OverwriteExisting replaces already-frozen values.
	OverwriteExisting bool
}

// Freeze locks the values in src into a frozen copy, optionally filtered.
// Keys whose values are empty are skipped and reported in Skipped.
func Freeze(src map[string]string, existing map[string]string, opts Options) (FreezeResult, error) {
	if src == nil {
		return FreezeResult{}, errors.New("envfreeze: src map must not be nil")
	}

	deny := toSet(opts.DenyKeys)
	allow := toSet(opts.AllowKeys)

	result := FreezeResult{
		Frozen: copyMap(existing),
	}

	for _, k := range sortedKeys(src) {
		if len(deny) > 0 && deny[k] {
			continue
		}
		if len(allow) > 0 && !allow[k] {
			continue
		}
		v := src[k]
		if v == "" {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if _, exists := result.Frozen[k]; exists && !opts.OverwriteExisting {
			continue
		}
		result.Frozen[k] = v
	}
	return result, nil
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
