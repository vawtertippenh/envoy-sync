// Package convert provides utilities to convert env maps between
// different formats such as YAML, TOML-style key=value, and JSON.
package convert

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format represents a supported conversion target format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatYAML   Format = "yaml"
	FormatExport Format = "export"
)

// Result holds the converted output string and the format used.
type Result struct {
	Format Format
	Output string
}

// Convert transforms an env map into the specified format string.
func Convert(env map[string]string, format Format) (Result, error) {
	switch format {
	case FormatDotenv:
		return Result{Format: format, Output: toDotenv(env)}, nil
	case FormatJSON:
		out, err := toJSON(env)
		if err != nil {
			return Result{}, err
		}
		return Result{Format: format, Output: out}, nil
	case FormatYAML:
		return Result{Format: format, Output: toYAML(env)}, nil
	case FormatExport:
		return Result{Format: format, Output: toExport(env)}, nil
	default:
		return Result{}, fmt.Errorf("unsupported format: %q", format)
	}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func toDotenv(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}
	return sb.String()
}

func toJSON(env map[string]string) (string, error) {
	b, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func toYAML(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "%s: \"%s\"\n", k, env[k])
	}
	return sb.String()
}

func toExport(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}
