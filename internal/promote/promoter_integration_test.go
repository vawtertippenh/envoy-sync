package promote_test

import (
	"os"
	"path/filepath"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/promote"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestPromoteIntegration_ParseAndPromote(t *testing.T) {
	srcFile := writeTempEnv(t, "FOO=from_staging\nBAR=shared\nSECRET=s3cr3t\n")
	dstFile := writeTempEnv(t, "BAR=local\nLOCAL_ONLY=yes\n")

	src, err := envfile.Parse(srcFile)
	if err != nil {
		t.Fatalf("parse src: %v", err)
	}
	dst, err := envfile.Parse(dstFile)
	if err != nil {
		t.Fatalf("parse dst: %v", err)
	}

	out, res, err := promote.Promote(src, dst, promote.Options{
		Keys:      []string{"FOO", "BAR"},
		Overwrite: false,
	})
	if err != nil {
		t.Fatalf("promote: %v", err)
	}

	// FOO should be promoted (new key)
	if out["FOO"] != "from_staging" {
		t.Errorf("expected FOO=from_staging, got %q", out["FOO"])
	}
	// BAR should be skipped (exists in dst, no overwrite)
	if out["BAR"] != "local" {
		t.Errorf("expected BAR=local (unchanged), got %q", out["BAR"])
	}
	// SECRET should not be promoted (not in Keys filter)
	if _, ok := out["SECRET"]; ok {
		t.Error("SECRET should not have been promoted")
	}
	// LOCAL_ONLY should survive
	if out["LOCAL_ONLY"] != "yes" {
		t.Errorf("expected LOCAL_ONLY=yes, got %q", out["LOCAL_ONLY"])
	}

	if len(res.Promoted) != 1 || res.Promoted[0] != "FOO" {
		t.Errorf("expected promoted=[FOO], got %v", res.Promoted)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "BAR" {
		t.Errorf("expected skipped=[BAR], got %v", res.Skipped)
	}
	_ = filepath.Base(srcFile) // suppress unused import warning
}
