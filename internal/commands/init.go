package commands

import (
	"fmt"
	"os"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the environment",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.MkdirAll(".envoy", 0700)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating .envoy directory:", err)
			return
		}
		cfg := config.Config{
			CurrentEnvironment: "development",
		}

		err = config.Save(cfg, ".envoy/config.json")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error saving configuration:", err)
			return
		}

		fmt.Println("Environment initialized")
	},
}
