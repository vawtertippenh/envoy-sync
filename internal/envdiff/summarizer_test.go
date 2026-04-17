package envdiff

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

func TestSummarize_NoDrift(t *testing.T) {
	s := Summarize(base(), base())
	if s.HasDrift() {
		t.Fatal("expected no drift")
	}
	a, r, m := s.Counts()
	if a != 0 || r != 0 || m != 0 {
		t.Fatalf("expected 0/0/0, got %d/%d/%d", a, r, m)
	}
}

func TestSummarize_Added(t *testing.T) {
	target := base()
	target["NEW_KEY"] = "newval"
	s := Summarize(base(), target)
	a, _, _ := s.Counts()
	if a != 1 {
		t.Fatalf("expected 1 added, got %d", a)
	}
	if s.Changes[0].Kind != Added || s.Changes[0].Key != "NEW_KEY" {
		t.Fatalf("unexpected change: %+v", s.Changes[0])
	}
}

func TestSummarize_Removed(t *testing.T) {
	target := base()
	delete(target, "SECRET")
	s := Summarize(base(), target)
	_, r, _ := s.Counts()
	if r != 1 {
		t.Fatalf("expected 1 removed, got %d", r)
	}
}

func TestSummarize_Modified(t *testing.T) {
	target := base()
	target["DB_HOST"] = "prod-db.example.com"
	s := Summarize(base(), target)
	_, _, m := s.Counts()
	if m != 1 {
		t.Fatalf("expected 1 modified, got %d", m)
	}
	for _, c := range s.Changes {
		if c.Key == "DB_HOST" {
			if c.OldValue != "localhost" || c.NewValue != "prod-db.example.com" {
				t.Fatalf("unexpected values: %+v", c)
			}
			return
		}
	}
	t.Fatal("DB_HOST change not found")
}

func TestSummarize_ChangesAreSorted(t *testing.T) {
	s := Summarize(
		map[string]string{"Z_KEY": "z", "A_KEY": "a"},
		map[string]string{"Z_KEY": "changed"},
	)
	if len(s.Changes) < 2 {
		t.Fatal("expected at least 2 changes")
	}
	if s.Changes[0].Key > s.Changes[1].Key {
		t.Fatal("changes not sorted by key")
	}
}
