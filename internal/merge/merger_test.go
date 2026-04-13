package merge

import (
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"APP_NAME": "envoy", "PORT": "8080"}
	b := map[string]string{"DEBUG": "true", "LOG_LEVEL": "info"}

	res, err := Merge([]map[string]string{a, b}, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.Env["APP_NAME"] != "envoy" || res.Env["DEBUG"] != "true" {
		t.Errorf("merged env missing expected keys: %v", res.Env)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := map[string]string{"PORT": "8080"}
	b := map[string]string{"PORT": "9090"}

	res, err := Merge([]map[string]string{a, b}, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080 (first wins), got %s", res.Env["PORT"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := map[string]string{"PORT": "8080"}
	b := map[string]string{"PORT": "9090"}

	res, err := Merge([]map[string]string{a, b}, StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["PORT"] != "9090" {
		t.Errorf("expected PORT=9090 (last wins), got %s", res.Env["PORT"])
	}
}

func TestMerge_StrategyError(t *testing.T) {
	a := map[string]string{"SECRET": "abc"}
	b := map[string]string{"SECRET": "xyz"}

	_, err := Merge([]map[string]string{a, b}, StrategyError)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_SameValueNoConflict(t *testing.T) {
	a := map[string]string{"HOST": "localhost"}
	b := map[string]string{"HOST": "localhost"}

	res, err := Merge([]map[string]string{a, b}, StrategyError)
	if err != nil {
		t.Fatalf("identical values should not conflict: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected 0 conflicts for identical values, got %d", len(res.Conflicts))
	}
}

func TestMerge_UnknownStrategy(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}

	_, err := Merge([]map[string]string{a, b}, Strategy("unknown"))
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res, err := Merge([]map[string]string{}, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
}
