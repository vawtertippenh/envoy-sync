package scope_test

import (
	"testing"

	"envoy-sync/internal/scope"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.local",
		"DB_PASSWORD": "secret",
		"LOG_LEVEL":   "info",
		"LOG_FORMAT":  "json",
	}
}

func TestScope_NoPatterns_AllMatched(t *testing.T) {
	r := scope.Scope(baseEnv(), scope.Options{})
	if len(r.Matched) != 6 {
		t.Errorf("expected 6 matched, got %d", len(r.Matched))
	}
	if len(r.Unmatched) != 0 {
		t.Errorf("expected 0 unmatched, got %d", len(r.Unmatched))
	}
}

func TestScope_ByPrefix(t *testing.T) {
	r := scope.Scope(baseEnv(), scope.Options{Prefixes: []string{"APP_"}})
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
	if _, ok := r.Matched["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in matched")
	}
	if _, ok := r.Matched["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in matched")
	}
}

func TestScope_BySuffix(t *testing.T) {
	r := scope.Scope(baseEnv(), scope.Options{Suffixes: []string{"_HOST"}})
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
}

func TestScope_StripPrefix(t *testing.T) {
	r := scope.Scope(baseEnv(), scope.Options{
		Prefixes: []string{"APP_"},
		Strip:    true,
	})
	if _, ok := r.Matched["HOST"]; !ok {
		t.Error("expected stripped key HOST")
	}
	if _, ok := r.Matched["PORT"]; !ok {
		t.Error("expected stripped key PORT")
	}
}

func TestScope_MultiplePrefix(t *testing.T) {
	r := scope.Scope(baseEnv(), scope.Options{Prefixes: []string{"APP_", "DB_"}})
	if len(r.Matched) != 4 {
		t.Errorf("expected 4 matched, got %d", len(r.Matched))
	}
	if len(r.Unmatched) != 2 {
		t.Errorf("expected 2 unmatched, got %d", len(r.Unmatched))
	}
}

func TestSummary(t *testing.T) {
	r := scope.Scope(baseEnv(), scope.Options{Prefixes: []string{"LOG_"}})
	s := r.Summary()
	if s != "2 matched, 4 unmatched" {
		t.Errorf("unexpected summary: %s", s)
	}
}
