package envindex_test

import (
	"testing"

	"github.com/yourorg/envoy-sync/internal/envindex"
)

var baseEnv = map[string]string{
	"APP_HOST":      "localhost",
	"APP_PORT":      "8080",
	"DB_HOST":       "db.local",
	"DB_PASSWORD":   "secret",
	"FEATURE_FLAG":  "true",
	"LOG_LEVEL":     "info",
	"APP_LOG_LEVEL": "debug",
}

func TestBuild_AllEntries(t *testing.T) {
	idx := envindex.Build(baseEnv)
	if len(idx.All()) != len(baseEnv) {
		t.Fatalf("expected %d entries, got %d", len(baseEnv), len(idx.All()))
	}
}

func TestBuild_SortedOrder(t *testing.T) {
	idx := envindex.Build(baseEnv)
	entries := idx.All()
	for i := 1; i < len(entries); i++ {
		if entries[i].Key < entries[i-1].Key {
			t.Errorf("entries not sorted at index %d: %s < %s", i, entries[i].Key, entries[i-1].Key)
		}
	}
}

func TestByPrefix_MatchesCorrectKeys(t *testing.T) {
	idx := envindex.Build(baseEnv)
	results := idx.ByPrefix("APP_")
	if len(results) != 3 {
		t.Fatalf("expected 3 APP_ keys, got %d", len(results))
	}
	for _, e := range results {
		if len(e.Key) < 4 || e.Key[:4] != "APP_" {
			t.Errorf("unexpected key: %s", e.Key)
		}
	}
}

func TestByPrefix_NoMatch(t *testing.T) {
	idx := envindex.Build(baseEnv)
	results := idx.ByPrefix("NOPE_")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestBySuffix_MatchesCorrectKeys(t *testing.T) {
	idx := envindex.Build(baseEnv)
	results := idx.BySuffix("_HOST")
	if len(results) != 2 {
		t.Fatalf("expected 2 _HOST keys, got %d", len(results))
	}
}

func TestBySubstring_MatchesCorrectKeys(t *testing.T) {
	idx := envindex.Build(baseEnv)
	results := idx.BySubstring("LOG")
	if len(results) != 2 {
		t.Fatalf("expected 2 LOG keys, got %d", len(results))
	}
}

func TestBySubstring_EmptySubstring_ReturnsAll(t *testing.T) {
	idx := envindex.Build(baseEnv)
	results := idx.BySubstring("")
	if len(results) != len(baseEnv) {
		t.Errorf("expected all %d entries, got %d", len(baseEnv), len(results))
	}
}

func TestBuild_EmptyEnv(t *testing.T) {
	idx := envindex.Build(map[string]string{})
	if len(idx.All()) != 0 {
		t.Errorf("expected empty index")
	}
}
