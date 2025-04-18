package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var (
	BackendType     string // Flag to select backend type
	SqliteDBPath    string // Flag for sqlite backend config
	JsonFilePath    string // Flag for jsonfile backend config
	MongoURI        string // Flag for mongodb backend config
	MongoDatabase   string // Flag for mongodb backend config
	MongoCollection string // Flag for mongodb backend config
)

// Config structure for loading defaults
type StoreConfig struct {
	BackendType     string `json:"backend_type"`
	SqliteDBPath    string `json:"sqlite_db_path"`
	JsonFilePath    string `json:"json_file_path"`
	MongoURI        string `json:"mongo_uri"`
	MongoDatabase   string `json:"mongo_database"`
	MongoCollection string `json:"mongo_collection"`
}

// LoadConfig loads config from ~/.secrets-cli.json if present
func LoadConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(home, ".secrets-cli.json")
	file, err := os.Open(configPath)
	if err != nil {
		// If config file does not exist, skip loading
		return nil
	}
	defer file.Close()

	var cfg StoreConfig
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	if BackendType == "" {
		BackendType = cfg.BackendType
	}
	if SqliteDBPath == "" {
		SqliteDBPath = cfg.SqliteDBPath
	}
	if JsonFilePath == "" {
		JsonFilePath = cfg.JsonFilePath
	}
	if MongoURI == "" {
		MongoURI = cfg.MongoURI
	}
	if MongoDatabase == "" {
		MongoDatabase = cfg.MongoDatabase
	}
	if MongoCollection == "" {
		MongoCollection = cfg.MongoCollection
	}
	return nil
}

// getSecretStore is a helper function to create and initialize the chosen backend.
func GetSecretStore() (SecretStore, error) {
	// Load defaults from config file if not already set
	//_ = LoadConfig()

	var s SecretStore
	var err error

	switch BackendType {
	case "sqlite":
		s, err = NewSQLiteStore(SqliteDBPath)
	case "jsonfile":
		s, err = NewJSONFileStore(JsonFilePath)
	case "mongodb-placeholder":
		s, err = NewMongoDBStore(MongoURI, MongoDatabase, MongoCollection)
	default:
		return nil, fmt.Errorf("unknown backend type: %s", BackendType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create store instance: %w", err)
	}

	if err := s.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize store backend: %w", err)
	}

	return s, nil
}
