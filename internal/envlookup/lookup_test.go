package envlookup_test

import (
	"strings"
	"testing"

	"envoy-sync/internal/envlookup"
)

var base = map[string]string{
	"APP_NAME":     "myapp",
	"DB_PASSWORD":  "s3cr3t",
	"API_KEY":      "abc123",
	"DEBUG":        "true",
}

func TestLookup_AllKeys(t *testing.T) {
	results := envlookup.Lookup(base, envlookup.Options{})
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	// Results should be sorted.
	if results[0].Key != "API_KEY" {
		t.Errorf("expected first key API_KEY, got %s", results[0].Key)
	}
}

func TestLookup_SpecificKeys(t *testing.T) {
	results := envlookup.Lookup(base, envlookup.Options{Keys: []string{"APP_NAME", "MISSING"}})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Found || results[0].Value != "myapp" {
		t.Errorf("expected APP_NAME=myapp found=true")
	}
	if results[1].Found {
		t.Errorf("expected MISSING to be not found")
	}
}

func TestLookup_MaskSensitive(t *testing.T) {
	results := envlookup.Lookup(base, envlookup.Options{
		Keys:          []string{"DB_PASSWORD", "API_KEY", "APP_NAME"},
		MaskSensitive: true,
	})
	for _, r := range results {
		switch r.Key {
		case "DB_PASSWORD", "API_KEY":
			if r.Value != "***" || !r.Masked {
				t.Errorf("expected %s to be masked", r.Key)
			}
		case "APP_NAME":
			if r.Masked {
				t.Errorf("APP_NAME should not be masked")
			}
		}
	}
}

func TestLookup_CaseFold(t *testing.T) {
	results := envlookup.Lookup(base, envlookup.Options{
		Keys:     []string{"app_name"},
		CaseFold: true,
	})
	if len(results) != 1 || !results[0].Found {
		t.Fatal("expected case-folded match")
	}
	if results[0].Value != "myapp" {
		t.Errorf("unexpected value: %s", results[0].Value)
	}
}

func TestLookup_CaseFold_Disabled(t *testing.T) {
	results := envlookup.Lookup(base, envlookup.Options{
		Keys:     []string{"app_name"},
		CaseFold: false,
	})
	if results[0].Found {
		t.Error("expected no match without case fold")
	}
}

func TestRender_NotFound(t *testing.T) {
	results := envlookup.Lookup(base, envlookup.Options{Keys: []string{"GHOST"}})
	out := envlookup.Render(results)
	if !strings.Contains(out, "not found") {
		t.Errorf("expected 'not found' in output: %s", out)
	}
}

func TestRender_MaskedLabel(t *testing.T) {
	results := envlookup.Lookup(base, envlookup.Options{
		Keys:          []string{"DB_PASSWORD"},
		MaskSensitive: true,
	})
	out := envlookup.Render(results)
	if !strings.Contains(out, "[masked]") {
		t.Errorf("expected [masked] label in output: %s", out)
	}
}
