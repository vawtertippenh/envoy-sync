package envsplit

import (
	"testing"
)

func base() map[string]string {
	return map[string]string{
		"APP_HOST":  "localhost",
		"APP_PORT":  "8080",
		"DB_HOST":   "db",
		"DB_PASS":   "secret",
		"LOG_LEVEL": "info",
	}
}

func TestSplit_NoPrefixesError(t *testing.T) {
	_, err := Split(base(), Options{})
	if err == nil {
		t.Fatal("expected error for empty prefixes")
	}
}

func TestSplit_BasicGroups(t *testing.T) {
	parts, err := Split(base(), Options{Prefixes: []string{"APP_", "DB_"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	if _, ok := parts[0].Env["APP_HOST"]; !ok {
		t.Error("APP_HOST missing from APP_ group")
	}
	if _, ok := parts[1].Env["DB_HOST"]; !ok {
		t.Error("DB_HOST missing from DB_ group")
	}
}

func TestSplit_StripPrefix(t *testing.T) {
	parts, err := Split(base(), Options{Prefixes: []string{"APP_"}, StripPrefix: true})
	if err != nil {
		t.Fatal(err)
	}
	env := parts[0].Env
	if _, ok := env["HOST"]; !ok {
		t.Error("expected HOST after stripping APP_")
	}
	if _, ok := env["APP_HOST"]; ok {
		t.Error("APP_HOST should have been stripped")
	}
}

func TestSplit_Remainder(t *testing.T) {
	parts, err := Split(base(), Options{Prefixes: []string{"APP_", "DB_"}, KeepRemainder: true})
	if err != nil {
		t.Fatal(err)
	}
	last := parts[len(parts)-1]
	if last.Name != "_remainder" {
		t.Fatalf("expected _remainder part, got %s", last.Name)
	}
	if _, ok := last.Env["LOG_LEVEL"]; !ok {
		t.Error("LOG_LEVEL should be in remainder")
	}
}

func TestSplit_NoRemainderDropsUnmatched(t *testing.T) {
	parts, err := Split(base(), Options{Prefixes: []string{"APP_"}, KeepRemainder: false})
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range parts {
		if p.Name == "_remainder" {
			t.Error("remainder part should not be present")
		}
	}
}
