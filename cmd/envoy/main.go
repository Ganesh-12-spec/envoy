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

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing command:", err.Error())
		os.Exit(1)
	}
}
