package compare

import (
	"strings"
	"testing"
)

func TestAgainst_ExactMatch(t *testing.T) {
	tmpl := map[string]string{"HOST": "localhost", "PORT": "8080"}
	target := map[string]string{"HOST": "prod.example.com", "PORT": "443"}
	r := Against(tmpl, target)
	if len(r.Missing) != 0 || len(r.Extra) != 0 || len(r.Mismatch) != 0 {
		t.Errorf("expected no issues, got %+v", r)
	}
}

func TestAgainst_MissingKeys(t *testing.T) {
	tmpl := map[string]string{"HOST": "localhost", "SECRET": "x"}
	target := map[string]string{"HOST": "prod"}
	r := Against(tmpl, target)
	if len(r.Missing) != 1 || r.Missing[0] != "SECRET" {
		t.Errorf("expected SECRET missing, got %v", r.Missing)
	}
}

func TestAgainst_ExtraKeys(t *testing.T) {
	tmpl := map[string]string{"HOST": "localhost"}
	target := map[string]string{"HOST": "prod", "DEBUG": "true"}
	r := Against(tmpl, target)
	if len(r.Extra) != 1 || r.Extra[0] != "DEBUG" {
		t.Errorf("expected DEBUG extra, got %v", r.Extra)
	}
}

func TestAgainst_TypeMismatch(t *testing.T) {
	tmpl := map[string]string{"ENABLED": "true"}
	target := map[string]string{"ENABLED": "maybe"}
	r := Against(tmpl, target)
	if len(r.Mismatch) != 1 || r.Mismatch[0] != "ENABLED" {
		t.Errorf("expected ENABLED mismatch, got %v", r.Mismatch)
	}
}

func TestAgainst_NoMismatchSameBool(t *testing.T) {
	tmpl := map[string]string{"ENABLED": "true"}
	target := map[string]string{"ENABLED": "false"}
	r := Against(tmpl, target)
	if len(r.Mismatch) != 0 {
		t.Errorf("expected no mismatch for bool/bool, got %v", r.Mismatch)
	}
}

func TestSummary_Clean(t *testing.T) {
	r := MatchResult{}
	if !strings.Contains(r.Summary(), "✓") {
		t.Error("expected clean summary symbol")
	}
}

func TestSummary_WithIssues(t *testing.T) {
	r := MatchResult{
		Missing:  []string{"DB_PASS"},
		Extra:    []string{"OLD_KEY"},
		Mismatch: []string{"FLAG"},
	}
	s := r.Summary()
	if !strings.Contains(s, "missing") || !strings.Contains(s, "extra") || !strings.Contains(s, "mismatch") {
		t.Errorf("summary missing sections: %s", s)
	}
}
