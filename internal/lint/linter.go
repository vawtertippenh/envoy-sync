package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity represents the level of a lint issue.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Issue represents a single lint finding.
type Issue struct {
	Key      string
	Message  string
	Severity Severity
}

// Result holds all lint issues for a file.
type Result struct {
	Issues []Issue
}

// HasErrors returns true if any issue is an error.
func (r *Result) HasErrors() bool {
	for _, i := range r.Issues {
		if i.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary string.
func (r *Result) Summary() string {
	if len(r.Issues) == 0 {
		return "No lint issues found."
	}
	var sb strings.Builder
	for _, i := range r.Issues {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", strings.ToUpper(string(i.Severity)), i.Key, i.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}

var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Lint inspects a map of env key/value pairs and returns a Result.
func Lint(env map[string]string) Result {
	var issues []Issue

	for k, v := range env {
		// Keys must match UPPER_SNAKE_CASE
		if !validKeyRe.MatchString(k) {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "key should be UPPER_SNAKE_CASE (e.g. MY_VAR)",
				Severity: SeverityError,
			})
		}

		// Warn on empty values
		if strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "value is empty",
				Severity: SeverityWarning,
			})
		}

		// Warn on values that look like they contain unresolved references
		if strings.Contains(v, "${{") || strings.Contains(v, "%%") {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "value may contain unresolved template placeholder",
				Severity: SeverityWarning,
			})
		}
	}

	return Result{Issues: issues}
}
