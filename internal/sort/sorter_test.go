package sort

import (
	"testing"
)

var baseEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envoy",
	"APP_VERSION": "1.0.0",
	"Z_LAST":      "z",
	"A_FIRST":     "a",
}

func TestSort_Alpha(t *testing.T) {
	_, keys := Sort(baseEnv, Options{Strategy: Alpha})
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Errorf("alpha order violated: %q > %q", keys[i-1], keys[i])
		}
	}
}

func TestSort_AlphaDescending(t *testing.T) {
	_, keys := Sort(baseEnv, Options{Strategy: Alpha, Descending: true})
	for i := 1; i < len(keys); i++ {
		if keys[i-1] < keys[i] {
			t.Errorf("descending order violated: %q < %q", keys[i-1], keys[i])
		}
	}
}

func TestSort_Length(t *testing.T) {
	_, keys := Sort(baseEnv, Options{Strategy: Length})
	for i := 1; i < len(keys); i++ {
		if len(keys[i-1]) > len(keys[i]) {
			t.Errorf("length order violated: len(%q)=%d > len(%q)=%d",
				keys[i-1], len(keys[i-1]), keys[i], len(keys[i]))
		}
	}
}

func TestSort_Prefix(t *testing.T) {
	_, keys := Sort(baseEnv, Options{Strategy: Prefix})
	// All APP_ keys should appear before DB_ keys, and DB_ before Z_
	prefixOrder := map[string]int{"A": 0, "APP": 1, "DB": 2, "Z": 3}
	lastRank := -1
	for _, k := range keys {
		p := extractPrefix(k)
		rank, ok := prefixOrder[p]
		if !ok {
			continue
		}
		if rank < lastRank {
			t.Errorf("prefix group out of order at key %q", k)
		}
		if rank > lastRank {
			lastRank = rank
		}
	}
}

func TestSort_PreservesValues(t *testing.T) {
	out, keys := Sort(baseEnv, Options{Strategy: Alpha})
	if len(out) != len(baseEnv) {
		t.Fatalf("expected %d keys, got %d", len(baseEnv), len(out))
	}
	for _, k := range keys {
		if out[k] != baseEnv[k] {
			t.Errorf("value mismatch for %q: got %q, want %q", k, out[k], baseEnv[k])
		}
	}
}

func TestSort_EmptyEnv(t *testing.T) {
	out, keys := Sort(map[string]string{}, Options{Strategy: Alpha})
	if len(out) != 0 || len(keys) != 0 {
		t.Error("expected empty output for empty input")
	}
}

func TestExtractPrefix(t *testing.T) {
	cases := []struct{ key, want string }{
		{"DB_HOST", "DB"},
		{"APP_NAME", "APP"},
		{"NOPREFIX", "NOPREFIX"},
		{"_LEADING", "_LEADING"},
	}
	for _, c := range cases {
		got := extractPrefix(c.key)
		if got != c.want {
			t.Errorf("extractPrefix(%q) = %q, want %q", c.key, got, c.want)
		}
	}
}
