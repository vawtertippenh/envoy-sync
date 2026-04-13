package encrypt_test

import (
	"strings"
	"testing"

	"envoy-sync/internal/encrypt"
)

const testPassphrase = "super-secret-passphrase"

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	plaintext := "my-secret-value"
	enc, err := encrypt.Encrypt(plaintext, testPassphrase)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	if !strings.HasPrefix(enc, "enc:") {
		t.Errorf("expected 'enc:' prefix, got: %s", enc)
	}
	dec, err := encrypt.Decrypt(enc, testPassphrase)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if dec != plaintext {
		t.Errorf("expected %q, got %q", plaintext, dec)
	}
}

func TestEncrypt_ProducesUniqueValues(t *testing.T) {
	enc1, _ := encrypt.Encrypt("value", testPassphrase)
	enc2, _ := encrypt.Encrypt("value", testPassphrase)
	if enc1 == enc2 {
		t.Error("expected unique ciphertexts due to random nonce, got identical values")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	enc, _ := encrypt.Encrypt("secret", testPassphrase)
	_, err := encrypt.Decrypt(enc, "wrong-passphrase")
	if err == nil {
		t.Error("expected error when decrypting with wrong passphrase")
	}
}

func TestDecrypt_MissingPrefix(t *testing.T) {
	_, err := encrypt.Decrypt("plaintext-no-prefix", testPassphrase)
	if err == nil {
		t.Error("expected error for value without 'enc:' prefix")
	}
}

func TestIsEncrypted(t *testing.T) {
	enc, _ := encrypt.Encrypt("hello", testPassphrase)
	if !encrypt.IsEncrypted(enc) {
		t.Error("expected IsEncrypted to return true for encrypted value")
	}
	if encrypt.IsEncrypted("plain-value") {
		t.Error("expected IsEncrypted to return false for plain value")
	}
}

func TestEncryptMap_DecryptMap_Roundtrip(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"APP_NAME":    "envoy-sync",
	}
	encrypted, err := encrypt.EncryptMap(env, testPassphrase)
	if err != nil {
		t.Fatalf("EncryptMap error: %v", err)
	}
	for k, v := range encrypted {
		if !encrypt.IsEncrypted(v) {
			t.Errorf("key %q: expected encrypted value, got %q", k, v)
		}
	}
	decrypted, err := encrypt.DecryptMap(encrypted, testPassphrase)
	if err != nil {
		t.Fatalf("DecryptMap error: %v", err)
	}
	for k, want := range env {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: expected %q, got %q", k, want, got)
		}
	}
}

func TestDecryptMap_SkipsPlainValues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "envoy-sync",
		"PORT":     "8080",
	}
	decrypted, err := encrypt.DecryptMap(env, testPassphrase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decrypted["APP_NAME"] != "envoy-sync" {
		t.Errorf("expected plain value to pass through unchanged")
	}
}
