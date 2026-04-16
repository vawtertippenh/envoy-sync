package envset_test

import (
	"testing"

	"envoy-sync/internal/envset"
)

var envA = map[string]string{"A": "1", "B": "2", "C": "3"}
var envB = map[string]string{"B": "99", "C": "3", "D": "4"}

func TestUnion_CombinesKeys(t *testing.T) {
	r := envset.Union(envA, envB)
	if len(r.Env) != 4 {
		t.Fatalf("expected 4 keys, got %d", len(r.Env))
	}
	if r.Env["B"] != "99" {
		t.Errorf("expected B=99 (b wins), got %s", r.Env["B"])
	}
	if r.Env["A"] != "1" {
		t.Errorf("expected A=1, got %s", r.Env["A"])
	}
}

func TestUnion_KeysSorted(t *testing.T) {
	r := envset.Union(envA, envB)
	for i := 1; i < len(r.Keys); i++ {
		if r.Keys[i-1] > r.Keys[i] {
			t.Errorf("keys not sorted: %v", r.Keys)
		}
	}
}

func TestIntersect_CommonKeysOnly(t *testing.T) {
	r := envset.Intersect(envA, envB)
	if len(r.Env) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Env))
	}
	if _, ok := r.Env["B"]; !ok {
		t.Error("expected B in intersection")
	}
	if r.Env["B"] != "2" {
		t.Errorf("expected value from a, got %s", r.Env["B"])
	}
	if _, ok := r.Env["A"]; ok {
		t.Error("A should not be in intersection")
	}
}

func TestDifference_OnlyInA(t *testing.T) {
	r := envset.Difference(envA, envB)
	if len(r.Env) != 1 {
		t.Fatalf("expected 1 key, got %d", len(r.Env))
	}
	if _, ok := r.Env["A"]; !ok {
		t.Error("expected A in difference")
	}
}

func TestDifference_EmptyWhenSubset(t *testing.T) {
	sub := map[string]string{"A": "1"}
	r := envset.Difference(sub, envA)
	if len(r.Env) != 0 {
		t.Errorf("expected empty difference, got %v", r.Env)
	}
}

func TestUnion_DoesNotMutateInputs(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2", "Y": "3"}
	envset.Union(a, b)
	if a["X"] != "1" {
		t.Error("Union mutated input a")
	}
}
