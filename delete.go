package main

import (
	"errors"
	"fmt"
	"log" // Keep log for general logging, return error for cobra

	"secrets-cli/internal/store" // Adjust import path
	"secrets-cli/internal/key"   // Adjust import path

	"github.com/spf13/cobra"
)

var deleteKey string

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a secret by its key",
	Long:  `Deletes a secret and its encrypted value based on its key.`,
	Args:  cobra.ExactArgs(0), // Use flag
	RunE: func(cmd *cobra.Command, args []string) error {
		if deleteKey == "" {
			return fmt.Errorf("--key flag is required")
		}

		// Encryption key is not needed for deletion, but loading here
		// honors PersistentPreRunE check (can be skipped if not needed by interface)
		_, err := key.LoadKeyFromEnv()
		if err != nil {
			return fmt.Errorf("failed to load encryption key: %w", err)
		}

		// Get the selected store backend
		s, err := getSecretStore()
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
			return err // Return specific error if not found
		}
		if err != nil {
			return fmt.Errorf("failed to delete secret from store: %w", err)
		}

		fmt.Printf("Secret '%s' deleted successfully using backend '%s'.\n", deleteKey, backendType)
		return nil
	},
}

func init() {
	DeleteCmd.Flags().StringVarP(&deleteKey, "key", "k", "", "The key of the secret to delete")
	DeleteCmd.MarkFlagRequired("key") // Make the key flag mandatory
}
