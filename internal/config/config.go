package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// we might change the app name, we don't wanna search in 
// code we can just change in the const the gator part
const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configData, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("error reading file %s: %w", configFile, err)
	}

	var cfg Config
	err = json.Unmarshal(configData, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error decoding config file: %w", err)
	}

	return cfg, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error finding home dir: %w", err)
	}

	configPath := filepath.Join(home, configFileName) 
	// .Join instead of --> home + "/" + configFileName
	return configPath, nil
}

func write(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling object: %w", err)
	}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", configFilePath, err)
	}

	return nil
}
