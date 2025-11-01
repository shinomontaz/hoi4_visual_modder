package paradox

// Common types and utilities for Paradox scripting language

// Block represents a key-value block in Paradox scripts
type Block struct {
	Key   string
	Value interface{} // Can be string, number, or another Block
}

// IsKeyword checks if a string is a Paradox keyword
func IsKeyword(s string) bool {
	keywords := map[string]bool{
		"focus":                    true,
		"focus_tree":               true,
		"technologies":             true,
		"prerequisite":             true,
		"mutually_exclusive":       true,
		"available":                true,
		"bypass":                   true,
		"completion_reward":        true,
		"ai_will_do":               true,
		"allow":                    true,
		"path":                     true,
		"leads_to_tech":            true,
		"research_cost":            true,
		"on_research_complete":     true,
	}
	return keywords[s]
}
