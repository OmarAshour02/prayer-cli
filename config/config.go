package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	City   string `json:"city"`
}

func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "prayers", "config.json"), nil
}

func LoadOrInitConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultCfg := &Config{
			City:   "Makkah",
		}

		if err := SaveConfig(defaultCfg); err != nil {
			return nil, err
		}

		fmt.Println("Config initialized at:", path)
		return defaultCfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func UpdateConfig(updater func(cfg *Config)) error {
	cfg, err := LoadOrInitConfig()
	if err != nil {
		return err
	}

	updater(cfg)

	return SaveConfig(cfg)
}
