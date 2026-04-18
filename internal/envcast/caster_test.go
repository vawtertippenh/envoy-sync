package envcast

import (
	"sort"
	"testing"
)

func init() {
	// ensure sort is available in caster.go via this package
	_ = sort.Strings
}

var baseEnv = map[string]string{
	"PORT":    "8080",
	"DEBUG":   "true",
	"RATIO":   "3.14",
	"APP":     "envoy",
	"INVALID": "notanumber",
}

func TestCast_StringDefault(t *testing.T) {
	results := Cast(map[string]string{"NAME": "alice"}, Options{})
	if len(results) != 1 || results[0].Casted != "alice" {
		t.Fatalf("expected string cast, got %+v", results)
	}
}

func TestCast_IntSuccess(t *testing.T) {
	results := Cast(map[string]string{"PORT": "8080"}, Options{
		Rules: map[string]Type{"PORT": TypeInt},
	})
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].Casted.(int) != 8080 {
		t.Fatalf("expected 8080, got %v", results[0].Casted)
	}
}

func TestCast_IntFailure(t *testing.T) {
	results := Cast(map[string]string{"X": "abc"}, Options{
		Rules: map[string]Type{"X": TypeInt},
	})
	if results[0].Err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestCast_BoolSuccess(t *testing.T) {
	results := Cast(map[string]string{"DEBUG": "true"}, Options{
		Rules: map[string]Type{"DEBUG": TypeBool},
	})
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].Casted.(bool) != true {
		t.Fatal("expected true")
	}
}

func TestCast_FloatSuccess(t *testing.T) {
	results := Cast(map[string]string{"RATIO": "3.14"}, Options{
		Rules: map[string]Type{"RATIO": TypeFloat},
	})
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
}

func TestCast_SkipUnknown(t *testing.T) {
	results := Cast(map[string]string{"A": "1", "B": "2"}, Options{
		Rules:       map[string]Type{"A": TypeInt},
		SkipUnknown: true,
	})
	if len(results) != 1 || results[0].Key != "A" {
		t.Fatalf("expected only A, got %+v", results)
	}
}

func TestCast_ResultsAreSorted(t *testing.T) {
	env := map[string]string{"Z": "1", "A": "2", "M": "3"}
	results := Cast(env, Options{})
	if results[0].Key != "A" || results[1].Key != "M" || results[2].Key != "Z" {
		t.Fatalf("expected sorted results, got %v %v %v", results[0].Key, results[1].Key, results[2].Key)
	}
}

func TestCast_EmptyEnv(t *testing.T) {
	results := Cast(map[string]string{}, Options{})
	if len(results) != 0 {
		t.Fatalf("expected empty results for empty env, got %+v", results)
	}
}
