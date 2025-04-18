package store

import "fmt"

// SecretStore defines the interface for storing and retrieving encrypted secrets.
type SecretStore interface {
	// Init initializes the storage backend (e.g., connects to DB, creates file).
	Init() error
	// Close closes the storage backend resources (e.g., database connection).
	Close() error

	// Create stores a new encrypted value associated with a key.
	// Returns an error if the key already exists or storage fails.
	Create(key string, encryptedValue []byte) error
	// Read retrieves the encrypted value for a given key.
	// Returns an error if the key is not found or storage fails.
	Read(key string) ([]byte, error)
	// Update updates the encrypted value for an existing key.
	// Returns an error if the key is not found or storage fails.
	Update(key string, encryptedValue []byte) error
	// Delete removes a secret by its key.
	// Returns an error if the key is not found or storage fails.
	Delete(key string) error
	// ListKeys retrieves all available secret keys.
	ListKeys() ([]string, error)
}

// Common errors
var (
	ErrSecretNotFound       = fmt.Errorf("secret not found")
	ErrSecretAlreadyExists  = fmt.Errorf("secret already exists")
	ErrInvalidConfiguration = fmt.Errorf("invalid store configuration")
)
