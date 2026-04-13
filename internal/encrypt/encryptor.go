// Package encrypt provides AES-GCM encryption and decryption for .env values.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const encryptedPrefix = "enc:"

// deriveKey produces a 32-byte AES key from a passphrase using SHA-256.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// Encrypt encrypts plaintext using AES-GCM with the given passphrase.
// The result is base64-encoded and prefixed with "enc:".
func Encrypt(plaintext, passphrase string) (string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("encrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: create gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encrypt: generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a value previously encrypted with Encrypt.
// Returns an error if the value is not prefixed with "enc:".
func Decrypt(ciphertext, passphrase string) (string, error) {
	if !strings.HasPrefix(ciphertext, encryptedPrefix) {
		return "", errors.New("decrypt: value is not encrypted (missing 'enc:' prefix)")
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(ciphertext, encryptedPrefix))
	if err != nil {
		return "", fmt.Errorf("decrypt: base64 decode: %w", err)
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("decrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("decrypt: create gcm: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("decrypt: ciphertext too short")
	}
	plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: open gcm: %w", err)
	}
	return string(plaintext), nil
}

// IsEncrypted reports whether a value has the encrypted prefix.
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}

// EncryptMap encrypts all values in the provided map, returning a new map.
func EncryptMap(env map[string]string, passphrase string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		enc, err := Encrypt(v, passphrase)
		if err != nil {
			return nil, fmt.Errorf("encrypt map key %q: %w", k, err)
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptMap decrypts all encrypted values in the provided map, returning a new map.
func DecryptMap(env map[string]string, passphrase string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if !IsEncrypted(v) {
			out[k] = v
			continue
		}
		dec, err := Decrypt(v, passphrase)
		if err != nil {
			return nil, fmt.Errorf("decrypt map key %q: %w", k, err)
		}
		out[k] = dec
	}
	return out, nil
}
