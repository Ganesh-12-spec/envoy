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

var SetCmd = &cobra.Command{
	Use:   "set KEY VALUE",
	Short: "Set a secret",
	Args:  cobra.ExactArgs(2),
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

		ciphertext, nonce, err := crypto.Encrypt([]byte(args[1]), key)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error encrypting secret:", err)
			return
		}

		vaultPath := ".envoy/vault.json"

		vault, err := config.LoadVault(vaultPath)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintln(os.Stderr, "Error loading vault:", err)
				return
			}

			vault = crypto.Vault{
				Secrets: make(map[string]crypto.Secret),
			}
		}

		if vault.Secrets == nil {
			vault.Secrets = make(map[string]crypto.Secret)
		}

		vault.Secrets[args[0]] = crypto.Secret{
			Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
			Nonce:      base64.StdEncoding.EncodeToString(nonce),
		}

		err = config.SaveVault(vault, vaultPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error saving vault:", err)
			return
		}

		fmt.Println("Secret saved")
	},
}
