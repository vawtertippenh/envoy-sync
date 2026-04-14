// Package profile provides support for named environment profiles,
// allowing users to manage multiple env configurations (e.g. dev, staging, prod).
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
)

// Profile represents a named set of environment variables.
type Profile struct {
	Name string            `json:"name"`
	Env  map[string]string `json:"env"`
}

// Store holds multiple named profiles persisted to a JSON file.
type Store struct {
	Profiles map[string]*Profile `json:"profiles"`
}

// NewStore returns an empty profile store.
func NewStore() *Store {
	return &Store{Profiles: make(map[string]*Profile)}
}

// LoadStore reads a store from the given file path.
// If the file does not exist, an empty store is returned.
func LoadStore(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return NewStore(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("profile: read store: %w", err)
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("profile: parse store: %w", err)
	}
	if s.Profiles == nil {
		s.Profiles = make(map[string]*Profile)
	}
	return &s, nil
}

// SaveStore writes the store to the given file path.
func SaveStore(path string, s *Store) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("profile: marshal store: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("profile: write store: %w", err)
	}
	return nil
}

// Set adds or replaces a profile in the store.
func (s *Store) Set(name string, env map[string]string) {
	s.Profiles[name] = &Profile{Name: name, Env: env}
}

// Get retrieves a profile by name.
func (s *Store) Get(name string) (*Profile, bool) {
	p, ok := s.Profiles[name]
	return p, ok
}

// Delete removes a profile by name.
func (s *Store) Delete(name string) bool {
	_, ok := s.Profiles[name]
	delete(s.Profiles, name)
	return ok
}

// List returns sorted profile names.
func (s *Store) List() []string {
	names := make([]string, 0, len(s.Profiles))
	for n := range s.Profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
