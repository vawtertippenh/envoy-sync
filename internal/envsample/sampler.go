// Package envsampler provides functionality to sample a subset of keys
// from an env map, useful for generating representative examples or test fixtures.
package envsample

import (
	"fmt"
	"math/rand"
	"sort"
)

// Options controls how sampling is performed.
type Options struct {
	// N is the number of keys to sample. If 0 or >= len(env), all keys are returned.
	N int
	// Seed is used for deterministic sampling. If 0, sampling is non-deterministic.
	Seed int64
	// Deterministic forces use of the Seed value even when Seed is 0.
	Deterministic bool
	// IncludeKeys forces specific keys to always be included in the sample.
	IncludeKeys []string
}

// Result holds the sampled env map and metadata.
type Result struct {
	Env      map[string]string
	Total    int
	Sampled  int
	Forced   int
}

// Sample returns a subset of env keys according to the given options.
func Sample(env map[string]string, opts Options) (Result, error) {
	if len(env) == 0 {
		return Result{Env: map[string]string{}}, nil
	}

	forced := toSet(opts.IncludeKeys)
	for k := range forced {
		if _, ok := env[k]; !ok {
			return Result{}, fmt.Errorf("envsample: forced key %q not found in env", k)
		}
	}

	all := sortedKeys(env)
	total := len(all)

	n := opts.N
	if n <= 0 || n >= total {
		return Result{
			Env:     copyMap(env),
			Total:   total,
			Sampled: total,
			Forced:  len(forced),
		}, nil
	}

	var rng *rand.Rand
	if opts.Deterministic || opts.Seed != 0 {
		rng = rand.New(rand.NewSource(opts.Seed)) //nolint:gosec
	} else {
		rng = rand.New(rand.NewSource(rand.Int63())) //nolint:gosec
	}

	// Separate forced keys from candidates
	candidates := make([]string, 0, len(all))
	for _, k := range all {
		if !forced[k] {
			candidates = append(candidates, k)
		}
	}

	remaining := n - len(forced)
	if remaining < 0 {
		remaining = 0
	}
	if remaining > len(candidates) {
		remaining = len(candidates)
	}

	rng.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	picked := candidates[:remaining]

	out := map[string]string{}
	for k := range forced {
		out[k] = env[k]
	}
	for _, k := range picked {
		out[k] = env[k]
	}

	return Result{
		Env:     out,
		Total:   total,
		Sampled: len(out),
		Forced:  len(forced),
	}, nil
}

// Summary returns a human-readable summary of a Result.
func Summary(r Result) string {
	return fmt.Sprintf("sampled %d/%d keys (%d forced)", r.Sampled, r.Total, r.Forced)
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
