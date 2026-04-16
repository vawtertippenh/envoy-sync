package envdoc_test

import (
	"os"
	"strings"
	"testing"

	"envoy-sync/internal/envdoc"
	"envoy-sync/internal/envfile"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdoc-*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestDocIntegration_ParseAndDocument(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nDB_PASSWORD=secret\nPORT=8080\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	opts := envdoc.Options{
		Descriptions:  map[string]string{"APP_NAME": "App identifier"},
		RequiredKeys:  []string{"PORT"},
		SensitiveKeys: []string{"DB_PASSWORD"},
	}
	r := envdoc.Document(env, opts)
	if len(r.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(r.Entries))
	}
	out := envdoc.Render(r)
	if !strings.Contains(out, "App identifier") {
		t.Error("expected description in output")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected masked password in output")
	}
}

func TestDocIntegration_EmptyEnv(t *testing.T) {
	path := writeTempEnv(t, "")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	r := envdoc.Document(env, envdoc.Options{})
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
	out := envdoc.Render(r)
	if !strings.Contains(out, "| Key |") {
		t.Error("expected header even for empty env")
	}
}
