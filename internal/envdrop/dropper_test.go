package envdrop_test

import (
	"testing"

	"github.com/user/envoy-sync/internal/envdrop"
)

var baseEnv = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_PASSWORD": "secret",
	"DB_HOST":     "db",
	"LOG_LEVEL":   "info",
	"DEBUG":       "true",
}

func TestDrop_NoOptions(t *testing.T) {
	res := envdrop.Drop(baseEnv, envdrop.Options{})
	if len(res.Out) != len(baseEnv) {
		t.Errorf("expected %d keys, got %d", len(baseEnv), len(res.Out))
	}
	if len(res.Dropped) != 0 {
		t.Errorf("expected no dropped keys, got %v", res.Dropped)
	}
}

func TestDrop_ExactKeys(t *testing.T) {
	res := envdrop.Drop(baseEnv, envdrop.Options{Keys: []string{"DEBUG", "LOG_LEVEL"}})
	if _, ok := res.Out["DEBUG"]; ok {
		t.Error("DEBUG should have been dropped")
	}
	if len(res.Dropped) != 2 {
		t.Errorf("expected 2 dropped, got %d", len(res.Dropped))
	}
}

func TestDrop_ByPrefix(t *testing.T) {
	res := envdrop.Drop(baseEnv, envdrop.Options{Prefixes: []string{"APP_"}})
	if _, ok := res.Out["APP_HOST"]; ok {
		t.Error("APP_HOST should have been dropped")
	}
	if _, ok := res.Out["APP_PORT"]; ok {
		t.Error("APP_PORT should have been dropped")
	}
	if len(res.Dropped) != 2 {
		t.Errorf("expected 2 dropped, got %d", len(res.Dropped))
	}
}

func TestDrop_BySuffix(t *testing.T) {
	res := envdrop.Drop(baseEnv, envdrop.Options{Suffixes: []string{"_HOST"}})
	if _, ok := res.Out["APP_HOST"]; ok {
		t.Error("APP_HOST should be dropped")
	}
	if _, ok := res.Out["DB_HOST"]; ok {
		t.Error("DB_HOST should be dropped")
	}
	if len(res.Dropped) != 2 {
		t.Errorf("expected 2 dropped, got %d", len(res.Dropped))
	}
}

func TestDrop_DryRun(t *testing.T) {
	res := envdrop.Drop(baseEnv, envdrop.Options{Keys: []string{"DEBUG"}, DryRun: true})
	if _, ok := res.Out["DEBUG"]; !ok {
		t.Error("DryRun: DEBUG should still be in output")
	}
	if len(res.Dropped) != 1 || res.Dropped[0] != "DEBUG" {
		t.Errorf("expected Dropped=[DEBUG], got %v", res.Dropped)
	}
}

func TestDrop_BaseUnmutated(t *testing.T) {
	origLen := len(baseEnv)
	envdrop.Drop(baseEnv, envdrop.Options{Prefixes: []string{"DB_"}})
	if len(baseEnv) != origLen {
		t.Error("original env was mutated")
	}
}
