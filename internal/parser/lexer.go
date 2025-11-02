package parser

import "strings"

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Column  int
}

// TokenType represents the type of token
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenError
	
	// Literals
	TokenIdentifier // variable names, keywords
	TokenString     // "quoted string"
	TokenNumber     // 123, 45.67
	TokenDate       // 1939.1.1
	
	// Delimiters
	TokenLeftBrace  // {
	TokenRightBrace // }
	TokenEquals     // =
	TokenLessThan   // <
	TokenGreaterThan // >
	
	// Keywords (will be identified during lexing)
	TokenKeyword
)

// String returns string representation of TokenType
func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return "ERROR"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenString:
		return "STRING"
	case TokenNumber:
		return "NUMBER"
	case TokenDate:
		return "DATE"
	case TokenLeftBrace:
		return "{"
	case TokenRightBrace:
		return "}"
	case TokenEquals:
		return "="
	case TokenLessThan:
		return "<"
	case TokenGreaterThan:
		return ">"
	case TokenKeyword:
		return "KEYWORD"
	default:
		return "UNKNOWN"
	}
}

// Lexer tokenizes Paradox scripting language
type Lexer struct {
	input   string
	pos     int
	line    int
	column  int
	current rune
}

// NewLexer creates a new Lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	// Skip whitespace and comments
	l.skipWhitespace()
	
	// Skip comments
	if l.current == '#' {
		l.skipComment()
		l.skipWhitespace()
	}
	
	token := Token{
		Line:   l.line,
		Column: l.column,
	}
	
	// EOF
	if l.current == 0 {
		token.Type = TokenEOF
		return token
	}
	
	// Single character tokens
	switch l.current {
	case '{':
		token.Type = TokenLeftBrace
		token.Value = "{"
		l.readChar()
		return token
	case '}':
		token.Type = TokenRightBrace
		token.Value = "}"
		l.readChar()
		return token
	case '=':
		token.Type = TokenEquals
		token.Value = "="
		l.readChar()
		return token
	case '<':
		token.Type = TokenLessThan
		token.Value = "<"
		l.readChar()
		return token
	case '>':
		token.Type = TokenGreaterThan
		token.Value = ">"
		l.readChar()
		return token
	case '"':
		// String literal
		return l.readString()
	}
	
	// Variable references or numbers starting with @
	if l.current == '@' {
		// Check if it's followed by a letter (variable name) or digit (variable reference in expression)
		next := l.peekChar()
		if isLetter(next) || isDigit(next) || next == '_' {
			// Could be @VAR or @1918
			// If next is letter or underscore, treat as identifier
			// If next is digit, treat as number
			if isLetter(next) || next == '_' {
				return l.readIdentifier() // @SUPP, @VAR
			} else {
				return l.readNumber() // @1918
			}
		}
	}
	
	// Numbers (including dates and negative numbers)
	if isDigit(l.current) || (l.current == '-' && isDigit(l.peekChar())) {
		return l.readNumber()
	}
	
	// Identifiers and keywords
	if isLetter(l.current) || l.current == '_' {
		return l.readIdentifier()
	}
	
	// Unknown character
	token.Type = TokenError
	token.Value = string(l.current)
	l.readChar()
	return token
}

// readChar reads the next character
func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.current = 0 // EOF
	} else {
		l.current = rune(l.input[l.pos])
	}
	l.pos++
	l.column++
	
	if l.current == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar looks at the next character without consuming it
func (l *Lexer) peekChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r' {
		l.readChar()
	}
}

// skipComment skips comment lines starting with #
func (l *Lexer) skipComment() {
	if l.current == '#' {
		for l.current != '\n' && l.current != 0 {
			l.readChar()
		}
	}
}

// readString reads a quoted string
func (l *Lexer) readString() Token {
	token := Token{
		Type:   TokenString,
		Line:   l.line,
		Column: l.column,
	}
	
	l.readChar() // skip opening quote
	
	var value string
	for l.current != '"' && l.current != 0 {
		if l.current == '\\' {
			// Handle escape sequences
			l.readChar()
			if l.current != 0 {
				value += string(l.current)
				l.readChar()
			}
		} else {
			value += string(l.current)
			l.readChar()
		}
	}
	
	if l.current == '"' {
		l.readChar() // skip closing quote
	}
	
	token.Value = value
	return token
}

// readNumber reads a number (integer, float, date, or variable reference)
func (l *Lexer) readNumber() Token {
	token := Token{
		Line:   l.line,
		Column: l.column,
	}
	
	var value string
	
	// Handle variable references like @1918 or @SUPP
	if l.current == '@' {
		value += string(l.current)
		l.readChar()
		
		// Read identifier after @
		for isLetter(l.current) || isDigit(l.current) || l.current == '_' {
			value += string(l.current)
			l.readChar()
		}
		
		token.Type = TokenNumber
		token.Value = value
		return token
	}
	
	// Handle negative numbers
	if l.current == '-' {
		value += string(l.current)
		l.readChar()
	}
	
	// Read digits
	for isDigit(l.current) {
		value += string(l.current)
		l.readChar()
	}
	
	// Check for decimal point (float) or date separator
	if l.current == '.' {
		value += string(l.current)
		l.readChar()
		
		// Read more digits
		for isDigit(l.current) {
			value += string(l.current)
			l.readChar()
		}
		
		// Check if it's a date (has another dot)
		if l.current == '.' {
			value += string(l.current)
			l.readChar()
			
			// Read final part of date
			for isDigit(l.current) {
				value += string(l.current)
				l.readChar()
			}
			
			token.Type = TokenDate
			token.Value = value
			return token
		}
	}
	
	// It's a regular number
	token.Type = TokenNumber
	token.Value = value
	return token
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() Token {
	token := Token{
		Type:   TokenIdentifier,
		Line:   l.line,
		Column: l.column,
	}
	
	var value string
	
	// Allow @ at the beginning for variable names
	if l.current == '@' {
		value += string(l.current)
		l.readChar()
	}
	
	for isLetter(l.current) || isDigit(l.current) || l.current == '_' {
		value += string(l.current)
		l.readChar()
	}
	
	token.Value = value
	
	// Check if it's a keyword (but not if it starts with @)
	if !strings.HasPrefix(value, "@") && isKeyword(value) {
		token.Type = TokenKeyword
	}
	
	return token
}

// isDigit checks if a rune is a digit
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isLetter checks if a rune is a letter
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isKeyword checks if a string is a Paradox keyword
func isKeyword(s string) bool {
	keywords := map[string]bool{
		"focus_tree":           true,
		"focus":                true,
		"id":                   true,
		"icon":                 true,
		"x":                    true,
		"y":                    true,
		"cost":                 true,
		"prerequisite":         true,
		"mutually_exclusive":   true,
		"relative_position_id": true,
		"available":            true,
		"bypass":               true,
		"cancel_if_invalid":    true,
		"continue_if_invalid":  true,
		"completion_reward":    true,
		"technologies":         true,
		"research_cost":        true,
		"start_year":           true,
		"folder":               true,
		"position":             true,
		"categories":           true,
		"path":                 true,
		"leads_to_tech":        true,
		"name":                 true,
	}
	return keywords[s]
}
