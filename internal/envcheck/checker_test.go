package envcheck

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"HOST":     "localhost",
		"PORT":     "8080",
		"PASSWORD": "",
	}
}

func TestCheck_AllPresent(t *testing.T) {
	r := Check(baseEnv(), []string{"HOST", "PORT"})
	if !r.OK() {
		t.Fatalf("expected OK, got: %s", r.Summary())
	}
}

func TestCheck_MissingKey(t *testing.T) {
	r := Check(baseEnv(), []string{"HOST", "DB_URL"})
	if r.OK() {
		t.Fatal("expected not OK")
	}
	if len(r.Missing) != 1 || r.Missing[0] != "DB_URL" {
		t.Fatalf("unexpected missing: %v", r.Missing)
	}
}

func TestCheck_EmptyKey(t *testing.T) {
	r := Check(baseEnv(), []string{"PASSWORD"})
	if r.OK() {
		t.Fatal("expected not OK due to empty value")
	}
	if len(r.Empty) != 1 || r.Empty[0] != "PASSWORD" {
		t.Fatalf("unexpected empty: %v", r.Empty)
	}
}

func TestCheck_NoRequired(t *testing.T) {
	r := Check(baseEnv(), []string{})
	if !r.OK() {
		t.Fatal("expected OK with no required keys")
	}
}

func TestSummary_OK(t *testing.T) {
	r := Check(baseEnv(), []string{"HOST"})
	if r.Summary() != "all required keys present and non-empty" {
		t.Fatalf("unexpected summary: %s", r.Summary())
	}
}

func TestSummary_MissingAndEmpty(t *testing.T) {
	r := Check(baseEnv(), []string{"HOST", "PASSWORD", "SECRET"})
	s := r.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
	// should mention both missing and empty
	if r.Missing[0] != "SECRET" {
		t.Fatalf("expected SECRET missing, got %v", r.Missing)
	}
	if r.Empty[0] != "PASSWORD" {
		t.Fatalf("expected PASSWORD empty, got %v", r.Empty)
	}
}
