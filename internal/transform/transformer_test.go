package transform

import (
	"testing"
)

var baseEnv = map[string]string{
	"APP_NAME":  "myapp",
	"DB_HOST":   "  localhost  ",
	"API_TOKEN": "abc123",
	"LOG_LEVEL": "debug",
}

func TestTransform_Upper(t *testing.T) {
	rules := []Rule{{Key: "LOG_LEVEL", Op: "upper"}}
	res, err := Transform(baseEnv, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["LOG_LEVEL"] != "DEBUG" {
		t.Errorf("expected DEBUG, got %q", res.Env["LOG_LEVEL"])
	}
	if len(res.Changed) != 1 || res.Changed[0] != "LOG_LEVEL" {
		t.Errorf("expected Changed=[LOG_LEVEL], got %v", res.Changed)
	}
}

func TestTransform_Lower(t *testing.T) {
	env := map[string]string{"REGION": "US-EAST-1"}
	rules := []Rule{{Key: "REGION", Op: "lower"}}
	res, err := Transform(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["REGION"] != "us-east-1" {
		t.Errorf("expected us-east-1, got %q", res.Env["REGION"])
	}
}

func TestTransform_Trim(t *testing.T) {
	rules := []Rule{{Key: "DB_HOST", Op: "trim"}}
	res, err := Transform(baseEnv, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", res.Env["DB_HOST"])
	}
}

func TestTransform_Replace(t *testing.T) {
	env := map[string]string{"DB_URL": "postgres://localhost/mydb"}
	rules := []Rule{{Key: "DB_URL", Op: "replace", From: "localhost", To: "prod-db"}}
	res, err := Transform(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["DB_URL"] != "postgres://prod-db/mydb" {
		t.Errorf("unexpected value: %q", res.Env["DB_URL"])
	}
}

func TestTransform_Wildcard(t *testing.T) {
	env := map[string]string{"FOO": "  a  ", "BAR": "  b  "}
	rules := []Rule{{Key: "*", Op: "trim"}}
	res, err := Transform(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "a" || res.Env["BAR"] != "b" {
		t.Errorf("unexpected values: %v", res.Env)
	}
	if len(res.Changed) != 2 {
		t.Errorf("expected 2 changed keys, got %d", len(res.Changed))
	}
}

func TestTransform_UnknownOp(t *testing.T) {
	rules := []Rule{{Key: "APP_NAME", Op: "base64"}}
	_, err := Transform(baseEnv, rules)
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestTransform_NoChanges(t *testing.T) {
	rules := []Rule{{Key: "LOG_LEVEL", Op: "lower"}}
	res, err := Transform(baseEnv, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changed) != 0 {
		t.Errorf("expected no changes, got %v", res.Changed)
	}
}

func TestTransform_BaseUnmutated(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	rules := []Rule{{Key: "KEY", Op: "upper"}}
	_, _ = Transform(env, rules)
	if env["KEY"] != "value" {
		t.Error("original env was mutated")
	}
}
