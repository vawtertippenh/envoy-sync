package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
}

// Parse reads and parses a .env file at the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	ef := &EnvFile{Path: path}
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		// Skip blank lines
		if trimmed == "" {
			continue
		}

		// Full-line comment
		if strings.HasPrefix(trimmed, "#") {
			ef.Entries = append(ef.Entries, Entry{
				Comment: trimmed,
				Line:    lineNum,
			})
			continue
		}

		key, value, found := strings.Cut(trimmed, "=")
		if !found {
			return nil, fmt.Errorf("line %d: invalid format (missing '='): %q", lineNum, raw)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = stripQuotes(value)

		ef.Entries = append(ef.Entries, Entry{
			Key:   key,
			Value: value,
			Line:  lineNum,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return ef, nil
}

// ToMap converts the EnvFile entries to a key→value map (comments excluded).
func (ef *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		if e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	return m
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
