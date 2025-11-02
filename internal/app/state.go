package app

import (
	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
)

// State represents the application state
type State struct {
	// Project information
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
	return &State{
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
