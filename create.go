package main

import (
	"fmt"

	"secrets-cli/internal/crypto" // Adjust import path
	"secrets-cli/internal/key"    // Adjust import path
	// "secrets-cli/internal/store"  // Adjust import path

	"github.com/spf13/cobra"
)

var (
	createKey   string
	createValue string
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new secret",
	Long:  `Creates a new encrypted secret with the given key and value.`,
	Args:  cobra.ExactArgs(0), // Use flags
	RunE: func(cmd *cobra.Command, args []string) error {
		if createKey == "" || createValue == "" {
			return fmt.Errorf("both --key and --value flags are required")
		}

		encryptionKey, err := key.LoadKeyFromEnv()
		if err != nil {
			// Error handled by PersistentPreRunE, but good to check locally too.
			return fmt.Errorf("failed to load encryption key: %w", err)
		}

		// Get the selected store backend
		s, err := getSecretStore()
		if err != nil {
			return fmt.Errorf("failed to get store: %w", err)
		}
		defer s.Close() // Ensure store is closed

		// Encrypt the value *before* storing it
		encryptedValue, err := crypto.Encrypt([]byte(createValue), encryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}

		// Use the store interface to create the secret
		err = s.Create(createKey, encryptedValue)
		// if errors.Is(err, store.ErrSecretAlreadyExists) {
		// 	return err // Return the specific error
		// }
		if err != nil {
			return fmt.Errorf("failed to create secret in store: %w", err)
		}

		fmt.Printf("Secret '%s' created successfully using backend '%s'.\n", createKey, backendType)
		return nil
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&createKey, "key", "k", "", "The key for the secret (plaintext)")
	CreateCmd.Flags().StringVarP(&createValue, "value", "v", "", "The value of the secret (will be encrypted)")
	CreateCmd.MarkFlagRequired("key")
	CreateCmd.MarkFlagRequired("value")
}

