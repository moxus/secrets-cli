package main

import (
	"errors"
	"fmt"
	"log" // Keep log for general logging, return error for cobra
	"os"

	"secrets-cli/internal/key"   // Adjust import path
	"secrets-cli/internal/store" // Adjust import path

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete a secret by its key",
	Long:  `Deletes a secret and its encrypted value based on its key.`,
	Args:  cobra.ExactArgs(1), // Require exactly one argument
	RunE: func(cmd *cobra.Command, args []string) error {
		deleteKey := args[0]
		if deleteKey == "" {
			return fmt.Errorf("key argument is required")
		}

		// Encryption key is not needed for deletion, but loading here
		// honors PersistentPreRunE check (can be skipped if not needed by interface)
		_, err := key.LoadKeyFromEnv()
		if err != nil {
			return fmt.Errorf("failed to load encryption key: %w", err)
		}

		// Get the selected store backend
		s, err := store.GetSecretStore()
		if err != nil {
			return fmt.Errorf("failed to get store: %w", err)
		}
		defer func() {
			if closeErr := s.Close(); closeErr != nil {
				log.Printf("Error closing store connection: %v", closeErr)
			}
		}() // Ensure store is closed

		// Use the store interface to delete the secret
		err = s.Delete(deleteKey)
		if errors.Is(err, store.ErrSecretNotFound) {
			fmt.Fprintf(os.Stderr, "secret with key '%s' not found\n", deleteKey)
			os.Exit(1)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete secret from store: %v\n", err)
			os.Exit(1)
		}

		//fmt.Printf("Secret '%s' deleted successfully using backend '%s'.\n", deleteKey, store.BackendType)
		return nil
	},
}

func init() {
	// No flag needed for key anymore
}
