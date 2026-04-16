package copy

import (
	"testing"
)

var srcEnv = map[string]string{
	"DB_HOST": "localhost",
	"DB_PASS": "secret",
	"APP_PORT": "8080",
}

func TestCopy_AllKeys(t *testing.T) {
	dst := map[string]string{"EXISTING": "yes"}
	out, res := Copy(srcEnv, dst, Options{})
	if len(res.Copied) != 3 {
		t.Fatalf("expected 3 copied, got %d", len(res.Copied))
	}
	if out["EXISTING"] != "yes" {
		t.Error("existing key should be preserved")
	}
}

func TestCopy_FilterByKeys(t *testing.T) {
	out, res := Copy(srcEnv, map[string]string{}, Options{Keys: []string{"DB_HOST"}})
	if len(res.Copied) != 1 || out["DB_HOST"] != "localhost" {
		t.Error("expected only DB_HOST copied")
	}
}

func TestCopy_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := map[string]string{"DB_HOST": "remote"}
	out, res := Copy(srcEnv, dst, Options{})
	if out["DB_HOST"] != "remote" {
		t.Error("should not overwrite without flag")
	}
	if len(res.Skipped) == 0 {
		t.Error("expected skipped entry")
	}
}

func TestCopy_OverwriteExisting(t *testing.T) {
	dst := map[string]string{"DB_HOST": "remote"}
	out, res := Copy(srcEnv, dst, Options{Overwrite: true})
	if out["DB_HOST"] != "localhost" {
		t.Error("expected overwrite")
	}
	if len(res.Copied) == 0 {
		t.Error("expected copied entries")
	}
}

func TestCopy_WithPrefix(t *testing.T) {
	out, res := Copy(srcEnv, map[string]string{}, Options{Keys: []string{"DB_HOST"}, Prefix: "PROD_"})
	if out["PROD_DB_HOST"] != "localhost" {
		t.Errorf("expected prefixed key, got %v", out)
	}
	if len(res.Copied) != 1 {
		t.Error("expected 1 copied")
	}
}

func TestCopy_MissingKeySkipped(t *testing.T) {
	_, res := Copy(srcEnv, map[string]string{}, Options{Keys: []string{"NONEXISTENT"}})
	if len(res.Skipped) != 1 {
		t.Error("expected missing key in skipped")
	}
}

func TestCopy_SrcUnmutated(t *testing.T) {
	origLen := len(srcEnv)
	Copy(srcEnv, map[string]string{}, Options{})
	if len(srcEnv) != origLen {
		t.Error("src should not be mutated")
	}
}
