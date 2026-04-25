package envpivot

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":    "localhost",
		"CACHE_HOST": "localhost",
		"APP_ENV":    "production",
		"LOG_LEVEL":  "info",
		"DEBUG":      "false",
	}
}

func TestPivot_GroupsByValue(t *testing.T) {
	result := Pivot(baseEnv())
	if len(result.Groups) == 0 {
		t.Fatal("expected groups, got none")
	}
	// "localhost" should group DB_HOST and CACHE_HOST
	for _, g := range result.Groups {
		if g.Value == "localhost" {
			if len(g.Keys) != 2 {
				t.Errorf("expected 2 keys for 'localhost', got %d", len(g.Keys))
			}
			if g.Keys[0] != "CACHE_HOST" || g.Keys[1] != "DB_HOST" {
				t.Errorf("unexpected key order: %v", g.Keys)
			}
			return
		}
	}
	t.Error("no group found for value 'localhost'")
}

func TestPivot_SingletonCount(t *testing.T) {
	result := Pivot(baseEnv())
	// APP_ENV, LOG_LEVEL, DEBUG are unique values
	if result.Singletons != 3 {
		t.Errorf("expected 3 singletons, got %d", result.Singletons)
	}
}

func TestPivot_SharedCount(t *testing.T) {
	result := Pivot(baseEnv())
	if result.Shared != 1 {
		t.Errorf("expected 1 shared group, got %d", result.Shared)
	}
}

func TestPivot_EmptyEnv(t *testing.T) {
	result := Pivot(map[string]string{})
	if len(result.Groups) != 0 {
		t.Errorf("expected no groups for empty env, got %d", len(result.Groups))
	}
	if result.Singletons != 0 || result.Shared != 0 {
		t.Error("expected zero counts for empty env")
	}
}

func TestPivot_GroupsSortedByValue(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "zebra",
		"A_KEY": "apple",
		"M_KEY": "mango",
	}
	result := Pivot(env)
	if len(result.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(result.Groups))
	}
	if result.Groups[0].Value != "apple" || result.Groups[1].Value != "mango" || result.Groups[2].Value != "zebra" {
		t.Errorf("groups not sorted by value: %v", result.Groups)
	}
}

func TestSummary_NoEntries(t *testing.T) {
	s := Summary(Result{})
	if s != "no entries" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummary_WithShared(t *testing.T) {
	r := Result{Groups: []Group{{}, {}}, Singletons: 2, Shared: 1}
	s := Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
