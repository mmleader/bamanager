package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func NewConfig() *Config {
	return &Config{
		Width:  1024,
		Height: 768,
	}
}

func (c *Config) Load() error {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".bsmanager", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

func (c *Config) Save() error {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".bsmanager")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}
	configPath := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
