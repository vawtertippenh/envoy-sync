package envhash

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
		"SECRET":   "s3cr3t",
	}
}

func TestHash_Deterministic(t *testing.T) {
	env := baseEnv()
	r1 := Hash(env, nil)
	r2 := Hash(env, nil)
	if r1.Overall != r2.Overall {
		t.Errorf("expected same hash, got %s vs %s", r1.Overall, r2.Overall)
	}
}

func TestHash_DifferentEnv_DifferentHash(t *testing.T) {
	env1 := baseEnv()
	env2 := baseEnv()
	env2["PORT"] = "9090"
	r1 := Hash(env1, nil)
	r2 := Hash(env2, nil)
	if r1.Overall == r2.Overall {
		t.Error("expected different hashes for different envs")
	}
}

func TestHash_PerKeyHashes(t *testing.T) {
	env := baseEnv()
	r := Hash(env, nil)
	for k := range env {
		if _, ok := r.Keys[k]; !ok {
			t.Errorf("missing per-key hash for %s", k)
		}
	}
}

func TestHash_IncludeKeys_Subset(t *testing.T) {
	env := baseEnv()
	rAll := Hash(env, nil)
	rSub := Hash(env, []string{"APP_NAME", "PORT"})
	if rAll.Overall == rSub.Overall {
		t.Error("subset hash should differ from full hash")
	}
	if _, ok := rSub.Keys["SECRET"]; ok {
		t.Error("SECRET should not appear in subset hash")
	}
}

func TestEqual_SameHash(t *testing.T) {
	env := baseEnv()
	r1 := Hash(env, nil)
	r2 := Hash(env, nil)
	if !Equal(r1, r2) {
		t.Error("expected Equal to return true")
	}
}

func TestEqual_DifferentHash(t *testing.T) {
	env1 := baseEnv()
	env2 := map[string]string{"X": "1"}
	if Equal(Hash(env1, nil), Hash(env2, nil)) {
		t.Error("expected Equal to return false")
	}
}
