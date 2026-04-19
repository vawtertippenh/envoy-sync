package envfreeze

import (
	"testing"
)

var baseEnv = map[string]string{
	"APP_NAME": "myapp",
	"DB_PASS":  "secret",
	"DEBUG":    "",
	"PORT":     "8080",
}

func TestFreeze_AllKeys(t *testing.T) {
	res, err := Freeze(baseEnv, nil, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Frozen["APP_NAME"] != "myapp" {
		t.Errorf("expected myapp, got %s", res.Frozen["APP_NAME"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "DEBUG" {
		t.Errorf("expected DEBUG skipped, got %v", res.Skipped)
	}
}

func TestFreeze_DenyList(t *testing.T) {
	res, err := Freeze(baseEnv, nil, Options{DenyKeys: []string{"DB_PASS"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.Frozen["DB_PASS"]; ok {
		t.Error("DB_PASS should be excluded")
	}
	if res.Frozen["PORT"] != "8080" {
		t.Error("PORT should be frozen")
	}
}

func TestFreeze_AllowList(t *testing.T) {
	res, err := Freeze(baseEnv, nil, Options{AllowKeys: []string{"PORT"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Frozen) != 1 || res.Frozen["PORT"] != "8080" {
		t.Errorf("expected only PORT, got %v", res.Frozen)
	}
}

func TestFreeze_SkipsExistingWithoutOverwrite(t *testing.T) {
	existing := map[string]string{"PORT": "9999"}
	res, err := Freeze(baseEnv, existing, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Frozen["PORT"] != "9999" {
		t.Errorf("expected preserved 9999, got %s", res.Frozen["PORT"])
	}
}

func TestFreeze_OverwriteExisting(t *testing.T) {
	existing := map[string]string{"PORT": "9999"}
	res, err := Freeze(baseEnv, existing, Options{OverwriteExisting: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Frozen["PORT"] != "8080" {
		t.Errorf("expected overwritten 8080, got %s", res.Frozen["PORT"])
	}
}

func TestFreeze_NilSrcError(t *testing.T) {
	_, err := Freeze(nil, nil, Options{})
	if err == nil {
		t.Error("expected error for nil src")
	}
}
