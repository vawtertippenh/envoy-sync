package diff2

import (
	"testing"
)

func base() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"DB_HOST":  "localhost",
		"SECRET":   "abc123",
	}
}

func TestDiff_NoChanges(t *testing.T) {
	a := base()
	r := Diff(a, base())
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
	if len(r.Unchanged()) != 3 {
		t.Fatalf("expected 3 unchanged, got %d", len(r.Unchanged()))
	}
}

func (r Result) Unchanged() []Entry { return r.filter(Unchanged) }

func TestDiff_Added(t *testing.T) {
	a := base()
	b := base()
	b["NEW_KEY"] = "newval"
	r := Diff(a, b)
	if len(r.Added()) != 1 || r.Added()[0].Key != "NEW_KEY" {
		t.Fatal("expected NEW_KEY added")
	}
}

func TestDiff_Removed(t *testing.T) {
	a := base()
	b := base()
	delete(b, "DB_HOST")
	r := Diff(a, b)
	if len(r.Removed()) != 1 || r.Removed()[0].Key != "DB_HOST" {
		t.Fatal("expected DB_HOST removed")
	}
}

func TestDiff_Modified(t *testing.T) {
	a := base()
	b := base()
	b["APP_NAME"] = "otherapp"
	r := Diff(a, b)
	m := r.Modified()
	if len(m) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(m))
	}
	if m[0].OldVal != "myapp" || m[0].NewVal != "otherapp" {
		t.Fatal("unexpected modified values")
	}
}

func TestDiff_HasChanges(t *testing.T) {
	a := base()
	b := base()
	b["EXTRA"] = "x"
	if !Diff(a, b).HasChanges() {
		t.Fatal("expected HasChanges true")
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	r := Diff(map[string]string{}, map[string]string{})
	if r.HasChanges() {
		t.Fatal("expected no changes for empty maps")
	}
}
