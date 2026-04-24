package envbatch

import (
	"fmt"
	"sort"
)

// Batch holds a named slice of env key-value pairs.
type Batch struct {
	Name  string
	Items map[string]string
}

// Options controls how batching is performed.
type Options struct {
	Size      int  // max keys per batch
	SortKeys  bool // sort keys before batching
}

// Batch splits an env map into fixed-size chunks.
func BatchEnv(env map[string]string, opts Options) ([]Batch, error) {
	if opts.Size <= 0 {
		return nil, fmt.Errorf("batch size must be greater than zero")
	}

	keys := sortedKeys(env)
	if !opts.SortKeys {
		// use insertion order approximation via sorted for determinism
		keys = sortedKeys(env)
	}

	var batches []Batch
	for i := 0; i < len(keys); i += opts.Size {
		end := i + opts.Size
		if end > len(keys) {
			end = len(keys)
		}
		chunk := keys[i:end]
		items := make(map[string]string, len(chunk))
		for _, k := range chunk {
			items[k] = env[k]
		}
		batches = append(batches, Batch{
			Name:  fmt.Sprintf("batch_%d", len(batches)+1),
			Items: items,
		})
	}
	return batches, nil
}

// Summary returns a human-readable description of the batches.
func Summary(batches []Batch) string {
	if len(batches) == 0 {
		return "no batches produced"
	}
	total := 0
	for _, b := range batches {
		total += len(b.Items)
	}
	return fmt.Sprintf("%d batch(es), %d key(s) total", len(batches), total)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
