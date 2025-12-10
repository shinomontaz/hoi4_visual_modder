package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CountryFlagsParser parses country flags from history files
type CountryFlagsParser struct {
	gamePath string
	modPath  string
}

// NewCountryFlagsParser creates a new country flags parser
func NewCountryFlagsParser(gamePath, modPath string) *CountryFlagsParser {
	return &CountryFlagsParser{
		gamePath: gamePath,
		modPath:  modPath,
	}
}

// ParseCountryFlags loads flags from history/countries/<TAG>.txt or <TAG> - <Name>.txt
func (p *CountryFlagsParser) ParseCountryFlags(countryTag string) ([]string, error) {
	// Try mod path first
	if p.modPath != "" {
		flags, err := p.parseFlagsFromPath(p.modPath, countryTag)
		if err == nil && len(flags) > 0 {
			println("Loaded", len(flags), "flags from MOD for", countryTag)
			return flags, nil
		}
	}

	// Fallback to game path
	if p.gamePath != "" {
		flags, err := p.parseFlagsFromPath(p.gamePath, countryTag)
		if err == nil {
			println("Loaded", len(flags), "flags from GAME for", countryTag)
			return flags, nil
		}
		return nil, err
	}

	return []string{}, nil // No flags found, but not an error
}

// parseFlagsFromPath tries to find and parse country history file
func (p *CountryFlagsParser) parseFlagsFromPath(basePath, countryTag string) ([]string, error) {
	historyDir := filepath.Join(basePath, "history", "countries")

	// Check if directory exists
	if _, err := os.Stat(historyDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("history/countries directory not found: %s", historyDir)
	}

	// Try to find file matching the country tag
	// Possible formats:
	// 1. TAG.txt (e.g., GER.txt)
	// 2. TAG - Name.txt (e.g., GER - Germany.txt)

	var filePath string

	// Try exact match first
	exactPath := filepath.Join(historyDir, countryTag+".txt")
	if _, err := os.Stat(exactPath); err == nil {
		filePath = exactPath
	} else {
		// Try to find file starting with TAG
		files, err := filepath.Glob(filepath.Join(historyDir, countryTag+" - *.txt"))
		if err != nil {
			return nil, err
		}

		if len(files) > 0 {
			filePath = files[0] // Take first match
		} else {
			return nil, fmt.Errorf("history file not found for country: %s", countryTag)
		}
	}

	// Parse the file
	return p.parseHistoryFile(filePath)
}

// parseHistoryFile parses a history file and extracts set_country_flag statements
func (p *CountryFlagsParser) parseHistoryFile(filePath string) ([]string, error) {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	// Parse with existing parser
	parser := NewParser(string(content))
	program, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse history file: %w", err)
	}

	// Extract flags
	flags := make([]string, 0)
	flags = p.extractFlags(program.Statements, flags)

	return flags, nil
}

// extractFlags recursively extracts set_country_flag statements
func (p *CountryFlagsParser) extractFlags(statements []Statement, flags []string) []string {
	for _, stmt := range statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}

		// Check if this is set_country_flag = FLAG_NAME
		if assignStmt.Name.Value == "set_country_flag" {
			// Extract flag name
			flagName := p.extractFlagName(assignStmt.Value)
			if flagName != "" {
				flags = append(flags, flagName)
			}
		}

		// Recursively check nested blocks
		if blockStmt, ok := assignStmt.Value.(*BlockStatement); ok {
			flags = p.extractFlags(blockStmt.Statements, flags)
		}
	}

	return flags
}

// extractFlagName extracts flag name from expression
func (p *CountryFlagsParser) extractFlagName(expr Expression) string {
	switch v := expr.(type) {
	case *Identifier:
		return v.Value
	case *StringLiteral:
		return v.Value
	default:
		return ""
	}
}

// ParseAllFlags loads flags from all available sources and merges them
// This is useful for getting a complete picture of all flags
func (p *CountryFlagsParser) ParseAllFlags(countryTag string) ([]string, error) {
	flagsMap := make(map[string]bool)

	// Load from game
	if p.gamePath != "" {
		gameFlags, err := p.parseFlagsFromPath(p.gamePath, countryTag)
		if err == nil {
			for _, flag := range gameFlags {
				flagsMap[flag] = true
			}
		}
	}

	// Load from mod (can add or override)
	if p.modPath != "" {
		modFlags, err := p.parseFlagsFromPath(p.modPath, countryTag)
		if err == nil {
			for _, flag := range modFlags {
				flagsMap[flag] = true
			}
		}
	}

	// Convert to slice
	result := make([]string, 0, len(flagsMap))
	for flag := range flagsMap {
		result = append(result, flag)
	}

	return result, nil
}

// GetUnlockFlags returns only UNLOCK:* flags
func GetUnlockFlags(flags []string) []string {
	unlockFlags := make([]string, 0)
	for _, flag := range flags {
		if strings.HasPrefix(flag, "UNLOCK:") {
			unlockFlags = append(unlockFlags, flag)
		}
	}
	return unlockFlags
}

// GetCountrySpecificFlags returns flags specific to country (e.g., GER_air, SOV_armor)
func GetCountrySpecificFlags(flags []string, countryTag string) []string {
	prefix := countryTag + "_"
	specificFlags := make([]string, 0)
	for _, flag := range flags {
		if strings.HasPrefix(flag, prefix) {
			specificFlags = append(specificFlags, flag)
		}
	}
	return specificFlags
}
