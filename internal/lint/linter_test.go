package lint

import (
	"strings"
	"testing"
)

func TestLint_ValidEnv(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	result := Lint(env)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(result.Issues), result.Issues)
	}
	if result.HasErrors() {
		t.Error("expected HasErrors to be false")
	}
}

func TestLint_LowercaseKeyError(t *testing.T) {
	env := map[string]string{"app_name": "myapp"}
	result := Lint(env)
	if !result.HasErrors() {
		t.Error("expected error for lowercase key")
	}
	found := false
	for _, i := range result.Issues {
		if i.Key == "app_name" && i.Severity == SeverityError {
			found = true
		}
	}
	if !found {
		t.Error("expected error issue for key 'app_name'")
	}
}

func TestLint_EmptyValueWarning(t *testing.T) {
	env := map[string]string{"MY_VAR": ""}
	result := Lint(env)
	found := false
	for _, i := range result.Issues {
		if i.Key == "MY_VAR" && i.Severity == SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for empty value")
	}
}

func TestLint_UnresolvedPlaceholderWarning(t *testing.T) {
	env := map[string]string{"API_URL": "https://example.com/${{HOST}}"}
	result := Lint(env)
	found := false
	for _, i := range result.Issues {
		if i.Key == "API_URL" && i.Severity == SeverityWarning && strings.Contains(i.Message, "placeholder") {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for unresolved placeholder")
	}
}

func TestSummary_NoIssues(t *testing.T) {
	r := Result{}
	if r.Summary() != "No lint issues found." {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestSummary_WithIssues(t *testing.T) {
	r := Result{
		Issues: []Issue{
			{Key: "bad_key", Message: "key should be UPPER_SNAKE_CASE", Severity: SeverityError},
		},
	}
	summary := r.Summary()
	if !strings.Contains(summary, "[ERROR]") {
		t.Errorf("expected ERROR in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "bad_key") {
		t.Errorf("expected key name in summary, got: %s", summary)
	}
}
