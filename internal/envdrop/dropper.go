package envdrop

import "sort"

// Options controls which keys are dropped.
type Options struct {
	Keys     []string // exact key names to drop
	Prefixes []string // drop keys with these prefixes
	Suffixes []string // drop keys with these suffixes
	DryRun   bool
}

// Result holds the output of a Drop operation.
type Result struct {
	Out     map[string]string
	Dropped []string
}

// Drop removes keys from env according to opts.
func Drop(env map[string]string, opts Options) Result {
	keySet := toSet(opts.Keys)
	out := copyMap(env)
	var dropped []string

	for k := range env {
		if shouldDrop(k, keySet, opts.Prefixes, opts.Suffixes) {
			dropped = append(dropped, k)
			if !opts.DryRun {
				delete(out, k)
			}
		}
	}
	sort.Strings(dropped)
	return Result{Out: out, Dropped: dropped}
}

func shouldDrop(key string, keySet map[string]bool, prefixes, suffixes []string) bool {
	if keySet[key] {
		return true
	}
	for _, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	for _, s := range suffixes {
		if len(key) >= len(s) && key[len(key)-len(s):] == s {
			return true
		}
	}
	return false
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
