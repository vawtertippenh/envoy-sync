package envquote

import (
	"strings"
	"testing"
)

var baseEnv = map[string]string{
	"PLAIN":      "hello",
	"WITH_SPACE": "hello world",
	"WITH_HASH":  "val#comment",
	"EMPTY":      "",
	"ALREADY":    `"quoted"`,
}

func TestQuote_NoForce_OnlyNecessary(t *testing.T) {
	r := Quote(baseEnv, Options{Style: StyleDouble})
	if r.Env["PLAIN"] != "hello" {
		t.Errorf("expected plain value unchanged, got %q", r.Env["PLAIN"])
	}
	if !strings.HasPrefix(r.Env["WITH_SPACE"], "\"") {
		t.Errorf("expected WITH_SPACE to be quoted, got %q", r.Env["WITH_SPACE"])
	}
}

func TestQuote_ForceAll(t *testing.T) {
	env := map[string]string{"A": "simple", "B": "also"}
	r := Quote(env, Options{Style: StyleDouble, ForceAll: true})
	for k, v := range r.Env {
		if !strings.HasPrefix(v, "\"") || !strings.HasSuffix(v, "\"") {
			t.Errorf("key %s: expected double-quoted, got %q", k, v)
		}
	}
	if r.Quoted != 2 {
		t.Errorf("expected 2 quoted, got %d", r.Quoted)
	}
}

func TestQuote_SingleStyle(t *testing.T) {
	env := map[string]string{"MSG": "hello world"}
	r := Quote(env, Options{Style: StyleSingle})
	if r.Env["MSG"] != "'hello world'" {
		t.Errorf("expected single-quoted, got %q", r.Env["MSG"])
	}
}

func TestQuote_AutoStyle_PrefersDouble_WhenSinglePresent(t *testing.T) {
	env := map[string]string{"V": "it's here"}
	r := Quote(env, Options{Style: StyleAuto})
	if !strings.HasPrefix(r.Env["V"], "\"") {
		t.Errorf("expected double-quote for value with single quote, got %q", r.Env["V"])
	}
}

func TestQuote_StripExisting(t *testing.T) {
	env := map[string]string{"K": `"already quoted"`}
	r := Quote(env, Options{Style: StyleDouble, StripExisting: true, ForceAll: true})
	// Should strip then re-quote — no double-double-quoting
	if strings.Contains(r.Env["K"], `""`) {
		t.Errorf("double-quoting detected: %q", r.Env["K"])
	}
}

func TestQuote_CountsAccurate(t *testing.T) {
	env := map[string]string{
		"PLAIN": "nospace",
		"NEEDS": "has space",
	}
	r := Quote(env, Options{Style: StyleDouble})
	if r.Quoted != 1 {
		t.Errorf("expected 1 quoted, got %d", r.Quoted)
	}
	if r.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", r.Skipped)
	}
}

func TestRender_OutputFormat(t *testing.T) {
	r := Result{
		Env: map[string]string{"FOO": `"bar"`, "BAZ": "qux"},
	}
	out := Render(r)
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got: %s", out)
	}
	if !strings.Contains(out, `FOO="bar"`) {
		t.Errorf("expected FOO=\"bar\" in output, got: %s", out)
	}
}

func TestNeedsQuoting_Chars(t *testing.T) {
	cases := map[string]bool{
		"plain":      false,
		"with space": true,
		"hash#val":   true,
		"dollar$val": true,
		"tab\tval":   true,
	}
	for v, want := range cases {
		got := needsQuoting(v)
		if got != want {
			t.Errorf("needsQuoting(%q) = %v, want %v", v, got, want)
		}
	}
}
