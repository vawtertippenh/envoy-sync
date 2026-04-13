package diff

import (
	"strings"
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}

	r := Compare(a, b)

	if r.HasDifferences() {
		t.Error("expected no differences")
	}
	if len(r.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(r.Unchanged))
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "val"}
	b := map[string]string{"FOO": "bar"}

	r := Compare(a, b)

	if !r.HasDifferences() {
		t.Error("expected differences")
	}
	if _, ok := r.OnlyInA["ONLY_A"]; !ok {
		t.Error("expected ONLY_A in OnlyInA")
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "val"}

	r := Compare(a, b)

	if _, ok := r.OnlyInB["ONLY_B"]; !ok {
		t.Error("expected ONLY_B in OnlyInB")
	}
}

func TestCompare_Changed(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}

	r := Compare(a, b)

	pair, ok := r.Changed["FOO"]
	if !ok {
		t.Fatal("expected FOO in Changed")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected changed values: %v", pair)
	}
}

func TestSummary_NoDiff(t *testing.T) {
	a := map[string]string{"X": "1"}
	r := Compare(a, a)
	out := r.Summary("dev", "prod", false)
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-differences message, got: %s", out)
	}
}

func TestSummary_MaskSecrets(t *testing.T) {
	a := map[string]string{"SECRET": "mysecretvalue"}
	b := map[string]string{"SECRET": "othersecret"}

	r := Compare(a, b)
	out := r.Summary("dev", "prod", true)

	if strings.Contains(out, "mysecretvalue") || strings.Contains(out, "othersecret") {
		t.Error("secret values should be masked")
	}
	if !strings.Contains(out, "****") {
		t.Error("expected masked value '****' in output")
	}
}

func TestSummary_ShowsLabels(t *testing.T) {
	a := map[string]string{"ALPHA": "1"}
	b := map[string]string{"BETA": "2"}

	r := Compare(a, b)
	out := r.Summary("local", "remote", false)

	if !strings.Contains(out, "local") {
		t.Error("expected label 'local' in output")
	}
	if !strings.Contains(out, "remote") {
		t.Error("expected label 'remote' in output")
	}
}
