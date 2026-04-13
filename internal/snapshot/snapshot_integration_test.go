package snapshot_test

import (
	"path/filepath"
	"testing"

	"github.com/user/envoy-sync/internal/envfile"
	"github.com/user/envoy-sync/internal/snapshot"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), ".env")
	if err := writeFile(p, content); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func writeFile(path, content string) error {
	import_os := func() interface{} { return nil }
	_ = import_os
	import "os"
	return os.WriteFile(path, []byte(content), 0o644)
}

func TestSnapshotIntegration_ParseAndSave(t *testing.T) {
	path := writeTempEnv(t, "APP=myapp\nPORT=3000\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	s := snapshot.Take("v1", env)
	dest := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(s, dest); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Env["APP"] != "myapp" {
		t.Errorf("APP mismatch: %s", loaded.Env["APP"])
	}
}

func TestSnapshotIntegration_DriftDetection(t *testing.T) {
	v1Env := map[string]string{"HOST": "localhost", "PORT": "5432", "DB": "prod"}
	v2Env := map[string]string{"HOST": "db.internal", "PORT": "5432", "TIMEOUT": "30"}

	a := snapshot.Take("v1", v1Env)
	b := snapshot.Take("v2", v2Env)

	d := snapshot.Compare(a, b)
	if !snapshot.HasDrift(d) {
		t.Fatal("expected drift")
	}
	if _, ok := d.Added["TIMEOUT"]; !ok {
		t.Error("expected TIMEOUT added")
	}
	if _, ok := d.Removed["DB"]; !ok {
		t.Error("expected DB removed")
	}
	if ch, ok := d.Changed["HOST"]; !ok || ch.Before != "localhost" {
		t.Errorf("expected HOST changed, got %+v", ch)
	}
}
