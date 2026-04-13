package export_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/export"
	"envoy-sync/internal/mask"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestExportIntegration_DotenvRoundtrip(t *testing.T) {
	p := writeTempEnv(t, "APP=hello\nPORT=3000\n")
	env, err := envfile.Parse(p)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := export.Export(env, export.Options{Format: export.FormatDotenv, Sorted: true})
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if !strings.Contains(out, "APP=hello") || !strings.Contains(out, "PORT=3000") {
		t.Errorf("unexpected dotenv output: %s", out)
	}
}

func TestExportIntegration_MaskedJSON(t *testing.T) {
	p := writeTempEnv(t, "API_KEY=supersecret\nAPP=myapp\n")
	env, err := envfile.Parse(p)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	masked := mask.MaskMap(env, nil)
	out, err := export.Export(masked, export.Options{Format: export.FormatJSON, Sorted: true})
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected secret to be masked in output: %s", out)
	}
	if !strings.Contains(out, "APP") {
		t.Errorf("expected APP key in output: %s", out)
	}
}
