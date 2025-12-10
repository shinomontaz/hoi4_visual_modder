package app

import (
	"os"
	"path/filepath"

	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
	"github.com/shinomontaz/hoi4_visual_modder/internal/parser"
)

// TechnologyLoader loads and filters technologies
type TechnologyLoader struct {
	modPath  string
	gamePath string
}

// NewTechnologyLoader creates a new technology loader
func NewTechnologyLoader(modPath, gamePath string) *TechnologyLoader {
	return &TechnologyLoader{
		modPath:  modPath,
		gamePath: gamePath,
	}
}

// LoadTechnologiesForFolder loads all technologies that belong to a specific folder
func (tl *TechnologyLoader) LoadTechnologiesForFolder(folderName string) ([]*domain.Technology, error) {
	// Load all technologies from mod and game
	allTechnologies, err := tl.LoadAllTechnologies()
	if err != nil {
		return nil, err
	}

	// Filter technologies by folder
	filtered := make([]*domain.Technology, 0)
	for _, tech := range allTechnologies {
		if tl.techBelongsToFolder(tech, folderName) {
			filtered = append(filtered, tech)
		}
	}

	return filtered, nil
}

// LoadAllTechnologies loads all technologies from mod and game
// Priority: MOD first, then GAME fills missing technologies
func (tl *TechnologyLoader) LoadAllTechnologies() ([]*domain.Technology, error) {
	technologies := make(map[string]*domain.Technology) // key = tech ID

	// 1. Load from MOD first (highest priority)
	modLoaded := false
	if tl.modPath != "" {
		modTechs, err := tl.loadTechnologiesFromPath(tl.modPath)
		if err == nil && len(modTechs) > 0 {
			modLoaded = true
			for _, tech := range modTechs {
				technologies[tech.ID] = tech
			}
			println("Loaded", len(modTechs), "technologies from MOD")
		}
	}

	// 2. Load from GAME (fills missing technologies, does NOT override mod)
	if tl.gamePath != "" {
		gameTechs, err := tl.loadTechnologiesFromPath(tl.gamePath)
		if err == nil {
			addedFromGame := 0
			for _, tech := range gameTechs {
				// Only add if not already present from mod
				if _, exists := technologies[tech.ID]; !exists {
					technologies[tech.ID] = tech
					addedFromGame++
				}
			}
			if modLoaded {
				println("Added", addedFromGame, "missing technologies from GAME")
			} else {
				println("Loaded", len(gameTechs), "technologies from GAME")
			}
		}
	}

	// Convert map to slice
	result := make([]*domain.Technology, 0, len(technologies))
	for _, tech := range technologies {
		result = append(result, tech)
	}

	return result, nil
}

// loadTechnologiesFromPath loads technologies from a specific path
func (tl *TechnologyLoader) loadTechnologiesFromPath(basePath string) ([]*domain.Technology, error) {
	techDir := filepath.Join(basePath, "common", "technologies")

	// Check if directory exists
	if _, err := os.Stat(techDir); os.IsNotExist(err) {
		return nil, err
	}

	// Find all .txt files
	files, err := filepath.Glob(filepath.Join(techDir, "*.txt"))
	if err != nil {
		return nil, err
	}

	allTechs := make([]*domain.Technology, 0)

	// Parse each file
	for _, file := range files {
		techs, err := tl.parseTechnologyFile(file)
		if err != nil {
			// Skip files with errors
			println("Warning: Failed to parse", file, ":", err.Error())
			continue
		}
		allTechs = append(allTechs, techs...)
	}

	return allTechs, nil
}

// parseTechnologyFile parses a single technology file
func (tl *TechnologyLoader) parseTechnologyFile(filePath string) ([]*domain.Technology, error) {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse with lexer and parser
	p := parser.NewParser(string(content))
	program, err := p.Parse()
	if err != nil {
		return nil, err
	}

	// Convert AST to technologies
	techParser := parser.NewTechParser()
	technologies, err := techParser.ParseTechnologies(program)
	if err != nil {
		return nil, err
	}

	return technologies, nil
}

// techBelongsToFolder checks if a technology belongs to a specific folder
func (tl *TechnologyLoader) techBelongsToFolder(tech *domain.Technology, folderName string) bool {
	// Check if technology has this folder
	// Technology can belong to multiple folders
	// We need to check the Folder field in domain.Technology

	// For now, check if the Folder field matches
	// Note: Technology struct has a Folder field which is the folder name
	return tech.Folder == folderName
}
