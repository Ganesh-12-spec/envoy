package commands

import (
	"fmt"
	"strings"

	"github.com/Ganesh-12-spec/envoy/internal/audit"
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

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		action := strings.ToUpper(cmd.Name())

		var target string

		switch cmd.Name() {
		case "set", "get", "delete":
			if len(args) > 0 {
				target = args[0]
			}
		}

		if err := audit.Log(action, target); err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), "Warning: could not write audit log:", err)
		}
	},
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
