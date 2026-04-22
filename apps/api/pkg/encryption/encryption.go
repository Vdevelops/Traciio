package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidKey = errors.New("encryption key must be 32 bytes (AES-256)")
)

// Encrypt encrypts plaintext using AES-256-GCM
// Key can be either raw 32 bytes or base64-encoded 32 bytes
func Encrypt(plaintext, key string) (string, error) {
	// Try to decode as base64 first (common case)
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil || len(keyBytes) != 32 {
		// If base64 decode fails or length is wrong, try as raw bytes
		keyBytes = []byte(key)
		if len(keyBytes) != 32 {
			return "", ErrInvalidKey
		}
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
// Key can be either raw 32 bytes or base64-encoded 32 bytes
func Decrypt(ciphertext, key string) (string, error) {
	// Try to decode as base64 first (common case)
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil || len(keyBytes) != 32 {
		// If base64 decode fails or length is wrong, try as raw bytes
		keyBytes = []byte(key)
		if len(keyBytes) != 32 {
			return "", ErrInvalidKey
		}
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// GenerateKey generates a random 32-byte key for AES-256
func GenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

