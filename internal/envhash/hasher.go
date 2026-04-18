package envhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// HashResult holds the overall hash and per-key hashes.
type HashResult struct {
	Overall string
	Keys    map[string]string
}

// Hash computes a deterministic SHA-256 hash of the env map.
// If includeKeys is non-empty, only those keys are hashed.
func Hash(env map[string]string, includeKeys []string) HashResult {
	keys := selectKeys(env, includeKeys)
	sort.Strings(keys)

	perKey := make(map[string]string, len(keys))
	h := sha256.New()

	for _, k := range keys {
		v := env[k]
		entry := fmt.Sprintf("%s=%s\n", k, v)
		h.Write([]byte(entry))
		kh := sha256.Sum256([]byte(entry))
		perKey[k] = hex.EncodeToString(kh[:])
	}

	return HashResult{
		Overall: hex.EncodeToString(h.Sum(nil)),
		Keys:    perKey,
	}
}

// Equal returns true when two HashResults share the same overall hash.
func Equal(a, b HashResult) bool {
	return strings.EqualFold(a.Overall, b.Overall)
}

func selectKeys(env map[string]string, include []string) []string {
	if len(include) > 0 {
		set := make(map[string]struct{}, len(include))
		for _, k := range include {
			set[k] = struct{}{}
		}
		var out []string
		for k := range env {
			if _, ok := set[k]; ok {
				out = append(out, k)
			}
		}
		return out
	}
	out := make([]string, 0, len(env))
	for k := range env {
		out = append(out, k)
	}
	return out
}
