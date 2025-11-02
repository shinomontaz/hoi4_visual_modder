package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
)

// FocusParser converts AST to domain.Focus models
type FocusParser struct {
	variables map[string]string // Variable definitions (@X = 5)
}

// NewFocusParser creates a new FocusParser
func NewFocusParser() *FocusParser {
	return &FocusParser{
		variables: make(map[string]string),
	}
}

// ParseFocusTree parses a focus_tree block from AST
func (fp *FocusParser) ParseFocusTree(program *Program) ([]*domain.Focus, error) {
	if len(program.Statements) == 0 {
		return nil, fmt.Errorf("empty program")
	}
	
	// First pass: collect all variable definitions
	for _, stmt := range program.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		if strings.HasPrefix(assignStmt.Name.Value, "@") {
			fp.handleVariableDefinition(assignStmt)
		}
	}
	
	focuses := make([]*domain.Focus, 0)
	
	// Second pass: find and parse the focus_tree block
	for _, stmt := range program.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		// Check if it's the focus_tree block
		if assignStmt.Name.Value == "focus_tree" {
			block, ok := assignStmt.Value.(*BlockStatement)
			if !ok {
				return nil, fmt.Errorf("focus_tree value is not a block")
			}
			
			// Collect variables inside focus_tree block first
			for _, focusStmt := range block.Statements {
				focusAssign, ok := focusStmt.(*AssignmentStatement)
				if !ok {
					continue
				}
				
				if strings.HasPrefix(focusAssign.Name.Value, "@") {
					fp.handleVariableDefinition(focusAssign)
				}
			}
			
			// Now parse each focus in the block
			for _, focusStmt := range block.Statements {
				focusAssign, ok := focusStmt.(*AssignmentStatement)
				if !ok {
					continue
				}
				
				// Skip variable definitions and metadata (id, country, etc.)
				if strings.HasPrefix(focusAssign.Name.Value, "@") ||
					focusAssign.Name.Value == "id" ||
					focusAssign.Name.Value == "country" ||
					focusAssign.Name.Value == "continuous_focus_position" ||
					focusAssign.Name.Value == "reset_on_civilwar" {
					continue
				}
				
				// Parse focus
				if focusAssign.Name.Value == "focus" {
					focusBlock, ok := focusAssign.Value.(*BlockStatement)
					if !ok {
						continue // Not a focus definition, skip
					}
					
					focus, err := fp.parseFocus(focusBlock)
					if err != nil {
						return nil, fmt.Errorf("failed to parse focus: %w", err)
					}
					
					focuses = append(focuses, focus)
				}
			}
			
			break
		}
	}
	
	return focuses, nil
}

// handleVariableDefinition stores a variable definition
func (fp *FocusParser) handleVariableDefinition(stmt *AssignmentStatement) {
	varName := stmt.Name.Value
	
	// Get the value
	var value string
	switch v := stmt.Value.(type) {
	case *NumberLiteral:
		value = v.Value
	case *StringLiteral:
		value = v.Value
	case *Identifier:
		value = v.Value
	default:
		return
	}
	
	fp.variables[varName] = value
}

