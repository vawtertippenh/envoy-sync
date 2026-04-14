package generate_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-sync/internal/generate"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "myapp",
		"PORT":        "8080",
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
	}
}

func TestGenerate_DefaultPlaceholders(t *testing.T) {
	result := generate.Generate(baseEnv(), generate.Options{})
	if result["APP_NAME"] != "<YOUR_APP_NAME>" {
		t.Errorf("expected placeholder, got %q", result["APP_NAME"])
	}
	if result["PORT"] != "<YOUR_PORT>" {
		t.Errorf("expected placeholder, got %q", result["PORT"])
	}
}

func TestGenerate_SensitiveKeysMaskedWhenEnabled(t *testing.T) {
	result := generate.Generate(baseEnv(), generate.Options{MaskSensitive: true})
	if result["DB_PASSWORD"] == "supersecret" {
		t.Error("expected DB_PASSWORD to be masked")
	}
	if result["API_KEY"] == "abc123" {
		t.Error("expected API_KEY to be masked")
	}
}

func TestGenerate_SensitiveKeysPlaceholderWhenDisabled(t *testing.T) {
	result := generate.Generate(baseEnv(), generate.Options{MaskSensitive: false})
	if result["DB_PASSWORD"] != "<YOUR_DB_PASSWORD>" {
		t.Errorf("expected placeholder, got %q", result["DB_PASSWORD"])
	}
}

func TestGenerate_CustomPlaceholderFormat(t *testing.T) {
	result := generate.Generate(map[string]string{"HOST": "localhost"}, generate.Options{
		PlaceholderFormat: "CHANGE_ME_%s",
	})
	if result["HOST"] != "CHANGE_ME_HOST" {
		t.Errorf("unexpected placeholder: %q", result["HOST"])
	}
}

func TestGenerate_ExtraPatterns(t *testing.T) {
	result := generate.Generate(
		map[string]string{"MY_CUSTOM_TOKEN": "val"},
		generate.Options{MaskSensitive: true, ExtraPatterns: []string{"TOKEN"}},
	)
	if result["MY_CUSTOM_TOKEN"] == "val" {
		t.Error("expected MY_CUSTOM_TOKEN to be masked via extra pattern")
	}
}

func TestRender_ContainsSections(t *testing.T) {
	tmpl := map[string]string{
		"APP_NAME":    "<YOUR_APP_NAME>",
		"DB_PASSWORD": "****",
	}
	out := generate.Render(tmpl, nil)
	if !strings.Contains(out, "# Application settings") {
		t.Error("expected Application settings section")
	}
	if !strings.Contains(out, "# Sensitive / secret settings") {
		t.Error("expected Sensitive section")
	}
	if !strings.Contains(out, "APP_NAME=") {
		t.Error("expected APP_NAME in output")
	}
	if !strings.Contains(out, "DB_PASSWORD=") {
		t.Error("expected DB_PASSWORD in output")
	}
}

func TestRender_EmptyMap(t *testing.T) {
	out := generate.Render(map[string]string{}, nil)
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}
