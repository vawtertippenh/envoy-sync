package patch

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
		"DEBUG":    "false",
	}
}

func TestParseInstruction_Set(t *testing.T) {
	inst, err := ParseInstruction("set:FOO=bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inst.Op != OpSet || inst.Key != "FOO" || inst.Value != "bar" {
		t.Errorf("unexpected instruction: %+v", inst)
	}
}

func TestParseInstruction_Delete(t *testing.T) {
	inst, err := ParseInstruction("delete:PORT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inst.Op != OpDelete || inst.Key != "PORT" {
		t.Errorf("unexpected instruction: %+v", inst)
	}
}

func TestParseInstruction_Rename(t *testing.T) {
	inst, err := ParseInstruction("rename:DEBUG->VERBOSE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inst.Op != OpRename || inst.Key != "DEBUG" || inst.NewKey != "VERBOSE" {
		t.Errorf("unexpected instruction: %+v", inst)
	}
}

func TestParseInstruction_Invalid(t *testing.T) {
	cases := []string{
		"nocolon",
		"set:NOEQUALSSIGN",
		"delete:",
		"rename:NOARROW",
		"unknown:KEY=VAL",
	}
	for _, c := range cases {
		_, err := ParseInstruction(c)
		if err == nil {
			t.Errorf("expected error for %q but got none", c)
		}
	}
}

func TestApply_Set(t *testing.T) {
	result := Apply(baseEnv(), []Instruction{{Op: OpSet, Key: "NEW_KEY", Value: "hello"}})
	if result.Env["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", result.Env["NEW_KEY"])
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(result.Applied))
	}
}

func TestApply_Delete(t *testing.T) {
	result := Apply(baseEnv(), []Instruction{{Op: OpDelete, Key: "PORT"}})
	if _, ok := result.Env["PORT"]; ok {
		t.Error("expected PORT to be deleted")
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(result.Applied))
	}
}

func TestApply_DeleteMissingKey(t *testing.T) {
	result := Apply(baseEnv(), []Instruction{{Op: OpDelete, Key: "NONEXISTENT"}})
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
}

func TestApply_Rename(t *testing.T) {
	result := Apply(baseEnv(), []Instruction{{Op: OpRename, Key: "DEBUG", NewKey: "VERBOSE"}})
	if _, ok := result.Env["DEBUG"]; ok {
		t.Error("expected DEBUG to be removed after rename")
	}
	if result.Env["VERBOSE"] != "false" {
		t.Errorf("expected VERBOSE=false, got %q", result.Env["VERBOSE"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	Apply(env, []Instruction{{Op: OpDelete, Key: "PORT"}})
	if _, ok := env["PORT"]; !ok {
		t.Error("original env was mutated")
	}
}
