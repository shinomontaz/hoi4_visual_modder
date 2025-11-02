package parser

// Node represents a node in the AST
type Node interface {
	TokenLiteral() string
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// AssignmentStatement represents an assignment (key = value)
type AssignmentStatement struct {
	Token Token // The '=' token
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Value }

// BlockStatement represents a block { ... }
type BlockStatement struct {
	Token      Token // The '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) expressionNode()      {} // Block can be used as expression
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Value }

// Identifier represents an identifier
type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Value }

// StringLiteral represents a string literal
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Value }

// NumberLiteral represents a number literal
type NumberLiteral struct {
	Token Token
	Value string
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Value }

// DateLiteral represents a date literal (1939.1.1)
type DateLiteral struct {
	Token Token
	Value string
}

func (dl *DateLiteral) expressionNode()      {}
func (dl *DateLiteral) TokenLiteral() string { return dl.Token.Value }

// ArrayLiteral represents an array of values
type ArrayLiteral struct {
	Token    Token // The first token of the array
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Value }

// ObjectLiteral represents an object/block with key-value pairs
type ObjectLiteral struct {
	Token Token // The '{' token
	Pairs map[string]Expression
}

func (ol *ObjectLiteral) expressionNode()      {}
func (ol *ObjectLiteral) TokenLiteral() string { return ol.Token.Value }
