package domain

// Technology represents a technology in HOI4
type Technology struct {
	// Identification
	ID string
	
	// Position
	Position Position
	Folder   string // Folder/category name
	
	// Access and Prerequisites
	Allow      string   // Conditions for technology to appear
	Categories []string // Technology categories
	
	// Effects
	Effects map[string]map[string]float64 // category -> modifier -> value
	
	// Paths (connections to other techs)
	Paths []TechPath
	
	// Research Properties
	ResearchCost      float64
	XPResearchType    string  // army/navy/air
	XPBoostCost       int
	XPResearchBonus   float64
	
	// Completion
	OnResearchComplete string // Effects on completion
	
	// AI
	AIWillDo         string
	AIResearchWeights map[string]float64
	
	// Special
	XOR            []string // Mutually exclusive technologies
	EnableTactic   string
	EnableBuilding string
}

// TechPath represents a connection to another technology
type TechPath struct {
	LeadsToTech      string  // Target technology ID
	ResearchCostCoeff float64 // Cost multiplier for this path
}

// NewTechnology creates a new Technology with required fields
func NewTechnology(id string, x, y int, folder string) *Technology {
	return &Technology{
		ID:                id,
		Position:          NewPosition(x, y),
		Folder:            folder,
		Categories:        make([]string, 0),
		Effects:           make(map[string]map[string]float64),
		Paths:             make([]TechPath, 0),
		ResearchCost:      1.0,
		AIResearchWeights: make(map[string]float64),
		XOR:               make([]string, 0),
	}
}

// Validate checks if the technology is valid
func (t *Technology) Validate() []string {
	errors := make([]string, 0)
	
	if t.ID == "" {
		errors = append(errors, "technology ID cannot be empty")
	}
	
	if t.ResearchCost <= 0 {
		errors = append(errors, "research cost must be positive")
	}
	
	if t.Folder == "" {
		errors = append(errors, "technology must have a folder")
	}
	
	return errors
}

// AddPath adds a path to another technology
func (t *Technology) AddPath(targetID string, costCoeff float64) {
	t.Paths = append(t.Paths, TechPath{
		LeadsToTech:      targetID,
		ResearchCostCoeff: costCoeff,
	})
}

// AddEffect adds an effect modifier
func (t *Technology) AddEffect(category, modifier string, value float64) {
	if t.Effects[category] == nil {
		t.Effects[category] = make(map[string]float64)
	}
	t.Effects[category][modifier] = value
}

// IsExclusiveWith checks if this tech is mutually exclusive with another
func (t *Technology) IsExclusiveWith(otherID string) bool {
	for _, id := range t.XOR {
		if id == otherID {
			return true
		}
	}
	return false
}
