package parser

import (
	"testing"
)

func TestLexer_BasicTokens(t *testing.T) {
	input := `{ } = < >`
	
	lexer := NewLexer(input)
	
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokenLeftBrace, "{"},
		{TokenRightBrace, "}"},
		{TokenEquals, "="},
		{TokenLessThan, "<"},
		{TokenGreaterThan, ">"},
		{TokenEOF, ""},
	}
	
	for i, tt := range tests {
		token := lexer.NextToken()
		
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, token.Type)
		}
		
		if token.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong token value. expected=%q, got=%q",
				i, tt.expectedValue, token.Value)
		}
	}
}

func TestLexer_Identifiers(t *testing.T) {
	input := `tech_support support_folder engineers_tech`
	
	lexer := NewLexer(input)
	
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokenIdentifier, "tech_support"},
		{TokenIdentifier, "support_folder"},
		{TokenIdentifier, "engineers_tech"},
		{TokenEOF, ""},
	}
	
	for i, tt := range tests {
		token := lexer.NextToken()
		
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, token.Type)
		}
		
		if token.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong token value. expected=%q, got=%q",
				i, tt.expectedValue, token.Value)
		}
	}
}

func TestLexer_Keywords(t *testing.T) {
	input := `focus_tree focus id icon position`
	
	lexer := NewLexer(input)
	
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokenKeyword, "focus_tree"},
		{TokenKeyword, "focus"},
		{TokenKeyword, "id"},
		{TokenKeyword, "icon"},
		{TokenKeyword, "position"},
		{TokenEOF, ""},
	}
	
	for i, tt := range tests {
		token := lexer.NextToken()
		
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, token.Type)
		}
		
		if token.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong token value. expected=%q, got=%q",
				i, tt.expectedValue, token.Value)
		}
	}
}

func TestLexer_Numbers(t *testing.T) {
	input := `1.0 1918 -5 45.67 @1918 @SUPP`
	
	lexer := NewLexer(input)
	
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokenNumber, "1.0"},
		{TokenNumber, "1918"},
		{TokenNumber, "-5"},
		{TokenNumber, "45.67"},
		{TokenNumber, "@1918"},
		{TokenIdentifier, "@SUPP"}, // @SUPP is identifier (letters after @)
		{TokenEOF, ""},
	}
	
	for i, tt := range tests {
		token := lexer.NextToken()
		
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, token.Type)
		}
		
		if token.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong token value. expected=%q, got=%q",
				i, tt.expectedValue, token.Value)
		}
	}
}

func TestLexer_Strings(t *testing.T) {
	input := `"hello world" "GFX_focus_icon"`
	
	lexer := NewLexer(input)
	
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokenString, "hello world"},
		{TokenString, "GFX_focus_icon"},
		{TokenEOF, ""},
	}
	
	for i, tt := range tests {
		token := lexer.NextToken()
		
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, token.Type)
		}
		
		if token.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong token value. expected=%q, got=%q",
				i, tt.expectedValue, token.Value)
		}
	}
}

func TestLexer_Comments(t *testing.T) {
	input := `tech_support # this is a comment
research_cost = 1.0`
	
	lexer := NewLexer(input)
	
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokenIdentifier, "tech_support"},
		{TokenKeyword, "research_cost"},
		{TokenEquals, "="},
		{TokenNumber, "1.0"},
		{TokenEOF, ""},
	}
	
	for i, tt := range tests {
		token := lexer.NextToken()
		
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, token.Type)
		}
		
		if token.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong token value. expected=%q, got=%q",
				i, tt.expectedValue, token.Value)
		}
	}
}

func TestLexer_ComplexStructure(t *testing.T) {
	input := `technologies = {
	tech_support = {
		research_cost = 1.0
		position = { x = @SUPP y = @1918 }
	}
}`
	
	lexer := NewLexer(input)
	
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokenKeyword, "technologies"},
		{TokenEquals, "="},
		{TokenLeftBrace, "{"},
		{TokenIdentifier, "tech_support"},
		{TokenEquals, "="},
		{TokenLeftBrace, "{"},
		{TokenKeyword, "research_cost"},
		{TokenEquals, "="},
		{TokenNumber, "1.0"},
		{TokenKeyword, "position"},
		{TokenEquals, "="},
		{TokenLeftBrace, "{"},
		{TokenKeyword, "x"},
		{TokenEquals, "="},
		{TokenIdentifier, "@SUPP"}, // @SUPP is identifier
		{TokenKeyword, "y"},
		{TokenEquals, "="},
		{TokenNumber, "@1918"}, // @1918 is number
		{TokenRightBrace, "}"},
		{TokenRightBrace, "}"},
		{TokenRightBrace, "}"},
		{TokenEOF, ""},
	}
	
	for i, tt := range tests {
		token := lexer.NextToken()
		
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, token.Type)
		}
		
		if token.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong token value. expected=%q, got=%q",
				i, tt.expectedValue, token.Value)
		}
	}
}

func TestLexer_Date(t *testing.T) {
	input := `1939.1.1 1945.5.9`
	
	lexer := NewLexer(input)
	
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokenDate, "1939.1.1"},
		{TokenDate, "1945.5.9"},
		{TokenEOF, ""},
	}
	
	for i, tt := range tests {
		token := lexer.NextToken()
		
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, token.Type)
		}
		
		if token.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong token value. expected=%q, got=%q",
				i, tt.expectedValue, token.Value)
		}
	}
}
