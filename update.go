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

var (
	updateKey      string
	updateNewValue string
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing secret",
	Long:  `Updates the value of an existing secret with the given key.`,
	Args:  cobra.ExactArgs(0), // Use flags
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateKey == "" || updateNewValue == "" {
			return fmt.Errorf("both --key and --new-value flags are required")
		}

		encryptionKey, err := key.LoadKeyFromEnv()
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

		// Encrypt the new value *before* storing it
		encryptedValue, err := crypto.Encrypt([]byte(updateNewValue), encryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt new value: %w", err)
		}

		// Use the store interface to update the secret
		err = s.Update(updateKey, encryptedValue)
		if errors.Is(err, store.ErrSecretNotFound) {
			return err // Return specific error if not found
		}
		if err != nil {
			return fmt.Errorf("failed to update secret in store: %w", err)
		}

		fmt.Printf("Secret '%s' updated successfully using backend '%s'.\n", updateKey, backendType)
		return nil
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&updateKey, "key", "k", "", "The key of the secret to update")
	UpdateCmd.Flags().StringVar(&updateNewValue, "new-value", "", "The new value for the secret (will be encrypted)")
	UpdateCmd.MarkFlagRequired("key")
	UpdateCmd.MarkFlagRequired("new-value")
}
