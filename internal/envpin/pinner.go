// Package envpin provides functionality to pin env keys to specific values
// and detect when those values drift from their pinned state.
package envpin

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Pin represents a single pinned key-value pair.
type Pin struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PinFile holds a set of pinned entries.
type PinFile struct {
	Pins []Pin `json:"pins"`
}

// Result holds the outcome of a pin check for a single key.
type Result struct {
	Key      string
	Pinned   string
	Actual   string
	Missing  bool
	Drifted  bool
}

// SavePins writes a pin file from the given env map (subset of keys).
func SavePins(env map[string]string, keys []string, path string) error {
	pf := PinFile{}
	for _, k := range keys {
		v, ok := env[k]
		if !ok {
			return fmt.Errorf("key %q not found in env", k)
		}
		pf.Pins = append(pf.Pins, Pin{Key: k, Value: v})
	}
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPins reads a pin file from disk.
func LoadPins(path string) (PinFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PinFile{}, fmt.Errorf("could not read pin file: %w", err)
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return PinFile{}, fmt.Errorf("invalid pin file: %w", err)
	}
	return pf, nil
}

// Check compares pinned values against the provided env map.
func Check(pf PinFile, env map[string]string) []Result {
	results := make([]Result, 0, len(pf.Pins))
	for _, p := range pf.Pins {
		actual, ok := env[p.Key]
		r := Result{Key: p.Key, Pinned: p.Value, Actual: actual}
		if !ok {
			r.Missing = true
		} else if actual != p.Value {
			r.Drifted = true
		}
		results = append(results, r)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Key < results[j].Key })
	return results
}

// HasDrift returns true if any result is drifted or missing.
func HasDrift(results []Result) bool {
	for _, r := range results {
		if r.Drifted || r.Missing {
			return true
		}
	}
	return false
}
