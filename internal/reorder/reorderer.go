package reorder

import "sort"

// Strategy defines how keys should be reordered.
type Strategy string

const (
	StrategyAlpha    Strategy = "alpha"
	StrategyTemplate Strategy = "template"
	StrategyCustom   Strategy = "custom"
)

// Options configures the reorder operation.
type Options struct {
	Strategy Strategy
	// Order is used with StrategyCustom or StrategyTemplate to define key order.
	Order []string
	// PutUnknownLast places keys not in Order at the end (template/custom only).
	PutUnknownLast bool
}

// Result holds the reordered keys and metadata.
type Result struct {
	Ordered []string
	Env     map[string]string
	Unknown []string // keys not present in Order (template/custom only)
}

// Reorder returns a Result with keys sorted according to opts.
func Reorder(env map[string]string, opts Options) (Result, error) {
	switch opts.Strategy {
	case StrategyAlpha:
		return reorderAlpha(env), nil
	case StrategyTemplate, StrategyCustom:
		return reorderByOrder(env, opts), nil
	default:
		return reorderAlpha(env), nil
	}
}

func reorderAlpha(env map[string]string) Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return Result{Ordered: keys, Env: copyMap(env)}
}

func reorderByOrder(env map[string]string, opts Options) Result {
	seen := make(map[string]bool)
	ordered := []string{}
	for _, k := range opts.Order {
		if _, ok := env[k]; ok {
			ordered = append(ordered, k)
			seen[k] = true
		}
	}
	unknown := []string{}
	for k := range env {
		if !seen[k] {
			unknown = append(unknown, k)
		}
	}
	sort.Strings(unknown)
	if opts.PutUnknownLast {
		ordered = append(ordered, unknown...)
	}
	return Result{Ordered: ordered, Env: copyMap(env), Unknown: unknown}
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
