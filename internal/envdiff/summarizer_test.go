package envdiff

import (
	"testing"
)

var base = map[string]string{
	"HOST": "localhost",
	"PORT": "5432",
	"SECRET": "abc",
}

func TestSummarize_NoDrift(t *testing.T) {
	r := Summarize(base, base)
	if r.HasDrift() {
		t.Fatal("expected no drift")
	}
	if len(r.Added()) != 0 || len(r.Removed()) != 0 || len(r.Modified()) != 0 {
		t.Fatal("expected no changes")
	}
}

func TestSummarize_Added(t *testing.T) {
	target := map[string]string{
		"HOST":    "localhost",
		"PORT":    "5432",
		"SECRET":  "abc",
		"NEW_KEY": "new",
	}
	r := Summarize(base, target)
	if len(r.Added()) != 1 {
		t.Fatalf("expected 1 added, got %d", len(r.Added()))
	}
	if r.Added()[0].Key != "NEW_KEY" {
		t.Errorf("unexpected key: %s", r.Added()[0].Key)
	}
}

func TestSummarize_Removed(t *testing.T) {
	target := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	r := Summarize(base, target)
	if len(r.Removed()) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(r.Removed()))
	}
	if r.Removed()[0].Key != "SECRET" {
		t.Errorf("unexpected key: %s", r.Removed()[0].Key)
	}
}

func TestSummarize_Modified(t *testing.T) {
	target := map[string]string{
		"HOST":   "remotehost",
		"PORT":   "5432",
		"SECRET": "abc",
	}
	r := Summarize(base, target)
	if len(r.Modified()) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(r.Modified()))
	}
	m := r.Modified()[0]
	if m.Key != "HOST" || m.OldValue != "localhost" || m.NewValue != "remotehost" {
		t.Errorf("unexpected modification: %+v", m)
	}
}

func TestSummarize_HasDrift(t *testing.T) {
	r := Summarize(base, map[string]string{})
	if !r.HasDrift() {
		t.Fatal("expected drift")
	}
}
