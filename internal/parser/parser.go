package parser

import (
	"fmt"
)

// Parser builds an AST from tokens
type Parser struct {
	lexer   *Lexer
	current Token
	peek    Token
	errors  []string
}

// NewParser creates a new Parser
func NewParser(input string) *Parser {
	p := &Parser{
		lexer: NewLexer(input),
	}

	// Read two tokens to initialize current and peek
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

// Parse parses the input and returns an AST
func (p *Parser) Parse() (*Program, error) {
	program := &Program{
		Statements: []Statement{},
	}

	for p.current.Type != TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parsing errors: %v", p.errors)
	}

	return program, nil
}

// parseStatement parses a statement
func (p *Parser) parseStatement() Statement {
	// Check if it's an assignment (identifier = value OR @1918 = value OR "string" = value)
	if p.current.Type == TokenIdentifier || p.current.Type == TokenKeyword || p.current.Type == TokenNumber || p.current.Type == TokenString {
		if p.peek.Type == TokenEquals {
			return p.parseAssignmentStatement()
		}
	}

	return nil
}

// parseAssignmentStatement parses an assignment statement
func (p *Parser) parseAssignmentStatement() *AssignmentStatement {
	// Name can be Identifier, Keyword, or Number (for @1918 = 0)
	stmt := &AssignmentStatement{
		Name: &Identifier{
			Token: p.current,
			Value: p.current.Value,
		},
	}

	// Expect '='
	if !p.expectPeek(TokenEquals) {
		return nil
	}

	stmt.Token = p.current

	// Parse the value
	p.nextToken()
	stmt.Value = p.parseExpression()

	return stmt
}

// parseExpression parses an expression
func (p *Parser) parseExpression() Expression {
	switch p.current.Type {
	case TokenString:
		return &StringLiteral{
			Token: p.current,
			Value: p.current.Value,
		}
	case TokenNumber:
		return &NumberLiteral{
			Token: p.current,
			Value: p.current.Value,
		}
	case TokenDate:
		return &DateLiteral{
			Token: p.current,
			Value: p.current.Value,
		}
	case TokenIdentifier, TokenKeyword:
		return &Identifier{
			Token: p.current,
			Value: p.current.Value,
		}
	case TokenLeftBrace:
		return p.parseBlockExpression()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token: %s at line %d", p.current.Type, p.current.Line))
		return nil
	}
}

// parseBlockExpression parses a block expression { ... }
func (p *Parser) parseBlockExpression() Expression {
	block := &BlockStatement{
		Token:      p.current,
		Statements: []Statement{},
	}

	p.nextToken()

	// Parse statements until we hit '}'
	for p.current.Type != TokenRightBrace && p.current.Type != TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	if p.current.Type != TokenRightBrace {
		p.errors = append(p.errors, fmt.Sprintf("expected '}', got %s at line %d", p.current.Type, p.current.Line))
		return nil
	}

	return block
}

// expectPeek checks if the next token is of the expected type
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peek.Type == t {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("expected %s, got %s at line %d", t, p.peek.Type, p.peek.Line))
	return false
}

// Errors returns the parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}
