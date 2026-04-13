package validate

import (
	"strings"
	"testing"
)

func TestValidate_AllValid(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	report := Validate("test.env", env, []string{"DATABASE_URL", "PORT"})
	if report.HasErrors() {
		t.Errorf("expected no errors, got: %s", report.Summary())
	}
	if len(report.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(report.Results))
	}
}

func TestValidate_LowercaseKeyWarn(t *testing.T) {
	env := map[string]string{
		"my_var": "value",
	}
	report := Validate("test.env", env, nil)
	if report.HasErrors() {
		t.Error("expected no errors for lowercase key warning")
	}
	if len(report.Results) == 0 {
		t.Fatal("expected at least one warning")
	}
	if report.Results[0].Level != "warn" {
		t.Errorf("expected warn level, got %s", report.Results[0].Level)
	}
}

func TestValidate_EmptyValueWarn(t *testing.T) {
	env := map[string]string{
		"API_KEY": "",
	}
	report := Validate("test.env", env, nil)
	found := false
	for _, r := range report.Results {
		if r.Key == "API_KEY" && r.Level == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected warn for empty API_KEY value")
	}
}

func TestValidate_MissingRequiredError(t *testing.T) {
	env := map[string]string{
		"PORT": "3000",
	}
	report := Validate("test.env", env, []string{"PORT", "DATABASE_URL"})
	if !report.HasErrors() {
		t.Error("expected error for missing required key")
	}
	found := false
	for _, r := range report.Results {
		if r.Key == "DATABASE_URL" && r.Level == "error" {
			found = true
		}
	}
	if !found {
		t.Error("expected error result for DATABASE_URL")
	}
}

func TestReport_Summary_NoIssues(t *testing.T) {
	report := &Report{File: "prod.env", Results: []Result{}}
	if !strings.Contains(report.Summary(), "all entries valid") {
		t.Errorf("unexpected summary: %s", report.Summary())
	}
}

func TestReport_Summary_WithIssues(t *testing.T) {
	report := &Report{
		File: "prod.env",
		Results: []Result{
			{Key: "SECRET", Message: "required key is missing", Level: "error"},
		},
	}
	summary := report.Summary()
	if !strings.Contains(summary, "[ERROR]") {
		t.Errorf("expected [ERROR] in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "SECRET") {
		t.Errorf("expected key name in summary, got: %s", summary)
	}
}
