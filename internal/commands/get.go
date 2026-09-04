package commands

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/Ganesh-12-spec/envoy/internal/crypto"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var GetCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Get a secret",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(".envoy/config.json")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error loading configuration:", err)
			return
		}

		fmt.Print("Enter master password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading password:", err)
			return
		}
		fmt.Println()

		salt, err := base64.StdEncoding.DecodeString(cfg.Salt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error decoding salt:", err)
			return
		}

		key, err := crypto.DeriveKey(password, salt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error deriving key:", err)
			return
		}

		vaultPath := ".envoy/vault.json"

		vault, err := config.LoadVault(vaultPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error loading vault:", err)
			return
		}

		secret, ok := vault.Secrets[args[0]]
		if !ok {
			fmt.Fprintln(os.Stderr, "Error: secret not found")
			return
		}

		ciphertext, err := base64.StdEncoding.DecodeString(secret.Ciphertext)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error decoding ciphertext:", err)
			return
		}

		nonce, err := base64.StdEncoding.DecodeString(secret.Nonce)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error decoding nonce:", err)
			return
		}

		plaintext, err := crypto.Decrypt(ciphertext, nonce, key)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error decrypting secret:", err)
			return
		}

		fmt.Println(string(plaintext))
	},
}
