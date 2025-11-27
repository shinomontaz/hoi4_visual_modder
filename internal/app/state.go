package app

import (
	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
)

// State represents the application state
type State struct {
	// Configuration
	Config *AppConfig

	// Mod and Game installations
	ModDescriptor    *ModDescriptor
	GameInstallation *GameInstallation

	// Country context
	CountryContext *CountryContext

	// Project information (legacy, will be replaced)
	ModPath          string
	BasePath         string   // Root directory of the mod
	SelectedFilePath string   // Full path to selected file
	FileType         FileType // Type of selected file
	CurrentMode      Mode

	// File browser state
	AvailableFiles []FileInfo
	SelectedFile   *FileInfo
	FileContent    string

	// Loaded data
	FocusTree      *domain.FocusTree
	TechnologyTree *domain.TechnologyTree

	// UI state
	SelectedNodeID string
	CameraX        float64
	CameraY        float64
	Zoom           float64
}

// Mode represents the editing mode
type Mode int

const (
	ModeNone Mode = iota
	ModeFocusTree
	ModeTechnology
)

// NewState creates a new application state
func NewState() *State {
	// Load config
	config, err := LoadConfig()
	if err != nil {
		// If config load fails, use default
		config = DefaultConfig()
	}

	return &State{
		Config:      config,
		CurrentMode: ModeNone,
		Zoom:        1.0,
	}
}

// SetModPath sets the mod directory path
func (s *State) SetModPath(path string) {
	s.ModPath = path
}

// SetMode sets the current editing mode
func (s *State) SetMode(mode Mode) {
	s.CurrentMode = mode
}

// LoadFocusTree loads a focus tree
func (s *State) LoadFocusTree(tree *domain.FocusTree) {
	s.FocusTree = tree
	s.CurrentMode = ModeFocusTree
}

// LoadTechnologyTree loads a technology tree
func (s *State) LoadTechnologyTree(tree *domain.TechnologyTree) {
	s.TechnologyTree = tree
	s.CurrentMode = ModeTechnology
}

// SetAvailableFiles sets the list of available files
func (s *State) SetAvailableFiles(files []FileInfo) {
	s.AvailableFiles = files
}

// SelectFile selects a file and loads its content
func (s *State) SelectFile(file *FileInfo) {
	s.SelectedFile = file
}

// SetFileContent sets the content of the selected file
func (s *State) SetFileContent(content string) {
	s.FileContent = content
}

// SetBasePath sets the base path of the mod
func (s *State) SetBasePath(path string) {
	s.BasePath = path
}

// LoadFile loads a file and sets all related state
func (s *State) LoadFile(filePath string) error {
	basePath, fileType, content, err := LoadModFile(filePath)
	if err != nil {
		return err
	}

	s.SelectedFilePath = filePath
	s.BasePath = basePath
	s.FileType = fileType
	s.FileContent = content

	// Also set ModPath for backwards compatibility
	s.ModPath = basePath

	return nil
}

// SetModDescriptor sets the mod descriptor and updates config
func (s *State) SetModDescriptor(mod *ModDescriptor) error {
	s.ModDescriptor = mod

	// Update legacy fields for backwards compatibility
	s.BasePath = mod.ModFolderPath
	s.ModPath = mod.ModFolderPath

	// Save to config
	if s.Config != nil {
		return s.Config.UpdateModPath(mod.FilePath)
	}

	return nil
}

// SetGameInstallation sets the game installation and updates config
func (s *State) SetGameInstallation(game *GameInstallation) error {
	s.GameInstallation = game

	// Save to config
	if s.Config != nil {
		return s.Config.UpdateGamePath(game.Path)
	}

	return nil
}

// GetModPath returns the mod folder path (new or legacy)
func (s *State) GetModPath() string {
	if s.ModDescriptor != nil {
		return s.ModDescriptor.ModFolderPath
	}
	return s.BasePath
}

// GetGamePath returns the game installation path
func (s *State) GetGamePath() string {
	if s.GameInstallation != nil {
		return s.GameInstallation.Path
	}
	if s.Config != nil {
		return s.Config.GamePath
	}
	return ""
}

// SetCountryContext sets the country context
func (s *State) SetCountryContext(country *domain.BookmarkCountry) {
	modPath := s.GetModPath()
	gamePath := s.GetGamePath()

	s.CountryContext = NewCountryContext(country, modPath, gamePath)

	// Save to config
	if s.Config != nil {
		s.Config.UpdateLastCountry(country.Tag)
	}
}

// GetCountryContext returns the current country context
func (s *State) GetCountryContext() *CountryContext {
	return s.CountryContext
}
