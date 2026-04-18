package envtag

import (
	"strings"
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"APP_NAME":    "myapp",
		"APP_PORT":    "8080",
		"LOG_LEVEL":   "info",
	}
}

func TestTag_NoRules_DefaultTag(t *testing.T) {
	results := Tag(baseEnv(), Options{DefaultTag: "misc"})
	for _, r := range results {
		if len(r.Tags) != 1 || r.Tags[0].Key != "misc" {
			t.Errorf("expected default tag for %s", r.EnvKey)
		}
	}
}

func TestTag_PrefixPattern(t *testing.T) {
	opts := Options{
		Tags: map[string][]string{
			"database": {"DB_*"},
			"app":      {"APP_*"},
		},
	}
	results := Tag(baseEnv(), opts)
	tagMap := map[string]string{}
	for _, r := range results {
		if len(r.Tags) > 0 {
			tagMap[r.EnvKey] = r.Tags[0].Key
		}
	}
	if tagMap["DB_HOST"] != "database" {
		t.Errorf("expected DB_HOST -> database, got %s", tagMap["DB_HOST"])
	}
	if tagMap["APP_NAME"] != "app" {
		t.Errorf("expected APP_NAME -> app, got %s", tagMap["APP_NAME"])
	}
}

func TestTag_ExactMatch(t *testing.T) {
	opts := Options{
		Tags: map[string][]string{
			"logging": {"LOG_LEVEL"},
		},
	}
	results := Tag(baseEnv(), opts)
	for _, r := range results {
		if r.EnvKey == "LOG_LEVEL" {
			if len(r.Tags) != 1 || r.Tags[0].Key != "logging" {
				t.Errorf("expected logging tag for LOG_LEVEL")
			}
			return
		}
	}
	t.Error("LOG_LEVEL not found in results")
}

func TestTag_Untagged_NoDefault(t *testing.T) {
	opts := Options{
		Tags: map[string][]string{
			"database": {"DB_*"},
		},
	}
	results := Tag(baseEnv(), opts)
	for _, r := range results {
		if r.EnvKey == "LOG_LEVEL" && len(r.Tags) != 0 {
			t.Errorf("expected LOG_LEVEL untagged")
		}
	}
}

func TestRender_Output(t *testing.T) {
	results := []Result{
		{EnvKey: "DB_HOST", Tags: []Tag{{Key: "database"}}},
		{EnvKey: "LOG_LEVEL", Tags: []Tag{}},
	}
	out := Render(results)
	if !strings.Contains(out, "DB_HOST: [database]") {
		t.Errorf("expected DB_HOST line, got: %s", out)
	}
	if !strings.Contains(out, "LOG_LEVEL: (untagged)") {
		t.Errorf("expected LOG_LEVEL untagged line, got: %s", out)
	}
}
