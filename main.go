package main

import (
	"fmt"
	"log"
	"os"

	"secrets-cli/internal/key"   // Adjust import path
	"secrets-cli/internal/store" // Adjust import path

	"github.com/spf13/cobra"
)

var (
	backendType     string // Flag to select backend type
	sqliteDBPath    string // Flag for sqlite backend config
	jsonFilePath    string // Flag for jsonfile backend config
	mongoURI        string // Flag for mongodb backend config
	mongoDatabase   string // Flag for mongodb backend config
	mongoCollection string // Flag for mongodb backend config
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "secrets-cli",
		Short: "Secure Secrets Storage CLI with multiple backends",
		Long: `A command-line tool to manage encrypted key-value secrets
using different storage backends (sqlite, jsonfile, mongodb-placeholder).
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Check if encryption key is available before most commands
			// (Skip for generate-key command)
			if cmd.Name() != "generate-key" {
				_, err := key.LoadKeyFromEnv()
				if err != nil {
					// Log the error but let the command's RunE handle the exit
					log.Printf("Encryption key error: %v", err)
					return err // Cobra will print the error and exit
				}
			}
			return nil
		},
	}

	// Add persistent flags for backend selection and configuration
	rootCmd.PersistentFlags().StringVar(&backendType, "backend", "sqlite", "Storage backend type (sqlite, jsonfile, mongodb-placeholder)")
	rootCmd.PersistentFlags().StringVar(&sqliteDBPath, "sqlite-db", "secrets.db", "SQLite database file path")
	rootCmd.PersistentFlags().StringVar(&jsonFilePath, "json-file", "secrets.json", "JSON file path")
	rootCmd.PersistentFlags().StringVar(&mongoURI, "mongo-uri", "mongodb://localhost:27017", "MongoDB connection URI")
	rootCmd.PersistentFlags().StringVar(&mongoDatabase, "mongo-db", "secrets", "MongoDB database name")
	rootCmd.PersistentFlags().StringVar(&mongoCollection, "mongo-collection", "secrets", "MongoDB collection name")

	// Add subcommands
	//rootCmd.AddCommand(generateKeyCmd)
	rootCmd.AddCommand(CreateCmd)
	rootCmd.AddCommand(ReadCmd)
	rootCmd.AddCommand(UpdateCmd)
	rootCmd.AddCommand(DeleteCmd)
	rootCmd.AddCommand(ListCmd)

	if err := rootCmd.Execute(); err != nil {
		// Error handling is now mostly within RunE functions,
		// so simply exiting after cobra prints the error is fine.
		os.Exit(1)
	}
}

// getSecretStore is a helper function to create and initialize the chosen backend.
func getSecretStore() (store.SecretStore, error) {
	var s store.SecretStore
	var err error

	switch backendType {
	case "sqlite":
		s, err = store.NewSQLiteStore(sqliteDBPath)
	case "jsonfile":
		s, err = store.NewJSONFileStore(jsonFilePath)
	case "mongodb-placeholder":
		s, err = store.NewMongoDBStore(mongoURI, mongoDatabase, mongoCollection)
	default:
		return nil, fmt.Errorf("unknown backend type: %s", backendType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create store instance: %w", err)
	}

	if err := s.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize store backend: %w", err)
	}

	return s, nil
}
