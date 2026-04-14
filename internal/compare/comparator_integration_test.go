package compare_test

import (
	"os"
	"path/filepath"
	"testing"

	"envoy-sync/internal/compare"
	"envoy-sync/internal/envfile"
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

func TestCompareIntegration_TemplateVsTarget(t *testing.T) {
	tmplPath := writeTempEnv(t, "HOST=localhost\nPORT=8080\nDEBUG=true\n")
	targetPath := writeTempEnv(t, "HOST=prod.example.com\nPORT=443\nEXTRA=yes\n")

	tmpl, err := envfile.Parse(tmplPath)
	if err != nil {
		t.Fatal(err)
	}
	target, err := envfile.Parse(targetPath)
	if err != nil {
		t.Fatal(err)
	}

	r := compare.Against(tmpl, target)

	if len(r.Missing) != 1 || r.Missing[0] != "DEBUG" {
		t.Errorf("expected DEBUG missing, got %v", r.Missing)
	}
	if len(r.Extra) != 1 || r.Extra[0] != "EXTRA" {
		t.Errorf("expected EXTRA extra, got %v", r.Extra)
	}
	if len(r.Mismatch) != 0 {
		t.Errorf("expected no mismatch, got %v", r.Mismatch)
	}
}

func TestCompareIntegration_PerfectMatch(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nLOG_LEVEL=info\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	r := compare.Against(env, env)
	if len(r.Missing) != 0 || len(r.Extra) != 0 || len(r.Mismatch) != 0 {
		t.Errorf("expected perfect match, got %+v", r)
	}
}

func TestCompareIntegration_EmptyTarget(t *testing.T) {
	tmplPath := writeTempEnv(t, "HOST=localhost\nPORT=8080\n")
	targetPath := writeTempEnv(t, "")

	tmpl, err := envfile.Parse(tmplPath)
	if err != nil {
		t.Fatal(err)
	}
	target, err := envfile.Parse(targetPath)
	if err != nil {
		t.Fatal(err)
	}

	r := compare.Against(tmpl, target)

	if len(r.Missing) != 2 {
		t.Errorf("expected 2 missing keys, got %v", r.Missing)
	}
	if len(r.Extra) != 0 {
		t.Errorf("expected no extra keys, got %v", r.Extra)
	}
}
