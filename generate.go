package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"secrets-cli/internal/crypto"
	"secrets-cli/internal/key"
	"secrets-cli/internal/store"

	"github.com/spf13/cobra"
)

var (
	genUppercase      bool
	genLowercase      bool
	genNumbers        bool
	genUpdateIfExists bool
)

var GenerateCmd = &cobra.Command{
	Use:     "generate [key] [length]",
	Aliases: []string{"gen"},
	Short:   "Generate and store a random password",
	Long: `Generates a random password of the specified length using the selected character sets
and stores it as a secret under the given key.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		createKey := args[0]
		length, err := parseLengthArg(args[1])
		if err != nil {
			return err
		}
		if createKey == "" {
			return fmt.Errorf("key argument is required")
		}
		if length <= 0 {
			return fmt.Errorf("password length must be positive")
		}

		charset := ""
		if genUppercase {
			charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		}
		if genLowercase {
			charset += "abcdefghijklmnopqrstuvwxyz"
		}
		if genNumbers {
			charset += "0123456789"
		}
		if charset == "" {
			return fmt.Errorf("at least one of --uppercase, --lowercase, or --numbers must be set")
		}

		password, err := generatePassword(length, charset)
		if err != nil {
			return fmt.Errorf("failed to generate password: %w", err)
		}

		encryptionKey, err := key.LoadKeyFromEnv()
		if err != nil {
			return fmt.Errorf("failed to load encryption key: %w", err)
		}

		s, err := store.GetSecretStore()
		if err != nil {
			return fmt.Errorf("failed to get store: %w", err)
		}
		defer s.Close()

		encryptedValue, err := crypto.Encrypt([]byte(password), encryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt password: %w", err)
		}

		if genUpdateIfExists {
			err = s.Update(createKey, encryptedValue)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to update secret in store: %v\n", err)
				os.Exit(1)
			}
			return nil
		}

		err = s.Create(createKey, encryptedValue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create secret in store: %v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

func parseLengthArg(arg string) (int, error) {
	var length int
	_, err := fmt.Sscanf(arg, "%d", &length)
	if err != nil {
		return 0, fmt.Errorf("invalid length argument: %v", err)
	}
	return length, nil
}

func generatePassword(length int, charset string) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = charset[int(b)%len(charset)]
	}
	return string(bytes), nil
}

func init() {
	GenerateCmd.Flags().BoolVarP(&genUppercase, "uppercase", "u", false, "Include uppercase letters")
	GenerateCmd.Flags().BoolVarP(&genLowercase, "lowercase", "l", false, "Include lowercase letters")
	GenerateCmd.Flags().BoolVarP(&genNumbers, "numbers", "n", false, "Include numbers")
	GenerateCmd.Flags().BoolVar(&genUpdateIfExists, "update", false, "Update the secret if it already exists")
}
