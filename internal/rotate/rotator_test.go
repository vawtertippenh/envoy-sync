package rotate

import (
	"strings"
	"testing"
)

var baseEnv = map[string]string{
	"APP_NAME":    "myapp",
	"DB_PASSWORD": "old-db-pass",
	"API_SECRET":  "old-api-secret",
	"PORT":        "8080",
}

func TestRotate_SensitiveKeysChanged(t *testing.T) {
	out, res, err := Rotate(baseEnv, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range res.Rotated {
		if out[k] == baseEnv[k] {
			t.Errorf("key %q was not rotated", k)
		}
	}
}

func TestRotate_NonSensitiveKeysUnchanged(t *testing.T) {
	out, _, err := Rotate(baseEnv, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should not be rotated, got %q", out["APP_NAME"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("PORT should not be rotated, got %q", out["PORT"])
	}
}

func TestRotate_OriginalMapUnmutated(t *testing.T) {
	orig := map[string]string{"SECRET_KEY": "original"}
	_, _, err := Rotate(orig, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orig["SECRET_KEY"] != "original" {
		t.Error("original map was mutated")
	}
}

func TestRotate_FilterByKeys(t *testing.T) {
	out, res, err := Rotate(baseEnv, Options{Keys: []string{"DB_PASSWORD"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 1 || res.Rotated[0] != "DB_PASSWORD" {
		t.Errorf("expected only DB_PASSWORD rotated, got %v", res.Rotated)
	}
	if out["API_SECRET"] != baseEnv["API_SECRET"] {
		t.Error("API_SECRET should not have been rotated")
	}
}

func TestRotate_DryRun(t *testing.T) {
	out, res, err := Rotate(baseEnv, Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) == 0 {
		t.Error("expected rotated keys to be reported in dry-run")
	}
	for _, k := range res.Rotated {
		if out[k] != baseEnv[k] {
			t.Errorf("dry-run: key %q value should be unchanged", k)
		}
	}
}

func TestRotate_CustomLength(t *testing.T) {
	out, res, err := Rotate(baseEnv, Options{Length: 8})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range res.Rotated {
		// hex of 8 bytes = 16 chars
		if len(out[k]) != 16 {
			t.Errorf("key %q: expected 16 hex chars, got %d", k, len(out[k]))
		}
	}
}

func TestRotate_GeneratedValuesAreHex(t *testing.T) {
	out, res, err := Rotate(baseEnv, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const hexChars = "0123456789abcdef"
	for _, k := range res.Rotated {
		for _, ch := range out[k] {
			if !strings.ContainsRune(hexChars, ch) {
				t.Errorf("key %q: non-hex character %q in generated value", k, ch)
			}
		}
	}
}
