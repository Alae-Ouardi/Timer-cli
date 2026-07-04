package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigMissingReturnsDefault(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}
	if cfg.DefaultTheme != "obsidian" {
		t.Errorf("LoadConfig().DefaultTheme = %q, want %q", cfg.DefaultTheme, "obsidian")
	}
}

func TestSaveThenLoadRoundTrips(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg := Config{DefaultTheme: "matrix"}
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() unexpected error: %v", err)
	}

	got, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}
	if got.DefaultTheme != "matrix" {
		t.Errorf("LoadConfig().DefaultTheme = %q, want %q", got.DefaultTheme, "matrix")
	}

	path := filepath.Join(dir, "timer-cli", "config.json")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected config file at %s: %v", path, err)
	}
}

func TestConfigDirFallsBackToHomeConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", home)

	dir, err := configDir()
	if err != nil {
		t.Fatalf("configDir() unexpected error: %v", err)
	}
	want := filepath.Join(home, ".config", "timer-cli")
	if dir != want {
		t.Errorf("configDir() = %q, want %q", dir, want)
	}
}

func TestLoadConfigCorruptFileReturnsError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfgDir := filepath.Join(dir, "timer-cli")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte("not json"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := LoadConfig(); err == nil {
		t.Error("LoadConfig() expected error for corrupt file, got nil")
	}
}
