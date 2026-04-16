package reorder_test

import (
	"os"
	"testing"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/reorder"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "reorder-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestReorderIntegration_ParseAndReorder(t *testing.T) {
	path := writeTempEnv(t, "ZEBRA=z\nAPPLE=a\nMANGO=m\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	res, err := reorder.Reorder(env, reorder.Options{Strategy: reorder.StrategyAlpha})
	if err != nil {
		t.Fatalf("reorder error: %v", err)
	}
	if res.Ordered[0] != "APPLE" {
		t.Errorf("expected APPLE first, got %s", res.Ordered[0])
	}
	if res.Env["ZEBRA"] != "z" {
		t.Errorf("value mismatch for ZEBRA")
	}
}

func TestReorderIntegration_TemplateOrder(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=prod\nSECRET_KEY=abc\n")
	env, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	order := []string{"APP_ENV", "DB_HOST", "DB_PORT"}
	res, err := reorder.Reorder(env, reorder.Options{
		Strategy:       reorder.StrategyTemplate,
		Order:          order,
		PutUnknownLast: true,
	})
	if err != nil {
		t.Fatalf("reorder error: %v", err)
	}
	if res.Ordered[0] != "APP_ENV" {
		t.Errorf("expected APP_ENV first, got %s", res.Ordered[0])
	}
	if len(res.Unknown) != 1 || res.Unknown[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY as unknown, got %v", res.Unknown)
	}
}
