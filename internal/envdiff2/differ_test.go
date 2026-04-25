package envdiff2

import "testing"

func base() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"DB_HOST":  "localhost",
		"DB_PASS":  "secret",
	}
}

func TestDiff_NoChanges(t *testing.T) {
	a := base()
	b := base()
	r := Diff(a, b, false)
	if r.HasDiff() {
		t.Fatal("expected no diff")
	}
}

func TestDiff_Added(t *testing.T) {
	a := base()
	b := base()
	b["NEW_KEY"] = "value"
	r := Diff(a, b, false)
	if !r.HasDiff() {
		t.Fatal("expected diff")
	}
	if len(r.Changes) != 1 || r.Changes[0].Kind != Added {
		t.Fatalf("expected one Added change, got %+v", r.Changes)
	}
}

func TestDiff_Removed(t *testing.T) {
	a := base()
	b := base()
	delete(b, "APP_NAME")
	r := Diff(a, b, false)
	if len(r.Changes) != 1 || r.Changes[0].Kind != Removed {
		t.Fatalf("expected one Removed change, got %+v", r.Changes)
	}
}

func TestDiff_Modified(t *testing.T) {
	a := base()
	b := base()
	b["DB_HOST"] = "remotehost"
	r := Diff(a, b, false)
	if len(r.Changes) != 1 || r.Changes[0].Kind != Modified {
		t.Fatalf("expected one Modified change, got %+v", r.Changes)
	}
	if r.Changes[0].OldVal != "localhost" || r.Changes[0].NewVal != "remotehost" {
		t.Fatalf("unexpected values: %+v", r.Changes[0])
	}
}

func TestDiff_IncludeUnchanged(t *testing.T) {
	a := base()
	b := base()
	b["NEW_KEY"] = "x"
	r := Diff(a, b, true)
	unchangedCount := 0
	for _, c := range r.Changes {
		if c.Kind == Unchanged {
			unchangedCount++
		}
	}
	if unchangedCount != 3 {
		t.Fatalf("expected 3 unchanged, got %d", unchangedCount)
	}
}

func TestDiff_SortedKeys(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "2"}
	b := map[string]string{"Z": "1", "A": "3"}
	r := Diff(a, b, false)
	if r.Changes[0].Key != "A" {
		t.Fatalf("expected sorted keys, first=%s", r.Changes[0].Key)
	}
}
