package lint_test

import (
	"os"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/lint"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLintIntegration_CleanFile(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nPORT=8080\nDEBUG=false\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	result := lint.Lint(env)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues for clean file, got: %v", result.Issues)
	}
}

func TestLintIntegration_MixedFile(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nbad_key=oops\nEMPTY_VAL=\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	result := lint.Lint(env)

	hasError := false
	hasWarning := false
	for _, i := range result.Issues {
		if i.Severity == lint.SeverityError {
			hasError = true
		}
		if i.Severity == lint.SeverityWarning {
			hasWarning = true
		}
	}
	if !hasError {
		t.Error("expected at least one error issue")
	}
	if !hasWarning {
		t.Error("expected at least one warning issue")
	}
}
