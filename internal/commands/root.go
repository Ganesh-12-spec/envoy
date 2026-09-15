package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "envoy",
	Short: "A secure CLI for managing encrypted secrets",
	Long: `Envoy is a secure command-line secret manager.

Store, retrieve, list, delete, export, and import encrypted secrets
using a master password.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	RootCmd.AddCommand(SetCmd)
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(DeleteCmd)
	RootCmd.AddCommand(ExportCmd)
	RootCmd.AddCommand(ImportCmd)
}

func printError(err error) {
	fmt.Fprintln(RootCmd.ErrOrStderr(), "Error:", err)
}
