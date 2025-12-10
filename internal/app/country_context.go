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
	CountryFlags    []string             // Country flags from history files
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

	// Load country flags (needed for folder filtering)
	ctx.loadCountryFlags()

	// Load localizations
	ctx.loadLocalizations()

	// Resolve focus path
	ctx.resolveFocusPath()

	// Resolve tech folders (uses country flags)
	ctx.resolveTechFolders()

	// Load all technologies once (cached)
	ctx.loadAllTechnologies()

	return ctx
}

// loadCountryFlags loads flags from history/countries files
func (ctx *CountryContext) loadCountryFlags() {
	flagsParser := parser.NewCountryFlagsParser(ctx.GamePath, ctx.ModPath)
	flags, err := flagsParser.ParseCountryFlags(ctx.Country.Tag)
	if err != nil {
		println("Warning: Failed to load country flags:", err.Error())
		ctx.CountryFlags = make([]string, 0)
		return
	}

	ctx.CountryFlags = flags

	// If this is a major country, add all UNLOCK:* flags for technology folders
	if ctx.Country.IsMajor {
		unlockFlags := []string{
			"UNLOCK:electronics_folder",
			"UNLOCK:nuclear_folder",
			"UNLOCK:infantry_folder",
			"UNLOCK:support_folder",
			"UNLOCK:armor_folder",
			"UNLOCK:artillery_folder",
			"UNLOCK:land_doctrine_folder",
			"UNLOCK:air_doctrine_folder",
			"UNLOCK:naval_doctrine_folder",
		}

		// Add unlock flags if not already present
		for _, unlockFlag := range unlockFlags {
			hasFlag := false
			for _, existingFlag := range ctx.CountryFlags {
				if existingFlag == unlockFlag {
					hasFlag = true
					break
				}
			}
			if !hasFlag {
				ctx.CountryFlags = append(ctx.CountryFlags, unlockFlag)
			}
		}
		println("Loaded", len(flags), "country flags for", ctx.Country.Tag, "(major: added UNLOCK flags)")
	} else {
		println("Loaded", len(flags), "country flags for", ctx.Country.Tag)
	}
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
// Filters folders based on country flags and overlay logic
func (ctx *CountryContext) resolveTechFolders() {
	// Parse technology_folders with detailed information
	tagsParser := parser.NewTechnologyTagsParser(ctx.GamePath, ctx.ModPath)
	allFolders, err := tagsParser.ParseTechnologyFoldersDetailed()
	if err != nil {
		println("Warning: Failed to parse technology_folders:", err.Error())
		ctx.TechFolders = make([]string, 0)
		return
	}

	println("Found", len(allFolders), "total technology folders")

	// Create condition evaluator
	evaluator := NewConditionEvaluator(ctx.CountryFlags)

	// Filter folders based on availability conditions
	availableFolders := make([]string, 0)
	overlayMap := make(map[string]bool) // base_folder -> overlay_active

	for _, folder := range allFolders {
		if folder.IsOverlay {
			// Check if overlay should be shown
			baseName := folder.Name[:len(folder.Name)-len("_overlay_folder")]
			if folder.Available != nil {
				// Overlay shown when condition is FALSE (i.e., folder is locked)
				overlayActive := !evaluator.Evaluate(folder.Available)
				overlayMap[baseName] = overlayActive
			}
		} else {
			// Regular folder - check if available
			if folder.Available == nil || evaluator.Evaluate(folder.Available) {
				availableFolders = append(availableFolders, folder.Name)
			}
		}
	}

	// Remove folders that have active overlays (locked folders)
	filtered := make([]string, 0)
	for _, folderName := range availableFolders {
		if !overlayMap[folderName] {
			filtered = append(filtered, folderName)
		}
	}

	ctx.TechFolders = filtered
	println("Available technology folders:", len(filtered), "out of", len(allFolders))
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

// HasFlag checks if country has a specific flag
func (ctx *CountryContext) HasFlag(flagName string) bool {
	for _, flag := range ctx.CountryFlags {
		if flag == flagName {
			return true
		}
	}
	return false
}

// HasAnyFlag checks if country has any of the specified flags
func (ctx *CountryContext) HasAnyFlag(flags []string) bool {
	for _, flag := range flags {
		if ctx.HasFlag(flag) {
			return true
		}
	}
	return false
}

// GetUnlockFlags returns only UNLOCK:* flags
func (ctx *CountryContext) GetUnlockFlags() []string {
	return parser.GetUnlockFlags(ctx.CountryFlags)
}

// GetCountrySpecificFlags returns flags specific to this country (e.g., GER_air, SOV_armor)
func (ctx *CountryContext) GetCountrySpecificFlags() []string {
	return parser.GetCountrySpecificFlags(ctx.CountryFlags, ctx.Country.Tag)
}
