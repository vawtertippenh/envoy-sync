package envsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
)

// Signature holds the HMAC-SHA256 signature and the list of signed keys.
type Signature struct {
	Digest string   `json:"digest"`
	Keys   []string `json:"keys"`
}

// Sign produces an HMAC-SHA256 signature over the canonical form of env.
// If keys is non-empty only those keys are included; otherwise all keys are signed.
func Sign(env map[string]string, secret string, keys []string) (Signature, error) {
	if secret == "" {
		return Signature{}, errors.New("envsign: secret must not be empty")
	}

	target := keys
	if len(target) == 0 {
		for k := range env {
			target = append(target, k)
		}
	}
	sort.Strings(target)

	h := hmac.New(sha256.New, []byte(secret))
	for _, k := range target {
		v, ok := env[k]
		if !ok {
			return Signature{}, fmt.Errorf("envsign: key %q not found in env", k)
		}
		fmt.Fprintf(h, "%s=%s\n", k, v)
	}

	return Signature{
		Digest: hex.EncodeToString(h.Sum(nil)),
		Keys:   target,
	}, nil
}

// Verify checks that sig matches the env using secret.
func Verify(env map[string]string, secret string, sig Signature) (bool, error) {
	got, err := Sign(env, secret, sig.Keys)
	if err != nil {
		return false, err
	}
	return hmac.Equal([]byte(got.Digest), []byte(sig.Digest)), nil
}
