package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	CurrentEnvironment string `json:"current_environment"`
}

func Save(cfg Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
