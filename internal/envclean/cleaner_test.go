package envclean

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":  "myapp",
		"DB_HOST":   "  localhost  ",
		"EMPTY_KEY": "",
		"API_KEY":   "secret",
	}
}

func TestClean_NoOptions(t *testing.T) {
	env := baseEnv()
	r := Clean(env, Options{})
	if len(r.Env) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(r.Env))
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected no removed keys, got %v", r.Removed)
	}
}

func TestClean_RemoveEmpty(t *testing.T) {
	r := Clean(baseEnv(), Options{RemoveEmpty: true})
	if _, ok := r.Env["EMPTY_KEY"]; ok {
		t.Error("expected EMPTY_KEY to be removed")
	}
	if len(r.Removed) != 1 || r.Removed[0] != "EMPTY_KEY" {
		t.Errorf("unexpected removed list: %v", r.Removed)
	}
}

func TestClean_TrimWhitespace(t *testing.T) {
	r := Clean(baseEnv(), Options{TrimWhitespace: true})
	if r.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected trimmed value, got %q", r.Env["DB_HOST"])
	}
}

func TestClean_TrimAndRemoveEmpty(t *testing.T) {
	env := map[string]string{
		"KEY_A": "  ",
		"KEY_B": "value",
	}
	r := Clean(env, Options{TrimWhitespace: true, RemoveEmpty: true})
	if _, ok := r.Env["KEY_A"]; ok {
		t.Error("expected KEY_A removed after trim")
	}
	if r.Env["KEY_B"] != "value" {
		t.Error("expected KEY_B preserved")
	}
}

func TestClean_OriginalUnmutated(t *testing.T) {
	env := baseEnv()
	Clean(env, Options{RemoveEmpty: true, TrimWhitespace: true})
	if env["EMPTY_KEY"] != "" {
		t.Error("original map should not be mutated")
	}
	if env["DB_HOST"] != "  localhost  " {
		t.Error("original whitespace should be preserved")
	}
}

func TestClean_EmptyEnv(t *testing.T) {
	r := Clean(map[string]string{}, Options{RemoveEmpty: true, TrimWhitespace: true})
	if len(r.Env) != 0 {
		t.Error("expected empty result")
	}
}
