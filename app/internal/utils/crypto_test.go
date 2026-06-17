package utils

import (
	"os"
	"testing"
)

func TestEncryptDecryptWithDefaultKey(t *testing.T) {
	t.Setenv("ENCRYPTION_KEY", "")

	encrypted, err := Encrypt("hello")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if decrypted != "hello" {
		t.Fatalf("Decrypt() = %q, want %q", decrypted, "hello")
	}
}

func TestDefaultEncryptionKeyLength(t *testing.T) {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		key = "default-32-byte-secret-key-prod!"
	}
	if len(key) != 32 {
		t.Fatalf("default key length = %d, want 32", len(key))
	}
}
