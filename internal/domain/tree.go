package domain

// FocusTree represents a complete national focus tree
type FocusTree struct {
	ID                       string
	Country                  string // Country tag or filter
	ContinuousFocusPosition  *Position
	Default                  bool
	ResetOnCivilWar          bool
	Focuses                  map[string]*Focus // ID -> Focus
}

// NewFocusTree creates a new FocusTree
func NewFocusTree(id string) *FocusTree {
	return &FocusTree{
		ID:      id,
		Focuses: make(map[string]*Focus),
	}
}

// AddFocus adds a focus to the tree
func (ft *FocusTree) AddFocus(focus *Focus) {
	ft.Focuses[focus.ID] = focus
}

// GetFocus retrieves a focus by ID
func (ft *FocusTree) GetFocus(id string) (*Focus, bool) {
	focus, exists := ft.Focuses[id]
	return focus, exists
}

// Validate validates the entire focus tree
func (ft *FocusTree) Validate() []string {
	errors := make([]string, 0)
	
	// Validate each focus
	for _, focus := range ft.Focuses {
		focusErrors := focus.Validate()
		errors = append(errors, focusErrors...)
	}
	
	// Check for circular dependencies
	if circularErrors := ft.checkCircularDependencies(); len(circularErrors) > 0 {
		errors = append(errors, circularErrors...)
	}
	
	// Check for invalid prerequisites
	if prereqErrors := ft.checkPrerequisites(); len(prereqErrors) > 0 {
		errors = append(errors, prereqErrors...)
	}
	
	// Check for position conflicts
	if posErrors := ft.checkPositionConflicts(); len(posErrors) > 0 {
		errors = append(errors, posErrors...)
	}
	
	return errors
}

// checkCircularDependencies detects circular dependency chains
func (ft *FocusTree) checkCircularDependencies() []string {
	errors := make([]string, 0)
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	
	var visit func(id string) bool
	visit = func(id string) bool {
		visited[id] = true
		recStack[id] = true
		
		focus, exists := ft.Focuses[id]
		if !exists {
			return false
		}
		
		for _, prereqGroup := range focus.Prerequisites {
			for _, prereqID := range prereqGroup {
				if !visited[prereqID] {
					if visit(prereqID) {
						return true
					}
				} else if recStack[prereqID] {
					errors = append(errors, "circular dependency detected involving: "+id+" -> "+prereqID)
					return true
				}
			}
		}
		
		recStack[id] = false
		return false
	}
	
	for id := range ft.Focuses {
		if !visited[id] {
			visit(id)
		}
	}
	
	return errors
}

// checkPrerequisites validates that all prerequisite IDs exist
func (ft *FocusTree) checkPrerequisites() []string {
	errors := make([]string, 0)
	
	for _, focus := range ft.Focuses {
		for _, prereqGroup := range focus.Prerequisites {
			for _, prereqID := range prereqGroup {
				if _, exists := ft.Focuses[prereqID]; !exists {
					errors = append(errors, "focus "+focus.ID+" references non-existent prerequisite: "+prereqID)
				}
			}
		}
		
		for _, exclusiveID := range focus.MutuallyExclusive {
			if _, exists := ft.Focuses[exclusiveID]; !exists {
				errors = append(errors, "focus "+focus.ID+" references non-existent mutually exclusive focus: "+exclusiveID)
			}
		}
	}
	
	return errors
}

// checkPositionConflicts detects focuses at the same position
func (ft *FocusTree) checkPositionConflicts() []string {
	errors := make([]string, 0)
	positions := make(map[Position][]string)
	
	for _, focus := range ft.Focuses {
		pos := focus.Position
		positions[pos] = append(positions[pos], focus.ID)
	}
	
	for pos, ids := range positions {
		if len(ids) > 1 {
			errors = append(errors, "position conflict at ("+string(rune(pos.X))+","+string(rune(pos.Y))+"): multiple focuses")
		}
	}
	
	return errors
}

// TechnologyTree represents a collection of technologies
type TechnologyTree struct {
	Technologies map[string]*Technology // ID -> Technology
	Folders      map[string][]string    // Folder -> Technology IDs
}

// NewTechnologyTree creates a new TechnologyTree
func NewTechnologyTree() *TechnologyTree {
	return &TechnologyTree{
		Technologies: make(map[string]*Technology),
		Folders:      make(map[string][]string),
	}
}

// AddTechnology adds a technology to the tree
func (tt *TechnologyTree) AddTechnology(tech *Technology) {
	tt.Technologies[tech.ID] = tech
	tt.Folders[tech.Folder] = append(tt.Folders[tech.Folder], tech.ID)
}

// GetTechnology retrieves a technology by ID
func (tt *TechnologyTree) GetTechnology(id string) (*Technology, bool) {
	tech, exists := tt.Technologies[id]
	return tech, exists
}

// Validate validates the entire technology tree
func (tt *TechnologyTree) Validate() []string {
	errors := make([]string, 0)
	
	// Validate each technology
	for _, tech := range tt.Technologies {
		techErrors := tech.Validate()
		errors = append(errors, techErrors...)
	}
	
	// Check for invalid paths
	if pathErrors := tt.checkPaths(); len(pathErrors) > 0 {
		errors = append(errors, pathErrors...)
	}
	
	return errors
}

// checkPaths validates that all path targets exist
func (tt *TechnologyTree) checkPaths() []string {
	errors := make([]string, 0)
	
	for _, tech := range tt.Technologies {
		for _, path := range tech.Paths {
			if _, exists := tt.Technologies[path.LeadsToTech]; !exists {
				errors = append(errors, "technology "+tech.ID+" has path to non-existent tech: "+path.LeadsToTech)
			}
		}
		
		for _, exclusiveID := range tech.XOR {
			if _, exists := tt.Technologies[exclusiveID]; !exists {
				errors = append(errors, "technology "+tech.ID+" references non-existent exclusive tech: "+exclusiveID)
			}
		}
	}
	
	return errors
}
