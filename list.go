package main

import (
	"fmt"
	"log" // Keep log for general logging, return error for cobra
	"os"
	"sort"

	"secrets-cli/internal/key"   // Adjust import path
	"secrets-cli/internal/store" // Adjust import path

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all secret keys",
	Aliases: []string{"ls"},
	Long:    `Retrieves and lists the keys of all available secrets in the store.`,
	Args:    cobra.NoArgs, // No arguments expected
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := key.LoadKeyFromEnv()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load encryption key: %v\n", err)
			os.Exit(1)
		}

		s, err := store.GetSecretStore()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get store: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			if closeErr := s.Close(); closeErr != nil {
				log.Printf("Error closing store connection: %v", closeErr)
			}
		}()

		keys, err := s.ListKeys()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to list secrets from store: %v\n", err)
			os.Exit(1)
		}

		if len(keys) == 0 {
			fmt.Printf("No secrets found in backend '%s'.\n", store.BackendType)
		} else {
			sort.Strings(keys)
			for _, key := range keys {
				fmt.Printf("%s\n", key)
			}
		}

		return nil
	},
}

func init() {
	// No flags or arguments for the list command
}
