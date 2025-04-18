package main

import (
	"fmt"
	"log" // Keep log for general logging, return error for cobra
	"sort"

	"secrets-cli/internal/key"   // Adjust import path

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secret keys",
	Long:  `Retrieves and lists the keys of all available secrets in the store.`,
	Args:  cobra.NoArgs, // No arguments expected
	RunE: func(cmd *cobra.Command, args []string) error {
		// Encryption key is not needed for listing keys, but loading here
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

		// Use the store interface to list keys
		keys, err := s.ListKeys()
		if err != nil {
			return fmt.Errorf("failed to list secrets from store: %w", err)
		}

		if len(keys) == 0 {
			fmt.Printf("No secrets found in backend '%s'.\n", backendType)
		} else {
			fmt.Printf("Available secrets (keys) using backend '%s':\n", backendType)
			sort.Strings(keys) // Sort keys for consistent output
			for _, key := range keys {
				fmt.Printf("- %s\n", key)
			}
		}

		return nil
	},
}

func init() {
	// No flags or arguments for the list command
}
