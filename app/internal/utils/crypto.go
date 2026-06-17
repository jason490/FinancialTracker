package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"github.com/labstack/gommon/log"
)

// Encrypt encrypts a string using AES-GCM and returns a base64 encoded string.
// It requires an ENCRYPTION_KEY environment variable.
func Encrypt(plaintext string) (string, error) {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		keyStr = "default-32-byte-secret-key-prod!" // exactly 32 bytes for AES-256
		log.Warn("No Encryption in env file!!!")
	}
	key := []byte(keyStr)
	if len(key) != 32 {
		return "", errors.New("ENCRYPTION_KEY must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	// RawURLEncoding avoids '=' padding that OAuth redirect URLs may strip.
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64 encoded AES-GCM ciphertext.
func Decrypt(cryptoText string) (string, error) {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		keyStr = "default-32-byte-secret-key-prod!"
		log.Warn("No Encryption in env file!!!")
	}
	key := []byte(keyStr)
	if len(key) != 32 {
		return "", errors.New("ENCRYPTION_KEY must be 32 bytes for AES-256")
	}

	ciphertext, err := base64.RawURLEncoding.DecodeString(cryptoText)
	if err != nil {
		ciphertext, err = base64.URLEncoding.DecodeString(cryptoText)
		if err != nil {
			return "", err
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
