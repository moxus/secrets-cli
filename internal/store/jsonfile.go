package store

import (
	"encoding/base64" // <--- Add this line
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync" // For potential future concurrency needs
)

// JSONFileStore implements the SecretStore interface using a simple JSON file.
type JSONFileStore struct {
	FilePath string
	mu       sync.Mutex // Protects file access
	// Store secrets as map[plaintext_key] -> encrypted_value_base64
	// Storing as base64 in JSON makes it more readable,
	// but requires base64 encoding/decoding during save/load.
	// Alternatively, could store raw bytes if JSON supports it well (less standard).
	// Let's stick to base64 for JSON compatibility.
}

// NewJSONFileStore creates a new JSONFileStore instance.
func NewJSONFileStore(filePath string) (*JSONFileStore, error) {
	if filePath == "" {
		return nil, fmt.Errorf("%w: JSON file path cannot be empty", ErrInvalidConfiguration)
	}
	return &JSONFileStore{FilePath: filePath}, nil
}

// Init ensures the file exists (creates empty JSON object if not).
func (s *JSONFileStore) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.FilePath); os.IsNotExist(err) {
		// File doesn't exist, create it with an empty JSON object
		err := os.WriteFile(s.FilePath, []byte("{}"), 0600) // Owner read/write
		if err != nil {
			return fmt.Errorf("failed to create JSON file: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat JSON file: %w", err)
	}
	return nil
}

// Close does nothing for a file-based store.
func (s *JSONFileStore) Close() error {
	return nil // No resources to close
}

// loadData reads and unmarshals the JSON file.
func (s *JSONFileStore) loadData() (map[string][]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := make(map[string][]byte)

	// Read file content
	content, err := os.ReadFile(s.FilePath)
	if (err != nil) {
		// If file doesn't exist, treat as empty store
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	// If file is empty or contains only whitespace, treat as empty JSON object
	if len(content) == 0 || len([]byte(string(content))) == 0 {
		return data, nil
	}

	// Unmarshal JSON (reading base64 strings)
	var base64Data map[string]string
	if err := json.Unmarshal(content, &base64Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}


	// Decode base64 values
	for k, v := range base64Data {
		decodedValue, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 value for key '%s': %w", k, err)
		}
		data[k] = decodedValue
	}

	return data, nil
}

// saveData marshals and writes data to the JSON file.
func (s *JSONFileStore) saveData(data map[string][]byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Encode values to base64 for JSON
	base64Data := make(map[string]string)
	for k, v := range data {
		base64Data[k] = base64.StdEncoding.EncodeToString(v)
	}

	// Marshal data to JSON
	content, err := json.MarshalIndent(base64Data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	// Write to a temporary file in the same directory as the target file and rename for atomic update
	tmpFile, err := os.CreateTemp(filepath.Dir(s.FilePath), "secrets-json-")
	if err != nil {
		return fmt.Errorf("failed to create temp file for JSON save: %w", err)
	}
	tmpFilePath := tmpFile.Name()

	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFilePath) // Clean up temp file
		return fmt.Errorf("failed to write to temp JSON file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFilePath)
		return fmt.Errorf("failed to close temp JSON file: %w", err)
	}

	// Ensure permissions are restrictive on the temp file before rename
	if err := os.Chmod(tmpFilePath, 0600); err != nil {
		os.Remove(tmpFilePath)
		return fmt.Errorf("failed to set permissions on temp JSON file: %w", err)
	}

	// Rename temp file to replace original
	if err := os.Rename(tmpFilePath, s.FilePath); err != nil {
		os.Remove(tmpFilePath)
		return fmt.Errorf("failed to rename temp JSON file: %w", err)
	}

	return nil
}

// Create stores a new encrypted value.
func (s *JSONFileStore) Create(key string, encryptedValue []byte) error {
	data, err := s.loadData()
	if err != nil {
		return err
	}

	if _, exists := data[key]; exists {
		return fmt.Errorf("%w: secret with key '%s'", ErrSecretAlreadyExists, key)
	}

	data[key] = encryptedValue
	return s.saveData(data)
}

// Read retrieves an encrypted value.
func (s *JSONFileStore) Read(key string) ([]byte, error) {
	data, err := s.loadData()
	if err != nil {
		return nil, err
	}

	value, exists := data[key]
	if !exists {
		return nil, fmt.Errorf("%w: secret with key '%s'", ErrSecretNotFound, key)
	}

	return value, nil
}

// Update updates an existing encrypted value.
func (s *JSONFileStore) Update(key string, encryptedValue []byte) error {
	data, err := s.loadData()
	if err != nil {
		return err
	}

	if _, exists := data[key]; !exists {
		return fmt.Errorf("%w: secret with key '%s'", ErrSecretNotFound, key)
	}

	data[key] = encryptedValue
	return s.saveData(data)
}

// Delete removes a secret.
func (s *JSONFileStore) Delete(key string) error {
	data, err := s.loadData()
	if err != nil {
		return err
	}

	if _, exists := data[key]; !exists {
		return fmt.Errorf("%w: secret with key '%s'", ErrSecretNotFound, key)
	}

	delete(data, key)
	return s.saveData(data)
}

// ListKeys lists all available keys.
func (s *JSONFileStore) ListKeys() ([]string, error) {
	data, err := s.loadData()
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return keys, nil
}
