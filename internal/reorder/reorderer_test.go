package reorder

import (
	"testing"
)

var baseEnv = map[string]string{
	"ZEBRA":   "1",
	"APPLE":   "2",
	"MANGO":   "3",
	"BANANA":  "4",
}

func TestReorder_Alpha(t *testing.T) {
	res, err := Reorder(baseEnv, Options{Strategy: StrategyAlpha})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"APPLE", "BANANA", "MANGO", "ZEBRA"}
	for i, k := range res.Ordered {
		if k != expected[i] {
			t.Errorf("pos %d: got %q want %q", i, k, expected[i])
		}
	}
}

func TestReorder_Custom(t *testing.T) {
	order := []string{"MANGO", "APPLE", "ZEBRA", "BANANA"}
	res, err := Reorder(baseEnv, Options{Strategy: StrategyCustom, Order: order, PutUnknownLast: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Ordered) != 4 {
		t.Fatalf("expected 4 keys, got %d", len(res.Ordered))
	}
	if res.Ordered[0] != "MANGO" {
		t.Errorf("expected MANGO first, got %s", res.Ordered[0])
	}
}

func TestReorder_Template_UnknownLast(t *testing.T) {
	order := []string{"APPLE", "BANANA"}
	res, err := Reorder(baseEnv, Options{Strategy: StrategyTemplate, Order: order, PutUnknownLast: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Ordered[0] != "APPLE" || res.Ordered[1] != "BANANA" {
		t.Errorf("expected APPLE, BANANA first two, got %v", res.Ordered[:2])
	}
	if len(res.Unknown) != 2 {
		t.Errorf("expected 2 unknown keys, got %d", len(res.Unknown))
	}
}

func TestReorder_Template_UnknownExcluded(t *testing.T) {
	order := []string{"APPLE"}
	res, _ := Reorder(baseEnv, Options{Strategy: StrategyTemplate, Order: order, PutUnknownLast: false})
	if len(res.Ordered) != 1 {
		t.Errorf("expected 1 ordered key, got %d", len(res.Ordered))
	}
	if len(res.Unknown) != 3 {
		t.Errorf("expected 3 unknown, got %d", len(res.Unknown))
	}
}

func TestReorder_DefaultStrategy(t *testing.T) {
	res, err := Reorder(baseEnv, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Ordered) != 4 {
		t.Errorf("expected 4 keys, got %d", len(res.Ordered))
	}
}

func TestReorder_EnvUnmutated(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	Reorder(env, Options{Strategy: StrategyAlpha})
	if len(env) != 2 {
		t.Error("original env was mutated")
	}
}
