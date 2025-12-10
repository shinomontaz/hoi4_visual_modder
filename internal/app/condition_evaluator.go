package app

import "github.com/shinomontaz/hoi4_visual_modder/internal/parser"

// ConditionEvaluator evaluates technology folder availability conditions
type ConditionEvaluator struct {
	countryFlags []string
	dlcs         []string // For future DLC support
	isMajor      bool     // For future major_country support
}

// NewConditionEvaluator creates a new condition evaluator
func NewConditionEvaluator(flags []string) *ConditionEvaluator {
	return &ConditionEvaluator{
		countryFlags: flags,
		dlcs:         make([]string, 0),
		isMajor:      false,
	}
}

// Evaluate checks if all conditions in AvailableCondition are met
func (e *ConditionEvaluator) Evaluate(condition *parser.AvailableCondition) bool {
	if condition == nil {
		return true // No conditions = always available
	}

	// All conditions must be true (AND logic)
	for _, cond := range condition.Conditions {
		if !e.evaluateCondition(cond) {
			return false
		}
	}

	return true
}

// evaluateCondition recursively evaluates a single condition
func (e *ConditionEvaluator) evaluateCondition(cond *parser.Condition) bool {
	switch cond.Type {
	case "has_country_flag":
		hasFlag := e.hasFlag(cond.Value)
		if cond.Negated {
			return !hasFlag
		}
		return hasFlag

	case "NOT":
		// NOT = { ... } means ALL children must be FALSE
		// If ANY child is TRUE, then NOT fails
		for _, child := range cond.Children {
			if e.evaluateCondition(child) {
				return false // One child is true, so NOT fails
			}
		}
		return true // All children are false, so NOT succeeds

	case "has_dlc":
		// For now, assume no DLCs installed
		hasDLC := e.hasDLC(cond.Value)
		if cond.Negated {
			return !hasDLC
		}
		return hasDLC

	case "major_country":
		// For now, assume not a major country
		if cond.Negated {
			return !e.isMajor
		}
		return e.isMajor

	default:
		// Unknown condition type - assume true to be permissive
		return true
	}
}

// hasFlag checks if country has a specific flag
func (e *ConditionEvaluator) hasFlag(flagName string) bool {
	for _, flag := range e.countryFlags {
		if flag == flagName {
			return true
		}
	}
	return false
}

// hasDLC checks if DLC is installed (placeholder for future)
func (e *ConditionEvaluator) hasDLC(dlcName string) bool {
	for _, dlc := range e.dlcs {
		if dlc == dlcName {
			return true
		}
	}
	return false
}
