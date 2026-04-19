package envclean_test

import (
	"os"
	"testing"

	"envoy-sync/internal/envclean"
	"envoy-sync/internal/envfile"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envclean-*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestCleanIntegration_ParseAndClean(t *testing.T) {
	path := writeTempEnv(t, "APP=hello\nEMPTY=\nDB_HOST=  db  \n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	r := envclean.Clean(env, envclean.Options{
		RemoveEmpty:    true,
		TrimWhitespace: true,
	})
	if _, ok := r.Env["EMPTY"]; ok {
		t.Error("expected EMPTY removed")
	}
	if r.Env["DB_HOST"] != "db" {
		t.Errorf("expected trimmed DB_HOST, got %q", r.Env["DB_HOST"])
	}
	if r.Env["APP"] != "hello" {
		t.Errorf("expected APP=hello, got %q", r.Env["APP"])
	}
}

func TestCleanIntegration_AllClean(t *testing.T) {
	path := writeTempEnv(t, "KEY=value\nOTHER=data\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	r := envclean.Clean(env, envclean.Options{Rem)
	if len(r.Removed) != 0 {
		t.Errorf("expected nothing removed, got %v", r.Removed)
	}
	if len(r.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Env))
	}
}
