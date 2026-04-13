package interpolate

import (
	"strings"
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	res := Interpolate(env)
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", res.Env["HOST"])
	}
	if len(res.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", res.Warnings)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"DB_URL": "postgres://${HOST}:${PORT}/mydb",
	}
	res := Interpolate(env)
	want := "postgres://localhost:5432/mydb"
	if res.Env["DB_URL"] != want {
		t.Errorf("expected %q, got %q", want, res.Env["DB_URL"])
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	env := map[string]string{
		"DOMAIN": "example.com",
		"API_URL": "https://$DOMAIN/api",
	}
	res := Interpolate(env)
	want := "https://example.com/api"
	if res.Env["API_URL"] != want {
		t.Errorf("expected %q, got %q", want, res.Env["API_URL"])
	}
}

func TestInterpolate_UndefinedVariable(t *testing.T) {
	env := map[string]string{
		"URL": "https://${UNDEFINED_HOST}/path",
	}
	res := Interpolate(env)
	if len(res.Warnings) == 0 {
		t.Error("expected a warning for undefined variable")
	}
	if !strings.Contains(res.Warnings[0], "UNDEFINED_HOST") {
		t.Errorf("warning should mention variable name, got: %s", res.Warnings[0])
	}
	// original reference should remain unchanged
	if !strings.Contains(res.Env["URL"], "UNDEFINED_HOST") {
		t.Errorf("unresolved reference should stay in value, got: %s", res.Env["URL"])
	}
}

func TestInterpolate_MultipleReferences(t *testing.T) {
	env := map[string]string{
		"USER":    "admin",
		"PASS":    "secret",
		"DB_DSN":  "${USER}:${PASS}@localhost",
	}
	res := Interpolate(env)
	want := "admin:secret@localhost"
	if res.Env["DB_DSN"] != want {
		t.Errorf("expected %q, got %q", want, res.Env["DB_DSN"])
	}
	if len(res.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", res.Warnings)
	}
}

func TestExtractVarName(t *testing.T) {
	cases := []struct{ input, want string }{
		{"${FOO}", "FOO"},
		{"$BAR", "BAR"},
	}
	for _, c := range cases {
		got := extractVarName(c.input)
		if got != c.want {
			t.Errorf("extractVarName(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
