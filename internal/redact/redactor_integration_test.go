package redact_test

import (
	"os"
	"testing"

	"github.com/envoy-sync/internal/envfile"
	"github.com/envoy-sync/internal/redact"
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
	return f.Name()
}

func TestRedactIntegration_ParseAndRedact(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nDB_PASSWORD=secret\nPORT=9000\n")
	defer os.Remove(path)

	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result := redact.Redact(env, nil, nil)

	if result.Redacted["DB_PASSWORD"] != redact.DefaultReplacement {
		t.Errorf("expected DB_PASSWORD redacted, got %s", result.Redacted["DB_PASSWORD"])
	}
	if result.Redacted["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %s", result.Redacted["APP_NAME"])
	}
	if result.Redacted["PORT"] != "9000" {
		t.Errorf("expected PORT=9000, got %s", result.Redacted["PORT"])
	}
}

func TestRedactIntegration_CustomRuleFromFile(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nSECRET_TOKEN=tok123\nHOST=localhost\n")
	defer os.Remove(path)

	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	rules := []redact.Rule{{Key: "HOST", Replacement: "<hidden>"}}
	result := redact.Redact(env, rules, nil)

	if result.Redacted["HOST"] != "<hidden>" {
		t.Errorf("expected HOST=<hidden>, got %s", result.Redacted["HOST"])
	}
	if result.Redacted["SECRET_TOKEN"] != redact.DefaultReplacement {
		t.Errorf("expected SECRET_TOKEN redacted")
	}
}
