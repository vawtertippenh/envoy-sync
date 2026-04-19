package envmerge_test

import (
	"os"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envmerge"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envmerge-*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestMergeIntegration_TwoFiles(t *testing.T) {
	pA := writeTempEnv(t, "FOO=from_a\nSHARED=a\n")
	pB := writeTempEnv(t, "BAR=from_b\nSHARED=b\n")

	envA, err := envfile.Parse(pA)
	if err != nil {
		t.Fatal(err)
	}
	envB, err := envfile.Parse(pB)
	if err != nil {
		t.Fatal(err)
	}

	res, err := envmerge.Merge([]map[string]string{envA, envB}, envmerge.Options{Strategy: envmerge.StrategyLast})
	if err != nil {
		t.Fatal(err)
	}
	if res.Env["FOO"] != "from_a" {
		t.Errorf("expected from_a, got %s", res.Env["FOO"])
	}
	if res.Env["BAR"] != "from_b" {
		t.Errorf("expected from_b, got %s", res.Env["BAR"])
	}
	if res.Env["SHARED"] != "b" {
		t.Errorf("expected b, got %s", res.Env["SHARED"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMergeIntegration_EmptyFile(t *testing.T) {
	pA := writeTempEnv(t, "KEY=val\n")
	pB := writeTempEnv(t, "")

	envA, _ := envfile.Parse(pA)
	envB, _ := envfile.Parse(pB)

	res, err := envmerge.Merge([]map[string]string{envA, envB}, envmerge.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Env["KEY"] != "val" {
		t.Errorf("expected val")
	}
}