// parseFocus parses a single focus block
func (fp *FocusParser) parseFocus(block *BlockStatement) (*domain.Focus, error) {
	focus := &domain.Focus{
		Cost:          70, // Default
		Prerequisites: make([][]string, 0),
		MutuallyExclusive: make([]string, 0),
		SearchFilters: make([]string, 0),
	}
	
	// Parse each field in the focus block
	for _, stmt := range block.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		fieldName := assignStmt.Name.Value
		
		switch fieldName {
		case "id":
			if str, ok := assignStmt.Value.(*StringLiteral); ok {
				focus.ID = str.Value
			} else if id, ok := assignStmt.Value.(*Identifier); ok {
				focus.ID = id.Value
			}
			
		case "icon":
			if str, ok := assignStmt.Value.(*StringLiteral); ok {
				focus.Icon = str.Value
			} else if id, ok := assignStmt.Value.(*Identifier); ok {
				focus.Icon = id.Value
			}
			
		case "cost":
			if num, ok := assignStmt.Value.(*NumberLiteral); ok {
				cost, err := strconv.Atoi(num.Value)
				if err == nil {
					focus.Cost = cost
				}
			}
			
		case "x":
			if posValue := fp.parsePositionValue(assignStmt.Value); posValue != 0 {
				focus.Position.X = posValue
			}
			
		case "y":
			if posValue := fp.parsePositionValue(assignStmt.Value); posValue != 0 {
				focus.Position.Y = posValue
			}
			
		case "relative_position_id":
			if str, ok := assignStmt.Value.(*StringLiteral); ok {
				focus.RelativePositionID = str.Value
			} else if id, ok := assignStmt.Value.(*Identifier); ok {
				focus.RelativePositionID = id.Value
			}
			
		case "prerequisite":
			if prereqBlock, ok := assignStmt.Value.(*BlockStatement); ok {
				prereq := fp.parsePrerequisite(prereqBlock)
				if len(prereq) > 0 {
					focus.Prerequisites = append(focus.Prerequisites, prereq)
				}
			}
			
		case "mutually_exclusive":
			if mexBlock, ok := assignStmt.Value.(*BlockStatement); ok {
				focus.MutuallyExclusive = fp.parseMutuallyExclusive(mexBlock)
			}
			
		case "cancel_if_invalid":
			if id, ok := assignStmt.Value.(*Identifier); ok {
				focus.CancelIfInvalid = (id.Value == "yes" || id.Value == "true")
			}
			
		case "continue_if_invalid":
			if id, ok := assignStmt.Value.(*Identifier); ok {
				focus.ContinueIfInvalid = (id.Value == "yes" || id.Value == "true")
			}
			
		case "available_if_capitulated":
			if id, ok := assignStmt.Value.(*Identifier); ok {
				focus.AvailableIfCapitulated = (id.Value == "yes" || id.Value == "true")
			}
			
		case "available":
			// Store as raw string for now
			focus.Available = fp.blockToString(assignStmt.Value)
			
		case "bypass":
			focus.Bypass = fp.blockToString(assignStmt.Value)
			
		case "completion_reward":
			focus.CompletionReward = fp.blockToString(assignStmt.Value)
			
		case "ai_will_do":
			focus.AIWillDo = fp.blockToString(assignStmt.Value)
			
		case "search_filters":
			if searchBlock, ok := assignStmt.Value.(*BlockStatement); ok {
				focus.SearchFilters = fp.parseSearchFilters(searchBlock)
			}
		}
	}
	
	return focus, nil
}

// parsePositionValue parses a position value (number or variable reference)
func (fp *FocusParser) parsePositionValue(expr Expression) int {
	var rawValue string
	
	switch v := expr.(type) {
	case *NumberLiteral:
		rawValue = v.Value
	case *Identifier:
		rawValue = v.Value
	default:
		return 0
	}
	
	// Resolve variable if needed
	resolvedValue := fp.resolveVariable(rawValue)
	val, err := strconv.Atoi(resolvedValue)
	if err != nil {
		return 0
	}
	
	return val
}

// parsePrerequisite parses a prerequisite block
func (fp *FocusParser) parsePrerequisite(block *BlockStatement) []string {
	prereqs := make([]string, 0)
	
	for _, stmt := range block.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		if assignStmt.Name.Value == "focus" {
			if str, ok := assignStmt.Value.(*StringLiteral); ok {
				prereqs = append(prereqs, str.Value)
			} else if id, ok := assignStmt.Value.(*Identifier); ok {
				prereqs = append(prereqs, id.Value)
			}
		}
	}
	
	return prereqs
}

// parseMutuallyExclusive parses a mutually_exclusive block
func (fp *FocusParser) parseMutuallyExclusive(block *BlockStatement) []string {
	mex := make([]string, 0)
	
	for _, stmt := range block.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		if assignStmt.Name.Value == "focus" {
			if str, ok := assignStmt.Value.(*StringLiteral); ok {
				mex = append(mex, str.Value)
			} else if id, ok := assignStmt.Value.(*Identifier); ok {
				mex = append(mex, id.Value)
			}
		}
	}
	
	return mex
}

// parseSearchFilters parses a search_filters block
func (fp *FocusParser) parseSearchFilters(block *BlockStatement) []string {
	filters := make([]string, 0)
	
	for _, stmt := range block.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		// Search filters are identifiers
		filters = append(filters, assignStmt.Name.Value)
	}
	
	return filters
}

// blockToString converts a block or expression to string (for raw storage)
func (fp *FocusParser) blockToString(expr Expression) string {
	// For now, just return a placeholder
	// In a full implementation, we would reconstruct the block syntax
	switch v := expr.(type) {
	case *BlockStatement:
		return "{...}" // Placeholder
	case *StringLiteral:
		return v.Value
	case *Identifier:
		return v.Value
	default:
		return ""
	}
}

// resolveVariable resolves a variable reference (@VAR) to its value
func (fp *FocusParser) resolveVariable(value string) string {
	if strings.HasPrefix(value, "@") {
		if resolved, ok := fp.variables[value]; ok {
			return resolved
		}
		// Variable not found, return original
		return value
	}
	return value
}
