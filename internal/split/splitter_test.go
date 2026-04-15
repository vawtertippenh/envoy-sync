package split

import (
	"testing"
)

var baseEnv = map[string]string{
	"APP_HOST":   "localhost",
	"APP_PORT":   "8080",
	"DB_HOST":    "db.local",
	"DB_PASS":    "secret",
	"LOG_LEVEL":  "info",
	"UNTAGGED":   "value",
}

func TestSplit_BasicGroups(t *testing.T) {
	opts := Options{
		Prefixes: map[string]string{
			"app": "APP_",
			"db":  "DB_",
		},
	}
	groups, err := Split(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	appGroup := findGroup(groups, "app")
	if appGroup == nil {
		t.Fatal("expected group 'app'")
	}
	if appGroup.Env["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost")
	}
}

func TestSplit_StripPrefix(t *testing.T) {
	opts := Options{
		StripPrefix: true,
		Prefixes:    map[string]string{"db": "DB_"},
	}
	groups, err := Split(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dbGroup := findGroup(groups, "db")
	if dbGroup == nil {
		t.Fatal("expected group 'db'")
	}
	if _, ok := dbGroup.Env["HOST"]; !ok {
		t.Errorf("expected stripped key HOST, got %v", dbGroup.Env)
	}
	if _, ok := dbGroup.Env["DB_HOST"]; ok {
		t.Errorf("expected prefix to be stripped from DB_HOST")
	}
}

func TestSplit_Remainder(t *testing.T) {
	opts := Options{
		Prefixes:  map[string]string{"app": "APP_", "db": "DB_"},
		Remainder: "other",
	}
	groups, err := Split(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	other := findGroup(groups, "other")
	if other == nil {
		t.Fatal("expected remainder group 'other'")
	}
	if _, ok := other.Env["LOG_LEVEL"]; !ok {
		t.Errorf("expected LOG_LEVEL in remainder")
	}
	if _, ok := other.Env["UNTAGGED"]; !ok {
		t.Errorf("expected UNTAGGED in remainder")
	}
}

func TestSplit_NoPrefixesError(t *testing.T) {
	_, err := Split(baseEnv, Options{})
	if err == nil {
		t.Fatal("expected error for empty prefixes")
	}
}

func TestSplit_EmptyGroupWhenNoMatch(t *testing.T) {
	opts := Options{
		Prefixes: map[string]string{"cache": "CACHE_"},
	}
	groups, err := Split(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cache := findGroup(groups, "cache")
	if cache == nil {
		t.Fatal("expected group 'cache'")
	}
	if len(cache.Env) != 0 {
		t.Errorf("expected empty cache group, got %v", cache.Env)
	}
}

func findGroup(groups []Group, name string) *Group {
	for i := range groups {
		if groups[i].Name == name {
			return &groups[i]
		}
	}
	return nil
}
