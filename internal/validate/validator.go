package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Result holds the outcome of validating a single env entry.
type Result struct {
	Key     string
	Message string
	Level   string // "error" or "warn"
}

// Report aggregates all validation results for a file.
type Report struct {
	File    string
	Results []Result
}

// HasErrors returns true if any result is at error level.
func (r *Report) HasErrors() bool {
	for _, res := range r.Results {
		if res.Level == "error" {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary of the report.
func (r *Report) Summary() string {
	if len(r.Results) == 0 {
		return fmt.Sprintf("%s: all entries valid", r.File)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s: %d issue(s) found\n", r.File, len(r.Results)))
	for _, res := range r.Results {
		sb.WriteString(fmt.Sprintf("  [%s] %s: %s\n", strings.ToUpper(res.Level), res.Key, res.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}

var validKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Validate checks a parsed env map for common issues.
// It warns on lowercase keys and errors on empty values for required keys.
func Validate(file string, env map[string]string, required []string) *Report {
	report := &Report{File: file}

	for key, value := range env {
		if !validKeyPattern.MatchString(key) {
			report.Results = append(report.Results, Result{
				Key:     key,
				Message: "key should be uppercase with underscores only (e.g. MY_VAR)",
				Level:   "warn",
			})
		}
		if strings.TrimSpace(value) == "" {
			report.Results = append(report.Results, Result{
				Key:     key,
				Message: "value is empty",
				Level:   "warn",
			})
		}
	}

	for _, req := range required {
		if _, ok := env[req]; !ok {
			report.Results = append(report.Results, Result{
				Key:     req,
				Message: "required key is missing",
				Level:   "error",
			})
		}
	}

	return report
}
