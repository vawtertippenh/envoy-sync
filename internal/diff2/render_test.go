package diff2

import (
	"strings"
	"testing"
)

func TestRender_Added(t *testing.T) {
	r := Diff(map[string]string{}, map[string]string{"KEY": "val"})
	out := Render(r, RenderOptions{})
	if !strings.Contains(out, "+ KEY=val") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestRender_Removed(t *testing.T) {
	r := Diff(map[string]string{"KEY": "val"}, map[string]string{})
	out := Render(r, RenderOptions{})
	if !strings.Contains(out, "- KEY=val") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestRender_Modified(t *testing.T) {
	r := Diff(map[string]string{"KEY": "old"}, map[string]string{"KEY": "new"})
	out := Render(r, RenderOptions{})
	if !strings.Contains(out, "~ KEY: old -> new") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestRender_MaskValues(t *testing.T) {
	r := Diff(map[string]string{"KEY": "old"}, map[string]string{"KEY": "new"})
	out := Render(r, RenderOptions{MaskValues: true})
	if strings.Contains(out, "old") || strings.Contains(out, "new") {
		t.Fatal("expected values to be masked")
	}
	if !strings.Contains(out, "***") {
		t.Fatal("expected mask placeholder")
	}
}

func TestRender_ShowUnchanged(t *testing.T) {
	r := Diff(map[string]string{"KEY": "val"}, map[string]string{"KEY": "val"})
	without := Render(r, RenderOptions{ShowUnchanged: false})
	with := Render(r, RenderOptions{ShowUnchanged: true})
	if strings.Contains(without, "KEY") {
		t.Fatal("should not show unchanged when disabled")
	}
	if !strings.Contains(with, "KEY") {
		t.Fatal("should show unchanged when enabled")
	}
}

func TestSummary(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"A": "1", "C": "3"}
	r := Diff(a, b)
	s := Summary(r)
	if !strings.Contains(s, "added=1") || !strings.Contains(s, "removed=1") {
		t.Fatalf("unexpected summary: %s", s)
	}
}
