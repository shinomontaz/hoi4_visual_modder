package app

import (
	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
)

// State represents the application state
type State struct {
	// Project information
	ModPath     string
	CurrentMode Mode
	
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
