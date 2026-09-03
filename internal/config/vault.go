package config

import (
	"encoding/json"
	"os"

	"github.com/Ganesh-12-spec/envoy/internal/crypto"
)

func SaveVault(vault crypto.Vault, path string) error {
	data, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0600)
	if err != nil {
		return err
	}

	return nil
}

func LoadVault(path string) (crypto.Vault, error) {
	var vault crypto.Vault

	data, err := os.ReadFile(path)
	if err != nil {
		return vault, err
	}

	err = json.Unmarshal(data, &vault)
	if err != nil {
		return vault, err
	}

	return vault, nil
}
