package app

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo represents metadata about a discovered file
type FileInfo struct {
	Path         string    // Full absolute path
	Name         string    // Filename without path
	RelativePath string    // Path relative to Base_path
	Size         int64     // File size in bytes
	ModTime      time.Time // Last modification time
	Category     string    // "focus" or "technology"
}

// FileScanner scans mod directories for .txt files
type FileScanner struct {
	basePath string
}

// NewFileScanner creates a new FileScanner
func NewFileScanner(basePath string) *FileScanner {
	return &FileScanner{
		basePath: basePath,
	}
}

// ScanAll scans both focus and technology directories
func (fs *FileScanner) ScanAll() ([]FileInfo, error) {
	files := make([]FileInfo, 0)
	
	// Scan national focus files
	focusFiles, err := fs.scanDirectory("common/national_focus", "focus")
	if err == nil {
		files = append(files, focusFiles...)
	}
	
	// Scan technology files
	techFiles, err := fs.scanDirectory("common/technologies", "technology")
	if err == nil {
		files = append(files, techFiles...)
	}
	
	return files, nil
}

// ScanFocusFiles scans only national focus directory
func (fs *FileScanner) ScanFocusFiles() ([]FileInfo, error) {
	return fs.scanDirectory("common/national_focus", "focus")
}

// ScanTechnologyFiles scans only technology directory
func (fs *FileScanner) ScanTechnologyFiles() ([]FileInfo, error) {
	return fs.scanDirectory("common/technologies", "technology")
}

// scanDirectory scans a specific directory for .txt files
func (fs *FileScanner) scanDirectory(relativeDir string, category string) ([]FileInfo, error) {
	fullPath := filepath.Join(fs.basePath, relativeDir)
	
	// Check if directory exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return []FileInfo{}, nil // Return empty list if directory doesn't exist
	}
	
	files := make([]FileInfo, 0)
	
	// Walk through directory
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files with errors
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Only process .txt files
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".txt") {
			return nil
		}
		
		// Get relative path from base
		relPath, err := filepath.Rel(fs.basePath, path)
		if err != nil {
			relPath = path
		}
		
		files = append(files, FileInfo{
			Path:         path,
			Name:         info.Name(),
			RelativePath: relPath,
			Size:         info.Size(),
			ModTime:      info.ModTime(),
			Category:     category,
		})
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return files, nil
}

// ValidateModDirectory checks if the path looks like a valid HOI4 mod directory
func ValidateModDirectory(path string) bool {
	// Check for common/ directory
	commonPath := filepath.Join(path, "common")
	if _, err := os.Stat(commonPath); os.IsNotExist(err) {
		return false
	}
	
	// Check for at least one of the expected subdirectories
	focusPath := filepath.Join(path, "common", "national_focus")
	techPath := filepath.Join(path, "common", "technologies")
	
	focusExists := false
	techExists := false
	
	if _, err := os.Stat(focusPath); err == nil {
		focusExists = true
	}
	
	if _, err := os.Stat(techPath); err == nil {
		techExists = true
	}
	
	return focusExists || techExists
}
