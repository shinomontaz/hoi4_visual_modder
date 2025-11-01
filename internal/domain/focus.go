package domain

// Focus represents a national focus in HOI4
type Focus struct {
	// Identification
	ID   string
	Icon string

	// Position
	Position             Position
	RelativePositionID   string // Optional: position relative to another focus
	
	// Dependencies
	Prerequisites      [][]string // Each inner slice is an OR group, outer is AND
	MutuallyExclusive  []string   // IDs of mutually exclusive focuses
	
	// Properties
	Cost                    int     // Days to complete
	Available               string  // Conditions (stored as raw string for now)
	Bypass                  string  // Bypass conditions
	CancelIfInvalid         bool
	ContinueIfInvalid       bool
	AvailableIfCapitulated  bool
	
	// Rewards and AI
	CompletionReward string   // Raw reward block
	AIWillDo         string   // AI priority block
	SearchFilters    []string // UI search filters
}

// NewFocus creates a new Focus with required fields
func NewFocus(id string, x, y int) *Focus {
	return &Focus{
		ID:            id,
		Position:      NewPosition(x, y),
		Cost:          70, // Default cost
		Prerequisites: make([][]string, 0),
	}
}

// Validate checks if the focus is valid
func (f *Focus) Validate() []string {
	errors := make([]string, 0)
	
	if f.ID == "" {
		errors = append(errors, "focus ID cannot be empty")
	}
	
	if f.Cost <= 0 {
		errors = append(errors, "focus cost must be positive")
	}
	
	return errors
}

// HasPrerequisite checks if this focus has any prerequisites
func (f *Focus) HasPrerequisite() bool {
	return len(f.Prerequisites) > 0
}

// IsMutuallyExclusiveWith checks if this focus is mutually exclusive with another
func (f *Focus) IsMutuallyExclusiveWith(otherID string) bool {
	for _, id := range f.MutuallyExclusive {
		if id == otherID {
			return true
		}
	}
	return false
}
