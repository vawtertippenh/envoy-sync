package scope_test

import (
	"os"
	"path/filepath"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/scope"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestScopeIntegration_ParseAndScope(t *testing.T) {
	p := writeTempEnv(t, `
APP_HOST=localhost
APP_PORT=8080
DB_HOST=db.local
DB_PASS=secret
LOG_LEVEL=debug
`)
	env, err := envfile.Parse(p)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	r := scope.Scope(env, scope.Options{Prefixes: []string{"APP_"}})
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
	if r.Matched["APP_HOST"] != "localhost" {
		t.Errorf("unexpected value for APP_HOST: %s", r.Matched["APP_HOST"])
	}
}

func TestScopeIntegration_StripAndRoundtrip(t *testing.T) {
	p := writeTempEnv(t, `
SVC_NAME=myservice
SVC_VERSION=1.2.3
OTHER_KEY=ignore
`)
	env, err := envfile.Parse(p)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	r := scope.Scope(env, scope.Options{
		Prefixes: []string{"SVC_"},
		Strip:    true,
	})

	if r.Matched["NAME"] != "myservice" {
		t.Errorf("expected NAME=myservice, got %s", r.Matched["NAME"])
	}
	if r.Matched["VERSION"] != "1.2.3" {
		t.Errorf("expected VERSION=1.2.3, got %s", r.Matched["VERSION"])
	}
	if len(r.Unmatched) != 1 {
		t.Errorf("expected 1 unmatched, got %d", len(r.Unmatched))
	}
}
