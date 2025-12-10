package parser

import (
	"fmt"
	"os"
	"path/filepath"
)

// TechnologyTagsParser parses technology_tags files to extract technology folders
type TechnologyTagsParser struct {
	gamePath string
	modPath  string
}

// NewTechnologyTagsParser creates a new technology tags parser
func NewTechnologyTagsParser(gamePath, modPath string) *TechnologyTagsParser {
	return &TechnologyTagsParser{
		gamePath: gamePath,
		modPath:  modPath,
	}
}

// ParseTechnologyFolders extracts the list of technology folders
// Priority: MOD first, then GAME (mod can override or add folders)
func (p *TechnologyTagsParser) ParseTechnologyFolders() ([]string, error) {
	folders := make(map[string]bool) // Use map to avoid duplicates

	// 1. Parse from MOD path first (highest priority)
	if p.modPath != "" {
		modFolders, err := p.parseFoldersFromPath(p.modPath)
		if err == nil && len(modFolders) > 0 {
			for _, folder := range modFolders {
				folders[folder] = true
			}
		}
	}

	// 2. Parse from GAME path (fallback if mod has no folders, or to add missing ones)
	if p.gamePath != "" {
		gameFolders, err := p.parseFoldersFromPath(p.gamePath)
		if err == nil {
			for _, folder := range gameFolders {
				// Add game folders (mod folders take priority if already present)
				folders[folder] = true
			}
		}
	}

	// If no folders found at all, return error
	if len(folders) == 0 {
		return nil, fmt.Errorf("no technology folders found in mod or game")
	}

	// Convert map to slice
	result := make([]string, 0, len(folders))
	for folder := range folders {
		result = append(result, folder)
	}

	return result, nil
}

// parseFoldersFromPath parses technology_folders from a specific path
func (p *TechnologyTagsParser) parseFoldersFromPath(basePath string) ([]string, error) {
	tagsDir := filepath.Join(basePath, "common", "technology_tags")

	// Check if directory exists
	if _, err := os.Stat(tagsDir); os.IsNotExist(err) {
		return nil, err
	}

	// Find all .txt files
	files, err := filepath.Glob(filepath.Join(tagsDir, "*.txt"))
	if err != nil {
		return nil, err
	}

	folders := make([]string, 0)

	// Parse each file
	for _, file := range files {
		fileFolders, err := p.parseTagsFile(file)
		if err != nil {
			// Skip files with errors
			continue
		}
		folders = append(folders, fileFolders...)
	}

	return folders, nil
}

// parseTagsFile parses a single technology_tags file
func (p *TechnologyTagsParser) parseTagsFile(filePath string) ([]string, error) {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse using existing parser
	parser := NewParser(string(content))
	program, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	folders := make([]string, 0)

	// Find "technology_folders" block
	for _, stmt := range program.Statements {
		if assign, ok := stmt.(*AssignmentStatement); ok {
			if assign.Name.Value == "technology_folders" {
				if block, ok := assign.Value.(*BlockStatement); ok {
					folders = p.extractFolders(block)
				}
			}
		}
	}

	return folders, nil
}

// extractFolders extracts folder names from technology_folders block
func (p *TechnologyTagsParser) extractFolders(block *BlockStatement) []string {
	folders := make([]string, 0)

	for _, stmt := range block.Statements {
		// Folders are listed as identifiers (not assignments)
		// Example: infantry_folder
		if assign, ok := stmt.(*AssignmentStatement); ok {
			// Sometimes folders are written as "folder_name = yes"
			folders = append(folders, assign.Name.Value)
		} else {
			// Or just as identifiers in the block
			// This would need additional handling in the parser
			// For now, we'll rely on assignment format
		}
	}

	return folders
}
