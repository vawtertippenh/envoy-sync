package envcast

import (
	"fmt"
	"strconv"
	"strings"
)

// Type represents a target cast type.
type Type string

const (
	TypeString  Type = "string"
	TypeInt     Type = "int"
	TypeFloat   Type = "float"
	TypeBool    Type = "bool"
)

// Result holds the cast result for a single key.
type Result struct {
	Key      string
	Original string
	Casted   interface{}
	Type     Type
	Err      error
}

// Options configures the Cast operation.
type Options struct {
	// Rules maps key names to desired types.
	Rules map[string]Type
	// SkipUnknown ignores keys not in Rules.
	SkipUnknown bool
}

// Cast attempts to convert env values to typed Go values based on Rules.
func Cast(env map[string]string, opts Options) []Result {
	results := make([]Result, 0, len(env))
	for _, key := range sortedKeys(env) {
		val := env[key]
		t, ok := opts.Rules[key]
		if !ok {
			if opts.SkipUnknown {
				continue
			}
			t = TypeString
		}
		r := Result{Key: key, Original: val, Type: t}
		r.Casted, r.Err = castValue(val, t)
		results = append(results, r)
	}
	return results
}

func castValue(val string, t Type) (interface{}, error) {
	switch t {
	case TypeString:
		return val, nil
	case TypeInt:
		v, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			return nil, fmt.Errorf("cannot cast %q to int", val)
		}
		return v, nil
	case TypeFloat:
		v, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return nil, fmt.Errorf("cannot cast %q to float", val)
		}
		return v, nil
	case TypeBool:
		v, err := strconv.ParseBool(strings.TrimSpace(val))
		if err != nil {
			return nil, fmt.Errorf("cannot cast %q to bool", val)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown type %q", t)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
