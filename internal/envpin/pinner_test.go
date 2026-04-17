package envpin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func tempPinPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "pins.json")
}

func TestSaveLoad_Roundtrip(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	path := tempPinPath(t)
	if err := SavePins(env, []string{"DB_HOST", "PORT"}, path); err != nil {
		t.Fatalf("SavePins: %v", err)
	}
	pf, err := LoadPins(path)
	if err != nil {
		t.Fatalf("LoadPins: %v", err)
	}
	if len(pf.Pins) != 2 {
		t.Fatalf("expected 2 pins, got %d", len(pf.Pins))
	}
}

func TestSavePins_MissingKey(t *testing.T) {
	env := map[string]string{"A": "1"}
	err := SavePins(env, []string{"MISSING"}, tempPinPath(t))
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestCheck_NoDrift(t *testing.T) {
	pf := PinFile{Pins: []Pin{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}}
	env := map[string]string{"A": "1", "B": "2"}
	results := Check(pf, env)
	for _, r := range results {
		if r.Drifted || r.Missing {
			t.Errorf("unexpected drift for key %s", r.Key)
		}
	}
}

func TestCheck_Drifted(t *testing.T) {
	pf := PinFile{Pins: []Pin{{Key: "A", Value: "old"}}}
	env := map[string]string{"A": "new"}
	results := Check(pf, env)
	if !results[0].Drifted {
		t.Error("expected drift")
	}
}

func TestCheck_Missing(t *testing.T) {
	pf := PinFile{Pins: []Pin{{Key: "GONE", Value: "x"}}}
	results := Check(pf, map[string]string{})
	if !results[0].Missing {
		t.Error("expected missing")
	}
}

func TestHasDrift_True(t *testing.T) {
	results := []Result{{Key: "X", Drifted: true}}
	if !HasDrift(results) {
		t.Error("expected drift")
	}
}

func TestHasDrift_False(t *testing.T) {
	results := []Result{{Key: "X", Pinned: "v", Actual: "v"}}
	if HasDrift(results) {
		t.Error("expected no drift")
	}
}

func TestLoadPins_InvalidJSON(t *testing.T) {
	p := tempPinPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0644)
	_, err := LoadPins(p)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadPins_MissingFile(t *testing.T) {
	_, err := LoadPins("/nonexistent/pins.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSaveLoad_PinValues(t *testing.T) {
	env := map[string]string{"SECRET": "abc123"}
	p := tempPinPath(t)
	_ = SavePins(env, []string{"SECRET"}, p)
	data, _ := os.ReadFile(p)
	var pf PinFile
	_ = json.Unmarshal(data, &pf)
	if pf.Pins[0].Value != "abc123" {
		t.Errorf("expected abc123, got %s", pf.Pins[0].Value)
	}
}
