package promote

import (
	"testing"
)

func TestPromote_AllKeys(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	dst := map[string]string{}
	out, res, err := Promote(src, dst, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected values in output: %v", out)
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}
	out, res, err := Promote(src, dst, Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if out["FOO"] != "old" {
		t.Errorf("expected old value preserved, got %q", out["FOO"])
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}
	out, res, err := Promote(src, dst, Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwrite) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(res.Overwrite))
	}
	if out["FOO"] != "new" {
		t.Errorf("expected new value, got %q", out["FOO"])
	}
}

func TestPromote_FilterByKeys(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	dst := map[string]string{}
	_, res, err := Promote(src, dst, Options{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d: %v", len(res.Promoted), res.Promoted)
	}
}

func TestPromote_DryRun(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	dst := map[string]string{}
	out, res, err := Promote(src, dst, Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 {
		t.Errorf("expected 1 in promoted list, got %d", len(res.Promoted))
	}
	if _, ok := out["FOO"]; ok {
		t.Error("dry run should not modify destination")
	}
}

func TestPromote_NilSource(t *testing.T) {
	_, _, err := Promote(nil, map[string]string{}, Options{})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{
		Promoted:  []string{"A", "B"},
		Skipped:   []string{"C"},
		Overwrite: []string{},
	}
	s := r.Summary()
	if s != "promoted=2 skipped=1 overwritten=0" {
		t.Errorf("unexpected summary: %q", s)
	}
}
