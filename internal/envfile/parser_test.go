package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ef.ToMap()
	if m["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", m["APP_ENV"])
	}
	if m["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false, got %q", m["DEBUG"])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/mydb"
SECRET='hunter2'
`)
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ef.ToMap()
	if m["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("unexpected DB_URL: %q", m["DB_URL"])
	}
	if m["SECRET"] != "hunter2" {
		t.Errorf("unexpected SECRET: %q", m["SECRET"])
	}
}

func TestParse_CommentsSkipped(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\nFOO=bar\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ef.ToMap()
	if _, ok := m[""]; ok {
		t.Error("comment line should not produce an empty key entry")
	}
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", m["FOO"])
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := Parse(path)
	if err == nil {
		t.Error("expected error for invalid line, got nil")
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestStripQuotes(t *testing.T) {
	cases := []struct{ in, want string }{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`noQuotes`, "noQuotes"},
		{`"mixed'`, `"mixed'`},
	}
	for _, c := range cases {
		if got := stripQuotes(c.in); got != c.want {
			t.Errorf("stripQuotes(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
