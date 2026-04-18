package envgroup

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"AWS_KEY":     "abc",
		"AWS_SECRET":  "xyz",
		"APP_NAME":    "myapp",
		"UNMATCHED":   "value",
	}
}

func TestGroupBy_BasicPrefixes(t *testing.T) {
	groups := GroupBy(baseEnv(), Options{
		Prefixes: map[string]string{"db": "DB_", "aws": "AWS_"},
	})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "aws" || len(groups[0].Keys) != 2 {
		t.Errorf("unexpected aws group: %+v", groups[0])
	}
	if groups[1].Name != "db" || len(groups[1].Keys) != 2 {
		t.Errorf("unexpected db group: %+v", groups[1])
	}
}

func TestGroupBy_Remainder(t *testing.T) {
	groups := GroupBy(baseEnv(), Options{
		Prefixes:  map[string]string{"db": "DB_"},
		Remainder: "other",
	})
	names := map[string]bool{}
	for _, g := range groups {
		names[g.Name] = true
	}
	if !names["other"] {
		t.Error("expected 'other' remainder group")
	}
	for _, g := range groups {
		if g.Name == "other" && len(g.Keys) != 4 {
			t.Errorf("expected 4 remainder keys, got %d", len(g.Keys))
		}
	}
}

func TestGroupBy_NoRemainderDropsUnmatched(t *testing.T) {
	groups := GroupBy(baseEnv(), Options{
		Prefixes: map[string]string{"db": "DB_"},
	})
	if len(groups) != 1 || groups[0].Name != "db" {
		t.Errorf("unexpected groups: %+v", groups)
	}
}

func TestGroupBy_EmptyEnv(t *testing.T) {
	groups := GroupBy(map[string]string{}, Options{
		Prefixes: map[string]string{"db": "DB_"},
	})
	if len(groups) != 0 {
		t.Errorf("expected no groups, got %d", len(groups))
	}
}

func TestGroupBy_KeysAreSorted(t *testing.T) {
	groups := GroupBy(baseEnv(), Options{
		Prefixes: map[string]string{"db": "DB_"},
	})
	if groups[0].Keys[0] != "DB_HOST" || groups[0].Keys[1] != "DB_PORT" {
		t.Errorf("keys not sorted: %v", groups[0].Keys)
	}
}
