package commands

import (
	"fmt"
	"os"

	"github.com/Ganesh-12-spec/envoy/internal/lock"
	"github.com/Ganesh-12-spec/envoy/internal/paths"
	"github.com/spf13/cobra"
)

var ExportCmd = &cobra.Command{
	Use:   "export BACKUP_FILE",
	Short: "Export an encrypted vault backup",
	Long: `Export the encrypted vault to a backup file.

The exported file contains encrypted secrets and does not contain
plaintext secret values.`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		backupPath := args[0]
		vaultPath := paths.Vault()

		fileLock, err := lock.Acquire(paths.Lock())
		if err != nil {
			return fmt.Errorf("acquiring vault lock: %w", err)
		}

		data, err := os.ReadFile(vaultPath)
		if err != nil {
			fileLock.Release()
			return fmt.Errorf("reading vault: %w", err)
		}

		if err := fileLock.Release(); err != nil {
			return fmt.Errorf("releasing vault lock: %w", err)
		}

		if err := os.WriteFile(backupPath, data, 0600); err != nil {
			return fmt.Errorf("writing backup: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Vault exported to %s\n", backupPath)

		return nil
	},
}
