package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
	"github.com/shinomontaz/hoi4_visual_modder/internal/parser"
)

// CountryContext holds the context for a selected country
type CountryContext struct {
	Country         *domain.BookmarkCountry
	ModPath         string
	GamePath        string
	FocusPath       string               // Path to national focus file
	TechFolders     []string             // Available technology folders (IDs)
	Localizations   map[string]string    // Localized strings
	AllTechnologies []*domain.Technology // All loaded technologies (cached)
}

// NewCountryContext creates a new country context
func NewCountryContext(country *domain.BookmarkCountry, modPath, gamePath string) *CountryContext {
	ctx := &CountryContext{
		Country:       country,
		ModPath:       modPath,
		GamePath:      gamePath,
		TechFolders:   make([]string, 0),
		Localizations: make(map[string]string),
	}

	// Load localizations
	ctx.loadLocalizations()

	// Resolve focus path
	ctx.resolveFocusPath()

	// Resolve tech folders
	ctx.resolveTechFolders()

	// Load all technologies once (cached)
	ctx.loadAllTechnologies()

	return ctx
}

// loadAllTechnologies loads all technologies once and caches them
func (ctx *CountryContext) loadAllTechnologies() {
	loader := NewTechnologyLoader(ctx.ModPath, ctx.GamePath)
	technologies, err := loader.LoadAllTechnologies()
	if err != nil {
		println("Warning: Failed to load technologies:", err.Error())
		ctx.AllTechnologies = make([]*domain.Technology, 0)
		return
	}

	ctx.AllTechnologies = technologies
	println("Loaded", len(technologies), "technologies into cache")
}

// loadLocalizations loads localization strings
func (ctx *CountryContext) loadLocalizations() {
	locParser := parser.NewLocalizationParser(ctx.ModPath, ctx.GamePath, "english")
	localizations, err := locParser.LoadLocalizations()
	if err != nil {
		println("Warning: Failed to load localizations:", err.Error())
		ctx.Localizations = make(map[string]string)
		return
	}

	ctx.Localizations = localizations
	println("Loaded", len(localizations), "localization strings")
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

// resolveTechFolders finds available technology folders from technology_tags
func (ctx *CountryContext) resolveTechFolders() {
	// Parse technology_folders from technology_tags files
	tagsParser := parser.NewTechnologyTagsParser(ctx.GamePath, ctx.ModPath)
	folders, err := tagsParser.ParseTechnologyFolders()
	if err != nil {
		println("Warning: Failed to parse technology_folders:", err.Error())
		// Fallback: if parsing fails, return empty list
		ctx.TechFolders = make([]string, 0)
		return
	}

	println("Found", len(folders), "technology folders")
	ctx.TechFolders = folders
}

// GetFocusPath returns the path to the national focus file
func (ctx *CountryContext) GetFocusPath() (string, error) {
	if ctx.FocusPath == "" {
		return "", fmt.Errorf("no national focus file found for %s", ctx.Country.Tag)
	}
	return ctx.FocusPath, nil
}

// LoadTechnologiesForFolder filters cached technologies for a specific folder
func (ctx *CountryContext) LoadTechnologiesForFolder(folderName string) ([]*domain.Technology, error) {
	// Filter from cached technologies
	filtered := make([]*domain.Technology, 0)
	for _, tech := range ctx.AllTechnologies {
		if tech.Folder == folderName {
			filtered = append(filtered, tech)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no technologies found for folder: %s", folderName)
	}

	println("Filtered", len(filtered), "technologies for folder:", folderName)
	return filtered, nil
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

// GetLocalizedFolderName returns the localized name for a technology folder
func (ctx *CountryContext) GetLocalizedFolderName(folderID string) string {
	// Build localization key: folder_id + "_name"
	locKey := folderID + "_name"

	// Try to get localized name
	if localizedName, ok := ctx.Localizations[locKey]; ok {
		return localizedName
	}

	// Fallback: format folder ID (remove "_folder" suffix and capitalize)
	displayName := folderID
	if len(displayName) > 7 && displayName[len(displayName)-7:] == "_folder" {
		displayName = displayName[:len(displayName)-7]
	}

	// Capitalize first letter
	if len(displayName) > 0 && displayName[0] >= 'a' && displayName[0] <= 'z' {
		displayName = string(displayName[0]-32) + displayName[1:]
	}

	return displayName
}
