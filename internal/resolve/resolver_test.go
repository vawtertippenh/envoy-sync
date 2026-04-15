package resolve

import (
	"testing"
)

func TestResolve_BaseOnly(t *testing.T) {
	base := map[string]string{"APP_ENV": "dev", "PORT": "8080"}
	r := Resolve(base, nil)
	if r.Env["APP_ENV"] != "dev" {
		t.Errorf("expected dev, got %s", r.Env["APP_ENV"])
	}
	if len(r.Overrides) != 0 {
		t.Errorf("expected no overrides, got %v", r.Overrides)
	}
}

func TestResolve_SingleSource(t *testing.T) {
	base := map[string]string{"APP_ENV": "dev", "PORT": "8080"}
	sources := []Source{
		{Name: "prod", Values: map[string]string{"APP_ENV": "production"}},
	}
	r := Resolve(base, sources)
	if r.Env["APP_ENV"] != "production" {
		t.Errorf("expected production, got %s", r.Env["APP_ENV"])
	}
	if r.Env["PORT"] != "8080" {
		t.Errorf("expected 8080, got %s", r.Env["PORT"])
	}
	if r.Overrides["APP_ENV"] != "prod" {
		t.Errorf("expected override source prod, got %s", r.Overrides["APP_ENV"])
	}
}

func TestResolve_LastSourceWins(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost"}
	sources := []Source{
		{Name: "staging", Values: map[string]string{"DB_HOST": "staging-db"}},
		{Name: "prod", Values: map[string]string{"DB_HOST": "prod-db"}},
	}
	r := Resolve(base, sources)
	if r.Env["DB_HOST"] != "prod-db" {
		t.Errorf("expected prod-db, got %s", r.Env["DB_HOST"])
	}
	if r.Overrides["DB_HOST"] != "prod" {
		t.Errorf("expected override source prod, got %s", r.Overrides["DB_HOST"])
	}
}

func TestResolve_NewKeyFromSource(t *testing.T) {
	base := map[string]string{"APP_ENV": "dev"}
	sources := []Source{
		{Name: "extra", Values: map[string]string{"NEW_KEY": "new_value"}},
	}
	r := Resolve(base, sources)
	if r.Env["NEW_KEY"] != "new_value" {
		t.Errorf("expected new_value, got %s", r.Env["NEW_KEY"])
	}
	if r.Overrides["NEW_KEY"] != "extra" {
		t.Errorf("expected override source extra, got %s", r.Overrides["NEW_KEY"])
	}
}

func TestResolve_BaseUnmutated(t *testing.T) {
	base := map[string]string{"APP_ENV": "dev"}
	sources := []Source{
		{Name: "prod", Values: map[string]string{"APP_ENV": "production"}},
	}
	Resolve(base, sources)
	if base["APP_ENV"] != "dev" {
		t.Errorf("base map was mutated, expected dev got %s", base["APP_ENV"])
	}
}

func TestSummary_NoOverrides(t *testing.T) {
	r := Result{Env: map[string]string{"A": "1"}, Overrides: map[string]string{}}
	lines := Summary(r)
	if len(lines) != 0 {
		t.Errorf("expected empty summary, got %v", lines)
	}
}

func TestSummary_WithOverrides(t *testing.T) {
	r := Result{
		Env: map[string]string{"A": "1", "B": "2"},
		Overrides: map[string]string{"A": "prod", "B": "staging"},
	}
	lines := Summary(r)
	if len(lines) != 2 {
		t.Errorf("expected 2 summary lines, got %d", len(lines))
	}
	if lines[0] != "A <- prod" {
		t.Errorf("unexpected line: %s", lines[0])
	}
}
