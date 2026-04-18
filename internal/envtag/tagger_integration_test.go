package envtag_test

import (
	"os"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envtag"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envtag-*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestTagIntegration_ParseAndTag(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PASS=secret\nAPP_PORT=8080\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	opts := envtag.Options{
		Tags: map[string][]string{
			"database": {"DB_*"},
			"app":      {"APP_*"},
		},
		DefaultTag: "other",
	}
	results := envtag.Tag(env, opts)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if len(r.Tags) == 0 {
			t.Errorf("key %s should have at least one tag", r.EnvKey)
		}
	}
}

func TestTagIntegration_RenderNonEmpty(t *testing.T) {
	path := writeTempEnv(t, "LOG_LEVEL=debug\nSECRET_KEY=abc\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	results := envtag.Tag(env, envtag.Options{DefaultTag: "general"})
	out := envtag.Render(results)
	if out == "" {
		t.Error("expected non-empty render output")
	}
}
