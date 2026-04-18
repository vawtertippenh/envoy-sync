package envfmt

import (
	"strings"
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"db_host": "localhost",
		"DB_PORT": "5432",
		"api_key": "secret",
	}
}

func TestFormat_NoOptions(t *testing.T) {
	r := Format(baseEnv(), Style{})
	if r.Total != 3 {
		t.Fatalf("expected 3 total, got %d", r.Total)
	}
	if r.Changed != 0 {
		t.Fatalf("expected 0 changed, got %d", r.Changed)
	}
}

func TestFormat_UppercaseKeys(t *testing.T) {
	r := Format(map[string]string{"db_host": "localhost"}, Style{UppercaseKeys: true})
	if r.Changed != 1 {
		t.Fatalf("expected 1 changed, got %d", r.Changed)
	}
	if !strings.Contains(r.Lines[0], "DB_HOST") {
		t.Errorf("expected uppercased key, got %q", r.Lines[0])
	}
}

func TestFormat_QuoteAllValues(t *testing.T) {
	r := Format(map[string]string{"HOST": "localhost"}, Style{QuoteAllValues: true})
	if !strings.Contains(r.Lines[0], `"localhost"`) {
		t.Errorf("expected quoted value, got %q", r.Lines[0])
	}
}

func TestFormat_AlreadyQuoted_NoDoubleQuote(t *testing.T) {
	r := Format(map[string]string{"HOST": `"localhost"`}, Style{QuoteAllValues: true})
	if strings.Contains(r.Lines[0], `""localhost""`) {
		t.Errorf("double-quoted value detected: %q", r.Lines[0])
	}
}

func TestFormat_SpaceAroundEqual(t *testing.T) {
	r := Format(map[string]string{"KEY": "val"}, Style{SpaceAroundEqual: true})
	if !strings.Contains(r.Lines[0], "KEY = val") {
		t.Errorf("expected space around =, got %q", r.Lines[0])
	}
}

func TestFormat_SortKeys(t *testing.T) {
	env := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	r := Format(env, Style{SortKeys: true})
	if !strings.HasPrefix(r.Lines[0], "A_KEY") {
		t.Errorf("expected A_KEY first, got %q", r.Lines[0])
	}
}

func TestRender_TrailingNewline(t *testing.T) {
	r := Format(map[string]string{"X": "1"}, Style{})
	out := Render(r)
	if !strings.HasSuffix(out, "\n") {
		t.Errorf("expected trailing newline")
	}
}
