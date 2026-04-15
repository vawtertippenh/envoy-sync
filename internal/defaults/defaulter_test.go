package defaults

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"LOG_LEVEL": "info",
	}
}

func TestApply_NoRules(t *testing.T) {
	env := baseEnv()
	res := Apply(env, nil)
	if len(res.Applied) != 0 {
		t.Errorf("expected no applied, got %v", res.Applied)
	}
	if res.Env["APP_NAME"] != "myapp" {
		t.Error("existing key should be unchanged")
	}
}

func TestApply_MissingKeyGetsDefault(t *testing.T) {
	env := baseEnv()
	rules := []Rule{{Key: "PORT", Value: "8080"}}
	res := Apply(env, rules)
	if res.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", res.Env["PORT"])
	}
	if len(res.Applied) != 1 || res.Applied[0] != "PORT" {
		t.Errorf("expected applied=[PORT], got %v", res.Applied)
	}
}

func TestApply_ExistingKeySkipped(t *testing.T) {
	env := baseEnv()
	rules := []Rule{{Key: "APP_NAME", Value: "default-name"}}
	res := Apply(env, rules)
	if res.Env["APP_NAME"] != "myapp" {
		t.Error("existing key should not be overwritten without Override")
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "APP_NAME" {
		t.Errorf("expected skipped=[APP_NAME], got %v", res.Skipped)
	}
}

func TestApply_OverrideReplaces(t *testing.T) {
	env := baseEnv()
	rules := []Rule{{Key: "LOG_LEVEL", Value: "debug", Override: true}}
	res := Apply(env, rules)
	if res.Env["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", res.Env["LOG_LEVEL"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected applied=[LOG_LEVEL], got %v", res.Applied)
	}
}

func TestApply_EmptyValueGetsDefault(t *testing.T) {
	env := map[string]string{"DB_HOST": ""}
	rules := []Rule{{Key: "DB_HOST", Value: "localhost"}}
	res := Apply(env, rules)
	if res.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", res.Env["DB_HOST"])
	}
}

func TestApply_OriginalUnmutated(t *testing.T) {
	env := baseEnv()
	rules := []Rule{{Key: "NEW_KEY", Value: "val"}}
	Apply(env, rules)
	if _, ok := env["NEW_KEY"]; ok {
		t.Error("original env map should not be mutated")
	}
}

func TestSummary_NoDefaults(t *testing.T) {
	r := Result{}
	if r.Summary() != "no defaults applied" {
		t.Errorf("unexpected summary: %q", r.Summary())
	}
}

func TestSummary_WithDefaults(t *testing.T) {
	r := Result{Applied: []string{"PORT", "TIMEOUT"}}
	s := r.Summary()
	if s == "" || s == "no defaults applied" {
		t.Errorf("expected non-empty summary with keys, got %q", s)
	}
}
