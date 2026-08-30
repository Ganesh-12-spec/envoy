package main

import (
	"os"

	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "envoy",
	Short: "A secure environment for running applications",
}
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Environment initialized")
	},
}

func main() {
	rootCmd.AddCommand(initCmd)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing command:", err.Error())
		os.Exit(1)
	}
}
