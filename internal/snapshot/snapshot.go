// Package snapshot provides functionality to capture and compare .env file
// snapshots over time, enabling drift detection between saved states.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an env map.
type Snapshot struct {
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

// Diff describes the difference between two snapshots.
type Diff struct {
	Added   map[string]string `json:"added"`
	Removed map[string]string `json:"removed"`
	Changed map[string]Change `json:"changed"`
}

// Change holds the before/after values for a modified key.
type Change struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

// Take creates a new Snapshot from the given env map and label.
func Take(label string, env map[string]string) Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return Snapshot{
		Label:     label,
		Timestamp: time.Now().UTC(),
		Env:       copy,
	}
}

// Save writes a snapshot to the given file path as JSON.
func Save(s Snapshot, path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return s, nil
}

// Compare returns the Diff between snapshot a (baseline) and snapshot b (current).
func Compare(a, b Snapshot) Diff {
	d := Diff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]Change),
	}
	for k, v := range b.Env {
		if old, ok := a.Env[k]; !ok {
			d.Added[k] = v
		} else if old != v {
			d.Changed[k] = Change{Before: old, After: v}
		}
	}
	for k, v := range a.Env {
		if _, ok := b.Env[k]; !ok {
			d.Removed[k] = v
		}
	}
	return d
}

// HasDrift reports whether the Diff contains any changes.
func HasDrift(d Diff) bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
