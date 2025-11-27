package domain

// Bookmark represents a scenario bookmark (e.g., 1936, 1939)
type Bookmark struct {
	Name           string             // "GATHERING_STORM_NAME"
	Description    string             // Description key
	Date           string             // "1936.1.1.12"
	DefaultCountry string             // "GER"
	Countries      []*BookmarkCountry // List of playable countries
}

// BookmarkCountry represents a country in a bookmark
type BookmarkCountry struct {
	Tag      string   // "GER", "SOV", "USA", etc.
	Name     string   // Display name (or tag if no localization)
	History  string   // Description key
	Ideology string   // "fascism", "communism", "democratic", etc.
	IsMajor  bool     // Major power flag
	Ideas    []string // Starting national ideas
	Focuses  []string // Starting focuses
}

// GetDisplayName returns the display name or tag
func (bc *BookmarkCountry) GetDisplayName() string {
	if bc.Name != "" {
		return bc.Name
	}
	return bc.Tag
}

// GetTypeLabel returns "Major" or "Minor"
func (bc *BookmarkCountry) GetTypeLabel() string {
	if bc.IsMajor {
		return "Major"
	}
	return "Minor"
}
