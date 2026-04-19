package envdiff

import (
	"strings"
	"testing"
)

func TestRender_Added(t *testing.T) {
	r := Summarize(map[string]string{}, map[string]string{"FOO": "bar"})
	lines := Render(r, RenderOptions{})
	if len(lines) != 1 || !strings.HasPrefix(lines[0], "+ FOO") {
		t.Errorf("unexpected output: %v", lines)
	}
}

func TestRender_Removed(t *testing.T) {
	r := Summarize(map[string]string{"FOO": "bar"}, map[string]string{})
	lines := Render(r, RenderOptions{})
	if len(lines) != 1 || !strings.HasPrefix(lines[0], "- FOO") {
		t.Errorf("unexpected output: %v", lines)
	}
}

func TestRender_Modified(t *testing.T) {
	r := Summarize(map[string]string{"FOO": "old"}, map[string]string{"FOO": "new"})
	lines := Render(r, RenderOptions{})
	if len(lines) != 1 || !strings.Contains(lines[0], "old -> new") {
		t.Errorf("unexpected output: %v", lines)
	}
}

func TestRender_MaskValues(t *testing.T) {
	r := Summarize(map[string]string{}, map[string]string{"SECRET": "topsecret"})
	lines := Render(r, RenderOptions{MaskValues: true})
	if !strings.Contains(lines[0], masked) {
		t.Errorf("expected masked value, got: %s", lines[0])
	}
}

func TestRender_ShowUnchanged(t *testing.T) {
	r := Summarize(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	lines := Render(r, RenderOptions{ShowUnchanged: true})
	if len(lines) != 1 || !strings.HasPrefix(lines[0], "  A") {
		t.Errorf("expected unchanged line, got: %v", lines)
	}
}

func TestRender_HideUnchanged(t *testing.T) {
	r := Summarize(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	lines := Render(r, RenderOptions{ShowUnchanged: false})
	if len(lines) != 0 {
		t.Errorf("expected no lines, got: %v", lines)
	}
}

func TestSummary_Format(t *testing.T) {
	r := Summarize(
		map[string]string{"A": "1", "B": "2"},
		map[string]string{"A": "changed", "C": "3"},
	)
	s := Summary(r)
	if s != "+1 -1 ~1" {
		t.Errorf("unexpected summary: %s", s)
	}
}
