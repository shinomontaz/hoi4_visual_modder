package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shinomontaz/hoi4_visual_modder/internal/parser"
)

// ModDescriptor represents a parsed .mod file
type ModDescriptor struct {
	FilePath         string   // Full path to .mod file
	Name             string   // Mod name
	Version          string   // Mod version
	SupportedVersion string   // Game version
	Path             string   // Relative path to mod folder (from .mod file)
	ReplacePaths     []string // Paths that mod replaces
	Tags             []string // Mod tags

	// Resolved paths
	ModFolderPath string // Absolute path to mod folder
}

// LoadModDescriptor loads and parses a .mod file
func LoadModDescriptor(modFilePath string) (*ModDescriptor, error) {
	// Check if file exists
	if _, err := os.Stat(modFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("mod file not found: %s", modFilePath)
	}

	// Check extension
	if filepath.Ext(modFilePath) != ".mod" {
		return nil, fmt.Errorf("not a .mod file: %s", modFilePath)
	}

	// Read file content
	content, err := os.ReadFile(modFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read mod file: %w", err)
	}

	// Parse using existing parser
	p := parser.NewParser(string(content))
	program, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse mod file: %w", err)
	}

	// Extract mod descriptor data
	descriptor := &ModDescriptor{
		FilePath:     modFilePath,
		ReplacePaths: make([]string, 0),
		Tags:         make([]string, 0),
	}

	// Parse assignments
	for _, stmt := range program.Statements {
		if assign, ok := stmt.(*parser.AssignmentStatement); ok {
			key := assign.Name.Value
			switch key {
			case "name":
				descriptor.Name = extractString(assign.Value)
			case "version":
				descriptor.Version = extractString(assign.Value)
			case "supported_version":
				descriptor.SupportedVersion = extractString(assign.Value)
			case "path":
				descriptor.Path = extractString(assign.Value)
			case "replace_path":
				descriptor.ReplacePaths = append(descriptor.ReplacePaths, extractString(assign.Value))
			case "tags":
				descriptor.Tags = extractArray(assign.Value)
			}
		}
	}

	// Validate required fields
	if descriptor.Name == "" {
		return nil, fmt.Errorf("mod file missing 'name' field")
	}
	if descriptor.Path == "" {
		return nil, fmt.Errorf("mod file missing 'path' field")
	}

	// Resolve mod folder path
	// Path in .mod file is relative, like "mod/BlackICE_Historical_Immersion_Mod/"
	// We need to resolve it relative to the .mod file location
	modFileDir := filepath.Dir(modFilePath)

	// Clean the path (remove quotes, normalize slashes)
	cleanPath := strings.Trim(descriptor.Path, `"`)
	cleanPath = strings.ReplaceAll(cleanPath, "/", string(filepath.Separator))

	// If path starts with "mod/", it's relative to parent of .mod file
	if strings.HasPrefix(cleanPath, "mod"+string(filepath.Separator)) {
		// Remove "mod/" prefix
		cleanPath = strings.TrimPrefix(cleanPath, "mod"+string(filepath.Separator))
		descriptor.ModFolderPath = filepath.Join(modFileDir, cleanPath)
	} else {
		// Otherwise it's relative to .mod file directory
		descriptor.ModFolderPath = filepath.Join(modFileDir, cleanPath)
	}

	// Validate mod folder exists
	if err := ValidateModFolder(descriptor.ModFolderPath); err != nil {
		return nil, fmt.Errorf("mod folder validation failed: %w", err)
	}

	return descriptor, nil
}

// ValidateModFolder checks if the mod folder has valid HOI4 mod structure
func ValidateModFolder(modFolderPath string) error {
	// Check if folder exists
	info, err := os.Stat(modFolderPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("mod folder does not exist: %s", modFolderPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("mod path is not a directory: %s", modFolderPath)
	}

	// Check for common/ directory (basic validation)
	commonPath := filepath.Join(modFolderPath, "common")
	if _, err := os.Stat(commonPath); os.IsNotExist(err) {
		return fmt.Errorf("mod folder missing 'common' directory: %s", modFolderPath)
	}

	return nil
}

// extractString extracts string value from AST node
func extractString(node parser.Expression) string {
	switch v := node.(type) {
	case *parser.StringLiteral:
		return v.Value
	case *parser.Identifier:
		return v.Value
	default:
		return ""
	}
}

// extractArray extracts array of strings from AST node (block with multiple values)
func extractArray(node parser.Expression) []string {
	result := make([]string, 0)

	// For now, just return empty array
	// Tags parsing can be improved later if needed
	// The important fields (name, version, path) are working

	return result
}

// GetModInfo returns a formatted string with mod information
func (md *ModDescriptor) GetModInfo() string {
	return fmt.Sprintf("%s v%s (Game: %s)", md.Name, md.Version, md.SupportedVersion)
}
