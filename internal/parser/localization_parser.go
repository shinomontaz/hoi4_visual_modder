package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LocalizationParser parses HOI4 localization files
type LocalizationParser struct {
	modPath  string
	gamePath string
	language string // e.g., "english", "russian"
}

// NewLocalizationParser creates a new localization parser
func NewLocalizationParser(modPath, gamePath, language string) *LocalizationParser {
	if language == "" {
		language = "english" // Default to English
	}
	return &LocalizationParser{
		modPath:  modPath,
		gamePath: gamePath,
		language: language,
	}
}

// LoadLocalizations loads all localizations for the specified language
// Priority: mod overrides game
func (p *LocalizationParser) LoadLocalizations() (map[string]string, error) {
	localizations := make(map[string]string)

	// 1. Load from GAME first (base layer)
	if p.gamePath != "" {
		gameLocalizations, err := p.loadFromPath(p.gamePath)
		if err == nil {
			for key, value := range gameLocalizations {
				localizations[key] = value
			}
		}
	}

	// 2. Load from MOD (overrides game)
	if p.modPath != "" {
		modLocalizations, err := p.loadFromPath(p.modPath)
		if err == nil {
			for key, value := range modLocalizations {
				localizations[key] = value // Override game values
			}
		}
	}

	return localizations, nil
}

// loadFromPath loads localizations from a specific base path
func (p *LocalizationParser) loadFromPath(basePath string) (map[string]string, error) {
	locDir := filepath.Join(basePath, "localisation", p.language)

	// Check if directory exists
	if _, err := os.Stat(locDir); os.IsNotExist(err) {
		// Try without language subdirectory (some mods use flat structure)
		locDir = filepath.Join(basePath, "localisation")
		if _, err := os.Stat(locDir); os.IsNotExist(err) {
			return nil, err
		}
	}

	// Find all localization files: *_l_<language>.yml
	pattern := filepath.Join(locDir, "*_l_"+p.language+".yml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// Also try flat structure without language subdirectory
	if len(files) == 0 {
		pattern = filepath.Join(basePath, "localisation", "*_l_"+p.language+".yml")
		files, _ = filepath.Glob(pattern)
	}

	localizations := make(map[string]string)

	// Parse each file
	for _, file := range files {
		fileLocalizations, err := p.parseLocalizationFile(file)
		if err != nil {
			// Skip files with errors
			continue
		}

		// Merge localizations
		for key, value := range fileLocalizations {
			localizations[key] = value
		}
	}

	return localizations, nil
}

// parseLocalizationFile parses a single .yml localization file
func (p *LocalizationParser) parseLocalizationFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	localizations := make(map[string]string)
	scanner := bufio.NewScanner(file)

	// Skip first line (language marker like "l_english:")
	inLanguageBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for language block start (e.g., "l_english:")
		if strings.HasPrefix(line, "l_"+p.language+":") {
			inLanguageBlock = true
			continue
		}

		// Parse key-value pairs
		if inLanguageBlock {
			key, value := p.parseLocalizationLine(line)
			if key != "" {
				localizations[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return localizations, nil
}

// parseLocalizationLine parses a single localization line
// Format: key:version "value"
// Example: infantry_folder_name:0 "Infantry"
func (p *LocalizationParser) parseLocalizationLine(line string) (string, string) {
	// Find the colon that separates key from version
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return "", ""
	}

	key := strings.TrimSpace(line[:colonIndex])

	// Find the opening quote
	quoteStart := strings.Index(line, "\"")
	if quoteStart == -1 {
		return "", ""
	}

	// Find the closing quote
	quoteEnd := strings.LastIndex(line, "\"")
	if quoteEnd == -1 || quoteEnd == quoteStart {
		return "", ""
	}

	value := line[quoteStart+1 : quoteEnd]

	return key, value
}

// GetLocalization returns a localized string for a key, or the key itself if not found
func GetLocalization(localizations map[string]string, key string) string {
	if value, ok := localizations[key]; ok {
		return value
	}
	return key
}
