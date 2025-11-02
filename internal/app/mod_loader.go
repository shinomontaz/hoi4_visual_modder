package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileType represents the type of HOI4 file
type FileType int

const (
	FileTypeUnknown FileType = iota
	FileTypeFocus
	FileTypeTechnology
)

// String returns string representation of FileType
func (ft FileType) String() string {
	switch ft {
	case FileTypeFocus:
		return "National Focus"
	case FileTypeTechnology:
		return "Technology"
	default:
		return "Unknown"
	}
}

// DetectBasePath extracts the mod root directory from a file path
// Example: "E:/mods/my_mod/common/national_focus/brazil.txt" → "E:/mods/my_mod"
func DetectBasePath(filePath string) (string, error) {
	// Normalize path separators
	filePath = filepath.Clean(filePath)

	// Split path into parts
	parts := strings.Split(filePath, string(filepath.Separator))

	// Find "common" directory in path
	commonIndex := -1
	for i, part := range parts {
		if strings.ToLower(part) == "common" {
			commonIndex = i
			break
		}
	}

	if commonIndex == -1 {
		return "", fmt.Errorf("file is not in a valid HOI4 mod structure (missing 'common' directory)")
	}

	// Base path is everything before "common"
	if commonIndex == 0 {
		return "", fmt.Errorf("invalid path: 'common' directory at root")
	}

	// Join all parts before "common"
	basePath := filepath.Join(parts[:commonIndex]...)
	
	// On Windows, filepath.Join removes the colon from drive letter (C: becomes C)
	// We need to restore it
	if len(parts) > 0 && len(parts[0]) >= 2 && parts[0][1] == ':' {
		// Drive letter exists, ensure it's properly formatted
		// filepath.Join already handles the path correctly on Windows
		basePath = parts[0] + string(filepath.Separator) + filepath.Join(parts[1:commonIndex]...)
	}

	return basePath, nil
}

// ValidateModStructure checks if the given path is a valid HOI4 mod directory
func ValidateModStructure(basePath string) error {
	// Check if base path exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return fmt.Errorf("mod directory does not exist: %s", basePath)
	}

	// Check for common/ directory
	commonPath := filepath.Join(basePath, "common")
	if _, err := os.Stat(commonPath); os.IsNotExist(err) {
		return fmt.Errorf("invalid mod structure: missing 'common' directory")
	}

	// Check for at least one of the expected subdirectories
	focusPath := filepath.Join(basePath, "common", "national_focus")
	techPath := filepath.Join(basePath, "common", "technologies")

	focusExists := false
	techExists := false

	_, err := os.Stat(focusPath)
	if err == nil {
		focusExists = true
	}
	_, err = os.Stat(techPath)
	if err == nil {
		techExists = true
	}

	if !focusExists && !techExists {
		return fmt.Errorf("invalid mod structure: missing 'national_focus' or 'technologies' directories")
	}

	return nil
}

// DetectFileType determines if the file is a focus or technology file
func DetectFileType(filePath string) (FileType, error) {
	// Normalize path
	filePath = filepath.Clean(filePath)
	normalizedPath := strings.ToLower(filePath)

	// Check if path contains national_focus
	if strings.Contains(normalizedPath, "national_focus") {
		return FileTypeFocus, nil
	}

	// Check if path contains technologies
	if strings.Contains(normalizedPath, "technologies") {
		return FileTypeTechnology, nil
	}

	return FileTypeUnknown, fmt.Errorf("cannot determine file type from path")
}

// LoadModFile loads a mod file and returns its content along with metadata
func LoadModFile(filePath string) (basePath string, fileType FileType, content string, err error) {
	// Detect base path
	basePath, err = DetectBasePath(filePath)
	if err != nil {
		return "", FileTypeUnknown, "", fmt.Errorf("failed to detect base path: %w", err)
	}

	// Validate mod structure
	err = ValidateModStructure(basePath)
	if err != nil {
		return "", FileTypeUnknown, "", fmt.Errorf("invalid mod structure: %w", err)
	}

	// Detect file type
	fileType, err = DetectFileType(filePath)
	if err != nil {
		return "", FileTypeUnknown, "", fmt.Errorf("failed to detect file type: %w", err)
	}

	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", FileTypeUnknown, "", fmt.Errorf("failed to read file: %w", err)
	}

	content = string(data)

	return basePath, fileType, content, nil
}
