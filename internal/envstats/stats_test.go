package envstats_test

import (
	"strings"
	"testing"

	"envoy-sync/internal/envstats"
)

var baseEnv = map[string]string{
	"APP_NAME":     "myapp",
	"APP_VERSION":  "1.0.0",
	"DB_PASSWORD":  "s3cr3t",
	"DB_HOST":      "localhost",
	"API_TOKEN":    "tok123",
	"EMPTY_VALUE":  "",
}

func TestAnalyze_Total(t *testing.T) {
	s := envstats.Analyze(baseEnv, nil)
	if s.Total != 6 {
		t.Errorf("expected Total=6, got %d", s.Total)
	}
}

func TestAnalyze_EmptyAndNonEmpty(t *testing.T) {
	s := envstats.Analyze(baseEnv, nil)
	if s.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", s.Empty)
	}
	if s.NonEmpty != 5 {
		t.Errorf("expected NonEmpty=5, got %d", s.NonEmpty)
	}
}

func TestAnalyze_SensitiveKeys(t *testing.T) {
	s := envstats.Analyze(baseEnv, nil)
	// DB_PASSWORD and API_TOKEN are sensitive by default
	if s.Sensitive != 2 {
		t.Errorf("expected Sensitive=2, got %d", s.Sensitive)
	}
}

func TestAnalyze_ExtraSensitivePattern(t *testing.T) {
	s := envstats.Analyze(baseEnv, []string{"VERSION"})
	// DB_PASSWORD, API_TOKEN, APP_VERSION
	if s.Sensitive != 3 {
		t.Errorf("expected Sensitive=3 with extra pattern, got %d", s.Sensitive)
	}
}

func TestAnalyze_Prefixes(t *testing.T) {
	s := envstats.Analyze(baseEnv, nil)
	if s.Prefixes["APP"] != 2 {
		t.Errorf("expected APP prefix count=2, got %d", s.Prefixes["APP"])
	}
	if s.Prefixes["DB"] != 2 {
		t.Errorf("expected DB prefix count=2, got %d", s.Prefixes["DB"])
	}
}

func TestAnalyze_AvgLength(t *testing.T) {
	env := map[string]string{"A": "ab", "B": "abcd"}
	s := envstats.Analyze(env, nil)
	if s.AvgLength != 3.0 {
		t.Errorf("expected AvgLength=3.0, got %.1f", s.AvgLength)
	}
}

func TestAnalyze_EmptyEnv(t *testing.T) {
	s := envstats.Analyze(map[string]string{}, nil)
	if s.Total != 0 {
		t.Errorf("expected Total=0 for empty env")
	}
	if s.AvgLength != 0 {
		t.Errorf("expected AvgLength=0 for empty env")
	}
}

func TestSummary_ContainsFields(t *testing.T) {
	s := envstats.Analyze(baseEnv, nil)
	out := envstats.Summary(s)
	for _, want := range []string{"Total", "Empty", "Sensitive", "Avg", "Prefixes"} {
		if !strings.Contains(out, want) {
			t.Errorf("Summary missing field %q", want)
		}
	}
}
