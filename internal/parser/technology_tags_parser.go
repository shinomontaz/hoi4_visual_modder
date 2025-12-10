package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TechFolder represents a technology folder with its metadata
type TechFolder struct {
	Name      string
	Ledger    string // army, navy, air, civilian
	Available *AvailableCondition
	IsOverlay bool
}

// AvailableCondition represents availability conditions for a folder
type AvailableCondition struct {
	Conditions []*Condition
}

// Condition represents a single condition (has_country_flag, NOT, etc.)
type Condition struct {
	Type     string // "has_country_flag", "NOT", "has_dlc", "major_country"
	Value    string
	Negated  bool
	Children []*Condition // for nested conditions like NOT { ... }
}

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

// ParseTechnologyFoldersDetailed extracts detailed folder information including conditions
func (p *TechnologyTagsParser) ParseTechnologyFoldersDetailed() ([]*TechFolder, error) {
	folders := make(map[string]*TechFolder) // Use map to merge mod+game

	// 1. Parse from GAME path first (base)
	if p.gamePath != "" {
		gameFolders, err := p.parseFoldersDetailedFromPath(p.gamePath)
		if err == nil {
			for _, folder := range gameFolders {
				folders[folder.Name] = folder
			}
		}
	}

	// 2. Parse from MOD path (override/add)
	if p.modPath != "" {
		modFolders, err := p.parseFoldersDetailedFromPath(p.modPath)
		if err == nil {
			for _, folder := range modFolders {
				folders[folder.Name] = folder // Mod overrides game
			}
		}
	}

	if len(folders) == 0 {
		return nil, fmt.Errorf("no technology folders found")
	}

	// Convert to slice
	result := make([]*TechFolder, 0, len(folders))
	for _, folder := range folders {
		result = append(result, folder)
	}

	return result, nil
}

// parseFoldersDetailedFromPath parses detailed folder info from path
func (p *TechnologyTagsParser) parseFoldersDetailedFromPath(basePath string) ([]*TechFolder, error) {
	tagsDir := filepath.Join(basePath, "common", "technology_tags")

	if _, err := os.Stat(tagsDir); os.IsNotExist(err) {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(tagsDir, "*.txt"))
	if err != nil {
		return nil, err
	}

	folders := make([]*TechFolder, 0)

	for _, file := range files {
		fileFolders, err := p.parseTagsFileDetailed(file)
		if err != nil {
			continue
		}
		folders = append(folders, fileFolders...)
	}

	return folders, nil
}

// parseTagsFileDetailed parses a file with detailed folder information
func (p *TechnologyTagsParser) parseTagsFileDetailed(filePath string) ([]*TechFolder, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	parser := NewParser(string(content))
	program, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	folders := make([]*TechFolder, 0)

	for _, stmt := range program.Statements {
		if assign, ok := stmt.(*AssignmentStatement); ok {
			if assign.Name.Value == "technology_folders" {
				if block, ok := assign.Value.(*BlockStatement); ok {
					folders = p.extractFoldersDetailed(block)
				}
			}
		}
	}

	return folders, nil
}

// extractFoldersDetailed extracts detailed folder information
func (p *TechnologyTagsParser) extractFoldersDetailed(block *BlockStatement) []*TechFolder {
	folders := make([]*TechFolder, 0)

	for _, stmt := range block.Statements {
		if assign, ok := stmt.(*AssignmentStatement); ok {
			folderName := assign.Name.Value
			folder := &TechFolder{
				Name:      folderName,
				IsOverlay: strings.HasSuffix(folderName, "_overlay_folder"),
			}

			// Parse folder block
			if folderBlock, ok := assign.Value.(*BlockStatement); ok {
				p.parseFolderBlock(folder, folderBlock)
			}

			folders = append(folders, folder)
		}
	}

	return folders
}

// parseFolderBlock parses the contents of a folder block
func (p *TechnologyTagsParser) parseFolderBlock(folder *TechFolder, block *BlockStatement) {
	for _, stmt := range block.Statements {
		if assign, ok := stmt.(*AssignmentStatement); ok {
			switch assign.Name.Value {
			case "ledger":
				if ident, ok := assign.Value.(*Identifier); ok {
					folder.Ledger = ident.Value
				}
			case "available":
				if availBlock, ok := assign.Value.(*BlockStatement); ok {
					folder.Available = p.parseAvailableCondition(availBlock)
				}
			}
		}
	}
}

// parseAvailableCondition parses available = { ... } block
func (p *TechnologyTagsParser) parseAvailableCondition(block *BlockStatement) *AvailableCondition {
	condition := &AvailableCondition{
		Conditions: make([]*Condition, 0),
	}

	for _, stmt := range block.Statements {
		if assign, ok := stmt.(*AssignmentStatement); ok {
			cond := p.parseCondition(assign)
			if cond != nil {
				condition.Conditions = append(condition.Conditions, cond)
			}
		}
	}

	return condition
}

// parseCondition parses a single condition
func (p *TechnologyTagsParser) parseCondition(assign *AssignmentStatement) *Condition {
	condName := assign.Name.Value

	switch condName {
	case "has_country_flag":
		// has_country_flag = FLAG_NAME
		flagName := p.extractValue(assign.Value)
		return &Condition{
			Type:  "has_country_flag",
			Value: flagName,
		}

	case "NOT":
		// NOT = { ... }
		if block, ok := assign.Value.(*BlockStatement); ok {
			children := make([]*Condition, 0)
			for _, stmt := range block.Statements {
				if childAssign, ok := stmt.(*AssignmentStatement); ok {
					childCond := p.parseCondition(childAssign)
					if childCond != nil {
						childCond.Negated = true
						children = append(children, childCond)
					}
				}
			}
			return &Condition{
				Type:     "NOT",
				Children: children,
			}
		}

	case "has_dlc":
		// has_dlc = "DLC Name"
		dlcName := p.extractValue(assign.Value)
		return &Condition{
			Type:  "has_dlc",
			Value: dlcName,
		}

	case "major_country":
		// major_country = yes
		return &Condition{
			Type:  "major_country",
			Value: "yes",
		}
	}

	return nil
}

// extractValue extracts string value from expression
func (p *TechnologyTagsParser) extractValue(expr Expression) string {
	switch v := expr.(type) {
	case *Identifier:
		return v.Value
	case *StringLiteral:
		return v.Value
	default:
		return ""
	}
}
