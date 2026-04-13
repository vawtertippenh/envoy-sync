package redact_test

import (
	"strings"
	"testing"

	"github.com/envoy-sync/internal/redact"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":     "myapp",
		"DB_PASSWORD":  "secret123",
		"API_KEY":      "key-abc",
		"PORT":         "8080",
		"DEBUG":        "true",
	}
}

func TestRedact_SensitiveKeysDefaultReplacement(t *testing.T) {
	result := redact.Redact(baseEnv(), nil, nil)
	if result.Redacted["DB_PASSWORD"] != redact.DefaultReplacement {
		t.Errorf("expected DB_PASSWORD to be redacted, got %s", result.Redacted["DB_PASSWORD"])
	}
	if result.Redacted["API_KEY"] != redact.DefaultReplacement {
		t.Errorf("expected API_KEY to be redacted, got %s", result.Redacted["API_KEY"])
	}
}

func TestRedact_NonSensitiveKeysUnchanged(t *testing.T) {
	result := redact.Redact(baseEnv(), nil, nil)
	if result.Redacted["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %s", result.Redacted["APP_NAME"])
	}
	if result.Redacted["PORT"] != "8080" {
		t.Errorf("expected PORT unchanged, got %s", result.Redacted["PORT"])
	}
}

func TestRedact_CustomRuleOverridesDefault(t *testing.T) {
	rules := []redact.Rule{{Key: "APP_NAME", Replacement: "***"}}
	result := redact.Redact(baseEnv(), rules, nil)
	if result.Redacted["APP_NAME"] != "***" {
		t.Errorf("expected APP_NAME=***, got %s", result.Redacted["APP_NAME"])
	}
}

func TestRedact_CustomRuleEmptyReplacementUsesDefault(t *testing.T) {
	rules := []redact.Rule{{Key: "PORT", Replacement: ""}}
	result := redact.Redact(baseEnv(), rules, nil)
	if result.Redacted["PORT"] != redact.DefaultReplacement {
		t.Errorf("expected PORT to use default replacement, got %s", result.Redacted["PORT"])
	}
}

func TestRedact_AffectedListPopulated(t *testing.T) {
	result := redact.Redact(baseEnv(), nil, nil)
	if len(result.Affected) == 0 {
		t.Error("expected affected list to be non-empty")
	}
}

func TestRedact_Summary_NoAffected(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	result := redact.Redact(env, nil, nil)
	if result.Summary() != "No keys redacted." {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestRedact_Summary_WithAffected(t *testing.T) {
	result := redact.Redact(baseEnv(), nil, nil)
	s := result.Summary()
	if !strings.Contains(s, "redacted") {
		t.Errorf("expected summary to contain 'redacted', got: %s", s)
	}
}

func TestRedact_ExtraPatterns(t *testing.T) {
	env := map[string]string{"MY_CUSTOM_CERT": "certvalue", "PORT": "8080"}
	result := redact.Redact(env, nil, []string{"cert"})
	if result.Redacted["MY_CUSTOM_CERT"] != redact.DefaultReplacement {
		t.Errorf("expected MY_CUSTOM_CERT to be redacted via extra pattern")
	}
}
