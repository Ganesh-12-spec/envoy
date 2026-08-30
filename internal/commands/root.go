package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = cobra.Command{
	Use:   "envoy",
	Short: "A secure environment for running applications",
}
