package commands

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the environment",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".envoy/config.json"); err == nil {
			fmt.Fprintln(os.Stderr, "Error: envoy is already initialized")
			return
		}

		fmt.Print("Enter master password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading password:", err)
			return
		}
		fmt.Println()

		salt := make([]byte, 16)
		_, err = rand.Read(salt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error generating salt:", err)
			return
		}

		hash := sha256.Sum256(append(password, salt...))

		err = os.MkdirAll(".envoy", 0700)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating .envoy directory:", err)
			return
		}

		cfg := config.Config{
			CurrentEnvironment: "development",
			Salt:               base64.StdEncoding.EncodeToString(salt),
			PasswordHash:       base64.StdEncoding.EncodeToString(hash[:]),
		}

		err = config.Save(cfg, ".envoy/config.json")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error saving configuration:", err)
			return
		}

		fmt.Println("Environment initialized")
	},
}
