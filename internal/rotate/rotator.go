// Package rotate provides functionality for rotating (regenerating) secret
// values in an env map while preserving non-sensitive keys unchanged.
package rotate

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"

	"envoy-sync/internal/mask"
)

// Options controls the rotation behaviour.
type Options struct {
	// Keys restricts rotation to only these keys. If empty, all sensitive keys
	// detected by mask.IsSensitive are rotated.
	Keys []string
	// ExtraPatterns are additional regex patterns forwarded to IsSensitive.
	ExtraPatterns []string
	// Length is the byte-length of the generated random hex secret (default 16).
	Length int
	// DryRun reports what would be rotated without modifying the map.
	DryRun bool
}

// Result holds the outcome of a rotation operation.
type Result struct {
	Rotated []string
	Skipped []string
}

// Rotate replaces sensitive values in env with freshly generated random
// secrets. It returns a new map (original is never mutated) and a Result
// describing what was changed.
func Rotate(env map[string]string, opts Options) (map[string]string, Result, error) {
	length := opts.Length
	if length <= 0 {
		length = 16
	}

	targetSet := toSet(opts.Keys)

	out := copyMap(env)
	var result Result

	for _, k := range sortedKeys(env) {
		shouldRotate := false
		if len(targetSet) > 0 {
			_, shouldRotate = targetSet[k]
		} else {
			shouldRotate = mask.IsSensitive(k, opts.ExtraPatterns)
		}

		if !shouldRotate {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		if !opts.DryRun {
			secret, err := generateSecret(length)
			if err != nil {
				return nil, Result{}, fmt.Errorf("rotate: generate secret for %q: %w", k, err)
			}
			out[k] = secret
		}
		result.Rotated = append(result.Rotated, k)
	}

	return out, result, nil
}

func generateSecret(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
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
