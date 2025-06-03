package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cfgPath := filepath.Join(homeDir, configFileName)

	return cfgPath, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	if err := cfg.write(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) write() error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(cfgPath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}
