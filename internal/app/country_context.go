package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
)

// CountryContext holds the context for a selected country
type CountryContext struct {
	Country     *domain.BookmarkCountry
	ModPath     string
	GamePath    string
	FocusPath   string   // Path to national focus file
	TechFolders []string // Available technology folders
}

// NewCountryContext creates a new country context
func NewCountryContext(country *domain.BookmarkCountry, modPath, gamePath string) *CountryContext {
	ctx := &CountryContext{
		Country:     country,
		ModPath:     modPath,
		GamePath:    gamePath,
		TechFolders: make([]string, 0),
	}

	// Resolve focus path
	ctx.resolveFocusPath()

	// Resolve tech folders
	ctx.resolveTechFolders()

	return ctx
}

// resolveFocusPath finds the national focus file for this country
func (ctx *CountryContext) resolveFocusPath() {
	tag := ctx.Country.Tag

	// Try mod path first
	if ctx.ModPath != "" {
		focusFile := filepath.Join(ctx.ModPath, "common", "national_focus", tag+"_focus.txt")
		if _, err := os.Stat(focusFile); err == nil {
			ctx.FocusPath = focusFile
			return
		}

		// Try alternative naming: lowercase
		focusFile = filepath.Join(ctx.ModPath, "common", "national_focus", tag+"_focus.txt")
		if _, err := os.Stat(focusFile); err == nil {
			ctx.FocusPath = focusFile
			return
		}
	}

	// Try game path
	if ctx.GamePath != "" {
		focusFile := filepath.Join(ctx.GamePath, "common", "national_focus", tag+"_focus.txt")
		if _, err := os.Stat(focusFile); err == nil {
			ctx.FocusPath = focusFile
			return
		}
	}

	// Focus file not found - that's okay, not all countries have custom focus trees
	ctx.FocusPath = ""
}

// resolveTechFolders finds available technology folders
func (ctx *CountryContext) resolveTechFolders() {
	// Technology folders are typically organized by category
	// Common folders: armor, infantry, artillery, support, naval, air, industry, electronics

	basePath := ctx.ModPath
	if basePath == "" {
		basePath = ctx.GamePath
	}

	if basePath == "" {
		return
	}

	techDir := filepath.Join(basePath, "common", "technologies")

	// Check if directory exists
	if _, err := os.Stat(techDir); os.IsNotExist(err) {
		return
	}

	// List all .txt files in technologies directory
	files, err := filepath.Glob(filepath.Join(techDir, "*.txt"))
	if err != nil {
		return
	}

	// Extract folder names (filenames without extension)
	for _, file := range files {
		filename := filepath.Base(file)
		folderName := filename[:len(filename)-4] // Remove .txt
		ctx.TechFolders = append(ctx.TechFolders, folderName)
	}
}

// GetFocusPath returns the path to the national focus file
func (ctx *CountryContext) GetFocusPath() (string, error) {
	if ctx.FocusPath == "" {
		return "", fmt.Errorf("no national focus file found for %s", ctx.Country.Tag)
	}
	return ctx.FocusPath, nil
}

// GetTechPath returns the path to a technology file by folder name
func (ctx *CountryContext) GetTechPath(folderName string) (string, error) {
	basePath := ctx.ModPath
	if basePath == "" {
		basePath = ctx.GamePath
	}

	if basePath == "" {
		return "", fmt.Errorf("no mod or game path set")
	}

	techFile := filepath.Join(basePath, "common", "technologies", folderName+".txt")

	// Check if file exists
	if _, err := os.Stat(techFile); os.IsNotExist(err) {
		return "", fmt.Errorf("technology file not found: %s", folderName)
	}

	return techFile, nil
}

// HasFocusTree returns true if the country has a custom focus tree
func (ctx *CountryContext) HasFocusTree() bool {
	return ctx.FocusPath != ""
}

// GetDisplayName returns the display name of the country
func (ctx *CountryContext) GetDisplayName() string {
	return ctx.Country.GetDisplayName()
}

// GetTag returns the country tag
func (ctx *CountryContext) GetTag() string {
	return ctx.Country.Tag
}
