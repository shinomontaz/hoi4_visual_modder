package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shinomontaz/hoi4_visual_modder/internal/domain"
)

// TechParser converts AST to domain.Technology models
type TechParser struct {
	variables map[string]string // Variable definitions (@1918 = 0)
}

// NewTechParser creates a new TechParser
func NewTechParser() *TechParser {
	return &TechParser{
		variables: make(map[string]string),
	}
}

// ParseTechnologies parses a technologies block from AST
func (tp *TechParser) ParseTechnologies(program *Program) ([]*domain.Technology, error) {
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
			tp.handleVariableDefinition(assignStmt)
		}
	}
	
	technologies := make([]*domain.Technology, 0)
	
	// Second pass: find and parse the technologies block
	for _, stmt := range program.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		// Check if it's the technologies block
		if assignStmt.Name.Value == "technologies" {
			block, ok := assignStmt.Value.(*BlockStatement)
			if !ok {
				return nil, fmt.Errorf("technologies value is not a block")
			}
			
			// Collect variables inside technologies block first
			for _, techStmt := range block.Statements {
				techAssign, ok := techStmt.(*AssignmentStatement)
				if !ok {
					continue
				}
				
				if strings.HasPrefix(techAssign.Name.Value, "@") {
					tp.handleVariableDefinition(techAssign)
				}
			}
			
			// Now parse each technology in the block
			for _, techStmt := range block.Statements {
				techAssign, ok := techStmt.(*AssignmentStatement)
				if !ok {
					continue
				}
				
				// Skip variable definitions
				if strings.HasPrefix(techAssign.Name.Value, "@") {
					continue
				}
				
				// Parse technology
				techBlock, ok := techAssign.Value.(*BlockStatement)
				if !ok {
					continue // Not a tech definition, skip
				}
				
				tech, err := tp.parseTechnology(techAssign.Name.Value, techBlock)
				if err != nil {
					return nil, fmt.Errorf("failed to parse technology %s: %w", techAssign.Name.Value, err)
				}
				
				technologies = append(technologies, tech)
			}
			
			break
		}
	}
	
	return technologies, nil
}

// handleVariableDefinition stores a variable definition
func (tp *TechParser) handleVariableDefinition(stmt *AssignmentStatement) {
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
	
	tp.variables[varName] = value
}

// parseTechnology parses a single technology block
func (tp *TechParser) parseTechnology(id string, block *BlockStatement) (*domain.Technology, error) {
	tech := &domain.Technology{
		ID:                id,
		Categories:        make([]string, 0),
		Effects:           make(map[string]map[string]float64),
		Paths:             make([]domain.TechPath, 0),
		ResearchCost:      1.0,
		AIResearchWeights: make(map[string]float64),
		XOR:               make([]string, 0),
	}
	
	// Parse each field in the technology block
	for _, stmt := range block.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		fieldName := assignStmt.Name.Value
		
		switch fieldName {
		case "research_cost":
			if num, ok := assignStmt.Value.(*NumberLiteral); ok {
				cost, err := strconv.ParseFloat(num.Value, 64)
				if err == nil {
					tech.ResearchCost = cost
				}
			}
			
		case "start_year":
			// Store as metadata, not used in domain model currently
			
		case "folder":
			if folderBlock, ok := assignStmt.Value.(*BlockStatement); ok {
				tp.parseFolder(tech, folderBlock)
			}
			
		case "categories":
			if catBlock, ok := assignStmt.Value.(*BlockStatement); ok {
				tech.Categories = tp.parseCategories(catBlock)
			}
			
		case "path":
			if pathBlock, ok := assignStmt.Value.(*BlockStatement); ok {
				path := tp.parsePath(pathBlock)
				if path != nil {
					tech.Paths = append(tech.Paths, *path)
				}
			}
			
		case "xor":
			if xorBlock, ok := assignStmt.Value.(*BlockStatement); ok {
				tech.XOR = tp.parseXOR(xorBlock)
			}
			
		case "xp_research_type":
			if str, ok := assignStmt.Value.(*StringLiteral); ok {
				tech.XPResearchType = str.Value
			} else if id, ok := assignStmt.Value.(*Identifier); ok {
				tech.XPResearchType = id.Value
			}
			
		case "xp_boost_cost":
			if num, ok := assignStmt.Value.(*NumberLiteral); ok {
				cost, err := strconv.Atoi(num.Value)
				if err == nil {
					tech.XPBoostCost = cost
				}
			}
			
		case "xp_research_bonus":
			if num, ok := assignStmt.Value.(*NumberLiteral); ok {
				bonus, err := strconv.ParseFloat(num.Value, 64)
				if err == nil {
					tech.XPResearchBonus = bonus
				}
			}
		}
	}
	
	return tech, nil
}

