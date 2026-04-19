package envmerge

import (
	"testing"
)

func TestMerge_LastWins(t *testing.T) {
	a := map[string]string{"FOO": "a", "BAR": "b"}
	b := map[string]string{"FOO": "z", "BAZ": "c"}
	res, err := Merge([]map[string]string{a, b}, Options{Strategy: StrategyLast})
	if err != nil {
		t.Fatal(err)
	}
	if res.Env["FOO"] != "z" {
		t.Errorf("expected z, got %s", res.Env["FOO"])
	}
	if res.Env["BAR"] != "b" {
		t.Errorf("expected b, got %s", res.Env["BAR"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0].Key != "FOO" {
		t.Errorf("expected one conflict on FOO")
	}
}

func TestMerge_FirstWins(t *testing.T) {
	a := map[string]string{"FOO": "first"}
	b := map[string]string{"FOO": "second"}
	res, err := Merge([]map[string]string{a, b}, Options{Strategy: StrategyFirst})
	if err != nil {
		t.Fatal(err)
	}
	if res.Env["FOO"] != "first" {
		t.Errorf("expected first, got %s", res.Env["FOO"])
	}
}

func TestMerge_ErrorOnConflict(t *testing.T) {
	a := map[string]string{"KEY": "v1"}
	b := map[string]string{"KEY": "v2"}
	_, err := Merge([]map[string]string{a, b}, Options{Strategy: StrategyError})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMerge_NoConflict(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	res, err := Merge([]map[string]string{a, b}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
	if res.Env["A"] != "1" || res.Env["B"] != "2" {
		t.Errorf("unexpected env: %v", res.Env)
	}
}

func TestMerge_WithPrefix(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	res, err := Merge([]map[string]string{a}, Options{Prefix: "APP_"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.Env["APP_FOO"]; !ok {
		t.Errorf("expected APP_FOO key")
	}
}

func TestMerge_SameValueNoConflict(t *testing.T) {
	a := map[string]string{"X": "same"}
	b := map[string]string{"X": "same"}
	res, err := Merge([]map[string]string{a, b}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("same value should not be a conflict")
	}
}
