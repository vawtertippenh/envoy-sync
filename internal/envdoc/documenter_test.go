package envdoc

import (
	"strings"
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "secret",
		"PORT":        "8080",
	}
}

func TestDocument_AllKeys(t *testing.T) {
	r := Document(baseEnv(), Options{})
	if len(r.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(r.Entries))
	}
}

func TestDocument_SensitiveMasked(t *testing.T) {
	r := Document(baseEnv(), Options{SensitiveKeys: []string{"DB_PASSWORD"}})
	for _, e := range r.Entries {
		if e.Key == "DB_PASSWORD" {
			if e.Value != "***" {
				t.Errorf("expected masked value, got %q", e.Value)
			}
			if !e.Sensitive {
				t.Error("expected Sensitive=true")
			}
		}
	}
}

func TestDocument_RequiredFlag(t *testing.T) {
	r := Document(baseEnv(), Options{RequiredKeys: []string{"PORT"}})
	for _, e := range r.Entries {
		if e.Key == "PORT" && !e.Required {
			t.Error("expected Required=true for PORT")
		}
	}
}

func TestDocument_Descriptions(t *testing.T) {
	opts := Options{Descriptions: map[string]string{"APP_NAME": "The application name"}}
	r := Document(baseEnv(), opts)
	for _, e := range r.Entries {
		if e.Key == "APP_NAME" && e.Description != "The application name" {
			t.Errorf("unexpected description: %q", e.Description)
		}
	}
}

func TestDocument_SortedKeys(t *testing.T) {
	r := Document(baseEnv(), Options{})
	if r.Entries[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME first, got %s", r.Entries[0].Key)
	}
}

func TestRender_ContainsHeaders(t *testing.T) {
	r := Document(baseEnv(), Options{})
	out := Render(r)
	if !strings.Contains(out, "| Key |") {
		t.Error("expected markdown header in output")
	}
}

func TestRender_ContainsKeys(t *testing.T) {
	r := Document(baseEnv(), Options{})
	out := Render(r)
	for _, key := range []string{"APP_NAME", "DB_PASSWORD", "PORT"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in rendered output", key)
		}
	}
}
