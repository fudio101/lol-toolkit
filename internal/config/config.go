package config

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
)

// Embed the config file into the binary
// Create configs/embedded.json with your API key before building
//
//go:embed embedded.json
var embeddedConfig []byte

// Config holds the application configuration
type Config struct {
	RiotAPIKey string `json:"riot_api_key"`
	Region     string `json:"region"`
}

// Default returns a default configuration
func Default() *Config {
	return &Config{
		RiotAPIKey: "",
		Region:     "vn2", // Vietnam region as default
	}
}

// LoadEmbedded loads configuration embedded in the binary
func LoadEmbedded() (*Config, error) {
	if len(embeddedConfig) == 0 {
		return Default(), nil
	}

	var cfg Config
	if err := json.Unmarshal(embeddedConfig, &cfg); err != nil {
		return Default(), nil
	}

	return &cfg, nil
}

// configPath returns the path to the user config file
func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := filepath.Join(configDir, "lol-toolkit")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appConfigDir, "config.json"), nil
}

// Load loads configuration with priority:
// 1. User config file (if exists) - allows override
// 2. Embedded config (compiled into binary)
// 3. Default config
func Load() (*Config, error) {
	// First, try to load user config (allows runtime override)
	path, err := configPath()
	if err == nil {
		data, err := os.ReadFile(path)
		if err == nil {
			var cfg Config
			if err := json.Unmarshal(data, &cfg); err == nil {
				return &cfg, nil
			}
		}
	}

	// Fall back to embedded config
	return LoadEmbedded()
}

// Save saves configuration to the user config file
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
