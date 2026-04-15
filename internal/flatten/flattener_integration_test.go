package flatten_test

import (
	"os"
	"strings"
	"testing"

	"github.com/yourorg/envoy-sync/internal/flatten"
)

func writeTempJSON(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "flatten-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestFlattenIntegration_JSONFileRoundtrip(t *testing.T) {
	jsonContent := `{
  "database": {
    "host": "db.prod.internal",
    "port": 5432,
    "ssl": true
  },
  "app": {
    "name": "envoy-sync",
    "debug": false
  }
}`
	path := writeTempJSON(t, jsonContent)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	out, err := flatten.FlattenJSON(data, flatten.Options{Separator: "_", UpperCase: true})
	if err != nil {
		t.Fatalf("FlattenJSON: %v", err)
	}
	expected := map[string]string{
		"DATABASE_HOST": "db.prod.internal",
		"DATABASE_PORT": "5432",
		"DATABASE_SSL":  "true",
		"APP_NAME":      "envoy-sync",
		"APP_DEBUG":     "false",
	}
	for k, want := range expected {
		if got := out[k]; got != want {
			t.Errorf("key %s: want %q, got %q", k, want, got)
		}
	}
}

func TestFlattenIntegration_RenderOutput(t *testing.T) {
	input := map[string]interface{}{
		"service": map[string]interface{}{
			"url":     "https://api.example.com",
			"timeout": float64(30),
		},
	}
	flat, err := flatten.Flatten(input, flatten.Options{Separator: "_", UpperCase: true})
	if err != nil {
		t.Fatalf("Flatten: %v", err)
	}
	rendered := flatten.Render(flat)
	if !strings.Contains(rendered, "SERVICE_URL=https://api.example.com") {
		t.Errorf("missing SERVICE_URL in output:\n%s", rendered)
	}
	if !strings.Contains(rendered, "SERVICE_TIMEOUT=30") {
		t.Errorf("missing SERVICE_TIMEOUT in output:\n%s", rendered)
	}
}
