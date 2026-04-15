package filter_test

import (
	"os"
	"path/filepath"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/filter"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestFilterIntegration_ParseAndFilter(t *testing.T) {
	path := writeTempEnv(t, `
DB_HOST=localhost
DB_PORT=5432
APP_NAME=myapp
APP_ENV=production
SECRET_KEY=abc123
`)
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	res, err := filter.Filter(env, filter.Options{Patterns: []string{"DB_*"}})
	if err != nil {
		t.Fatalf("filter error: %v", err)
	}
	if res.Matched != 2 {
		t.Errorf("expected 2 DB_ keys, got %d", res.Matched)
	}
	if res.Env["DB_HOST"] != "localhost" {
		t.Errorf("unexpected DB_HOST value: %s", res.Env["DB_HOST"])
	}
}

func TestFilterIntegration_RegexFromFile(t *testing.T) {
	path := writeTempEnv(t, `
AWS_ACCESS_KEY=AKIA123
AWS_SECRET=topsecret
GCP_PROJECT=my-project
AZURE_TENANT=tenant-id
`)
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	res, err := filter.Filter(env, filter.Options{Regex: "^AWS_"})
	if err != nil {
		t.Fatalf("filter error: %v", err)
	}
	if res.Matched != 2 {
		t.Errorf("expected 2 AWS_ keys, got %d", res.Matched)
	}
	if res.Dropped != 2 {
		t.Errorf("expected 2 dropped, got %d", res.Dropped)
	}
}
