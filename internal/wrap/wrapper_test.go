package wrap

import (
	"strings"
	"testing"
)

func TestWrap_ShortValues_NoWrap(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	lines := Wrap(env, Options{Width: 80})
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if strings.Contains(l, `\`) {
			t.Errorf("unexpected continuation in short line: %s", l)
		}
	}
}

func TestWrap_LongValue_SplitsCorrectly(t *testing.T) {
	env := map[string]string{
		"MY_KEY": strings.Repeat("A", 100),
	}
	lines := Wrap(env, Options{Width: 40})
	if len(lines) < 2 {
		t.Fatalf("expected multiple lines, got %d", len(lines))
	}
	// All lines except the last should end with continuation char.
	for i, l := range lines[:len(lines)-1] {
		if !strings.HasSuffix(l, `\`) {
			t.Errorf("line %d missing continuation: %q", i, l)
		}
	}
	// Last line should NOT end with continuation.
	last := lines[len(lines)-1]
	if strings.HasSuffix(last, `\`) {
		t.Errorf("last line should not have continuation: %q", last)
	}
}

func TestWrap_ContinuationLines_HaveIndent(t *testing.T) {
	env := map[string]string{
		"K": strings.Repeat("B", 120),
	}
	lines := Wrap(env, Options{Width: 40, Indent: "  "})
	for _, l := range lines[1:] {
		if !strings.HasPrefix(l, "  ") && !strings.HasPrefix(l, "  ") {
			t.Errorf("continuation line missing indent: %q", l)
		}
	}
}

func TestWrap_ReassemblesOriginalValue(t *testing.T) {
	original := strings.Repeat("XY", 60)
	env := map[string]string{"DATA": original}
	lines := Wrap(env, Options{Width: 30, Indent: ""})

	// Strip key= prefix from first line and continuation chars, then join.
	var parts []string
	for i, l := range lines {
		if i == 0 {
			// Remove "DATA=" prefix
			l = strings.TrimPrefix(l, "DATA=")
		}
		l = strings.TrimSuffix(l, `\`)
		parts = append(parts, l)
	}
	reassembled := strings.Join(parts, "")
	if reassembled != original {
		t.Errorf("reassembled value mismatch\ngot: %s\nwant: %s", reassembled, original)
	}
}

func TestWrap_DefaultOptions(t *testing.T) {
	env := map[string]string{"A": "short"}
	lines := Wrap(env, Options{})
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != "A=short" {
		t.Errorf("unexpected line: %q", lines[0])
	}
}

func TestRender_JoinsWithNewlines(t *testing.T) {
	lines := []string{"A=1", "B=2", "C=3"}
	out := Render(lines)
	if out != "A=1\nB=2\nC=3" {
		t.Errorf("unexpected render output: %q", out)
	}
}

func TestWrap_EmptyEnv(t *testing.T) {
	lines := Wrap(map[string]string{}, Options{})
	if len(lines) != 0 {
		t.Errorf("expected no lines for empty env, got %d", len(lines))
	}
}
