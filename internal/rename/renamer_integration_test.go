package rename_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/rename"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRenameIntegration_ParseAndRename(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=production\n")

	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	out, r := rename.Rename(env, "DB_HOST", "DATABASE_HOST", rename.Options{})
	if r.Skipped {
		t.Fatalf("unexpected skip: %s", r.Reason)
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", out["DATABASE_HOST"])
	}
}

func TestRenameIntegration_RenameManyRoundtrip(t *testing.T) {
	path := writeTempEnv(t, "OLD_KEY=value1\nANOTHER_OLD=value2\n")

	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	pairs := [][2]string{{"OLD_KEY", "NEW_KEY"}, {"ANOTHER_OLD", "ANOTHER_NEW"}}
	out, results := rename.RenameMany(env, pairs, rename.Options{})

	for _, r := range results {
		if r.Skipped {
			t.Errorf("unexpected skip for %s->%s: %s", r.OldKey, r.NewKey, r.Reason)
		}
	}

	if out["NEW_KEY"] != "value1" || out["ANOTHER_NEW"] != "value2" {
		t.Error("values not correctly transferred after rename")
	}

	var sb strings.Builder
	for k, v := range out {
		sb.WriteString(k + "=" + v + "\n")
	}
	if strings.Contains(sb.String(), "OLD_KEY") || strings.Contains(sb.String(), "ANOTHER_OLD") {
		t.Error("old keys should not appear in output")
	}
}
