package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete KEY",
	Short: "Delete a secret from the vault",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		vaultPath := ".envoy/vault.json"

		// 1. Load the vault
		vault, err := config.LoadVault(vaultPath)
		if err != nil {
			return err
		}

		// 2. Check whether the secret exists
		if _, exists := vault.Secrets[key]; !exists {
			return fmt.Errorf("secret %q not found", key)
		}

		// 3. Ask for confirmation
		fmt.Printf("Are you sure you want to delete %q? (y/n): ", key)

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		input = strings.ToLower(strings.TrimSpace(input))

		// 4. If not yes → do nothing
		if input != "y" && input != "yes" {
			fmt.Println("Deletion cancelled.")
			return nil
		}

		// 5. Delete the secret
		delete(vault.Secrets, key)

		// 6. Save the updated vault
		if err := config.SaveVault(vault, vaultPath); err != nil {
			return err
		}

		// 7. Print success
		fmt.Printf("Secret %q deleted successfully\n", key)

		return nil
	},
}
