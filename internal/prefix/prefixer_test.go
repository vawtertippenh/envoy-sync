package prefix

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"NAME": "mydb",
	}
}

func TestApply_AddPrefix(t *testing.T) {
	result := Apply(baseEnv(), Options{Prefix: "DB_"})
	if _, ok := result.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to exist")
	}
	if _, ok := result.Env["HOST"]; ok {
		t.Error("original key HOST should not exist")
	}
	if result.Changed != 3 {
		t.Errorf("expected 3 changed, got %d", result.Changed)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"DB_NAME": "mydb",
	}
	result := Apply(env, Options{Prefix: "DB_", Strip: true})
	if v, ok := result.Env["HOST"]; !ok || v != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", v)
	}
	if result.Changed != 3 {
		t.Errorf("expected 3 changed, got %d", result.Changed)
	}
}

func TestApply_StripPrefix_KeyWithoutPrefix_Kept(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"OTHER":   "value",
	}
	result := Apply(env, Options{Prefix: "DB_", Strip: true, SkipMissing: false})
	if _, ok := result.Env["OTHER"]; !ok {
		t.Error("expected OTHER to be kept when SkipMissing is false")
	}
	if result.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", result.Changed)
	}
}

func TestApply_StripPrefix_KeyWithoutPrefix_Skipped(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"OTHER":   "value",
	}
	result := Apply(env, Options{Prefix: "DB_", Strip: true, SkipMissing: true})
	if _, ok := result.Env["OTHER"]; ok {
		t.Error("expected OTHER to be skipped when SkipMissing is true")
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
}

func TestApply_EmptyPrefix(t *testing.T) {
	result := Apply(baseEnv(), Options{Prefix: ""})
	if len(result.Env) != len(baseEnv()) {
		t.Errorf("expected same number of keys, got %d", len(result.Env))
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	result := Apply(map[string]string{}, Options{Prefix: "APP_"})
	if len(result.Env) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result.Env))
	}
	if result.Changed != 0 {
		t.Errorf("expected 0 changed, got %d", result.Changed)
	}
}
