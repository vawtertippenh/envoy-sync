package mask_test

import (
	"testing"

	"envoy-sync/internal/mask"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	cases := []struct {
		key      string
		expected bool
	}{
		{"DB_PASSWORD", true},
		{"API_SECRET", true},
		{"AUTH_TOKEN", true},
		{"AWS_ACCESS_KEY", true},
		{"DATABASE_URL", false},
		{"APP_NAME", false},
		{"PORT", false},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := mask.IsSensitive(tc.key, nil)
			if got != tc.expected {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.expected)
			}
		})
	}
}

func TestIsSensitive_ExtraPatterns(t *testing.T) {
	got := mask.IsSensitive("MY_CUSTOM_CERT", []string{"CERT"})
	if !got {
		t.Error("expected MY_CUSTOM_CERT to be sensitive with extra pattern CERT")
	}
}

func TestMaskValue_SensitiveKey(t *testing.T) {
	result := mask.MaskValue("DB_PASSWORD", "supersecret", nil)
	if result != mask.MaskedValue {
		t.Errorf("expected masked value, got %q", result)
	}
}

func TestMaskValue_NonSensitiveKey(t *testing.T) {
	result := mask.MaskValue("APP_ENV", "production", nil)
	if result != "production" {
		t.Errorf("expected original value, got %q", result)
	}
}

func TestMaskMap(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "s3cr3t",
		"API_TOKEN":   "tok_abc123",
		"PORT":        "8080",
	}

	masked := mask.MaskMap(env, nil)

	if masked["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should not be masked, got %q", masked["APP_NAME"])
	}
	if masked["PORT"] != "8080" {
		t.Errorf("PORT should not be masked, got %q", masked["PORT"])
	}
	if masked["DB_PASSWORD"] != mask.MaskedValue {
		t.Errorf("DB_PASSWORD should be masked, got %q", masked["DB_PASSWORD"])
	}
	if masked["API_TOKEN"] != mask.MaskedValue {
		t.Errorf("API_TOKEN should be masked, got %q", masked["API_TOKEN"])
	}
}
