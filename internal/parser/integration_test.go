package parser

import (
	"os"
	"testing"
)

func TestParser_RealTechnologyFile(t *testing.T) {
	// Read the test technology file
	content, err := os.ReadFile("../../common/technologies/test_tech.txt")
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}
	
	parser := NewParser(string(content))
	program, err := parser.Parse()
	
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	
	if len(program.Statements) == 0 {
		t.Fatal("program.Statements is empty")
	}
	
	// Should have at least the technologies block
	techStmt, ok := program.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("First statement is not *AssignmentStatement. got=%T", program.Statements[0])
	}
	
	if techStmt.Name.Value != "technologies" {
		t.Errorf("Expected 'technologies' block, got=%s", techStmt.Name.Value)
	}
	
	techBlock, ok := techStmt.Value.(*BlockStatement)
	if !ok {
		t.Fatalf("technologies value is not *BlockStatement. got=%T", techStmt.Value)
	}
	
	// Should have variable definitions and tech blocks
	if len(techBlock.Statements) == 0 {
		t.Fatal("technologies block is empty")
	}
	
	t.Logf("Successfully parsed %d statements in technologies block", len(techBlock.Statements))
	
	// Count tech definitions (should have tech_support, tech_engineers, tech_combat_engineers)
	techCount := 0
	for _, stmt := range techBlock.Statements {
		if assignStmt, ok := stmt.(*AssignmentStatement); ok {
			if _, ok := assignStmt.Value.(*BlockStatement); ok {
				// This is a tech definition (has a block value)
				techCount++
				t.Logf("Found tech: %s", assignStmt.Name.Value)
			}
		}
	}
	
	if techCount < 3 {
		t.Errorf("Expected at least 3 tech definitions, got=%d", techCount)
	}
}

func TestParser_VariableDefinitions(t *testing.T) {
	input := `@1918 = 0
@1934 = 2
@SUPP = 6`
	
	parser := NewParser(input)
	program, err := parser.Parse()
	
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	
	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got=%d", len(program.Statements))
	}
	
	// Check first variable
	stmt1 := program.Statements[0].(*AssignmentStatement)
	if stmt1.Name.Value != "@1918" {
		t.Errorf("Expected '@1918', got=%s", stmt1.Name.Value)
	}
	
	num1 := stmt1.Value.(*NumberLiteral)
	if num1.Value != "0" {
		t.Errorf("Expected '0', got=%s", num1.Value)
	}
}
