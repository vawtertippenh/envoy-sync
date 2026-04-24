package envsample

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "myapp",
		"APP_VERSION": "1.0.0",
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
		"LOG_LEVEL":   "info",
		"DEBUG":       "false",
	}
}

func TestSample_EmptyEnv(t *testing.T) {
	r, err := Sample(map[string]string{}, Options{N: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Env) != 0 {
		t.Errorf("expected empty result, got %d keys", len(r.Env))
	}
}

func TestSample_NGreaterThanTotal(t *testing.T) {
	env := baseEnv()
	r, err := Sample(env, Options{N: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Sampled != len(env) {
		t.Errorf("expected all %d keys, got %d", len(env), r.Sampled)
	}
}

func TestSample_NZeroReturnsAll(t *testing.T) {
	env := baseEnv()
	r, err := Sample(env, Options{N: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Sampled != len(env) {
		t.Errorf("expected all keys, got %d", r.Sampled)
	}
}

func TestSample_Deterministic(t *testing.T) {
	env := baseEnv()
	opts := Options{N: 3, Seed: 42, Deterministic: true}

	r1, err := Sample(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := Sample(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(r1.Env) != len(r2.Env) {
		t.Fatalf("expected same sample size, got %d vs %d", len(r1.Env), len(r2.Env))
	}
	for k := range r1.Env {
		if _, ok := r2.Env[k]; !ok {
			t.Errorf("key %q in first sample but not second", k)
		}
	}
}

func TestSample_ForcedKeysAlwaysIncluded(t *testing.T) {
	env := baseEnv()
	opts := Options{N: 2, Seed: 99, Deterministic: true, IncludeKeys: []string{"DB_PASSWORD", "API_KEY"}}

	r, err := Sample(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Env["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD to be in sample")
	}
	if _, ok := r.Env["API_KEY"]; !ok {
		t.Error("expected API_KEY to be in sample")
	}
	if r.Forced != 2 {
		t.Errorf("expected Forced=2, got %d", r.Forced)
	}
}

func TestSample_ForcedKeyMissing(t *testing.T) {
	env := baseEnv()
	opts := Options{N: 3, IncludeKeys: []string{"NONEXISTENT_KEY"}}
	_, err := Sample(env, opts)
	if err == nil {
		t.Error("expected error for missing forced key")
	}
}

func TestSample_OriginalUnmutated(t *testing.T) {
	env := baseEnv()
	orig := copyMap(env)
	_, _ = Sample(env, Options{N: 3, Seed: 1, Deterministic: true})
	for k, v := range orig {
		if env[k] != v {
			t.Errorf("original env mutated at key %q", k)
		}
	}
}

func TestSummary_Format(t *testing.T) {
	r := Result{Total: 10, Sampled: 4, Forced: 1}
	s := Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
	expected := "sampled 4/10 keys (1 forced)"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
