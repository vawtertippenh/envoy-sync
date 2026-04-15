package template

import (
	"strings"
	"testing"
)

func TestFill_AllProvided(t *testing.T) {
	tmpl := map[string]string{
		"DB_HOST": "<required>",
		"DB_PORT": "5432",
		"APP_ENV": "<optional>",
	}
	values := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5433",
		"APP_ENV": "production",
	}
	r := Fill(tmpl, values)
	if len(r.Missing) != 0 {
		t.Errorf("expected no missing, got %v", r.Missing)
	}
	if r.Filled["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", r.Filled["DB_HOST"])
	}
	if r.Filled["DB_PORT"] != "5433" {
		t.Errorf("expected DB_PORT=5433, got %s", r.Filled["DB_PORT"])
	}
}

func TestFill_MissingRequired(t *testing.T) {
	tmpl := map[string]string{
		"SECRET_KEY": "<required>",
		"TIMEOUT":    "30",
	}
	values := map[string]string{
		"TIMEOUT": "60",
	}
	r := Fill(tmpl, values)
	if len(r.Missing) != 1 || r.Missing[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY missing, got %v", r.Missing)
	}
}

func TestFill_UnusedValues(t *testing.T) {
	tmpl := map[string]string{
		"APP_NAME": "<required>",
	}
	values := map[string]string{
		"APP_NAME": "envoy",
		"EXTRA_KEY": "ignored",
	}
	r := Fill(tmpl, values)
	if len(r.Unused) != 1 || r.Unused[0] != "EXTRA_KEY" {
		t.Errorf("expected EXTRA_KEY unused, got %v", r.Unused)
	}
}

func TestFill_OptionalKeepsDefault(t *testing.T) {
	tmpl := map[string]string{
		"LOG_LEVEL": "info",
	}
	values := map[string]string{}
	r := Fill(tmpl, values)
	if r.Filled["LOG_LEVEL"] != "info" {
		t.Errorf("expected default info, got %s", r.Filled["LOG_LEVEL"])
	}
	if len(r.Missing) != 0 {
		t.Errorf("unexpected missing: %v", r.Missing)
	}
}

func TestSummary_NoIssues(t *testing.T) {
	r := Result{Filled: map[string]string{"A": "1"}}
	out := Summary(r)
	if !strings.Contains(out, "no issues") {
		t.Errorf("expected no issues message, got: %s", out)
	}
}

func TestSummary_WithIssues(t *testing.T) {
	r := Result{
		Filled:  map[string]string{},
		Missing: []string{"DB_PASS"},
		Unused:  []string{"OLD_KEY"},
	}
	out := Summary(r)
	if !strings.Contains(out, "missing required") {
		t.Errorf("expected missing section, got: %s", out)
	}
	if !strings.Contains(out, "unused") {
		t.Errorf("expected unused section, got: %s", out)
	}
}
