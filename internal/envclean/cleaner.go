package envclean

import "sort"

// Options controls what Clean removes or normalizes.
type Options struct {
	RemoveEmpty    bool
	RemoveComments bool
	TrimWhitespace bool
	DeduplicateKeys bool
}

// Result holds the cleaned env map and metadata.
type Result struct {
	Env     map[string]string
	Removed []string
}

// Clean applies the given options to env and returns a cleaned copy.
func Clean(env map[string]string, opts Options) Result {
	out := copyMap(env)
	var removed []string

	seen := map[string]bool{}
	for _, k := range sortedKeys(env) {
		v := env[k]

		if opts.TrimWhitespace {
			v = trimSpace(v)
			out[k] = v
		}

		if opts.RemoveEmpty && v == "" {
			delete(out, k)
			removed = append(removed, k)
			continue
		}

		if opts.DeduplicateKeys {
			if seen[k] {
				delete(out, k)
				removed = append(removed, k)
				continue
			}
			seen[k] = true
		}
	}

	return Result{Env: out, Removed: removed}
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
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
