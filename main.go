package main

import (
	"fmt"
	"log"
	"os"

	"secrets-cli/internal/key"   // Adjust import path
	"secrets-cli/internal/store" // Adjust import path

	// Adjust import path
	"github.com/spf13/cobra"
)

func main() {
	// Load config file to initialize backend parameters before flags are parsed
	err := store.LoadConfig()
	if err != nil {
		fmt.Println("error loading config", err)
		panic(err)
	}

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
	rootCmd.PersistentFlags().StringVar(&store.BackendType, "backend", store.BackendType, "Storage backend type (sqlite, jsonfile, mongodb-placeholder)")
	rootCmd.PersistentFlags().StringVar(&store.SqliteDBPath, "sqlite-db", store.SqliteDBPath, "SQLite database file path")
	rootCmd.PersistentFlags().StringVar(&store.JsonFilePath, "json-file", store.JsonFilePath, "JSON file path")
	rootCmd.PersistentFlags().StringVar(&store.MongoURI, "mongo-uri", store.MongoURI, "MongoDB connection URI")
	rootCmd.PersistentFlags().StringVar(&store.MongoDatabase, "mongo-db", store.MongoDatabase, "MongoDB database name")
	rootCmd.PersistentFlags().StringVar(&store.MongoCollection, "mongo-collection", store.MongoCollection, "MongoDB collection name")

	// Add subcommands
	//rootCmd.AddCommand(generateKeyCmd)
	rootCmd.AddCommand(CreateCmd)
	rootCmd.AddCommand(ReadCmd)
	rootCmd.AddCommand(DeleteCmd)
	rootCmd.AddCommand(ListCmd)
	rootCmd.AddCommand(GenerateCmd)

	if err := rootCmd.Execute(); err != nil {
		// Error handling is now mostly within RunE functions,
		// so simply exiting after cobra prints the error is fine.
		os.Exit(1)
	}
}
