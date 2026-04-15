package strip

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "myapp",
		"APP_SECRET":  "s3cr3t",
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "pass",
		"DEBUG":       "true",
	}
}

func TestStrip_NoOptions(t *testing.T) {
	res := Strip(baseEnv(), Options{})
	if len(res.Env) != 5 {
		t.Fatalf("expected 5 keys, got %d", len(res.Env))
	}
	if len(res.Removed) != 0 {
		t.Fatalf("expected no removals, got %v", res.Removed)
	}
}

func TestStrip_ExplicitKeys(t *testing.T) {
	res := Strip(baseEnv(), Options{Keys: []string{"DEBUG", "APP_SECRET"}})
	if _, ok := res.Env["DEBUG"]; ok {
		t.Error("DEBUG should have been removed")
	}
	if _, ok := res.Env["APP_SECRET"]; ok {
		t.Error("APP_SECRET should have been removed")
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestStrip_ByPrefix(t *testing.T) {
	res := Strip(baseEnv(), Options{RemovePrefixes: []string{"DB_"}})
	if _, ok := res.Env["DB_HOST"]; ok {
		t.Error("DB_HOST should have been removed")
	}
	if _, ok := res.Env["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been removed")
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
	if _, ok := res.Env["APP_NAME"]; !ok {
		t.Error("APP_NAME should remain")
	}
}

func TestStrip_BySuffix(t *testing.T) {
	res := Strip(baseEnv(), Options{RemoveSuffixes: []string{"_SECRET", "_PASSWORD"}})
	if _, ok := res.Env["APP_SECRET"]; ok {
		t.Error("APP_SECRET should have been removed")
	}
	if _, ok := res.Env["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been removed")
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestStrip_RemovedSorted(t *testing.T) {
	res := Strip(baseEnv(), Options{RemovePrefixes: []string{"APP_", "DB_"}})
	for i := 1; i < len(res.Removed); i++ {
		if res.Removed[i-1] > res.Removed[i] {
			t.Fatalf("removed list not sorted: %v", res.Removed)
		}
	}
}

func TestStrip_OriginalUnmutated(t *testing.T) {
	env := baseEnv()
	Strip(env, Options{Keys: []string{"DEBUG"}})
	if _, ok := env["DEBUG"]; !ok {
		t.Error("original map should not be mutated")
	}
}
