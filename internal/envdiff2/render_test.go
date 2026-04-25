package envdiff2

import (
	"strings"
	"testing"
)

func TestRender_Added(t *testing.T) {
	r := Result{Changes: []Change{{Key: "FOO", Kind: Added, NewVal: "bar"}}}
	out := Render(r, RenderOptions{})
	if !strings.Contains(out, "+ FOO=bar") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestRender_Removed(t *testing.T) {
	r := Result{Changes: []Change{{Key: "FOO", Kind: Removed, OldVal: "bar"}}}
	out := Render(r, RenderOptions{})
	if !strings.Contains(out, "- FOO=bar") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestRender_Modified(t *testing.T) {
	r := Result{Changes: []Change{{Key: "HOST", Kind: Modified, OldVal: "a", NewVal: "b"}}}
	out := Render(r, RenderOptions{})
	if !strings.Contains(out, "~ HOST: a -> b") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestRender_MaskSensitive(t *testing.T) {
	r := Result{Changes: []Change{{Key: "DB_PASSWORD", Kind: Added, NewVal: "s3cr3t"}}}
	out := Render(r, RenderOptions{MaskSensitive: true})
	if strings.Contains(out, "s3cr3t") {
		t.Fatalf("sensitive value should be masked, got: %q", out)
	}
}

func TestRender_ShowUnchanged(t *testing.T) {
	r := Result{Changes: []Change{{Key: "APP", Kind: Unchanged, OldVal: "x", NewVal: "x"}}}
	out := Render(r, RenderOptions{ShowUnchanged: true})
	if !strings.Contains(out, "  APP=x") {
		t.Fatalf("expected unchanged line, got: %q", out)
	}
}

func TestSummary_Counts(t *testing.T) {
	r := Result{Changes: []Change{
		{Kind: Added},
		{Kind: Added},
		{Kind: Removed},
		{Kind: Modified},
	}}
	s := Summary(r)
	if !strings.Contains(s, "added=2") || !strings.Contains(s, "removed=1") || !strings.Contains(s, "modified=1") {
		t.Fatalf("unexpected summary: %s", s)
	}
}

func TestSummary_NoDiff(t *testing.T) {
	r := Result{}
	s := Summary(r)
	if s != "added=0 removed=0 modified=0" {
		t.Fatalf("unexpected summary: %s", s)
	}
}
