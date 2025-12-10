package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
// Priority: MOD technologies override GAME technologies by ID
func (tl *TechnologyLoader) LoadAllTechnologies() ([]*domain.Technology, error) {
	technologies := make(map[string]*domain.Technology) // key = tech ID

	// 1. Load from MOD first (highest priority)
	modTechCount := 0
	if tl.modPath != "" {
		modFileMap, err := tl.loadTechnologiesFromPath(tl.modPath)
		if err == nil && len(modFileMap) > 0 {
			for _, techs := range modFileMap {
				for _, tech := range techs {
					technologies[tech.ID] = tech
					modTechCount++
				}
			}
			println("Loaded", modTechCount, "technologies from MOD (", len(modFileMap), "files )")
		}
	}

	// 2. Load from GAME (fills missing technologies by ID, does NOT override mod)
	if tl.gamePath != "" {
		gameFileMap, err := tl.loadTechnologiesFromPath(tl.gamePath)
		if err == nil {
			addedFromGame := 0
			for _, techs := range gameFileMap {
				for _, tech := range techs {
					// Only add if not already present from mod
					if _, exists := technologies[tech.ID]; !exists {
						technologies[tech.ID] = tech
						addedFromGame++
					}
				}
			}
			println("Added", addedFromGame, "missing technologies from GAME")
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
func (tl *TechnologyLoader) loadTechnologiesFromPath(basePath string) (map[string][]*domain.Technology, error) {
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

	// Map: filename -> technologies
	fileMap := make(map[string][]*domain.Technology)

	// Parse each file
	for _, file := range files {
		techs, err := tl.parseTechnologyFile(file)
		if err != nil {
			// Skip files with errors
			println("Warning: Failed to parse", file, ":", err.Error())
			continue
		}
		fileName := filepath.Base(file)
		fileMap[fileName] = techs
	}

	return fileMap, nil
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

	// Debug: print file name and tech count only for electronics
	fileName := filepath.Base(filePath)
	if len(technologies) > 0 {
		isElectronics := false
		for _, tech := range technologies {
			if tech.Folder == "electronics_folder" {
				isElectronics = true
				break
			}
		}

		if isElectronics {
			println("Parsed", fileName, ":", len(technologies), "technologies")
			// Show first few techs with their positions
			for i, tech := range technologies {
				if i >= 3 {
					break
				}
				xStr := fmt.Sprintf("%d", tech.Position.X)
				if tech.Position.XVar != "" {
					xStr = fmt.Sprintf("%s=%d", tech.Position.XVar, tech.Position.X)
				}
				yStr := fmt.Sprintf("%d", tech.Position.Y)
				if tech.Position.YVar != "" {
					yStr = fmt.Sprintf("%s=%d", tech.Position.YVar, tech.Position.Y)
				}
				println("  -", tech.ID, "X:", xStr, "Y:", yStr)
			}
		}
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

// DetectSubTrees groups technologies into sub-trees by X-coordinate gaps
func (tl *TechnologyLoader) DetectSubTrees(folderName string, technologies []*domain.Technology) []*domain.SubTree {
	if len(technologies) == 0 {
		return nil
	}

	println("\n=== DetectSubTrees for folder:", folderName, "===")
	println("Total technologies:", len(technologies))

	// Step 1: Collect unique X values
	uniqueXValues := make(map[int]bool)
	for _, tech := range technologies {
		uniqueXValues[tech.Position.X] = true
		// Show coordinates with aliases
		xStr := fmt.Sprintf("%d", tech.Position.X)
		if tech.Position.XVar != "" {
			xStr = fmt.Sprintf("%s=%d", tech.Position.XVar, tech.Position.X)
		}
		yStr := fmt.Sprintf("%d", tech.Position.Y)
		if tech.Position.YVar != "" {
			yStr = fmt.Sprintf("%s=%d", tech.Position.YVar, tech.Position.Y)
		}
		println("  Tech:", tech.ID, "| X:", xStr, "| Y:", yStr)
	}

	// Step 2: Sort unique X values
	xValues := make([]int, 0, len(uniqueXValues))
	for x := range uniqueXValues {
		xValues = append(xValues, x)
	}
	sort.Ints(xValues)

	// Step 3: Find gaps and create sub-trees
	subTrees := make([]*domain.SubTree, 0)
	currentRange := []int{xValues[0]}

	for i := 1; i < len(xValues); i++ {
		gap := xValues[i] - xValues[i-1]

		if gap > 5 {
			// Gap found - create sub-tree for current range
			subTree := createSubTreeForRange(
				currentRange[0],
				currentRange[len(currentRange)-1],
				technologies,
				folderName,
			)
			subTrees = append(subTrees, subTree)

			// Start new range
			currentRange = []int{xValues[i]}
		} else {
			currentRange = append(currentRange, xValues[i])
		}
	}

	// Add last range
	if len(currentRange) > 0 {
		subTree := createSubTreeForRange(
			currentRange[0],
			currentRange[len(currentRange)-1],
			technologies,
			folderName,
		)
		subTrees = append(subTrees, subTree)
	}

	println("\nDetected", len(subTrees), "sub-trees:")
	for i, st := range subTrees {
		println(fmt.Sprintf("  SubTree %d: %s | X range: %d..%d | Technologies: %d",
			i+1, st.Name, st.XMin, st.XMax, len(st.Technologies)))
	}
	println("=== End DetectSubTrees ===\n")

	return subTrees
}

// createSubTreeForRange creates a sub-tree for a given X range
func createSubTreeForRange(
	xMin, xMax int,
	allTechs []*domain.Technology,
	folderName string,
) *domain.SubTree {
	// Filter technologies in range
	techs := make([]*domain.Technology, 0)
	categorySet := make(map[string]bool)
	varSet := make(map[string]bool)

	for _, tech := range allTechs {
		if tech.Position.X >= xMin && tech.Position.X <= xMax {
			techs = append(techs, tech)
			for _, cat := range tech.Categories {
				categorySet[cat] = true
			}
			// Collect X variables used in this range
			if tech.Position.XVar != "" {
				varSet[tech.Position.XVar] = true
			}
		}
	}

	categories := mapKeysToSlice(categorySet)
	name := identifySubTreeName(categorySet, folderName)

	// If we have variable names, add them to the name for clarity
	if len(varSet) > 0 {
		vars := mapKeysToSlice(varSet)
		sort.Strings(vars)
		if len(vars) <= 3 {
			// Show variable names if not too many
			varStr := strings.Join(vars, ", ")
			name = fmt.Sprintf("%s (%s)", name, varStr)
		}
	}

	return &domain.SubTree{
		Name:         name,
		XMin:         xMin,
		XMax:         xMax,
		Technologies: techs,
		Categories:   categories,
	}
}

// identifySubTreeName identifies sub-tree name by categories
func identifySubTreeName(categories map[string]bool, folderName string) string {
	categoryNames := map[string]string{
		"electronics":    "Electronic Engineering",
		"computing_tech": "Electronic Engineering",
		"radar_tech":     "Electronic Engineering",
		"radio_tech":     "Electronic Engineering",
		"rocketry":       "Experimental Rockets",
		"mot_rockets":    "Experimental Rockets",
		"jet_technology": "Jets & Aircraft Engines",
		"jet_engine":     "Jets & Aircraft Engines",
		"nuclear":        "Atomic Research",
		"land_doctrine":  "Land Doctrine",
		"air_doctrine":   "Air Doctrine",
		"naval_doctrine": "Naval Doctrine",
	}

	for cat := range categories {
		if name, ok := categoryNames[cat]; ok {
			return name
		}
	}

	// Fallback: use folder name
	return folderName
}

// mapKeysToSlice converts map keys to slice
func mapKeysToSlice(m map[string]bool) []string {
	result := make([]string, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result
}
