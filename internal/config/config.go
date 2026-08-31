package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	CurrentEnvironment string `json:"current_environment"`
}

func Save(cfg Config, path string)  error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
func Load(path string) (Config, error) {
	var cfg Config
	data,err := os.ReadFile(path)
	if err != nil {
		
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		
		return cfg, err
	}
	return cfg, nil
}