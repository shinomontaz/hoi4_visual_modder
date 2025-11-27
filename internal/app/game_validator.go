package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GameInstallation represents a validated HOI4 game installation
type GameInstallation struct {
	Path       string // Absolute path to game folder
	Version    string // Game version (if detectable)
	IsValid    bool   // Validation status
	Executable string // Name of executable (hoi4.exe, hoi4.app, etc.)
}

// AutoDetectGamePath tries to find HOI4 installation automatically
func AutoDetectGamePath() (string, error) {
	var possiblePaths []string

	switch runtime.GOOS {
	case "windows":
		possiblePaths = []string{
			`C:\Program Files (x86)\Steam\steamapps\common\Hearts of Iron IV`,
			`C:\Program Files\Steam\steamapps\common\Hearts of Iron IV`,
			`D:\Steam\steamapps\common\Hearts of Iron IV`,
			`E:\Steam\steamapps\common\Hearts of Iron IV`,
			`F:\Steam\steamapps\common\Hearts of Iron IV`,
		}
	case "darwin": // macOS
		possiblePaths = []string{
			filepath.Join(os.Getenv("HOME"), "Library/Application Support/Steam/steamapps/common/Hearts of Iron IV"),
		}
	case "linux":
		possiblePaths = []string{
			filepath.Join(os.Getenv("HOME"), ".steam/steam/steamapps/common/Hearts of Iron IV"),
			filepath.Join(os.Getenv("HOME"), ".local/share/Steam/steamapps/common/Hearts of Iron IV"),
		}
	}

	// Try each path
	for _, path := range possiblePaths {
		if game, err := ValidateGameInstallation(path); err == nil && game.IsValid {
			return path, nil
		}
	}

	return "", fmt.Errorf("HOI4 installation not found in common locations")
}

// ValidateGameInstallation checks if the given path is a valid HOI4 installation
func ValidateGameInstallation(gamePath string) (*GameInstallation, error) {
	game := &GameInstallation{
		Path:    gamePath,
		IsValid: false,
	}

	// Check if folder exists
	info, err := os.Stat(gamePath)
	if os.IsNotExist(err) {
		return game, fmt.Errorf("game folder does not exist: %s", gamePath)
	}
	if !info.IsDir() {
		return game, fmt.Errorf("game path is not a directory: %s", gamePath)
	}

	// Determine executable name based on OS
	var executableName string
	switch runtime.GOOS {
	case "windows":
		executableName = "hoi4.exe"
	case "darwin":
		executableName = "hoi4.app"
	case "linux":
		executableName = "hoi4"
	default:
		executableName = "hoi4"
	}

	// Check for executable
	execPath := filepath.Join(gamePath, executableName)
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		return game, fmt.Errorf("game executable not found: %s", execPath)
	}

	game.Executable = executableName

	// Check for required folders
	requiredFolders := []string{"common", "gfx", "history"}
	for _, folder := range requiredFolders {
		folderPath := filepath.Join(gamePath, folder)
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			return game, fmt.Errorf("game folder missing required directory '%s': %s", folder, gamePath)
		}
	}

	// All checks passed
	game.IsValid = true

	// Try to detect version (optional, can fail)
	game.Version = detectGameVersion(gamePath)

	return game, nil
}

// detectGameVersion tries to detect game version from launcher-settings.json or other files
func detectGameVersion(gamePath string) string {
	// Try to read launcher-settings.json
	launcherPath := filepath.Join(gamePath, "launcher-settings.json")
	if content, err := os.ReadFile(launcherPath); err == nil {
		// Simple string search for version (not full JSON parsing)
		contentStr := string(content)
		if idx := findVersion(contentStr); idx != "" {
			return idx
		}
	}

	// If version detection fails, return unknown
	return "Unknown"
}

// findVersion is a simple helper to extract version string
func findVersion(content string) string {
	// Look for version pattern like "1.14.2"
	// This is a simplified version, could be improved
	return "" // For now, return empty
}

// GetGameInfo returns formatted game information
func (gi *GameInstallation) GetGameInfo() string {
	if gi.Version != "" && gi.Version != "Unknown" {
		return fmt.Sprintf("HOI4 v%s", gi.Version)
	}
	return "Hearts of Iron IV"
}
