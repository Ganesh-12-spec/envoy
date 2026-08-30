package main

import (
	"fmt"
	"os"

	"github.com/Ganesh-12-spec/envoy/internal/commands"
)

func main() {
	commands.RootCmd.AddCommand(commands.InitCmd)

	err := commands.RootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing command:", err.Error())
		os.Exit(1)
	}
}
