// Package flatten provides utilities for flattening nested env-like
// structures (e.g. JSON objects) into flat KEY=VALUE pairs, with
// configurable separator and prefix support.
package flatten

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Options controls how flattening behaves.
type Options struct {
	// Separator between nested key segments (default "_").
	Separator string
	// Prefix prepended to every output key.
	Prefix string
	// UpperCase converts all keys to upper-case.
	UpperCase bool
}

// Flatten converts a nested map (parsed from JSON) into a flat
// map[string]string suitable for use as an env file.
func Flatten(input map[string]interface{}, opts Options) (map[string]string, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	out := make(map[string]string)
	if err := flatten(input, opts.Prefix, opts.Separator, opts.UpperCase, out); err != nil {
		return nil, err
	}
	return out, nil
}

// FlattenJSON parses raw JSON and flattens it into a flat env map.
func FlattenJSON(data []byte, opts Options) (map[string]string, error) {
	var nested map[string]interface{}
	if err := json.Unmarshal(data, &nested); err != nil {
		return nil, fmt.Errorf("flatten: invalid JSON: %w", err)
	}
	return Flatten(nested, opts)
}

// Render serialises a flat map as dotenv-style lines, sorted by key.
func Render(flat map[string]string) string {
	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		v := flat[k]
		if strings.ContainsAny(v, " \t\n") {
			v = `"` + v + `"`
		}
		sb.WriteString(k + "=" + v + "\n")
	}
	return sb.String()
}

func flatten(node map[string]interface{}, prefix, sep string, upper bool, out map[string]string) error {
	for k, v := range node {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		if upper {
			key = strings.ToUpper(key)
		}
		switch val := v.(type) {
		case map[string]interface{}:
			if err := flatten(val, key, sep, upper, out); err != nil {
				return err
			}
		case string:
			out[key] = val
		case bool:
			out[key] = fmt.Sprintf("%t", val)
		case float64:
			out[key] = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", val), "0"), ".")
		case nil:
			out[key] = ""
		default:
			out[key] = fmt.Sprintf("%v", val)
		}
	}
	return nil
}
