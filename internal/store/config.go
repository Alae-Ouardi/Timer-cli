package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"timer/internal/theme"
)

// Config holds user preferences persisted across runs.
type Config struct {
	DefaultTheme string `json:"default_theme"`
}

func defaultConfig() Config {
	return Config{DefaultTheme: theme.Default().Name}
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// LoadConfig loads the config file, returning defaults (with no error) if
// the file does not exist yet.
func LoadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, fmt.Errorf("resolve config path: %w", err)
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return defaultConfig(), nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// SaveConfig writes the config file, creating the config directory if
// necessary.
func SaveConfig(cfg Config) error {
	dir, err := ensureConfigDir()
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "config.json"), data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
