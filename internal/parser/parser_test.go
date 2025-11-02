package parser

import (
	"testing"
)

func TestParser_SimpleAssignment(t *testing.T) {
	input := `research_cost = 1.0`
	
	parser := NewParser(input)
	program, err := parser.Parse()
	
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *AssignmentStatement. got=%T", program.Statements[0])
	}
	
	if stmt.Name.Value != "research_cost" {
		t.Errorf("stmt.Name.Value not 'research_cost'. got=%s", stmt.Name.Value)
	}
	
	numLit, ok := stmt.Value.(*NumberLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not *NumberLiteral. got=%T", stmt.Value)
	}
	
	if numLit.Value != "1.0" {
		t.Errorf("numLit.Value not '1.0'. got=%s", numLit.Value)
	}
}

func TestParser_StringAssignment(t *testing.T) {
	input := `icon = "GFX_focus_icon"`
	
	parser := NewParser(input)
	program, err := parser.Parse()
	
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *AssignmentStatement. got=%T", program.Statements[0])
	}
	
	strLit, ok := stmt.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not *StringLiteral. got=%T", stmt.Value)
	}
	
	if strLit.Value != "GFX_focus_icon" {
		t.Errorf("strLit.Value not 'GFX_focus_icon'. got=%s", strLit.Value)
	}
}

func TestParser_BlockAssignment(t *testing.T) {
	input := `position = { x = 5 y = 10 }`
	
	parser := NewParser(input)
	program, err := parser.Parse()
	
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *AssignmentStatement. got=%T", program.Statements[0])
	}
	
	block, ok := stmt.Value.(*BlockStatement)
	if !ok {
		t.Fatalf("stmt.Value is not *BlockStatement. got=%T", stmt.Value)
	}
	
	if len(block.Statements) != 2 {
		t.Fatalf("block.Statements does not contain 2 statements. got=%d", len(block.Statements))
	}
	
	// Check x = 5
	xStmt, ok := block.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("block.Statements[0] is not *AssignmentStatement. got=%T", block.Statements[0])
	}
	
	if xStmt.Name.Value != "x" {
		t.Errorf("xStmt.Name.Value not 'x'. got=%s", xStmt.Name.Value)
	}
	
	xNum, ok := xStmt.Value.(*NumberLiteral)
	if !ok {
		t.Fatalf("xStmt.Value is not *NumberLiteral. got=%T", xStmt.Value)
	}
	
	if xNum.Value != "5" {
		t.Errorf("xNum.Value not '5'. got=%s", xNum.Value)
	}
}

func TestParser_NestedBlocks(t *testing.T) {
	input := `tech_support = {
		folder = {
			name = support_folder
			position = { x = 5 y = 10 }
		}
	}`
	
	parser := NewParser(input)
	program, err := parser.Parse()
	
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *AssignmentStatement. got=%T", program.Statements[0])
	}
	
	if stmt.Name.Value != "tech_support" {
		t.Errorf("stmt.Name.Value not 'tech_support'. got=%s", stmt.Name.Value)
	}
	
	outerBlock, ok := stmt.Value.(*BlockStatement)
	if !ok {
		t.Fatalf("stmt.Value is not *BlockStatement. got=%T", stmt.Value)
	}
	
	if len(outerBlock.Statements) != 1 {
		t.Fatalf("outerBlock.Statements does not contain 1 statement. got=%d", len(outerBlock.Statements))
	}
	
	folderStmt, ok := outerBlock.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("outerBlock.Statements[0] is not *AssignmentStatement. got=%T", outerBlock.Statements[0])
	}
	
	if folderStmt.Name.Value != "folder" {
		t.Errorf("folderStmt.Name.Value not 'folder'. got=%s", folderStmt.Name.Value)
	}
	
	innerBlock, ok := folderStmt.Value.(*BlockStatement)
	if !ok {
		t.Fatalf("folderStmt.Value is not *BlockStatement. got=%T", folderStmt.Value)
	}
	
	if len(innerBlock.Statements) != 2 {
		t.Fatalf("innerBlock.Statements does not contain 2 statements. got=%d", len(innerBlock.Statements))
	}
}

func TestParser_VariableReferences(t *testing.T) {
	input := `position = { x = @SUPP y = @1918 }`
	
	parser := NewParser(input)
	program, err := parser.Parse()
	
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	
	stmt := program.Statements[0].(*AssignmentStatement)
	block := stmt.Value.(*BlockStatement)
	
	xStmt := block.Statements[0].(*AssignmentStatement)
	xId := xStmt.Value.(*Identifier) // @SUPP is identifier
	
	if xId.Value != "@SUPP" {
		t.Errorf("xId.Value not '@SUPP'. got=%s", xId.Value)
	}
	
	yStmt := block.Statements[1].(*AssignmentStatement)
	yNum := yStmt.Value.(*NumberLiteral) // @1918 is number
	
	if yNum.Value != "@1918" {
		t.Errorf("yNum.Value not '@1918'. got=%s", yNum.Value)
	}
}

func TestParser_ComplexTechnology(t *testing.T) {
	input := `technologies = {
		tech_support = {
			research_cost = 1.0
			start_year = 1918
			folder = {
				name = support_folder
				position = { x = @SUPP y = @1918 }
			}
		}
	}`
	
	parser := NewParser(input)
	program, err := parser.Parse()
	
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	// Verify top-level structure
	techStmt := program.Statements[0].(*AssignmentStatement)
	if techStmt.Name.Value != "technologies" {
		t.Errorf("techStmt.Name.Value not 'technologies'. got=%s", techStmt.Name.Value)
	}
	
	techBlock := techStmt.Value.(*BlockStatement)
	if len(techBlock.Statements) != 1 {
		t.Fatalf("techBlock.Statements does not contain 1 statement. got=%d", len(techBlock.Statements))
	}
	
	// Verify tech_support block
	supportStmt := techBlock.Statements[0].(*AssignmentStatement)
	if supportStmt.Name.Value != "tech_support" {
		t.Errorf("supportStmt.Name.Value not 'tech_support'. got=%s", supportStmt.Name.Value)
	}
	
	supportBlock := supportStmt.Value.(*BlockStatement)
	if len(supportBlock.Statements) != 3 {
		t.Fatalf("supportBlock.Statements does not contain 3 statements. got=%d", len(supportBlock.Statements))
	}
}
