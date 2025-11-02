package main

import (
	"fmt"
	"os"

	"github.com/shinomontaz/hoi4_visual_modder/internal/parser"
)

func main() {
	content, err := os.ReadFile("common/technologies/test_tech.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	p := parser.NewParser(string(content))
	program, err := p.Parse()
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	techParser := parser.NewTechParser()
	
	// Debug: print variables
	fmt.Println("=== Parsing Technologies ===")
	
	technologies, err := techParser.ParseTechnologies(program)
	if err != nil {
		fmt.Printf("TechParser error: %v\n", err)
		return
	}

	fmt.Printf("\nParsed %d technologies:\n", len(technologies))
	for i, tech := range technologies {
		fmt.Printf("\n%d. %s\n", i+1, tech.ID)
		fmt.Printf("   Research Cost: %.1f\n", tech.ResearchCost)
		fmt.Printf("   Folder: %s\n", tech.Folder)
		fmt.Printf("   Position: (%d, %d)\n", tech.Position.X, tech.Position.Y)
		fmt.Printf("   Categories: %v\n", tech.Categories)
	}
}
