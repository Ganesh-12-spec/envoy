package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ExportCmd = &cobra.Command{
	Use:   "export <backup-file>",
	Short: "Export an encrypted backup of the vault",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		backupPath := args[0]

		vaultPath := ".envoy/vault.json"

		data, err := os.ReadFile(vaultPath)
		if err != nil {
			return fmt.Errorf("failed to read vault: %w", err)
		}

		err = os.WriteFile(backupPath, data, 0600)
		if err != nil {
			return fmt.Errorf("failed to write backup: %w", err)
		}

		fmt.Printf("vault exported to %s\n", backupPath)

		return nil
	},
}
