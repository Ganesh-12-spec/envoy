package main

import (
	"fmt"
	"os"

	"github.com/Ganesh-12-spec/envoy/internal/commands"
)

func main() {
	commands.RootCmd.AddCommand(commands.InitCmd)

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
