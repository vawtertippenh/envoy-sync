package envnorm

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"db_host": "localhost",
		"API_KEY":  "  secret  ",
		"empty":    "",
		"Port":     "8080",
	}
}

func TestNormalize_NoOptions(t *testing.T) {
	env := baseEnv()
	r := Normalize(env, Options{})
	if len(r.Env) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(r.Env))
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(r.Changes))
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	r := Normalize(baseEnv(), Options{UppercaseKeys: true})
	if _, ok := r.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to exist")
	}
	if _, ok := r.Env["db_host"]; ok {
		t.Error("expected db_host to be removed")
	}
	found := false
	for _, c := range r.Changes {
		if c.Action == "uppercase_key" && c.OldKey == "db_host" {
			found = true
		}
	}
	if !found {
		t.Error("expected uppercase_key change for db_host")
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	r := Normalize(baseEnv(), Options{TrimValues: true})
	if r.Env["API_KEY"] != "secret" {
		t.Errorf("expected trimmed value, got %q", r.Env["API_KEY"])
	}
	found := false
	for _, c := range r.Changes {
		if c.Action == "trim_value" && c.Key == "API_KEY" {
			found = true
		}
	}
	if !found {
		t.Error("expected trim_value change for API_KEY")
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	r := Normalize(baseEnv(), Options{RemoveEmpty: true})
	if _, ok := r.Env["empty"]; ok {
		t.Error("expected empty key to be removed")
	}
	found := false
	for _, c := range r.Changes {
		if c.Action == "removed_empty" && c.Key == "empty" {
			found = true
		}
	}
	if !found {
		t.Error("expected removed_empty change for empty key")
	}
}

func TestNormalize_CombinedOptions(t *testing.T) {
	r := Normalize(baseEnv(), Options{
		UppercaseKeys: true,
		TrimValues:    true,
		RemoveEmpty:   true,
	})
	if _, ok := r.Env["EMPTY"]; ok {
		t.Error("EMPTY should have been removed")
	}
	if r.Env["API_KEY"] != "secret" {
		t.Errorf("expected trimmed API_KEY, got %q", r.Env["API_KEY"])
	}
	if r.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", r.Env["DB_HOST"])
	}
}
