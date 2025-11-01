package parser

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
	// TODO: Implement tokenization logic
	// This is a placeholder for the lexer implementation
	return Token{Type: TokenEOF, Line: l.line, Column: l.column}
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
