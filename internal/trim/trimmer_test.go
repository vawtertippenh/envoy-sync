package trim_test

import (
	"testing"

	"envoy-sync/internal/trim"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "myapp",
		"DEBUG":       "true",
		"EMPTY_KEY":   "",
		"DB_PASSWORD": "secret",
		"PORT":        "8080",
	}
}

func TestTrim_NoOptions(t *testing.T) {
	r := trim.Trim(baseEnv(), trim.Options{})
	if len(r.Kept) != 5 {
		t.Errorf("expected 5 kept, got %d", len(r.Kept))
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(r.Removed))
	}
}

func TestTrim_RemoveEmpty(t *testing.T) {
	r := trim.Trim(baseEnv(), trim.Options{RemoveEmpty: true})
	if _, ok := r.Kept["EMPTY_KEY"]; ok {
		t.Error("expected EMPTY_KEY to be removed")
	}
	if len(r.Removed) != 1 || r.Removed[0] != "EMPTY_KEY" {
		t.Errorf("expected [EMPTY_KEY] in removed, got %v", r.Removed)
	}
}

func TestTrim_DenyList(t *testing.T) {
	opts := trim.Options{DenyList: []string{"DEBUG", "PORT"}}
	r := trim.Trim(baseEnv(), opts)
	if _, ok := r.Kept["DEBUG"]; ok {
		t.Error("DEBUG should have been removed")
	}
	if _, ok := r.Kept["PORT"]; ok {
		t.Error("PORT should have been removed")
	}
	if len(r.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(r.Removed))
	}
}

func TestTrim_AllowList(t *testing.T) {
	opts := trim.Options{AllowList: []string{"APP_NAME", "PORT"}}
	r := trim.Trim(baseEnv(), opts)
	if len(r.Kept) != 2 {
		t.Errorf("expected 2 kept, got %d", len(r.Kept))
	}
	if r.Kept["APP_NAME"] != "myapp" {
		t.Errorf("unexpected value for APP_NAME: %s", r.Kept["APP_NAME"])
	}
}

func TestTrim_AllowListAndRemoveEmpty(t *testing.T) {
	opts := trim.Options{
		AllowList:   []string{"APP_NAME", "EMPTY_KEY", "PORT"},
		RemoveEmpty: true,
	}
	r := trim.Trim(baseEnv(), opts)
	// EMPTY_KEY is in allow list but should still be removed by RemoveEmpty
	if _, ok := r.Kept["EMPTY_KEY"]; ok {
		t.Error("EMPTY_KEY should be removed even when in allow list")
	}
	if len(r.Kept) != 2 {
		t.Errorf("expected 2 kept, got %d", len(r.Kept))
	}
}

func TestTrim_EmptyEnv(t *testing.T) {
	r := trim.Trim(map[string]string{}, trim.Options{RemoveEmpty: true})
	if len(r.Kept) != 0 {
		t.Errorf("expected 0 kept, got %d", len(r.Kept))
	}
}
