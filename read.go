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

var ReadCmd = &cobra.Command{
	Use:     "read [key]",
	Short:   "Read a secret by its key",
	Aliases: []string{"get"},
	Long:    `Retrieves and decrypts a secret value based on its key.`,
	Args:    cobra.ExactArgs(1), // Require exactly one argument
	RunE: func(cmd *cobra.Command, args []string) error {
		readKey := args[0]
		if readKey == "" {
			return fmt.Errorf("key argument is required")
		}

		encryptionKey, err := key.LoadKeyFromEnv()
		if err != nil {
			return fmt.Errorf("failed to load encryption key: %w", err)
		}

		s, err := store.GetSecretStore()
		if err != nil {
			return fmt.Errorf("failed to get store: %w", err)
		}
		defer func() {
			if closeErr := s.Close(); closeErr != nil {
				log.Printf("Error closing store connection: %v", closeErr)
			}
		}()

		encryptedValue, err := s.Read(readKey)
		if errors.Is(err, store.ErrSecretNotFound) {
			return err
		}
		if err != nil {
			return fmt.Errorf("failed to read secret from store: %w", err)
		}

		secretValue, err := crypto.Decrypt(encryptedValue, encryptionKey)
		if err != nil {
			return fmt.Errorf("failed to decrypt value for key '%s': %w", readKey, err)
		}

		fmt.Printf("Secret '%s': %s\n", readKey, string(secretValue))
		return nil
	},
}

func init() {
	// No flag needed for key anymore
}
