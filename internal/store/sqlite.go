package store

import (
	"database/sql"
	"errors"
	"fmt"

	 //_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
	//"github.com/mattn/go-sqlite3"    // Import sqlite3 for specific error codes
	_ "modernc.org/sqlite"
)

const (
	sqliteTableName = "secrets"
)

// SQLiteStore implements the SecretStore interface for a SQLite database.
type SQLiteStore struct {
	DBPath string
	db     *sql.DB // Database connection
}

// NewSQLiteStore creates a new SQLiteStore instance.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("%w: SQLite database path cannot be empty", ErrInvalidConfiguration)
	}
	return &SQLiteStore{DBPath: dbPath}, nil
}

// Init connects to the database and creates the table if it doesn't exist.
func (s *SQLiteStore) Init() error {
	dbConn, err := sql.Open("sqlite", s.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	// Pinging is a good practice to verify the connection
	if err = dbConn.Ping(); err != nil {
		dbConn.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}
	s.db = dbConn

	// Create table if not exists
	// 'key' is TEXT and UNIQUE NOT NULL (implicitly creates an index)
	// 'value' is BLOB to store the encrypted binary data
	query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            key TEXT UNIQUE NOT NULL,
            value BLOB NOT NULL
        )`, sqliteTableName)

	_, err = s.db.Exec(query)
	if err != nil {
		s.Close() // Close connection on error
		return fmt.Errorf("failed to create table '%s': %w", sqliteTableName, err)
	}

	return nil
}

// Close closes the database connection.
func (s *SQLiteStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Create stores a new encrypted value.
func (s *SQLiteStore) Create(key string, encryptedValue []byte) error {
	query := fmt.Sprintf("INSERT INTO %s (key, value) VALUES (?, ?)", sqliteTableName)
	_, err := s.db.Exec(query, key, encryptedValue)

	if err != nil {
		return fmt.Errorf("sqlite create failed: %w", err)
	}

	return nil
}

// Read retrieves an encrypted value.
func (s *SQLiteStore) Read(key string) ([]byte, error) {
	query := fmt.Sprintf("SELECT value FROM %s WHERE key = ?", sqliteTableName)
	row := s.db.QueryRow(query, key)

	var encryptedValue []byte
	err := row.Scan(&encryptedValue)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: secret with key '%s'", ErrSecretNotFound, key)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlite read failed: %w", err)
	}

	return encryptedValue, nil
}

// Update updates an existing encrypted value.
func (s *SQLiteStore) Update(key string, encryptedValue []byte) error {
	query := fmt.Sprintf("UPDATE %s SET value = ? WHERE key = ?", sqliteTableName)
	result, err := s.db.Exec(query, encryptedValue, key)
	if err != nil {
		return fmt.Errorf("sqlite update failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqlite update get rows affected failed: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: secret with key '%s'", ErrSecretNotFound, key)
	}

	return nil
}

// Delete removes a secret.
func (s *SQLiteStore) Delete(key string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE key = ?", sqliteTableName)
	result, err := s.db.Exec(query, key)
	if err != nil {
		return fmt.Errorf("sqlite delete failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqlite delete get rows affected failed: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: secret with key '%s'", ErrSecretNotFound, key)
	}

	return nil
}

// ListKeys lists all available keys.
func (s *SQLiteStore) ListKeys() ([]string, error) {
	query := fmt.Sprintf("SELECT key FROM %s", sqliteTableName)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("sqlite list keys failed: %w", err)
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("sqlite list keys scan failed: %w", err)
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlite list keys row iteration error: %w", err)
	}

	return keys, nil
}
