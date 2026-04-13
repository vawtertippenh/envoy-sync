package convert

import (
	"strings"
	"testing"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "localhost",
	"SECRET":   "s3cr3t",
}

func TestConvert_Dotenv(t *testing.T) {
	res, err := Convert(sampleEnv, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Format != FormatDotenv {
		t.Errorf("expected format dotenv, got %s", res.Format)
	}
	for k, v := range sampleEnv {
		expected := k + "=" + v
		if !strings.Contains(res.Output, expected) {
			t.Errorf("expected output to contain %q", expected)
		}
	}
}

func TestConvert_JSON(t *testing.T) {
	res, err := Convert(sampleEnv, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(res.Output, "{") {
		t.Errorf("expected JSON output to start with '{', got: %s", res.Output)
	}
	if !strings.Contains(res.Output, "\"APP_ENV\"") {
		t.Error("expected JSON output to contain APP_ENV key")
	}
}

func TestConvert_YAML(t *testing.T) {
	res, err := Convert(sampleEnv, FormatYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "APP_ENV: ") {
		t.Error("expected YAML output to contain 'APP_ENV: '")
	}
}

func TestConvert_Export(t *testing.T) {
	res, err := Convert(sampleEnv, FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "export APP_ENV=") {
		t.Error("expected export output to contain 'export APP_ENV='")
	}
}

func TestConvert_UnknownFormat(t *testing.T) {
	_, err := Convert(sampleEnv, Format("xml"))
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestConvert_EmptyEnv(t *testing.T) {
	res, err := Convert(map[string]string{}, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "" {
		t.Errorf("expected empty output for empty env, got: %q", res.Output)
	}
}

func TestConvert_DotenvOrdering(t *testing.T) {
	env := map[string]string{"Z_KEY": "last", "A_KEY": "first", "M_KEY": "mid"}
	res, err := Convert(env, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(res.Output), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to start with A_KEY, got: %s", lines[0])
	}
}
