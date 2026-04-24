package envsqueeze_test

import (
	"testing"

	"github.com/yourorg/envoy-sync/internal/envsqueeze"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":  "myapp",
		"APP_ALIAS": "myapp", // duplicate value
		"DEBUG":     "",
		"SECRET":    "CHANGE_ME",
		"PORT":      "8080",
	}
}

func TestSqueeze_NoOptions(t *testing.T) {
	res := envsqueeze.Squeeze(baseEnv(), envsqueeze.Options{})
	if len(res.Env) != 5 {
		t.Fatalf("expected 5 keys, got %d", len(res.Env))
	}
	if len(res.Dropped) != 0 {
		t.Fatalf("expected no drops, got %v", res.Dropped)
	}
}

func TestSqueeze_RemoveEmpty(t *testing.T) {
	res := envsqueeze.Squeeze(baseEnv(), envsqueeze.Options{RemoveEmpty: true})
	if _, ok := res.Env["DEBUG"]; ok {
		t.Error("DEBUG should have been removed")
	}
	if len(res.Dropped) != 1 || res.Dropped[0] != "DEBUG" {
		t.Errorf("expected [DEBUG] dropped, got %v", res.Dropped)
	}
}

func TestSqueeze_RemovePlaceholders(t *testing.T) {
	res := envsqueeze.Squeeze(baseEnv(), envsqueeze.Options{RemovePlaceholders: true})
	if _, ok := res.Env["SECRET"]; ok {
		t.Error("SECRET with CHANGE_ME should have been removed")
	}
	if len(res.Dropped) != 1 || res.Dropped[0] != "SECRET" {
		t.Errorf("expected [SECRET] dropped, got %v", res.Dropped)
	}
}

func TestSqueeze_DedupeValues(t *testing.T) {
	res := envsqueeze.Squeeze(baseEnv(), envsqueeze.Options{DedupeValues: true})
	// APP_ALIAS comes after APP_NAME alphabetically, so APP_ALIAS should be dropped
	if _, ok := res.Env["APP_ALIAS"]; ok {
		t.Error("APP_ALIAS should have been removed as duplicate value")
	}
	if _, ok := res.Env["APP_NAME"]; !ok {
		t.Error("APP_NAME should be kept as canonical key")
	}
}

func TestSqueeze_AllOptions(t *testing.T) {
	res := envsqueeze.Squeeze(baseEnv(), envsqueeze.Options{
		RemoveEmpty:        true,
		RemovePlaceholders: true,
		DedupeValues:       true,
	})
	expectedKeys := []string{"APP_NAME", "PORT"}
	if len(res.Env) != len(expectedKeys) {
		t.Fatalf("expected %d keys, got %d: %v", len(expectedKeys), len(res.Env), res.Env)
	}
	for _, k := range expectedKeys {
		if _, ok := res.Env[k]; !ok {
			t.Errorf("expected key %q to be present", k)
		}
	}
}

func TestSqueeze_OriginalUnmutated(t *testing.T) {
	env := baseEnv()
	envsqueeze.Squeeze(env, envsqueeze.Options{RemoveEmpty: true, DedupeValues: true})
	if len(env) != 5 {
		t.Errorf("original map should not be mutated, got %d keys", len(env))
	}
}

func TestSqueeze_DroppedSorted(t *testing.T) {
	res := envsqueeze.Squeeze(baseEnv(), envsqueeze.Options{
		RemoveEmpty:        true,
		RemovePlaceholders: true,
		DedupeValues:       true,
	})
	for i := 1; i < len(res.Dropped); i++ {
		if res.Dropped[i] < res.Dropped[i-1] {
			t.Errorf("dropped list not sorted: %v", res.Dropped)
			break
		}
	}
}
