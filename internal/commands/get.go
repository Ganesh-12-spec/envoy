package commands

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/Ganesh-12-spec/envoy/internal/crypto"
	"github.com/Ganesh-12-spec/envoy/internal/lock"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var GetCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Retrieve a secret",
	Long: `Retrieve and decrypt a secret from the Envoy vault.

KEY may be a simple name or a namespaced name such as:
  envoy get database/password
  envoy get api/key`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(".envoy/config.json")
		if err != nil {
			return fmt.Errorf("loading configuration: %w", err)
		}

		fmt.Fprint(cmd.OutOrStdout(), "Enter master password: ")

		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("reading password: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout())

		salt, err := base64.StdEncoding.DecodeString(cfg.Salt)
		if err != nil {
			return fmt.Errorf("decoding salt: %w", err)
		}

		hash := sha256.Sum256(append(password, salt...))

		storedHash, err := base64.StdEncoding.DecodeString(cfg.PasswordHash)
		if err != nil {
			return fmt.Errorf("decoding password hash: %w", err)
		}

		if subtle.ConstantTimeCompare(hash[:], storedHash) != 1 {
			return fmt.Errorf("incorrect master password")
		}

		key, err := crypto.DeriveKey(password, salt)
		if err != nil {
			return fmt.Errorf("deriving encryption key: %w", err)
		}

		fileLock, err := lock.Acquire(".envoy/vault.lock")
		if err != nil {
			return fmt.Errorf("acquiring vault lock: %w", err)
		}
		defer fileLock.Release()

		vault, err := config.LoadVault(".envoy/vault.json")
		if err != nil {
			return fmt.Errorf("loading vault: %w", err)
		}

		secret, ok := vault.Secrets[args[0]]
		if !ok {
			return fmt.Errorf("secret %q not found", args[0])
		}

		ciphertext, err := base64.StdEncoding.DecodeString(secret.Ciphertext)
		if err != nil {
			return fmt.Errorf("decoding ciphertext: %w", err)
		}

		nonce, err := base64.StdEncoding.DecodeString(secret.Nonce)
		if err != nil {
			return fmt.Errorf("decoding nonce: %w", err)
		}

		plaintext, err := crypto.Decrypt(ciphertext, nonce, key)
		if err != nil {
			return fmt.Errorf("incorrect master password")
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(plaintext))

		return nil
	},
}
