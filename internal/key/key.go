package key

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

const (
	// EnvKeyName is the environment variable name where the encryption key is expected.
	EnvKeyName = "SECRETS_ENCRYPTION_KEY"
	// SecretBoxKeySize is the required key size for nacl/secretbox (32 bytes).
	SecretBoxKeySize = 32
)

// GenerateKey generates a new random 32-byte key suitable for nacl/secretbox.
// It returns the key as a base64 encoded string and an error.
func GenerateKey() (string, error) {
	key := make([]byte, SecretBoxKeySize)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// LoadKeyFromEnv loads the encryption key from the specified environment variable.
// It expects the key to be base64 encoded and returns the raw byte key.
func LoadKeyFromEnv() ([]byte, error) {
	keyBase64 := os.Getenv(EnvKeyName)
	if keyBase64 == "" {
		return nil, fmt.Errorf("encryption key environment variable '%s' is not set", EnvKeyName)
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key from environment: %w", err)
	}

	if len(key) < SecretBoxKeySize {
		// Extend key with zeros if too short
		extended := make([]byte, SecretBoxKeySize)
		copy(extended, key)
		key = extended
	} else if len(key) > SecretBoxKeySize {
		// Trim key if too long
		key = key[:SecretBoxKeySize]
	}

	if len(key) != SecretBoxKeySize {
		return nil, fmt.Errorf("invalid key: could not adjust to required size %d bytes", SecretBoxKeySize)
	}

	return key, nil
}

// --- Helper function to save/load from file (FOR DEMO/GENERATION ONLY) ---
// This is not recommended for production as discussed.
// You could use this helper in a 'generate-key' command.

func SaveKeyToFile(path string, key []byte) error {
	// Ensure file permissions are restrictive (owner read/write only)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key file for writing: %w", err)
	}
	defer file.Close()

	encodedKey := base64.StdEncoding.EncodeToString(key)
	_, err = file.WriteString(encodedKey)
	if err != nil {
		return fmt.Errorf("failed to write key to file: %w", err)
	}

	return nil
}

func LoadKeyFromFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	keyBase64 := string(content)
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key from file: %w", err)
	}

	if len(key) != SecretBoxKeySize {
		return nil, fmt.Errorf("invalid key size in file: expected %d bytes, got %d bytes after decoding", SecretBoxKeySize, len(key))
	}

	return key, nil
}
