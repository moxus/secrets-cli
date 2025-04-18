package main

import (
	"fmt"

	"secrets-cli/internal/crypto" // Adjust import path
	"secrets-cli/internal/key"    // Adjust import path
	"secrets-cli/internal/store"  // Adjust import path

	"github.com/spf13/cobra"
)

var updateIfExists bool

var CreateCmd = &cobra.Command{
	Use:     "create [key] [value]",
	Short:   "Create a new secret",
	Aliases: []string{"add", "new", "save", "set"},
	Long:    `Creates a new encrypted secret with the given key and value.`,
	Args:    cobra.ExactArgs(2), // Require exactly two arguments
	RunE: func(cmd *cobra.Command, args []string) error {
		createKey := args[0]
		createValue := args[1]
		if createKey == "" || createValue == "" {
			return fmt.Errorf("both key and value arguments are required")
		}

		encryptionKey, err := key.LoadKeyFromEnv()
		if err != nil {
			return fmt.Errorf("failed to load encryption key: %w", err)
		}

		s, err := store.GetSecretStore()
		if err != nil {
			return fmt.Errorf("failed to get store: %w", err)
		}
		defer s.Close() // Ensure store is closed

		encryptedValue, err := crypto.Encrypt([]byte(createValue), encryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}

		if updateIfExists {
			err = s.Update(createKey, encryptedValue)
			if err != nil {
				return fmt.Errorf("failed to update secret in store: %w", err)
			}
			fmt.Printf("Secret '%s' updated successfully using backend '%s'.\n", createKey, store.BackendType)
			return nil
		}

		err = s.Create(createKey, encryptedValue)
		if err != nil {
			return fmt.Errorf("failed to create secret in store: %w", err)
		}

		fmt.Printf("Secret '%s' created successfully using backend '%s'.\n", createKey, store.BackendType)
		return nil
	},
}

func init() {
	CreateCmd.Flags().BoolVar(&updateIfExists, "update", false, "Update the secret if it already exists")
}
