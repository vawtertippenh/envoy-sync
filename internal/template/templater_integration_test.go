package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-sync/internal/envfile"
	"github.com/user/envoy-sync/internal/template"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestTemplateIntegration_ParseAndFill(t *testing.T) {
	tmplPath := writeTempEnv(t, "DB_HOST=<required>\nDB_PORT=5432\nAPP_ENV=<optional>\n")
	valPath := writeTempEnv(t, "DB_HOST=db.prod.internal\nAPP_ENV=production\n")

	tmplEnv, err := envfile.Parse(tmplPath)
	if err != nil {
		t.Fatalf("parse template: %v", err)
	}
	valEnv, err := envfile.Parse(valPath)
	if err != nil {
		t.Fatalf("parse values: %v", err)
	}

	r := template.Fill(tmplEnv, valEnv)

	if len(r.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", r.Missing)
	}
	if r.Filled["DB_HOST"] != "db.prod.internal" {
		t.Errorf("unexpected DB_HOST: %s", r.Filled["DB_HOST"])
	}
	if r.Filled["DB_PORT"] != "5432" {
		t.Errorf("expected default DB_PORT=5432, got %s", r.Filled["DB_PORT"])
	}
}

func TestTemplateIntegration_MissingRequired(t *testing.T) {
	tmplPath := writeTempEnv(t, "API_KEY=<required>\nRETRIES=3\n")
	valPath := writeTempEnv(t, "RETRIES=5\n")

	tmplEnv, _ := envfile.Parse(tmplPath)
	valEnv, _ := envfile.Parse(valPath)

	r := template.Fill(tmplEnv, valEnv)

	if len(r.Missing) != 1 || r.Missing[0] != "API_KEY" {
		t.Errorf("expected API_KEY in missing, got %v", r.Missing)
	}

	summary := template.Summary(r)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
