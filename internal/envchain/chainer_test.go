package envchain

import (
	"strings"
	"testing"
)

func TestChain_LastWins(t *testing.T) {
	links := []Link{
		{Name: "base", Env: map[string]string{"A": "1", "B": "2"}},
		{Name: "override", Env: map[string]string{"B": "99", "C": "3"}},
	}
	r := Chain(links, false)
	if r.Env["A"] != "1" {
		t.Errorf("expected A=1, got %s", r.Env["A"])
	}
	if r.Env["B"] != "99" {
		t.Errorf("expected B=99, got %s", r.Env["B"])
	}
	if r.Origin["B"] != "override" {
		t.Errorf("expected origin override for B, got %s", r.Origin["B"])
	}
	if r.Env["C"] != "3" {
		t.Errorf("expected C=3, got %s", r.Env["C"])
	}
}

func TestChain_StopOnFirst(t *testing.T) {
	links := []Link{
		{Name: "base", Env: map[string]string{"A": "1", "B": "2"}},
		{Name: "override", Env: map[string]string{"B": "99", "C": "3"}},
	}
	r := Chain(links, true)
	if r.Env["B"] != "2" {
		t.Errorf("expected B=2 (first wins), got %s", r.Env["B"])
	}
	if r.Origin["B"] != "base" {
		t.Errorf("expected origin base for B, got %s", r.Origin["B"])
	}
	if r.Env["C"] != "3" {
		t.Errorf("expected C=3, got %s", r.Env["C"])
	}
}

func TestChain_EmptyLinks(t *testing.T) {
	r := Chain([]Link{}, false)
	if len(r.Env) != 0 {
		t.Errorf("expected empty env")
	}
}

func TestSummary_IncludesOrigin(t *testing.T) {
	r := Result{
		Env:    map[string]string{"FOO": "bar"},
		Origin: map[string]string{"FOO": "prod"},
	}
	lines := Summary(r)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "prod") {
		t.Errorf("expected origin in summary, got %s", lines[0])
	}
}

func TestSummary_SortedKeys(t *testing.T) {
	r := Result{
		Env:    map[string]string{"Z": "1", "A": "2"},
		Origin: map[string]string{"Z": "x", "A": "y"},
	}
	lines := Summary(r)
	if !strings.HasPrefix(lines[0], "A") {
		t.Errorf("expected A first, got %s", lines[0])
	}
}
