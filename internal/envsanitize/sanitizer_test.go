package envsanitize

import (
	"strings"
	"testing"
)

var baseEnv = map[string]string{
	"APP_NAME":    "  myapp  ",
	"DB_PASSWORD": "secret\x01value",
	"MULTI_LINE":  "line1\nline2",
	"NORMAL":      "hello",
}

func TestSanitize_NoOptions(t *testing.T) {
	res := Sanitize(baseEnv, Options{})
	if len(res.Changed) != 0 {
		t.Errorf("expected no changes, got %v", res.Changed)
	}
	if res.Env["APP_NAME"] != "  myapp  " {
		t.Errorf("unexpected modification")
	}
}

func TestSanitize_TrimWhitespace(t *testing.T) {
	res := Sanitize(baseEnv, Options{TrimWhitespace: true})
	if res.Env["APP_NAME"] != "myapp" {
		t.Errorf("expected trimmed value, got %q", res.Env["APP_NAME"])
	}
	found := false
	for _, k := range res.Changed {
		if k == "APP_NAME" {
			found = true
		}
	}
	if !found {
		t.Error("APP_NAME should appear in Changed")
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	res := Sanitize(baseEnv, Options{StripControlChars: true})
	if strings.ContainsRune(res.Env["DB_PASSWORD"], '\x01') {
		t.Error("control char should have been stripped")
	}
	if res.Env["DB_PASSWORD"] != "secretvalue" {
		t.Errorf("unexpected value: %q", res.Env["DB_PASSWORD"])
	}
}

func TestSanitize_ReplaceNewlines(t *testing.T) {
	res := Sanitize(baseEnv, Options{ReplaceNewlines: true})
	if strings.Contains(res.Env["MULTI_LINE"], "\n") {
		t.Error("newline should have been replaced")
	}
	if res.Env["MULTI_LINE"] != `line1\nline2` {
		t.Errorf("unexpected value: %q", res.Env["MULTI_LINE"])
	}
}

func TestSanitize_MaxLength(t *testing.T) {
	env := map[string]string{"KEY": "abcdefghij"}
	res := Sanitize(env, Options{MaxLength: 5})
	if res.Env["KEY"] != "abcde" {
		t.Errorf("expected truncated value, got %q", res.Env["KEY"])
	}
	if len(res.Changed) != 1 || res.Changed[0] != "KEY" {
		t.Errorf("expected KEY in Changed, got %v", res.Changed)
	}
}

func TestSanitize_MaxLength_Zero_NoTruncation(t *testing.T) {
	env := map[string]string{"KEY": strings.Repeat("x", 500)}
	res := Sanitize(env, Options{MaxLength: 0})
	if len(res.Env["KEY"]) != 500 {
		t.Error("value should not be truncated when MaxLength is 0")
	}
}

func TestSanitize_OriginalUnmutated(t *testing.T) {
	orig := map[string]string{"K": "  v  "}
	Sanitize(orig, Options{TrimWhitespace: true})
	if orig["K"] != "  v  " {
		t.Error("original map should not be mutated")
	}
}
