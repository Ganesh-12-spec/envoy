package commands

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Ganesh-12-spec/envoy/internal/config"
	"github.com/Ganesh-12-spec/envoy/internal/paths"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Envoy vault",
	Long: `Initialize Envoy in the current directory.

This creates the .envoy directory and stores the vault configuration,
including a random salt and password hash.`,
	Args: cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(paths.Config()); err == nil {
			return fmt.Errorf("envoy is already initialized")
		}

		fmt.Fprint(cmd.OutOrStdout(), "Enter master password: ")

		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("reading password: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout())

		salt := make([]byte, 16)

		_, err = rand.Read(salt)
		if err != nil {
			return fmt.Errorf("generating salt: %w", err)
		}

		hash := sha256.Sum256(append(password, salt...))

		if err := os.MkdirAll(paths.EnvoyDir(), 0700); err != nil {
			return fmt.Errorf("creating Envoy directory: %w", err)
		}

		cfg := config.Config{
			CurrentEnvironment: "development",
			Salt:               base64.StdEncoding.EncodeToString(salt),
			PasswordHash:       base64.StdEncoding.EncodeToString(hash[:]),
		}

		if err := config.Save(cfg, paths.Config()); err != nil {
			return fmt.Errorf("saving configuration: %w", err)
		}

		if err := paths.Init(); err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Environment initialized")

		return nil
	},
}
