package clone_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-sync/internal/clone"
)

func tempDest(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), ".env.clone")
}

func TestClone_WritesAllKeys(t *testing.T) {
	src := map[string]string{"APP_NAME": "envoy", "PORT": "8080"}
	dest := tempDest(t)

	res, err := clone.Clone(src, dest, clone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written != 2 {
		t.Errorf("expected 2 written, got %d", res.Written)
	}
	if res.Masked != 0 {
		t.Errorf("expected 0 masked, got %d", res.Masked)
	}

	data, _ := os.ReadFile(dest)
	if !strings.Contains(string(data), "APP_NAME=envoy") {
		t.Errorf("expected APP_NAME=envoy in output")
	}
}

func TestClone_MasksSensitiveKeys(t *testing.T) {
	src := map[string]string{"DB_PASSWORD": "secret123", "APP_ENV": "production"}
	dest := tempDest(t)

	res, err := clone.Clone(src, dest, clone.Options{MaskSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Masked != 1 {
		t.Errorf("expected 1 masked, got %d", res.Masked)
	}

	data, _ := os.ReadFile(dest)
	if strings.Contains(string(data), "secret123") {
		t.Errorf("sensitive value should be masked")
	}
	if !strings.Contains(string(data), "***") {
		t.Errorf("expected placeholder *** in output")
	}
}

func TestClone_CustomPlaceholder(t *testing.T) {
	src := map[string]string{"API_KEY": "abc", "HOST": "localhost"}
	dest := tempDest(t)

	_, err := clone.Clone(src, dest, clone.Options{MaskSensitive: true, Placeholder: "REDACTED"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(dest)
	if !strings.Contains(string(data), "REDACTED") {
		t.Errorf("expected custom placeholder REDACTED in output")
	}
}

func TestClone_EmptyDestPathError(t *testing.T) {
	_, err := clone.Clone(map[string]string{"K": "V"}, "", clone.Options{})
	if err == nil {
		t.Fatal("expected error for empty dest path")
	}
}

func TestClone_SortedOutput(t *testing.T) {
	src := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	dest := tempDest(t)

	_, err := clone.Clone(src, dest, clone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(dest)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
