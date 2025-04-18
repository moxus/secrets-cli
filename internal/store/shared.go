package store

import "fmt"

var (
	BackendType     string // Flag to select backend type
	SqliteDBPath    string // Flag for sqlite backend config
	JsonFilePath    string // Flag for jsonfile backend config
	MongoURI        string // Flag for mongodb backend config
	MongoDatabase   string // Flag for mongodb backend config
	MongoCollection string // Flag for mongodb backend config
)

// getSecretStore is a helper function to create and initialize the chosen backend.
func GetSecretStore() (SecretStore, error) {
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
