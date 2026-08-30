package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Environment initialized")
	},
}
