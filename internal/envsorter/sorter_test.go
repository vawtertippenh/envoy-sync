package envsorter_test

import (
	"testing"

	envsorter "envoy-sync/internal/envsorter"
)

func baseEnv() map[string]string {
	return map[string]string{
		"ZEBRA":       "last",
		"ALPHA":       "first",
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"APP_NAME":    "myapp",
		"APP_VERSION": "1.0",
		"TIMEOUT":     "30",
	}
}

func TestSort_Alpha(t *testing.T) {
	res := envsorter.Sort(baseEnv(), envsorter.Options{Strategy: envsorter.StrategyAlpha})
	if len(res.Keys) != 7 {
		t.Fatalf("expected 7 keys, got %d", len(res.Keys))
	}
	if res.Keys[0] != "ALPHA" {
		t.Errorf("expected first key ALPHA, got %s", res.Keys[0])
	}
	if res.Keys[len(res.Keys)-1] != "ZEBRA" {
		t.Errorf("expected last key ZEBRA, got %s", res.Keys[len(res.Keys)-1])
	}
}

func TestSort_AlphaDescending(t *testing.T) {
	res := envsorter.Sort(baseEnv(), envsorter.Options{Strategy: envsorter.StrategyAlpha, Descending: true})
	if res.Keys[0] != "ZEBRA" {
		t.Errorf("expected first key ZEBRA, got %s", res.Keys[0])
	}
}

func TestSort_Length(t *testing.T) {
	res := envsorter.Sort(baseEnv(), envsorter.Options{Strategy: envsorter.StrategyLength})
	for i := 1; i < len(res.Keys); i++ {
		if len(res.Keys[i]) < len(res.Keys[i-1]) {
			t.Errorf("keys not sorted by length at index %d: %s < %s", i, res.Keys[i], res.Keys[i-1])
		}
	}
}

func TestSort_Prefix(t *testing.T) {
	res := envsorter.Sort(baseEnv(), envsorter.Options{Strategy: envsorter.StrategyPrefix, PrefixSep: "_"})
	// APP_ keys should appear before DB_ keys, which appear before standalone keys
	appSeen, dbSeen := false, false
	for _, k := range res.Keys {
		if len(k) >= 3 && k[:3] == "APP" {
			appSeen = true
			if dbSeen {
				t.Errorf("APP key %s appeared after DB key", k)
			}
		}
		if len(k) >= 2 && k[:2] == "DB" {
			dbSeen = true
		}
	}
	if !appSeen || !dbSeen {
		t.Error("expected both APP and DB prefixed keys")
	}
}

func TestSort_Value(t *testing.T) {
	env := map[string]string{"C": "charlie", "A": "alpha", "B": "beta"}
	res := envsorter.Sort(env, envsorter.Options{Strategy: envsorter.StrategyValue})
	if res.Keys[0] != "A" {
		t.Errorf("expected A first (value alpha), got %s", res.Keys[0])
	}
}

func TestSort_PreservesValues(t *testing.T) {
	env := baseEnv()
	res := envsorter.Sort(env, envsorter.Options{})
	for k, v := range env {
		if res.Env[k] != v {
			t.Errorf("value mismatch for %s: want %s got %s", k, v, res.Env[k])
		}
	}
}

func TestSort_EmptyEnv(t *testing.T) {
	res := envsorter.Sort(map[string]string{}, envsorter.Options{})
	if len(res.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(res.Keys))
	}
}
