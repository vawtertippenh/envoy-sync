package envsplit_test

import (
	"os"
	"path/filepath"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envsplit"
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

func TestSplitIntegration_ParseAndSplit(t *testing.T) {
	p := writeTempEnv(t, "APP_HOST=localhost\nAPP_PORT=9000\nDB_URL=postgres://localhost/db\nSECRET=abc\n")
	env, err := envfile.Parse(p)
	if err != nil {
		t.Fatal(err)
	}
	parts, err := envsplit.Split(env, envsplit.Options{
		Prefixes:      []string{"APP_", "DB_"},
		KeepRemainder: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts (APP_, DB_, _remainder), got %d", len(parts))
	}
}

func TestSplitIntegration_StripPrefixRoundtrip(t *testing.T) {
	p := writeTempEnv(t, "SVC_NAME=auth\nSVC_PORT=8081\n")
	env, err := envfile.Parse(p)
	if err != nil {
		t.Fatal(err)
	}
	parts, err := envsplit.Split(env, envsplit.Options{
		Prefixes:    []string{"SVC_"},
		StripPrefix: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if v := parts[0].Env["NAME"]; v != "auth" {
		t.Errorf("expected NAME=auth, got %q", v)
	}
}
