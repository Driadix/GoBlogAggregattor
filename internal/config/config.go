package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ensure FD is released

	cfg := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("malformed json in config: %w", err)
	}

	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "gator-tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := tmpFile.Chmod(0600); err != nil {
		tmpFile.Close()
		return fmt.Errorf("setting permissions: %w", err)
	}

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("writing to temp file: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		return fmt.Errorf("syncing temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}
