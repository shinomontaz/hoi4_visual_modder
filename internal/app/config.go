package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AppConfig stores application configuration
type AppConfig struct {
	ModFilePath  string `json:"mod_file_path"` // Path to selected .mod file
	GamePath     string `json:"game_path"`     // Path to HOI4 installation
	LastCountry  string `json:"last_country"`  // Last selected country tag
	WindowWidth  int    `json:"window_width"`  // Window width
	WindowHeight int    `json:"window_height"` // Window height
}

// DefaultConfig returns default configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		ModFilePath:  "",
		GamePath:     "",
		LastCountry:  "",
		WindowWidth:  1280,
		WindowHeight: 720,
	}
}

// GetConfigPath returns the path to config file
func GetConfigPath() (string, error) {
	// Get user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	// Create app config directory
	appConfigDir := filepath.Join(configDir, "hoi4_visual_modder")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(appConfigDir, "config.json"), nil
}

// LoadConfig loads configuration from file
func LoadConfig() (*AppConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config doesn't exist, return default
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to read config: %w", err)
	}

	// Parse JSON
	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultConfig(), fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(config *AppConfig) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// UpdateModPath updates mod file path and saves config
func (c *AppConfig) UpdateModPath(modFilePath string) error {
	c.ModFilePath = modFilePath
	return SaveConfig(c)
}

// UpdateGamePath updates game path and saves config
func (c *AppConfig) UpdateGamePath(gamePath string) error {
	c.GamePath = gamePath
	return SaveConfig(c)
}

// UpdateLastCountry updates last selected country and saves config
func (c *AppConfig) UpdateLastCountry(countryTag string) error {
	c.LastCountry = countryTag
	return SaveConfig(c)
}

// UpdateWindowSize updates window size and saves config
func (c *AppConfig) UpdateWindowSize(width, height int) error {
	c.WindowWidth = width
	c.WindowHeight = height
	return SaveConfig(c)
}
