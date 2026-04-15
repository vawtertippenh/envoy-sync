package prefix

import "sort"

// Options controls how prefix addition or removal is applied.
type Options struct {
	Prefix  string
	Strip   bool // if true, remove the prefix instead of adding it
	SkipMissing bool // if true, skip keys that don't have the prefix when stripping
}

// Result holds the transformed env map and metadata.
type Result struct {
	Env     map[string]string
	Changed int
	Skipped int
}

// Apply adds or removes a prefix from all keys in the env map.
func Apply(env map[string]string, opts Options) Result {
	result := Result{
		Env: make(map[string]string, len(env)),
	}

	for _, k := range sortedKeys(env) {
		v := env[k]
		if opts.Strip {
			if len(k) > len(opts.Prefix) && k[:len(opts.Prefix)] == opts.Prefix {
				newKey := k[len(opts.Prefix):]
				result.Env[newKey] = v
				result.Changed++
			} else {
				if opts.SkipMissing {
					result.Skipped++
					continue
				}
				result.Env[k] = v
			}
		} else {
			newKey := opts.Prefix + k
			result.Env[newKey] = v
			result.Changed++
		}
	}

	return result
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
