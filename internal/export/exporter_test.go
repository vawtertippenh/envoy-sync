package export

import (
	"strings"
	"testing"
)

var sampleEnv = map[string]string{
	"APP_NAME": "envoy",
	"PORT":     "8080",
	"SECRET":   "***",
}

func TestExport_Dotenv(t *testing.T) {
	out, err := Export(sampleEnv, Options{Format: FormatDotenv, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APP_NAME=") {
		t.Errorf("expected first line to start with APP_NAME=, got %s", lines[0])
	}
}

func TestExport_JSON(t *testing.T) {
	out, err := Export(map[string]string{"KEY": "val"}, Options{Format: FormatJSON, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"KEY": "val"`) {
		t.Errorf("expected JSON to contain key/value, got: %s", out)
	}
}

func TestExport_Shell(t *testing.T) {
	out, err := Export(map[string]string{"HOME": "/root"}, Options{Format: FormatShell, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export HOME=") {
		t.Errorf("expected shell export statement, got: %s", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	_, err := Export(sampleEnv, Options{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExport_EmptyEnv(t *testing.T) {
	out, err := Export(map[string]string{}, Options{Format: FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}
