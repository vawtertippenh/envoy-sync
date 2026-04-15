package dedupe_test

import (
	"testing"

	"github.com/yourorg/envoy-sync/internal/dedupe"
)

func lines(ss ...string) []string { return ss }

func TestDedupe_NoDuplicates(t *testing.T) {
	input := lines("FOO=bar", "BAZ=qux", "PORT=8080")
	res := dedupe.Dedupe(input, false)
	if len(res.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %v", res.Duplicates)
	}
	if res.Env["FOO"] != "bar" || res.Env["BAZ"] != "qux" || res.Env["PORT"] != "8080" {
		t.Fatalf("unexpected env map: %v", res.Env)
	}
}

func TestDedupe_KeepFirst(t *testing.T) {
	input := lines("FOO=first", "BAR=hello", "FOO=second")
	res := dedupe.Dedupe(input, false)
	if len(res.Duplicates) != 1 || res.Duplicates[0] != "FOO" {
		t.Fatalf("expected duplicate FOO, got %v", res.Duplicates)
	}
	if res.Env["FOO"] != "first" {
		t.Fatalf("expected first value to be kept, got %q", res.Env["FOO"])
	}
}

func TestDedupe_KeepLast(t *testing.T) {
	input := lines("FOO=first", "BAR=hello", "FOO=second")
	res := dedupe.Dedupe(input, true)
	if len(res.Duplicates) != 1 || res.Duplicates[0] != "FOO" {
		t.Fatalf("expected duplicate FOO, got %v", res.Duplicates)
	}
	if res.Env["FOO"] != "second" {
		t.Fatalf("expected last value to be kept, got %q", res.Env["FOO"])
	}
}

func TestDedupe_SkipsComments(t *testing.T) {
	input := lines("# this is a comment", "FOO=bar", "", "  # another", "BAZ=1")
	res := dedupe.Dedupe(input, false)
	if len(res.Duplicates) != 0 {
		t.Fatalf("unexpected duplicates: %v", res.Duplicates)
	}
	if len(res.Env) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(res.Env))
	}
}

func TestDedupe_MultipleDuplicates(t *testing.T) {
	input := lines("A=1", "B=2", "A=3", "B=4", "C=5", "A=6")
	res := dedupe.Dedupe(input, false)
	if len(res.Duplicates) != 2 {
		t.Fatalf("expected 2 duplicate keys, got %v", res.Duplicates)
	}
	if res.Env["A"] != "1" || res.Env["B"] != "2" || res.Env["C"] != "5" {
		t.Fatalf("unexpected env values: %v", res.Env)
	}
}

func TestSummary_NoDuplicates(t *testing.T) {
	res := dedupe.Result{Env: map[string]string{"X": "1"}, Duplicates: nil}
	got := res.Summary()
	if got != "No duplicate keys found." {
		t.Fatalf("unexpected summary: %q", got)
	}
}

func TestSummary_WithDuplicates(t *testing.T) {
	res := dedupe.Result{
		Env:        map[string]string{"A": "1"},
		Duplicates: []string{"A", "B"},
	}
	got := res.Summary()
	if got == "" {
		t.Fatal("expected non-empty summary")
	}
	for _, k := range []string{"A", "B"} {
		if !contains(got, k) {
			t.Errorf("summary missing key %q: %s", k, got)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
