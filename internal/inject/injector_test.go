package inject

import (
	"strings"
	"testing"
)

var baseEnv = map[string]string{
	"APP_NAME": "myapp",
	"APP_ENV":  "production",
}

func TestInject_NewKeys(t *testing.T) {
	pairs := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, res := Inject(baseEnv, pairs, Options{})

	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if len(res.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(res.Injected))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}
}

func TestInject_SkipsExistingWithoutOverwrite(t *testing.T) {
	pairs := map[string]string{"APP_NAME": "other"}
	out, res := Inject(baseEnv, pairs, Options{Overwrite: false})

	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", out["APP_NAME"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME in skipped, got %v", res.Skipped)
	}
}

func TestInject_OverwriteExisting(t *testing.T) {
	pairs := map[string]string{"APP_NAME": "newapp"}
	out, res := Inject(baseEnv, pairs, Options{Overwrite: true})

	if out["APP_NAME"] != "newapp" {
		t.Errorf("expected APP_NAME=newapp, got %q", out["APP_NAME"])
	}
	if len(res.Overwrite) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(res.Overwrite))
	}
}

func TestInject_WithPrefix(t *testing.T) {
	pairs := map[string]string{"HOST": "db.local"}
	out, res := Inject(baseEnv, pairs, Options{Prefix: "DB_"})

	if out["DB_HOST"] != "db.local" {
		t.Errorf("expected DB_HOST=db.local, got %q", out["DB_HOST"])
	}
	if len(res.Injected) != 1 || res.Injected[0] != "DB_HOST" {
		t.Errorf("unexpected injected list: %v", res.Injected)
	}
}

func TestInject_BaseUnmutated(t *testing.T) {
	original := map[string]string{"KEY": "val"}
	pairs := map[string]string{"NEW": "x"}
	Inject(original, pairs, Options{})
	if _, ok := original["NEW"]; ok {
		t.Error("base map should not be mutated")
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{
		Injected:  []string{"FOO"},
		Overwrite: []string{"BAR"},
		Skipped:   []string{"BAZ"},
	}
	s := r.Summary()
	if !strings.Contains(s, "injected: FOO") {
		t.Errorf("summary missing injected: %s", s)
	}
	if !strings.Contains(s, "overwritten: BAR") {
		t.Errorf("summary missing overwritten: %s", s)
	}
	if !strings.Contains(s, "skipped") {
		t.Errorf("summary missing skipped: %s", s)
	}
}

func TestResult_Summary_Empty(t *testing.T) {
	r := Result{}
	if r.Summary() != "nothing to inject" {
		t.Errorf("expected 'nothing to inject', got %q", r.Summary())
	}
}
