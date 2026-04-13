package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "envoy",
		"PORT":     "8080",
		"DEBUG":    "false",
	}
}

func TestTake_CopiesEnv(t *testing.T) {
	env := baseEnv()
	s := Take("test", env)
	env["PORT"] = "9999" // mutate original
	if s.Env["PORT"] != "8080" {
		t.Errorf("expected snapshot to be independent, got %s", s.Env["PORT"])
	}
	if s.Label != "test" {
		t.Errorf("expected label 'test', got %s", s.Label)
	}
}

func TestSaveLoad_Roundtrip(t *testing.T) {
	s := Take("roundtrip", baseEnv())
	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := Save(s, tmp); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Label != s.Label {
		t.Errorf("label mismatch: got %s", loaded.Label)
	}
	if loaded.Env["APP_NAME"] != "envoy" {
		t.Errorf("env mismatch: got %s", loaded.Env["APP_NAME"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(tmp, []byte("not json"), 0o644)
	_, err := Load(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCompare_NoDrift(t *testing.T) {
	a := Take("a", baseEnv())
	b := Take("b", baseEnv())
	d := Compare(a, b)
	if HasDrift(d) {
		t.Error("expected no drift")
	}
}

func TestCompare_Added(t *testing.T) {
	a := Take("a", baseEnv())
	newEnv := baseEnv()
	newEnv["NEW_KEY"] = "value"
	b := Take("b", newEnv)
	d := Compare(a, b)
	if _, ok := d.Added["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY in Added")
	}
}

func TestCompare_Removed(t *testing.T) {
	a := Take("a", baseEnv())
	newEnv := baseEnv()
	delete(newEnv, "DEBUG")
	b := Take("b", newEnv)
	d := Compare(a, b)
	if _, ok := d.Removed["DEBUG"]; !ok {
		t.Error("expected DEBUG in Removed")
	}
}

func TestCompare_Changed(t *testing.T) {
	a := Take("a", baseEnv())
	newEnv := baseEnv()
	newEnv["PORT"] = "9090"
	b := Take("b", newEnv)
	d := Compare(a, b)
	ch, ok := d.Changed["PORT"]
	if !ok {
		t.Fatal("expected PORT in Changed")
	}
	if ch.Before != "8080" || ch.After != "9090" {
		t.Errorf("unexpected change: %+v", ch)
	}
}
