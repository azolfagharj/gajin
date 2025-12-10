package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// EnvTokenKey is the environment variable key for GitHub token.
	EnvTokenKey = "GH_TOKEN_WITH_ACTIONS_WRITE"
)

// LoadConfig loads configuration from a YAML file.
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Load token from environment variable if not set in config
	if cfg.GitHub.Token == "" {
		cfg.GitHub.Token = os.Getenv(EnvTokenKey)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// LoadConfigFromPath loads configuration from a path, expanding it if needed.
func LoadConfigFromPath(path string) (*Config, error) {
	expandedPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand config path: %w", err)
	}

	return LoadConfig(expandedPath)
}

