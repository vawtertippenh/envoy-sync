package envreplace

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"API_URL":  "http://localhost:8080",
		"APP_NAME": "my-app",
	}
}

func TestReplace_NoRules(t *testing.T) {
	env := baseEnv()
	res, err := Replace(env, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.ChangedKeys) != 0 {
		t.Errorf("expected no changes, got %v", res.ChangedKeys)
	}
	if res.Env["DB_HOST"] != "localhost" {
		t.Errorf("value should be unchanged")
	}
}

func TestReplace_SimpleString(t *testing.T) {
	env := baseEnv()
	rules := []Rule{{Find: "localhost", Replace: "prod.example.com"}}
	res, err := Replace(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["DB_HOST"] != "prod.example.com" {
		t.Errorf("expected DB_HOST to be updated, got %q", res.Env["DB_HOST"])
	}
	if res.Env["API_URL"] != "http://prod.example.com:8080" {
		t.Errorf("expected API_URL to be updated, got %q", res.Env["API_URL"])
	}
	if len(res.ChangedKeys) != 2 {
		t.Errorf("expected 2 changed keys, got %d: %v", len(res.ChangedKeys), res.ChangedKeys)
	}
}

func TestReplace_OriginalUnmutated(t *testing.T) {
	env := baseEnv()
	rules := []Rule{{Find: "localhost", Replace: "remote"}}
	_, err := Replace(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DB_HOST"] != "localhost" {
		t.Errorf("original map was mutated")
	}
}

func TestReplace_RegexRule(t *testing.T) {
	env := map[string]string{
		"VERSION": "v1.2.3",
		"TAG":     "release-v2.0.0",
	}
	rules := []Rule{{Find: `v\d+\.\d+\.\d+`, Replace: "vX.Y.Z", Regex: true}}
	res, err := Replace(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["VERSION"] != "vX.Y.Z" {
		t.Errorf("expected VERSION replaced, got %q", res.Env["VERSION"])
	}
	if res.Env["TAG"] != "release-vX.Y.Z" {
		t.Errorf("expected TAG replaced, got %q", res.Env["TAG"])
	}
}

func TestReplace_InvalidRegex(t *testing.T) {
	env := baseEnv()
	rules := []Rule{{Find: "[invalid", Replace: "x", Regex: true}}
	_, err := Replace(env, rules)
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestReplace_EmptyFindError(t *testing.T) {
	env := baseEnv()
	rules := []Rule{{Find: "", Replace: "something"}}
	_, err := Replace(env, rules)
	if err == nil {
		t.Error("expected error for empty Find, got nil")
	}
}

func TestReplace_MultipleRules(t *testing.T) {
	env := map[string]string{"URL": "http://old-host:9000/api"}
	rules := []Rule{
		{Find: "old-host", Replace: "new-host"},
		{Find: "9000", Replace: "443"},
		{Find: "http", Replace: "https"},
	}
	res, err := Replace(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["URL"] != "https://new-host:443/api" {
		t.Errorf("expected full replacement, got %q", res.Env["URL"])
	}
}
