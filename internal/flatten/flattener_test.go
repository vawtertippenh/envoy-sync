package flatten

import (
	"strings"
	"testing"
)

func TestFlatten_Simple(t *testing.T) {
	input := map[string]interface{}{
		"HOST": "localhost",
		"PORT": float64(5432),
	}
	out, err := Flatten(input, Options{UpperCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", out["PORT"])
	}
}

func TestFlatten_Nested(t *testing.T) {
	input := map[string]interface{}{
		"db": map[string]interface{}{
			"host": "127.0.0.1",
			"port": float64(3306),
		},
	}
	out, err := Flatten(input, Options{Separator: "_", UpperCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "127.0.0.1" {
		t.Errorf("expected DB_HOST=127.0.0.1, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "3306" {
		t.Errorf("expected DB_PORT=3306, got %q", out["DB_PORT"])
	}
}

func TestFlatten_Prefix(t *testing.T) {
	input := map[string]interface{}{"key": "val"}
	out, err := Flatten(input, Options{Prefix: "APP", Separator: "_", UpperCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_KEY"] != "val" {
		t.Errorf("expected APP_KEY=val, got %v", out)
	}
}

func TestFlatten_BoolAndNil(t *testing.T) {
	input := map[string]interface{}{"ENABLED": true, "EMPTY": nil}
	out, err := Flatten(input, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENABLED"] != "true" {
		t.Errorf("expected true, got %q", out["ENABLED"])
	}
	if out["EMPTY"] != "" {
		t.Errorf("expected empty string, got %q", out["EMPTY"])
	}
}

func TestFlattenJSON_Valid(t *testing.T) {
	data := []byte(`{"app":{"name":"envoy","debug":true}}`)
	out, err := FlattenJSON(data, Options{Separator: "_", UpperCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "envoy" {
		t.Errorf("expected APP_NAME=envoy, got %q", out["APP_NAME"])
	}
}

func TestFlattenJSON_Invalid(t *testing.T) {
	_, err := FlattenJSON([]byte(`not json`), Options{})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestRender_SortedAndQuoted(t *testing.T) {
	flat := map[string]string{
		"Z_KEY": "last",
		"A_KEY": "with spaces",
	}
	out := Render(flat)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != `A_KEY="with spaces"` {
		t.Errorf("unexpected first line: %q", lines[0])
	}
	if lines[1] != "Z_KEY=last" {
		t.Errorf("unexpected second line: %q", lines[1])
	}
}
