package audit

import (
	"testing"
)

func TestRecord_AddsEntry(t *testing.T) {
	l := &Log{}
	l.Record("API_KEY", KindMasked, "value masked")

	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries))
	}
	e := l.Entries[0]
	if e.Key != "API_KEY" {
		t.Errorf("expected key API_KEY, got %s", e.Key)
	}
	if e.Kind != KindMasked {
		t.Errorf("expected kind masked, got %s", e.Kind)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSummary_NoEvents(t *testing.T) {
	l := &Log{}
	got := l.Summary()
	want := "audit log: no events recorded"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestSummary_WithEvents(t *testing.T) {
	l := &Log{}
	l.Record("DB_PASS", KindMasked, "masked")
	l.Record("NEW_VAR", KindAdded, "added")
	l.Record("OLD_VAR", KindRemoved, "removed")
	l.Record("HOST", KindChanged, "changed")

	summary := l.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	for _, sub := range []string{"added=1", "changed=1", "masked=1", "removed=1"} {
		if !containsStr(summary, sub) {
			t.Errorf("expected summary to contain %q, got: %s", sub, summary)
		}
	}
}

func TestFilterByKind(t *testing.T) {
	l := &Log{}
	l.Record("SECRET", KindMasked, "masked")
	l.Record("TOKEN", KindMasked, "masked")
	l.Record("HOST", KindAdded, "added")

	masked := l.FilterByKind(KindMasked)
	if len(masked) != 2 {
		t.Errorf("expected 2 masked entries, got %d", len(masked))
	}

	added := l.FilterByKind(KindAdded)
	if len(added) != 1 {
		t.Errorf("expected 1 added entry, got %d", len(added))
	}

	removed := l.FilterByKind(KindRemoved)
	if len(removed) != 0 {
		t.Errorf("expected 0 removed entries, got %d", len(removed))
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
