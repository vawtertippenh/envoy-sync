package strip_test

import (
	"os"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/strip"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "strip-*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestStripIntegration_ParseAndStrip(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nAPP_SECRET=s3cr3t\nDB_HOST=localhost\nDEBUG=true\n")

	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	res := strip.Strip(env, strip.Options{
		Keys:           []string{"DEBUG"},
		RemoveSuffixes: []string{"_SECRET"},
	})

	if _, ok := res.Env["DEBUG"]; ok {
		t.Error("DEBUG should be stripped")
	}
	if _, ok := res.Env["APP_SECRET"]; ok {
		t.Error("APP_SECRET should be stripped")
	}
	if res.Env["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should remain, got %q", res.Env["APP_NAME"])
	}
	if res.Env["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST should remain, got %q", res.Env["DB_HOST"])
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %v", res.Removed)
	}
}

func TestStripIntegration_PrefixStrip(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=prod\n")

	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	res := strip.Strip(env, strip.Options{RemovePrefixes: []string{"DB_"}})

	if len(res.Env) != 1 {
		t.Fatalf("expected 1 key remaining, got %d", len(res.Env))
	}
	if res.Env["APP_ENV"] != "prod" {
		t.Errorf("APP_ENV should remain")
	}
}
