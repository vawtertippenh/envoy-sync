// Package patch provides functionality to apply key-value patches to env maps.
package patch

import (
	"fmt"
	"strings"
)

// Op represents a patch operation type.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
	OpRename Op = "rename"
)

// Instruction describes a single patch operation.
type Instruction struct {
	Op      Op
	Key     string
	Value   string // used by OpSet
	NewKey  string // used by OpRename
}

// Result holds the outcome of applying a patch.
type Result struct {
	Env     map[string]string
	Applied []string
	Skipped []string
}

// ParseInstruction parses a patch instruction string of the form:
//   set:KEY=VALUE
//   delete:KEY
//   rename:OLD_KEY->NEW_KEY
func ParseInstruction(s string) (Instruction, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return Instruction{}, fmt.Errorf("invalid instruction %q: missing op prefix", s)
	}
	op, payload := Op(parts[0]), parts[1]
	switch op {
	case OpSet:
		kv := strings.SplitN(payload, "=", 2)
		if len(kv) != 2 {
			return Instruction{}, fmt.Errorf("set instruction missing '=': %q", s)
		}
		return Instruction{Op: OpSet, Key: kv[0], Value: kv[1]}, nil
	case OpDelete:
		if payload == "" {
			return Instruction{}, fmt.Errorf("delete instruction missing key: %q", s)
		}
		return Instruction{Op: OpDelete, Key: payload}, nil
	case OpRename:
		pair := strings.SplitN(payload, "->", 2)
		if len(pair) != 2 || pair[0] == "" || pair[1] == "" {
			return Instruction{}, fmt.Errorf("rename instruction must be OLD->NEW: %q", s)
		}
		return Instruction{Op: OpRename, Key: pair[0], NewKey: pair[1]}, nil
	default:
		return Instruction{}, fmt.Errorf("unknown op %q in instruction %q", op, s)
	}
}

// Apply applies a slice of Instructions to a copy of env and returns a Result.
func Apply(env map[string]string, instructions []Instruction) Result {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var applied, skipped []string

	for _, inst := range instructions {
		switch inst.Op {
		case OpSet:
			out[inst.Key] = inst.Value
			applied = append(applied, fmt.Sprintf("set %s", inst.Key))
		case OpDelete:
			if _, ok := out[inst.Key]; ok {
				delete(out, inst.Key)
				applied = append(applied, fmt.Sprintf("delete %s", inst.Key))
			} else {
				skipped = append(skipped, fmt.Sprintf("delete %s (not found)", inst.Key))
			}
		case OpRename:
			if val, ok := out[inst.Key]; ok {
				out[inst.NewKey] = val
				delete(out, inst.Key)
				applied = append(applied, fmt.Sprintf("rename %s->%s", inst.Key, inst.NewKey))
			} else {
				skipped = append(skipped, fmt.Sprintf("rename %s (not found)", inst.Key))
			}
		}
	}

	return Result{Env: out, Applied: applied, Skipped: skipped}
}
