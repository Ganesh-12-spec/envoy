package commands

import (
	"fmt"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/Ganesh-12-spec/envoy/internal/crypto"
	"github.com/spf13/cobra"
)

var ImportCmd = &cobra.Command{
	Use:   "import BACKUP_FILE",
	Short: "Import and merge an encrypted vault backup",
	Long: `Import an encrypted vault backup and merge its secrets
into the current vault.

If a secret already exists, the imported version replaces it.`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		backupPath := args[0]
		vaultPath := ".envoy/vault.json"

		backupVault, err := config.LoadVault(backupPath)
		if err != nil {
			return fmt.Errorf("loading backup: %w", err)
		}

		currentVault, err := config.LoadVault(vaultPath)
		if err != nil {
			return fmt.Errorf("loading current vault: %w", err)
		}

		if currentVault.Secrets == nil {
			currentVault.Secrets = make(map[string]crypto.Secret)
		}

		for key, secret := range backupVault.Secrets {
			currentVault.Secrets[key] = secret
		}

		if err := config.SaveVault(currentVault, vaultPath); err != nil {
			return fmt.Errorf("saving vault: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Vault imported from %s\n", backupPath)

		return nil
	},
}
