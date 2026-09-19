package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/Ganesh-12-spec/envoy/internal/lock"
	"github.com/Ganesh-12-spec/envoy/internal/paths"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete KEY",
	Short: "Delete a secret from the vault",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		vaultPath := paths.Vault()

		fileLock, err := lock.Acquire(paths.Lock())
		if err != nil {
			return fmt.Errorf("acquiring vault lock: %w", err)
		}
		defer fileLock.Release()

		vault, err := config.LoadVault(vaultPath)
		if err != nil {
			return fmt.Errorf("loading vault: %w", err)
		}

		if _, exists := vault.Secrets[key]; !exists {
			return fmt.Errorf("secret %q not found", key)
		}

		fmt.Printf("Are you sure you want to delete %q? (y/n): ", key)

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading confirmation: %w", err)
		}

		input = strings.ToLower(strings.TrimSpace(input))

		if input != "y" && input != "yes" {
			fmt.Println("Deletion cancelled.")
			return nil
		}

		delete(vault.Secrets, key)

		if err := config.SaveVault(vault, vaultPath); err != nil {
			return fmt.Errorf("saving vault: %w", err)
		}

		fmt.Printf("Secret %q deleted successfully\n", key)

		return nil
	},
}
