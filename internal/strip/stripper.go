package strip

import "sort"

// Options controls which entries are stripped from the env map.
type Options struct {
	// RemoveComments strips keys that start with '#' (shouldn't normally exist
	// after parsing, but defensive).
	RemoveComments bool
	// RemovePrefixes removes any key whose prefix matches one of the given strings.
	RemovePrefixes []string
	// RemoveSuffixes removes any key whose suffix matches one of the given strings.
	RemoveSuffixes []string
	// Keys is an explicit list of keys to remove.
	Keys []string
}

// Result holds the stripped env map and metadata.
type Result struct {
	Env     map[string]string
	Removed []string
}

// Strip removes entries from env according to opts and returns a Result.
func Strip(env map[string]string, opts Options) Result {
	explicit := toSet(opts.Keys)
	out := make(map[string]string, len(env))
	var removed []string

	for k, v := range env {
		if shouldRemove(k, opts, explicit) {
			removed = append(removed, k)
			continue
		}
		out[k] = v
	}

	sort.Strings(removed)
	return Result{Env: out, Removed: removed}
}

func shouldRemove(key string, opts Options, explicit map[string]bool) bool {
	if explicit[key] {
		return true
	}
	if opts.RemoveComments && len(key) > 0 && key[0] == '#' {
		return true
	}
	for _, p := range opts.RemovePrefixes {
		if len(p) > 0 && len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	for _, s := range opts.RemoveSuffixes {
		if len(s) > 0 && len(key) >= len(s) && key[len(key)-len(s):] == s {
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
