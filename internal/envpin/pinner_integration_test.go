package envpin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envoy-sync/internal/envfile"
	"github.com/yourusername/envoy-sync/internal/envpin"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPinIntegration_SaveAndCheck(t *testing.T) {
	envPath := writeTempEnv(t, "DB_HOST=prod-db\nPORT=5432\nSECRET=abc\n")
	env, err := envfile.Parse(envPath)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	pinPath := filepath.Join(t.TempDir(), "pins.json")
	if err := envpin.SavePins(env, []string{"DB_HOST", "PORT"}, pinPath); err != nil {
		t.Fatalf("SavePins: %v", err)
	}
	pf, err := envpin.LoadPins(pinPath)
	if err != nil {
		t.Fatalf("LoadPins: %v", err)
	}
	results := envpin.Check(pf, env)
	if envpin.HasDrift(results) {
		t.Error("expected no drift on fresh save")
	}
}

func TestPinIntegration_DriftDetected(t *testing.T) {
	envPath := writeTempEnv(t, "DB_HOST=prod-db\nPORT=5432\n")
	env, _ := envfile.Parse(envPath)
	pinPath := filepath.Join(t.TempDir(), "pins.json")
	_ = envpin.SavePins(env, []string{"DB_HOST", "PORT"}, pinPath)

	// simulate env change
	env["DB_HOST"] = "staging-db"

	pf, _ := envpin.LoadPins(pinPath)
	results := envpin.Check(pf, env)
	if !envpin.HasDrift(results) {
		t.Error("expected drift after env change")
	}
	var drifted []string
	for _, r := range results {
		if r.Drifted {
			drifted = append(drifted, r.Key)
		}
	}
	if len(drifted) != 1 || drifted[0] != "DB_HOST" {
		t.Errorf("unexpected drifted keys: %v", drifted)
	}
}
