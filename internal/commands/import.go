package commands

import (
	"fmt"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/Ganesh-12-spec/envoy/internal/crypto"
	"github.com/spf13/cobra"
)

var ImportCmd = &cobra.Command{
	Use:   "import BACKUP_FILE",
	Short: "Import an encrypted vault backup",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		backupPath := args[0]
		vaultPath := ".envoy/vault.json"

		// Load backup vault
		backupVault, err := config.LoadVault(backupPath)
		if err != nil {
			return err
		}

		// Load current vault
		currentVault, err := config.LoadVault(vaultPath)
		if err != nil {
			return err
		}

		if currentVault.Secrets == nil {
			currentVault.Secrets = make(map[string]crypto.Secret)
		}

		// Merge backup secrets into current vault
		for key, secret := range backupVault.Secrets {
			currentVault.Secrets[key] = secret
		}

		// Save merged vault
		if err := config.SaveVault(currentVault, vaultPath); err != nil {
			return err
		}

		fmt.Printf("Vault imported from %s\n", backupPath)

		return nil
	},
}
