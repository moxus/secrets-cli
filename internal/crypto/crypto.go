package crypto

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	// NonceSize is the required nonce size for nacl/secretbox (24 bytes).
	NonceSize = 24
	// TagSize is the size of the authentication tag appended by Seal.
	// The nacl/secretbox package uses a 16-byte authentication tag,
	// though this isn't exposed as a public constant.
	TagSize = 16
)

// Encrypt encrypts plaintext using the provided secretbox key.
// It prepends a random nonce to the ciphertext and returns the combined []byte.
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size for encryption: expected 32 bytes, got %d", len(key))
	}

	var nonce [NonceSize]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Seal appends the tag and returns nonce || ciphertext || tag
	// We prepend the nonce manually so we can easily split it off later
	ciphertext := secretbox.Seal(nonce[:], plaintext, &nonce, (*[32]byte)(key))

	return ciphertext, nil
}

// Decrypt decrypts ciphertext (with prepended nonce) using the provided secretbox key.
// It returns the original plaintext.
func Decrypt(ciphertextWithNonce []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size for decryption: expected 32 bytes, got %d", len(key))
	}

	if len(ciphertextWithNonce) < NonceSize+TagSize {
		return nil, fmt.Errorf("ciphertext is too short to contain nonce and tag")
	}

	var nonce [NonceSize]byte
	copy(nonce[:], ciphertextWithNonce[:NonceSize])

	// The actual ciphertext with tag starts after the nonce
	ciphertext := ciphertextWithNonce[NonceSize:]

	// Open verifies the tag and decrypts. It returns the plaintext or an error.
	plaintext, ok := secretbox.Open(nil, ciphertext, &nonce, (*[32]byte)(key))
	if !ok {
		return nil, fmt.Errorf("decryption failed (authentication tag mismatch or corrupted data)")
	}

	return plaintext, nil
}
