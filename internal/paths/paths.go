package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const rcFile = ".envoyrc"

func EnvoyDir() string {
	data, err := os.ReadFile(rcFile)
	if err == nil {
		path := strings.TrimSpace(string(data))
		if path != "" {
			return path
		}
	}

	return ".envoy"
}

func Config() string {
	return filepath.Join(EnvoyDir(), "config.json")
}

func Vault() string {
	return filepath.Join(EnvoyDir(), "vault.json")
}

func Lock() string {
	return filepath.Join(EnvoyDir(), "vault.lock")
}

func Audit() string {
	return filepath.Join(EnvoyDir(), "audit.log")
}

func Init() error {
	if err := os.WriteFile(rcFile, []byte(".envoy\n"), 0600); err != nil {
		return fmt.Errorf("creating %s: %w", rcFile, err)
	}

	return nil
}
