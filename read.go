package main

import (
	"errors"
	"fmt"
	"log" // Keep log for general logging, return error for cobra

	"secrets-cli/internal/crypto" // Adjust import path
	"secrets-cli/internal/key"    // Adjust import path
	"secrets-cli/internal/store"  // Adjust import path

	"github.com/spf13/cobra"
)

var readKey string

var ReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read a secret by its key",
	Long:  `Retrieves and decrypts a secret value based on its key.`,
	Args:  cobra.ExactArgs(0), // Use flag
	RunE: func(cmd *cobra.Command, args []string) error {
		if readKey == "" {
			return fmt.Errorf("--key flag is required")
		}

		encryptionKey, err := key.LoadKeyFromEnv()
		if err != nil {
			// Error handled by PersistentPreRunE, but returning here ensures clean exit
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

		// Use the store interface to read the encrypted value
		encryptedValue, err := s.Read(readKey)
		if errors.Is(err, store.ErrSecretNotFound) {
			return err // Return the specific error if the secret isn't found
		}
		if err != nil {
			return fmt.Errorf("failed to read secret from store: %w", err)
		}

		// Decrypt the retrieved value
		secretValue, err := crypto.Decrypt(encryptedValue, encryptionKey)
		if err != nil {
			return fmt.Errorf("encryptedValue '%w'", encryptedValue)
			return fmt.Errorf("failed to decrypt value for key 999 '%s': %w", readKey, err)
		}

		fmt.Printf("Secret '%s': %s\n", readKey, string(secretValue))
		return nil
	},
}

func init() {
	ReadCmd.Flags().StringVarP(&readKey, "key", "k", "", "The key of the secret to read")
	ReadCmd.MarkFlagRequired("key") // Make the key flag mandatory
}
