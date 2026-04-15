package defaults_test

import (
	"os"
	"path/filepath"
	"testing"

	"envoy-sync/internal/defaults"
	"envoy-sync/internal/envfile"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestDefaultsIntegration_ParseAndApply(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nLOG_LEVEL=\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	rules := []defaults.Rule{
		{Key: "LOG_LEVEL", Value: "warn"},
		{Key: "PORT", Value: "3000"},
		{Key: "APP_NAME", Value: "fallback"},
	}

	res := defaults.Apply(env, rules)

	if res.Env["LOG_LEVEL"] != "warn" {
		t.Errorf("expected LOG_LEVEL=warn, got %q", res.Env["LOG_LEVEL"])
	}
	if res.Env["PORT"] != "3000" {
		t.Errorf("expected PORT=3000, got %q", res.Env["PORT"])
	}
	if res.Env["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp (unchanged), got %q", res.Env["APP_NAME"])
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d: %v", len(res.Applied), res.Applied)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "APP_NAME" {
		t.Errorf("expected skipped=[APP_NAME], got %v", res.Skipped)
	}
}

func TestDefaultsIntegration_OverrideAll(t *testing.T) {
	path := writeTempEnv(t, "TIMEOUT=30\nRETRIES=3\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	rules := []defaults.Rule{
		{Key: "TIMEOUT", Value: "60", Override: true},
		{Key: "RETRIES", Value: "5", Override: true},
	}

	res := defaults.Apply(env, rules)

	if res.Env["TIMEOUT"] != "60" {
		t.Errorf("expected TIMEOUT=60, got %q", res.Env["TIMEOUT"])
	}
	if res.Env["RETRIES"] != "5" {
		t.Errorf("expected RETRIES=5, got %q", res.Env["RETRIES"])
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(res.Applied))
	}
}
