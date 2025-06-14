package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DefaultPath string `yaml:"defaultPath"`
}

// Load loads the configuration from the YAML config file
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "rmtk", "conf.yml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found")
	}

	// Read the config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse YAML
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// GetDefaultPath returns the default path from the config file or empty string if not valid
func GetDefaultPath() string {
	config, err := Load()
	if err != nil {
		return ""
	}

	path := strings.TrimSpace(config.DefaultPath)
	if path == "" {
		return ""
	}

	// Get home directory for tilde expansion
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Expand ~ to home directory if present
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(homeDir, path[2:])
	}

	// Check if the path exists and is a directory
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return path
	}

	return ""
}