// parseFolder parses a folder block
func (tp *TechParser) parseFolder(tech *domain.Technology, block *BlockStatement) {
	for _, stmt := range block.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		switch assignStmt.Name.Value {
		case "name":
			if str, ok := assignStmt.Value.(*StringLiteral); ok {
				tech.Folder = str.Value
			} else if id, ok := assignStmt.Value.(*Identifier); ok {
				tech.Folder = id.Value
			}
			
		case "position":
			if posBlock, ok := assignStmt.Value.(*BlockStatement); ok {
				tech.Position = tp.parsePosition(posBlock)
			}
		}
	}
}

// parsePosition parses a position block { x = 5 y = 10 }
func (tp *TechParser) parsePosition(block *BlockStatement) domain.Position {
	pos := domain.Position{}
	
	for _, stmt := range block.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		var value int
		var rawValue string
		
		// Get the raw value (could be NumberLiteral or Identifier)
		switch v := assignStmt.Value.(type) {
		case *NumberLiteral:
			rawValue = v.Value
		case *Identifier:
			rawValue = v.Value
		default:
			continue
		}
		
		// Resolve variable if needed
		resolvedValue := tp.resolveVariable(rawValue)
		val, err := strconv.Atoi(resolvedValue)
		if err == nil {
			value = val
		}
		
		switch assignStmt.Name.Value {
		case "x":
			pos.X = value
		case "y":
			pos.Y = value
		}
	}
	
	return pos
}

// parseCategories parses a categories block
func (tp *TechParser) parseCategories(block *BlockStatement) []string {
	categories := make([]string, 0)
	seen := make(map[string]bool)
	
	// In HOI4 files, categories are listed as: categories = { cat1 cat2 cat3 }
	// Our parser sees them as assignments (cat1 = cat2, cat2 = cat3, etc.)
	// So we extract both the name and value
	for _, stmt := range block.Statements {
		if assignStmt, ok := stmt.(*AssignmentStatement); ok {
			// Add the assignment name (left side)
			if !seen[assignStmt.Name.Value] {
				categories = append(categories, assignStmt.Name.Value)
				seen[assignStmt.Name.Value] = true
			}
			
			// Add the value if it's an identifier (right side)
			if id, ok := assignStmt.Value.(*Identifier); ok {
				if !seen[id.Value] {
					categories = append(categories, id.Value)
					seen[id.Value] = true
				}
			}
		}
	}
	
	return categories
}

// parsePath parses a path block
func (tp *TechParser) parsePath(block *BlockStatement) *domain.TechPath {
	path := &domain.TechPath{
		ResearchCostCoeff: 1.0,
	}
	
	for _, stmt := range block.Statements {
		assignStmt, ok := stmt.(*AssignmentStatement)
		if !ok {
			continue
		}
		
		switch assignStmt.Name.Value {
		case "leads_to_tech":
			if id, ok := assignStmt.Value.(*Identifier); ok {
				path.LeadsToTech = id.Value
			} else if str, ok := assignStmt.Value.(*StringLiteral); ok {
				path.LeadsToTech = str.Value
			}
			
		case "research_cost_coeff":
			if num, ok := assignStmt.Value.(*NumberLiteral); ok {
				coeff, err := strconv.ParseFloat(num.Value, 64)
				if err == nil {
					path.ResearchCostCoeff = coeff
				}
			}
		}
	}
	
	return path
}

// parseXOR parses an XOR block (mutually exclusive technologies)
func (tp *TechParser) parseXOR(block *BlockStatement) []string {
	xor := make([]string, 0)
	
	for _, stmt := range block.Statements {
		if assignStmt, ok := stmt.(*AssignmentStatement); ok {
			xor = append(xor, assignStmt.Name.Value)
		}
	}
	
	return xor
}

// resolveVariable resolves a variable reference (@VAR) to its value
func (tp *TechParser) resolveVariable(value string) string {
	if strings.HasPrefix(value, "@") {
		if resolved, ok := tp.variables[value]; ok {
			return resolved
		}
		// Variable not found, return original
		return value
	}
	return value
}
