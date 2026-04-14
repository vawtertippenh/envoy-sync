package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"envoy-sync/internal/profile"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "profiles.json")
}

func TestSet_AndGet(t *testing.T) {
	s := profile.NewStore()
	s.Set("dev", map[string]string{"APP_ENV": "development", "PORT": "3000"})
	p, ok := s.Get("dev")
	if !ok {
		t.Fatal("expected profile 'dev' to exist")
	}
	if p.Env["APP_ENV"] != "development" {
		t.Errorf("expected APP_ENV=development, got %s", p.Env["APP_ENV"])
	}
}

func TestGet_Missing(t *testing.T) {
	s := profile.NewStore()
	_, ok := s.Get("prod")
	if ok {
		t.Fatal("expected missing profile to return false")
	}
}

func TestDelete_Existing(t *testing.T) {
	s := profile.NewStore()
	s.Set("staging", map[string]string{"X": "1"})
	removed := s.Delete("staging")
	if !removed {
		t.Fatal("expected Delete to return true for existing key")
	}
	_, ok := s.Get("staging")
	if ok {
		t.Fatal("expected profile to be gone after delete")
	}
}

func TestDelete_Missing(t *testing.T) {
	s := profile.NewStore()
	removed := s.Delete("ghost")
	if removed {
		t.Fatal("expected Delete to return false for missing key")
	}
}

func TestList_SortedNames(t *testing.T) {
	s := profile.NewStore()
	s.Set("prod", map[string]string{})
	s.Set("dev", map[string]string{})
	s.Set("staging", map[string]string{})
	names := s.List()
	expected := []string{"dev", "prod", "staging"}
	for i, n := range expected {
		if names[i] != n {
			t.Errorf("expected names[%d]=%s, got %s", i, n, names[i])
		}
	}
}

func TestSaveLoad_Roundtrip(t *testing.T) {
	path := tempStorePath(t)
	s := profile.NewStore()
	s.Set("dev", map[string]string{"KEY": "value"})
	if err := profile.SaveStore(path, s); err != nil {
		t.Fatalf("SaveStore: %v", err)
	}
	loaded, err := profile.LoadStore(path)
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	p, ok := loaded.Get("dev")
	if !ok {
		t.Fatal("expected 'dev' profile after load")
	}
	if p.Env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", p.Env["KEY"])
	}
}

func TestLoadStore_MissingFile(t *testing.T) {
	s, err := profile.LoadStore("/nonexistent/path/profiles.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.List()) != 0 {
		t.Fatal("expected empty store for missing file")
	}
}

func TestLoadStore_InvalidJSON(t *testing.T) {
	path := tempStorePath(t)
	_ = os.WriteFile(path, []byte("not-json{"), 0o644)
	_, err := profile.LoadStore(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
