package filter

import (
	"testing"
)

var baseEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"APP_PORT":    "8080",
	"APP_DEBUG":   "true",
	"LOG_LEVEL":   "info",
}

func TestFilter_NoPatterns_AllKept(t *testing.T) {
	res, err := Filter(baseEnv, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != len(baseEnv) {
		t.Errorf("expected %d matched, got %d", len(baseEnv), res.Matched)
	}
	if res.Dropped != 0 {
		t.Errorf("expected 0 dropped, got %d", res.Dropped)
	}
}

func TestFilter_PrefixPattern(t *testing.T) {
	res, err := Filter(baseEnv, Options{Patterns: []string{"DB_*"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != 2 {
		t.Errorf("expected 2 matched, got %d", res.Matched)
	}
	if _, ok := res.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestFilter_SuffixPattern(t *testing.T) {
	res, err := Filter(baseEnv, Options{Patterns: []string{"*_LEVEL"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != 1 {
		t.Errorf("expected 1 matched, got %d", res.Matched)
	}
	if _, ok := res.Env["LOG_LEVEL"]; !ok {
		t.Error("expected LOG_LEVEL in result")
	}
}

func TestFilter_ExactPattern(t *testing.T) {
	res, err := Filter(baseEnv, Options{Patterns: []string{"APP_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != 1 {
		t.Errorf("expected 1 matched, got %d", res.Matched)
	}
}

func TestFilter_Regex(t *testing.T) {
	res, err := Filter(baseEnv, Options{Regex: "^(DB|LOG)_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != 3 {
		t.Errorf("expected 3 matched, got %d", res.Matched)
	}
}

func TestFilter_InvalidRegex(t *testing.T) {
	_, err := Filter(baseEnv, Options{Regex: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestFilter_Invert(t *testing.T) {
	res, err := Filter(baseEnv, Options{Patterns: []string{"DB_*"}, Invert: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != 3 {
		t.Errorf("expected 3 matched, got %d", res.Matched)
	}
	if _, ok := res.Env["DB_HOST"]; ok {
		t.Error("DB_HOST should have been dropped")
	}
}

func TestFilter_DroppedCount(t *testing.T) {
	res, err := Filter(baseEnv, Options{Patterns: []string{"APP_*"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Dropped != 3 {
		t.Errorf("expected 3 dropped, got %d", res.Dropped)
	}
}
