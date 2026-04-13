package rename

import "fmt"

// Result represents the outcome of a single rename operation.
type Result struct {
	OldKey  string
	NewKey  string
	Value   string
	Skipped bool
	Reason  string
}

// Options controls rename behaviour.
type Options struct {
	// Overwrite allows renaming into an existing key, replacing its value.
	Overwrite bool
}

// Rename renames oldKey to newKey in env, returning a new map and a Result.
// The original map is not mutated.
func Rename(env map[string]string, oldKey, newKey string, opts Options) (map[string]string, Result) {
	out := copyMap(env)

	val, exists := out[oldKey]
	if !exists {
		return out, Result{
			OldKey:  oldKey,
			NewKey:  newKey,
			Skipped: true,
			Reason:  fmt.Sprintf("key %q not found", oldKey),
		}
	}

	if _, conflict := out[newKey]; conflict && !opts.Overwrite {
		return out, Result{
			OldKey:  oldKey,
			NewKey:  newKey,
			Value:   val,
			Skipped: true,
			Reason:  fmt.Sprintf("key %q already exists; use --overwrite to replace", newKey),
		}
	}

	delete(out, oldKey)
	out[newKey] = val

	return out, Result{
		OldKey: oldKey,
		NewKey: newKey,
		Value:  val,
	}
}

// RenameMany applies multiple renames sequentially.
func RenameMany(env map[string]string, pairs [][2]string, opts Options) (map[string]string, []Result) {
	current := copyMap(env)
	results := make([]Result, 0, len(pairs))
	for _, p := range pairs {
		var r Result
		current, r = Rename(current, p[0], p[1], opts)
		results = append(results, r)
	}
	return current, results
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
