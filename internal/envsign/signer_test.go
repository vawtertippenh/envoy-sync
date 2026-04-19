package envsign

import (
	"testing"
)

var baseEnv = map[string]string{
	"APP_NAME": "envoy",
	"APP_ENV":  "production",
	"SECRET":   "s3cr3t",
}

func TestSign_Deterministic(t *testing.T) {
	s1, err := Sign(baseEnv, "mykey", nil)
	if err != nil {
		t.Fatal(err)
	}
	s2, err := Sign(baseEnv, "mykey", nil)
	if err != nil {
		t.Fatal(err)
	}
	if s1.Digest != s2.Digest {
		t.Errorf("expected deterministic digest, got %s vs %s", s1.Digest, s2.Digest)
	}
}

func TestSign_DifferentSecret_DifferentDigest(t *testing.T) {
	s1, _ := Sign(baseEnv, "key1", nil)
	s2, _ := Sign(baseEnv, "key2", nil)
	if s1.Digest == s2.Digest {
		t.Error("expected different digests for different secrets")
	}
}

func TestSign_SubsetOfKeys(t *testing.T) {
	sig, err := Sign(baseEnv, "mykey", []string{"APP_NAME"})
	if err != nil {
		t.Fatal(err)
	}
	if len(sig.Keys) != 1 || sig.Keys[0] != "APP_NAME" {
		t.Errorf("unexpected keys: %v", sig.Keys)
	}
}

func TestSign_MissingKey_Error(t *testing.T) {
	_, err := Sign(baseEnv, "mykey", []string{"MISSING_KEY"})
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestSign_EmptySecret_Error(t *testing.T) {
	_, err := Sign(baseEnv, "", nil)
	if err == nil {
		t.Error("expected error for empty secret")
	}
}

func TestVerify_ValidSignature(t *testing.T) {
	sig, _ := Sign(baseEnv, "mykey", nil)
	ok, err := Verify(baseEnv, "mykey", sig)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected valid signature")
	}
}

func TestVerify_TamperedEnv(t *testing.T) {
	sig, _ := Sign(baseEnv, "mykey", nil)
	tampered := map[string]string{
		"APP_NAME": "hacked",
		"APP_ENV":  "production",
		"SECRET":   "s3cr3t",
	}
	ok, err := Verify(tampered, "mykey", sig)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected invalid signature for tampered env")
	}
}
