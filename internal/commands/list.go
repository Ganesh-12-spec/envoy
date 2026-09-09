package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets",

	Run: func(cmd *cobra.Command, args []string) {
		vaultPath := ".envoy/vault.json"

		vault, err := config.LoadVault(vaultPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error loading vault:", err)
			return
		}

		groups := make(map[string][]string)

		for name := range vault.Secrets {
			parts := strings.SplitN(name, "/", 2)

			if len(parts) == 2 {
				namespace := parts[0]
				secretName := parts[1]

				groups[namespace] = append(groups[namespace], secretName)
			} else {
				groups["default"] = append(groups["default"], name)
			}
		}

		namespaces := make([]string, 0, len(groups))

		for namespace := range groups {
			namespaces = append(namespaces, namespace)
		}

		sort.Strings(namespaces)

		fmt.Println("Secrets:")

		for _, namespace := range namespaces {
			fmt.Printf("\n%s/\n", namespace)

			sort.Strings(groups[namespace])

			for _, secretName := range groups[namespace] {
				fmt.Printf("  %s\n", secretName)
			}
		}
	},
}
