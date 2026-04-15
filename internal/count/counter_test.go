package count

import (
	"strings"
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":       "localhost",
		"DB_PASSWORD":   "secret",
		"DB_PORT":       "5432",
		"AWS_ACCESS_KEY": "AKIA123",
		"AWS_SECRET":    "abc",
		"APP_DEBUG":     "",
		"APP_NAME":      "myapp",
	}
}

func TestCount_Total(t *testing.T) {
	r := Count(baseEnv(), Options{})
	if r.Total != 7 {
		t.Errorf("expected Total=7, got %d", r.Total)
	}
}

func TestCount_EmptyVsNonEmpty(t *testing.T) {
	r := Count(baseEnv(), Options{})
	if r.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", r.Empty)
	}
	if r.NonEmpty != 6 {
		t.Errorf("expected NonEmpty=6, got %d", r.NonEmpty)
	}
}

func TestCount_SensitiveKeys(t *testing.T) {
	r := Count(baseEnv(), Options{})
	// DB_PASSWORD, AWS_ACCESS_KEY, AWS_SECRET should match
	if r.Sensitive != 3 {
		t.Errorf("expected Sensitive=3, got %d", r.Sensitive)
	}
}

func TestCount_ExtraSensitivePattern(t *testing.T) {
	env := map[string]string{
		"APP_API_TOKEN": "tok",
		"APP_CERT":      "pem",
	}
	r := Count(env, Options{SensitivePatterns: []string{"CERT"}})
	if r.Sensitive != 2 {
		t.Errorf("expected Sensitive=2, got %d", r.Sensitive)
	}
}

func TestCount_Prefixes(t *testing.T) {
	r := Count(baseEnv(), Options{})
	if r.Prefixes["DB"] != 3 {
		t.Errorf("expected DB prefix count=3, got %d", r.Prefixes["DB"])
	}
	if r.Prefixes["AWS"] != 2 {
		t.Errorf("expected AWS prefix count=2, got %d", r.Prefixes["AWS"])
	}
	if r.Prefixes["APP"] != 2 {
		t.Errorf("expected APP prefix count=2, got %d", r.Prefixes["APP"])
	}
}

func TestCount_EmptyEnv(t *testing.T) {
	r := Count(map[string]string{}, Options{})
	if r.Total != 0 || r.Empty != 0 || r.Sensitive != 0 {
		t.Errorf("expected all zeros for empty env, got %+v", r)
	}
}

func TestSummary_ContainsFields(t *testing.T) {
	r := Count(baseEnv(), Options{})
	s := Summary(r)
	for _, want := range []string{"Total", "Empty", "Sensitive", "Prefixes", "DB:", "AWS:"} {
		if !strings.Contains(s, want) {
			t.Errorf("Summary missing %q\n%s", want, s)
		}
	}
}

func TestSummary_NoPrefix(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": ""}
	r := Count(env, Options{})
	s := Summary(r)
	if strings.Contains(s, "Prefixes:") {
		t.Errorf("Summary should not contain Prefixes section when none exist")
	}
}
